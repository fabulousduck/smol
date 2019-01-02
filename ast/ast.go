package ast

import (
	"fmt"
	"os"

	"github.com/fabulousduck/proto/src/types"
	"github.com/fabulousduck/smol/lexer"
)

//ReleaseStatement is an instruction that frees a variable from the stack
type ReleaseStatement struct {
	Variable Node
}

func (r ReleaseStatement) GetNodeName() string {
	return "releaseStatement"
}

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
	Args []Node
}

func (fc FunctionCall) GetNodeName() string {
	return "functionCall"
}

//Parser contains the final AST and forms a base for all ast generating functions
type Parser struct {
	Tokens         []lexer.Token
	Ast            []Node
	Filename       string
	TokensConsumed int
}

//NewParser returns a new Parser instance with the given file
func NewParser(filename string, tokens []lexer.Token) *Parser {
	p := new(Parser)
	p.Filename = filename
	p.TokensConsumed = 0
	p.Tokens = tokens
	return p
}

//Parse takes a set of tokens and generates an AST from them
func (p *Parser) Parse() ([]Node, int) {
	nodes := []Node{}
	for p.TokensConsumed < len(p.Tokens) {
		switch p.currentToken().Type {
		case "variable_assignment":
			p.advance()
			nodes = append(nodes, p.createVariable())
		case "function_definition":
			p.advance()
			nodes = append(nodes, p.createFunction())
		case "left_not_right":
			p.advance()
			nodes = append(nodes, p.createLNR())
		case "print_integer":
			p.advance()
			nodes = append(nodes, p.createStatement("PRI"))
		case "print_break":
			p.advance()
			nodes = append(nodes, p.createSingleWordStatement("BRK"))
		case "print_ascii":
			p.advance()
			nodes = append(nodes, p.createStatement("PRU"))
		case "increment_value":
			p.advance()
			nodes = append(nodes, p.createStatement("INC"))
		case "set_variable":
			p.advance()
			nodes = append(nodes, p.createSetStatement())
		case "addition":
			fallthrough
		case "subtraction":
			fallthrough
		case "power_of":
			fallthrough
		case "multiplication":
			fallthrough
		case "division":
			nodes = append(nodes, p.createMathStatement())
		case "close_block":
			p.advance()
			return nodes, p.TokensConsumed
		case "string": //we are allowed to assume a character or string means a function call since every thing else gets tagged as either a variable or statement
			fallthrough
		case "character":
			nodes = append(nodes, p.createFunctionCall())
		case "equals":
			fallthrough
		case "not_equals":
			fallthrough
		case "less_than":
			fallthrough
		case "greater_than":
			nodes = append(nodes, p.createComparison())
		case "release":
			p.advance()
			nodes = append(nodes, p.createReleaseStatement())
		case "switch":
			p.advance()
			nodes = append(nodes, p.createSwitchStatement())
		case "case":
			p.advance()
			nodes = append(nodes, p.createSwitchCase())
		case "end_of_switch":
			p.advance()
			nodes = append(nodes, p.createEOSStatement())
		case "end_of_file":
			return nodes, p.TokensConsumed
		default:
			//TODO: make an error for this
			fmt.Println("Unknown token type found.")
		}

	}

	return nodes, p.TokensConsumed
}

func (p *Parser) createReleaseStatement() *ReleaseStatement {
	r := new(ReleaseStatement)

	p.expectCurrent([]string{"string", "character"})
	r.Variable = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"semicolon"})
	p.advance()

	return r
}

func (p *Parser) createEOSStatement() *Eos {
	eos := new(Eos)

	p.expectCurrent([]string{"double_dot"})
	p.advance()

	eosParser := NewParser(p.Filename, p.Tokens[p.TokensConsumed:])
	body, consumed := eosParser.Parse()
	eos.Body = body

	p.advanceN(consumed)

	return eos
}

func (p *Parser) createSwitchCase() *SwitchCase {
	sc := new(SwitchCase)

	p.expectCurrent([]string{"character", "string", "integer"})
	sc.MatchValue = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"double_dot"})
	p.advance()

	switchParser := NewParser(p.Filename, p.Tokens[p.TokensConsumed:])
	body, consumed := switchParser.Parse()
	sc.Body = body

	p.advanceN(consumed)

	return sc
}

func (p *Parser) createSwitchStatement() *SwitchStatement {
	st := new(SwitchStatement)

	p.expectCurrent([]string{"left_bracket"})
	p.advance()

	p.expectCurrent([]string{"character", "string", "integer"})
	st.MatchValue = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"right_bracket"})
	p.advance()

	p.expectCurrent([]string{"double_dot"})
	p.advance()

	switchParser := NewParser(p.Filename, p.Tokens[p.TokensConsumed:])

	body, consumed := switchParser.Parse()
	st.Cases = body
	p.advanceN(consumed)

	return st
}

func (p *Parser) createComparison() *Comparison {
	ch := new(Comparison)

	ch.Operator = p.currentToken().Value
	p.advance()

	p.expectCurrent([]string{"left_bracket"})
	p.advance()

	p.expectCurrent([]string{"character", "string", "integer"})
	ch.LHS = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"comma"})
	p.advance()

	p.expectCurrent([]string{"character", "string", "integer"})
	ch.RHS = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"right_bracket"})
	p.advance()

	p.expectCurrent([]string{"double_dot"})
	p.advance()

	comparisonParser := NewParser(p.Filename, p.Tokens[p.TokensConsumed:])

	body, consumed := comparisonParser.Parse()
	ch.Body = body
	p.advanceN(consumed)

	return ch
}

