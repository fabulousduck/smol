package ast

import (
	"os"

	"github.com/fabulousduck/proto/src/types"
	"github.com/fabulousduck/smol/lexer"
)

func (p *Parser) expectCurrent(expectedValues []string) {
	currentToken := p.currentToken()

	if !types.Contains(currentToken.Type, expectedValues) {
		lexer.ThrowSemanticError(&currentToken, expectedValues, p.Filename)
		os.Exit(65)
	}
}

func (p *Parser) expectNext(expectedValues []string) {
	nextToken := p.nextToken()
	if !types.Contains(nextToken.Type, expectedValues) {
		lexer.ThrowSemanticError(&nextToken, expectedValues, p.Filename)
		os.Exit(65)
	}
}

func (p *Parser) nextExists() bool {
	//+1 because we have to account for arrays starting at 0
	return p.TokensConsumed+1 < len(p.Tokens)
}

func (p *Parser) currentToken() lexer.Token {
	return p.Tokens[p.TokensConsumed]
}

func (p *Parser) nextToken() lexer.Token {
	return p.Tokens[p.TokensConsumed+1]
}

func (p *Parser) advance() {
	p.TokensConsumed++
}

func (p *Parser) advanceN(n int) {
	p.TokensConsumed += n
}

//NodeIsVariable allows for nice statements like if NodeIsVariable(node) {}
func NodeIsVariable(node Node) bool {
	return node.GetNodeName() == "statVar"
}

//GetAllocationType determines if a type must be stack or heap allocated
func (v *Variable) GetAllocationType() string {
	types := map[string]string{
		"Uint16": "stack",
		"Uint32": "stack",
		"Uint64": "stack",
		"String": "heap",
		"Bool":   "stack",
		"Char":   "stack",
	}

	if val, ok := types[v.Type]; ok {
		return val
	}

	return "heap"
}
