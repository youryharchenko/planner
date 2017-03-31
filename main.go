package main

import (
	"log"

	"github.com/youryharchenko/planner/pl"
)

func main() {

	src9 := `
	{prog
			(a b c)
			{set c {sum$int 1 .b}}
			{set b {sum$int 1 .a}}
			{set a 1}
			{exit .c}
	}
	`
	log.Println("Prog SourceStream():", pl.Begin().Eval(pl.ParseFromString("<STRING>", src9)...))

}
