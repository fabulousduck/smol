package gameboy

import ir "github.com/fabulousduck/smol/ir/gameboy"

/*
Generator contains all the info needed to generator the
actual opcodes for a gameboy game
*/
type Generator struct {
	irGenerator *ir.Generator
	filename    string
}

/*
Init returns a pointer to a new Generator struct
*/
func Init(irGenerator *ir.Generator, filename string) *Generator {
	g := new(Generator)
	g.irGenerator = irGenerator
	g.filename = filename
	return g
}

/*
Generate turns an ir into an actual binary
*/
func (g *Generator) Generate() {}
