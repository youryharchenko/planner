package pl

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type RefType int

const (
	LocalValue RefType = iota
	GlobalValue
	LocalSegmentValue
	GlobalSegmentValue
)

var RefTypeString = []string{".", ":", "!.", "!:"}

type FuncType int

const (
	BuiltIn FuncType = iota
	UserDef
)

type FuncClass int

const (
	FSubr FuncClass = iota
	Subr
)

type Vars struct {
	ctx  map[IdentNode]chan Node
	ret  chan Node
	exit chan Node
	cont bool
	next *Vars

	lock sync.RWMutex
}

type Env struct {
	globalVars *Vars
	localVars  *Vars
	current    *Vars

	lock sync.RWMutex
}

type RefNode struct {
	NodeType
	val  string
	mode RefType
	ref  IdentNode
}

func newRefNode(val string) RefNode {
	switch val[0] {
	case '.':
		return RefNode{val: val, mode: LocalValue, ref: newIdentNode(val[1:])}
	case ':':
		return RefNode{val: val, mode: GlobalValue, ref: newIdentNode(val[1:])}
	}
	return newRefNode(":<unexpected reference char>")
}

func (expr RefNode) Value(env *Env) Node {
	switch expr.mode {
	case LocalValue:
		vars := env.current
		for {
			vars.lock.RLock()
			if ch, ok := vars.ctx[expr.ref]; ok {
				if ch != nil {
					var val Node
					select {
					case val = <-ch:
						log.Println("Ref Value:", expr.ref, val)
						ch <- val
						return val
					case <-time.After(time.Second * 5):
						log.Println("Ref Value timeout", expr.ref)
						return newIdentNode("<timeout>")
					}
				} else {
					fmt.Println(fmt.Sprintf("Variable %s <unassigned>", expr.ref.String()))
					return newIdentNode("<unassigned>")
				}
			}

			if vars.next == nil {
				fmt.Println(fmt.Sprintf("Variable %s <unbound>", expr.ref.String()))
				return newIdentNode("<unbound>")
			}
			nvars := vars.next
			vars.lock.RUnlock()
			vars = nvars
		}
	case GlobalValue:
		env.globalVars.lock.RLock()
		defer env.globalVars.lock.RUnlock()
		if ch, ok := env.globalVars.ctx[expr.ref]; ok {
			if ch != nil {
				var val Node
				select {
				case val = <-ch:
					log.Println("Ref Value:", expr.ref, val)
					ch <- val
					return val
				case <-time.After(time.Second * 5):
					log.Println("Ref Value timeout", expr.ref)
					return newIdentNode("<timeout>")
				}
			} else {
				fmt.Println(fmt.Sprintf("Variable %s <unassigned>", expr.ref.String()))
				return newIdentNode("<unassigned>")
			}
		} else {
			fmt.Println(fmt.Sprintf("Variable %s <unbound>", expr.ref.String()))
			return newIdentNode("<unbound>")
		}
	}
	return newIdentNode("<unexpected>")
}

func (expr RefNode) String() string {
	return fmt.Sprintf("%s%s", RefTypeString[expr.mode], expr.ref.String())
}

func (expr RefNode) Copy() Node {
	return newRefNode(expr.val)
}

type Func struct {
	NodeType
	mode  FuncType
	class FuncClass
	bi    func(*Env, []Node) Node
	ud    *Lambda
}

func (expr Func) Value(env *Env) Node {
	return expr
}

func (expr Func) String() string {
	return fmt.Sprintf("%v, %v", expr.mode, expr.class)
}

func (expr Func) Copy() Node {
	return expr
}

type Lambda struct {
	env  *Env
	arg  Node
	body []Node
}

