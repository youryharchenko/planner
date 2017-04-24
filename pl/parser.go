package pl

import (
	"fmt"
	"go/token"
	"log"
	"strconv"
	"strings"
)

type Node interface {
	Type() NodeType
	// Position() Pos
	String() string
	Value(*Vars) Node
	Copy() Node
}

type Pos int

func (p Pos) Position() Pos {
	return p
}

type NodeType int

func (t NodeType) Type() NodeType {
	return t
}

const (
	NodeIdent NodeType = iota
	NodeString
	NodeNumber
	NodeCall
	NodeVector
	NodeList
	NodeRef
	NodeFunc
	NodePair
	NodeObj
)

type IdentNode struct {
	// Pos
	NodeType
	Ident string
}

func (node IdentNode) Copy() Node {
	return newIdentNode(node.Ident)
}

func (node IdentNode) String() string {
	if node.Ident == "nil" {
		return "()"
	}

	return node.Ident
}

func (node IdentNode) Value(v *Vars) Node {
	return node
}

type StringNode struct {
	// Pos
	NodeType
	Val string
}

func (node StringNode) Copy() Node {
	return newStringNode(node.Val)
}

func (node StringNode) String() string {
	//return "\"" + node.Val + "\""
	return node.Val
}

func (node StringNode) Value(v *Vars) Node {
	return node
}

type NumberNode struct {
	// Pos
	NodeType
	Val        string
	Int        int64
	Float      float64
	NumberType token.Token
}

func (node NumberNode) Copy() Node {
	return &NumberNode{NodeType: node.Type(), Val: node.Val, NumberType: node.NumberType}
}

func (node NumberNode) String() string {
	return node.Val
}

func (node NumberNode) Value(v *Vars) Node {
	return node
}

type VectorNode struct {
	// Pos
	NodeType
	Nodes []Node
}

func (node VectorNode) Copy() Node {
	vect := VectorNode{NodeType: node.Type(), Nodes: make([]Node, len(node.Nodes))}
	for i, v := range node.Nodes {
		vect.Nodes[i] = v.Copy()
	}
	return vect
}

func (node VectorNode) String() string {

	s := "["
	b := ""
	for _, n := range node.Nodes {
		s += b + n.String()
		b = " "
	}
	return s + "]"
}

func (node VectorNode) Value(v *Vars) Node {
	vect := VectorNode{NodeType: node.Type(), Nodes: make([]Node, len(node.Nodes))}
	for i, val := range node.Nodes {
		vect.Nodes[i] = val.Value(v)
	}
	return vect
}

type ListNode struct {
	// Pos
	NodeType
	Head *PairNode
	//Nodes []Node
}

func (node ListNode) Copy() Node {
	//vect := ListNode{NodeType: node.Type(), Nodes: make([]Node, len(node.Nodes))}
	//for i, v := range node.Nodes {
	//	vect.Nodes[i] = v.Copy()
	//}
	//return vect
	var pair *PairNode
	list := ListNode{NodeType: node.Type(), Head: nil}
	for pair = node.Head; pair != nil; pair = pair.Second {
		npair := newPairNode(pair.First, list.Head)
		list.Head = &npair
	}
	return list
}

func (node ListNode) Rev() ListNode {
	return node.Copy().(ListNode)
}

func (node ListNode) String() string {
	//s := fmt.Sprint(node.Nodes)
	//return "(" + s[1:len(s)-1] + ")"
	//rev := node.Rev()
	rev := node
	var pair *PairNode
	s := "("
	b := ""
	for pair = rev.Head; pair != nil; pair = pair.Second {
		s += b + pair.First.String()
		b = " "
	}
	return s + ")"
}

func (node ListNode) Value(v *Vars) Node {
	//vect := ListNode{NodeType: node.Type(), Nodes: make([]Node, len(node.Nodes))}
	//for i, val := range node.Nodes {
	//	vect.Nodes[i] = val.Value(v)
	//}
	var pair *PairNode
	list := ListNode{NodeType: node.Type(), Head: nil}
	for pair = node.Head; pair != nil; pair = pair.Second {
		npair := newPairNode(pair.First.Value(v), list.Head)
		list.Head = &npair
	}
	return list
}

func (node ListNode) Nodes(n int64) Node {

	i := int64(0)
	for pair := node.Head; pair != nil; pair = pair.Second {
		if i == n {
			return pair.First
		}
		i++
	}
	log.Panicf("ListNodes>>Nodes: out of range: %d", i)
	return node
}

func (node ListNode) Tail(n int64) ListNode {

	i := int64(0)
	for pair := node.Head; pair != nil; pair = pair.Second {
		//log.Println(i, n)
		if i == n {
			return ListNode{NodeType: node.NodeType, Head: pair}
		}
		if pair.Second == nil {
			return newListNode()
		}
		i++
	}
	log.Panicf("ListNodes>>Tail: out of range: %d", i)
	return node
}

func (node ListNode) Len() int64 {

	i := int64(0)
	for pair := node.Head; pair != nil; pair = pair.Second {
		i++
	}

	return i
}

func (node ListNode) Cons(n Node) ListNode {
	pair := newPairNode(n, node.Head)
	node.Head = &pair
	return node
}

