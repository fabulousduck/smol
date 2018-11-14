package main

import (
	"os"

	"github.com/fabulousduck/smol"
	"github.com/fabulousduck/smol/repl"
)

func main() {
	s := smol.NewSmol()

	if len(os.Args) > 1 {
		s.RunFile(os.Args[1])
	} else {
		//TODO
		repl.Repl(s)
	}

}
