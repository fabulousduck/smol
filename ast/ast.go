package ast

import (
	"fmt"
	"os"

	"github.com/fabulousduck/proto/src/types"
	"github.com/fabulousduck/smol/lexer"
)

//Node is a wrapper interface that AST nodes can implement
type Node interface {
	GetNodeName() string //GetNodeName Gets the identifier of a AST node describing what it is

}

//NumLit represents a numeric litteral.
type NumLit struct {
	Value string
}

func (nm NumLit) GetNodeName() string {
	return "numLit"
}

//Eos is a special node in a switch statement that is called if defined when no cases match the given value
type Eos struct {
	Body []Node
}

func (eos Eos) GetNodeName() string {
	return "end_of_switch"
}

//SwitchCase is a block definiton that is run when the MatchValue is matched
type SwitchCase struct {
	MatchValue Node
	Body       []Node
}

func (sc SwitchCase) GetNodeName() string {
	return "switchCase"
}

//SwitchStatement matches the Matchcase against all cases defined in Cases of the switchcase
//if one matches, the body of that case will be executed.
//if a EOS is defined within the body, the EOS body will be run if no case matches the matchvalue
type SwitchStatement struct {
	MatchValue Node
	Cases      []Node
}

func (st SwitchStatement) GetNodeName() string {
	return "switchStatement"
}

//StatVar contains the value of a static variable.
//These statVars are used when a variable is being referenced
//where the Value is the name of the variable referenced
type StatVar struct {
	Value string
}

func (sv StatVar) GetNodeName() string {
	return "statVar"
}

//Function is a standard function definition containing the name, parameters and body of the function
type Function struct {
	Name   string
	Params []string
	Body   []Node
}

func (f Function) GetNodeName() string {
	return "function"
}

//Variable is a construct used to create a new variable.
//This is the struct that will be pushed to the stack
type Variable struct {
	Name  string
	Value Node
}

//Anb is the while loop of smol. It will keep executing its body until LHS equals RHS
type Anb struct {
	LHS  Node
	RHS  Node
	Body []Node
}

func (anb Anb) GetNodeName() string {
	return "anb"
}

//Comparison is an expression type that is used when an operation like GT, LT, EQ or NEQ is called
type Comparison struct {
	Operator string
	LHS      Node
	RHS      Node
	Body     []Node
}

func (c Comparison) GetNodeName() string {
	return "comparison"
}

//SetStatement is used when a value needs to be set to a variable. Instructions that could make use of this are SET
type SetStatement struct {
	LHS string
	MHS Node
	RHS Node
}

func (ss SetStatement) GetNodeName() string {
	return "setStatement"
}

//MathStatement contains info needed to execute a mathematical statement like ADD, SUB, MUL and DIV
type MathStatement struct {
	LHS string
	MHS Node
	RHS Node
}

func (ms MathStatement) GetNodeName() string {
	return "mathStatement"
}

//Statement is a general statement container for all other statements that do not fall under math and logic for example MEM
type Statement struct {
	LHS string
	RHS Node
}

func (s Statement) GetNodeName() string {
	return "statement"
}

func (v Variable) GetNodeName() string {
	return "variable"
}

//FunctionCall specifies a function call and the arguments given
type FunctionCall struct {
	Name string
	Args []string
}

func (fc FunctionCall) GetNodeName() string {
	return "functionCall"
}

//Parser contains the final AST and forms a base for all ast generating functions
type Parser struct {
	Ast      []Node
	Filename string
}

//NewParser returns a new Parser instance with the given file
func NewParser(filename string) *Parser {
	p := new(Parser)
	p.Filename = filename
	return p
}

