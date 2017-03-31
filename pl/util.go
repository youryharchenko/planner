package pl

import "fmt"

func (env *Env) new_current_local(vars ListNode) {
	env.current = &Vars{ctx: map[IdentNode]chan Node{}, next: env.current}
	for _, elm := range vars.Nodes {
		switch elm.(type) {
		case IdentNode:
			env.current.ctx[elm.(IdentNode)] = make(chan Node, 1)
		case ListNode:
			if llist := elm.(ListNode); len(llist.Nodes) == 2 {
				word := llist.Nodes[0].(IdentNode)
				env.current.ctx[word] = makeVar(llist.Nodes[1])
			}
		}
	}
}

func (env *Env) del_current_local() {
	env.current = env.current.next
}

func (env *Env) run_stmt(args []Node) {
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

func (env *Env) run_fold(f *Func, val Node, list []Node) {
	if env.current.cont && len(list) >= 1 {
		newVal := applyFunc(f, []Node{val, list[0]}, env)
		if len(list) == 1 {
			env.current.ret <- newVal
		} else {
			go env.run_fold(f, newVal, list[1:])
		}
	}
}

func (env *Env) run_map(f *Func, new_list []Node, list []Node) {
	if env.current.cont && len(list) >= 1 {
		new_list = append(new_list, applyFunc(f, []Node{list[0]}, env))
		if len(list) == 1 {
			env.current.ret <- newListNode(new_list)
		} else {
			go env.run_map(f, new_list, list[1:])
		}
	}
}

func findFunc(word IdentNode, env *Env) *Func {
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

func applyFunc(f *Func, args []Node, env *Env) Node {
	switch f.mode {
	case BuiltIn:
		var list []Node
		if f.class == FSubr {
			list = args
		} else {
			list = []Node{}
			for _, elm := range args {
				list = append(list, elm.Value(env))
			}
		}
		return f.bi(env, list)
	}
	return newIdentNode("<unexpected>")
}

func makeVar(expr Node) chan Node {
	//log.Println("makeVar", expr)
	ch := make(chan Node, 1)
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
