package pl

import (
	"fmt"
	"testing"
)

func TestWord(t *testing.T) {
	word1 := NewWord("a")
	word2 := NewWord("a")
	word3 := NewWord("b")
	if word1.String() != word2.String() {
		t.Error(fmt.Sprintf("Expected string 'a' equals string 'a', got string '%s' not equal string '%s'", word1.String(), word2.String()))
	}
	if word2.String() == word3.String() {
		t.Error(fmt.Sprintf("Expected string 'a' not equals string 'b', got string '%s' equals string '%s'", word2.String(), word3.String()))
	}
	if word1 != word2 {
		t.Error(fmt.Sprintf("Expected word 'a' equals word 'a', got word '%v' not equal word '%v'", word1, word2))
	}
	if word2 == word3 {
		t.Error(fmt.Sprintf("Expected word 'a' not equals word 'b', got word '%v' equals word '%v'", word2, word3))
	}

}

func TestPlist(t *testing.T) {
	plist1 := NewPlist(NewWord("a"), NewWord("b"), NewWord("c"))
	if plist1.String() != "[a b c ]" {
		t.Error(fmt.Sprintf("Expected string '[a b c ]' equals string '[a b c ]', got string '%s' not equal string '%s'", plist1.String(), "[a b c ]"))
	}

	result := Begin().
		Eval(
			NewPlist(
				NewWord("prog"),
				NewLlist(NewWord("X"), NewLlist(NewWord("Y"), NewWord("ValueOfY")), NewWord("Z")),
				NewPlist(NewWord("set"), NewWord("X"), NewWord("ValueOfX")),
				NewPlist(NewWord("set"), NewWord("Z"), NewLlist(NewRef(LocalValue, NewWord("X")), NewRef(LocalValue, NewWord("Y")))),
				NewRef(LocalValue, NewWord("Z")),
			),
		).String()

	if result != "(ValueOfX ValueOfY )" {
		t.Error(fmt.Sprintf("Expected string '(ValueOfX ValueOfY )' equals string '(ValueOfX ValueOfY )', got string '%s' not equal string '%s'", result, "(ValueOfX ValueOfY )"))
	}
}

func TestLlist(t *testing.T) {
	llist1 := NewLlist(NewWord("x"), NewWord("y"), NewWord("z"))
	if llist1.String() != "(x y z )" {
		t.Error(fmt.Sprintf("Expected string '(x y z )' equals string '(x y z )', got string '%s' not equal string '%s'", llist1.String(), "(x y z )"))
	}
}
