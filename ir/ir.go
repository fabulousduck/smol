package ir

import (
	"strconv"

	"github.com/davecgh/go-spew/spew"

	"github.com/fabulousduck/smol/ast"
	"github.com/fabulousduck/smol/errors"
	"github.com/fabulousduck/smol/ir/memtable"
	"github.com/fabulousduck/smol/ir/registertable"
)

type instruction interface {
	GetInstructionName() string
	Opcodeable() bool
	usesVariableSpace() bool
}

/*
Mov R1 into R2

Chip-8 knows MOV in the form of 6XNN
where X is the register and NN is the memory address
*/
type MOV struct {
	R1, R2 int
	ANNN   bool
}

func (m MOV) GetInstructionName() string {
	return "MOV"
}

func (m MOV) Opcodeable() bool {
	return true
}

func (m MOV) usesVariableSpace() bool {
	return true
}

/*
JMP FROM TO
*/
type JMP struct {
	From, To int
}

func (j JMP) GetInstructionName() string {
	return "JMP"
}

func (j JMP) Opcodeable() bool {
	return true
}

func (j JMP) usesVariableSpace() bool {
	return false
}

/*
PLOT X Y

this is the IR instruction for the draw opcode itself
*/
type PLOT struct {
	X, Y, H int
}

func (p PLOT) GetInstructionName() string {
	return "PLOT"
}

func (p PLOT) Opcodeable() bool {
	return true
}

func (p PLOT) usesVariableSpace() bool {
	return false
}

/*
SET Val ADDR

set a value at a given memory address
*/
type SET struct {
	Val, Addr int
}

func (s SET) GetInstructionName() string {
	return "SET"
}

func (s SET) Opcodeable() bool {
	return false
}

func (s SET) usesVariableSpace() bool {
	return true
}

//Generator contains all the basic information needed
//to transform an AST into a chip-8 ROM
type Generator struct {
	filename                                     string
	nodesConsumed                                int
	memorySize                                   int
	IRegisterIndex, plotXRegister, plotYRegister int
	Ir                                           []instruction
	memTable                                     memtable.MemTable
	regTable                                     registertable.RegisterTable
}

//NewGenerator inits the generator
func NewGenerator(filename string) *Generator {
	g := new(Generator)
	g.memTable = make(memtable.MemTable)
	g.regTable = make(registertable.RegisterTable)
	g.filename = filename
	g.nodesConsumed = 0
	g.memorySize = 4096 - 0x200 //0x200 is reserved space that we cannot use
	g.IRegisterIndex = 0xF
	g.plotXRegister = 0xE
	g.plotYRegister = 0xD
	g.regTable.Init()

	return g
}

/*
Generate interprets the AST and makes an IR from it
*/
func (g *Generator) Generate(AST []ast.Node) {
	for i := 0; i < len(AST); i++ {
		nodeType := AST[i].GetNodeName()
		switch nodeType {
		case "variable":
			variable := AST[i].(*ast.Variable)

			//check if its a reference
			if ast.NodeIsVariable(variable.Value) {
				g.Ir = append(g.Ir, g.newSetInstruction(variable, 0, true))
			} else {
				variableValue, _ := strconv.Atoi(variable.Value.(*ast.NumLit).Value)
				g.Ir = append(g.Ir, g.newSetInstruction(variable, variableValue, false))
			}
		case "statement":

		case "anb":

		case "function":

		case "functionCall":

		case "setStatement":

		case "mathStatement":

		case "comparison":

		case "switchStatement":

		case "plotStatement":
			plotStatement := AST[i].(*ast.PlotStatement)
			g.Ir = append(g.Ir, g.newPlotInstructionSet(plotStatement))
		}
	}

	g.wrapCodeInLoop()

	spew.Dump(g)
}

