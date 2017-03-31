package main

import (
	"log"
	"strings"

	"github.com/youryharchenko/planner/pl"
)

func main() {
	/*
			w := pl.NewWord("Hello")
			log.Println("Word String():", w.String())
			log.Println("Word Eval():", pl.Begin().Eval(w).String())

			ref := pl.NewRef(pl.LocalValue, w)
			set := pl.NewPlist(pl.NewWord("set"), w, pl.NewWord("World!!!"))
			log.Println("Set String():", set.String())
			log.Println("Ref String():", ref.String())

			llist := pl.NewLlist(w, ref)

			log.Println("Ref Eval():", pl.Begin().Eval(set, ref, llist).String())

			log.Println("Ref Eval():",
				pl.Begin().
					Eval(
						pl.NewLlist(
							pl.NewPlist(pl.NewWord("set"), pl.NewWord("Hello"), pl.NewWord("Youry")),
							pl.NewRef(pl.LocalValue, pl.NewWord("Hello")),
							pl.NewPlist(pl.NewWord("quote"), pl.NewPlist(pl.NewWord("set"), pl.NewWord("Hello"), pl.NewWord("Youry"))),
							pl.NewPlist(pl.NewWord("quote"), pl.NewRef(pl.LocalValue, pl.NewWord("Hello"))),
						),
					).String(),
			)

		log.Println("Prog Eval():",
			pl.Begin().
				Eval(
					pl.NewPlist(
						pl.NewWord("prog"),
						pl.NewLlist(pl.NewWord("X"), pl.NewLlist(pl.NewWord("Y"), pl.NewWord("ValueOfY")), pl.NewWord("Z")),
						pl.NewPlist(pl.NewWord("set"), pl.NewWord("X"), pl.NewWord("ValueOfX")),
						pl.NewPlist(pl.NewWord("set"), pl.NewWord("Z"), pl.NewLlist(pl.NewRef(pl.LocalValue, pl.NewWord("X")), pl.NewRef(pl.LocalValue, pl.NewWord("Y")))),
						pl.NewRef(pl.LocalValue, pl.NewWord("Z")),
					),
				).String(),
		)
	*/
	src0 := `
	{prog (a b c) {set c {sum$int 1 .b}} {set b {sum$int 1 .a}} {set a 1} .c}
	`
	log.Println("Prog SourceStream():", pl.Begin().SourceStream(strings.NewReader(src0)))

	src1 := `
		{prog (X (Y ValueOfY) Z) {set X ValueOfX} {set Z (.X .Y)} .Z}
		`
	log.Println("Prog SourceStream():", pl.Begin().SourceStream(strings.NewReader(src1)))

	src2 := `
		{prog (X (Y 1) Z) {set X 2.5} {set Z {sum$int .X .Y}} .Z}
		`
	log.Println("Prog SourceStream():", pl.Begin().SourceStream(strings.NewReader(src2)))

	src3 := `
		{prog
			(X (Y 1) Z)
			{set X 2.5}
			{set Z ({sum$int {set X .X} {set Y .Y}} {sum$float .X .Y})}
			.Z
		}
		`
	log.Println("Prog SourceStream():", pl.Begin().SourceStream(strings.NewReader(src3)))

	src4 := `
		{prog
			((X 1) (Y 2) (Z 3))
			(fold sum (1 2 3) {fold sum$float 0 (.X .Y .Z)})
		}
		`
	log.Println("Prog SourceStream():", pl.Begin().SourceStream(strings.NewReader(src4)))

	src5 := `
		{prog
			(X Y Z)
			(fold sum  {print ({set X 1} {set Y 2} {set Z 3})} is {fold sum$int 0 (.X .Y .Z)} )
		}
		`
	log.Println("Prog SourceStream():", pl.Begin().SourceStream(strings.NewReader(src5)))

	src6 := `
		{prog
			()
			{map print ((A B) (C D))}
		}
		`
	log.Println("Prog SourceStream():", pl.Begin().SourceStream(strings.NewReader(src6)))

	src7 := `
		{prog
			(X Y Z)
			(fold prod {print ({set X 1} {set Y 2} {set Z 3})} is {fold prod$int 1 (.X .Y .Z)} )
		}
		`
	log.Println("Prog SourceStream():", pl.Begin().SourceStream(strings.NewReader(src7)))

	src9 := `
	{prog (a b c) {set c {sum$int 1 .b}} {set b {sum$int 1 .a}} {set a 1} {exit .c} }
	`
	log.Println("Prog SourceStream():", pl.Begin().SourceStream(strings.NewReader(src9)))

}
