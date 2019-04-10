package ir

import (
	"github.com/fabulousduck/smol/ast"
	"github.com/fabulousduck/smol/errors"
	"github.com/fabulousduck/smol/ir/memtable"
)

/*
Mov R1 into R2

Chip-8 knows MOV in the form of 6XNN
where X is the register and NN is the memory address
*/
type MOV struct {
	R1, R2 int
	ANNN   bool
}

func (m MOV) GetInstructionName() string {
	return "MOV"
}

func (m MOV) Opcodeable() bool {
	return true
}

func (m MOV) usesVariableSpace() bool {
	return true
}

/*

	newMovInstructionFromLoose takes a loose set of values and turns them into
	a MOV instruction

	R1 must be a register
	R2 can either be a register or a memory address
	if R2 is a memory address the field "R2IsAddr" will be set to true
*/
func (g *Generator) newMovInstructionFromLoose(R1 int, R2 int, ANNN bool) MOV {
	instr := MOV{R1, R2, false}
	if memtable.IsValidMemRegion(R1) {
		errors.RegisterAdressModeFailure(R1)
	}

	instr.ANNN = ANNN
	return instr
}

/*
	same as newMovInstructionFromLoose.
	It does the extraction of variable values for you
*/
func (g *Generator) newMovInstruction(v *ast.Variable) MOV {
	instr := MOV{}
	return instr
}
