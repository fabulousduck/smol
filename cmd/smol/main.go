package main

import (
	"flag"
	"fmt"

	"github.com/fabulousduck/smol"
	"github.com/fabulousduck/smol/repl"
)

func main() {
	s := smol.NewSmol()
	flagPtr := flag.Bool("c", false, "compile to a chip-8 rom")
	targetPtr := flag.String("t", "", "compile target CPU ")
	filenamePtr := flag.String("file", "", "input file for the interpreter")

	flag.Parse()

	if *filenamePtr != "" {
		fmt.Printf("compiling for target arch : %s\n", *targetPtr)
		s.RunFile(*filenamePtr, *flagPtr, *targetPtr)
	} else {
		//TODO add on the fly rom compilation to repls
		repl.Repl(s)
	}

}
