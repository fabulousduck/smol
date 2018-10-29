package main

import (
	"os"

	"github.com/fabulousduck/smol"
)

func main() {
	l := smol.NewLosp()

	if len(os.Args) > 1 {
		l.RunFile(os.Args[1])
	} else {
		l.Repl()
	}

}
