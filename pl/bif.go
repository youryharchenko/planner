package pl

import "fmt"

func quote(env *Env, args []Expression) Expression {
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

func prog(env *Env, args []Expression) Expression {
	vars := args[0].(Llist)
	env.new_current_local(vars)

	env.current.ret = make(chan Expression)
	env.current.cont = true
	go env.run_stmt(args[1:])
	ret := <-env.current.ret

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
