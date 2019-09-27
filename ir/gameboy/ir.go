package ir

/*
Generator is a structure that contains information
required to generate an IR for the gameboy system
*/
type Generator struct {
}

/*
Init creates a new generator instance
that can be used to transform an AST into an
IR in gameboy format
*/
func Init() *Generator {
	g := new(Generator)
	return g
}
