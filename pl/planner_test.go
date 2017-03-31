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

}