//Parse takes a set of tokens and generates an AST from them
func (p *Parser) Parse(tokens []lexer.Token) ([]Node, int) {
	nodes := []Node{}
	for i := 0; i < len(tokens); {
		switch tokens[i].Type {
		case "variable_assignment":
			node, tokensConsumed := p.createVariable(tokens, i+1)
			i += tokensConsumed + 1
			nodes = append(nodes, node)
		case "function_definition":
			node, tokensConsumed := p.createFunctionHeader(tokens, i+1)
			i += tokensConsumed + 1
			body, consumed := p.Parse(tokens[i:])
			node.Body = body
			i += consumed
			nodes = append(nodes, node)
		case "left_not_right":
			node, tokensConsumed := p.createLNR(tokens, i+1)
			i += tokensConsumed + 1
			body, consumed := p.Parse(tokens[i:])
			node.Body = body
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
		case "set_variable":
			node, tokensConsumed := p.createSetStatement(tokens, i+1, "SET")
			i += tokensConsumed + 1
			nodes = append(nodes, node)
		case "addition":
			fallthrough
		case "subtraction":
			fallthrough
		case "power_of":
			fallthrough
		case "multiplication":
			fallthrough
		case "division":
			node, tokensConsumed := p.createMathStatement(tokens, i)
			i += tokensConsumed
			nodes = append(nodes, node)
		case "close_block":
			i++
			return nodes, i
		case "string":
			node, tokensConsumed := p.createFunctionCall(tokens, i)
			i += tokensConsumed
			nodes = append(nodes, node)
		case "character":
			node, tokensConsumed := p.createFunctionCall(tokens, i)
			i += tokensConsumed
			nodes = append(nodes, node)
		case "equals":
			fallthrough
		case "not_equals":
			fallthrough
		case "less_than":
			fallthrough
		case "greater_than":
			node, tokensConsumed := p.createComparisonHeader(tokens, i)
			i += tokensConsumed
			body, consumed := p.Parse(tokens[i:])
			node.Body = body
			i += consumed
			nodes = append(nodes, node)
		case "switch":
			node, tokensConsumed := p.createSwitchStatement(tokens, i+1)
			i += tokensConsumed
			nodes = append(nodes, node)
		case "case":
			node, tokensConsumed := p.createSwitchCase(tokens, i+1)
			nodes = append(nodes, node)
			i += tokensConsumed
		case "end_of_switch":
			node, tokensConsumed := p.createEOSStatement(tokens, i+1)
			nodes = append(nodes, node)
			i += tokensConsumed
		default:
			//TODO: make an error for this
			fmt.Println("Unknown token type found.")
		}
	}

	return nodes, len(tokens)
}

func (p *Parser) createEOSStatement(tokens []lexer.Token, index int) (*Eos, int) {
	eos := new(Eos)
	tokensConsumed := 0

	p.expect([]string{"double_dot"}, tokens[index+tokensConsumed])
	tokensConsumed++

	body, consumed := p.Parse(tokens[index+tokensConsumed:])
	eos.Body = body
	tokensConsumed += consumed + 1

	return eos, tokensConsumed
}

func (p *Parser) createSwitchCase(tokens []lexer.Token, index int) (*SwitchCase, int) {
	sc := new(SwitchCase)
	tokensConsumed := 0

	p.expect([]string{"character", "string", "integer"}, tokens[index+tokensConsumed])
	sc.MatchValue = createLit(tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"double_dot"}, tokens[index+tokensConsumed])
	tokensConsumed++

	body, consumed := p.Parse(tokens[index+tokensConsumed:])
	sc.Body = body
	tokensConsumed += consumed + 1

	return sc, tokensConsumed
}

func (p *Parser) createSwitchStatement(tokens []lexer.Token, index int) (*SwitchStatement, int) {
	st := new(SwitchStatement)
	tokensConsumed := 0

	p.expect([]string{"left_bracket"}, tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"character", "string", "integer"}, tokens[index+tokensConsumed])
	st.MatchValue = createLit(tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"right_bracket"}, tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"double_dot"}, tokens[index+tokensConsumed])
	tokensConsumed++

	body, consumed := p.Parse(tokens[index+tokensConsumed:])
	st.Cases = body
	tokensConsumed += consumed

	return st, tokensConsumed
}

func (p *Parser) createComparisonHeader(tokens []lexer.Token, index int) (*Comparison, int) {
	ch := new(Comparison)
	tokensConsumed := 0

	ch.Operator = tokens[index+tokensConsumed].Value
	tokensConsumed++

	p.expect([]string{"left_bracket"}, tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"character", "string", "integer"}, tokens[index+tokensConsumed])
	ch.LHS = createLit(tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"comma"}, tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"character", "string", "integer"}, tokens[index+tokensConsumed])
	ch.RHS = createLit(tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"right_bracket"}, tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"double_dot"}, tokens[index+tokensConsumed])
	tokensConsumed++

	return ch, tokensConsumed
}

func (p *Parser) createMathStatement(tokens []lexer.Token, index int) (Node, int) {
	ms := new(MathStatement)
	tokensConsumed := 0

	ms.LHS = tokens[index].Value
	tokensConsumed++

	p.expect([]string{"character", "string"}, tokens[index+tokensConsumed])
	ms.MHS = createLit(tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"character", "string", "integer"}, tokens[index+tokensConsumed])
	ms.RHS = createLit(tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"semicolon"}, tokens[index+tokensConsumed])
	tokensConsumed++

	return ms, tokensConsumed
}

func (p *Parser) createSingleWordStatement(tokens []lexer.Token, index int, t string) (*Statement, int) {
	s := new(Statement)
	tokensConsumed := 0

	s.LHS = t
	p.expect([]string{"semicolon"}, tokens[index+tokensConsumed])
	tokensConsumed++

	return s, tokensConsumed
}

