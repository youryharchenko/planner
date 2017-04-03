package pl

import (
	"fmt"
	"testing"
)

func TestLang(t *testing.T) {
	if res := Begin().Eval(ParseFromString("<STRING>", "atom"+"\n")...); res.String() != "atom" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "atom", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "1"+"\n")...); res.String() != "1" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "1", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "1.5"+"\n")...); res.String() != "1.5" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "1.5", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "(1 2 3)"+"\n")...); res.String() != "(1 2 3)" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "(1 2 3)", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "[1 2 3]"+"\n")...); res.String() != "[1 2 3]" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "[1 2 3]", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "{sum$int 2 3}"+"\n")...); res.String() != "5" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "5", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "{sum$int -2 3}"+"\n")...); res.String() != "1" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "1", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "{prod$int 1 2 3}"+"\n")...); res.String() != "6" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "6", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "{fold prod$int 1.0 (1 2 3 4)}"+"\n")...); res.String() != "24" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "24", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "{fold sum$float 0 (1 2 3 4)}"+"\n")...); res.String() != "10.000000" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "10.000000", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "{prog () {fold sum$float 0 (1 2 3 4)}}"+"\n")...); res.String() != "10.000000" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "10.000000", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "{prog (X Y Z) {set X 1} {set Y 2} {set Z 3} {fold prod$int 1 (.X .Y .Z)}}"+"\n")...); res.String() != "6" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "6", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "{prog ((X 3.7) (Y 5.4) (Z 7.2)) {fold prod$float 1 (.X .Y .Z)}}"+"\n")...); res.String() != "143.856000" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "6", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "{prog (X Y Z) {set X 1} {set Y 2} {set Z 3} {fold sub$int 6 (.X .Y .Z)}}"+"\n")...); res.String() != "0" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "0", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "{prog ((X 3.7) (Y 5.4) (Z 7.2)) {fold sub$float 100 (.X .Y .Z)}}"+"\n")...); res.String() != "83.700000" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "83.700000", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "{prog ((f fold) (X 3.7) (Y 5.4) (Z 7.2)) {.f sub$float 100 (.X .Y .Z)}}"+"\n")...); res.String() != "83.700000" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "83.700000", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "{def f (lambda *p {fold sum$int 0 .p})} {f 1 2 3 4 5}"+"\n")...); res.String() != "15" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "15", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "{def f (lambda *p {fold sum$int 0 .p})} {f {sum$int 1 2} 3 4 5}"+"\n")...); res.String() != "15" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "15", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "{map type ({quote{quote{print z}}} {quote{quote .z}} X (a b c) 12 \"Hello\" [1 2 3])}"+"\n")...); res.String() != "(Call Ref Id List Num Str Vect)" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "(Call Id List Num Str Vect)", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "({eq 1 2} {eq 3 3} {eq 3 3.0} {eq () ()} {eq {quote{print z}} {quote{print z}}} {eq {quote .z} {quote .z}} {eq [1 2 3] (1 2 3)})"+"\n")...); res.String() != "(() T () T T T ())" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "(() T () T T T ())", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "({neq 1 2} {neq 3 3} {neq 3 3.0} {neq () ()} {neq {quote{print z}} {quote{print z}}} {neq {quote .z} {quote .z}} {neq [1 2 3] (1 2 3)})"+"\n")...); res.String() != "(T () T () () () T)" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "(T () T () () () T)", res))
	} else {
		fmt.Printf("%v\n", res)
	}

	if res := Begin().Eval(ParseFromString("<STRING>", "({not A} {not 3} {not (3 3.0)} {not ()})"+"\n")...); res.String() != "(() () () T)" {
		t.Error(fmt.Sprintf("Expected result '%s', got string '%s'", "(() () () T)", res))
	} else {
		fmt.Printf("%v\n", res)
	}

}