func (g *Generator) newPlotInstructionSet(plotStatement *ast.PlotStatement) *PLOT {
	/*
		chip-8's pixel placement system works on a 8x8 sprite.
		It draws whatever the byte represents in binary.
		so 10000000 (I.E 0x80) only displays one pixel in the top left corner
		and 00000000 displays an empty row.

		chip-8 allows for sprite overlap. Although it might be more
		effificient to add to the sprite already existing when we
		want to draw a new pixel in its 8x8 vacinity

		first we check if this pixel representor has been set
	*/
	plotInstr := new(PLOT)
	topLeftPixel := 0x80
	topLeftPixelMemoryName := "PIXEL_BUFFER_REP"

	/*
		Since we only draw a single pixel, the height of the sprite can always be one
	*/
	plotInstr.H = 1

	if g.memTable.LookupVariable(topLeftPixelMemoryName, true) == nil {
		spew.Dump("lookup failed")
		g.Ir = append(g.Ir, g.newSetInstructionFromLoose(topLeftPixelMemoryName, topLeftPixel))
	}
	pixelBufferVariable := g.memTable.LookupVariable(topLeftPixelMemoryName, true)

	/*
		fill the I register with the memory address of the single pixel value
		the emulator will read the sprite data from

		only do this if the I register is not there already
	*/
	if g.regTable[g.IRegisterIndex].Value != pixelBufferVariable.Addr {
		g.Ir = append(g.Ir, g.newMovInstructionFromLoose(g.IRegisterIndex, pixelBufferVariable.Addr, true))
	}

	/*
		actually set the I register
	*/
	IRegister := g.regTable[g.IRegisterIndex]
	IRegister.Value = pixelBufferVariable.Addr
	g.regTable[g.IRegisterIndex] = IRegister

	/*
		if the node uses variables, we will need to resolve those
		otherwise, we simply use the integer value of the plot statement
	*/
	if ast.NodeIsVariable(plotStatement.X) {
		variableName := plotStatement.X.(*ast.StatVar).Value
		variableTableEntry := g.memTable.LookupVariable(variableName, true)
		g.Ir = append(g.Ir, g.newMovInstructionFromLoose(g.plotXRegister, variableTableEntry.Value, false))
	} else {
		variableValue := plotStatement.X.(*ast.NumLit).Value
		intValue, _ := strconv.Atoi(variableValue)
		g.Ir = append(g.Ir, g.newMovInstructionFromLoose(g.plotXRegister, intValue, false))

	}

	/*
		Check y variable node
	*/
	if ast.NodeIsVariable(plotStatement.Y) {
		variableName := plotStatement.Y.(*ast.StatVar).Value
		variableTableEntry := g.memTable.LookupVariable(variableName, true)
		g.Ir = append(g.Ir, g.newMovInstructionFromLoose(g.plotYRegister, variableTableEntry.Value, false))
	} else {
		variableValue := plotStatement.Y.(*ast.NumLit).Value
		intValue, _ := strconv.Atoi(variableValue)
		g.Ir = append(g.Ir, g.newMovInstructionFromLoose(g.plotYRegister, intValue, false))
	}

	plotInstr.X = g.plotXRegister
	plotInstr.Y = g.plotYRegister

	return plotInstr
}

/*

	newMovInstructionFromLoose takes a loose set of values and turns them into
	a MOV instruction

	R1 must be a register
	R2 can either be a register or a memory address
	if R2 is a memory address the field "R2IsAddr" will be set to true
*/
func (g *Generator) newMovInstructionFromLoose(R1 int, R2 int, ANNN bool) MOV {
	instr := MOV{R1, R2, false}
	if memtable.IsValidMemRegion(R1) {
		errors.RegisterAdressModeFailure(R1)
	}

	instr.ANNN = ANNN
	return instr
}

/*
	same as newMovInstructionFromLoose.
	It does the extraction of variable values for you
*/
func (g *Generator) newMovInstruction(v *ast.Variable) MOV {
	instr := MOV{}
	return instr
}

/*
This is used for MEM operations.
This means that it is used for setting variables into memory before any opcode is executed
Once a variable is needed for an operation. for instance INC or ANB or any other. It will
be retreived by MOV and put into a free register
*/
func (g *Generator) newSetInstruction(v *ast.Variable, varValue int, resolve bool) SET {
	isntr := SET{}

	if resolve {
		resolutionName := v.Value.(*ast.StatVar).Value
		if val, ok := g.memTable[resolutionName]; ok {
			varValue = val.Value
		} else {
			errors.UndefinedVariableError(resolutionName)
		}
	}

	region := g.memTable.Put(v.Name, varValue)
	isntr.Addr = region.Addr
	isntr.Val = varValue
	return isntr
}

func (g *Generator) newSetInstructionFromLoose(name string, varValue int) SET {
	instr := SET{}
	region := g.memTable.Put(name, varValue)
	instr.Addr = region.Addr
	instr.Val = varValue
	return instr
}

func (g *Generator) newJumpInstructionFromLoose(from int, to int) JMP {
	return JMP{from, to}
}

/*
compressMemoryLayout relocates all variables next to the opcodes to reduce the size of the rom
*/
func (g *Generator) compressMemoryLayout() {
	variablesReplaced := 0

	//make sure the game does not start reading variable space
	g.wrapCodeInLoop()

	//get the end position of the opcodes
	endOpcodeSpace := len(g.Ir) * 2

	//move all variables closer
	for i := 0; i < len(g.Ir); i++ {
		if g.Ir[i].usesVariableSpace() {
			switch g.Ir[i].GetInstructionName() {
			case "SET":
				newPostion := endOpcodeSpace + variablesReplaced
				cast := g.Ir[i].(SET)
				cast.Addr = newPostion
				memTableVariable := g.memTable.FindByAddr(cast.Addr)
				g.memTable.Move(memTableVariable, newPostion, true)
				variablesReplaced++
				break
			case "MOV":
				newPostion := endOpcodeSpace + variablesReplaced
				cast := g.Ir[i].(MOV)
				cast.R2 = newPostion
				memTableVariable := g.memTable.FindByAddr(cast.R2)
				g.memTable.Move(memTableVariable, newPostion, true)
				variablesReplaced++
				break
			}
		}
	}
}

/*
   to make sure the machine does not start reading into variable space which is located after the opcodes
   it needs to execute, we need to move the PC back to the start of the progams opcodes

   this is done using a MOV call to set PC back to 0x200 which is the start of the opcode space
*/
func (g *Generator) wrapCodeInLoop() {
	opcodeSpaceEnd := len(g.Ir) * 4 //each instruction is 4 bytes
	resetInstruction := g.newJumpInstructionFromLoose(opcodeSpaceEnd, 0x200)

	g.Ir = append(g.Ir, resetInstruction)

}
