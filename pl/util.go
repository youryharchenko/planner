package pl

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func (v *Vars) new_current_local(name string, vars VectorNode) *Vars {
	//env.lock.Lock()

	nv := &Vars{
		name:  name,
		deep:  v.deep + 1,
		vars:  map[IdentNode]chan Node{},
		next:  v,
		ret:   make(chan Node),
		exit:  make(chan Node),
		err:   make(chan Node),
		debug: v.debug,
		//cont: true,
		lock: sync.RWMutex{},
	}

	//env.lock.Unlock()

	for _, elm := range vars.Nodes {
		switch elm.(type) {
		case IdentNode:
			nv.set_var_chan(elm.(IdentNode), makeVar(nil))
			//nv.lock.Lock()
			//nv.ctx[elm.(IdentNode)] = makeVar(nil) //make(chan Node, 1)
			//nv.lock.Unlock()
		case VectorNode:
			if llist := elm.(VectorNode); len(llist.Nodes) == 2 {
				word := llist.Nodes[0].(IdentNode)
				nv.set_var_chan(word, makeVar(&llist.Nodes[1]))
				//nv.lock.Lock()
				//nv.ctx[word] = makeVar(&llist.Nodes[1])
				//nv.lock.Unlock()
			}
		}
	}
	if v.debug == true {
		log.Printf("new_current_local: new vars: %v", nv)
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
		for key, val := range cv.vars {
			v.set_var_chan(key, val)
			//v.lock.Lock()
			//v.ctx[key] = val
			//v.lock.Unlock()
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

func (v *Vars) is_bound(id IdentNode) bool {
	vars := v
	for {
		if ch := vars.get_var_chan(id); ch != nil {
			return true
		}
		if vars.next == nil {
			return false
		}
		nvars := vars.next
		vars = nvars
	}
}

func (v *Vars) is_assigned(id IdentNode) bool {
	vars := v
	for {
		if ch := vars.get_var_chan(id); ch != nil {
			if len(ch) > 0 {
				return true
			} else {
				return false
			}
		}
		if vars.next == nil {
			log.Panicf("variable %s <unbound>, deep: %d, ctx: %s", id.String(), v.deep, v.name)
			return false
		}
		nvars := vars.next
		vars = nvars
	}
}

func (v *Vars) assign(id IdentNode, val Node) Node {
	vars := v
	for {
		if ch := vars.get_var_chan(id); ch != nil {
			ch <- val
			return val
		}
		if vars.next == nil {
			log.Panicf("variable %s <unbound>, deep: %d, ctx: %s", id.String(), v.deep, v.name)
			return newIdentNode("<unbound>")
		}
		nvars := vars.next
		vars = nvars
	}
}

func (v *Vars) reassign(id IdentNode, val Node) Node {
	vars := v
	for {
		if ch := vars.get_var_chan(id); ch != nil {
			if len(ch) > 0 {
				<-ch
			}
			ch <- val
			return val
		}
		if vars.next == nil {
			log.Panicf("variable %s <unbound>, deep: %d, ctx: %s", id.String(), v.deep, v.name)
			return newIdentNode("<unbound>")
		}
		nvars := vars.next
		vars = nvars
	}
}

func (v *Vars) run_is(pat Node, expr Node) bool {

	switch pat.Type() {
	case NodeNumber:
		if pat.Type() == expr.Type() && pat.String() == expr.String() {
			return true
		} else {
			return false
		}
	case NodeIdent:
		ident := pat.(IdentNode)
		sident := ident.String()

		if sident[0] == '*' {
			vident := newIdentNode(sident[1:])
			v.reassign(vident, expr)
			return true
		} else {
			if pat.Type() == expr.Type() && pat.String() == expr.String() {
				return true
			} else {
				return false
			}
		}
	case NodeCall:
		val := pat.Value(v)
		if val.Type() == expr.Type() && val.String() == expr.String() {
			return true
		} else {
			return false
		}
	case NodeRef:
		ref := pat.(RefNode)

		if v.is_assigned(ref.ref) {
			val := pat.Value(v)
			if val.Type() == expr.Type() && val.String() == expr.String() {
				return true
			} else {
				return false
			}
		} else {
			v.assign(ref.ref, expr)
			return true
		}
	case NodeList, NodeVector:
		if pat.Type() == expr.Type() && pat.String() == expr.String() {
			return true
		} else {
			return false
		}
	}
	return false
}

func (v *Vars) run_stmt(args []Node) {
	//env.current.lock.RLock()
	//log.Printf("run_stmt: %v", args)
	//if env.current.cont && len(args) >= 1 {
	if len(args) >= 1 {
		//env.current.lock.RUnlock()
		if len(args) == 1 {
			val := args[0].Value(v)
			//log.Printf("run_stmt: val: %v", val)
			v.ret <- val
		} else {
			go v.run_stmt(args[1:])
			args[0].Value(v)
		}
	} else {
		//env.current.lock.RUnlock()
	}
}

func (v *Vars) run_catch(expr Node) {

	//v.run_stmt([]Node{expr})
	if v.debug == true {
		log.Println("run_catch: expr: ", expr.String())
	}
	ret := expr.Value(v)
	//log.Println("run_catch: ret: ", ret.String())
	v.ret <- ret

}

func (v *Vars) wait_catch_return(on_err Node) Node {
	var err, ret Node
	if v.debug {
		log.Printf("wait_catch_return: started, deep: %d, ctx: %s", v.deep, v)
	}
	for ret == nil {
	Loop:
		for {
			select {
			case ret = <-v.ret:
				if v.debug == true {
					log.Println("wait_catch_return: select ret", ret)
				}
				break Loop
			case ret = <-v.exit:
				if v.debug == true {
					log.Println("wait_catch_return: select exit", ret)
				}
				break Loop
			case err = <-v.err:
				//nv := v.next.new_current_local("catch on_error", newVectNode([]Node{}))
				//go v.next.run_stmt([]Node{on_err})
				//ret = nv.wait_return()
				//nv.del_current_local()
				if v.debug == true {
					log.Println("wait_catch_return: select err", err)
				}
				ret = on_err.Value(v)
				if v.debug == true {
					log.Println("wait_catch_return: ret", ret)
				}
				//v.ret <- ret
				break Loop
			case <-time.After(time.Second * 10):
				//ret = newIdentNode("<timeout>")
				v.lock.RLock()
				v.err <- newStringNode(fmt.Sprintf("wait_return: select timeout, deep: %d, ctx: %s", v.deep, v.name))
				v.lock.RUnlock()
				//v.printTrace()
				//log.Panicf("wait_return: select timeout, deep: %d, ctx: %s", v.deep, v.name)
				break Loop
			}
		}
	}
	if v.debug == true {
		log.Printf("wait_catch_return: finished, deep: %d, ctx: %s, ret: %v", v.deep, v.name, ret)
	}
	return ret
}

func (v *Vars) wait_return() Node {
	var ret, err Node
	//log.Printf("wait_return: started, deep: %d, ctx: %s (%v)", v.deep, v.name, v)
Loop:
	for {
		select {
		case ret = <-v.ret:
			if v.debug == true {
				log.Println("wait_return: select ret", ret)
			}
			break Loop
		case ret = <-v.exit:
			if v.debug == true {
				log.Println("wait_return: select exit", ret)
			}
			break Loop
		case err = <-v.err:
			v.lock.RLock()
			if v.next != nil {
				//log.Println(err.String())
				if v.debug == true {
					log.Printf("wait_return: transfer error: %s up, from deep: %d, ctx: %s (%v), to deep: %d, ctx: %s", err.String(), v.deep, v.name, v, v.next.deep, v.next.name)
				}
				//v.printTrace()
				v.next.err <- err
				v.lock.RUnlock()
			} else {
				v.lock.RUnlock()
			}
			//ret = err
			break Loop
		case <-time.After(time.Second * 10):
			//ret = newIdentNode("<timeout>")
			v.lock.RLock()
			v.err <- newStringNode(fmt.Sprintf("wait_return: select timeout, deep: %d, ctx: %s", v.deep, v.name))
			v.lock.RUnlock()
			//v.printTrace()
			//log.Panicf("wait_return: select timeout, deep: %d, ctx: %s", v.deep, v.name)
			break Loop
		}
	}
	return ret
}

func (v *Vars) raise_error(err string) {
	v.lock.RLock()
	v.err <- newStringNode(err)
	v.lock.RUnlock()
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
		newVal := applyFunc(f, []Node{val, list.Node(0)}, v)
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
		new_list = new_list.Cons(applyFunc(f, []Node{list.Node(0)}, v))
		if list.Len() == 1 {
			v.ret <- new_list.Rev()
		} else {
			go v.run_map(f, new_list, list.Tail(1))
		}
	}
}

func (v *Vars) get_var_chan(key IdentNode) chan Node {
	v.lock.RLock()
	if ch, ok := v.vars[key]; ok {
		v.lock.RUnlock()
		return ch
	} else {
		v.lock.RUnlock()
		return nil
	}
}

func (v *Vars) set_var_chan(key IdentNode, val chan Node) {
	v.lock.Lock()
	v.vars[key] = val
	v.lock.Unlock()
}

func findFunc(word IdentNode, v *Vars) *Func {

	for i := 0; i < 3; i++ {
		//env.lock.RLock()
		vars := v
		//env.lock.RUnlock()

		//var f Func

		for {
			//v.lock.RLock()
			if ch := vars.get_var_chan(word); ch != nil {
				//v.lock.RUnlock()
				//if ch != nil {

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
			} //else {
			//log.Panicf("variable %s <unassigned>, deep: %d, ctx: %s", word.String(), v.deep, v.name)
			//}

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
			log.Println("applyFunc:BuiltIn", f.Type(), f.mode, f.name, args, v.deep, v)
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
			log.Println("applyFunc:UserDef", f.Type(), f.mode, f.name, args, v.deep, v)
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
