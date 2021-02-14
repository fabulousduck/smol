package ast

import (
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/fabulousduck/proto/src/types"
	"github.com/fabulousduck/smol/errors"
	"github.com/fabulousduck/smol/lexer"
)

//PlotStatement is a statement that contains all info needed to draw a pixel to the screen
type PlotStatement struct {
	X, Y Node
}

func (p PlotStatement) GetNodeName() string {
	return "plotStatement"
}

//FreeStatement is an instruction that frees a variable from the stack
type FreeStatement struct {
	Variable Node
}

func (r FreeStatement) GetNodeName() string {
	return "freeStatement"
}

type DirectOperation struct {
	Variable  Node
	Operation string
}

func (do DirectOperation) GetNodeName() string {
	return "directOperation"
}

//Node is a wrapper interface that AST nodes can implement
type Node interface {
	GetNodeName() string //GetNodeName Gets the identifier of a AST node describing what it is

}

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
	Type  string
	Value Node
}

//WhileNot is the while loop of smol. It will keep executing its body until LHS equals RHS
type WhileNot struct {
	LHS  Node
	RHS  Node
	Body []Node
}

func (whileNot WhileNot) GetNodeName() string {
	return "whileNot"
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

type UseStatement struct {
	name string
}

func (us UseStatement) GetNodeName() string {
	return "useStatement"
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

//PrintCall specifies a call to the inbuilt print function
type PrintCall struct {
	Printable Node
}

func (pc PrintCall) GetNodeName() string {
	return "printCall"
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
		case "plot":
			p.advance()
			nodes = append(nodes, p.createPlot())
		case "use":
			p.advance()
			nodes = append(nodes, p.createUse())
		case "variable_type":
			nodes = append(nodes, p.createVariable())
		case "function_definition":
			p.advance()
			nodes = append(nodes, p.createFunction())
		case "while_not":
			p.advance()
			nodes = append(nodes, p.createWhileNot())
		case "print":
			p.advance()
			nodes = append(nodes, p.createPrintCall())
		case "set_variable":
			p.advance()
			nodes = append(nodes, p.createSetStatement())
		case "close_block":
			p.advance()
			return nodes, p.TokensConsumed
		//string and character loose can be either a function call or a direct operation on the variable such as a++s
		case "string":
			fallthrough
		case "character":
			//its either a function
			if p.nextToken().Type == "left_parenthesis" {
				nodes = append(nodes, p.createFunctionCall())
				//or a direct operation
			} else {
				nodes = append(nodes, p.createDirectOperation())
			}
		case "equals":
			fallthrough
		case "not_equals":
			fallthrough
		case "less_than":
			fallthrough
		case "greater_than":
			nodes = append(nodes, p.createComparison())
		case "free":
			p.advance()
			nodes = append(nodes, p.createFreeStatement())
		case "switch":
			p.advance()
			nodes = append(nodes, p.createSwitchStatement())
			p.advance()
		case "case":
			p.advance()
			nodes = append(nodes, p.createSwitchCase())
		case "end_of_switch":
			p.advance()
			nodes = append(nodes, p.createEOSStatement())
		case "end_of_file":
			return nodes, p.TokensConsumed
		default:
			errors.UnknownTypeError()
		}

	}
	return nodes, p.TokensConsumed
}

func (p *Parser) createDirectOperation() *DirectOperation {
	do := new(DirectOperation)

	do.Variable = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"direct_variable_operation"})
	do.Operation = p.currentToken().Value
	p.advance()

	return do
}

func (p *Parser) createPrintCall() *PrintCall {
	pc := new(PrintCall)

	p.expectCurrent([]string{"left_parenthesis"})
	p.advance()

	p.expectCurrent([]string{"character", "string", "integer"})
	pc.Printable = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"right_parenthesis"})
	p.advance()
	return pc
}

func (p *Parser) createPlot() *PlotStatement {
	ps := new(PlotStatement)

	p.expectCurrent([]string{"left_parenthesis"})
	p.advance()

	p.expectCurrent([]string{"character", "string", "integer"})
	ps.X = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"comma"})
	p.advance()

	p.expectCurrent([]string{"character", "string", "integer"})
	ps.Y = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"right_parenthesis"})
	p.advance()

	return ps
}

