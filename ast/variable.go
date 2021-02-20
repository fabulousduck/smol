package ast

import "github.com/fabulousduck/smol/lexer"

//StringLit represents a string litteral
type StringLit struct {
	Value string
}

func (sl StringLit) GetNodeName() string {
	return "stringLit"
}

//BoolLit represents a boolean litteral being either "True" or "False"
type BoolLit struct {
	Value string
}

func (bm BoolLit) GetNodeName() string {
	return "boolLit"
}

//NumLit represents a numeric litteral.
type NumLit struct {
	Value string
}

func (nm NumLit) GetNodeName() string {
	return "numLit"
}

/*
createVariable reads tokens to create a variable
It adheres to the following structure

<type> <name> <value>

*/
func (p *Parser) createVariable() *Variable {
	variable := new(Variable)

	p.expectCurrent([]string{"variable_type"})
	variable.Type = p.currentToken().Value
	p.advance()

	p.expectCurrent([]string{"character", "string"})
	variable.Name = p.currentToken().Value
	p.advance()

	p.expectCurrent([]string{"equals"})
	p.advance()
	variable.ValueExpression = p.readExpression()
	return variable
}

func createLit(token lexer.Token) Node {
	switch token.Type {
	case "integer":
		nl := new(NumLit)
		nl.Value = token.Value
		return nl
	case "boolean_keyword":
		bl := new(BoolLit)
		bl.Value = token.Value
		return bl
	case "string_litteral":
		sl := new(StringLit)
		sl.Value = token.Value
		return sl
	default:
		sv := new(StatVar)
		sv.Value = token.Value
		return sv
	}
}
