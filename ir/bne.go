package ir

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
