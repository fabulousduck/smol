package lexer

import "testing"

const (
	testVariable = "MEM A 10;"
)

func TestVariable(T *testing.T) {
	testProgram := "MEM A 10;"
	expectedResults := []Token{}
	l := NewLexer(testProgram, "TESTING")
}
