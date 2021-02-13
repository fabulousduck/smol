package lexer

import (
	"bytes"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"

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
	l.currentIndex = 0
	l.currentLine = 1
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
			currTok.Value = l.peekTypesN([]string{"integer", "character"})
		case "integer":
			currTok.Value = l.peekTypesN([]string{"integer"})
		case "comment":
			l.readComment()
			l.advance()
			l.currentCol = 0
			continue
		case "double_quote":
			l.advance()
			currTok.Value = l.peekUntil([]string{"double_quote"})
			currTok.Type = "string_litteral"
		case "equals":
			if l.peek() == "=" {
				currTok.Value = "=="
				currTok.Type = "comparison"
				l.advance()
			}
			l.advance()
		case "plus":
			if l.peek() == "+" {
				currTok.Value = "++"
				currTok.Type = "direct_variable_operation"
				l.advance()
			} else if l.peek() == "=" {
				currTok.Value = "+="
				currTok.Type = "direct_variable_operation"
				l.advance()
			}
			l.advance()
		case "dash":
			if l.peek() == "-" {
				currTok.Value = "--"
				currTok.Type = "direct_variable_operation"
				l.advance()
			}
			l.advance()
		case "star":
			fallthrough
		case "division":
			fallthrough
		case "less_than":
			fallthrough
		case "greater_than":
			fallthrough
		case "comma":
			fallthrough
		case "left_bracket":
			fallthrough
		case "right_bracket":
			fallthrough
		case "left_parenthesis":
			fallthrough
		case "right_parenthesis":
			fallthrough
		case "double_dot":
			fallthrough
		case "exponent":
			fallthrough
		case "semicolon":
			l.advance()
		case "undefined_symbol":
			errors.Report(l.currentLine, l.FileName, fmt.Sprintf("undefined symbol \"%s\" used", currTok.Value))
			os.Exit(65)
		case "newline":
			l.currentCol = 0
			l.currentLine++
			l.advance()
			continue
		case "ignoreable":
			l.currentCol = 0
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

func (l *Lexer) peek() string {
	if l.currentIndex+1 <= len(l.Program) {
		return string(l.Program[l.currentIndex+1])
	}
	return ""
}

func (l *Lexer) peekUntil(types []string) string {
	var currentString bytes.Buffer

	//loop until the current character type (t) is in types
	for t := determineType(l.currentChar()); !contains(t, types); t = determineType(l.currentChar()) {
		char := l.currentChar()

		//we do this to avoid index out of range errors
		if l.currentIndex+1 >= len(l.Program) {
			spew.Dump("out of range clause")
			currentString.WriteString(char)
			l.advance()

			return currentString.String()
		}
		currentString.WriteString(char)
		l.advance()
	}
	//advance over the untill symbol
	l.advance()
	return currentString.String()
}

func (l *Lexer) peekTypesN(types []string) string {
	var currentString bytes.Buffer

	for t := determineType(l.currentChar()); contains(t, types); t = determineType(l.currentChar()) {
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
