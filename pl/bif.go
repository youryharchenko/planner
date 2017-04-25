package pl

import (
	"go/token"
	"log"
	"math"
)

func absfloat(v *Vars, args []Node) Node {
	var d float64
	switch args[0].(NumberNode).NumberType {
	case token.INT:
		d = math.Abs(float64(args[0].(NumberNode).Int))
	case token.FLOAT:
		d = math.Abs(args[0].(NumberNode).Float)
	}
	return newFloat(d)
}

func and(v *Vars, args []Node) Node {
	nv := v.new_current_local("and", newVectNode([]Node{}))

	go nv.run_and(args[:])

	//ret := <-env.current.ret
	ret := nv.wait_return()

	nv.del_current_local()
	return ret
}

func car(v *Vars, args []Node) Node {
	switch args[0].Type() {
	case NodeList:
		return args[0].(ListNode).Head.First
	default:
		return args[0].(VectorNode).Nodes[0]
	}

}

func catch(v *Vars, args []Node) Node {

	nv := v.new_current_local("catch", newVectNode([]Node{}))

	go nv.run_catch(args[0])

	ret := nv.wait_catch_return(args[1])

	nv.del_current_local()

	if v.debug {
		log.Printf("BIF>> catch: ret: %v", ret)
	}
	return ret
}

func cdr(v *Vars, args []Node) Node {
	switch args[0].Type() {
	case NodeList:
		return args[0].(ListNode).Tail(1)
	default:
		return args[0].(VectorNode).Nodes[1]
	}

}

func cond(v *Vars, args []Node) Node {
	nv := v.new_current_local("cond", newVectNode([]Node{}))

	go nv.run_cond(args[:])

	ret := nv.wait_return()

	nv.del_current_local()
	if v.debug {
		log.Printf("BIF>> cond: ret: %v", ret)
	}
	return ret
}

func cons(v *Vars, args []Node) Node {
	switch args[1].Type() {
	case NodeList:
		if args[1].String() == "()" {
			return newListNode().Cons(args[0])
		} else {
			return args[1].(ListNode).Cons(args[0])
		}
	default:
		return newVectNode([]Node{args[0], args[1]})
	}
}

func cos(v *Vars, args []Node) Node {
	var s float64
	switch args[0].(NumberNode).NumberType {
	case token.INT:
		s = math.Cos(float64(args[0].(NumberNode).Int))
	case token.FLOAT:
		s = math.Cos(args[0].(NumberNode).Float)
	}
	return newFloat(s)
}

func debug_(v *Vars, args []Node) Node {
	if args[0].String() == "()" {
		v.next.debug = false
	} else {
		v.next.debug = true
	}
	return args[0]
}

func def(v *Vars, args []Node) Node {
	ident := args[0].(IdentNode)
	ret := args[1].Value(v)
	if v.debug == true {
		log.Printf("BIF>> def: ident: %s, value: %s, ctx: %v,", ident.String(), ret.String(), v)
	}
	//v.next.lock.Lock()
	//v.next.ctx[ident] = makeVar(&ret)
	//v.next.lock.Unlock()
	v.next.set_var_chan(ident, makeVar(&ret))
	return ret
}

func divfloat(v *Vars, args []Node) Node {
	var d float64
	switch args[0].(NumberNode).NumberType {
	case token.INT:
		d = float64(args[0].(NumberNode).Int)
	case token.FLOAT:
		d = args[0].(NumberNode).Float
	}

	for _, arg := range args[1:] {
		switch arg.(NumberNode).NumberType {
		case token.INT:
			d /= float64(arg.(NumberNode).Int)
		case token.FLOAT:
			d /= arg.(NumberNode).Float
		}
	}
	return newFloat(d)
}

func divint(v *Vars, args []Node) Node {
	var d int64
	switch args[0].(NumberNode).NumberType {
	case token.INT:
		d = args[0].(NumberNode).Int
	case token.FLOAT:
		d = round(args[0].(NumberNode).Float)
	}

	for _, arg := range args[1:] {
		switch arg.(NumberNode).NumberType {
		case token.INT:
			d /= arg.(NumberNode).Int
		case token.FLOAT:
			d /= round(arg.(NumberNode).Float)
		}
	}
	return newInt(d)
}