func (p *Parser) createFunctionCall(tokens []lexer.Token, index int) (*FunctionCall, int) {
	fc := new(FunctionCall)
	tokensConsumed := 0
	p.expect([]string{"string", "character"}, tokens[index+tokensConsumed])
	fc.Name = tokens[index+tokensConsumed].Value
	tokensConsumed++

	p.expect([]string{"left_bracket"}, tokens[index+tokensConsumed])
	tokensConsumed++

	for currentToken := tokens[index+tokensConsumed]; currentToken.Type != "right_bracket"; currentToken = tokens[index+tokensConsumed] {
		if currentToken.Type == "comma" {
			p.expect([]string{"character", "string", "integer"}, tokens[index+tokensConsumed+1])
			tokensConsumed++
			continue
		}
		fc.Args = append(fc.Args, currentToken.Value)
		tokensConsumed++
	}

	tokensConsumed++
	p.expect([]string{"semicolon"}, tokens[index+tokensConsumed])
	tokensConsumed++
	return fc, tokensConsumed
}

func (p *Parser) createSetStatement(tokens []lexer.Token, index int, t string) (*SetStatement, int) {
	ss := new(SetStatement)
	tokensConsumed := 0

	ss.LHS = t

	p.expect([]string{"character", "string"}, tokens[index+tokensConsumed])
	ss.MHS = createLit(tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"integer", "character", "string"}, tokens[index+tokensConsumed])
	ss.RHS = createLit(tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"semicolon"}, tokens[index+tokensConsumed])
	tokensConsumed++

	return ss, tokensConsumed
}

func (p *Parser) createStatement(tokens []lexer.Token, index int, t string) (*Statement, int) {
	s := new(Statement)
	tokensConsumed := 0

	s.LHS = t

	p.expect([]string{"string", "character", "integer"}, tokens[index+tokensConsumed])
	s.RHS = createLit(tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"semicolon"}, tokens[index+tokensConsumed])
	tokensConsumed++

	return s, tokensConsumed
}

func createLit(token lexer.Token) Node {
	if token.Type == "integer" {
		nm := new(NumLit)
		nm.Value = token.Value
		return nm
	}
	sv := new(StatVar)
	sv.Value = token.Value
	return sv
}

func (p *Parser) createLNR(tokens []lexer.Token, index int) (*Anb, int) {
	anb := new(Anb)
	tokensConsumed := 0

	p.expect([]string{"left_bracket"}, tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"character", "integer", "string"}, tokens[index+tokensConsumed])

	anb.LHS = createLit(tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"comma"}, tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"character", "integer", "string"}, tokens[index+tokensConsumed])
	anb.RHS = createLit(tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"right_bracket"}, tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"double_dot"}, tokens[index+tokensConsumed])
	tokensConsumed++

	return anb, tokensConsumed
}

func (p *Parser) createFunctionHeader(tokens []lexer.Token, index int) (*Function, int) {
	f := new(Function)
	tokensConsumed := 0
	p.expect([]string{"string", "character"}, tokens[index+tokensConsumed])
	f.Name = tokens[index+tokensConsumed].Value
	tokensConsumed++

	p.expect([]string{"left_arrow", "double_dot"}, tokens[index+tokensConsumed])
	if tokens[index+tokensConsumed].Type == "double_dot" {
		tokensConsumed++
		return f, tokensConsumed
	}
	tokensConsumed++
	for currentToken := tokens[index+tokensConsumed]; currentToken.Type != "right_arrow"; currentToken = tokens[index+tokensConsumed] {
		p.expect([]string{"string", "character", "comma"}, currentToken)
		if currentToken.Type == "comma" {
			p.expect([]string{"character", "string"}, tokens[index+tokensConsumed+1])
			tokensConsumed++
			continue
		}
		f.Params = append(f.Params, currentToken.Value)
		tokensConsumed++
	}
	tokensConsumed++
	p.expect([]string{"double_dot"}, tokens[index+tokensConsumed])
	tokensConsumed++
	return f, tokensConsumed
}

func (p *Parser) createVariable(tokens []lexer.Token, index int) (*Variable, int) {
	variable := new(Variable)
	tokensConsumed := 0
	expectedNameTypes := []string{
		"character",
		"string",
	}

	p.expect(expectedNameTypes, tokens[index+tokensConsumed])
	variable.Name = tokens[index+tokensConsumed].Value
	tokensConsumed++

	expectedValueTypes := []string{"integer", "character", "string"}
	p.expect(expectedValueTypes, tokens[index+tokensConsumed])
	variable.Value = createLit(tokens[index+tokensConsumed])
	tokensConsumed++

	p.expect([]string{"semicolon"}, tokens[index+tokensConsumed])
	tokensConsumed++
	return variable, tokensConsumed
}

func (p *Parser) expect(expectedValues []string, token lexer.Token) {
	if !types.Contains(token.Type, expectedValues) {
		lexer.ThrowSemanticError(&token, expectedValues, p.Filename)
		os.Exit(65)
	}
}
