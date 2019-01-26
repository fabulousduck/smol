package ir

import (
	"fmt"
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

//Generator contains all the basic information needed
//to transform an AST into a chip-8 ROM
type Generator struct {
	filename                                     string
	nodesConsumed                                int
	memorySize                                   int
	IRegisterIndex, plotXRegister, plotYRegister int
	BNEXRegister                                 int
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
	g.BNEXRegister = 0xC
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
			statement := AST[i].(*ast.Statement)
			g.Ir = append(g.Ir, g.handleStatement(statement))
		case "anb":
			instruction := AST[i].(*ast.Anb)
			if ast.NodeIsVariable(instruction.LHS) {
				variable := instruction.LHS.(*ast.StatVar)
				bnexreg := registertable.Register{g.memTable.LookupVariable(variable.Value, true).Value, variable.Value}
				g.regTable[g.BNEXRegister] = bnexreg
			} else {
				lhsValue, _ := strconv.Atoi(instruction.LHS.(*ast.NumLit).Value)
				bnexreg := registertable.Register{lhsValue, ""}
				g.regTable[g.BNEXRegister] = bnexreg

			}

			anbInstructionStart := 0x200 + (len(g.Ir) * 2) // we need to multiply by 2 because each instruction is 2 bytes long
			fmt.Printf("NIGGERS")
			spew.Dump(anbInstructionStart)
			g.Generate(instruction.Body)
			g.Ir = append(g.Ir, g.newBNEInstruction(instruction))
			g.Ir = append(g.Ir, g.newJumpInstructionFromLoose(anbInstructionStart))
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

}

func (g *Generator) handleStatement(statement *ast.Statement) instruction {
	var newInstr instruction
	switch statement.LHS {
	case "INC":
		if ast.NodeIsVariable(statement.RHS) {
			registerIndex := g.regTable.Find(statement.RHS.(*ast.StatVar).Value)
			if registerIndex == -1 {
				errors.ROMModError()
			}
			newRegState := registertable.Register{g.regTable[registerIndex].Value + 1, g.regTable[registerIndex].Name}
			g.regTable[registerIndex] = newRegState
			newInstr = g.newAddInstruction(registerIndex, 1)
		}
	}
	return newInstr
}

/*
   to make sure the machine does not start reading into variable space which is located after the opcodes
   it needs to execute, we need to move the PC back to the start of the progams opcodes

   this is done using a MOV call to set PC back to 0x200 which is the start of the opcode space
*/
func (g *Generator) WrapCodeInLoop() {
	resetInstruction := g.newJumpInstructionFromLoose(0x200)
	g.Ir = append(g.Ir, resetInstruction)
}
