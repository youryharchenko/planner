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
		Test{"{def square (lambda (x) {prod$int .x .x})}", "(lambda (x) {prod$int .x .x})"},
		Test{"{square 21}", "441"},
		Test{"{square {sum$int 2 5}}", "49"},
		Test{"{square {square 3}}", "81"},
		Test{"{def sum-of-squares (lambda (x y) {sum$int {square .x} {square .y}})}", "(lambda (x y) {sum$int {square .x} {square .y}})"},
		Test{"{sum-of-squares 3 4}", "25"},
		Test{"{def f (lambda (a) {sum-of-squares {sum$int .a 1} {prod$int .a 2}})}{f 5}", "136"},
		Test{"{def if (lambda (c *t *e) {cond (.c {eval .t}) (T {eval .e})})}", "(lambda (c *t *e) {cond (.c {eval .t}) (T {eval .e})})"},
		// Calc SQRT
		Test{"{def square$float (lambda (x) {prod$float .x .x})}", "(lambda (x) {prod$float .x .x})"},
		Test{"{def sqrt-iter (lambda (guess x) {if {good-enough .guess .x} .guess {sqrt-iter {improve-guess .guess .x} .x}})}", "(lambda (guess x) {if {good-enough .guess .x} .guess {sqrt-iter {improve-guess .guess .x} .x}})"},
		Test{"{def improve-guess (lambda (guess x) {average .guess {div$float .x .guess}})}", "(lambda (guess x) {average .guess {div$float .x .guess}})"},
		Test{"{def average (lambda (x y) {div$float {sum$float .x .y} 2.0})}", "(lambda (x y) {div$float {sum$float .x .y} 2.0})"},
		Test{"{def good-enough (lambda (guess x) {lt$float {abs$float {sub$float {square$float .guess} .x}} 0.001})}", "(lambda (guess x) {lt$float {abs$float {sub$float {square$float .guess} .x}} 0.001})"},
		Test{"{def sqrt$float (lambda (x) {sqrt-iter 1.0 .x})}", "(lambda (x) {sqrt-iter 1.0 .x})"},
		Test{"{sqrt$float 9.0}", "3.000092"},
		Test{"{sqrt$float 50.0}", "7.071068"},
		//
		Test{"{map sqrt$float {map square$float (1 2 3 4 5 6 7 8 9)}}", "(1.0 2.000000 3.000092 4.000001 5.000023 6.000000 7.000000 8.000002 9.000011)"},
		// Calc factorial - linear recursive process
		Test{"{def factor-lrp (lambda (n) {if {lt$int .n 2} 1 {prod$int .n {factor-lrp {sub$int .n 1}}}})}", "(lambda (n) {if {lt$int .n 2} 1 {prod$int .n {factor-lrp {sub$int .n 1}}}})"},
		Test{"{factor-lrp 6}", "720"},
		// Calc factorial - linear iterative process
		Test{"{def factor-iter (lambda (product counter max-count) {if {gt$int .counter .max-count} .product {factor-iter {prod$int .counter .product} {sum$int .counter 1} .max-count}})}", "(lambda (product counter max-count) {if {gt$int .counter .max-count} .product {factor-iter {prod$int .counter .product} {sum$int .counter 1} .max-count}})"},
		Test{"{def factor-lip (lambda (n) {factor-iter 1 1 .n})}", "(lambda (n) {factor-iter 1 1 .n})"},
		Test{"{factor-lip 6}", "720"},
		// Calc Fibonacci numbers iteratively
		Test{"{def fib (lambda (n) {fib-iter 1 0 .n})}", "(lambda (n) {fib-iter 1 0 .n})"},
		Test{"{def fib-iter (lambda (a b count) {if {lt$int .count 1} .b {fib-iter {sum$int .a .b} .a {sub$int .count 1}}})}", "(lambda (a b count) {if {lt$int .count 1} .b {fib-iter {sum$int .a .b} .a {sub$int .count 1}}})"},
		Test{"{fib 8}", "21"},
		// Example: Counting change
		Test{"{def count-change (lambda (amount) {cc .amount 5})}", "(lambda (amount) {cc .amount 5})"},
		Test{"{def cc (lambda (amount kinds-of-coins) {cond ({eq$int .amount 0} 1) ({or {lt$int .amount 0} {eq$int .kinds-of-coins 0}} 0) (else {sum$int {cc .amount {sub$int .kinds-of-coins 1}} {cc {sub$int .amount {first-denomination .kinds-of-coins}} .kinds-of-coins}})})}", "(lambda (amount kinds-of-coins) {cond ({eq$int .amount 0} 1) ({or {lt$int .amount 0} {eq$int .kinds-of-coins 0}} 0) (else {sum$int {cc .amount {sub$int .kinds-of-coins 1}} {cc {sub$int .amount {first-denomination .kinds-of-coins}} .kinds-of-coins}})})"},
		Test{"{def first-denomination (lambda (kinds-of-coins) {cond ({eq$int .kinds-of-coins 1} 1) ({eq$int .kinds-of-coins 2} 5) ({eq$int .kinds-of-coins 3} 10) ({eq$int .kinds-of-coins 4} 25) ({eq$int .kinds-of-coins 5} 50)})}", "(lambda (kinds-of-coins) {cond ({eq$int .kinds-of-coins 1} 1) ({eq$int .kinds-of-coins 2} 5) ({eq$int .kinds-of-coins 3} 10) ({eq$int .kinds-of-coins 4} 25) ({eq$int .kinds-of-coins 5} 50)})"},
		Test{"{count-change 1}", "1"},
		Test{"{count-change 100}", "292"},
		// Procedures as Arguments
		Test{"{def sum (lambda (term a next b) {if {gt$int .a .b} 0 {sum$int {term .a} {sum .term {next .a} .next .b}}})}", "(lambda (term a next b) {if {gt$int .a .b} 0 {sum$int {term .a} {sum .term {next .a} .next .b}}})"},
		Test{"{def inc (lambda (n) {sum$int .n 1})}", "(lambda (n) {sum$int .n 1})"},
		Test{"{def identity (lambda (x) .x)}", "(lambda (x) .x)"},
		Test{"{def cube (lambda (x) {prod$int .x .x .x})}", "(lambda (x) {prod$int .x .x .x})"},
		Test{"{def sum-int (lambda (a b) {sum identity .a inc .b})}", "(lambda (a b) {sum identity .a inc .b})"},
		Test{"{def sum-cube (lambda (a b) {sum cube .a inc .b})}", "(lambda (a b) {sum cube .a inc .b})"},
		Test{"{sum-int 1 10}", "55"},
		Test{"{sum-cube 1 10}", "3025"},
		Test{"{def sumf (lambda (term a next b) {if {gt$float .a .b} 0 {sum$float {term .a} {sumf .term {next .a} .next .b}}})}", "(lambda (term a next b) {if {gt$float .a .b} 0 {sum$float {term .a} {sumf .term {next .a} .next .b}}})"},
		Test{"{def pi-sum (lambda (a b) {def pi-term (lambda (x) {div$float 1.0 {prod$float .x {sum$float .x 2}}})} {def pi-next (lambda (x) {sum$float .x 4})} {sumf pi-term .a pi-next .b})}", "(lambda (a b) {def pi-term (lambda (x) {div$float 1.0 {prod$float .x {sum$float .x 2}}})} {def pi-next (lambda (x) {sum$float .x 4})} {sumf pi-term .a pi-next .b})"},
		Test{"{prod$float 8 {pi-sum 1 1000}}", "3.139593"},
		Test{"{def integral (lambda (f a b dx) {def add-dx (lambda (x) {sum$float .x .dx})} {prod$float {sumf .f {sum$float .a {div$float .dx 2.0}} add-dx .b} .dx})}", "(lambda (f a b dx) {def add-dx (lambda (x) {sum$float .x .dx})} {prod$float {sumf .f {sum$float .a {div$float .dx 2.0}} add-dx .b} .dx})"},
		Test{"{integral cube 0 1 0.01}", "0.500000"},
		//
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