func (p *Parser) createUse() *UseStatement {
	us := new(UseStatement)

	p.expectCurrent([]string{"string"})
	us.name = p.currentToken().Value
	p.advance()

	p.expectCurrent([]string{"semicolon"})
	p.advance()

	return us
}

func (p *Parser) createFreeStatement() *FreeStatement {
	r := new(FreeStatement)

	p.expectCurrent([]string{"string", "character"})
	r.Variable = createLit(p.currentToken())
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

	p.expectCurrent([]string{"left_parenthesis"})
	p.advance()

	p.expectCurrent([]string{"character", "string", "integer"})
	st.MatchValue = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"right_parenthesis"})
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

	p.expectCurrent([]string{"left_parenthesis"})
	p.advance()

	p.expectCurrent([]string{"character", "string", "integer"})
	ch.LHS = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"comma"})
	p.advance()

	p.expectCurrent([]string{"character", "string", "integer"})
	ch.RHS = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"right_parenthesis"})
	p.advance()

	p.expectCurrent([]string{"double_dot"})
	p.advance()

	comparisonParser := NewParser(p.Filename, p.Tokens[p.TokensConsumed:])

	body, consumed := comparisonParser.Parse()
	ch.Body = body
	p.advanceN(consumed)

	return ch
}

func (p *Parser) createFunctionCall() *FunctionCall {
	fc := new(FunctionCall)

	p.expectCurrent([]string{"string", "character"})
	fc.Name = p.currentToken().Value
	p.advance()

	p.expectCurrent([]string{"left_parenthesis"})
	p.advance()

	for currentToken := p.currentToken(); currentToken.Type != "right_parenthesis"; currentToken = p.currentToken() {
		if currentToken.Type == "comma" {
			p.expectNext([]string{"character", "string", "integer"})
			p.advance()
			continue
		}
		fc.Args = append(fc.Args, createLit(currentToken))
		p.advance()
	}

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

func (p *Parser) createWhileNot() *WhileNot {
	whileNot := new(WhileNot)

	p.expectCurrent([]string{"left_parenthesis"})
	p.advance()

	p.expectCurrent([]string{"character", "integer", "string"})

	whileNot.LHS = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"comma"})
	p.advance()

	p.expectCurrent([]string{"character", "integer", "string"})
	whileNot.RHS = createLit(p.currentToken())
	p.advance()

	p.expectCurrent([]string{"right_parenthesis"})
	p.advance()

	p.expectCurrent([]string{"double_dot"})
	p.advance()

	whileNotParser := NewParser(p.Filename, p.Tokens[p.TokensConsumed:])

	body, consumed := whileNotParser.Parse()
	whileNot.Body = body
	p.advanceN(consumed)

	return whileNot
}

func (p *Parser) createFunction() *Function {
	f := new(Function)

	p.expectCurrent([]string{"string", "character"})
	f.Name = p.currentToken().Value
	p.advance()

	p.expectCurrent([]string{"left_parenthesis"})
	p.advance()

	for currentToken := p.currentToken(); currentToken.Type != "right_parenthesis"; currentToken = p.currentToken() {
		p.expectCurrent([]string{"string", "character", "comma"})

		//we expect there to be another parameter when we see a comma
		if currentToken.Type == "comma" {
			p.expectNext([]string{"character", "string"})
			p.advance()
			continue
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

	switch variable.Type {
	case "Bool":
		p.expectCurrent([]string{"boolean_keyword"})
	case "String":
		p.expectCurrent([]string{"string_litteral"})
	case "Uint32":
		p.expectCurrent([]string{"integer"})
	case "Uint64":
		p.expectCurrent([]string{"integer"})
	default:
		errors.UnknownVariableTypeError(variable.Type)
	}

	variable.Value = createLit(p.currentToken())
	p.advance()

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

func (p *Parser) expectCurrent(expectedValues []string) {
	currentToken := p.currentToken()

	if !types.Contains(currentToken.Type, expectedValues) {
		spew.Dump(currentToken)
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
