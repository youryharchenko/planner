package pl

import (
	"log"
	"math"
	"time"
)

func absfloat(v *Vars, args []Node) Node {
	//var d float64
	//switch args[0].(NumberNode).NumberType {
	//case token.INT:
	//	d = math.Abs(float64(args[0].(NumberNode).Int))
	//case token.FLOAT:
	//	d = math.Abs(args[0].(NumberNode).Float)
	//}
	return newFloat(math.Abs(number_to_float(args[0].(NumberNode))))
}

func and(v *Vars, args []Node) Node {
	nv := v.new_current_local("and", newVectNode([]Node{}))
	go nv.run_and(args[:])
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
	//var s float64
	//switch args[0].(NumberNode).NumberType {
	//case token.INT:
	//	s = math.Cos(float64(args[0].(NumberNode).Int))
	//case token.FLOAT:
	//	s = math.Cos(args[0].(NumberNode).Float)
	//}
	return newFloat(math.Cos(number_to_float(args[0].(NumberNode))))
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
	v.next.set_var_chan(ident, makeVar(&ret))
	return ret
}

func divfloat(v *Vars, args []Node) Node {
	d := number_to_float(args[0].(NumberNode))
	//var d float64
	//switch args[0].(NumberNode).NumberType {
	//case token.INT:
	//	d = float64(args[0].(NumberNode).Int)
	//case token.FLOAT:
	//	d = args[0].(NumberNode).Float
	//}

	for _, arg := range args[1:] {
		d /= number_to_float(arg.(NumberNode))
		//	switch arg.(NumberNode).NumberType {
		//	case token.INT:
		//		d /= float64(arg.(NumberNode).Int)
		//	case token.FLOAT:
		//		d /= arg.(NumberNode).Float
		//	}
	}
	return newFloat(d)
}

