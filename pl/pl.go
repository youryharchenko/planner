package pl

import (
	"fmt"
	"io"
	"log"
)

var Q = func(quietly_ignored ...interface{}) {}

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
	parser     *Parser
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

type Pair struct {
	head Expression
	tail Expression
}

func NewPair(head Expression, tail Expression) Pair {
	return Pair{head: head, tail: tail}
}

func (expr Pair) Value(env *Env) Expression {
	return expr
}

func (expr Pair) String() string {
	return expr.head.String() + " : " + expr.tail.String()
}

type Number interface {
	Float() float64
	Int() int64
}

type Int struct {
	number int64
}

func NewInt(number int64) Int {
	return Int{number: number}
}

func (expr Int) Value(env *Env) Expression {
	return expr
}

func (expr Int) String() string {
	return fmt.Sprintf("%d", expr.number)
}

type Float struct {
	number float64
}

func NewFloat(number float64) Float {
	return Float{number: number}
}

func (expr Float) Value(env *Env) Expression {
	return expr
}

func (expr Float) String() string {
	return fmt.Sprintf("%f", expr.number)
}

type Comment struct {
	text string
}

func NewComment(text string) Comment {
	return Comment{text: text}
}

func (expr Comment) Value(env *Env) Expression {
	return expr
}

func (expr Comment) String() string {
	return expr.text
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

type Sentinel struct {
	val int
}

func (expr Sentinel) Value(env *Env) Expression {
	return expr
}

func (expr Sentinel) String() string {
	return fmt.Sprintf("Sentinel:%d", expr.val)
}

// these are values now so that they also have addresses.
var ExprNull = &Sentinel{val: 0}
var ExprEnd = &Sentinel{val: 1}
var ExprMarker = &Sentinel{val: 2}

var ExprIntSize = 64
var ExprFloatSize = 64

func Begin() *Env {
	global := Vars{ctx: map[Word]Expression{}, next: nil}
	local := Vars{ctx: map[Word]Expression{}, next: nil}

	global.ctx[NewWord("quote")] = Func{mode: BuiltIn, class: FSubr, bi: quote}
	global.ctx[NewWord("prog")] = Func{mode: BuiltIn, class: FSubr, bi: prog}
	global.ctx[NewWord("set")] = Func{mode: BuiltIn, class: Subr, bi: set}
	global.ctx[NewWord("sumint")] = Func{mode: BuiltIn, class: Subr, bi: sumint}

	env := &Env{
		globalVars: &global,
		localVars:  &local,
		current:    &local,
	}

	env.parser = env.NewParser()

	return env
}

func (env *Env) Eval(args ...Expression) Expression {
	var ret Expression
	for _, expr := range args {
		ret = expr.Value(env)
	}
	return ret
}

func (env *Env) SourceExpressions(expressions []Expression) Expression {
	log.Println("SourceExpression: started")
	var result Expression
	for _, expr := range expressions {
		log.Println("Source:", expr.String())
		result = expr.Value(env)
		log.Println("Result:", result)
	}
	return result
}

func (env *Env) SourceStream(stream io.RuneScanner) Expression {
	//log.Println("SourceStream: started")
	env.parser.Start()
	env.parser.ResetAddNewInput(stream)
	expressions, err := env.parser.ParseTokens()
	if err != nil {
		return NewWord(fmt.Sprintf(
			"Error parsing on line %d: %v\n", env.parser.lexer.Linenum(), err))
	}

	return env.SourceExpressions(expressions)
}
