package main

import (
	"flag"

	"github.com/fabulousduck/smol"
	"github.com/fabulousduck/smol/repl"
)

func main() {
	s := smol.NewSmol()
	ch8FlagPointer := flag.Bool("ch8", false, "compile to a chip-8 rom")
	gameboyFlagPointer := flag.Bool("gb", false, "compile to a gameboy binary")
	filenamePtr := flag.String("file", "", "input file for the interpreter")

	flag.Parse()

	if *filenamePtr != "" {
		s.RunFile(*filenamePtr, *ch8FlagPointer, *gameboyFlagPointer)
	} else {
		//TODO add on the fly rom compilation to repls
		repl.Repl(s)
	}

}
