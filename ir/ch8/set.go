package ir

/*
	note: this thing is internal as fuck. there is no opcode for this. this is something we just do. not the emulator or a system for that matter
*/

import (
	"strconv"

	"github.com/fabulousduck/smol/ast"
	"github.com/fabulousduck/smol/errors"
	"github.com/fabulousduck/smol/ir/ch8/registertable"
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
	region := g.memTable.Put(name, value, 1)
	instr.Addr = region.Addr
	instr.Val = value
	return instr
}

func (g *Generator) newSpecificRegisterSet(registerIndex int, value int, name string) SETREG {
	instr := SETREG{}
	g.regTable[registerIndex] = registertable.Register{value, name}
	instr.Index = registerIndex
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

func (g *Generator) createVariableOperationInstructions(variable *ast.Variable) {
	//check if its a reference
	if ast.NodeIsVariable(variable.Value) {
		//if it is a reference, we get the original value,
		//and copy it over into a new register with the name of the new variable
		variableValue := variable.Value.(*ast.StatVar)
		emptyRegister := g.regTable.FindEmptyRegister()
		originalRegister := g.regTable.Find(variableValue.Value)
		g.regTable[emptyRegister] = registertable.Register{g.regTable[originalRegister].Value, variable.Name}
		g.Ir = append(g.Ir, g.newRegCpy(originalRegister, emptyRegister))
	} else {
		variableValue, _ := strconv.Atoi(variable.Value.(*ast.NumLit).Value)
		g.Ir = append(g.Ir, g.newSetRegisterInstructionFromLoose(variable.Name, variableValue))
	}
}

func (g *Generator) createSetStatement(instruction *ast.SetStatement) {
	castVariable := instruction.MHS.(*ast.StatVar)

	//find the register in which the variable is currently stored
	variableRegister := g.regTable.Find(castVariable.Value)

	//if the rhs of the set statement is a variable too, we need to get its value first
	//and then embed a register copy instruction
	if ast.NodeIsVariable(instruction.RHS) {
		referenceVariableRegister := g.regTable.Find(instruction.RHS.(*ast.StatVar).Value)
		g.Ir = append(g.Ir, g.newRegCpy(referenceVariableRegister, variableRegister))
	} else {
		//otherwise, we need to set the value of the register to the right hand side value
		castVal, _ := strconv.Atoi(instruction.RHS.(*ast.NumLit).Value)
		g.Ir = append(g.Ir, g.newSetRegisterInstructionFromLoose(castVariable.Value, castVal))

	}
}
