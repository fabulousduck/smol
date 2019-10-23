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

//UnknownFunctionName is an error when a lookup on a function is done but none could be found
func UnknownFunctionName(name string) {
	fmt.Printf("tried to call unknown function: %s\n", name)
}

//IlligalRegisterAccess is thrown by the register table when it detects the compilers accesses a non existant register
func IlligalRegisterAccess(register int) {
	fmt.Printf("illigal access of register: %d", register)
}

//UnAssignedMemoryLookupError is an error for the IR to throw when it wants to find a variable by addr in the memtable that does not exist
func UnAssignedMemoryLookupError() {
	fmt.Printf("tried to look up register that is empty while expecting it to be full")
}

//UnknownTypeError is an error for the AST generator for when it encounters a token that it does not have a name for
func UnknownTypeError() {
	fmt.Printf("Unknown token type found.")
}

//LitteralFree error can be thrown when the programmer wants to free a number litteral
func LitteralFree() {
	fmt.Printf("Cannot release a number litteral\n")
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

//ROMModError can be thrown when a variable modification is called on a variable that is not loaded into a register.
//the user is most likely attempting to change rom here
func ROMModError() {
	fmt.Printf("Trying to modify variable that is not loaded into memory")
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

//OutOfRegistersError is and error that indicates someone tried to assign more values than is allowed by the bytecode generator
func OutOfRegistersError() {
	fmt.Printf("Tried to store more variables than available registers (15) ")
}

//OutOfMemoryError can be thrown when the compiler has no more space to place a variable
func OutOfMemoryError() {
	fmt.Printf("Out of memory error")
}

/*
MemAddrAdressModeFailure is an internal error where a MOV statement tried to use
a memory address as arg one

TODO: rename this to something more appropriate
TODO: maybe even make a separate errors package for internal errors
*/
func RegisterAdressModeFailure(attemptedRegisterIndex int) {
	fmt.Printf("Invalid MOV to register [%d]. 0xF boundary exceeded", attemptedRegisterIndex)
}

/*
s
*/
