package ast

import (
	"os"

	"github.com/davecgh/go-spew/spew"
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
	Symbol   Symbol
	LHS, RHS Node
}

//GetNodeName is a generic function that allows subtypes of a node in the AST
func (o Operator) GetNodeName() string {
	return "operator"
}

//Symbol is an operator by itself
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
		expr := expressionParser.parseExpression()
		p.rpnToTree(expr)
		return expr, delimFound
	}

	return createExpression([]lexer.Token{}), delimFound
}

func (p *Parser) rpnToTree(e Expression) []Node {
	tree := []Node{}
	workingStack := invertStack(e.Tokens)
	spew.Dump(workingStack)
	for 0 < len(e.Tokens) {
		ctop := top(workingStack)
		if lexer.IsOperator(ctop.Type) {
			operator := Operator{}
			operator.Symbol = CreateSymbol(ctop.Value)
			//pop the operator of the input stack
			workingStack = (workingStack)[:len(workingStack)-1]
			//lhs
			operator.LHS = tree[len(tree)-1]
			//pop lhs off tree
			tree = (tree)[:len(tree)-1]
			//rhs
			operator.RHS = tree[len(tree)-1]
			//pop rhs off the tree
			tree = (tree)[:len(tree)-1]
			//append operator to the tree
			tree = append(tree, operator)
		} else {
			switch ctop.Type {
			case "integer":
				tree = append(tree, CreateLitteral(ctop.Value, ctop.Type))
			case "string":
				fallthrough
			case "character":
				tree = append(tree, createVariableReference(ctop.Value))
			}
		}
	}

	return tree
}

/*
LEFT OFF HERE.

PROBLEM:
	WHEN WE PARSE AN EXPRESSION, IT COULD BE THAT THERE IS A FUNCTION CALL IN THERE
	IF THERE IS, WE NEED TO GET ALL THE ARGUMENTS TO THAT FUNCTION CALL IN ORDER TO MAKE A FUNCTIONCALL NODE
	WE DO NOT KNOW WHERE ARGUMENTS START AND END WHEN THIS PROCESS OF RPN GENERATION IS DONE

POSSIBLE SOLUTION:
	- IF WE FIND A FUNCTION NAME, READEXPRUNTILL([]STRING{","}) UNTIL WE HIT A )
	  THESE WILL BE SUB EXPRESSION WRAPPED IN A TOKEN WITH TYPE SUBEXPRESSION AND VALUE OF THE EXPRESSION IN RPN FORM
	  THEN WE PARSE ALL THOSE EXPRESSIONS THROUGH THE RPN PARSER INDIVIDUALLY
	  ONCE DONE, WE CAN THEN RPN PARSE THE FUNCTION CALL AND ITS ARGUMENTS TO END UP WITH VALID RPN
	  THEN CONTINUE THROUGH THE PARSER UNTIL DONE
*/

func (p *Parser) parseExpression() Expression {
	operatorStack := []lexer.Token{}
	outputQueue := []lexer.Token{}
	for p.TokensConsumed < len(p.Tokens) {
		token := p.currentToken()
		spew.Dump(token)
		switch token.Type {
		case "comma":
			// subExpr := lexer.Token{}
			operatorStack = append(operatorStack, token)
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
		case "function_name":
			operatorStack = append(operatorStack, token)
			p.advance()
		case "character":
			fallthrough
		case "string":
			outputQueue = append(outputQueue, token)
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

func invertStack(stack []lexer.Token) []lexer.Token {
	nstack := []lexer.Token{}
	for i := len(stack) - 1; i >= 0; i-- {
		nstack = append(nstack, stack[i])
	}

	return nstack
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
