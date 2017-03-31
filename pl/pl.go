package pl

import (
	"fmt"
	"log"
	"time"
)

//var Q = func(quietly_ignored ...interface{}) {}

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
	ctx  map[IdentNode]chan Node
	ret  chan Node
	exit chan Node
	cont bool
	next *Vars
}

type Env struct {
	//parser     *Parser
	globalVars *Vars
	localVars  *Vars
	current    *Vars
}

/*
type Node interface {
	Value(env *Env) Node
	String() string
}
*/

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
			vars = vars.next
		}
	case GlobalValue:
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

/*
type Word struct {
	word string
}

func NewWord(word string) Word {
	return Word{word: word}
}

func (expr Word) Value(env *Env) Node {
	return expr
}

func (expr Word) String() string {
	return expr.word
}

type Pair struct {
	head Node
	tail Node
}

func NewPair(head Node, tail Node) Pair {
	return Pair{head: head, tail: tail}
}

func (expr Pair) Value(env *Env) Node {
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

func (expr Int) Value(env *Env) Node {
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

func (expr Float) Value(env *Env) Node {
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

func (expr Comment) Value(env *Env) Node {
	return expr
}

func (expr Comment) String() string {
	return expr.text
}


type List interface {
	Head() Node
	Tail() List
}
*/

type Func struct {
	NodeType
	mode  FuncType
	class FuncClass
	bi    func(*Env, []Node) Node
	ud    Node
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

/*
type Sentinel struct {
	val int
}

func (expr Sentinel) Value(env *Env) Node {
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
*/

func Begin() *Env {
	global := Vars{ctx: map[IdentNode]chan Node{}, next: nil}
	local := Vars{ctx: map[IdentNode]chan Node{}, next: nil}

	global.ctx[newIdentNode("def")] = makeVar(Func{mode: BuiltIn, class: FSubr, bi: def})
	global.ctx[newIdentNode("div$float")] = makeVar(Func{mode: BuiltIn, class: Subr, bi: divfloat})
	global.ctx[newIdentNode("div$int")] = makeVar(Func{mode: BuiltIn, class: Subr, bi: divint})
	global.ctx[newIdentNode("exit")] = makeVar(Func{mode: BuiltIn, class: Subr, bi: exit})
	global.ctx[newIdentNode("fold")] = makeVar(Func{mode: BuiltIn, class: Subr, bi: fold})
	global.ctx[newIdentNode("map")] = makeVar(Func{mode: BuiltIn, class: Subr, bi: fmap})
	global.ctx[newIdentNode("quote")] = makeVar(Func{mode: BuiltIn, class: FSubr, bi: quote})
	global.ctx[newIdentNode("print")] = makeVar(Func{mode: BuiltIn, class: Subr, bi: print})
	global.ctx[newIdentNode("prod$float")] = makeVar(Func{mode: BuiltIn, class: Subr, bi: prodfloat})
	global.ctx[newIdentNode("prod$int")] = makeVar(Func{mode: BuiltIn, class: Subr, bi: prodint})
	global.ctx[newIdentNode("prog")] = makeVar(Func{mode: BuiltIn, class: FSubr, bi: prog})
	global.ctx[newIdentNode("set")] = makeVar(Func{mode: BuiltIn, class: Subr, bi: set})
	global.ctx[newIdentNode("sub$float")] = makeVar(Func{mode: BuiltIn, class: Subr, bi: subfloat})
	global.ctx[newIdentNode("sub$int")] = makeVar(Func{mode: BuiltIn, class: Subr, bi: subint})
	global.ctx[newIdentNode("sum$float")] = makeVar(Func{mode: BuiltIn, class: Subr, bi: sumfloat})
	global.ctx[newIdentNode("sum$int")] = makeVar(Func{mode: BuiltIn, class: Subr, bi: sumint})

	env := &Env{
		globalVars: &global,
		localVars:  &local,
		current:    &local,
	}

	//env.parser = env.NewParser()

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

//func (env *Env) SourceStream(stream io.RuneScanner) Node {
//log.Println("SourceStream: started")
//env.parser.Start()
//env.parser.ResetAddNewInput(stream)
//Nodes, err := env.parser.ParseTokens()
//if err != nil {
//	return NewWord(fmt.Sprintf(
//		"Error parsing on line %d: %v\n", env.parser.lexer.Linenum(), err))
//}

//return env.SourceNodes(Nodes)
//}