func eq(v *Vars, args []Node) Node {
	if args[0].String() == args[1].String() {
		return newIdentNode("T")
	} else {
		return newListNode()
	}
}

func eqint(v *Vars, args []Node) Node {
	var d1, d2 int64
	switch args[0].(NumberNode).NumberType {
	case token.INT:
		d1 = args[0].(NumberNode).Int
	case token.FLOAT:
		d1 = round(args[0].(NumberNode).Float)
	}
	switch args[1].(NumberNode).NumberType {
	case token.INT:
		d2 = args[1].(NumberNode).Int
	case token.FLOAT:
		d2 = round(args[1].(NumberNode).Float)
	}
	//log.Println(d1, d2, d1-d2)
	if d1 == d2 {
		return newIdentNode("T")
	} else {
		return newListNode()
	}
}

func error(v *Vars, args []Node) Node {
	err := newStringNode(args[0].String())
	v.lock.RLock()
	v.err <- err
	v.lock.RUnlock()
	return err
}

func ete(v *Vars, args []Node) Node {
	var pref string
	var ref RefNode

	switch args[0].Type() {
	case NodeIdent:
		if args[1].Type() == NodeRef {
			ref = args[1].(RefNode)
			pref = RefTypeString[ref.mode]

			return newRefNode(pref + args[0].String())
		} else {
			return args[0]
		}
	case NodeRef:
		if args[1].Type() == NodeRef {
			ref = args[1].(RefNode)
			pref = RefTypeString[ref.mode]

			return newRefNode(pref + args[0].(RefNode).ref.String())
		} else {
			return args[0].(RefNode).ref
		}
	case NodeList:
		list := args[0].(ListNode)
		switch args[1].Type() {
		case NodeVector:
			return newVectNode(list.Nodes())
		case NodeCall:
			return newCallNode(list.Nodes())
		default:
			return list
		}
	case NodeVector:
		vect := args[0].(VectorNode)
		switch args[1].Type() {
		case NodeList:
			return newListNodeFromSlice(vect.Nodes).Rev()
		case NodeCall:
			return newCallNode(vect.Nodes)
		default:
			return vect
		}
	case NodeCall:
		call := args[0].(CallNode)
		switch args[1].Type() {
		case NodeList:
			return newListNodeFromSlice(append([]Node{call.Callee}, call.Args...)).Rev()
		case NodeVector:
			return newVectNode(append([]Node{call.Callee}, call.Args...))
		default:
			return call
		}
	}
	return nil
}

func eval(v *Vars, args []Node) Node {
	//log.Println(args[0])
	return args[0].Value(v)
}

func exit(v *Vars, args []Node) Node {
	//env.current.lock.Lock()
	//env.current.cont = false
	//env.current.lock.Unlock()

	v.exit <- args[0]
	return args[0]
}

func fold(v *Vars, args []Node) Node {
	word := args[0].(IdentNode)
	init := args[1]
	var list ListNode
	switch args[2].Type() {
	case NodeList:
		list = args[2].(ListNode)
	case NodeVector:
		list = newListNodeFromSlice(args[2].(VectorNode).Nodes[:])
	}

	f := findFunc(word, v)

	nv := v.new_current_local("fold", newVectNode([]Node{}))

	go nv.run_fold(f, init, list)

	//ret := <-env.current.ret
	ret := nv.wait_return()

	nv.del_current_local()
	return ret
}

func fmap(v *Vars, args []Node) Node {
	word := args[0].(IdentNode)
	list := args[1].(ListNode)
	f := findFunc(word, v)

	nv := v.new_current_local("map", newVectNode([]Node{}))

	go nv.run_map(f, newListNode(), list)

	//ret := <-env.current.ret
	ret := nv.wait_return()

	nv.del_current_local()
	return ret.(ListNode)
}

