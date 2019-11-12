package ast

import (
	"os"

	"github.com/fabulousduck/smol/errors"
	"github.com/fabulousduck/smol/lexer"
)

//Expression is a set of tokens in shunting yard format
type Expression struct {
	nodes []Node
}

//GetNodeName so its valid on the node interface
//and we can ask what type it is later
func (e Expression) GetNodeName() string {
	return "expression"
}

//Litteral is a node type for static values such as integers and string litterals
type Litteral struct {
	value string
}

func (l Litteral) GetNodeName() string {
	return "litteral"
}

//CreateExpression returns a new pointer to an expression
func CreateExpression(expressionNodes []Node) *Expression {
	expression := new(Expression)
	expression.nodes = expressionNodes
	return expression
}

//CreateLitteral creates a new Litteral struct
func CreateLitteral(value string) Litteral {
	litteral := Litteral{value: value}
	return litteral
}

/*
readExpression turns a set of nodes into a expression using shunting yard
*/
func (p *Parser) readExpression() *Expression {
	expression := new(Expression)
	expressionLine := p.currentToken().Line
	expressionTokens := []lexer.Token{}

	for currTok := p.currentToken(); currTok.Line == expressionLine; currTok = p.currentToken() {
		expressionTokens = append(expressionTokens, currTok)
	}

	switch len(expressionTokens) {
	case 0:
		errors.ExpectedExpressionError()
		os.Exit(65)
	case 1:
		if lexer.IsLitteral(expressionTokens[0]) {
			return CreateExpression([]Node{CreateLitteral(expressionTokens[0].Value)})
		}
		errors.InvalidOperatorError()
		os.Exit(65)
	default:
		return CreateExpression(applyShuntingYard(expressionTokens))
	}

	return expression
}

/*
readExpressionUntil allows for parsing and expression with a defined
symbol as an end boundary
*/
func (p *Parser) readExpressionUntil() []*Node {
	expressionAST := []*Node{}

	return expressionAST
}

func applyShuntingYard(tokens []lexer.Token) []Node {
	operatorStack := []Node{}
	outputQueue := []Node{}

	for _, token := range tokens {

	}

	return outputQueue
}
