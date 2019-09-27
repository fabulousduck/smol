package ir

import (
	"strconv"

	"github.com/fabulousduck/smol/ast"
	"github.com/fabulousduck/smol/ir/ch8/registertable"
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
	Lhs, Rhs int
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
	BNERR is the same as BNE but the RHS is also a register

	opcode: 5XY0
	5: indentifier
	X: lhs register
	Y: rhs register
*/
type BNERR struct {
	Lhs, Rhs int
}

func (b BNERR) GetInstructionName() string {
	return "BNERR"
}

func (b BNERR) Opcodeable() bool {
	return true
}

func (b BNERR) usesVariableSpace() bool {
	return false
}

func (g *Generator) newBNEInstructionFromLoose(R1 int, rhs int) BNE {
	instr := BNE{R1, rhs}

	return instr
}

func (g *Generator) newBNERRInstructionFromLoose(R1 int, R2 int) BNERR {
	return BNERR{R1, R2}
}

func (g *Generator) createAnbInstructions(instruction *ast.Anb) {
	// we need to multiply by 2 because each instruction is 2 bytes long
	//doing this here is a little trick to get the value to keep getting updated from its origin register
	anbInstructionStart := 0x200 + (len(g.Ir) * 2)
	g.Generate(instruction.Body)

	if ast.NodeIsVariable(instruction.LHS) {
		variable := instruction.LHS.(*ast.StatVar)
		lhsVariableRegister := g.regTable.Find(variable.Value)
		g.regTable[g.BNEXRegister] = registertable.Register{g.regTable[lhsVariableRegister].Value, "BNEX"}
		g.Ir = append(g.Ir, g.newRegCpy(lhsVariableRegister, g.BNEXRegister))
	} else {
		lhsValue, _ := strconv.Atoi(instruction.LHS.(*ast.NumLit).Value)
		bnexreg := registertable.Register{lhsValue, "BNEX"}
		g.regTable[g.BNEXRegister] = bnexreg
		g.Ir = append(g.Ir, g.newSpecificRegisterSet(g.BNEXRegister, lhsValue, "BNEX"))
	}

	if ast.NodeIsVariable(instruction.RHS) {
		variable := instruction.RHS.(*ast.StatVar)
		rhsVariableRegister := g.regTable.Find(variable.Value)
		g.Ir = append(g.Ir, g.newBNERRInstructionFromLoose(g.BNEXRegister, rhsVariableRegister))
	} else {
		rhsValue, _ := strconv.Atoi(instruction.RHS.(*ast.NumLit).Value)
		g.Ir = append(g.Ir, g.newBNEInstructionFromLoose(g.BNEXRegister, rhsValue))
	}
	g.Ir = append(g.Ir, g.newJumpInstructionFromLoose(anbInstructionStart))
}
