package losp

import "fmt"

func report(line int, where string, message string) {
	//TODO: make this nicer
	fmt.Printf("\n[%s|%d] %s\n\n", where, line, message)
}
