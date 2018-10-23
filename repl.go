package losp

import (
	"bufio"
	"fmt"
	"os"
)

//Repl : Activates a new interactive REPL reading from STDIN
func (losp *Losp) Repl() {
	fmt.Printf("Losp repl v0.1\nUse ^C to exit\n\n")
	for true {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("losp> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		losp.run(text, "repl")
		losp.HadError = false
	}
}
