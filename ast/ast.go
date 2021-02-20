package ast

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/fabulousduck/smol/errors"
	"github.com/fabulousduck/smol/lexer"
)

//Include is a means to include other source code files
type IncludeStatement struct {
	Name string
}

func (i IncludeStatement) GetNodeName() string {
	return "IncludeStatement"
}

//IfStatement is a conditional block that has an expression and a body
type IfStatement struct {
	Condition Node
	Body      []Node
}

func (i IfStatement) GetNodeName() string {
	return "IfStatement"
}

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
	Name            string
	Type            string
	Value           Node
	ValueExpression Expression
}

//SetStatement is used when a value needs to be set to a variable. Instructions that could make use of this are SET
type SetStatement struct {
	MHS Node
	RHS Node
}

func (ss SetStatement) GetNodeName() string {
	return "setStatement"
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
func (p *Parser) Parse(delim string) ([]Node, int) {
	nodes := []Node{}
	for p.TokensConsumed < len(p.Tokens) {
		if delim != "" && p.currentToken().Type == delim {
			p.advance()
			return nodes, p.TokensConsumed
		}
		currentTokenType := p.currentToken().Type
		switch currentTokenType {
		case "include":
			p.advance()
			//ResolveInclude is special because it includes the nodes from the file
			//to the AST itself
			p.resolveInclude()
		case "plot":
			p.advance()
			nodes = append(nodes, p.createPlot())
		case "variable_type":
			nodes = append(nodes, p.createVariable())
		case "function_definition":
			p.advance()
			nodes = append(nodes, p.createFunction())
		case "print":
			p.advance()
			nodes = append(nodes, p.createPrintCall())
		case "set_variable":
			p.advance()
			nodes = append(nodes, p.createSetStatement())
		case "if_statement":
			p.advance()
			nodes = append(nodes, p.createIfStatement())
		case "close_block":
			p.advance()
			return nodes, p.TokensConsumed
		//string and character by themselves can be either a function call or a direct operation on the variable such as a++s
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
			spew.Dump(p.currentToken())
			errors.UnknownTypeError()
		}

	}

	spew.Dump(nodes)
	return nodes, p.TokensConsumed
}

func (p *Parser) resolveInclude() {
	includeStatement := new(IncludeStatement)

	p.expectCurrent([]string{"string_litteral"})
	includeStatement.Name = p.currentToken().Value
	p.advance()

	var sb strings.Builder
	pwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("failed PWD call %s", err.Error()))
	}

	//get the path to the current dir
	sb.WriteString(pwd)
	sb.WriteString("/")

	entryFolderPathSplit := strings.Split(p.Filename, "/")
	if len(entryFolderPathSplit) > 1 {
		//pop filename off
		entryFolderPath := entryFolderPathSplit[:len(entryFolderPathSplit)-1]
		//add folder path
		sb.WriteString(strings.Join(entryFolderPath[:], "/"))
		sb.WriteString("/")
	}
	sb.WriteString(includeStatement.Name)

	//resolve the file
	file, err := ioutil.ReadFile(sb.String())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	libLexer := lexer.NewLexer(includeStatement.Name, string(file))
	libLexer.Lex()
	libParser := NewParser(includeStatement.Name, libLexer.Tokens)
	libParser.Ast, _ = libParser.Parse("")
	for i := 0; i < len(libParser.Ast); i++ {
		p.Ast = append(p.Ast, libParser.Ast[i])
	}
}

func (p *Parser) createIfStatement() *IfStatement {
	ifStatement := new(IfStatement)

	p.expectCurrent([]string{"left_parenthesis"})
	p.advance()

	condition, _ := p.readExpressionUntil([]string{")"})
	spew.Dump(condition)
	ifStatement.Condition = condition
	p.advance()

	p.expectCurrent([]string{"double_dot"})
	p.advance()

	ifBodyParser := NewParser(p.Filename, p.Tokens[p.TokensConsumed:])
	body, consumed := ifBodyParser.Parse("")
	ifStatement.Body = body
	p.advanceN(consumed)

	spew.Dump(ifStatement)

	return ifStatement
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

	p.expectCurrent([]string{"character", "string", "integer", "string_litteral"})
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

//TODO: make this into a parenthesized function
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
	body, consumed := eosParser.Parse("")
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
	body, consumed := switchParser.Parse("")
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

	body, consumed := switchParser.Parse("")
	st.Cases = body
	p.advanceN(consumed)

	return st
}

func (p *Parser) createFunctionCall() *FunctionCall {
	fc := new(FunctionCall)

	p.expectCurrent([]string{"string", "character"})
	fc.Name = p.currentToken().Value
	p.advance()

	p.expectCurrent([]string{"left_parenthesis"})
	p.advance()

	if p.currentToken().Type == "right_parenthesis" {
		p.advance()
		return fc
	}

	for currentToken := p.currentToken(); currentToken.Type != "right_parenthesis"; currentToken = p.currentToken() {
		exprList, delimFound := p.readExpressionUntil([]string{",", ")"})
		fc.Args = append(fc.Args, exprList)
		if delimFound == "," {
			p.advance()
			continue
		}
		p.advance()
		break
	}

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
	body, consumed := functionParser.Parse("")
	f.Body = body
	p.advanceN(consumed)
	return f
}
