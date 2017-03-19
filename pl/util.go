package pl

func (env *Env) new_current_local(vars Llist) {
	env.current = &Vars{ctx: map[Word]Expression{}, next: env.current}
	for _, elm := range vars.list {
		switch elm.(type) {
		case Word:
			env.current.ctx[elm.(Word)] = nil
		case Llist:
			if llist := elm.(Llist); len(llist.list) == 2 {
				word := llist.list[0].(Word)
				env.current.ctx[word] = llist.list[1]
			}
		}
	}
}

func (env *Env) del_current_local() {
	env.current = env.current.next
}

func (env *Env) run_stmt(args []Expression) {
	if env.current.cont && len(args) >= 1 {
		val := args[0].Value(env)
		if len(args) == 1 {
			env.current.ret <- val
		} else {
			env.run_stmt(args[1:])
		}
	}
}
