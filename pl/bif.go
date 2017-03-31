package pl

import (
	"fmt"
	"go/token"
	"log"
	"time"
)

func def(env *Env, args []Node) Node {
	ident := args[0].(IdentNode)
	val := args[1].Value(env)
	env.globalVars.ctx[ident] = makeVar(val)
	return val
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
	return newFloatNode(fmt.Sprintf("%f", d))
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
	return newIntNode(fmt.Sprintf("%d", d))
}

func exit(env *Env, args []Node) Node {
	env.current.cont = false
	env.current.exit <- args[0]
	return args[0]
}

func fold(env *Env, args []Node) Node {
	word := args[0].(IdentNode)
	init := args[1]
	list := args[2].(ListNode).Nodes
	f := findFunc(word, env)

	env.new_current_local(newListNode([]Node{}))

	env.current.ret = make(chan Node)
	env.current.cont = true

	go env.run_fold(f, init, list[:])

	ret := <-env.current.ret

	env.del_current_local()
	return ret
}

func fmap(env *Env, args []Node) Node {
	word := args[0].(IdentNode)
	list := args[1].(ListNode).Nodes
	f := findFunc(word, env)

	new_list := []Node{}
	env.new_current_local(newListNode(new_list))

	env.current.ret = make(chan Node)
	env.current.cont = true

	go env.run_map(f, new_list[:], list[:])

	ret := <-env.current.ret

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

	env.current.ret = make(chan Node)
	env.current.exit = make(chan Node)
	env.current.cont = true

	go env.run_stmt(args[1:])

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
	env.del_current_local()
	return ret
}

func quote(env *Env, args []Node) Node {
	return args[0]
}

func set(env *Env, args []Node) Node {
	word := args[0].(IdentNode)
	vars := env.current
	for {
		if _, ok := vars.ctx[word]; ok {
			vars.ctx[word] <- args[1]
			return args[1]
		}
		if vars.next == nil {
			fmt.Println(fmt.Sprintf("Variable %s <unbound>", word.String()))
			return newIdentNode("<unbound>")
		}
		vars = vars.next
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
	return newFloatNode(fmt.Sprintf("%f", p))
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
	return newIntNode(fmt.Sprintf("%d", p))
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
	return newFloatNode(fmt.Sprintf("%f", s))
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
	return newIntNode(fmt.Sprintf("%d", s))
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
	return newFloatNode(fmt.Sprintf("%f", s))
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
	return newIntNode(fmt.Sprintf("%d", s))
}
