package smol

import (
	"bytes"
	"os"

	"github.com/davecgh/go-spew/spew"
)

type token struct {
	Value, Type string
	Line, Col   int
}

type lexer struct {
	tokens                                []token
	currentIndex, currentLine, currentCol int
}

func (l *lexer) lex(sourceCode string, filename string) {

	for l.currentIndex < len(sourceCode) {
		currentChar := string(sourceCode[l.currentIndex])
		currTok := new(token)
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
			curr := l.currentIndex
			l.readComment(sourceCode)
			l.currentCol = l.currentIndex - curr
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
			spew.Dump(currentChar)
			report(l.currentLine, filename, "undefined symbol used")
			os.Exit(65)
		case "WIN_NEWLINE":
			fallthrough
		case "NEWLINE":
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
			l.tokens = append(l.tokens, *currTok)
		}
	}
	l.tagKeywords()
}

func (l *lexer) readComment(program string) {
	l.currentIndex++
	for t := determineType(string(program[l.currentIndex])); t != "NEWLINE" && t != "WIN_NEWLINE"; t = determineType(string(program[l.currentIndex])) {
		l.currentIndex++
	}
}

func (l *lexer) peekTypeN(typeName string, program string) string {
	var currentString bytes.Buffer
	for t := determineType(string(program[l.currentIndex])); t == typeName; t = determineType(string(program[l.currentIndex])) {
		currentString.WriteString(string(program[l.currentIndex]))
		l.currentIndex++
	}
	return currentString.String()
}

func (l *lexer) tagKeywords() {
	for i := 0; i < len(l.tokens); i++ {
		if l.tokens[i].Type == "CHAR" {
			l.tokens[i].Type = getKeyword(&l.tokens[i])
		}
	}
}
