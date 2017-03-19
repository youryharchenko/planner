package pl

import "fmt"

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
	Userdef
)

type FuncClass int

const (
	FSubr FuncClass = iota
	Subr
)

type Vars struct {
	ctx  map[Word]Expression
	ret  chan Expression
	cont bool
	next *Vars
}

type Env struct {
	globalVars *Vars
	localVars  *Vars
	current    *Vars
}

type Expression interface {
	Value(env *Env) Expression
	String() string
}

type Ref struct {
	mode RefType
	ref  Word
}

func NewRef(refType RefType, word Word) Ref {
	return Ref{mode: refType, ref: word}
}

func (expr Ref) Value(env *Env) Expression {
	switch expr.mode {
	case LocalValue:
		vars := env.current
		for {
			if val, ok := vars.ctx[expr.ref]; ok {
				if val != nil {
					return val
				} else {
					fmt.Println(fmt.Sprintf("Variable %s <unassigned>", expr.ref.String()))
					return NewWord("<unassigned>")
				}
			}
			if vars.next == nil {
				fmt.Println(fmt.Sprintf("Variable %s <unbound>", expr.ref.String()))
				return NewWord("<unbound>")
			}
			vars = vars.next
		}
	}
	return NewWord("<unexpected>")
}

func (expr Ref) String() string {
	return fmt.Sprintf("%s%s", RefTypeString[expr.mode], expr.ref.String())
}

type Word struct {
	word string
}

func NewWord(word string) Word {
	return Word{word: word}
}

func (expr Word) Value(env *Env) Expression {
	return expr
}

func (expr Word) String() string {
	return expr.word
}

type Number interface {
	Float() float64
	Int() int64
}

type Int struct {
	number int64
}

type Float struct {
	number float64
}

/*
type List interface {
	Head() Expression
	Tail() List
}
*/

type Func struct {
	mode  FuncType
	class FuncClass
	bi    func(*Env, []Expression) Expression
	ud    Expression
}

func (expr Func) Value(env *Env) Expression {
	return expr
}

func (expr Func) String() string {
	return fmt.Sprintf("%v", expr.mode)
}

func Begin() *Env {
	global := Vars{ctx: map[Word]Expression{}, next: nil}
	local := Vars{ctx: map[Word]Expression{}, next: nil}

	global.ctx[NewWord("quote")] = Func{mode: BuiltIn, class: FSubr, bi: quote}
	global.ctx[NewWord("prog")] = Func{mode: BuiltIn, class: FSubr, bi: prog}
	global.ctx[NewWord("set")] = Func{mode: BuiltIn, class: Subr, bi: set}

	return &Env{globalVars: &global, localVars: &local, current: &local}
}

func (env *Env) Eval(args ...Expression) Expression {
	var ret Expression
	for _, expr := range args {
		ret = expr.Value(env)
	}
	return ret
}
