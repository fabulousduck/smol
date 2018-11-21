package smol

import (
	"io/ioutil"
	"os"

	"github.com/fabulousduck/smol/ast"
	"github.com/fabulousduck/smol/interpreter"
	"github.com/fabulousduck/smol/lexer"
)

//Smol : Defines the global attributes of the interpreter
type Smol struct {
	Tokens   []*lexer.Token
	HadError bool //TODO: use this
}

//NewSmol : Creates a new Smol instance
func NewSmol() *Smol {
	return new(Smol)
}

//RunFile : Interprets a given file
func (smol *Smol) RunFile(filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	smol.Run(string(file), filename)
	if smol.HadError {
		os.Exit(65)
	}
}

//Run exectues a given script
func (smol *Smol) Run(sourceCode string, filename string) {
	l := lexer.NewLexer(filename, sourceCode)
	l.Lex()
	p := ast.NewParser(filename, l.Tokens)
	//We can ignore the second return value here as it is the amount of tokens consumed.
	//We do not need this here
	p.Ast, _ = p.Parse()

	i := interpreter.NewInterpreter()
	i.Interpret(p.Ast)
}

//RunRepl is the same as Run except it allows you to pass our own interpreter so we can keep the state
//of it after the line has been executed
func (smol *Smol) RunRepl(sourceCode string, filename string, statefullInterpreter *interpreter.Interpreter) {
	l := lexer.NewLexer(filename, sourceCode)
	l.Lex()
	p := ast.NewParser(filename, l.Tokens)

	//Add an EOF token so semicolon errors dont index out of range
	p.Tokens = append(l.Tokens, lexer.Token{
		Value: "EOF",
		Type:  "EOF",
	})
	//We can ignore the second return value here as it is the amount of tokens consumed.
	//We do not need this here
	p.Ast, _ = p.Parse()
	statefullInterpreter.Interpret(p.Ast)
}
