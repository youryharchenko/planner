package pl

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestOther(t *testing.T) {
	tests := []Test{
		// JSON
		Test{fmt.Sprintf("{def json1 %s}", json1_string), strings.Replace(json1_string, "\\\"", "\"", -1)},
		Test{"{get$json .json1}", "([a 1.000000] [b \"test\"] [c T] [d (1.000000 2.000000 3.000000)] [e ([x 10.300000] [y 25.800000] [z (T ())])] [f ()] [g ()])"},
		// Actor
		//Test{"{def actor-def {omega [mess] {print .mess}}}", "omega"},
		//Test{"{let [[a {start actor-def actor-test}]] {send .a \"Hello, world\"} {send .a Ok}}", "actor-test"},
		Test{"{def ping-def {omega [mess] {print {cdr .mess}} {send {car .mess} {cons .me \"Ping\"}}}}", "omega"},
		Test{"{def pong-def {omega [mess] {print {cdr .mess}} {send {car .mess} {cons .me \"Pong\"}}}}", "omega"},
		Test{"{def ping {start ping-def Ping}}", "Ping"},
		Test{"{def pong {start pong-def Pong}}", "Pong"},
		Test{"{send :pong {cons .ping \"Start\"}}", "Pong"},
		Test{"{send :pong stop}", "Pong"},
		Test{"{send :ping stop}", "Ping"},
		//Test{"{let [[ping {start ping-def Ping}] [pong {start pong-def Pong}]] {send .pong {cons .ping \"Start\"}} {send .pong stop} {send .ping stop}}", "Ping"},
		//Test{"", ""},

	}

	env := Begin()
	for i, test := range tests {
		//log.Println(i, test.text, "->", test.res)
		if res := env.Eval(ParseFromString("<STRING>", test.text+"\n")...); res.String() != test.res {
			t.Error(fmt.Sprintf("#%d: (%s)  Expected result '%s', got string '%s'", i, test.text, test.res, res))
		} else {
			//fmt.Printf("%v\n", res)
		}
		time.Sleep(time.Millisecond)
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
