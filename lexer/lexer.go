package lexer

import (
	"bytes"
	"fmt"
	"os"

	"github.com/fabulousduck/smol/errors"
)

//Token contains all info about a specific token from syntax
type Token struct {
	Value, Type string
	Line, Col   int
}

//Lexer contains all the info needed for the lexer to generate a set of usable tokens
type Lexer struct {
	Tokens                                []Token
	currentIndex, currentLine, currentCol int
	FileName                              string
}

//NewLexer creates a new instance of a lexer stuct
func NewLexer(filename string) *Lexer {
	l := new(Lexer)
	l.FileName = filename
	return l
}

//Lex takes a sourcecode string and transforms it into usable tokens to build an AST with
func (l *Lexer) Lex(sourceCode string) {

	for l.currentIndex < len(sourceCode) {
		currentChar := string(sourceCode[l.currentIndex])
		currTok := new(Token)
		currTok.Line = l.currentLine
		currTok.Col = l.currentCol
		currTok.Type = determineType(currentChar)
		appendToken := true
		switch currTok.Type {
		case "character":
			currTok.Value = l.peekTypeN("character", sourceCode)
			l.currentCol += len(currTok.Value)
		case "integer":
			currTok.Value = l.peekTypeN("integer", sourceCode)
			l.currentCol += len(currTok.Value)
		case "comment":
			appendToken = false
			l.readComment(sourceCode)
			l.currentCol = 0
		case "left_arrow":
			currTok.Value = "<"
			currTok.Type = "left_arrow"
			l.currentCol++
			l.currentIndex++
		case "right_arrow":
			currTok.Value = ">"
			currTok.Type = "right_arrow"
			l.currentCol++
			l.currentIndex++
		case "comma":
			currTok.Value = ","
			currTok.Type = "comma"
			l.currentCol++
			l.currentIndex++
		case "left_bracket":
			currTok.Value = "["
			currTok.Type = "left_bracket"
			l.currentCol++
			l.currentIndex++
		case "right_bracket":
			currTok.Value = "]"
			currTok.Type = "right_bracket"
			l.currentCol++
			l.currentIndex++
		case "double_dot":
			currTok.Value = ":"
			currTok.Type = "double_dot"
			l.currentCol++
			l.currentIndex++
		case "semicolon":
			currTok.Value = ";"
			currTok.Type = "semicolon"
			l.currentCol++
			l.currentIndex++
		case "space":
			l.currentCol++
			l.currentIndex++
			appendToken = false
		case "undefined_symbol":
			errors.Report(l.currentLine, l.FileName, "undefined symbol used")
			os.Exit(65)
		case "newline":
			l.currentCol = 0
			l.currentLine++
			l.currentIndex++
			appendToken = false
		case "tab":
			l.currentCol++
			l.currentIndex++
			appendToken = false
		}

		if appendToken {
			l.Tokens = append(l.Tokens, *currTok)
		}
	}
	l.tagKeywords()
}

func (l *Lexer) readComment(program string) {
	l.currentIndex++
	for t := determineType(string(program[l.currentIndex])); t != "newline"; t = determineType(string(program[l.currentIndex])) {
		l.currentIndex++
	}
}

func (l *Lexer) peekTypeN(typeName string, program string) string {
	var currentString bytes.Buffer
	for t := determineType(string(program[l.currentIndex])); t == typeName; t = determineType(string(program[l.currentIndex])) {
		if l.currentIndex+1 >= len(program) {
			currentString.WriteString(string(program[l.currentIndex]))
			l.currentIndex++
			return currentString.String()
		}
		currentString.WriteString(string(program[l.currentIndex]))
		l.currentIndex++
	}
	return currentString.String()
}

func (l *Lexer) tagKeywords() {
	for i := 0; i < len(l.Tokens); i++ {
		if l.Tokens[i].Type == "character" {
			l.Tokens[i].Type = getKeyword(&l.Tokens[i])
		}
	}
}

//ThrowSemanticError can be used when an error occurs while generating an AST and not at interpret time
func ThrowSemanticError(token *Token, expected []string, filename string) {
	errors.Report(
		token.Line,
		filename,
		fmt.Sprintf("expected one of [%s]. got %s",
			errors.ConcatVariables(expected, ", "),
			token.Type),
	)
}
