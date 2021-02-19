package ir

import (
	"os"
	"strconv"

	"github.com/fabulousduck/smol/ast"
	"github.com/fabulousduck/smol/errors"
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

func (g *Generator) newPlotInstructionSet(plotStatement *ast.PlotStatement) PLOT {
	targetCPU := g.targetCPU

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
	plotInstr := PLOT{}
	topLeftPixel := 0x80
	topLeftPixelMemoryName := "PIXEL_BUFFER_REP"

	/*
		Since we only draw a single pixel, the height of the sprite can always be one
	*/
	plotInstr.H = 1

	/*
		check if the single pixel has been set or not
	*/
	if g.memTable.LookupVariable(topLeftPixelMemoryName, true) == nil {
		g.Ir = append(g.Ir, g.newSetMemoryLocationFromLoose(topLeftPixelMemoryName, topLeftPixel))
	}
	pixelBufferVariable := g.memTable.LookupVariable(topLeftPixelMemoryName, true)

	/*
		fill the I register with the memory address of the single pixel value
		the emulator will read the sprite data from

		only do this if the I register is not there already
	*/
	if g.regTable[targetCPU["IRegisterIndex"]].Value != pixelBufferVariable.Addr {
		g.Ir = append(g.Ir, g.newMovInstructionFromLoose(targetCPU["IRegisterIndex"], pixelBufferVariable.Addr, true))
	}

	/*
		actually set the I register
	*/
	IRegister := g.regTable[targetCPU["IRegisterIndex"]]
	IRegister.Value = pixelBufferVariable.Addr
	g.regTable[targetCPU["IRegisterIndex"]] = IRegister

	/*
		if the node uses variables, we will need to resolve those
		otherwise, we simply use the integer value of the plot statement
	*/
	if ast.NodeIsVariable(plotStatement.X) {
		variableName := plotStatement.X.(*ast.StatVar).Value
		//first we check if the variable is loaded into a register somewhere
		registerLoadedValue := g.regTable.Find(variableName)
		if registerLoadedValue == -1 {
			//if it is not in any register. the variable does not exist and we error out
			errors.UndefinedVariableError(variableName)
			os.Exit(65)
		} else {
			//if it is variable loaded
			g.Ir = append(g.Ir, g.newRegCpy(registerLoadedValue, targetCPU["plotXRegister"]))
		}
	} else {
		variableValue := plotStatement.X.(*ast.NumLit).Value
		intValue, _ := strconv.Atoi(variableValue)
		g.Ir = append(g.Ir, g.newSpecificRegisterSet(targetCPU["plotXRegister"], intValue, "plotXRegister"))
	}

	if ast.NodeIsVariable(plotStatement.Y) {
		variableName := plotStatement.Y.(*ast.StatVar).Value

		//first we check if the variable is loaded into a register somewhere
		registerLoadedValue := g.regTable.Find(variableName)
		if registerLoadedValue == -1 {
			//if it is not in any register. the variable does not exist and we error out
			errors.UndefinedVariableError(variableName)
			os.Exit(65)
		} else {
			//if it is variable loaded
			g.Ir = append(g.Ir, g.newRegCpy(registerLoadedValue, targetCPU["plotYRegister"]))
		}
	} else {
		variableValue := plotStatement.Y.(*ast.NumLit).Value
		intValue, _ := strconv.Atoi(variableValue)
		g.Ir = append(g.Ir, g.newSpecificRegisterSet(targetCPU["plotYRegister"], intValue, "plotYRegister"))
	}

	plotInstr.X = targetCPU["plotXRegister"]
	plotInstr.Y = targetCPU["plotYRegister"]
	return plotInstr
}
