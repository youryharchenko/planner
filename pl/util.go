package pl

import "fmt"

func (env *Env) new_current_local(vars Llist) {
	env.current = &Vars{ctx: map[Word]chan Expression{}, next: env.current}
	for _, elm := range vars.list {
		switch elm.(type) {
		case Word:
			env.current.ctx[elm.(Word)] = make(chan Expression, 1)
		case Llist:
			if llist := elm.(Llist); len(llist.list) == 2 {
				word := llist.list[0].(Word)
				env.current.ctx[word] = makeVar(llist.list[1])
			}
		}
	}
}

func (env *Env) del_current_local() {
	env.current = env.current.next
}

func (env *Env) run_stmt(args []Expression) {
	if env.current.cont && len(args) >= 1 {
		if len(args) == 1 {
			val := args[0].Value(env)
			env.current.ret <- val
		} else {
			go env.run_stmt(args[1:])
			args[0].Value(env)
		}
	}
}

func (env *Env) run_fold(f *Func, val Expression, list []Expression) {
	if env.current.cont && len(list) >= 1 {
		newVal := applyFunc(f, []Expression{val, list[0]}, env)
		if len(list) == 1 {
			env.current.ret <- newVal
		} else {
			go env.run_fold(f, newVal, list[1:])
		}
	}
}

func findFunc(word Word, env *Env) *Func {
	vars := env.current
	var f Func
	for {
		if ch, ok := vars.ctx[word]; ok {
			val := <-ch
			ch <- val
			f = val.(Func)
			goto Apply
		}
		if vars.next == nil {
			break
		}
		vars = vars.next
	}
	if ch, ok := env.globalVars.ctx[word]; ok {
		val := <-ch
		ch <- val
		f = val.(Func)
	} else {
		fmt.Println(fmt.Sprintf("Function %s <unbound>", word.String()))
		return nil
	}
Apply:
	return &f
}

func applyFunc(f *Func, args []Expression, env *Env) Expression {
	switch f.mode {
	case BuiltIn:
		var list []Expression
		if f.class == FSubr {
			list = args
		} else {
			list = []Expression{}
			for _, elm := range args {
				list = append(list, elm.Value(env))
			}
		}
		return f.bi(env, list)
	}
	return NewWord("<unexpected>")
}

func makeVar(expr Expression) chan Expression {
	//log.Println("makeVar", expr)
	ch := make(chan Expression, 1)
	ch <- expr
	//log.Println("makeVar", ch)
	return ch
}

func round(val float64) int64 {
	if val < 0 {
		return int64(val - 0.5)
	}
	return int64(val + 0.5)
}
