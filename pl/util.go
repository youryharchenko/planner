package pl

import (
	"log"
	"sync"
	"time"
)

func (v *Vars) new_current_local(name string, vars VectorNode) *Vars {
	//env.lock.Lock()

	nv := &Vars{
		name:  name,
		deep:  v.deep + 1,
		ctx:   map[IdentNode]chan Node{},
		next:  v,
		ret:   make(chan Node),
		exit:  make(chan Node),
		debug: v.debug,
		//cont: true,
		lock: sync.RWMutex{},
	}

	//env.lock.Unlock()

	for _, elm := range vars.Nodes {
		switch elm.(type) {
		case IdentNode:
			nv.lock.Lock()
			nv.ctx[elm.(IdentNode)] = makeVar(nil) //make(chan Node, 1)
			nv.lock.Unlock()
		case VectorNode:
			if llist := elm.(VectorNode); len(llist.Nodes) == 2 {
				word := llist.Nodes[0].(IdentNode)

				nv.lock.Lock()
				nv.ctx[word] = makeVar(&llist.Nodes[1])
				nv.lock.Unlock()
			}
		}
	}
	return nv
}

func (v *Vars) merge(a *Vars) *Vars {
	cv := a
	for {
		cv.lock.RLock()
		if cv.next == nil {
			cv.lock.RUnlock()
			break
		}
		for key, val := range cv.ctx {
			v.lock.Lock()
			v.ctx[key] = val
			v.lock.Unlock()
		}
		cv.lock.RUnlock()
		cv = cv.next
	}
	return v
}

func (v *Vars) del_current_local() {
	//env.lock.Lock()
	//env.current = env.current.next
	//env.lock.Unlock()
}

func (v *Vars) run_stmt(args []Node) {
	//env.current.lock.RLock()

	//if env.current.cont && len(args) >= 1 {
	if len(args) >= 1 {
		//env.current.lock.RUnlock()
		if len(args) == 1 {
			val := args[0].Value(v)
			v.ret <- val
		} else {
			go v.run_stmt(args[1:])
			args[0].Value(v)
		}
	} else {
		//env.current.lock.RUnlock()
	}
}

func (v *Vars) wait_return() Node {
	var ret Node
Loop:
	for {
		select {
		case ret = <-v.ret:
			//log.Println("prog: select ret", ret)
			break Loop
		case ret = <-v.exit:
			//log.Println("prog: select exit", ret)
			break Loop
		case <-time.After(time.Second * 20):
			ret = newIdentNode("<timeout>")
			v.printTrace()
			log.Panicf("wait_return: select timeout, deep: %d, ctx: %s", v.deep, v.name)
			break Loop
		}
	}
	return ret
}

func (v *Vars) run_cond(args []Node) {
	list := args[0].(VectorNode)

	if val := list.Nodes[0].Value(v); val.String() == "()" && len(args) > 1 {
		//env.current.lock.RLock()
		//if env.current.cont {
		//env.current.lock.RUnlock()
		go v.run_cond(args[1:])
		//} else {
		//env.current.lock.RUnlock()
		//}
	} else {
		var ret Node
		if val.String() == "()" {
			ret = newListNode()
		} else {
			//env.current.lock.Lock()
			//env.current.cont = false
			//env.current.lock.Unlock()

			nv := v.new_current_local("cond clause", newVectNode([]Node{}))

			go nv.run_stmt(list.Nodes[1:])

			ret = nv.wait_return()

			//env.del_current_local()
		}

		v.ret <- ret

	}
}

func (v *Vars) run_or(args []Node) {

	if val := args[0].Value(v); val.String() == "()" {
		//env.current.lock.RLock()
		//if env.current.cont && len(args) >= 1 {
		if len(args) >= 1 {
			//env.current.lock.RUnlock()
			if len(args) == 1 {
				v.ret <- val
			} else {
				go v.run_or(args[1:])
			}
		} else {
			//env.current.lock.RUnlock()
		}
	} else {

		//env.current.lock.Lock()
		//env.current.cont = false
		//env.current.lock.Unlock()

		v.ret <- val
	}
}

func (v *Vars) run_and(args []Node) {

	if val := args[0].Value(v); val.String() != "()" {
		//env.current.lock.RLock()
		//if env.current.cont && len(args) >= 1 {
		if len(args) >= 1 {
			//env.current.lock.RUnlock()
			if len(args) == 1 {
				v.ret <- val
			} else {
				go v.run_and(args[1:])
			}
		} else {
			//env.current.lock.RUnlock()
		}
	} else {

		//env.current.lock.Lock()
		//env.current.cont = false
		//env.current.lock.Unlock()

		v.ret <- val
	}
}

