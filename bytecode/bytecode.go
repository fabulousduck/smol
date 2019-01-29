package bytecode

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"

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
		spew.Dump(instructionType)
		fmt.Printf("\n")
		switch instructionType {
		case "SETREG":
			setRegInstruction := g.ir.Ir[i].(ir.SETREG)
			g.embedSetRegister(setRegInstruction, romFile)
		case "SETMEM":
			setMemInstruction := g.ir.Ir[i].(ir.SETMEM)
			//write a single byte into memory
			file.WriteBytes(romFile, []byte{byte(uint8(setMemInstruction.Val))}, true, int64(setMemInstruction.Addr))
			break
		case "RegCpy":
			regCpyInstruction := g.ir.Ir[i].(ir.RegCpy)
			g.embedRegCpy(regCpyInstruction, romFile)
		case "MOV":
			movInstruction := g.ir.Ir[i].(ir.MOV)
			if movInstruction.ANNN {
				g.embedANNN(movInstruction, romFile)
				break
			}
			g.embedMOV(movInstruction, romFile)
		case "BNE":
			bneInstruction := g.ir.Ir[i].(ir.BNE)
			g.embedBNE(bneInstruction, romFile)
		case "BNERR":
			bnerrInstruction := g.ir.Ir[i].(ir.BNERR)
			g.embedBNERR(bnerrInstruction, romFile)
		case "PLOT":
			plotInstruction := g.ir.Ir[i].(ir.PLOT)
			g.embedPLOT(plotInstruction, romFile)
		case "JMP":
			jmpInstruction := g.ir.Ir[i].(ir.Jump)
			g.embedJMP(jmpInstruction, romFile)

		}
	}
	return
}

func (g *Generator) embedBNE(instruction ir.BNE, romFile *os.File) {
	baseByte := 0x3
	baseByte = baseByte<<4 | instruction.Lhs
	secondaryByte := instruction.Rhs
	opcode := []byte{clampUint8(baseByte), clampUint8(secondaryByte)}
	file.WriteBytes(romFile, opcode, false, 0)
}

func (g *Generator) embedBNERR(instruction ir.BNERR, romFile *os.File) {
	baseByte := 0x5
	baseByte = baseByte<<4 | instruction.Lhs
	secondaryByte := instruction.Rhs<<4 | 0
	opcode := []byte{clampUint8(baseByte), clampUint8(secondaryByte)}
	file.WriteBytes(romFile, opcode, false, 0)
}

/*
	opcode: 8XY0
	X: register where the original value is stored
	Y: register where the value is to be copied to
*/
func (g *Generator) embedRegCpy(instruction ir.RegCpy, romFile *os.File) {
	baseByte := 0x8
	baseByte = baseByte<<4 | instruction.From
	secondaryByte := instruction.To<<4 | 0

	opcode := []byte{clampUint8(baseByte), clampUint8(secondaryByte)}
	file.WriteBytes(romFile, opcode, false, 0)
}

/*
	opcode: 6XNN
	X: index of the register that the value will be placed in
	NN: the value to be placed in the register
*/
func (g *Generator) embedSetRegister(instruction ir.SETREG, romFile *os.File) {
	baseByte := 0x6
	baseByte = baseByte<<4 | instruction.Index
	secondaryByte := instruction.Val
	opcode := []byte{clampUint8(baseByte), clampUint8(secondaryByte)}
	file.WriteBytes(romFile, opcode, false, 0)
}

/*
	opcode: 1NNN
	1: identifier
	NNN: address to jump to
*/
func (g *Generator) embedJMP(instruction ir.Jump, romFile *os.File) {
	baseByte := 0x1
	baseByte = baseByte<<4 | shiftRight(instruction.To)

	secondaryByte := instruction.To & 0x0FF
	file.WriteBytes(romFile, []byte{clampUint8(baseByte), clampUint8(secondaryByte)}, false, 0)
}

/*
	opcode: ANNN
	A: identifier
	NNNN: address to move into I
*/
func (g *Generator) embedANNN(instruction ir.MOV, romFile *os.File) {
	instruction.R2 += 0x200 //generate and address that is relative to the machine, not the file
	baseByte := 0xA
	secondaryByte := 0x00
	baseByte = baseByte<<4 | shiftRight(instruction.R2)

	secondaryByte = instruction.R2 & 0x0FF

	file.WriteBytes(romFile, []byte{clampUint8(baseByte), clampUint8(secondaryByte)}, false, 0)
}

/*
opcode: DXYN
D: identifier
X: register index containing the X coordinate
Y: register index containing the Y coordinate
N: number of columns to draw
*/

func (g *Generator) embedPLOT(instruction ir.PLOT, romFile *os.File) {

	baseByte := 0xD
	baseByte = baseByte<<4 | instruction.X

	secondaryByte := instruction.Y<<4 | instruction.H

	file.WriteBytes(romFile, []byte{clampUint8(baseByte), clampUint8(secondaryByte)}, false, 0)
}

/*
	opcode: 6XNN
	6: identifier
	X: register index
	NN: value to be moved into register
*/
func (g *Generator) embedMOV(instruction ir.MOV, romFile *os.File) {

	baseByte := 0x6

	baseByte = baseByte<<4 | instruction.R1
	file.WriteBytes(romFile, []byte{clampUint8(baseByte), clampUint8(instruction.R2)}, false, 0)
}

/*
	shifts the left most value in an int to the right
	this is used when we need to split up a number to append it
	to the right in a byte
*/
func shiftRight(value int) int {
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
