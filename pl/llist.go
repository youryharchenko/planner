package pl

type Llist struct {
	list []Expression
}

func NewLlist(args ...Expression) Llist {
	list := []Expression{}
	for _, arg := range args {
		//fmt.Println(arg.String())
		list = append(list, arg)
	}
	return Llist{list: list}
}

func (expr Llist) Value(env *Env) Expression {
	list := []Expression{}
	for _, elm := range expr.list {
		//fmt.Println(arg.String())
		list = append(list, elm.Value(env))
	}
	return Llist{list: list}
}

func (expr Llist) String() string {
	ret := "("
	for _, elm := range expr.list {
		ret += elm.String() + " "
	}
	return ret + ")"
}