func gtfloat(v *Vars, args []Node) Node {
	var d1, d2 float64
	switch args[0].(NumberNode).NumberType {
	case token.INT:
		d1 = float64(args[0].(NumberNode).Int)
	case token.FLOAT:
		d1 = args[0].(NumberNode).Float
	}
	switch args[1].(NumberNode).NumberType {
	case token.INT:
		d2 = float64(args[1].(NumberNode).Int)
	case token.FLOAT:
		d2 = args[1].(NumberNode).Float
	}
	//log.Println(d1, d2, d1-d2)
	if d1 > d2 {
		return newIdentNode("T")
	} else {
		return newListNode()
	}
}

func gtint(v *Vars, args []Node) Node {
	var d1, d2 int64
	switch args[0].(NumberNode).NumberType {
	case token.INT:
		d1 = args[0].(NumberNode).Int
	case token.FLOAT:
		d1 = round(args[0].(NumberNode).Float)
	}
	switch args[1].(NumberNode).NumberType {
	case token.INT:
		d2 = args[1].(NumberNode).Int
	case token.FLOAT:
		d2 = round(args[1].(NumberNode).Float)
	}
	//log.Println(d1, d2, d1-d2)
	if d1 > d2 {
		return newIdentNode("T")
	} else {
		return newListNode()
	}
}

func getjson(v *Vars, args []Node) Node {
	return newObjNode(args[0].(StringNode))
}

func lambda(v *Vars, args []Node) Node {
	return makeLambda("lambda", v, args[0], args[1:])
}

func ltfloat(v *Vars, args []Node) Node {
	var d1, d2 float64
	switch args[0].(NumberNode).NumberType {
	case token.INT:
		d1 = float64(args[0].(NumberNode).Int)
	case token.FLOAT:
		d1 = args[0].(NumberNode).Float
	}
	switch args[1].(NumberNode).NumberType {
	case token.INT:
		d2 = float64(args[1].(NumberNode).Int)
	case token.FLOAT:
		d2 = args[1].(NumberNode).Float
	}
	//log.Println(d1, d2, d1-d2)
	if d1 < d2 {
		return newIdentNode("T")
	} else {
		return newListNode()
	}
}

func ltint(v *Vars, args []Node) Node {
	var d1, d2 int64
	switch args[0].(NumberNode).NumberType {
	case token.INT:
		d1 = args[0].(NumberNode).Int
	case token.FLOAT:
		d1 = round(args[0].(NumberNode).Float)
	}
	switch args[1].(NumberNode).NumberType {
	case token.INT:
		d2 = args[1].(NumberNode).Int
	case token.FLOAT:
		d2 = round(args[1].(NumberNode).Float)
	}
	//log.Println(d1, d2, d1-d2)
	if d1 < d2 {
		return newIdentNode("T")
	} else {
		return newListNode()
	}
}

func neq(v *Vars, args []Node) Node {
	if args[0].String() != args[1].String() {
		return newIdentNode("T")
	} else {
		return newListNode()
	}
}

func not(v *Vars, args []Node) Node {
	if args[0].String() == "()" {
		return newIdentNode("T")
	} else {
		return newListNode()
	}
}

func or(v *Vars, args []Node) Node {
	nv := v.new_current_local("or", newVectNode([]Node{}))

	go nv.run_or(args[:])

	//ret := <-env.current.ret
	ret := nv.wait_return()

	nv.del_current_local()
	return ret
}

func print(v *Vars, args []Node) Node {
	for _, arg := range args {
		log.Println(arg.String())
	}
	return args[len(args)-1]
}

func let(v *Vars, args []Node) Node {
	vars := args[0].(VectorNode)
	nvars := make([]Node, len(vars.Nodes))

	for i, n := range vars.Nodes {
		nvars[i] = n.Value(v)
	}
	nv := v.new_current_local("let", newVectNode(nvars))

	go nv.run_stmt(args[1:])

	ret := nv.wait_return()
	nv.del_current_local()
	return ret
}

func quote(v *Vars, args []Node) Node {
	return args[0]
}

