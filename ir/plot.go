package ir

import (
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/fabulousduck/smol/ast"
)

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

		//first we check if the variable is loaded into a register somewhere
		registerLoadedValue := g.regTable.Find(variableName)
		if registerLoadedValue == -1 {
			//if the variable is not register loaded
			registerLoadedValue = g.memTable.LookupVariable(variableName, false).Value
		} else {
			//if it is variable loaded
			registerValue := g.regTable[registerLoadedValue].Value
			g.Ir = append(g.Ir, g.newRegCpy(g.plotXRegister, registerValue))
		}

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
