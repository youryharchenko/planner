package pl

import (
	"fmt"
	"log"
	"testing"
)

func TestLang(t *testing.T) {
	tests := []Test{
		Test{"atom", "atom"},
		Test{"1", "1"},
		Test{"1.5", "1.5"},
		Test{"(1 2 3)", "(1 2 3)"},
		Test{"[1 2 3]", "[1 2 3]"},
		Test{"{sum$int 2 3}", "5"},
		Test{"{sum$int -2 3}", "1"},
		Test{"{prod$int 1 2 3}", "6"},
		Test{"{fold prod$int 1.0 (1 2 3 4)}", "24"},
		Test{"{fold sum$float 0 (1 2 3 4)}", "10.000000"},
		Test{"{prog [X Y Z] {set X 1} {set Y 2} {set Z 3} {fold prod$int 1 (.X .Y .Z)}}", "6"},
		Test{"{prog [[X 3.7] [Y 5.4] [Z 7.2]] {fold prod$float 1 (.X .Y .Z)}}", "143.856000"},
		Test{"{prog [X Y Z] {set X 1} {set Y 2} {set Z 3} {fold sub$int 6 (.X .Y .Z)}}", "0"},
		Test{"{prog [[X 3.7] [Y 5.4] [Z 7.2]] {fold sub$float 100 (.X .Y .Z)}}", "83.700000"},
		Test{"{prog [[f fold] [X 3.7] [Y 5.4] [Z 7.2]] {.f sub$float 100 (.X .Y .Z)}}", "83.700000"},
		Test{"{def f1 {lambda *p {fold sum$int 0 .p}}} {f1 1 2 3 4 5}", "15"},
		Test{"{def f2 {lambda *p {fold sum$int 0 .p}}}", "lambda"},
		Test{"{f2 {sum$int 1 2} 3 4 5}", "15"},
		Test{"{map type ({quote{quote{print z}}} {quote{quote .z}} X (a b c) 12 \"Hello\" [1 2 3])}", "(Call Ref Id List Num Str Vect)"},
		Test{"({eq 1 2} {eq 3 3} {eq 3 3.0} {eq () ()} {eq {quote{print z}} {quote{print z}}} {eq {quote .z} {quote .z}} {eq [1 2 3] (1 2 3)})", "(() T () T T T ())"},
		Test{"({neq 1 2} {neq 3 3} {neq 3 3.0} {neq () ()} {neq {quote{print z}} {quote{print z}}} {neq {quote .z} {quote .z}} {neq [1 2 3] (1 2 3)})", "(T () T () () () T)"},
		Test{"({not A} {not 3} {not (3 3.0)} {not ()})", "(() () () T)"},
		Test{"{or () () True () ()}", "True"},
		Test{"{and () () True () ()}", "()"},
		Test{"{or () () () ()}", "()"},
		Test{"{and A B C D E}", "E"},
		Test{"{prog [[X T] [Y ()]] {or .X .Y}}", "T"},
		Test{"{prog [[X T] [Y ()]] {and .X .Y}}", "()"},
		Test{"{cond [T {print True} False True] [() {print False} True False] [T {print Else} True False Else]}", "True"},
		Test{"{prog [[X T] [Y ()]] {cond [{not .X} Second First] [{not .Y} First Second]}}", "Second"},
		Test{"{prog [[X T] [Y ()]] {cond [{not .X} First] [.Y Second] [() Third]}}", "()"},
		Test{"{lt$float 9.0 9.0}", "()"},
		Test{"{lt$float 5.0 9.0}", "T"},
		Test{"{lt$float 9.0 5.0}", "()"},
		Test{"{gt$float 9.0 9.0}", "()"},
		Test{"{gt$float 5.0 9.0}", "()"},
		Test{"{gt$float 9.0 5.0}", "T"},
		Test{"{def if {lambda [c *t *e] {cond [.c {eval .t}] [else {eval .e}]}}}", "lambda"},
		Test{"{if {eq 1 1} {div$int 1 1} {div$int 1 0}}", "1"},
		Test{"{if {eq 1 0} {div$int 1 0} {div$int 1 1}}", "1"},
		Test{"{def square$float {lambda [x] {prod$float .x .x}}}", "lambda"},
		Test{"{abs$float {sub$float {square$float 3.0} 9.0}}", "0.000000"},
		Test{"{lt$float {abs$float {sub$float {square$float 3.0} 9.0}} 0.001}", "T"},
		Test{"{if {lt$int 0 1} {div$int 1 1} {div$int 1 0}}", "1"},
		Test{"{if {lt$float 0.001 0.002} {div$float 1 1} {div$float 1 0}}", "1.000000"},
		Test{"{if {lt$int 1 0} {div$int 1 0} {div$int 1 1}}", "1"},
		Test{"{if {lt$float 0.002 0.001} {div$float 1 0} {div$float 1 1}}", "1.000000"},
		Test{"{def cloj {lambda [p] {def fn {lambda [] .p}} {fn}}}", "lambda"},
		Test{"{cloj Hello}", "Hello"},
		//Test{"", ""},
	}

	env := Begin()
	for i, test := range tests {
		log.Println(i, test.text, "->", test.res)
		if res := env.Eval(ParseFromString("<STRING>", test.text+"\n")...); res.String() != test.res {
			t.Error(fmt.Sprintf("#%d: Expected result '%s', got string '%s'", i, test.res, res))
		} else {
			fmt.Printf("%v\n", res)
		}
	}
}
