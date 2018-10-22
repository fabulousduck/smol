package losp

import "fmt"

func (losp *Losp) err(line int, message string) {
	losp.report(line, "", message)
}

func (losp *Losp) report(line int, where string, message string) {
	//TODO: make this nicer
	fmt.Printf(" %s | at \n\t%d", message, line)
	losp.HadError = true
}
