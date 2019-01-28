package ir

/*
	note: this thing is internal as fuck. there is no opcode for this. this is something we just do. not the emulator or a system for that matter
*/

import (
	"github.com/fabulousduck/smol/ast"
	"github.com/fabulousduck/smol/errors"
)

/*
SET Val ADDR

set a value at a given memory address
*/
type SET struct {
	Val, Addr int
}

func (s SET) GetInstructionName() string {
	return "SET"
}

func (s SET) Opcodeable() bool {
	return false
}

func (s SET) usesVariableSpace() bool {
	return true
}

/*
This is used for MEM operations.
This means that it is used for setting variables into memory before any opcode is executed
Once a variable is needed for an operation. for instance INC or ANB or any other. It will
be retreived by MOV and put into a free register
*/
func (g *Generator) newSetInstruction(v *ast.Variable, varValue int, resolve bool) SET {
	isntr := SET{}

	if resolve {
		resolutionName := v.Value.(*ast.StatVar).Value
		if val, ok := g.memTable[resolutionName]; ok {
			varValue = val.Value
		} else {
			errors.UndefinedVariableError(resolutionName)
		}
	}

	region := g.memTable.Put(v.Name, varValue)
	isntr.Addr = region.Addr
	isntr.Val = varValue
	return isntr
}

func (g *Generator) newSetInstructionFromLoose(name string, varValue int) SET {
	instr := SET{}
	region := g.memTable.Put(name, varValue)
	instr.Addr = region.Addr
	instr.Val = varValue
	return instr
}
