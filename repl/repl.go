package smol

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fabulousduck/smol"
)

//Repl : Activates a new interactive REPL reading from STDIN
func Repl(s *smol.Smol) {
	fmt.Printf("Losp repl v0.1\nUse ^C to exit\n\n")
	for true {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("losp> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		s.Run(text, "repl")
		s.HadError = false
	}
}
