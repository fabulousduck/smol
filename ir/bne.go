package ir

import (
	"strconv"

	"github.com/fabulousduck/smol/ast"
)

/*
BNE is a simple structure that will skips the next instruction
if lhs does not equal rhs

opcode: 4XNN
4: identifier
X: lhs
NN: rhs
*/
type BNE struct {
	lhs, rhs int
}

func (b BNE) GetInstructionName() string {
	return "BNE"
}

func (b BNE) Opcodeable() bool {
	return true
}

func (b BNE) usesVariableSpace() bool {
	return false
}

/*
	0x200 | DDE1 #some arbitrary shit like a plot opcode
	0x202 | 4C0F #call to BNE
	0x204 | 1208 #jump to outside of the loop
	0x206 | 1200 # jump back to the beginning of the loop
	0x200 | code after loop
*/

func (g *Generator) newBNEInstruction(instruction ast.Anb) BNE {
	instr := BNE{}
	/*
		check if the lhs is a variable
		if so, resolve it
	*/
	if ast.NodeIsVariable(instruction.LHS) {
		variableName := instruction.LHS.(*ast.StatVar).Value
		variableValue := g.memTable.LookupVariable(variableName, false)
		instr.lhs = variableValue.Value
	} else {
		variableValue := instruction.LHS.(*ast.NumLit).Value
		intValue, _ := strconv.Atoi(variableValue)
		instr.lhs = intValue
	}

	//set the value to compare in the register
	//we need to do this every time since the memory adress will change
	g.Ir = append(g.Ir, g.newMovInstructionFromLoose(g.BNEXRegister, instr.lhs, false))

	/*
		check if the rhs is a variable
		if so, resolve it
	*/
	if ast.NodeIsVariable(instruction.LHS) {
		variableName := instruction.LHS.(*ast.StatVar).Value
		variableValue := g.memTable.LookupVariable(variableName, false)
		instr.lhs = variableValue.Value
	} else {
		variableValue := instruction.LHS.(*ast.NumLit).Value
		intValue, _ := strconv.Atoi(variableValue)
		instr.lhs = intValue
	}

	return instr
}
