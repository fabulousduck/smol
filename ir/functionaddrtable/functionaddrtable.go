package functionaddrtable

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
