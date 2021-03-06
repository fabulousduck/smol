package ir

type RegCpy struct {
	From, To int
}

func (j RegCpy) GetInstructionName() string {
	return "RegCpy"
}

func (j RegCpy) Opcodeable() bool {
	return true
}

func (j RegCpy) usesVariableSpace() bool {
	return false
}

/*
newRegCpy: from R1 into R2
*/
func (g *Generator) newRegCpy(R1 int, R2 int) RegCpy {
	return RegCpy{R1, R2}
}
