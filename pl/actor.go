package pl

import (
	"fmt"
	"log"
	"time"
)

type Actor struct {
	NodeType
	name string
	body *Omega
}

func (expr Actor) Value(v *Vars) Node {
	return expr
}

func (expr Actor) String() string {
	return fmt.Sprintf("%v", expr.name)
}

func (expr Actor) Copy() Node {
	return expr
}

type ActorInst struct {
	NodeType
	name  string
	actor *Actor
	in    chan Node
}

func (expr ActorInst) Value(v *Vars) Node {
	return expr
}

func (expr ActorInst) String() string {
	return fmt.Sprintf("%v", expr.name)
}

func (expr ActorInst) Copy() Node {
	return expr
}

type Omega struct {
	vars *Vars
	arg  Node
	body []Node
}

func (om *Omega) apply(name string, args []Node, v *Vars, me *ActorInst) bool {
	//log.Println("apply omega: args", args)
	var vars VectorNode

	switch om.arg.Type() {
	case NodeIdent:
		var ident IdentNode
		var param Node
		arg := om.arg.(IdentNode)

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
		lst := om.arg.(VectorNode)
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

	vars = newVectNode(append(vars.Nodes, newVectNode([]Node{newIdentNode("me"), *me})))
	//log.Println(name, vars, v.deep, v.name)
	//nv := fn.vars.new_current_local(name, vars)
	onv := v.new_current_local(name+"-def", newVectNode([]Node{})).merge(om.vars)
	//log.Println(fnv)
	go onv.wait_return()

	nv := onv.new_current_local(name+"-run", vars)

	if v.debug {
		//v.printTrace()
		//log.Println("=============")
		//fn.vars.printTrace()
	}

	go nv.run_stmt_sync(om.body)

	ret := nv.wait_omega_return()

	nv.del_current_local()

	return ret
}

func (v *Vars) wait_omega_return() bool {
	//log.Println("wait_omega_return started")
	ret := v.wait_return()
	//log.Println("wait_omega_return finished", ret.String())
	if ret.String() == "()" {
		return false
	} else {
		return true
	}
}

func (inst *ActorInst) stop() {
	inst.in <- newIdentNode("stop")
}

func (v *Vars) run_actor(inst *ActorInst) {
	//log.Printf("actor %s started", inst.name)
	//defer log.Printf("actor %s finished", inst.name)
Loop:
	for {
		select {
		case rec := <-inst.in:
			//log.Printf("actor %s recives message %s", inst.name, rec)
			switch {
			case rec.Type() == NodeIdent && rec.String() == "stop":
				//log.Printf("actor %s will stop by message", inst.name)
				break Loop
			default:
				if !inst.actor.body.apply("omega", []Node{rec}, v, inst) {
					//log.Printf("actor %s will stop by return", inst.name)
					break Loop
				}
			}
		}
	}

}

func makeOmega(name string, v *Vars, arg Node, body []Node) Actor {
	return Actor{NodeType: NodeActor, name: name, body: &Omega{vars: v, arg: arg, body: body}}
}

func newActorInst(actor Actor, name string) ActorInst {
	inst := ActorInst{NodeType: NodeActorInst, actor: &actor, name: name, in: make(chan Node, 10)}
	return inst
}

func findActor(id IdentNode, v *Vars) Actor {

	for i := 0; i < 3; i++ {
		vars := v
		for {
			if ch := vars.get_var_chan(id); ch != nil {
				var val Node
				select {
				case val = <-ch:
					ch <- val
				case <-time.After(time.Second * 5):
					log.Panicf("find actor timeout: %s, deep: %d, ctx: %s", id.String(), v.deep, v.name)
				}
				//log.Println("Function found", word, val)
				switch val.Type() {
				case NodeActor:
					a := val.(Actor)
					return a
				case NodeIdent:
					return findActor(val.(IdentNode), v)
				default:
					log.Panicf("findActor>> unexpected type, name:%s, type: %s, value: %s", id.String(), type_(v, []Node{val}), val.String())
				}

			} //else {
			//log.Panicf("variable %s <unassigned>, deep: %d, ctx: %s", word.String(), v.deep, v.name)
			//}

			if vars.next == nil {
				break
			}
			nvars := vars.next
			vars = nvars
		}
		time.Sleep(time.Millisecond * 1)
		v.printTrace()
		log.Printf("warning, wait actor %s, deep: %d, ctx: %s", id.String(), v.deep, v.name)
	}
	v.printTrace()
	log.Panicf("variable %s <unbound>, deep: %d, ctx: %s", id.String(), v.deep, v.name)
	return Actor{}
}
