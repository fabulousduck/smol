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
}

/*
Mov R1 into R2

Chip-8 knows MOV in the form of 6XNN
where X is the register and NN is the memory address
*/
type MOV struct {
	R1, R2   int
	R2IsAddr bool
}

func (m MOV) GetInstructionName() string {
	return "MOV"
}

func (m MOV) Opcodeable() bool {
	return true
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

/*
LDR MR to R2
*/
type LDR struct {
	MR, R2 int
}

func (l LDR) GetInstructionName() string {
	return "LDR"
}

func (l LDR) Opcodeable() bool {
	return true
}

/*
SET val to addr
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

//Generator contains all the basic information needed
//to transform an AST into a chip-8 ROM
type Generator struct {
	filename       string
	nodesConsumed  int
	memorySize     int
	IRegisterIndex int
	Ir             []instruction
	memTable       memtable.MemTable
	regTable       registertable.RegisterTable
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
		g.Ir = append(g.Ir, g.newSetInstructionFromLoose(topLeftPixelMemoryName, topLeftPixel))
	}
	pixelBufferVariable := g.memTable.LookupVariable(topLeftPixelMemoryName, true)

	/*
		fill the I register with the memory address of the single pixel value
		the emulator will read the sprite data from
	*/
	g.Ir = append(g.Ir, g.newMovInstructionFromLoose(g.IRegisterIndex, pixelBufferVariable.Addr))

	/*
		actually set the I register
	*/
	IRegister := g.regTable[g.IRegisterIndex]
	IRegister.Value = pixelBufferVariable.Addr
	g.regTable[g.IRegisterIndex] = IRegister

	/*
		check if the nodes in the plot statement has variables that need to be resolved
		and if so, resolve statement

		Check x variable node
	*/
	if ast.NodeIsVariable(plotStatement.X) {
		variableName := plotStatement.X.(*ast.StatVar).Value
		variableTableEntry := g.memTable.LookupVariable(variableName, true)
		plotInstr.X = variableTableEntry.Value
	} else {
		variableValue := plotStatement.X.(*ast.NumLit).Value
		intValue, _ := strconv.Atoi(variableValue)
		plotInstr.X = intValue
	}

	/*
		Check y variable node
	*/
	if ast.NodeIsVariable(plotStatement.Y) {
		variableName := plotStatement.Y.(*ast.StatVar).Value
		variableTableEntry := g.memTable.LookupVariable(variableName, true)
		plotInstr.Y = variableTableEntry.Value
	} else {
		variableValue := plotStatement.Y.(*ast.NumLit).Value
		intValue, _ := strconv.Atoi(variableValue)
		plotInstr.Y = intValue
	}

	return plotInstr
}

func (g *Generator) newISetInstruction(address int) {

}

/*
	used for FX65
*/
func (g *Generator) newLDRInstruction(v *ast.Variable) LDR {
	instr := LDR{}

	return instr
}

/*

	newMovInstructionFromLoose takes a loose set of values and turns them into
	a MOV instruction

	R1 must be a register
	R2 can either be a register or a memory address
	if R2 is a memory address the field "R2IsAddr" will be set to true
*/
func (g *Generator) newMovInstructionFromLoose(R1 int, R2 int) MOV {
	instr := MOV{R1, R2, false}
	if memtable.IsValidMemRegion(R1) {
		errors.RegisterAdressModeFailure(R1)
	}

	if memtable.IsValidMemRegion(R2) {
		instr.R2IsAddr = true
	}
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
