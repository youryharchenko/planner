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
	MatchBuiltIn
	MatchUserDef
)

type FuncClass int

const (
	FSubr FuncClass = iota
	Subr
)

type Vars struct {
	name string
	deep int64
	vars map[IdentNode]chan Node
	ret  chan Node
	exit chan Node
	err  chan Node
	//cont bool
	next *Vars

	debug bool
	rb    map[IdentNode]Node

	lock    sync.RWMutex
	lock_rb sync.RWMutex
}

func (v *Vars) findGlobal() *Vars {
	cv := v
	for {
		if cv.next == nil {
			return cv
		} else {
			cv = cv.next
		}
	}
}

func (v *Vars) printTrace() {
	cv := v
	for {
		if cv.next == nil {
			return
		}
		log.Printf("Trace>> deep: %d, ctx: %s %v", cv.deep, cv.name, cv.String())
		cv = cv.next
	}
}

func (v *Vars) String() string {
	s := fmt.Sprintf("name: %s, debug: %v, vars: [", v.name, v.debug)
	v.lock.RLock()
	for k, v := range v.vars {
		n := <-v
		v <- n
		s += k.String() + ": " + n.String() + ","
	}
	v.lock.RUnlock()
	return s + "]"
}

type Env struct {
	globalVars *Vars
	//localVars  *Vars
	//current *Vars

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
		return RefNode{NodeType: NodeRef, val: val, mode: LocalValue, ref: newIdentNode(val[1:])}
	case ':':
		return RefNode{NodeType: NodeRef, val: val, mode: GlobalValue, ref: newIdentNode(val[1:])}
	}
	log.Panicln(":<unexpected reference char>")
	return newRefNode(":<unexpected reference char>")
}

func (expr RefNode) Value(v *Vars) Node {
	switch expr.mode {
	case LocalValue:
		for i := 0; i < 3; i++ {
			vars := v
			for {
				if ch := vars.get_var_chan(expr.ref); ch != nil {
					var val Node
					select {
					case val = <-ch:
						//log.Println("Ref Value:", expr.ref, val)
						ch <- val
						return val
					case <-time.After(time.Second * 5):
						//log.Println("Ref Value timeout", expr.ref)
						//return newIdentNode("<timeout>")
						//log.Panicln("ref value timeout", expr.ref)
						log.Panicf("ref value timeout: %s, deep: %d, ctx: %s", expr.ref.String(), v.deep, v.name)
					}
				} //else {
				//fmt.Println(fmt.Sprintf("Variable %s <unassigned>", expr.ref.String()))
				//return newIdentNode("<unassigned>")
				//log.Panicf("Variable %s <unassigned>", expr.ref.String())
				//log.Panicf("variable %s <unassigned>, deep: %d, ctx: %s", expr.ref.String(), v.deep, v.name)
				//}

				if vars.next == nil {
					break
					//fmt.Println(fmt.Sprintf("Variable %s <unbound>", expr.ref.String()))
					//return newIdentNode("<unbound>")
					//log.Panicf("Variable %s <unbound>", expr.ref.String())
					//v.printTrace()
					//log.Panicf("variable %s <unbound>, deep: %d, ctx: %s(%v)", expr.ref.String(), v.deep, v.name, v.ctx)
				}
				nvars := vars.next
				//vars.lock.RUnlock()
				vars = nvars
			}
			time.Sleep(time.Millisecond * 1)
			log.Printf("warning, wait variable %s, deep: %d, ctx: %s", expr.ref.String(), v.deep, v.name)
		}
		log.Panicf("variable %s <unbound>, deep: %d, ctx: %s(%v)", expr.ref.String(), v.deep, v.name, v.vars)
	case GlobalValue:
		globalVars := v.findGlobal()
		//globalVars.lock.RLock()
		//defer globalVars.lock.RUnlock()
		//if ch, ok := globalVars.ctx[expr.ref]; ok {
		if ch := globalVars.get_var_chan(expr.ref); ch != nil {
			var val Node
			select {
			case val = <-ch:
				//log.Println("Ref Value:", expr.ref, val)
				ch <- val
				return val
			case <-time.After(time.Second * 5):
				//log.Println("Ref Value timeout", expr.ref)
				//return newIdentNode("<timeout>")
				//log.Panicln("Ref Value timeout", expr.ref)
				log.Panicf("ref value timeout: %s, deep: %d, ctx: %s", expr.ref.String(), v.deep, v.name)
			}
		} else {
			//fmt.Println(fmt.Sprintf("Variable %s <unbound>", expr.ref.String()))
			//return newIdentNode("<unbound>")
			//log.Panicf("Variable %s <unbound>", expr.ref.String())
			log.Panicf("variable %s <unbound>, deep: %d, ctx: %s", expr.ref.String(), v.deep, v.name)
		}
	}
	//return newIdentNode("<unexpected>")
	log.Panicln("<unexpected>")
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
	name  string
	mode  FuncType
	class FuncClass
	bi    func(*Vars, []Node) Node
	mbi   func(*Vars, []Node, Node) bool
	ud    *Lambda
	mud   *Kappa
}