func (v *Vars) run_fold(f *Func, val Node, list ListNode) {
	//env.current.lock.RLock()
	//defer env.current.lock.RUnlock()

	//if env.current.cont && len(list) >= 1 {
	if list.Len() >= 1 {
		newVal := applyFunc(f, []Node{val, list.Nodes(0)}, v)
		if list.Len() == 1 {
			v.ret <- newVal
		} else {
			go v.run_fold(f, newVal, list.Tail(1))
		}
	}
}

func (v *Vars) run_map(f *Func, new_list ListNode, list ListNode) {
	//env.current.lock.RLock()
	//defer env.current.lock.RUnlock()

	//if env.current.cont && len(list) >= 1 {
	//new_list := newListNode()
	//log.Println(list.String(), new_list.String())
	if list.Len() >= 1 {
		//new_list = append(new_list, applyFunc(f, []Node{list[0]}, v))
		new_list = new_list.Cons(applyFunc(f, []Node{list.Nodes(0)}, v))
		if list.Len() == 1 {
			v.ret <- new_list.Rev()
		} else {
			go v.run_map(f, new_list, list.Tail(1))
		}
	}
}

func findFunc(word IdentNode, v *Vars) *Func {

	for i := 0; i < 3; i++ {
		//env.lock.RLock()
		vars := v
		//env.lock.RUnlock()

		//var f Func

		for {
			v.lock.RLock()
			if ch, ok := vars.ctx[word]; ok {
				v.lock.RUnlock()
				if ch != nil {

					var val Node
					select {
					case val = <-ch:
						ch <- val
					case <-time.After(time.Second * 5):
						log.Panicf("find function timeout: %s, deep: %d, ctx: %s", word.String(), v.deep, v.name)
					}

					//log.Println("Function found", word, val)

					switch val.Type() {
					case NodeFunc:
						f := val.(Func)
						return &f
					case NodeIdent:
						if pf := findFunc(val.(IdentNode), v); pf != nil {
							return pf
						} else {
							return nil
						}
					default:
						log.Panicf("findFunc>> unexpected type, name:%s, type: %s, value: %s", word.String(), type_(v, []Node{val}), val.String())
					}
					// goto Apply
				} else {
					log.Panicf("variable %s <unassigned>, deep: %d, ctx: %s", word.String(), v.deep, v.name)
				}
			} else {
				v.lock.RUnlock()
			}

			if vars.next == nil {
				//vars.lock.RUnlock()
				break
			}
			nvars := vars.next
			//vars.lock.RUnlock()
			vars = nvars
			//vars.lock.RLock()
		}
		time.Sleep(time.Millisecond * 1)
		v.printTrace()
		log.Printf("warning, wait function %s, deep: %d, ctx: %s", word.String(), v.deep, v.name)
	}
	v.printTrace()
	log.Panicf("variable %s <unbound>, deep: %d, ctx: %s", word.String(), v.deep, v.name)
	//env.globalVars.lock.RLock()
	/*
		if ch, ok := env.globalVars.ctx[word]; ok {
			//env.globalVars.lock.RUnlock()
			val := <-ch
			ch <- val
			f = val.(Func)
		} else {
			//env.globalVars.lock.RUnlock()
			//fmt.Println(fmt.Sprintf("Function %s <unbound>", word.String()))
			log.Panicf("function %s <unbound>", word.String())
			return nil
		}
	*/
	//Apply:
	return nil
}

func applyFunc(f *Func, args []Node, v *Vars) Node {

	if f == nil {
		return newIdentNode("<unbound>")
	}
	switch f.mode {
	case BuiltIn:
		if v.debug {
			log.Println("applyFunc:BuiltIn", f.Type(), f.mode, f.name, args, v.deep, v.name)
		}
		var list []Node
		if f.class == FSubr {
			list = args
		} else {
			list = []Node{}
			for _, elm := range args {
				list = append(list, elm.Value(v))
			}
		}
		//log.Println(f.name, list, v.deep, v.name)
		return f.bi(v, list)
	case UserDef:
		if v.debug {
			log.Println("applyFunc:UserDef", f.Type(), f.mode, f.name, args, v.deep, v.name)
		}
		return f.ud.apply(f.name, args, v)
	}
	return newIdentNode("<unexpected>")
}

func makeLambda(name string, v *Vars, arg Node, body []Node) Func {
	return Func{NodeType: NodeFunc, name: name, mode: UserDef, ud: &Lambda{vars: v, arg: arg, body: body}}
}

func makeVar(expr *Node) chan Node {
	//log.Println("makeVar", expr)
	ch := make(chan Node, 1)
	if expr != nil {
		ch <- *expr
	}
	//log.Println("makeVar", ch)
	return ch
}

func makeFunc(expr Node) chan Node {
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
