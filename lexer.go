package losp

import (
	"bytes"
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
		switch currTok.Type {
		case "CHAR":
			currTok.Value = l.peekTypeN("CHAR", sourceCode)
			l.tokens = append(l.tokens, currTok)
		case "NUMB":
			l.currentIndex++
		case "LEFT_ARROW":
			l.currentIndex++
		case "RIGHT_ARROW":
			l.currentIndex++
		case "COMMA":
			l.currentIndex++
		case "LEFT_BRACE":
			l.currentIndex++
		case "RIGHT_BRACE":
			l.currentIndex++
		case "SEMI_COLON":
			l.currentIndex++
		case "SYMB":
			l.currentIndex++
		case "SPACE":
			l.currentIndex++
		case "UDEF":
			l.currentIndex++
		case "NEWLINE":
			l.currentIndex++
		case "TAB":
			l.currentIndex++
		default:
			panic(currentChar)
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
