package pl

import (
	"fmt"
	"log"
	"testing"
)

type Test struct {
	text, res string
}

func TestSICP(t *testing.T) {
	tests := []Test{
		Test{"486", "486"},
		Test{"{sum$int 137 349}", "486"},
		Test{"{sub$int 1000 334}", "666"},
		Test{"{prod$int 5 99}", "495"},
		Test{"{div$int 10 2}", "5"},
		Test{"{sum$float 2.7 10}", "12.700000"},
		Test{"{sum$int 21 35 12 7}", "75"},
		Test{"{prod$int 25 4 12}", "1200"},
		Test{"{sum$int {prod$int 3 5} {sub$int 10 6}}", "19"},
		Test{"{sum$int {prod$int 3 {sum$int {prod$int 2 4} {sum$int 3 5}}} {sum$int {sub$int 10 7} 6}}", "57"},
		Test{"{def size 2}", "2"},
		Test{":size", "2"},
		Test{"{prod$int 5 :size}", "10"},
		Test{"{def pi 3.14159}", "3.14159"},
		Test{"{def radius 10}", "10"},
		Test{"{prod$float :pi :radius :radius}", "314.159000"},
		Test{"{def circum {prod$float 2 :pi :radius}}", "62.831800"},
		Test{":circum", "62.831800"},
		Test{fmt.Sprintf("{def square %s}", lambda_square), "lambda"},
		Test{"{square 21}", "441"},
		Test{"{square {sum$int 2 5}}", "49"},
		Test{"{square {square 3}}", "81"},
		Test{fmt.Sprintf("{def sum-of-squares %s}", lambda_sum_of_squares), "lambda"},
		Test{"{sum-of-squares 3 4}", "25"},
		Test{fmt.Sprintf("{%s 5}", lambda_f), "136"},
		Test{"{def if {lambda [c *t *e] {cond [.c {eval .t}] [T {eval .e}]}}}", "lambda"},
		// Calc SQRT
		Test{"{def square$float {lambda [x] {prod$float .x .x}}}", "lambda"},
		Test{"{def sqrt-iter {lambda [guess x] {if {good-enough .guess .x} .guess {sqrt-iter {improve-guess .guess .x} .x}}}}", "lambda"},
		Test{"{def improve-guess {lambda [guess x] {average .guess {div$float .x .guess}}}}", "lambda"},
		Test{"{def average {lambda [x y] {div$float {sum$float .x .y} 2.0}}}", "lambda"},
		Test{"{def good-enough {lambda [guess x] {lt$float {abs$float {sub$float {square$float .guess} .x}} 0.001}}}", "lambda"},
		Test{"{def sqrt$float {lambda [x] {sqrt-iter 1.0 .x}}}", "lambda"},
		Test{"{sqrt$float 9.0}", "3.000092"},
		Test{"{sqrt$float 50.0}", "7.071068"},
		//
		Test{"{map sqrt$float {map square$float (1 2 3 4 5 6 7 8 9)}}", "(1.0 2.000000 3.000092 4.000001 5.000023 6.000000 7.000000 8.000002 9.000011)"},
		// Calc factorial - linear recursive process
		Test{"{def factor-lrp {lambda [n] {if {lt$int .n 2} 1 {prod$int .n {factor-lrp {sub$int .n 1}}}}}}", "lambda"},
		Test{"{factor-lrp 6}", "720"},
		// Calc factorial - linear iterative process
		Test{"{def factor-iter {lambda [product counter max-count] {if {gt$int .counter .max-count} .product {factor-iter {prod$int .counter .product} {sum$int .counter 1} .max-count}}}}", "lambda"},
		Test{"{def factor-lip {lambda [n] {factor-iter 1 1 .n}}}", "lambda"},
		Test{"{factor-lip 6}", "720"},
		// Calc Fibonacci numbers iteratively
		Test{"{def fib {lambda [n] {fib-iter 1 0 .n}}}", "lambda"},
		Test{"{def fib-iter {lambda [a b count] {if {lt$int .count 1} .b {fib-iter {sum$int .a .b} .a {sub$int .count 1}}}}}", "lambda"},
		Test{"{fib 8}", "21"},
		// Example: Counting change
		Test{"{def count-change {lambda [amount] {cc .amount 5}}}", "lambda"},
		Test{"{def cc {lambda [amount kinds-of-coins] {cond [{eq$int .amount 0} 1] [{or {lt$int .amount 0} {eq$int .kinds-of-coins 0}} 0] [else {sum$int {cc .amount {sub$int .kinds-of-coins 1}} {cc {sub$int .amount {first-denomination .kinds-of-coins}} .kinds-of-coins}}]}}}", "lambda"},
		Test{"{def first-denomination {lambda [kinds-of-coins] {cond [{eq$int .kinds-of-coins 1} 1] [{eq$int .kinds-of-coins 2} 5] [{eq$int .kinds-of-coins 3} 10] [{eq$int .kinds-of-coins 4} 25] [{eq$int .kinds-of-coins 5} 50]}}}", "lambda"},
		Test{"{count-change 1}", "1"},
		Test{"{count-change 100}", "292"},
		// Procedures as Arguments
		Test{"{def sum {lambda [term a next b] {if {gt$int .a .b} 0 {sum$int {term .a} {sum .term {next .a} .next .b}}}}}", "lambda"},
		Test{"{def inc {lambda [n] {sum$int .n 1}}}", "lambda"},
		Test{"{def identity {lambda [x] .x}}", "lambda"},
		Test{"{def cube {lambda [x] {prod$int .x .x .x}}}", "lambda"},
		Test{"{def sum-int {lambda [a b] {sum identity .a inc .b}}}", "lambda"},
		Test{"{def sum-cube {lambda [a b] {sum cube .a inc .b}}}", "lambda"},
		Test{"{sum-int 1 10}", "55"},
		Test{"{sum-cube 1 10}", "3025"},
		Test{"{def sumf {lambda [term a next b] {if {gt$float .a .b} 0 {sum$float {term .a} {sumf .term {next .a} .next .b}}}}}", "lambda"},
		Test{"{def pi-sum {lambda [a b] {def pi-term {lambda [x] {div$float 1.0 {prod$float .x {sum$float .x 2}}}}} {def pi-next {lambda [x] {sum$float .x 4}}} {sumf pi-term .a pi-next .b}}}", "lambda"},
		Test{"{prod$float 8 {pi-sum 1 1000}}", "3.139593"},
		Test{"{def integral {lambda [f a b dx] {def add-dx {lambda [x] {sum$float .x .dx}}} {prod$float {sumf .f {sum$float .a {div$float .dx 2.0}} add-dx .b} .dx}}}", "lambda"},
		Test{"{integral cube 0 1 0.01}", "0.500000"},
		Test{fmt.Sprintf("{prod$float 8 {%s 1 1000}}", lambda_pi_sum), "3.139593"},
		// let
		Test{"{let [[x 5]] {sum$int {let [[x 3]] {sum$int .x {prod$int .x 10}}} .x}}", "38"},
		Test{"{let [[x 2]] {let [[x 3] [y {sum$int .x 2}]] {prod$int .x .y}}}", "12"},
		// Finding roots of equations by the half-interval method
		Test{search_lambda, "lambda"},
		Test{close_enough_lambda, "lambda"},
		Test{half_interval_method_lambda, "lambda"},
		Test{"{def negative {lambda [x] {lt$float .x 0.0}}}", "lambda"},
		Test{"{def positive {lambda [x] {gt$float .x 0.0}}}", "lambda"},
		Test{"{half-interval-method sin 2.0 4.0}", "3.141113"},
		Test{"{half-interval-method {lambda [x] {sub$float {prod$float .x .x .x} {prod$float 2.0 .x} 3.0}} 1.0 2.0}", "1.893066"},
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

var lambda_sum_of_squares = `
	{lambda [x y]
		{sum$int
			{square .x}
			{square .y}
		}
	}
`
var lambda_square = `
	{lambda [x]
		{prod$int .x .x}
	}
`
var lambda_f = `
	{lambda [a]
		{sum-of-squares
			{sum$int .a 1}
			{prod$int .a 2}
		}
	}
`
var lambda_pi_sum = `
	{lambda [a b]
		{sumf
			{lambda [x]
				{div$float
					1.0
					{prod$float
						.x
						{sum$float
							.x
							2.0
						}
					}
				}
			}
			.a
			{lambda [x]
				{sum$float
					.x
					4.0
				}
			}
			.b
		}
	}
`
var search_lambda = `
{def search
	{lambda [f neg-point pos-point]
		{let [[midpoint {average .neg-point .pos-point}]]
			{if {close-enough .neg-point .pos-point}
				.midpoint
				{let [[test-value {f .midpoint}]]
					{cond
						[{positive .test-value} {search .f .neg-point .midpoint}]
						[{negative .test-value} {search .f .midpoint .pos-point}]
						[else .midpoint]
					}
				}
			}
		}
	}
}
`
var close_enough_lambda = `
{def close-enough
	{lambda [x y]
		{lt$float {abs$float {sub$float .x .y}} 0.001}
	}
}
`

var half_interval_method_lambda = `
{def half-interval-method
	{lambda [f a b]
		{let [[a-value {.f .a}] [b-value {.f .b}]]
			{cond
				[{and {negative .a-value} {positive .b-value}} {search .f .a .b}]
				[{and {negative .b-value} {positive .a-value}} {search .f .b .a}]
				[else {print "Values are not of opposite sign" a b}]
			}
		}
	}
}
`
