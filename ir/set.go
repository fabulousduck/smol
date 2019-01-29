package ir

/*
	note: this thing is internal as fuck. there is no opcode for this. this is something we just do. not the emulator or a system for that matter
*/

import (
	"github.com/fabulousduck/smol/ast"
	"github.com/fabulousduck/smol/errors"
)

type SETMEM struct {
	Val, Addr int
}

func (s SETMEM) GetInstructionName() string {
	return "SETMEM"
}

func (s SETMEM) Opcodeable() bool {
	return false
}

func (s SETMEM) usesVariableSpace() bool {
	return true
}

/*
SET Val ADDR

set a value in a given register address
*/
type SETREG struct {
	Val, Index int
}

func (s SETREG) GetInstructionName() string {
	return "SETREG"
}

func (s SETREG) Opcodeable() bool {
	return true
}

func (s SETREG) usesVariableSpace() bool {
	return true
}

func (g *Generator) newSetMemoryLocationFromLoose(name string, value int) SETMEM {
	instr := SETMEM{}
	region := g.memTable.Put(name, value)
	instr.Addr = region.Addr
	instr.Val = value
	return instr
}

/*
This is used for MEM operations.
Sets the declared variable in a free register
Will error out if it cant find a free register
*/
func (g *Generator) newSetRegisterInstruction(v *ast.Variable, varValue int, resolve bool) SETREG {
	instr := SETREG{}

	if resolve {
		resolutionName := v.Value.(*ast.StatVar).Value
		if val, ok := g.memTable[resolutionName]; ok {
			varValue = val.Value
		} else {
			errors.UndefinedVariableError(resolutionName)
		}
	}

	emptyRegisterAddress := g.regTable.FindEmptyRegister()
	g.regTable.PutRegisterValue(emptyRegisterAddress, varValue, v.Name)
	instr.Index = emptyRegisterAddress
	instr.Val = varValue
	return instr
}

func (g *Generator) newSetRegisterInstructionFromLoose(registerName string, varValue int) SETREG {
	instr := SETREG{}

	emptyRegisterAddress := g.regTable.FindEmptyRegister()
	g.regTable.PutRegisterValue(emptyRegisterAddress, varValue, registerName)
	instr.Index = emptyRegisterAddress
	instr.Val = varValue
	return instr
}