func (p *Parser) createMathStatement() *MathStatement {
	ms := new(MathStatement)

	ms.LHS = p.currentToken().Value
	p.advance()

	p.expectCurrent([]string{"character", "string"})
	ms.MHS = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"character", "string", "integer"})
	ms.RHS = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"semicolon"})
	p.advance()

	return ms
}

func (p *Parser) createSingleWordStatement(lhs string) *Statement {
	s := new(Statement)

	s.LHS = lhs
	p.expectCurrent([]string{"semicolon"})
	p.advance()

	return s
}

func (p *Parser) createFunctionCall() *FunctionCall {
	fc := new(FunctionCall)

	p.expectCurrent([]string{"string", "character"})
	fc.Name = p.currentToken().Value
	p.advance()

	p.expectCurrent([]string{"left_bracket"})
	p.advance()

	for currentToken := p.currentToken(); currentToken.Type != "right_bracket"; currentToken = p.currentToken() {
		if currentToken.Type == "comma" {
			p.expectNext([]string{"character", "string", "integer"})
			p.advance()
			continue
		}

		fc.Args = append(fc.Args, createLit(currentToken))
		p.advance()
	}

	//Todo: figure out why this advance needs to be here
	p.advance()
	p.expectCurrent([]string{"semicolon"})
	p.advance()
	return fc
}

func (p *Parser) createSetStatement() *SetStatement {
	ss := new(SetStatement)

	p.expectCurrent([]string{"character", "string"})
	ss.MHS = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"integer", "character", "string"})
	ss.RHS = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"semicolon"})
	p.advance()

	return ss
}

func (p *Parser) createStatement(lhs string) *Statement {
	s := new(Statement)

	s.LHS = lhs

	p.expectCurrent([]string{"string", "character", "integer"})
	s.RHS = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"semicolon"})
	p.advance()

	return s
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

func (p *Parser) createLNR() *Anb {
	anb := new(Anb)

	p.expectCurrent([]string{"left_bracket"})
	p.advance()

	p.expectCurrent([]string{"character", "integer", "string"})

	anb.LHS = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"comma"})
	p.advance()

	p.expectCurrent([]string{"character", "integer", "string"})
	anb.RHS = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"right_bracket"})
	p.advance()

	p.expectCurrent([]string{"double_dot"})
	p.advance()

	anbParser := NewParser(p.Filename, p.Tokens[p.TokensConsumed:])

	body, consumed := anbParser.Parse()
	anb.Body = body
	p.advanceN(consumed)

	return anb
}

func (p *Parser) createFunction() *Function {
	f := new(Function)

	p.expectCurrent([]string{"string", "character"})
	f.Name = p.currentToken().Value
	p.advance()

	p.expectCurrent([]string{"left_arrow", "double_dot"})
	if p.currentToken().Type == "double_dot" {
		p.advance()
		return f
	}
	p.advance()

	for currentToken := p.currentToken(); currentToken.Type != "right_arrow"; currentToken = p.currentToken() {
		p.expectCurrent([]string{"string", "character", "comma"})

		//we expect there to be another parameter when we see a comma
		if currentToken.Type == "comma" {
			p.expectNext([]string{"character", "string"})
			p.advance()
			continue
		}

		if currentToken.Type == "string" || currentToken.Type == "character" {

		}

		f.Params = append(f.Params, currentToken.Value)
		p.advance()

	}

	p.advance()
	p.expectCurrent([]string{"double_dot"})
	p.advance()

	functionParser := NewParser(p.Filename, p.Tokens[p.TokensConsumed:])

	body, consumed := functionParser.Parse()
	f.Body = body
	p.advanceN(consumed)

	return f
}

func (p *Parser) createVariable() *Variable {
	variable := new(Variable)

	p.expectCurrent([]string{"character", "string"})
	variable.Name = p.currentToken().Value
	p.advance()

	p.expectCurrent([]string{"integer", "character", "string"})
	variable.Value = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"semicolon"})
	p.advance()
	return variable
}

func (p *Parser) expectCurrent(expectedValues []string) {
	currentToken := p.currentToken()
	if !types.Contains(currentToken.Type, expectedValues) {

		lexer.ThrowSemanticError(&currentToken, expectedValues, p.Filename)
		os.Exit(65)
	}
}

func (p *Parser) expectNext(expectedValues []string) {
	nextToken := p.nextToken()
	if !types.Contains(nextToken.Type, expectedValues) {
		lexer.ThrowSemanticError(&nextToken, expectedValues, p.Filename)
		os.Exit(65)
	}
}

func (p *Parser) currentToken() lexer.Token {
	return p.Tokens[p.TokensConsumed]
}

func (p *Parser) nextToken() lexer.Token {
	return p.Tokens[p.TokensConsumed+1]
}

func (p *Parser) advance() {
	p.TokensConsumed++
}

func (p *Parser) advanceN(n int) {
	p.TokensConsumed += n
}

//NodeIsVariable allows for nice statements like if NodeIsVariable(node) {}
func NodeIsVariable(node Node) bool {
	return node.GetNodeName() == "statVar"
}