func (expr Func) Value(v *Vars) Node {
	return expr
}

func (expr Func) String() string {
	return fmt.Sprintf("%v", expr.name)
}

func (expr Func) Copy() Node {
	return expr
}

type Lambda struct {
	vars *Vars
	arg  Node
	body []Node
}

func (fn *Lambda) apply(name string, args []Node, v *Vars) Node {
	//log.Println("Lambda: args", args)
	var vars VectorNode

	switch fn.arg.Type() {
	case NodeIdent:
		var ident IdentNode
		var param Node
		arg := fn.arg.(IdentNode)

		if arg.Ident[0] == '*' {
			ident = newIdentNode(arg.Ident[1:])
			param = newVectNode(args)
		} else {
			ident = newIdentNode(arg.Ident)
			list := make([]Node, len(args))
			for i, a := range args {
				list[i] = a.Value(v)
			}
			param = newVectNode(list)
		}
		vars = newVectNode([]Node{newVectNode([]Node{ident, param})})
	case NodeVector:
		lst := fn.arg.(VectorNode)
		list := make([]Node, len(lst.Nodes))
		for i, a := range lst.Nodes {
			ident := a.(IdentNode)
			var param Node
			if ident.Ident[0] == '*' {
				ident = newIdentNode(ident.Ident[1:])
				param = args[i]
			} else {
				param = args[i].Value(v)
			}
			list[i] = newVectNode([]Node{ident, param})
		}
		vars = newVectNode(list)

	}
	//log.Println(name, vars, v.deep, v.name)
	//nv := fn.vars.new_current_local(name, vars)
	fnv := v.new_current_local(name+"-def", newVectNode([]Node{})).merge(fn.vars)
	//log.Println(fnv)
	go fnv.wait_return()

	nv := fnv.new_current_local(name+"-run", vars)

	if v.debug {
		//v.printTrace()
		//log.Println("=============")
		//fn.vars.printTrace()
	}

	go nv.run_stmt_sync(fn.body)

	ret := nv.wait_return()

	nv.del_current_local()
	return ret

}

type Kappa struct {
	vars *Vars
	arg  Node
	//expr Node
	body []Node
}

