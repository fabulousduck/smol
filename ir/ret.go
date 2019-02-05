package ir

/*
RET is used when exiting a function

has no contents. only used to signal to the bytecode generator
that it needs to put one in

opcode: 00EE
*/
type RET struct {
}

func (r RET) GetInstructionName() string {
	return "RET"
}

func (r RET) Opcodeable() bool {
	return false
}

func (r RET) usesVariableSpace() bool {
	return false
}

func (g *Generator) newRetInstruction() RET {
	return RET{}
}
