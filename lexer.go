package losp

import (
	"bytes"
	"fmt"
)

type token struct {
	Value, Type string
	Line, Col   int
}

type lexer struct {
	tokens                                []*token
	currentIndex, currentLine, currentCol int
}

func (l *lexer) lex(sourceCode string) {

	for l.currentIndex < len(sourceCode) {
		currentChar := string(sourceCode[l.currentIndex])
		currTok := new(token)
		currTok.Line = l.currentLine
		currTok.Type = determineType(currentChar)
		appendToken := true
		switch currTok.Type {
		case "CHAR":
			currTok.Value = l.peekTypeN("CHAR", sourceCode)
		case "NUMB":
			currTok.Value = l.peekTypeN("NUMB", sourceCode)
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
		case "LEFT_BRACE":
			currTok.Value = "["
			currTok.Type = "LEFT_BRACE"
			l.currentCol++
			l.currentIndex++
		case "RIGHT_BRACE":
			currTok.Value = "}"
			currTok.Type = "RIGHT_BRACE"
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
			fmt.Println(currentChar)
			panic("undefined char")
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
			l.tokens = append(l.tokens, currTok)
		}
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
