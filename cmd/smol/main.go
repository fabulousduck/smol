package main

import (
	"flag"

	"github.com/fabulousduck/smol"
)

func main() {
	s := smol.NewSmol()
	filenamePtr := flag.String("file", "", "input file for the interpreter")

	flag.Parse()
	s.RunFile(*filenamePtr)

}
