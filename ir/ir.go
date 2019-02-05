package ir

import (
	"os"

	"github.com/fabulousduck/smol/ast"
	"github.com/fabulousduck/smol/errors"
	"github.com/fabulousduck/smol/ir/functionaddrtable"
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
	functionAddrTable                            functionaddrtable.FunctionAddrTable
	nodesConsumed                                int
	memorySize                                   int
	functionSpaceStart                           int
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
	g.functionSpaceStart = 0xA00
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
			g.createVariableOperationInstructions(variable)
		case "statement":
			statement := AST[i].(*ast.Statement)
			g.Ir = append(g.Ir, g.handleStatement(statement))
		case "anb":
			instruction := AST[i].(*ast.Anb)
			g.createAnbInstructions(instruction)
		case "function":
			instruction := AST[i].(*ast.Function)
			g.createFunctionInstructions(instruction)
		case "functionCall":
			instruction := AST[i].(*ast.FunctionCall)
			g.createFunctionCallInstructions(instruction)
		case "setStatement":
			instruction := AST[i].(*ast.SetStatement)
			g.createSetStatement(instruction)
		case "mathStatement":

		case "comparison":

		case "switchStatement":

		case "plotStatement":
			plotStatement := AST[i].(*ast.PlotStatement)
			g.Ir = append(g.Ir, g.newPlotInstructionSet(plotStatement))
		}
	}
}

func (g *Generator) createFunctionInstructions(instruction *ast.Function) {

	currentInstructionCount := len(g.Ir)
	//find the location where the function will be placed
	functionAddr := 0x200 + (currentInstructionCount * 2)

	//put a new function on the function table so we know where can jump to to call it
	g.functionAddrTable = append(g.functionAddrTable, functionaddrtable.NewFunctionAddr(functionAddr, instruction.Name))

	//generate the function code
	g.Generate(instruction.Body)

	//insert a jump over function. before the

	//put in a return statement
	g.Ir = append(g.Ir, g.newRetInstruction())
}

func (g *Generator) handleStatement(s *ast.Statement) instruction {
	var instr instruction
	switch s.LHS {
	case "INC":

		if !ast.NodeIsVariable(s.RHS) {
			errors.LitIncrementError()
			os.Exit(65)
		}
		rhsVariable := s.RHS.(*ast.StatVar)
		variableRegisterTableIndex := g.regTable.Find(rhsVariable.Value)
		instr = g.newAddInstruction(variableRegisterTableIndex, 1)
	}
	return instr
}
