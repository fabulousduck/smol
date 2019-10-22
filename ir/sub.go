package ir

/*
SUB instruction


opcode: 8XY5
X: register to deduct a value from
Y: resister holding how much should be deducted from X
*/
type SUB struct {
	TargetRegister, AmountRegister int
}

func (s SUB) GetInstructionName() string {
	return "SUB"
}

func (s SUB) Opcodeable() bool {
	return true
}

func (s SUB) usesVariableSpace() bool {
	return false
}

/*
Sub is a little more complicated than ADD since there is no opcode to increment
a register with a negative value.

For this operation we must set a register (R2) to the value of amount.
Then we can use the 8XY5 opcode to deduct whatever value is in R2 from R1
*/
func (g *Generator) newSubInstruction(targetVariableTableIndex int, amount int) SUB {

	//First step is to check if there is a register free
	amountRegisterName := "amountRegister"
	amountRegister := g.regTable.FindEmptyRegister()
	subInstruction := SUB{targetVariableTableIndex, amountRegister}

	//modify the internal register table so it keeps track of things
	g.regTable.PutRegisterValue(amountRegister, amount, amountRegisterName)

	//create the instruction for the amount register setting
	g.Ir = append(g.Ir, g.newSpecificRegisterSet(targetVariableTableIndex, amount, amountRegisterName))

	//create the actual subtract instruction
	g.Ir = append(g.Ir, subInstruction)

	//create instruction to clear the amount register once we are done
	g.Ir = append(g.Ir, g.newSpecificRegisterSet(amountRegister, 0, ""))

	return subInstruction
}
