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
		case "CHAR":
			currTok.Value = l.peekTypeN("CHAR", sourceCode)
			l.currentCol += len(currTok.Value)
		case "NUMB":
			currTok.Value = l.peekTypeN("NUMB", sourceCode)
			l.currentCol += len(currTok.Value)
		case "COMMENT":
			appendToken = false
			l.readComment(sourceCode)
			l.currentCol = 0
		case "LEFT_ARROW":
			currTok.Value = "<"
			currTok.Type = "LEFT_ARROW"
			l.currentCol++
			l.currentIndex++
		case "RIGHT_ARROW":
			currTok.Value = ">"
			currTok.Type = "RIGHT_ARROW"
			l.currentCol++
			l.currentIndex++
		case "COMMA":
			currTok.Value = ","
			currTok.Type = "COMMA"
			l.currentCol++
			l.currentIndex++
		case "LEFT_BRACKET":
			currTok.Value = "["
			currTok.Type = "LEFT_BRACKET"
			l.currentCol++
			l.currentIndex++
		case "RIGHT_BRACKET":
			currTok.Value = "]"
			currTok.Type = "RIGHT_BRACKET"
			l.currentCol++
			l.currentIndex++
		case "DOUBLE_DOT":
			currTok.Value = ":"
			currTok.Type = "DOUBLE_DOT"
			l.currentCol++
			l.currentIndex++
		case "SEMI_COLON":
			currTok.Value = ";"
			currTok.Type = "SEMICOLON"
			l.currentCol++
			l.currentIndex++
		case "SPACE":
			l.currentCol++
			l.currentIndex++
			appendToken = false
		case "UDEF":
			errors.Report(l.currentLine, l.FileName, "undefined symbol used")
			os.Exit(65)
		case "WIN_NEWLINE":
			fallthrough
		case "UNIX_NEWLINE":
			l.currentCol = 0
			l.currentLine++
			l.currentIndex++
			appendToken = false
		case "TAB":
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
	for t := determineType(string(program[l.currentIndex])); t != "UNIX_NEWLINE" && t != "WIN_NEWLINE"; t = determineType(string(program[l.currentIndex])) {
		l.currentIndex++
	}
}

func (l *Lexer) peekTypeN(typeName string, program string) string {
	var currentString bytes.Buffer
	for t := determineType(string(program[l.currentIndex])); t == typeName; t = determineType(string(program[l.currentIndex])) {
		currentString.WriteString(string(program[l.currentIndex]))
		l.currentIndex++
	}
	return currentString.String()
}

func (l *Lexer) tagKeywords() {
	for i := 0; i < len(l.Tokens); i++ {
		if l.Tokens[i].Type == "CHAR" {
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
