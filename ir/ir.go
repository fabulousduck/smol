package ir

import (
	"os"

	"github.com/fabulousduck/smol/ast"
	"github.com/fabulousduck/smol/errors"
	"github.com/fabulousduck/smol/ir/functionaddrtable"
	"github.com/fabulousduck/smol/ir/memtable"
	"github.com/fabulousduck/smol/ir/registertable"

	"github.com/google/uuid"
)

//CPULayout are CPU specific variables that the ir needs.
//These are things like memory boundries and registers
type CPULayout map[string]int

//CH8CPULayout has specific offsets and variables for the CHIP8 CPU
var CH8CPULayout = CPULayout{
	"functionSpaceStart": 0xA00,
	"memorySize":         4094 - 0x200, //0x200 is reserved space on the CH8
	"IRegisterIndex":     0xF,
	"plotXRegister":      0xE,
	"plotYRegister":      0xD,
	"BNEXRegister":       0xC,
}

//Z80Layout has specific offsets and variables for the Z80 CPU
var Z80Layout = CPULayout{}

type instruction interface {
	GetInstructionName() string
	Opcodeable() bool
	usesVariableSpace() bool
}

//Generator contains all the basic information needed
//to transform an AST into a chip-8 ROM
type Generator struct {
	nodesConsumed     int
	target            string
	targetCPU         CPULayout
	Ir                []instruction
	memTable          memtable.MemTable
	regTable          registertable.RegisterTable
	functionAddrTable functionaddrtable.FunctionAddrTable
}

func getCPULayout(name string) CPULayout {

	switch name {
	case "CH8":
		return CH8CPULayout
		break
	case "Z80":
		return Z80Layout
	}

	errors.UnknownCPULayoutError(name)
	os.Exit(65)
	//this is just there so Go is happy about returns
	return CH8CPULayout
}

//NewGenerator inits the generator
func NewGenerator(target string) *Generator {
	g := new(Generator)
	g.memTable = make(memtable.MemTable)
	g.regTable = make(registertable.RegisterTable)
	g.nodesConsumed = 0
	g.target = target
	g.targetCPU = getCPULayout(target)
	g.regTable.Init()

	return g
}

/*
FindInstructionIndex looks up an instruction with given ID.
returns the first one it finds

returns -1 if nothing with that ID was found

currently only used for JMP instructions as those are the
only instructions that have ID's
*/
func (g *Generator) FindInstructionIndex(ID string) int {
	for i := 0; i < len(g.Ir); i++ {
		if g.Ir[i].GetInstructionName() == "Jump" {
			jumpInstrCast := g.Ir[i].(*Jump)
			if jumpInstrCast.ID == ID {
				return i
			}
		}
	}
	return -1
}

/*
Generate interprets the AST and makes an IR from it
*/
func (g *Generator) Generate(AST []ast.Node) {
	for i := 0; i < len(AST); i++ {
		node := AST[i]
		nodeType := node.GetNodeName()
		switch nodeType {
		case "variable":
			variable := node.(*ast.Variable)
			g.createVariableOperationInstructions(variable)
		case "statement":
			statement := node.(*ast.Statement)
			g.Ir = append(g.Ir, g.handleStatement(statement))
		case "function":
			instruction := node.(*ast.Function)
			g.createFunctionInstructions(instruction)
		case "directOperation":
			instruction := node.(*ast.DirectOperation)
			g.createDirectOperationInstructions(instruction)
		case "functionCall":
			instruction := node.(*ast.FunctionCall)
			g.createFunctionCallInstructions(instruction)
		case "setStatement":
			instruction := node.(*ast.SetStatement)
			g.createSetStatement(instruction)
		case "includeStatement":
			instruction := node.(*ast.IncludeStatement)
			g.createIncludeStatement(instruction)
		case "freeStatement":
			instruction := node.(*ast.FreeStatement)
			g.doFreeInstruction(instruction)
		case "switchStatement":

		case "plotStatement":
			plotStatement := AST[i].(*ast.PlotStatement)
			g.Ir = append(g.Ir, g.newPlotInstructionSet(plotStatement))
		}
	}
}

//doFreeInstruction does not actually embed a instruction to free a register
//it simply changes the internal compiler register table
func (g *Generator) doFreeInstruction(instruction *ast.FreeStatement) {
	variable := instruction.Variable.(*ast.StatVar)
	register := g.regTable.Find(variable.Value)
	g.regTable.PutRegisterValue(register, 0, "")
}

func (g *Generator) createFunctionInstructions(instruction *ast.Function) {

	//find the location where the function will be placed
	beforeGenerationInstructionCount := len(g.Ir)

	//create the jump instruction so it knows to jump over the function
	//when not called
	passJumpInstruction := g.newJumpInstructionFromLoose(0)

	uid := uuid.New()
	passJumpInstructionID := uid.String()
	passJumpInstruction.ID = passJumpInstructionID

	//save the byte addr before generating function code
	functionStartAddr := 0x200 + (beforeGenerationInstructionCount * 2)

	//put a new function on the function table so we know where can jump to to call it
	g.functionAddrTable = append(g.functionAddrTable, functionaddrtable.NewFunctionAddr(functionStartAddr, instruction.Name))

	//generate the function code
	g.Generate(instruction.Body)

	//find the jump back
	passJumpInstrIndex := g.FindInstructionIndex(passJumpInstructionID)

	//replace it with the new one that contains the proper address
	g.Ir[passJumpInstrIndex] = Jump{To: 0x200 + (len(g.Ir) * 2), ID: passJumpInstructionID}

	//put in a return statement
	g.Ir = append(g.Ir, g.newRetInstruction())
}

func (g *Generator) createDirectOperationInstructions(do *ast.DirectOperation) instruction {
	if !ast.NodeIsVariable(do.Variable) {
		errors.LitIncrementError()
		os.Exit(65)
	}
	rhsVariable := do.Variable.(*ast.StatVar)
	variableRegisterTableIndex := g.regTable.Find(rhsVariable.Value)
	if do.Operation == "++" {
		return g.newAddInstruction(variableRegisterTableIndex, 1)
	}

	return g.newSubInstruction(variableRegisterTableIndex, 1)
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
