package ast

import (
	"os"

	"github.com/fabulousduck/smol/errors"
	"github.com/fabulousduck/smol/lexer"
)

//Expression contains nodes in RPN form
type Expression struct {
	Tokens []lexer.Token
}

func createExpression(nodes []lexer.Token) Expression {
	expression := Expression{Tokens: nodes}
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

//Operator is a symbol that operates on one or more sides
type Operator struct {
	value string
}

//GetNodeName is a generic function that allows subtypes of a node in the AST
func (o Operator) GetNodeName() string {
	return "operator"
}

//Symbol is a character that is not and operator or a latin character / numeral
type Symbol struct {
	value string
}

//GetNodeName is a generic function that allows subtypes of a node in the AST
func (s Symbol) GetNodeName() string {
	return "symbol"
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

//CreateSymbol takes a string and returns it wrapped in a node like
func CreateSymbol(value string) Symbol {
	return Symbol{value: value}
}

/*
readExpression turns a set of nodes into a expression using shunting yard
*/
func (p *Parser) readExpression() Expression {
	expressionLine := p.currentToken().Line
	expressionTokens := []lexer.Token{}

	//gather all tokens of the expression into a slice
	for currTok := p.currentToken(); currTok.Line == expressionLine; currTok = p.currentToken() {
		expressionTokens = append(expressionTokens, currTok)
		p.advance()
		if !p.nextExists() {
			expressionTokens = append(expressionTokens, p.currentToken())
			break
		}
	}
	expressionParser := NewParser(p.Filename, expressionTokens)

	switch len(expressionTokens) {
	case 0:
		errors.ExpectedExpressionError()
		os.Exit(65)
		break
	case 1:
		if lexer.IsLitteral(expressionTokens[0]) {
			litteralToken := expressionTokens[0]
			return createExpression([]lexer.Token{litteralToken})
		}
		errors.InvalidOperatorError()
		os.Exit(65)
		break
	default:
		return expressionParser.parseExpression()
	}

	return createExpression([]lexer.Token{})
}

/*
readExpressionUntil allows for parsing and expression with a defined
symbol as an end boundary.

This is used for when we parse inside of function calls
*/
func (p *Parser) readExpressionUntil(tokValues []string) (Expression, string) {
	expressionTokens := []lexer.Token{}
	fnContext := false
	delimFound := ""

	for i := 0; i < len(p.Tokens); i++ {
		if containsStr(tokValues, p.currentToken().Value) && !fnContext {
			delimFound = p.currentToken().Value
			break
		}
		if p.currentToken().Value == "(" {
			fnContext = true
		}

		if p.currentToken().Value == ")" && fnContext {
			fnContext = false
		}
		expressionTokens = append(expressionTokens, p.currentToken())
		p.advance()

	}

	expressionParser := NewParser(p.Filename, expressionTokens)
	expressionParser.Tokens = expressionTokens

	switch len(expressionTokens) {
	case 0:
		errors.ExpectedExpressionError()
		os.Exit(65)
		break
	case 1:
		if lexer.IsLitteral(expressionTokens[0]) {
			return createExpression([]lexer.Token{expressionTokens[0]}), delimFound
		}
		errors.InvalidOperatorError()
		os.Exit(65)
		break
	default:
		return expressionParser.parseExpression(), delimFound
	}

	return createExpression([]lexer.Token{}), delimFound
}

func (p *Parser) parseExpression() Expression {
	operatorStack := []lexer.Token{}
	outputQueue := []lexer.Token{}
	for p.TokensConsumed < len(p.Tokens) {
		token := p.currentToken()
		switch token.Type {
		case "comma":
			p.advance()
			break
		case "integer":
			outputQueue = append(outputQueue, token)
			p.advance()
			break
		case "less_than":
			fallthrough
		case "greater_than":
			fallthrough
		case "exponent":
			fallthrough
		case "division":
			fallthrough
		case "star":
			fallthrough
		case "plus":
			fallthrough
		case "dash":
			if len(operatorStack) >= 1 {
				for len(operatorStack) != 0 {

					stackTopAttributes := lexer.GetOperatorAttributes(top(operatorStack).Type)
					tokenAttributes := lexer.GetOperatorAttributes(token.Type)
					hasHigherPrec := top(operatorStack).HasHigherPrec(token)
					eqRule := stackTopAttributes.Precedance == tokenAttributes.Precedance && tokenAttributes.Associativity == "left"
					parenNotTop := top(operatorStack).Type != "left_parenthesis"

					if (hasHigherPrec || eqRule) && parenNotTop {

						outputQueue = append(outputQueue, top(operatorStack))
						operatorStack = (operatorStack)[:len(operatorStack)-1]

					} else {
						break
					}
				}
			}
			operatorStack = append(operatorStack, token)
			p.advance()
			break
		case "left_parenthesis":
			operatorStack = append(operatorStack, token)
			p.advance()
			break
		case "right_parenthesis":
			for top(operatorStack).Value != "(" {
				outputQueue = append(outputQueue, top(operatorStack))
				operatorStack = (operatorStack)[:len(operatorStack)-1]
			}
			if top(operatorStack).Value == "(" {
				operatorStack = (operatorStack)[:len(operatorStack)-1]
			}
			p.advance()
			break
		case "character":
			fallthrough
		case "string":
			if p.nextExists() && p.nextToken().Type == "left_parenthesis" {
				outputQueue = append(outputQueue, token)
			} else {
				operatorStack = append(operatorStack, token)
			}

			p.advance()
			break

		}
	}
	if len(operatorStack) != 0 {
		for i := len(operatorStack); 0 < i; i-- {
			outputQueue = append(outputQueue, operatorStack[i-1])
			operatorStack = (operatorStack)[:len(operatorStack)-1]
		}
	}

	return createExpression(outputQueue)
}

func top(sl []lexer.Token) lexer.Token {
	return sl[len(sl)-1]
}

func containsStr(a []string, b string) bool {
	for i := 0; i < len(a); i++ {
		if a[i] == b {
			return true
		}
	}
	return false
}
