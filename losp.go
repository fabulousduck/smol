package losp

import (
	"io/ioutil"
	"os"
)

//Losp : Defines the global attributes of the interpreter
type Losp struct {
	Tokens   []*token
	HadError bool
}

//NewLosp : Creates a new Losp instance
func NewLosp() *Losp {
	return new(Losp)
}

//RunFile : Interprets a given file
func (losp *Losp) RunFile(filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	losp.run(string(file), filename)
	if losp.HadError {
		os.Exit(65)
	}
}

func (losp *Losp) run(sourceCode string, filename string) {
	l := new(lexer)
	l.lex(sourceCode, filename)
	p := NewParser(filename)
	p.ast, _ = p.parse(l.tokens)
	i := newInterpreter()
	i.interpret(p.ast)

}
