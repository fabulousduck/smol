package losp

import "os"

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

func (v variable) getNodeName() string {
	return "variable"
}

type functionCall struct {
	name   string
	params []string
	body   []node
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

func (p *parser) parse(tokens []token) []node {
	nodes := []node{}
	for i := 0; i < len(tokens); {
		//build something like parse body that can be called recursively
		//and see the entire program as the main body
		switch tokens[i].Type {
		case "variable_assignment":
			node, tokensConsumed := p.createVariable(tokens, i)
			i += tokensConsumed
			nodes = append(nodes, node)
		case "function_definition":
			node, tokensConsumed := p.createFunctionHeader(tokens, i+1)
			i += tokensConsumed
			node.body = p.parse(tokens[:i])
		case "left_not_right":
		case "print_statement":
		case "increment_value":
		case "string":
		case "CHAR":
		case "NUMB":
		case "LEFT_BRACE":
		case "RIGHT_BRACE":
		case "LEFT_ARROW":
		case "RIGHT_ARROW":
		case "DOUBLE_DOT":
		case "COMMA":
		case "SEMI_COLON":
		}
	}

	return nodes
}

func (p *parser) createFunctionHeader(tokens []token, index int) (*function, int) {
	f := new(function)
	tokensConsumed := 0
	p.expect([]string{"string", "CHAR"}, tokens[index+1])
	f.name = tokens[index+tokensConsumed].Value
	tokensConsumed++

	p.expect([]string{"LEFT_ARROW", "DOUBLE_DOT"}, tokens[index+2])
	if tokens[index+tokensConsumed].Type == "DOUBLE_DOT" {
		tokensConsumed++
		return f, tokensConsumed
	}
	for currentToken := tokens[index+tokensConsumed]; currentToken.Type != "RIGHT_ARROW"; currentToken = tokens[index+tokensConsumed] {
		p.expect([]string{"string", "CHAR", "COMMA"}, currentToken)
		if currentToken.Type == "COMMA" {
			p.expect([]string{"char", "string"}, tokens[index+tokensConsumed+1])
			tokensConsumed++
			continue
		}
		f.params = append(f.params, currentToken.Value)
		tokensConsumed++
	}
	return f, tokensConsumed
}

func (p *parser) createVariable(tokens []token, index int) (*variable, int) {
	variable := new(variable)
	tokensConsumed := 0
	expectedNameTypes := []string{
		"CHAR",
		"STRING",
	}
	p.expect(expectedNameTypes, tokens[index+1])
	variable.name = tokens[index+1].Value
	tokensConsumed++

	expectedValueTypes := []string{"NUMB"}
	p.expect(expectedValueTypes, tokens[index+2])
	variable.value = tokens[index+2].Value
	tokensConsumed++

	return variable, tokensConsumed
}

func (p *parser) expect(expectedValues []string, token token) {
	if !contains(token.Type, expectedValues) {
		throwSemanticError(&token, expectedValues, p.filename)
		os.Exit(65)
	}
}
