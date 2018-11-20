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

func newToken(line int, col int, value string) *Token {
	t := new(Token)
	t.Line = line
	t.Col = col
	t.Type = determineType(value)
	t.Value = value
	return t
}

//Lex takes a sourcecode string and transforms it into usable tokens to build an AST with
func (l *Lexer) Lex(sourceCode string) {

	for l.currentIndex < len(sourceCode) {
		currTok := newToken(l.currentLine, l.currentCol, string(sourceCode[l.currentIndex]))
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
			fallthrough
		case "right_arrow":
			fallthrough
		case "comma":
			fallthrough
		case "left_bracket":
			fallthrough
		case "right_bracket":
			fallthrough
		case "double_dot":
			fallthrough
		case "semicolon":
			l.advance()
		case "undefined_symbol":
			errors.Report(l.currentLine, l.FileName, "undefined symbol used")
			os.Exit(65)
		case "newline":
			l.currentCol = 0
			l.advance()
			appendToken = false
		case "ignoreable":
			l.advance()
			appendToken = false
		}

		if appendToken {
			l.Tokens = append(l.Tokens, *currTok)
		}
	}
	l.tagKeywords()
}

func (l *Lexer) advance() {
	l.currentCol++
	l.currentIndex++
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
