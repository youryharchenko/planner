package pl

import (
	"fmt"
	"testing"
	"time"
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
		Test{"{let [X Y Z] {set X 1} {set Y 2} {set Z 3} {fold prod$int 1 (.X .Y .Z)}}", "6"},
		Test{"{let [[X 3.7] [Y 5.4] [Z 7.2]] {fold prod$float 1 (.X .Y .Z)}}", "143.856000"},
		Test{"{let [X Y Z] {set X 1} {set Y 2} {set Z 3} {fold sub$int 6 (.X .Y .Z)}}", "0"},
		Test{"{let [[X 3.7] [Y 5.4] [Z 7.2]] {fold sub$float 100 (.X .Y .Z)}}", "83.700000"},
		Test{"{let [[f fold] [X 3.7] [Y 5.4] [Z 7.2]] {.f sub$float 100 (.X .Y .Z)}}", "83.700000"},
		Test{"{let-async [X Y Z] {set X 1} {set Y 2} {set Z 3} {fold prod$int 1 (.X .Y .Z)}}", "6"},
		Test{"{let-async [[X 3.7] [Y 5.4] [Z 7.2]] {fold prod$float 1 (.X .Y .Z)}}", "143.856000"},
		Test{"{let-async [X Y Z] {set X 1} {set Y 2} {set Z 3} {fold sub$int 6 (.X .Y .Z)}}", "0"},
		Test{"{let-async [[X 3.7] [Y 5.4] [Z 7.2]] {fold sub$float 100 (.X .Y .Z)}}", "83.700000"},
		Test{"{let-async [[f fold] [X 3.7] [Y 5.4] [Z 7.2]] {.f sub$float 100 (.X .Y .Z)}}", "83.700000"},
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
		Test{"{let [[X T] [Y ()]] {or .X .Y}}", "T"},
		Test{"{let [[X T] [Y ()]] {and .X .Y}}", "()"},
		Test{"{cond [T {print True} False True] [() {print False} True False] [T {print Else} True False Else]}", "True"},
		Test{"{let [[X T] [Y ()]] {cond [{not .X} Second First] [{not .Y} First Second]}}", "Second"},
		Test{"{let [[X T] [Y ()]] {cond [{not .X} First] [.Y Second] [() Third]}}", "()"},
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
		Test{"{remainder 10 4}", "2"},
		Test{"(+ - / * % # @ &)", "(+ - / * % # @ &)"},
		Test{"{map type (+ - / * % # @ & . ,)}", "(Id Id Id Id Id Id Id Id Id Id)"},
		Test{list_lambda, "lambda"},
		Test{"{memb & (+ - / * % # @ & . ,)}", "7"},
		Test{"{head 0 (+ - / * % # @ & . ,)}", "()"},
		Test{"{head 1 (+ - / * % # @ & . ,)}", "(+)"},
		Test{"{head 2 (+ - / * % # @ & . ,)}", "(+ -)"},
		Test{"{head 10 (+ - / * % # @ & . ,)}", "(+ - / * % # @ & . ,)"},
		Test{"{head 11 (+ - / * % # @ & . ,)}", "(+ - / * % # @ & . ,)"},
		Test{"{rest 0 (+ - / * % # @ & . ,)}", "(+ - / * % # @ & . ,)"},
		Test{"{rest 1 (+ - / * % # @ & . ,)}", "(- / * % # @ & . ,)"},
		Test{"{rest 2 (+ - / * % # @ & . ,)}", "(/ * % # @ & . ,)"},
		Test{"{rest 9 (+ - / * % # @ & . ,)}", "(,)"},
		Test{"{rest 10 (+ - / * % # @ & . ,)}", "()"},
		Test{"{rest 11 (+ - / * % # @ & . ,)}", "()"},
		Test{length_lambda, "lambda"},
		Test{poly_lambda, "lambda"},
		Test{"{trans-poly (X * Y + Z * V * W + U)}", "(+ (* X Y) (+ (* Z (* V W)) U))"},
		Test{oper_lambda, "lambda"},
		// Error
		Test{"{/ 1 1}", "1"},
		Test{"{/ 1 0}", "\"divide by 0\""},
		//Test{"{debug on}", "on"},
		Test{"{catch {/ 1 0} zerodivide}", "zerodivide"},
		Test{"{catch {/ 10 5} never}", "2"},
		// ete
		Test{"{ete {quote .x} a}", "x"},
		Test{"{ete x {quote .a}}", ".x"},
		Test{"{ete x y}", "x"},
		Test{"{ete {quote .x} {quote :a}}", ":x"},
		Test{"{ete [a b] {quote {x}}}", "{a b}"},
		Test{"{ete [a b] ()}", "(a b)"},
		Test{"{ete (a b) {quote {x}}}", "{a b}"},
		Test{"{ete (a b) []}", "[a b]"},
		Test{"{ete {quote {x y}} ()}", "(x y)"},
		Test{"{ete {quote {x y}} []}", "[x y]"},
		// is
		Test{"{is a a}", "T"},
		Test{"{is a 1}", "()"},
		Test{"{is 10.000001 10.000001}", "T"},
		Test{"{is 5 {sum$int 2 3}}", "T"},
		Test{"{is {car (2 3)} 2}", "T"},
		Test{"{let [[x (A B)]] {is .x (A B)}}", "T"},
		Test{"{let [[a 12]] {is .a 20}}", "()"},
		Test{"{let [x] {is .x (A + B)} .x}", "(A + B)"},
		Test{"{let [[x (A - B)]] {is .x (A + B)} .x}", "(A - B)"},
		Test{"{let [[x (A - B)]] {is *x (A + B)} .x}", "(A + B)"},
		Test{"{is [1 2 3] [1 2 3]}", "T"},
		Test{"{is [1 2 3 4] [1 2 3]}", "()"},
		Test{"{is [1 2 3 4] [1 2 3 3]}", "()"},
		Test{"{is [1 {car (2 3)} 3 4] [1 2 3 4]}", "T"},
		Test{"{is [1 2 3 4] [1 2 3 [4]]}", "()"},
		Test{"{is (1 2 3) (1 2 3)}", "T"},
		Test{"{is (1 2 3 4) (1 2 3)}", "()"},
		Test{"{is (1 2 3 4) (1 2 3 3)}", "()"},
		Test{"{is (1 {car (2 3)} 3 4) (1 2 3 4)}", "T"},
		Test{"{is (1 2 3 4) (1 2 3 [4])}", "()"},
		Test{"{let [x y] {is [.x *y] [(A + B) (A - B)]} .y}", "(A - B)"},
		Test{"{let [[x ()] [y ()]] {is [.x *y] [(A + B) (A - B)]} .y}", "()"},
		Test{"{let [x [y ()]] {is [.x *y] [(A + B) (A - B)]} .x}", "(A + B)"},
		Test{"{let [[x ()] [y ()]] {is [*x .y] [(A + B) (A - B)]}}", "()"},
		Test{"{let [[x ()] [y ()]] {is [*x .y] [(A + B) (A - B)]} .x}", "()"},
		Test{"{let [[x ()] [y (A - B)]] {is [*x .y] [(A + B) (A - B)]} .x}", "(A + B)"},
		Test{"{let [[x (A - B)] [y (A + B)]] {is [*x .y] [(A + B) (A - B)]} .x}", "(A - B)"},
		Test{"{let [[x ()]] {is [*x .x] [(A + B) (A + B)]} .x}", "(A + B)"},
		Test{"{let [[x ()]] {is [*x .x] [(A + B) (A - B)]} .x}", "()"},
		Test{"{is {?id} a}", "T"},
		Test{"{is {?id} 1}", "()"},
		Test{"{is {?num} 1}", "T"},
		Test{"{is {?num} 1.0}", "T"},
		Test{"{is {?list} 1.0}", "()"},
		Test{"{is {?list} ()}", "T"},
		Test{"{is {?list 0} ()}", "T"},
		Test{"{is {?list 1} ()}", "()"},
		Test{"{is {?list 5} (a b c d e)}", "T"},
		Test{"{is {?list} (a b c d e)}", "T"},
		Test{"{is {?list} [a b c d e]}", "()"},
		Test{"{is {?vect} [a b c d e]}", "T"},
		Test{"{is {?vect 5} [a b c d e]}", "T"},
		Test{"{is {?vect 5} (a b c d e)}", "()"},
		Test{"{is {?vect 4} [a b c d e]}", "()"},
		Test{"{is {?call} {quote {func c d e}}}", "T"},
		Test{"{is {?call 3} {quote {func c d e}}}", "T"},
		Test{"{is {?call 2} {quote {func c d e}}}", "()"},
		Test{"{is {?call 0} {quote {func}}}", "T"},
		Test{"{let [x] {is {?list .x} (A + B)} .x}", "3"},
		Test{"{let [x] {is {?call .x} {quote {func}}} .x}", "0"},
		Test{"{is [1 {?} 3] [1 2 3]}", "T"},
		Test{"{let [x] {is {?et {?list 3} .x} (A + B)} .x}", "(A + B)"},
		Test{"{let [[x null]] {is {?et {?list 2} .x} (A + B)} .x}", "null"},
		Test{"{is {?same [x] {?list 3} .x} (A + B)}", "T"},
		Test{"{let [x] {is {?aut (*x B) (A *x) (*x C)} (A C)} .x}", "C"},
		Test{"{let [[x null]] {is {?aut (.x B) (A .x) (.x C)} (A C)} .x}", "null"},
		Test{"{let [[x null]] {is {?aut (.x B) (A *x) (.x C)} (A C)} .x}", "C"},
		Test{"{let [[x null]] {is {?aut (.x B) (A .x) (*x C)} (A C)} .x}", "A"},
		Test{"{let [[x (A B C)]] {is {?one-of .x} A}}", "T"},
		Test{"{let [[x (A B C)]] {is {?one-of .x} D}}", "()"},
		Test{"{let [[x [1 2 3]]] {is {?one-of .x} {sum$int 1 2}}}", "T"},
		Test{"{let [y [x {quote (A .y)}]] {is {?pat .x} (A B)} .y}", "B"},
		Test{"{let [[y A] [x (A *y)]] {is {?pat .x} (A B)} .y}", "B"},
		Test{"{def ?pair {kappa [] {?vect 2}}}", "kappa"},
		Test{"{let [[x [A B]]] {is {?pair} .x}}", "T"},
		Test{"{let [[x [A B C]]] {is {?pair} .x}}", "()"},
		Test{"{def ?false {kappa [] {?list 0}}}", "kappa"},
		Test{"{let [[x ()]] {is {?false} .x}}", "T"},
		Test{"{let [[x A]] {is {?false} .x}}", "()"},
		Test{"{def ?true {kappa [] {?non {?false}}}}", "kappa"},
		Test{"{let [[x ()]] {is {?true} .x}}", "()"},
		Test{"{let [[x A]] {is {?true} .x}}", "T"},
		Test{"{def ?in {kappa [l] {?one-of .l}}}", "kappa"},
		Test{"{let [[x (A B C)]] {is {?in .x} A}}", "T"},
		Test{"{let [[x (A B C)]] {is {?in .x} D}}", "()"},
		Test{"{def ?not-in {kappa [l] {?non {?one-of .l}}}}", "kappa"},
		Test{"{let [[x (A B C)]] {is {?not-in .x} A}}", "()"},
		Test{"{let [[x (A B C)]] {is {?not-in .x} D}}", "T"},
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

var list_lambda = `
{def index
	{lambda [e items i]
		{cond
			[{not .items} ()]
			[{eq .e {car .items}} .i]
			[else  {index .e {cdr .items} {sum$int 1 .i}}]
		}
	}
}
{def beg-slice
	{lambda [n items new i]
		{cond
			[{not .items} .new]
			[{eq$int .n .i} .new]
			[else  {beg-slice .n {cdr .items} {cons {car .items} .new} {sum$int 1 .i}}]
		}
	}
}
{def end-slice
	{lambda [n items i]
		{cond
			[{not .items} ()]
			[{eq$int .n .i} .items]
			[else  {end-slice .n {cdr .items} {sum$int 1 .i}}]
		}
	}
}
{def reverse
	{lambda [items new i]
		{cond
			[{not .items} .new]
			[else  {reverse {cdr .items} {cons {car .items} .new} {sum$int 1 .i}}]
		}
	}
}
{def memb
	{lambda [e items]
		{index .e .items 0}
	}
}
{def head
	{lambda [n items]
		{reverse {beg-slice .n .items () 0} () 0}
	}
}
{def rest
	{lambda [n items]
		{end-slice .n .items 0}
	}
}
`

var poly_lambda = `
{def mono
	{lambda [m]
		{cond
			[{eq$int {length .m} 1} {car .m}]
			[else (* {car .m} {mono {rest 2 .m}})]
		}
	}
}
{def poly
	{lambda [p]
		{let [[k {memb + .p}]]
			{cond
				[{not .k} {mono .p}]
				[else (+ {mono {head .k .p}} {poly {rest {sum$int .k 1} .p}})]
			}
		}
	}
}
{def trans-poly
	{lambda [l]
		{cond
			[{eq$int {length .l} 1} .l]
			[else {poly .l}]
		}
	}
}
`
var oper_lambda = `
{def /
	{lambda [n d]
		{cond
			[{eq$int .d 0} {error "divide by 0"}]
			[else {div$int .n .d}]
		}
	}
}
`
