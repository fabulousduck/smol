package ast

func (p *Parser) readExpression() []*Node {
	expressionAST := []*Node{}
	expressionLine := p.currentToken().Line

	for currentExpressionToken := p.currentToken(); currentExpressionToken.Line == expressionLine; currentExpressionToken = p.currentToken() {

	}

	return expressionAST
}

/*
readExpressionUntil allows for parsing and expression with a defined
symbol as an end boundary
*/
func (p *Parser) readExpressionUntil() []*Node {
	expressionAST := []*Node{}

	return expressionAST
}
