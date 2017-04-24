package pl

import "github.com/tidwall/gjson"

type ObjNode struct {
	NodeType
	props ListNode
}

func newObjNode(src StringNode) ObjNode {
	//log.Println(src.String())
	s := src.String()
	l := len(s)
	json := gjson.Parse(s[1 : l-1])
	return makeObj(json, ObjNode{NodeType: NodeObj, props: newListNode()}).(ObjNode)
}

func makeObj(json gjson.Result, obj ObjNode) Node {
	//log.Println(json.Type, json.Value())
	switch json.Type.String() {
	case "JSON":
		o := ObjNode{NodeType: NodeObj, props: newListNode()}
		switch json.Value().(type) {
		case map[string]interface{}:
			json.ForEach(func(key gjson.Result, value gjson.Result) bool {
				o.props = o.props.Cons(newVectNode([]Node{newIdentNode(key.String()), makeObj(value, o)}))
				return true
			})
		case []interface{}:
			json.ForEach(func(key gjson.Result, value gjson.Result) bool {
				o.props = o.props.Cons(makeObj(value, o))
				return true
			})
		}
		obj.props = o.props.Rev()
		return obj
	case "String":
		return newStringNode("\"" + json.String() + "\"")
	case "Number":
		return newFloat(json.Float())
	case "True":
		return newIdentNode("T")
	case "False":
		return newListNode()
	case "Null":
		return newIdentNode("NULL")
	}
	return obj
}

func (obj ObjNode) Value(v *Vars) Node {
	return ObjNode{NodeType: NodeObj, props: obj.props.Value(v).(ListNode)}
}

func (obj ObjNode) String() string {
	return obj.props.String()
}

func (obj ObjNode) Copy() Node {
	return ObjNode{NodeType: NodeObj, props: obj.props.Copy().(ListNode)}
}
