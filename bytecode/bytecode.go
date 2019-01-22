package bytecode

import (
	"github.com/fabulousduck/smol/file"
	"github.com/fabulousduck/smol/ir"
)

/*
Generator holds the file pointer to the binary ROM
and other info relevant to generating the opcodes
*/
type Generator struct {
	filename string
	ir       *ir.Generator
}

/*
Init creates and fils a new bytecode generator struct
*/
func Init(ir *ir.Generator, filename string) *Generator {
	g := new(Generator)
	g.ir = ir
	g.filename = filename

	return g
}

/*
CreateRom generates a rom from an existing IR
*/
func (g *Generator) CreateRom() {
	romFile := file.Create(g.filename)

	for i := 0; i < len(g.ir.Ir); i++ {
		instructionType := g.ir.Ir[i].GetInstructionName()
		switch instructionType {
		case "SET":
			setInstruction := g.ir.Ir[i].(ir.SET)
			file.WriteBytes(romFile, []byte{byte(uint8(setInstruction.Val))}, true, int64(setInstruction.Addr))
			break
		case "MOV":
			movInstruction := g.ir.Ir[i].(ir.MOV)

			/*
				for 0xANNN
			*/
			if movInstruction.ANNN {
				movInstruction.R2 += 0x200 //generate and address that is relative to the machine, not the file
				baseByte := 0xA
				secondaryByte := 0x00
				baseByte = baseByte<<4 | shiftLeft(movInstruction.R2)

				secondaryByte = movInstruction.R2 & 0x0FF

				file.WriteBytes(romFile, []byte{clampUint8(baseByte), clampUint8(secondaryByte)}, false, 0)
				break
			}

			baseByte := 0x6

			baseByte = baseByte<<4 | movInstruction.R1
			file.WriteBytes(romFile, []byte{clampUint8(baseByte), clampUint8(movInstruction.R2)}, false, 0)
		case "PLOT":
			plotInstruction := g.ir.Ir[i].(*ir.PLOT)

			baseByte := 0xD
			baseByte = baseByte<<4 | plotInstruction.X

			secondaryByte := plotInstruction.Y<<4 | plotInstruction.H

			file.WriteBytes(romFile, []byte{clampUint8(baseByte), clampUint8(secondaryByte)}, false, 0)

		case "JMP":
			jmpInstruction := g.ir.Ir[i].(ir.JMP)

			baseByte := 0x1
			baseByte = baseByte<<4 | shiftLeft(jmpInstruction.To)

			secondaryByte := jmpInstruction.To & 0x0FF

			file.WriteBytes(romFile, []byte{clampUint8(baseByte), clampUint8(secondaryByte)}, false, 0)
		}
	}
	return
}

// ????

func shiftLeft(value int) int {
	if value <= 16 {
		return value
	} else if value < 256 {
		return (value & 0xF0 >> 4)
	} else if value < 4096 {
		return (value & 0xF00 >> 8)
	} else {
		return (value & 0xF000 >> 16)
	}
}

func clampUint8(variable int) byte {
	return byte(uint8(variable))
}
