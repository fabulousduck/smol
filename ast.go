package smol

import (
	"os"

	"github.com/davecgh/go-spew/spew"
)

type numLit struct {
	value string
}

func (nm numLit) getNodeName() string {
	return "numLit"
}

type statVar struct {
	value string
}

func (sv statVar) getNodeName() string {
	return "statVar"
}

type function struct {
	name   string
	params []string
	body   []node
}

func (f function) getNodeName() string {
	return "function"
}

type variable struct {
	name  string
	value string
}

type anb struct {
	lhs  node
	rhs  node
	body []node
}

func (anb anb) getNodeName() string {
	return "anb"
}

type statement struct {
	lhs string
	rhs node
}

func (s statement) getNodeName() string {
	return "statement"
}

func (v variable) getNodeName() string {
	return "variable"
}

type functionCall struct {
	name string
	args []string
}

func (fc functionCall) getNodeName() string {
	return "functionCall"
}

type node interface {
	getNodeName() string
}

type parser struct {
	ast      []node
	filename string
}

func NewParser(filename string) *parser {
	p := new(parser)
	p.filename = filename
	return p
}

func (p *parser) parse(tokens []token) ([]node, int) {
	nodes := []node{}
	for i := 0; i < len(tokens); {
		switch tokens[i].Type {
		case "variable_assignment":
			node, tokensConsumed := p.createVariable(tokens, i+1)
			i += tokensConsumed + 1
			nodes = append(nodes, node)
		case "function_definition":
			node, tokensConsumed := p.createFunctionHeader(tokens, i+1)
			i += tokensConsumed + 1
			body, consumed := p.parse(tokens[i:])
			node.body = body
			i += consumed
			nodes = append(nodes, node)
		case "left_not_right":
			node, tokensConsumed := p.createLNR(tokens, i+1)
			i += tokensConsumed + 1
			body, consumed := p.parse(tokens[i:])
			node.body = body
			i += consumed
			nodes = append(nodes, node)
		case "print_integer":
			node, tokensConsumed := p.createStatement(tokens, i+1, "PRI")
			i += tokensConsumed + 1
			nodes = append(nodes, node)
		case "print_break":
			node, tokensConsumed := p.createSingleWordStatement(tokens, i+1, "BRK")
			i += tokensConsumed + 1
			nodes = append(nodes, node)
		case "print_ascii":
			node, tokensConsumed := p.createStatement(tokens, i+1, "PRU")
			i += tokensConsumed + 1
			nodes = append(nodes, node)
		case "increment_value":
			node, tokensConsumed := p.createStatement(tokens, i+1, "INC")
			i += tokensConsumed + 1
			nodes = append(nodes, node)
		case "close_block":
			i++
			return nodes, i
		case "string":
			node, tokensConsumed := p.createFunctionCall(tokens, i)
			i += tokensConsumed
			nodes = append(nodes, node)
		case "CHAR":
			node, tokensConsumed := p.createFunctionCall(tokens, i)
			i += tokensConsumed
			nodes = append(nodes, node)
			spew.Dump(node)
		default:
			spew.Dump(tokens[i])
			spew.Dump("what the fuck ?")
		}
	}

	return nodes, len(tokens)
}

func (p *parser) createSingleWordStatement(tokens []token, index int, t string) (*statement, int) {
	s := new(statement)
	tokensConsumed := 0

	s.lhs = t

	p.expect([]string{"SEMICOLON"}, tokens[index+tokensConsumed])
	tokensConsumed++

	return s, tokensConsumed
}

