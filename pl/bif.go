package pl

import (
	"fmt"
	"log"
	"time"
)

func quote(env *Env, args []Expression) Expression {
	return args[0]
}

func exit(env *Env, args []Expression) Expression {
	env.current.cont = false
	env.current.exit <- args[0]
	return args[0]
}

func fold(env *Env, args []Expression) Expression {
	word := args[0].(Word)
	init := args[1]
	list := args[2].(Llist).list
	f := findFunc(word, env)

	env.new_current_local(NewLlist())

	env.current.ret = make(chan Expression)
	env.current.cont = true

	go env.run_fold(f, init, list[:])

	ret := <-env.current.ret

	env.del_current_local()
	return ret
}

func fmap(env *Env, args []Expression) Expression {
	word := args[0].(Word)
	list := args[1].(Llist).list
	f := findFunc(word, env)

	new_list := []Expression{}
	env.new_current_local(NewLlist())

	env.current.ret = make(chan Expression)
	env.current.cont = true

	go env.run_map(f, new_list[:], list[:])

	ret := <-env.current.ret

	env.del_current_local()
	return ret
}

func print(env *Env, args []Expression) Expression {
	for _, arg := range args {
		log.Println(arg.String())
	}
	return args[len(args)-1]
}

func prog(env *Env, args []Expression) Expression {
	vars := args[0].(Llist)
	env.new_current_local(vars)

	env.current.ret = make(chan Expression)
	env.current.exit = make(chan Expression)
	env.current.cont = true

	go env.run_stmt(args[1:])

	var ret Expression
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
			ret = NewWord("<timeout>")
			log.Println("prog: select timeout")
			break Loop
		}

	}
	env.del_current_local()
	return ret
}

func set(env *Env, args []Expression) Expression {
	word := args[0].(Word)
	vars := env.current
	for {
		if _, ok := vars.ctx[word]; ok {
			vars.ctx[word] <- args[1]
			return args[1]
		}
		if vars.next == nil {
			fmt.Println(fmt.Sprintf("Variable %s <unbound>", word.String()))
			return NewWord("<unbound>")
		}
		vars = vars.next
	}
	//return args[1]
}

func prodfloat(env *Env, args []Expression) Expression {
	p := float64(1)
	for _, arg := range args {
		switch arg.(type) {
		case Int:
			p *= float64(arg.(Int).number)
		case Float:
			p *= arg.(Float).number
		}
	}
	return NewFloat(p)
}

func prodint(env *Env, args []Expression) Expression {
	p := int64(1)
	for _, arg := range args {
		switch arg.(type) {
		case Int:
			p *= arg.(Int).number
		case Float:
			p *= round(arg.(Float).number)
		}
	}
	return NewInt(p)
}

func sumfloat(env *Env, args []Expression) Expression {
	s := float64(0)
	for _, arg := range args {
		switch arg.(type) {
		case Int:
			s += float64(arg.(Int).number)
		case Float:
			s += arg.(Float).number
		}
	}
	return NewFloat(s)
}

func sumint(env *Env, args []Expression) Expression {
	s := int64(0)
	for _, arg := range args {
		switch arg.(type) {
		case Int:
			s += arg.(Int).number
		case Float:
			s += round(arg.(Float).number)
		}
	}
	return NewInt(s)
}
