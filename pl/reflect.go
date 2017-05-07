package pl

import (
	"go/token"
	"reflect"
)

func gokindtype(v *Vars, args []Node) Node {
	return newIdentNode(args[0].(GoTypeNode).t.Kind().String())
}

func gokindvalue(v *Vars, args []Node) Node {
	return newIdentNode(args[0].(GoValueNode).v.Kind().String())
}

func gotypeof(v *Vars, args []Node) Node {
	return newGoTypeNode(reflect.TypeOf(args[0]))
}

func govalueof(v *Vars, args []Node) Node {
	return newGoValueNode(reflect.ValueOf(args[0]))
}

func gogettype(v *Vars, args []Node) Node {
	switch args[0].Type() {
	case NodeNumber:
		n := args[0].(NumberNode)
		switch n.NumberType {
		case token.INT:
			return newGoTypeNode(reflect.TypeOf(args[0].(NumberNode).Int))
		case token.FLOAT:
			return newGoTypeNode(reflect.TypeOf(args[0].(NumberNode).Float))
		}
	case NodeIdent:
		return newGoTypeNode(reflect.TypeOf(args[0].(IdentNode).Ident))
	case NodeString:
		return newGoTypeNode(reflect.TypeOf(args[0].(StringNode).Val))
	case NodeList:
		return newGoTypeNode(reflect.TypeOf(args[0].(ListNode).Head))
	case NodeVector:
		return newGoTypeNode(reflect.TypeOf(args[0].(VectorNode).Nodes))
	case NodeCall:
		return newGoTypeNode(reflect.TypeOf(args[0].(CallNode)))
	case NodeRef:
		return newGoTypeNode(reflect.TypeOf(args[0].(RefNode).val))
	case NodeFunc:
		return newGoTypeNode(reflect.TypeOf(args[0].(Func)))
	}
	return newGoTypeNode(reflect.TypeOf(args[0]))
}

func gogetvalue(v *Vars, args []Node) Node {
	switch args[0].Type() {
	case NodeNumber:
		n := args[0].(NumberNode)
		switch n.NumberType {
		case token.INT:
			return newGoValueNode(reflect.ValueOf(args[0].(NumberNode).Int))
		case token.FLOAT:
			return newGoValueNode(reflect.ValueOf(args[0].(NumberNode).Float))
		}
	case NodeIdent:
		return newGoValueNode(reflect.ValueOf(args[0].(IdentNode).Ident))
	case NodeString:
		return newGoValueNode(reflect.ValueOf(args[0].(StringNode).Val))
	case NodeList:
		return newGoValueNode(reflect.ValueOf(args[0].(ListNode).Head))
	case NodeVector:
		return newGoValueNode(reflect.ValueOf(args[0].(VectorNode).Nodes))
	case NodeCall:
		return newGoValueNode(reflect.ValueOf(args[0].(CallNode)))
	case NodeRef:
		return newGoValueNode(reflect.ValueOf(args[0].(RefNode).val))
	case NodeFunc:
		return newGoValueNode(reflect.ValueOf(args[0].(Func)))
	}
	return newGoValueNode(reflect.ValueOf(args[0]))
}

type GoTypeNode struct {
	// Pos
	NodeType
	t reflect.Type
}

func (node GoTypeNode) Copy() Node {
	return newGoTypeNode(node.t)
}

func (node GoTypeNode) String() string {
	return node.t.String()
}

func (node GoTypeNode) Value(v *Vars) Node {
	return node
}

func newGoTypeNode(t reflect.Type) GoTypeNode {
	return GoTypeNode{NodeType: NodeGoType, t: t}
}

type GoValueNode struct {
	// Pos
	NodeType
	v reflect.Value
}

func (node GoValueNode) Copy() Node {
	return newGoValueNode(node.v)
}

func (node GoValueNode) String() string {
	return node.v.String()
}

func (node GoValueNode) Value(v *Vars) Node {
	return node
}

func newGoValueNode(v reflect.Value) GoValueNode {
	return GoValueNode{NodeType: NodeGoValue, v: v}
}