func divint(v *Vars, args []Node) Node {
	d := number_to_int(args[0].(NumberNode))
	//var d int64
	//switch args[0].(NumberNode).NumberType {
	//case token.INT:
	//	d = args[0].(NumberNode).Int
	//case token.FLOAT:
	//	d = round(args[0].(NumberNode).Float)
	//}

	for _, arg := range args[1:] {
		d /= number_to_int(arg.(NumberNode))
		//switch arg.(NumberNode).NumberType {
		//case token.INT:
		//	d /= arg.(NumberNode).Int
		//case token.FLOAT:
		//	d /= round(arg.(NumberNode).Float)
		//}
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
	d1 := number_to_int(args[0].(NumberNode))
	d2 := number_to_int(args[1].(NumberNode))
	//var d1, d2 int64
	//switch args[0].(NumberNode).NumberType {
	//case token.INT:
	//	d1 = args[0].(NumberNode).Int
	//case token.FLOAT:
	//	d1 = round(args[0].(NumberNode).Float)
	//}
	//switch args[1].(NumberNode).NumberType {
	//case token.INT:
	//	d2 = args[1].(NumberNode).Int
	//case token.FLOAT:
	//	d2 = round(args[1].(NumberNode).Float)
	//}
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
	ret := nv.wait_return()
	nv.del_current_local()
	return ret.(ListNode)
}

func gtfloat(v *Vars, args []Node) Node {
	d1 := number_to_float(args[0].(NumberNode))
	d2 := number_to_float(args[1].(NumberNode))
	//var d1, d2 float64
	//switch args[0].(NumberNode).NumberType {
	//case token.INT:
	//	d1 = float64(args[0].(NumberNode).Int)
	//case token.FLOAT:
	//	d1 = args[0].(NumberNode).Float
	//}
	//switch args[1].(NumberNode).NumberType {
	//case token.INT:
	//	d2 = float64(args[1].(NumberNode).Int)
	//case token.FLOAT:
	//	d2 = args[1].(NumberNode).Float
	//}
	//log.Println(d1, d2, d1-d2)
	if d1 > d2 {
		return newIdentNode("T")
	} else {
		return newListNode()
	}
}

func gtint(v *Vars, args []Node) Node {
	d1 := number_to_int(args[0].(NumberNode))
	d2 := number_to_int(args[1].(NumberNode))
	//var d1, d2 int64
	//switch args[0].(NumberNode).NumberType {
	//case token.INT:
	//	d1 = args[0].(NumberNode).Int
	//case token.FLOAT:
	//	d1 = round(args[0].(NumberNode).Float)
	//}
	//switch args[1].(NumberNode).NumberType {
	//case token.INT:
	//	d2 = args[1].(NumberNode).Int
	//case token.FLOAT:
	//	d2 = round(args[1].(NumberNode).Float)
	//}
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

func is(v *Vars, args []Node) Node {
	f := newListNode()
	t := newIdentNode("T")

	pat := args[0]
	expr := args[1].Value(v)

	v.rb = map[IdentNode]Node{}

	if v.run_is(pat, expr) {
		return t
	} else {
		v.lock_rb.RLock()
		for key, node := range v.rb {
			//log.Printf("roll_back id: %s, val: %s", key, node)
			if node != nil {
				v.reassign(key, node)
			}
		}
		v.lock_rb.RUnlock()
		return f
	}
}

func kappa(v *Vars, args []Node) Node {
	return makeKappa("kappa", v, args[0], args[1:])
}

func lambda(v *Vars, args []Node) Node {
	return makeLambda("lambda", v, args[0], args[1:])
}

func lenvect(v *Vars, args []Node) Node {
	return newInt(int64(len(args[0].(VectorNode).Nodes)))
}

func ltfloat(v *Vars, args []Node) Node {
	d1 := number_to_float(args[0].(NumberNode))
	d2 := number_to_float(args[1].(NumberNode))
	//var d1, d2 float64
	//switch args[0].(NumberNode).NumberType {
	//case token.INT:
	//	d1 = float64(args[0].(NumberNode).Int)
	//case token.FLOAT:
	//	d1 = args[0].(NumberNode).Float
	//}
	//switch args[1].(NumberNode).NumberType {
	//case token.INT:
	//	d2 = float64(args[1].(NumberNode).Int)
	//case token.FLOAT:
	//	d2 = args[1].(NumberNode).Float
	//}
	//log.Println(d1, d2, d1-d2)
	if d1 < d2 {
		return newIdentNode("T")
	} else {
		return newListNode()
	}
}

func ltint(v *Vars, args []Node) Node {
	d1 := number_to_int(args[0].(NumberNode))
	d2 := number_to_int(args[1].(NumberNode))
	//var d1, d2 int64
	//switch args[0].(NumberNode).NumberType {
	//case token.INT:
	//	d1 = args[0].(NumberNode).Int
	//case token.FLOAT:
	//	d1 = round(args[0].(NumberNode).Float)
	//}
	//switch args[1].(NumberNode).NumberType {
	//case token.INT:
	//	d2 = args[1].(NumberNode).Int
	//case token.FLOAT:
	//	d2 = round(args[1].(NumberNode).Float)
	//}
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

func omega(v *Vars, args []Node) Node {
	return makeOmega("omega", v, args[0], args[1:])
}

func or(v *Vars, args []Node) Node {
	nv := v.new_current_local("or", newVectNode([]Node{}))
	go nv.run_or(args[:])
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
	go nv.run_stmt_sync(args[1:])
	ret := nv.wait_return()
	nv.del_current_local()
	return ret
}

func letasync(v *Vars, args []Node) Node {
	vars := args[0].(VectorNode)
	nvars := make([]Node, len(vars.Nodes))
	for i, n := range vars.Nodes {
		nvars[i] = n.Value(v)
	}
	nv := v.new_current_local("let", newVectNode(nvars))
	go nv.run_stmt_async(args[1:])
	ret := nv.wait_return()
	nv.del_current_local()
	return ret
}

func quote(v *Vars, args []Node) Node {
	if args[0].Type() == NodeList {
		return args[0].(ListNode).Rev()
	} else {
		return args[0]
	}
}

func set(v *Vars, args []Node) Node {
	word := args[0].(IdentNode)
	return v.assign(word, args[1])
}

func reset(v *Vars, args []Node) Node {
	word := args[0].(IdentNode)
	return v.reassign(word, args[1])
}

func prodfloat(v *Vars, args []Node) Node {
	p := float64(1)
	for _, arg := range args {
		p *= number_to_float(arg.(NumberNode))
		//switch arg.(NumberNode).NumberType {
		//case token.INT:
		//	p *= float64(arg.(NumberNode).Int)
		//case token.FLOAT:
		//	p *= arg.(NumberNode).Float
		//}
	}
	return newFloat(p)
}

func prodint(v *Vars, args []Node) Node {
	p := int64(1)
	for _, arg := range args {
		p *= number_to_int(arg.(NumberNode))
		//switch arg.(NumberNode).NumberType {
		//case token.INT:
		//	p *= arg.(NumberNode).Int
		//case token.FLOAT:
		//	p *= round(arg.(NumberNode).Float)
		//}
	}
	return newInt(p)
}

func remainder(v *Vars, args []Node) Node {
	s1 := number_to_int(args[0].(NumberNode))
	s2 := number_to_int(args[1].(NumberNode))
	//var s1, s2 int64
	//switch args[0].(NumberNode).NumberType {
	//case token.INT:
	//	s1 = args[0].(NumberNode).Int
	//case token.FLOAT:
	//	s1 = round(args[0].(NumberNode).Float)
	//}
	//switch args[1].(NumberNode).NumberType {
	//case token.INT:
	//	s2 = args[1].(NumberNode).Int
	//case token.FLOAT:
	//	s2 = round(args[1].(NumberNode).Float)
	//}
	return newInt(s1 % s2)
}

func send(v *Vars, args []Node) Node {
	//log.Println("send to actor: ", args[0], args[1])
	actor := args[0].(ActorInst)
	actor.in <- args[1]
	return actor
}

func sin(v *Vars, args []Node) Node {
	//var s float64
	//switch args[0].(NumberNode).NumberType {
	//case token.INT:
	//	s = math.Sin(float64(args[0].(NumberNode).Int))
	//case token.FLOAT:
	//	s = math.Sin(args[0].(NumberNode).Float)
	//}
	return newFloat(math.Sin(number_to_float(args[0].(NumberNode))))
}

func sleep(v *Vars, args []Node) Node {
	s := number_to_int(args[0].(NumberNode))
	//var s int64
	//switch args[0].(NumberNode).NumberType {
	//case token.INT:
	//	s = args[0].(NumberNode).Int
	//case token.FLOAT:
	//	s = round(args[0].(NumberNode).Float)
	//}
	<-time.After(time.Millisecond * time.Duration(s))
	return args[0]
}

func start(v *Vars, args []Node) Node {
	var actor ActorInst
	//log.Println("start actor ", args[0])
	switch args[0].Type() {
	case NodeIdent:
		actor = newActorInst(findActor(args[0].(IdentNode), v), args[1])
	case NodeActor:
		actor = newActorInst(args[0].(Actor), args[1])
	case NodeActorInst:
		actor = args[0].(ActorInst)
	}
	go v.run_actor(&actor)
	return actor
}

func subfloat(v *Vars, args []Node) Node {
	s := number_to_float(args[0].(NumberNode))
	//var s float64
	//switch args[0].(NumberNode).NumberType {
	//case token.INT:
	//	s = float64(args[0].(NumberNode).Int)
	//case token.FLOAT:
	//	s = args[0].(NumberNode).Float
	//}

	for _, arg := range args[1:] {
		s -= number_to_float(arg.(NumberNode))
		//	switch arg.(NumberNode).NumberType {
		//	case token.INT:
		//		s -= float64(arg.(NumberNode).Int)
		//	case token.FLOAT:
		//		s -= arg.(NumberNode).Float
		//	}
	}
	return newFloat(s)
}

func subint(v *Vars, args []Node) Node {
	s := number_to_int(args[0].(NumberNode))
	//var s int64
	//switch args[0].(NumberNode).NumberType {
	//case token.INT:
	//	s = args[0].(NumberNode).Int
	//case token.FLOAT:
	//	s = round(args[0].(NumberNode).Float)
	//}

	for _, arg := range args[1:] {
		s -= number_to_int(arg.(NumberNode))
		//switch arg.(NumberNode).NumberType {
		//case token.INT:
		//	s -= arg.(NumberNode).Int
		//case token.FLOAT:
		//	s -= round(arg.(NumberNode).Float)
		//}
	}
	return newInt(s)
}

func sumfloat(v *Vars, args []Node) Node {
	s := float64(0)
	for _, arg := range args {
		s += number_to_float(arg.(NumberNode))
		//switch arg.(NumberNode).NumberType {
		//case token.INT:
		//	s += float64(arg.(NumberNode).Int)
		//case token.FLOAT:
		//	s += arg.(NumberNode).Float
		//}
	}
	return newFloat(s)
}

func sumint(v *Vars, args []Node) Node {
	s := int64(0)
	for _, arg := range args {
		s += number_to_int(arg.(NumberNode))
		//switch arg.(NumberNode).NumberType {
		//case token.INT:
		//	s += arg.(NumberNode).Int
		//case token.FLOAT:
		//	s += round(arg.(NumberNode).Float)
		//}
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
	case NodeActor:
		t = "Actor"
	case NodeGoType:
		t = "GoType"
	case NodeGoValue:
		t = "GoValue"
	}
	return newIdentNode(t)
}
