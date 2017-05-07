package pl

import (
	"fmt"
	"testing"
	"time"
)

func TestReflect(t *testing.T) {
	tests := []Test{
		// TypeOf
		Test{"{go-type 1}", "pl.NumberNode"},
		Test{"{go-type A}", "pl.IdentNode"},
		Test{"{go-type \"A\"}", "pl.StringNode"},
		Test{"{go-type ()}", "pl.ListNode"},
		Test{"{go-type []}", "pl.VectorNode"},
		Test{"{go-type {quote {f}}}", "pl.CallNode"},
		Test{"{go-type {quote .x}}", "pl.RefNode"},
		Test{"{go-type {quote :y}}", "pl.RefNode"},
		Test{"{go-type {lambda [] x}}", "pl.Func"},
		Test{"{go-type {kappa [] x}}", "pl.Func"},
		Test{"{go-type {omega [] x}}", "pl.Actor"},
		Test{"{go-kind$type {go-type 1}}", "struct"},
		Test{"{go-kind$type {go-type {omega [] x}}}", "struct"},
		Test{"{go-value 1}", "<pl.NumberNode Value>"},
		Test{"{go-value A}", "<pl.IdentNode Value>"},
		Test{"{go-value \"A\"}", "<pl.StringNode Value>"},
		Test{"{go-value ()}", "<pl.ListNode Value>"},
		Test{"{go-value []}", "<pl.VectorNode Value>"},
		Test{"{go-value {quote {f}}}", "<pl.CallNode Value>"},
		Test{"{go-value {quote .x}}", "<pl.RefNode Value>"},
		Test{"{go-value {quote :y}}", "<pl.RefNode Value>"},
		Test{"{go-value {lambda [] x}}", "<pl.Func Value>"},
		Test{"{go-value {kappa [] x}}", "<pl.Func Value>"},
		Test{"{go-value {omega [] x}}", "<pl.Actor Value>"},
		Test{"{go-kind$value {go-value 1}}", "struct"},
		Test{"{go-kind$value {go-value {omega [] x}}}", "struct"},
		Test{"{go-get-type 1}", "int64"},
		Test{"{go-get-type 1.0}", "float64"},
		Test{"{go-get-type A}", "string"},
		Test{"{go-get-type \"A\"}", "string"},
		Test{"{go-get-type ()}", "*pl.PairNode"},
		Test{"{go-get-type []}", "[]pl.Node"},
		Test{"{go-get-type {quote {f}}}", "pl.CallNode"},
		Test{"{go-get-type {quote .x}}", "string"},
		Test{"{go-get-type {quote :x}}", "string"},
		Test{"{go-get-value 1}", "<int64 Value>"},
		Test{"{go-get-value 1.0}", "<float64 Value>"},
		Test{"{go-get-value A}", "A"},
		Test{"{go-get-value \"A\"}", "\"A\""},
		Test{"{go-get-value ()}", "<*pl.PairNode Value>"},
		Test{"{go-get-value []}", "<[]pl.Node Value>"},
		Test{"{go-get-value {quote {f}}}", "<pl.CallNode Value>"},
		Test{"{go-get-value {quote .x}}", ".x"},
		Test{"{go-get-value {quote :x}}", ":x"},
		Test{"{go-kind$value {go-get-value 1}}", "int64"},
		Test{"{go-kind$value {go-get-value 1.0}}", "float64"},
		Test{"{go-kind$value {go-get-value A}}", "string"},
		Test{"{go-kind$value {go-get-value []}}", "slice"},
		Test{"{go-kind$value {go-get-value ()}}", "ptr"},
		Test{"{go-type {go-value A}}", "pl.GoValueNode"},
		Test{"{go-value {go-type A}}", "<pl.GoTypeNode Value>"},
		Test{"{go-type {go-type A}}", "pl.GoTypeNode"},
		Test{"{go-kind$type {go-get-type 1}}", "int64"},
		Test{"{go-kind$type {go-get-type 1.0}}", "float64"},
		Test{"{go-kind$type {go-get-type A}}", "string"},
		Test{"{go-kind$type {go-get-type []}}", "slice"},
		Test{"{go-kind$type {go-get-type ()}}", "ptr"},
		//Test{"", ""},
	}

	env := Begin()
	for i, test := range tests {
		//log.Println(i, test.text, "->", test.res)
		if res := env.Eval(ParseFromString("<STRING>", test.text+"\n")...); res.String() != test.res {
			t.Error(fmt.Sprintf("#%d: (%s) Expected result '%s', got string '%s'", i, test.text, test.res, res))
		} else {
			//fmt.Printf("%v\n", res)
		}
		time.Sleep(time.Millisecond)
	}
}