type CallNode struct {
	// Pos
	NodeType
	Callee Node
	Args   []Node
}

func (node CallNode) Copy() Node {
	call := CallNode{NodeType: node.Type(), Callee: node.Callee.Copy(), Args: make([]Node, len(node.Args))}
	for i, v := range node.Args {
		call.Args[i] = v.Copy()
	}
	return call
}

func (node CallNode) String() string {
	args := fmt.Sprint(node.Args)
	return fmt.Sprintf("{%s %s}", node.Callee, args[1:len(args)-1])
}

func (node CallNode) Value(v *Vars) Node {
	//log.Println("CallNode.Node", node.String())
	fn := node.Callee.Value(v)
	//log.Println("CallNode.Value", fn.String())
	var f *Func

	switch fn.Type() {
	case NodeIdent:
		ident := fn.(IdentNode)
		f = findFunc(ident, v)
	case NodeFunc:
		ff := fn.(Func)
		f = &ff
	default:
		log.Println("<unexpected type CallNode>", fn.Type())
	}

	if f != nil {
		return applyFunc(f, node.Args[:], v)
	} else {
		return newIdentNode("<unbound>")
	}
}

var nilNode = newIdentNode("nil")

func ParseFromString(name, program string) []Node {
	return Parse(Lex(name, program))
}

func Parse(l *Lexer) []Node {
	return parser(l, make([]Node, 0), ' ')
}

func parser(l *Lexer, tree []Node, lookingFor rune) []Node {
	for item := l.nextItem(); item.Type != ItemEOF; {
		switch t := item.Type; t {
		case ItemIdent:
			tree = append(tree, newIdentNode(item.Value))
		case ItemRef:
			tree = append(tree, newRefNode(item.Value))
		case ItemString:
			tree = append(tree, newStringNode(item.Value))
		case ItemInt:
			tree = append(tree, newIntNode(item.Value))
		case ItemFloat:
			tree = append(tree, newFloatNode(item.Value))
		case ItemComplex:
			tree = append(tree, newComplexNode(item.Value))
		case ItemLeftCurl:
			tree = append(tree, newCallNode(parser(l, make([]Node, 0), '}')))
		case ItemLeftVect:
			tree = append(tree, newVectNode(parser(l, make([]Node, 0), ']')))
		case ItemLeftParen:
			tree = append(tree, newListNodeFromSlice(parser(l, make([]Node, 0), ')')))
		case ItemRightParen:
			if lookingFor != ')' {
				panic(fmt.Sprintf("unexpected \")\" [%d]", item.Pos))
			}
			return tree
		case ItemRightVect:
			if lookingFor != ']' {
				panic(fmt.Sprintf("unexpected \"]\" [%d]", item.Pos))
			}
			return tree
		case ItemRightCurl:
			if lookingFor != '}' {
				panic(fmt.Sprintf("unexpected \"}\" [%d]", item.Pos))
			}
			return tree
		case ItemError:
			println(item.Value)
		default:
			panic("Bad Item type")
		}
		item = l.nextItem()
	}

	return tree
}

func newIdentNode(name string) IdentNode {
	return IdentNode{NodeType: NodeIdent, Ident: name}
}

func newStringNode(val string) StringNode {
	return StringNode{NodeType: NodeString, Val: strings.Replace(val, "\\\"", "\"", -1)}
}

func newIntNode(val string) NumberNode {
	i, _ := strconv.ParseInt(val, 10, 64)
	return NumberNode{NodeType: NodeNumber, Val: val, Int: i, NumberType: token.INT}
}

func newFloatNode(val string) NumberNode {
	f, _ := strconv.ParseFloat(val, 64)
	return NumberNode{NodeType: NodeNumber, Val: val, Float: f, NumberType: token.FLOAT}
}

func newFloat(f float64) NumberNode {
	val := fmt.Sprintf("%f", f)
	return NumberNode{NodeType: NodeNumber, Val: val, Float: f, NumberType: token.FLOAT}
}

func newInt(i int64) NumberNode {
	val := fmt.Sprintf("%d", i)
	return NumberNode{NodeType: NodeNumber, Val: val, Int: i, NumberType: token.INT}
}

func newComplexNode(val string) NumberNode {
	return NumberNode{NodeType: NodeNumber, Val: val, NumberType: token.IMAG}
}

// We return Node here, because it could be that it's nil
func newCallNode(args []Node) Node {
	if len(args) > 0 {
		return CallNode{NodeType: NodeCall, Callee: args[0], Args: args[1:]}
	} else {
		return nilNode
	}
}

func newVectNode(content []Node) VectorNode {
	return VectorNode{NodeType: NodeVector, Nodes: content}
}

func newListNodeFromSlice(content []Node) ListNode {
	list := newListNode()
	var pair *PairNode = nil

	//for i := len(content); i > 0; i-- {
	//	node := newPairNode(content[i-1], pair)
	//	pair = &node
	//}
	for i := 0; i < len(content); i++ {
		node := newPairNode(content[i], pair)
		pair = &node
	}
	list.Head = pair

	return list
}

func newListNode() ListNode {
	return ListNode{NodeType: NodeList, Head: nil}
}

func newPairNode(first Node, second *PairNode) PairNode {
	return PairNode{NodeType: NodePair, First: first, Second: second}
}
