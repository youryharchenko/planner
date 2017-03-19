package main

import (
	"log"
	"planner/pl"
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
	*/
	log.Println("Prog Eval():",
		pl.Begin().
			Eval(
				pl.NewPlist(
					pl.NewWord("prog"),
					pl.NewLlist(pl.NewWord("X"), pl.NewLlist(pl.NewWord("Y"), pl.NewWord("Value of Y")), pl.NewWord("Z")),
					pl.NewPlist(pl.NewWord("set"), pl.NewWord("X"), pl.NewWord("Value of X")),
					pl.NewRef(pl.LocalValue, pl.NewWord("Y")),
				),
			).String(),
	)

}
