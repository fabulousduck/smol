package ir

import "github.com/fabulousduck/smol/ast"

/*
FNJMP is a special jump for the chip-8
that is meant for function calls

2NNN

NNN: address of the function on memory
*/
type FNJMP struct {
	addr int
}

func (f FNJMP) GetInstructionName() string {
	return "FNJMP"
}

func (f FNJMP) Opcodeable() bool {
	return true
}

func (f FNJMP) usesVariableSpace() bool {
	return false
}

func (g *Generator) newFNJMPInstruction(jmpAddr int) FNJMP {
	return FNJMP{jmpAddr}
}

func (g *Generator) createFunctionCallInstructions(instruction *ast.FunctionCall) {

	//lookup the function on the function table
	fnTableEntry := g.functionAddrTable.Find(instruction.Name)

	g.Ir = append(g.Ir, g.newFNJMPInstruction(fnTableEntry.Addr))
}