func (fn *Kappa) apply(name string, args []Node, expr Node, v *Vars) bool {
	//log.Println("Kappa: args", args)
	var vars VectorNode

	switch fn.arg.Type() {
	case NodeIdent:
		var ident IdentNode
		var param Node
		arg := fn.arg.(IdentNode)

		if arg.Ident[0] == '*' {
			ident = newIdentNode(arg.Ident[1:])
			param = newVectNode(args)
		} else {
			ident = newIdentNode(arg.Ident)
			list := make([]Node, len(args))
			for i, a := range args {
				list[i] = a.Value(v)
			}
			param = newVectNode(list)
		}
		vars = newVectNode([]Node{newVectNode([]Node{ident, param})})
	case NodeVector:
		lst := fn.arg.(VectorNode)
		list := make([]Node, len(lst.Nodes))
		for i, a := range lst.Nodes {
			ident := a.(IdentNode)
			var param Node
			if ident.Ident[0] == '*' {
				ident = newIdentNode(ident.Ident[1:])
				param = args[i]
			} else {
				param = args[i].Value(v)
			}
			list[i] = newVectNode([]Node{ident, param})
		}
		vars = newVectNode(list)

	}
	//log.Println(name, vars, v.deep, v.name)
	//nv := fn.vars.new_current_local(name, vars)
	fnv := v.new_current_local(name+"-def", newVectNode([]Node{})).merge(fn.vars)
	//log.Println(fnv)
	//go fnv.wait_return()

	nv := fnv.new_current_local(name+"-run", vars)

	ret := nv.run_is(fn.body[0], expr)

	if v.debug {
		//v.printTrace()
		//log.Println("=============")
		//fn.vars.printTrace()
	}

	//go nv.run_stmt_sync(fn.body)

	//ret := nv.wait_return()

	nv.del_current_local()
	return ret

}

type PairNode struct {
	NodeType
	First  Node
	Second *PairNode
}

func (node PairNode) Copy() Node {
	copy := node.Second
	return PairNode{NodeType: node.NodeType, First: node.First.Copy(), Second: copy}
}

func (node PairNode) String() string {
	return fmt.Sprintf("(%s %s)", node.First.String(), node.Second.String())
}

func (node PairNode) Value(v *Vars) Node {
	return PairNode{NodeType: node.NodeType, First: node.First.Value(v), Second: node.Second}
}

