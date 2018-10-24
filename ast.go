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
}

func (fc functionCall) getNodeName() string {
	return "functionCall"
}

type node interface {
	getNodeName() string
}

type parser struct {
	ast    []node
	expect []string
}

func NewParser() *parser {
	return new(parser)
}

func (p *parser) parse(tokens []token, filename string) {
	for i := 0; i < len(tokens); {
		var n node
		//build something like parse body that can be called recursively
		//and see the entire program as the main body
		switch tokens[i].Type {
		case "variable_assignment":
			nv, tokensConsumed := createVariable(tokens, i)
			i += tokensConsumed
			n = nv
		case "function_definition":
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
}

func createVariable(tokens []token, index int) (*variable, int) {
	variable := new(variable)
	tokensConsumed := 0
	expectedNameTypes := []string{
		"CHAR",
		"STRING",
	}
	if !contains(tokens[index+1].Type, expectedNameTypes) {
		throwSemanticError(&tokens[index+1], expectedNameTypes, "")
		os.Exit(65)
	}

	variable.name = tokens[index+1].Value
	tokensConsumed++

	expectedValueTypes := []string{"NUMB"}
	if !contains(tokens[index+2].Type, expectedValueTypes) {
		throwSemanticError(&tokens[index+2], expectedValueTypes, "")
		os.Exit(65)
	}
	variable.value = tokens[index+2].Value
	tokensConsumed++

	return variable, tokensConsumed
}
