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
	0x202 | 7C01 #increment the C register with one
	0x204 | 4C0F #call to BNE
	0x206 | 1208 #jump to outside of the loop
	0x208 | 1200 # jump back to the beginning of the loop
	0x20A | code after loop

	should the bne care what is in 0xC ?
	i dont think it does
	it just needs to check if whatever is in 0xC equals the RHS

	INC 0xC, 1
	BNE 0xC, 20
	JMP 0x208
	JMP 0x200
*/

func (g *Generator) newBNEInstruction(instruction *ast.Anb) BNE {
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
