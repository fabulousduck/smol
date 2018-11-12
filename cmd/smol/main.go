package main

import (
	"os"

	"github.com/fabulousduck/smol"
)

func main() {
	l := smol.NewSmol()

	if len(os.Args) > 1 {
		l.RunFile(os.Args[1])
	} else {
		//TODO
		// repl.Repl()
	}

}
