package registertable

import "github.com/fabulousduck/smol/errors"

/*
RegisterTable is a simple collection of registers so they can be indexed
*/
type RegisterTable map[int]Register

/*
Register simulates a basic CPU register

Value: the actual value in the register
Name: the name of the current variable in it
*/
type Register struct {
	Value int
	Name  string
}

/*
Find finds a variable on the memory table
Returns the index at which it is found
Returns -1 if the value cannot be found
*/
func (table RegisterTable) Find(name string) int {
	for i := 0; i < len(table); i++ {
		region := table[i]
		if region.Name == name {
			return i
		}
	}
	return -1
}

/*
Init fills a new table with empty registers
*/
func (table RegisterTable) Init() {
	for i := 0; i < 0x10; i++ {
		table[i] = Register{
			Value: 0,
			Name:  "",
		}
	}
}

/*
PutRegisterValue set the value of register to value
*/
func (table RegisterTable) PutRegisterValue(register int, value int) {
	if !isValidRegisterIndex(register) {
		errors.IlligalRegisterAccess(register)
	}

	table[register] = Register{value, table[register].Name}
}

func isValidRegisterIndex(registerIndex int) bool {
	return registerIndex < 0xF
}

/*
	somewhere we need to keep track of the variable in the ANB statement
	we need to do this because we need to increment it every loop successively

	the loops can be infinite, so there is a possibility that the variable in the BNE X register never gets adressed


	we need something like a "is currently in register" field on a variable, that if it is, we up the register value
*/
