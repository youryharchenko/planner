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
}
