package ir

import (
	"fmt"
	"os"
	"strconv"

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
	functionAddrTable                            []functionaddrtable.FunctionAddr
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
				//if it is a reference, we get the original value,
				//and copy it over into a new register with the name of the new variable
				variableValue := variable.Value.(*ast.StatVar)
				emptyRegister := g.regTable.FindEmptyRegister()
				originalRegister := g.regTable.Find(variableValue.Value)
				fmt.Printf("making register dupe with new name: %s", variable.Name)
				g.regTable[emptyRegister] = registertable.Register{g.regTable[originalRegister].Value, variable.Name}
				g.Ir = append(g.Ir, g.newRegCpy(originalRegister, emptyRegister))
			} else {
				variableValue, _ := strconv.Atoi(variable.Value.(*ast.NumLit).Value)
				g.Ir = append(g.Ir, g.newSetRegisterInstructionFromLoose(variable.Name, variableValue))
			}
		case "statement":
			statement := AST[i].(*ast.Statement)
			g.Ir = append(g.Ir, g.handleStatement(statement))
		case "anb":
			instruction := AST[i].(*ast.Anb)
			// we need to multiply by 2 because each instruction is 2 bytes long
			//doing this here is a little trick to get the value to keep getting updated from its origin register
			anbInstructionStart := 0x200 + (len(g.Ir) * 2)
			g.Generate(instruction.Body)

			if ast.NodeIsVariable(instruction.LHS) {
				variable := instruction.LHS.(*ast.StatVar)
				lhsVariableRegister := g.regTable.Find(variable.Value)
				g.regTable[g.BNEXRegister] = registertable.Register{g.regTable[lhsVariableRegister].Value, "BNEX"}
				g.Ir = append(g.Ir, g.newRegCpy(lhsVariableRegister, g.BNEXRegister))
			} else {
				lhsValue, _ := strconv.Atoi(instruction.LHS.(*ast.NumLit).Value)
				bnexreg := registertable.Register{lhsValue, "BNEX"}
				g.regTable[g.BNEXRegister] = bnexreg
				g.Ir = append(g.Ir, g.newSpecificRegisterSet(g.BNEXRegister, lhsValue, "BNEX"))
			}

			if ast.NodeIsVariable(instruction.RHS) {
				variable := instruction.RHS.(*ast.StatVar)
				rhsVariableRegister := g.regTable.Find(variable.Value)
				g.Ir = append(g.Ir, g.newBNERRInstructionFromLoose(g.BNEXRegister, rhsVariableRegister))
			} else {
				rhsValue, _ := strconv.Atoi(instruction.RHS.(*ast.NumLit).Value)
				g.Ir = append(g.Ir, g.newBNEInstructionFromLoose(g.BNEXRegister, rhsValue))
			}
			g.Ir = append(g.Ir, g.newJumpInstructionFromLoose(anbInstructionStart))
		case "function":
			instruction := AST[i].(*ast.Function)

			//find the location where the function will be placed
			functionAddr := 0x200 + (len(g.Ir) * 2)

			//put a new function on the function table so we know where can jump to to call it
			g.functionAddrTable = append(g.functionAddrTable, functionaddrtable.NewFunctionAddr(functionAddr, instruction.Name))

			//put the current register contents in memory so we dont lose it
			g.Ir = append(g.Ir, g.NewRGDInstruction())
		case "functionCall":

		case "setStatement":
			instruction := AST[i].(*ast.SetStatement)
			castVariable := instruction.MHS.(*ast.StatVar)
			variableRegister := g.regTable.Find(castVariable.Value)
			if ast.NodeIsVariable(instruction.RHS) {
				referenceVariableRegister := g.regTable.Find(instruction.RHS.(*ast.StatVar).Value)
				g.Ir = append(g.Ir, g.newRegCpy(referenceVariableRegister, variableRegister))
			} else {
				castVal, _ := strconv.Atoi(instruction.RHS.(*ast.NumLit).Value)
				g.Ir = append(g.Ir, g.newSetRegisterInstructionFromLoose(castVariable.Value, castVal))

			}
		case "mathStatement":

		case "comparison":

		case "switchStatement":

		case "plotStatement":
			plotStatement := AST[i].(*ast.PlotStatement)
			g.Ir = append(g.Ir, g.newPlotInstructionSet(plotStatement))
		}
	}

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

/*
compressMemoryLayout relocates all variables next to the opcodes to reduce the size of the rom
*/
// func (g *Generator) compressMemoryLayout() {
// 	variablesReplaced := 0

// 	//make sure the game does not start reading variable space
// 	g.wrapCodeInLoop()

// 	//get the end position of the opcodes
// 	endOpcodeSpace := len(g.Ir) * 2

// 	//move all variables closer
// 	for i := 0; i < len(g.Ir); i++ {
// 		if g.Ir[i].usesVariableSpace() {
// 			switch g.Ir[i].GetInstructionName() {
// 			case "SET":
// 				newPostion := endOpcodeSpace + variablesReplaced
// 				cast := g.Ir[i].(SET)
// 				cast.Addr = newPostion
// 				memTableVariable := g.memTable.FindByAddr(cast.Addr)
// 				g.memTable.Move(memTableVariable, newPostion, true)
// 				variablesReplaced++
// 				break
// 			case "MOV":
// 				newPostion := endOpcodeSpace + variablesReplaced
// 				cast := g.Ir[i].(MOV)
// 				cast.R2 = newPostion
// 				memTableVariable := g.memTable.FindByAddr(cast.R2)
// 				g.memTable.Move(memTableVariable, newPostion, true)
// 				variablesReplaced++
// 				break
// 			}
// 		}
// 	}
// }

/*
   to make sure the machine does not start reading into variable space which is located after the opcodes
   it needs to execute, we need to move the PC back to the start of the progams opcodes

   this is done using a MOV call to set PC back to 0x200 which is the start of the opcode space
*/
func (g *Generator) WrapCodeInLoop() {
	resetInstruction := g.newJumpInstructionFromLoose(0x200)
	g.Ir = append(g.Ir, resetInstruction)
}
