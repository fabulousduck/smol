package main

import (
	"os"

	"github.com/fabulousduck/losp"
)

func main() {
	l := losp.NewLosp()

	if len(os.Args) > 1 {
		l.RunFile(os.Args[1])
	} else {
		l.Repl()
	}

}
