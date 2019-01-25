package ir

/*
JMP FROM TO
*/
type JMP struct {
	To int
}

func (j JMP) GetInstructionName() string {
	return "JMP"
}

func (j JMP) Opcodeable() bool {
	return true
}

func (j JMP) usesVariableSpace() bool {
	return false
}

func (g *Generator) newJumpInstructionFromLoose(to int) JMP {
	return JMP{to}
}