func (p *parser) createFunctionCall(tokens []token, index int) (*functionCall, int) {
	fc := new(functionCall)
	tokensConsumed := 0
	p.expect([]string{"string", "CHAR"}, tokens[index+tokensConsumed])
	fc.name = tokens[index+tokensConsumed].Value
	tokensConsumed++

	p.expect([]string{"LEFT_BRACKET"}, tokens[index+tokensConsumed])
	tokensConsumed++

	for currentToken := tokens[index+tokensConsumed]; currentToken.Type != "RIGHT_BRACKET"; currentToken = tokens[index+tokensConsumed] {
		if currentToken.Type == "COMMA" {
			p.expect([]string{"CHAR", "string", "NUMB"}, tokens[index+tokensConsumed+1])
			tokensConsumed++
			continue
		}
		fc.args = append(fc.args, currentToken.Value)
		tokensConsumed++
	}

	tokensConsumed++
	p.expect([]string{"SEMICOLON"}, tokens[index+tokensConsumed])
	tokensConsumed++
	return fc, tokensConsumed
}

func (p *parser) createStatement(tokens []token, index int, t string) (*statement, int) {
	s := new(statement)
	tokensConsumed := 0

	s.lhs = t

	p.expect([]string{"string", "CHAR", "NUMB"}, tokens[index+tokensConsumed])
	s.rhs = createLit(tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"SEMICOLON"}, tokens[index+tokensConsumed])
	tokensConsumed++

	return s, tokensConsumed
}

func createLit(token token) node {
	if token.Type == "NUMB" {
		nm := new(numLit)
		nm.value = token.Value
		return nm
	}
	sv := new(statVar)
	sv.value = token.Value
	return sv
}

func (p *parser) createLNR(tokens []token, index int) (*anb, int) {
	anb := new(anb)
	tokensConsumed := 0

	p.expect([]string{"LEFT_BRACKET"}, tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"CHAR", "NUMB", "string"}, tokens[index+tokensConsumed])

	anb.lhs = createLit(tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"COMMA"}, tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"CHAR", "NUMB", "string"}, tokens[index+tokensConsumed])
	anb.rhs = createLit(tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"RIGHT_BRACKET"}, tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"DOUBLE_DOT"}, tokens[index+tokensConsumed])
	tokensConsumed++

	return anb, tokensConsumed
}

func (p *parser) createFunctionHeader(tokens []token, index int) (*function, int) {
	f := new(function)
	tokensConsumed := 0
	p.expect([]string{"string", "CHAR"}, tokens[index+tokensConsumed])
	f.name = tokens[index+tokensConsumed].Value
	tokensConsumed++

	p.expect([]string{"LEFT_ARROW", "DOUBLE_DOT"}, tokens[index+tokensConsumed])
	if tokens[index+tokensConsumed].Type == "DOUBLE_DOT" {
		tokensConsumed++
		return f, tokensConsumed
	}
	tokensConsumed++
	for currentToken := tokens[index+tokensConsumed]; currentToken.Type != "RIGHT_ARROW"; currentToken = tokens[index+tokensConsumed] {
		p.expect([]string{"string", "CHAR", "COMMA"}, currentToken)
		if currentToken.Type == "COMMA" {
			p.expect([]string{"CHAR", "string"}, tokens[index+tokensConsumed+1])
			tokensConsumed++
			continue
		}
		f.params = append(f.params, currentToken.Value)
		tokensConsumed++
	}
	tokensConsumed++
	p.expect([]string{"DOUBLE_DOT"}, tokens[index+tokensConsumed])
	tokensConsumed++
	return f, tokensConsumed
}

func (p *parser) createVariable(tokens []token, index int) (*variable, int) {
	variable := new(variable)
	tokensConsumed := 0
	expectedNameTypes := []string{
		"CHAR",
		"string",
	}
	p.expect(expectedNameTypes, tokens[index+tokensConsumed])
	variable.name = tokens[index+tokensConsumed].Value
	tokensConsumed++

	expectedValueTypes := []string{"NUMB"}
	p.expect(expectedValueTypes, tokens[index+tokensConsumed])
	variable.value = tokens[index+tokensConsumed].Value
	tokensConsumed++

	p.expect([]string{"SEMICOLON"}, tokens[index+tokensConsumed])
	tokensConsumed++
	return variable, tokensConsumed
}

func (p *parser) expect(expectedValues []string, token token) {
	if !contains(token.Type, expectedValues) {
		throwSemanticError(&token, expectedValues, p.filename)
		os.Exit(65)
	}
}
