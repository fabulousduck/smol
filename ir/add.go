package ir

/*
ADD instruction

opcode: 7XNN
X: register to add onto
NN: value to add onto value in register
*/
type ADD struct {
	Register, Value int
}

func (a ADD) GetInstructionName() string {
	return "ADD"
}

func (a ADD) Opcodeable() bool {
	return true
}

func (a ADD) usesVariableSpace() bool {
	return false
}

func (g *Generator) newAddInstruction(R1 int, value int) ADD {
	return ADD{R1, value}
}