func Begin() *Env {
	var name string
	global := Vars{
		name:  "global",
		deep:  0,
		vars:  map[IdentNode]chan Node{},
		err:   make(chan Node),
		ret:   make(chan Node),
		exit:  make(chan Node),
		next:  nil,
		debug: false,
	}
	//local := Vars{ctx: map[IdentNode]chan Node{}, next: nil}

	// BuiltIn
	name = "abs$float"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: absfloat}))
	name = "and"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: FSubr, bi: and}))
	name = "car"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: car}))
	name = "catch"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: FSubr, bi: catch}))
	name = "cdr"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: cdr}))
	name = "cond"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: FSubr, bi: cond}))
	name = "cons"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: cons}))
	name = "cos"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: cos}))
	name = "debug"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: debug_}))
	name = "def"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: FSubr, bi: def}))
	name = "div$float"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: divfloat}))
	name = "div$int"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: divint}))
	name = "eq"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: eq}))
	name = "eq$int"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: eqint}))
	name = "error"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: error}))
	name = "ete"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: ete}))
	name = "eval"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: eval}))
	name = "exit"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: exit}))
	name = "fold"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: fold}))
	name = "go-can-set"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: gocanset}))
	name = "go-elem"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: goelem}))
	name = "go-field"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: gofieldbyname}))
	name = "go-get$int"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: gogetint}))
	name = "go-get$string"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: gogetstring}))
	name = "go-get-type"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: gogettype}))
	name = "go-get-value"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: gogetvalue}))
	name = "go-indir"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: goindir}))
	name = "go-kind$type"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: gokindtype}))
	name = "go-kind$value"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: gokindvalue}))
	name = "go-new"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: gonew}))
	name = "go-set$int"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: gosetint}))
	name = "go-set$string"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: gosetstring}))
	name = "go-struct"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: gostructof}))
	name = "go-type"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: gotypeof}))
	name = "go-value"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: govalueof}))
	name = "gt$float"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: gtfloat}))
	name = "gt$int"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: gtint}))
	name = "get$json"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: getjson}))
	name = "is"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: FSubr, bi: is}))
	name = "kappa"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: FSubr, bi: kappa}))
	name = "lambda"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: FSubr, bi: lambda}))
	name = "len$vect"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: lenvect}))
	name = "let"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: FSubr, bi: let}))
	name = "let-async"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: FSubr, bi: letasync}))
	name = "lt$float"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: ltfloat}))
	name = "lt$int"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: ltint}))
	name = "map"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: fmap}))
	name = "neq"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: neq}))
	name = "not"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: not}))
	name = "omega"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: FSubr, bi: omega}))
	name = "or"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: FSubr, bi: or}))
	name = "quote"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: FSubr, bi: quote}))
	name = "print"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: print}))
	name = "prod$float"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: prodfloat}))
	name = "prod$int"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: prodint}))
	name = "remainder"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: remainder}))
	name = "reset"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: reset}))
	name = "send"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: send}))
	name = "set"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: set}))
	name = "sin"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: sin}))
	name = "sleep"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: sleep}))
	name = "start"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: start}))
	name = "sub$float"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: subfloat}))
	name = "sub$int"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: subint}))
	name = "sum$float"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: sumfloat}))
	name = "sum$int"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: sumint}))
	name = "type"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: BuiltIn, class: Subr, bi: type_}))

	// MatchBuiltIn
	name = "?"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: MatchBuiltIn, class: FSubr, mbi: m_one}))
	name = "?aut"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: MatchBuiltIn, class: FSubr, mbi: m_aut}))
	name = "?call"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: MatchBuiltIn, class: FSubr, mbi: m_call}))
	name = "?et"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: MatchBuiltIn, class: FSubr, mbi: m_et}))
	name = "?id"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: MatchBuiltIn, class: FSubr, mbi: m_id}))
	name = "?list"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: MatchBuiltIn, class: FSubr, mbi: m_list}))
	name = "?num"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: MatchBuiltIn, class: FSubr, mbi: m_num}))
	name = "?non"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: MatchBuiltIn, class: FSubr, mbi: m_non}))
	name = "?one-of"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: MatchBuiltIn, class: Subr, mbi: m_one_of}))
	name = "?pat"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: MatchBuiltIn, class: Subr, mbi: m_pat}))
	name = "?same"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: MatchBuiltIn, class: FSubr, mbi: m_same}))
	name = "?vect"
	global.set_var_chan(newIdentNode(name), makeFunc(Func{NodeType: NodeFunc, name: name, mode: MatchBuiltIn, class: FSubr, mbi: m_vect}))

	env := &Env{
		globalVars: &global,
		//localVars:  &local,
		//current: &global,
		lock: sync.RWMutex{},
	}

	return env
}

func (env *Env) Eval(args ...Node) Node {
	var ret Node

	nv := env.globalVars.new_current_local("global eval", newVectNode([]Node{}))

	go nv.run_stmt_async(args)
	ret = nv.wait_error_return()

	nv.del_current_local()

	if ret == nil {
		return newStringNode("return is nil")
	} else {
		return ret
	}
}

func (v *Vars) wait_error_return() Node {
	var ret, err Node
	ret = nil
	for ret == nil {
	Loop:
		for {
			select {
			case err = <-v.err:
				if v.debug == true {
					log.Println("wait_error_return: select err", err)
				}
				//log.Panicf("wait_error>> error: %s", err.String())
				ret = err
				break Loop
			case ret = <-v.ret:
				if v.debug == true {
					log.Println("wait_error_return: select ret", ret)
				}
				break Loop
			case ret = <-v.exit:
				if v.debug == true {
					log.Println("wait_error_return: select exit", ret)
				}
				break Loop
			case <-time.After(time.Second * 20):
				//v.printTrace()
				log.Panicf("wait_error: select timeout, deep: %d, ctx: %s", v.deep, v.name)
				break Loop
			}
		}
	}
	//if ret == nil {
	//	goto Loop
	//}
	return ret
}

/*
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
*/
