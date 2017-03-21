package pl

import "fmt"

func quote(env *Env, args []Expression) Expression {
	return args[0]
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
		if val, ok := vars.ctx[word]; ok {
			vars.ctx[word] = args[1]
			return val
		}
		if vars.next == nil {
			fmt.Println(fmt.Sprintf("Variable %s <unbound>", word.String()))
			return NewWord("<unbound>")
		}
		vars = vars.next
	}
	//return args[1]
}

func sumint(env *Env, args []Expression) Expression {
	s := int64(0)
	for _, arg := range args {
		s += arg.(Int).number
	}
	return NewInt(s)
}