func set(v *Vars, args []Node) Node {
	word := args[0].(IdentNode)

	//env.lock.RLock()
	vars := v
	//env.lock.RUnlock()

	for {
		//vars.lock.RLock()
		//if _, ok := vars.ctx[word]; ok {
		//	vars.ctx[word] <- args[1]
		//	return args[1]
		//}
		if ch := vars.get_var_chan(word); ch != nil {
			ch <- args[1]
			return args[1]
		}
		if vars.next == nil {
			//fmt.Println(fmt.Sprintf("Variable %s <unbound>", word.String()))
			log.Panicf("variable %s <unbound>, deep: %d, ctx: %s", word.String(), v.deep, v.name)
			return newIdentNode("<unbound>")
		}
		nvars := vars.next
		//vars.lock.Unlock()
		vars = nvars
	}
	//return args[1]
}

func prodfloat(v *Vars, args []Node) Node {
	p := float64(1)
	for _, arg := range args {
		switch arg.(NumberNode).NumberType {
		case token.INT:
			p *= float64(arg.(NumberNode).Int)
		case token.FLOAT:
			p *= arg.(NumberNode).Float
		}
	}
	return newFloat(p)
}

func prodint(v *Vars, args []Node) Node {
	p := int64(1)
	for _, arg := range args {
		switch arg.(NumberNode).NumberType {
		case token.INT:
			p *= arg.(NumberNode).Int
		case token.FLOAT:
			p *= round(arg.(NumberNode).Float)
		}
	}
	return newInt(p)
}

func remainder(v *Vars, args []Node) Node {
	var s1, s2 int64
	switch args[0].(NumberNode).NumberType {
	case token.INT:
		s1 = args[0].(NumberNode).Int
	case token.FLOAT:
		s1 = round(args[0].(NumberNode).Float)
	}
	switch args[1].(NumberNode).NumberType {
	case token.INT:
		s2 = args[1].(NumberNode).Int
	case token.FLOAT:
		s2 = round(args[1].(NumberNode).Float)
	}
	return newInt(s1 % s2)
}

func sin(v *Vars, args []Node) Node {
	var s float64
	switch args[0].(NumberNode).NumberType {
	case token.INT:
		s = math.Sin(float64(args[0].(NumberNode).Int))
	case token.FLOAT:
		s = math.Sin(args[0].(NumberNode).Float)
	}
	return newFloat(s)
}

func subfloat(v *Vars, args []Node) Node {
	var s float64
	switch args[0].(NumberNode).NumberType {
	case token.INT:
		s = float64(args[0].(NumberNode).Int)
	case token.FLOAT:
		s = args[0].(NumberNode).Float
	}

	for _, arg := range args[1:] {
		switch arg.(NumberNode).NumberType {
		case token.INT:
			s -= float64(arg.(NumberNode).Int)
		case token.FLOAT:
			s -= arg.(NumberNode).Float
		}
	}
	return newFloat(s)
}

func subint(v *Vars, args []Node) Node {
	var s int64
	switch args[0].(NumberNode).NumberType {
	case token.INT:
		s = args[0].(NumberNode).Int
	case token.FLOAT:
		s = round(args[0].(NumberNode).Float)
	}

	for _, arg := range args[1:] {
		switch arg.(NumberNode).NumberType {
		case token.INT:
			s -= arg.(NumberNode).Int
		case token.FLOAT:
			s -= round(arg.(NumberNode).Float)
		}
	}
	return newInt(s)
}

func sumfloat(v *Vars, args []Node) Node {
	s := float64(0)
	for _, arg := range args {
		switch arg.(NumberNode).NumberType {
		case token.INT:
			s += float64(arg.(NumberNode).Int)
		case token.FLOAT:
			s += arg.(NumberNode).Float
		}
	}
	return newFloat(s)
}

func sumint(v *Vars, args []Node) Node {
	s := int64(0)
	for _, arg := range args {
		switch arg.(NumberNode).NumberType {
		case token.INT:
			s += arg.(NumberNode).Int
		case token.FLOAT:
			s += round(arg.(NumberNode).Float)
		}
	}
	return newInt(s)
}

func type_(v *Vars, args []Node) Node {
	var t string
	switch args[0].Type() {
	case NodeCall:
		t = "Call"
	case NodeIdent:
		t = "Id"
	case NodeList:
		t = "List"
	case NodeNumber:
		t = "Num"
	case NodeRef:
		t = "Ref"
	case NodeString:
		t = "Str"
	case NodeVector:
		t = "Vect"
	}
	return newIdentNode(t)
}
