package main

import (
	"flag"

	"github.com/fabulousduck/smol"
)

func main() {
	s := smol.NewSmol()
	filenamePtr := flag.String("f", "", "input file for the interpreter")
	targetPtr := flag.String("t", "", "compile target CPU")

	flag.Parse()
	s.CompileFile(*filenamePtr, *targetPtr)

}
