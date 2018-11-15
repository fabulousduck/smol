package errors

import (
	"bytes"
	"fmt"
)

//Report can be used for appending info to the error message such as line number and file
func Report(line int, where string, message string) {
	//TODO: make this nicer
	fmt.Printf("\n[%s|%d] %s\n\n", where, line, message)
}

//ConcatVariables is a simple formatter to create a variable error string
func ConcatVariables(vars []string, sep string) string {
	var currentString bytes.Buffer
	for i := 0; i < len(vars); i++ {
		currentString.WriteString(string(fmt.Sprintf("%s%s", vars[i], sep)))
	}
	return currentString.String()
}

//UndefinedVariableError can be thrown at interpret time when a variable is not found on the local scope or higher level scopes
func UndefinedVariableError(variableName string) {
	//TODO: make this somewhat more informative
	fmt.Printf("Undefined varaible %s\n", variableName)
}

//LitAssignError can be used when the script tries to assign a new value to a litteral value
func LitAssignError() {
	fmt.Printf("Cannot assign new value to litteral value")
}

//LitIncrementError can be thrown when the script wants to call INC on a litteral. We do not support this as litterals are not expressions and we dont support returns yet
func LitIncrementError() {
	fmt.Printf("Cannot increment a num literal\n")
}

//UndefinedFunctionReferenceError can be thrown when the script tries to reference an error that is not defined
func UndefinedFunctionReferenceError(name string) {
	fmt.Printf("Cannot find function with name: %s\n", name)
}

//IncorrectFunctionParamCountError can be throw when more or less arguments are provided to a function than it asks for. We dont support argument defaulting so this is usefull
func IncorrectFunctionParamCountError(name string, given int, expected int) {
	fmt.Printf("function \"%s\" requires %d arguments. Got %d\n", name, expected, given)
}

//MathInvalidReceiverError can be thrown when the script wants to do a mathematical statement but does not have a receiver for the outcome as LHS
func MathInvalidReceiverError() {
	fmt.Printf("left hand side of mathematical operation must be variable")
}

//UnknownSwitchNode is thrown when something else than EOS or CAS is found as a top level definition in a switch
func UnknownSwitchNode() {
	fmt.Printf("unknown definition found in switch")
}

//EOFError allows us to throw an error when either the lexer or the AST generator runs out of tokens / characters to parse
//while it still expects there to be a token or character.
func EOFError() {
	fmt.Printf("EOF found in program execution.")
}
