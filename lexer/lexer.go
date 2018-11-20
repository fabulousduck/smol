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
	FileName, Program                     string
}

//NewLexer creates a new instance of a lexer stuct
func NewLexer(filename string, program string) *Lexer {
	l := new(Lexer)
	l.Program = program
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
func (l *Lexer) Lex() {

	for l.currentIndex < len(l.Program) {
		currTok := newToken(l.currentLine, l.currentCol, l.currentChar())
		switch currTok.Type {
		case "character":
			currTok.Value = l.peekTypeN("character")
			l.currentCol += len(currTok.Value)
		case "integer":
			currTok.Value = l.peekTypeN("integer")
			l.currentCol += len(currTok.Value)
		case "comment":
			l.readComment()
			l.advance()
			l.currentCol = 0
			continue
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
			continue
		case "ignoreable":
			l.advance()
			continue
		}

		l.Tokens = append(l.Tokens, *currTok)

	}
	l.tagKeywords()
}

func (l *Lexer) advance() {
	l.currentCol++
	l.currentIndex++
}

func (l *Lexer) readComment() {
	l.currentIndex++
	for t := determineType(l.currentChar()); t != "newline"; t = determineType(l.currentChar()) {
		l.currentIndex++
	}
}

func (l *Lexer) peekTypeN(typeName string) string {
	var currentString bytes.Buffer

	for t := determineType(l.currentChar()); t == typeName; t = determineType(l.currentChar()) {
		char := l.currentChar()

		//we do this to avoid index out of range errors
		if l.currentIndex+1 >= len(l.Program) {

			currentString.WriteString(char)
			l.advance()

			return currentString.String()
		}
		currentString.WriteString(char)
		l.advance()
	}

	return currentString.String()
}

func (l *Lexer) currentChar() string {
	return string(l.Program[l.currentIndex])
}

func (l *Lexer) tagKeywords() {
	for i, token := range l.Tokens {
		if token.Type == "character" {
			l.Tokens[i].Type = getKeyword(&token)
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
