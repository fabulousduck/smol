package ir

import (
	"strconv"

	"github.com/davecgh/go-spew/spew"

	"github.com/fabulousduck/smol/ast"
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
