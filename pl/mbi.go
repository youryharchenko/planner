package pl

func m_aut(v *Vars, args []Node, expr Node) bool {
	for _, n := range args {
		if v.run_is(n, expr) {
			return true
		}
	}
	return false
}

func m_call(v *Vars, args []Node, expr Node) bool {
	if expr.Type() == NodeCall {
		if len(args) == 0 {
			return true
		} else {
			return v.run_is(args[0], newInt(int64(len(expr.(CallNode).Args))))
		}
	}
	return false
}

func m_et(v *Vars, args []Node, expr Node) bool {
	for _, n := range args {
		if !v.run_is(n, expr) {
			return false
		}
	}
	return true
}

func m_id(v *Vars, args []Node, expr Node) bool {
	if expr.Type() == NodeIdent {
		return true
	}
	return false
}

func m_list(v *Vars, args []Node, expr Node) bool {
	if expr.Type() == NodeList {
		if len(args) == 0 {
			return true
		} else {
			return v.run_is(args[0], newInt(expr.(ListNode).Len()))
		}
	}
	return false
}

func m_num(v *Vars, args []Node, expr Node) bool {
	if expr.Type() == NodeNumber {
		return true
	}
	return false
}

func m_one(v *Vars, args []Node, expr Node) bool {
	return true
}

func m_same(v *Vars, args []Node, expr Node) bool {
	ret := true
	vars := args[0].(VectorNode)
	nvars := make([]Node, len(vars.Nodes))

	for i, n := range vars.Nodes {
		nvars[i] = n.Value(v)
	}
	nv := v.new_current_local("m_same", newVectNode(nvars))
	nv.rb = map[IdentNode]Node{}

	for _, n := range args[1:] {
		if !nv.run_is(n, expr) {
			ret = false
			break
		}
	}

	nv.del_current_local()
	return ret
}

func m_vect(v *Vars, args []Node, expr Node) bool {
	if expr.Type() == NodeVector {
		if len(args) == 0 {
			return true
		} else {
			return v.run_is(args[0], newInt(int64(len(expr.(VectorNode).Nodes))))
		}
	}
	return false
}
