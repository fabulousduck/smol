package losp

import "fmt"

type node struct {
	t    string
	body string
}

type parser struct {
	ast    []node
	expect []string
}

func NewParser() *parser {
	return new(parser)
}

func (p *parser) parse(tokens []token, filename string) {
	for i := 0; i < len(tokens); i++ {
		node := node{}
		if !contains(tokens[i].Type, p.expect) {
			report(
				tokens[i].Line,
				filename,
				fmt.Sprintf("expected one of [%s]. got %s",
					concatVariables(p.expect, ", "),
					tokens[i].Type))
		}
		switch tokens[i].Type {
		case "variable_assignment":
			node.t = "variable"
			p.expect = []string{"string", "CHAR"}
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
