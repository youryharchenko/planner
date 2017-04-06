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
