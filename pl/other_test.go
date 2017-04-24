package pl

import (
	"fmt"
	"strings"
	"testing"
)

func TestOther(t *testing.T) {
	tests := []Test{
		Test{fmt.Sprintf("{def json1 %s}", json1_string), strings.Replace(json1_string, "\\\"", "\"", -1)},
		Test{"{get$json .json1}", "([a 1.000000] [b \"test\"] [c T] [d (1.000000 2.000000 3.000000)] [e ([x 10.300000] [y 25.800000] [z (T ())])] [f ()] [g ()])"},
		//Test{"", ""},

	}

	env := Begin()
	for i, test := range tests {
		//log.Println(i, test.text, "->", test.res)
		if res := env.Eval(ParseFromString("<STRING>", test.text+"\n")...); res.String() != test.res {
			t.Error(fmt.Sprintf("#%d: Expected result '%s', got string '%s'", i, test.res, res))
		} else {
			//fmt.Printf("%v\n", res)
		}
	}
}

var json1_string = `"{
  \"a\": 1,
  \"b\": \"test\",
  \"c\": true,
  \"d\": [1, 2, 3],
  \"e\": {
    \"x\": 10.3,
    \"y\": 25.8,
    \"z\": [true, false]
  },
  \"f\": {},
  \"g\": [],
}"`