func (fn *Lambda) apply(args []Node, env *Env) Node {
	//log.Println(args)
	var vars ListNode
	switch fn.arg.Type() {
	case NodeIdent:
		var ident IdentNode
		var param Node
		arg := fn.arg.(IdentNode)

		if arg.Ident[0] == '*' {
			ident = newIdentNode(arg.Ident[1:])
			param = newListNode(args)
		} else {
			ident = newIdentNode(arg.Ident)
			list := make([]Node, len(args))
			for i, a := range args {
				list[i] = a.Value(env)
			}
			param = newListNode(list)
		}
		vars = newListNode([]Node{newListNode([]Node{ident, param})})

	case NodeList:
		lst := fn.arg.(ListNode)
		list := make([]Node, len(lst.Nodes))
		for i, a := range lst.Nodes {
			ident := a.(IdentNode)
			var param Node
			if ident.Ident[0] == '*' {
				ident = newIdentNode(ident.Ident[1:])
				param = args[i]
			} else {
				param = args[i].Value(env)
			}
			list[i] = newListNode([]Node{ident, param})
		}
		vars = newListNode(list)
	}

	fn.env.new_current_local(vars)

	go fn.env.run_stmt(fn.body)

	var ret Node
Loop:
	for {
		select {
		case ret = <-env.current.ret:
			log.Println("lambda: select ret", ret)
			break Loop
		case ret = <-env.current.exit:
			log.Println("lambda: select exit", ret)
			break Loop
		case <-time.After(time.Second * 5):
			ret = newIdentNode("<timeout>")
			log.Println("lambda: select timeout")
			break Loop
		}

	}

	fn.env.del_current_local()
	return ret

}

func Begin() *Env {

	global := Vars{ctx: map[IdentNode]chan Node{}, next: nil}
	local := Vars{ctx: map[IdentNode]chan Node{}, next: nil}

	global.ctx[newIdentNode("def")] = makeFunc(Func{mode: BuiltIn, class: FSubr, bi: def})
	global.ctx[newIdentNode("div$float")] = makeFunc(Func{mode: BuiltIn, class: Subr, bi: divfloat})
	global.ctx[newIdentNode("div$int")] = makeFunc(Func{mode: BuiltIn, class: Subr, bi: divint})
	global.ctx[newIdentNode("exit")] = makeFunc(Func{mode: BuiltIn, class: Subr, bi: exit})
	global.ctx[newIdentNode("fold")] = makeFunc(Func{mode: BuiltIn, class: Subr, bi: fold})
	global.ctx[newIdentNode("map")] = makeFunc(Func{mode: BuiltIn, class: Subr, bi: fmap})
	global.ctx[newIdentNode("quote")] = makeFunc(Func{mode: BuiltIn, class: FSubr, bi: quote})
	global.ctx[newIdentNode("print")] = makeFunc(Func{mode: BuiltIn, class: Subr, bi: print})
	global.ctx[newIdentNode("prod$float")] = makeFunc(Func{mode: BuiltIn, class: Subr, bi: prodfloat})
	global.ctx[newIdentNode("prod$int")] = makeFunc(Func{mode: BuiltIn, class: Subr, bi: prodint})
	global.ctx[newIdentNode("prog")] = makeFunc(Func{mode: BuiltIn, class: FSubr, bi: prog})
	global.ctx[newIdentNode("set")] = makeFunc(Func{mode: BuiltIn, class: Subr, bi: set})
	global.ctx[newIdentNode("sub$float")] = makeFunc(Func{mode: BuiltIn, class: Subr, bi: subfloat})
	global.ctx[newIdentNode("sub$int")] = makeFunc(Func{mode: BuiltIn, class: Subr, bi: subint})
	global.ctx[newIdentNode("sum$float")] = makeFunc(Func{mode: BuiltIn, class: Subr, bi: sumfloat})
	global.ctx[newIdentNode("sum$int")] = makeFunc(Func{mode: BuiltIn, class: Subr, bi: sumint})

	env := &Env{
		globalVars: &global,
		localVars:  &local,
		current:    &local,
		lock:       sync.RWMutex{},
	}

	return env
}

func (env *Env) Eval(args ...Node) Node {
	var ret Node
	for _, expr := range args {
		ret = expr.Value(env)
	}
	return ret
}

func (env *Env) SourceNodes(nodes []Node) Node {
	log.Println("SourceNode: started")
	var result Node
	for _, expr := range nodes {
		log.Println("Source:", expr.String())
		result = expr.Value(env)
		log.Println("Result:", result)
	}
	return result
}
