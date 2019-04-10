package functionaddrtable

import (
	"os"

	"github.com/fabulousduck/smol/errors"
)

//FunctionAddrTable is a simple type so we can do function mounting on it
type FunctionAddrTable []FunctionAddr

/*
FunctionAddr stores basic information about a function
and where it is stored in memory
*/
type FunctionAddr struct {
	Addr int
	Name string
}

/*
NewFunctionAddr returns a new filled FunctionAddr struct
*/
func NewFunctionAddr(addr int, name string) FunctionAddr {
	return FunctionAddr{addr, name}
}

/*
Find checks if a given function with name name exists in the function table
*/
func (table FunctionAddrTable) Find(name string) FunctionAddr {
	for i := 0; i < len(table); i++ {
		if table[i].Name == name {
			return table[i]
		}
	}
	errors.UnknownFunctionName(name)
	os.Exit(65)
	return FunctionAddr{}
}
