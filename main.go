package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/youryharchenko/planner/pl"
)

func main() {
	file := flag.String("f", "", "source filename")
	flag.Parse()
	fmt.Println("PLANNER: Started, v.0.1.19")
	fmt.Printf("Source file: %s\n", *file)
	if cont, err := ioutil.ReadFile(*file); err == nil {
		src := string(cont)
		fmt.Println(pl.Begin().Eval(pl.ParseFromString("<STRING>", src)...))
	} else {
		fmt.Println(err)
	}

	/*
		src9 := `
		{let
				[a b c]
				{set c {sum$int 1 .b}}
				{set b {sum$int 1 .a}}
				{set a 1}
				{print (.a + .b = .c)}
		}
		`
	*/

}
