package ir

/*
JMP FROM TO

has an ID field because sometimes we need to reference to it
for stuff like function hops where we need manipulate the call later
to set the jump to address
*/
type Jump struct {
	To int
	ID string
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
	return Jump{to, "0"}
}
