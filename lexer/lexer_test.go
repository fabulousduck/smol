package lexer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	testVariable = "MEM A 10;"
)

func TestVariable(T *testing.T) {
	testProgram := "MEM A 10;"
	l := NewLexer("TESTING")
	expectedResults := []Token{
		{
			Value: "MEM",
			Type:  "variable_assignment",
			Line:  0,
			Col:   0,
		},
		{
			Value: "A",
			Type:  "CHAR",
			Line:  0,
			Col:   4,
		},
		{
			Value: "10",
			Type:  "NUMB",
			Line:  0,
			Col:   6,
		},
		{
			Value: ";",
			Type:  "SEMICOLON",
			Line:  0,
			Col:   8,
		},
	}

	l.Lex(testProgram)

	for i := 0; i < len(expectedResults); i++ {
		if !cmp.Equal(expectedResults[i], l.Tokens[i]) {
			//TODO: give a bit more info on this
			T.Logf("\nTestVariableLex | failed to generate correct tokens for variable assignment")
			T.Fail()
		}
	}
}

func TestStringTypeDetermination(T *testing.T) {
	testNumber := "15"
	testChar := "boop"

	if DetermineStringType(testNumber) != "NUMB" {
		T.Logf("\nTestStringTypeDetermination | determined that %s is of type %s which is actually NUMB", testNumber, DetermineStringType(testNumber))
		T.Fail()
	}

	if DetermineStringType(testChar) != "CHAR" {
		T.Logf("\nTestStringTypeDetermination | determined that %s is of type %s which is actually CHAR", testNumber, DetermineStringType(testNumber))
		T.Fail()
	}
}

func TestDetermineType(T *testing.T) {
	values := []string{
		"1", "B", "<", ",", ">", ";", "[", "]", ":",
		"#", "\r", "\n", "\t", " ", "&",
	}
	expectedTypes := []string{
		"NUMB", "CHAR", "LEFT_ARROW", "COMMA", "RIGHT_ARROW", "SEMI_COLON",
		"LEFT_BRACKET", "RIGHT_BRACKET", "DOUBLE_DOT", "COMMENT", "WIN_NEWLINE",
		"NEWLINE", "TAB", "SPACE", "UDEF",
	}

	for i := 0; i < len(values); i++ {
		determinedType := determineType(values[i])
		if determinedType != expectedTypes[i] {
			T.Logf("\n TestDetermineType | type of %s wat determined to be %s. should be %s", values[i], determinedType, expectedTypes[i])
			T.Fail()
		}
	}
}
