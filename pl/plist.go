package pl

import "fmt"

type Plist struct {
	list []Expression
}

func NewPlist(args ...Expression) Plist {
	list := []Expression{}
	for _, arg := range args {
		//fmt.Println(arg.String())
		list = append(list, arg)
	}
	return Plist{list: list}
}

func (expr Plist) Value(env *Env) Expression {
	name := expr.list[0]
	word := name.Value(env).(Word)
	vars := env.current
	var f Func
	for {
		if val, ok := vars.ctx[word]; ok {
			f = val.(Func)
			goto Apply
		}
		if vars.next == nil {
			break
		}
		vars = vars.next
	}
	if val, ok := env.globalVars.ctx[word]; ok {
		f = val.(Func)
	} else {
		fmt.Println(fmt.Sprintf("Function %s <unbound>", word.String()))
		return NewWord("<unbound>")
	}

Apply:
	switch f.mode {
	case BuiltIn:
		var list []Expression
		if f.class == FSubr {
			list = expr.list[1:]
		} else {
			list = []Expression{}
			for _, elm := range expr.list[1:] {
				list = append(list, elm.Value(env))
			}
		}
		return f.bi(env, list)
	}
	return NewWord("<unexpected>")
}

func (expr Plist) String() string {
	ret := "["
	for _, elm := range expr.list {
		ret += elm.String() + " "
	}
	return ret + "]"
}
