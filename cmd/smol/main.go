package main

import (
	"flag"

	"github.com/fabulousduck/smol"
	"github.com/fabulousduck/smol/repl"
)

func main() {
	s := smol.NewSmol()
	flagPtr := flag.Bool("c", false, "compile to a chip-8 rom")
	filenamePtr := flag.String("file", "", "input file for the interpreter")

	flag.Parse()

	if *filenamePtr != "" {
		s.RunFile(*filenamePtr, *flagPtr)
	} else {
		//TODO add on the fly rom compilation to repls
		repl.Repl(s)
	}

}
