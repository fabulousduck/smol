package interpreter

import (
	"github.com/fabulousduck/smol/ast"
	"github.com/fabulousduck/smol/lexer"
)

/*
EvalExpression takes an expression in the form of a node tree from the AST.
It will return a new variable with the result of the expression.
*/
func EvalExpression(expression ast.Expression) string {
	expressionResult := ""
	stack := []lexer.Token{}

	for _, token := range expression.Tokens {
		switch token.Type {
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
			temp := stack[len(stack)-1]                   //top
			stack = (stack)[:len(stack)-1]                //pop
			res := exec(temp, token, stack[len(stack)-1]) //exec
			stack = (stack)[:len(stack)-1]                //pop
			stack = append(stack, res)                    //push
			break

		case "left_parenthesis":
			break
		case "right_parenthesis":
			break
		case "character":
			fallthrough
		case "string":
			fallthrough
		case "integer":
			stack = append(stack, token)

		default:
			//TODO: error: improper type found in expression
		}
	}
	return expressionResult
}

func exec(rhs lexer.Token, operator lexer.Token, lhs lexer.Token) lexer.Token {
	result := lexer.Token{}

	switch operator.Type {
	case "less_than":

	}

	return result
}
