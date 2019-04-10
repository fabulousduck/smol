package ir

/*
RGD 0 X

This instruction is used when we need to dump 0 through X
registers into the rom.

This is used when we need to do a temporary scope change for
functions
*/
type RGD struct {
	startLocation, count int
}

func (r RGD) GetInstructionName() string {
	return "RGD"
}

func (r RGD) Opcodeable() bool {
	return false
}

func (r RGD) usesVariableSpace() bool {
	return false
}

/*
NewRGDInstruction creates a new RGD instruction

uses the I register
*/
func (g *Generator) NewRGDInstruction(endRegister int) RGD {
	instr := RGD{}
	instr.count = endRegister

	//find a region for the dump to go into
	emptyRegionStart := g.memTable.FindNextEmptyAddr()
	instr.startLocation = emptyRegionStart

	//set I register to the start of the dump region
	g.regTable.PutRegisterValue(g.IRegisterIndex, emptyRegionStart, "I_REGISTER")
	g.Ir = append(g.Ir, g.newSpecificRegisterSet(g.IRegisterIndex, emptyRegionStart, "I_REGISTER"))

	//put all register values into a slice
	for i := 0; i < endRegister; i++ {
		g.memTable.Put(g.regTable[i].Name, g.regTable[i].Value, 1)
	}

	return instr
}
