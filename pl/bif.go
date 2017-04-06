package pl

import (
	"fmt"
	"go/token"
	"log"
	"math"
)

func absfloat(env *Env, args []Node) Node {
	var d float64
	switch args[0].(NumberNode).NumberType {
	case token.INT:
		d = math.Abs(float64(args[0].(NumberNode).Int))
	case token.FLOAT:
		d = math.Abs(args[0].(NumberNode).Float)
	}
	return newFloat(d)
}

func and(env *Env, args []Node) Node {
	env.new_current_local(newListNode([]Node{}))

	go env.run_and(args[:])

	//ret := <-env.current.ret
	ret := env.wait_return()

	env.del_current_local()
	return ret
}

func cond(env *Env, args []Node) Node {
	env.new_current_local(newListNode([]Node{}))

	go env.run_cond(args[:])

	ret := env.wait_return()

	env.del_current_local()
	return ret
}

func def(env *Env, args []Node) Node {
	ident := args[0].(IdentNode)
	val := args[1]
	var ret Node

	if val.Type() == NodeList {
		list := val.(ListNode)
		if list.Nodes[0].Type() == NodeIdent {
			id := list.Nodes[0].(IdentNode)
			switch id.Ident {
			case "lambda":
				val = Func{mode: UserDef, ud: &Lambda{env: env, arg: list.Nodes[1], body: list.Nodes[2:]}}
				ret = list
			default:
				val = args[1].Value(env)
				ret = val
			}
		} else {
			val = args[1].Value(env)
			ret = val
		}
	} else {
		val = args[1].Value(env)
		ret = val
	}
	/*
		env.globalVars.lock.Lock()
		env.globalVars.ctx[ident] = makeVar(&val)
		env.globalVars.lock.Unlock()
	*/
	env.current.lock.Lock()
	env.current.ctx[ident] = makeVar(&val)
	env.current.lock.Unlock()
	return ret
}

func divfloat(env *Env, args []Node) Node {
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

func divint(env *Env, args []Node) Node {
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

func eq(env *Env, args []Node) Node {
	if args[0].String() == args[1].String() {
		return newIdentNode("T")
	} else {
		return newListNode([]Node{})
	}
}

func eval(env *Env, args []Node) Node {
	//log.Println(args[0])
	return args[0].Value(env)
}

func exit(env *Env, args []Node) Node {
	env.current.lock.Lock()
	env.current.cont = false
	env.current.lock.Unlock()

	env.current.exit <- args[0]
	return args[0]
}

func fold(env *Env, args []Node) Node {
	word := args[0].(IdentNode)
	init := args[1]
	list := args[2].(ListNode).Nodes
	f := findFunc(word, env)

	env.new_current_local(newListNode([]Node{}))

	go env.run_fold(f, init, list[:])

	//ret := <-env.current.ret
	ret := env.wait_return()

	env.del_current_local()
	return ret
}

func fmap(env *Env, args []Node) Node {
	word := args[0].(IdentNode)
	list := args[1].(ListNode).Nodes
	f := findFunc(word, env)

	new_list := []Node{}
	env.new_current_local(newListNode(new_list))

	go env.run_map(f, new_list[:], list[:])

	//ret := <-env.current.ret
	ret := env.wait_return()

	env.del_current_local()
	return ret
}

func gtfloat(env *Env, args []Node) Node {
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
		return newListNode([]Node{})
	}
}

func gtint(env *Env, args []Node) Node {
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
		return newListNode([]Node{})
	}
}

func ltfloat(env *Env, args []Node) Node {
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
		return newListNode([]Node{})
	}
}

func ltint(env *Env, args []Node) Node {
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
		return newListNode([]Node{})
	}
}

func neq(env *Env, args []Node) Node {
	if args[0].String() != args[1].String() {
		return newIdentNode("T")
	} else {
		return newListNode([]Node{})
	}
}

func not(env *Env, args []Node) Node {
	if args[0].String() == "()" {
		return newIdentNode("T")
	} else {
		return newListNode([]Node{})
	}
}

func or(env *Env, args []Node) Node {
	env.new_current_local(newListNode([]Node{}))

	go env.run_or(args[:])

	//ret := <-env.current.ret
	ret := env.wait_return()

	env.del_current_local()
	return ret
}

func print(env *Env, args []Node) Node {
	for _, arg := range args {
		log.Println(arg.String())
	}
	return args[len(args)-1]
}

func prog(env *Env, args []Node) Node {
	vars := args[0].(ListNode)
	env.new_current_local(vars)

	go env.run_stmt(args[1:])
	/*
			var ret Node
		Loop:
			for {
				select {
				case ret = <-env.current.ret:
					log.Println("prog: select ret", ret)
					break Loop
				case ret = <-env.current.exit:
					log.Println("prog: select exit", ret)
					break Loop
				case <-time.After(time.Second * 5):
					ret = newIdentNode("<timeout>")
					log.Println("prog: select timeout")
					break Loop
				}

			}
	*/
	ret := env.wait_return()
	env.del_current_local()
	return ret
}

func quote(env *Env, args []Node) Node {
	return args[0]
}

func set(env *Env, args []Node) Node {
	word := args[0].(IdentNode)

	env.lock.RLock()
	vars := env.current
	env.lock.RUnlock()

	for {
		vars.lock.RLock()
		if _, ok := vars.ctx[word]; ok {
			vars.ctx[word] <- args[1]
			return args[1]
		}
		if vars.next == nil {
			fmt.Println(fmt.Sprintf("Variable %s <unbound>", word.String()))
			return newIdentNode("<unbound>")
		}
		nvars := vars.next
		vars.lock.Unlock()
		vars = nvars
	}
	//return args[1]
}

func prodfloat(env *Env, args []Node) Node {
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

func prodint(env *Env, args []Node) Node {
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

func subfloat(env *Env, args []Node) Node {
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

func subint(env *Env, args []Node) Node {
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

func sumfloat(env *Env, args []Node) Node {
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

func sumint(env *Env, args []Node) Node {
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

func type_(env *Env, args []Node) Node {
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
