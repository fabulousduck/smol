package gameboy

/*
Generator contains all the info needed to generator the
actual opcodes for a gameboy game
*/
type Generator struct {
}

/*
Init returns a pointer to a new Generator struct
*/
func Init() *Generator {
	g := new(Generator)
	return g
}
