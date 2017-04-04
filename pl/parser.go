package pl

import (
	"fmt"
	"go/token"
	"strconv"
)

type Node interface {
	Type() NodeType
	// Position() Pos
	String() string
	Value(*Env) Node
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

func (node IdentNode) Value(env *Env) Node {
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
	return node.Val
}

func (node StringNode) Value(env *Env) Node {
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

func (node NumberNode) Value(env *Env) Node {
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
	return fmt.Sprint(node.Nodes)
}

func (node VectorNode) Value(env *Env) Node {
	vect := VectorNode{NodeType: node.Type(), Nodes: make([]Node, len(node.Nodes))}
	for i, v := range node.Nodes {
		vect.Nodes[i] = v.Value(env)
	}
	return vect
}

type ListNode struct {
	// Pos
	NodeType
	Nodes []Node
}

func (node ListNode) Copy() Node {
	vect := ListNode{NodeType: node.Type(), Nodes: make([]Node, len(node.Nodes))}
	for i, v := range node.Nodes {
		vect.Nodes[i] = v.Copy()
	}
	return vect
}

func (node ListNode) String() string {
	s := fmt.Sprint(node.Nodes)
	return "(" + s[1:len(s)-1] + ")"
}

func (node ListNode) Value(env *Env) Node {
	vect := ListNode{NodeType: node.Type(), Nodes: make([]Node, len(node.Nodes))}
	for i, v := range node.Nodes {
		vect.Nodes[i] = v.Value(env)
	}
	return vect
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

func (node CallNode) Value(env *Env) Node {
	name := node.Callee
	ident := name.Value(env).(IdentNode)

	f := findFunc(ident, env)
	if f != nil {
		return applyFunc(f, node.Args[:], env)
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
			tree = append(tree, newListNode(parser(l, make([]Node, 0), ')')))
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
	return StringNode{NodeType: NodeString, Val: val}
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

func newListNode(content []Node) ListNode {
	return ListNode{NodeType: NodeList, Nodes: content}
}
