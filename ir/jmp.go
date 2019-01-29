package ir

/*
JMP FROM TO
*/
type Jump struct {
	To int
}

func (j Jump) GetInstructionName() string {
	return "Jump"
}

func (j Jump) Opcodeable() bool {
	return true
}

func (j Jump) usesVariableSpace() bool {
	return false
}

func (g *Generator) newJumpInstructionFromLoose(to int) Jump {
	return Jump{to}
}
