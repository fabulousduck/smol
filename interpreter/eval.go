package interpreter

import "github.com/fabulousduck/smol/ast"

/*
EvalExpression takes an expression in the form of a node tree from the AST.
It will return a new variable with the result of the expression.
*/
func EvalExpression(nodes []*ast.Node) *ast.Variable {
	expressionResult := new(ast.Variable)

	return expressionResult
}
