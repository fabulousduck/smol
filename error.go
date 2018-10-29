package smol

import (
	"bytes"
	"fmt"
)

func report(line int, where string, message string) {
	//TODO: make this nicer
	fmt.Printf("\n[%s|%d] %s\n\n", where, line, message)
}

func concatVariables(vars []string, sep string) string {
	var currentString bytes.Buffer
	for i := 0; i < len(vars); i++ {
		currentString.WriteString(string(fmt.Sprintf("%s%s", vars[i], sep)))
	}
	return currentString.String()
}

func throwSemanticError(token *token, expected []string, filename string) {
	report(
		token.Line,
		filename,
		fmt.Sprintf("expected one of [%s]. got %s",
			concatVariables(expected, ", "),
			token.Type))
}

func undefinedVariableError(variableName string) {
	//TODO: make this somewhat more informative
	fmt.Printf("Undefined varaible %s\n", variableName)
}

func litAssignError() {
	fmt.Printf("Cannot assign new value to litteral value")
}

func litIncrementError() {
	fmt.Printf("Cannot increment a num literal\n")
}

func undefinedFunctionReferenceError(name string) {
	fmt.Printf("Cannot find function with name: %s\n", name)
}

func incorrectFunctionParamCountError(name string, given int, expected int) {
	fmt.Printf("function \"%s\" requires %d arguments. Got %d\n", name, expected, given)
}

func additionInvalidReceiverError() {
	fmt.Printf("left hand side of mathematical operation must be variable")
}
