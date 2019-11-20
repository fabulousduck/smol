package ast

import (
	"fmt"
	"os"

	"github.com/fabulousduck/smol/errors"
	"github.com/fabulousduck/smol/lexer"
)

//Expression contains nodes in RPN form
type Expression struct {
	Nodes []Node
}

func createExpression(nodes []Node) *Expression {
	expression := new(Expression)
	expression.Nodes = nodes
	return expression
}

//GetNodeName so its valid on the node interface
//and we can ask what type it is later
func (e Expression) GetNodeName() string {
	return "expression"
}

//VariableReference contains the name of a referenced variable during AST generation
type VariableReference struct {
	name string
}

//GetNodeName so its valid on the node interface
//and we can ask what type it is later
func (vr VariableReference) GetNodeName() string {
	return "variableReference"
}

func createVariableReference(name string) VariableReference {
	return VariableReference{name: name}
}

//Litteral is a node type for static values such as integers and string litterals
type Litteral struct {
	ltype string
	value string
}

//GetNodeName is a generic function that allows subtypes of a node in the AST
func (l Litteral) GetNodeName() string {
	return "litteral"
}

//CreateLitteral creates a new Litteral struct
func CreateLitteral(value string, ltype string) Litteral {
	litteral := Litteral{value: value, ltype: ltype}
	return litteral
}

/*
readExpression turns a set of nodes into a expression using shunting yard
*/
func (p *Parser) readExpression() *Expression {
	expression := new(Expression)
	expressionLine := p.currentToken().Line
	expressionTokens := []lexer.Token{}

	fmt.Printf("expression line: %d\n", expressionLine)

	for currTok := p.currentToken(); currTok.Line == expressionLine; currTok = p.currentToken() {
		expressionTokens = append(expressionTokens, currTok)
		p.advance()
		if !p.nextExists() {
			break
		}
	}

	expressionParser := NewParser(p.Filename, expressionTokens)

	switch len(expressionTokens) {
	case 0:
		errors.ExpectedExpressionError()
		os.Exit(65)
	case 1:
		if lexer.IsLitteral(expressionTokens[0]) {
			litteralToken := expressionTokens[0]
			return createExpression([]Node{CreateLitteral(litteralToken.Value, litteralToken.Type)})
		}
		errors.InvalidOperatorError()
		os.Exit(65)
	default:
		return expressionParser.parseExpression()
	}

	return expression
}

/*
readExpressionUntil allows for parsing and expression with a defined
symbol as an end boundary
*/
func (p *Parser) readExpressionUntil() *Expression {
	expressionAST := createExpression([]Node{})

	return expressionAST
}

func (p *Parser) parseExpression() *Expression {
	operatorStack := []Node{}
	outputQueue := []Node{}

	for p.TokensConsumed < len(p.Tokens) {
		token := p.currentToken()

		switch token.Type {
		case "integer":
			outputQueue = append(operatorStack, CreateLitteral(token.Value, token.Type))
			p.advance()
			break
		case "string":
			//its a function
			if p.nextToken().Type == "left_parenthesis" {
				outputQueue = append(outputQueue, p.createFunctionCall())
				break
			}
			//we treat variable names the same as integers because all other variable values are not allowed in expressions
			operatorStack = append(operatorStack, createVariableReference(token.Value))
			p.advance()
			break

		}
	}

	return createExpression(outputQueue)
}
