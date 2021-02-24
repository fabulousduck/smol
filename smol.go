package smol

import (
	"io/ioutil"
	"os"

	"github.com/fabulousduck/smol/bytecode"

	"github.com/fabulousduck/smol/ast"
	"github.com/fabulousduck/smol/ir"
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
func (smol *Smol) CompileFile(filename string, target string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	smol.Compile(string(file), filename, target)
	if smol.HadError {
		os.Exit(65)
	}
}

//Run exectues a given script
func (smol *Smol) Compile(sourceCode string, filename string, target string) {
	l := lexer.NewLexer(filename, sourceCode)
	l.Lex()
	p := ast.NewParser(filename, l.Tokens)
	//We can ignore the second return value here as it is the amount of tokens consumed.
	//We do not need this here
	ast, _ := p.Parse("")
	g := ir.NewGenerator(target)
	g.Generate(ast)
	bg := bytecode.Init(g, filename)
	bg.CreateRom()
	return
}
