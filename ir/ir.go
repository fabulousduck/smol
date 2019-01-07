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
Mov R1 to R2
*/
type MOV struct {
	R1, R2 int
}

func (m MOV) GetInstructionName() string {
	return "MOV"
}

func (m MOV) Opcodeable() bool {
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
	val, addr int
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
	filename                 string
	nodesConsumed            int
	memorySize, memoryOffset int
	ir                       []instruction
	memTable                 memtable.MemTable
	regTable                 registertable.RegisterTable
}

//NewGenerator inits the generator
func NewGenerator(filename string) *Generator {
	g := new(Generator)
	g.memTable = make(memtable.MemTable)
	g.regTable = make(registertable.RegisterTable)
	g.filename = filename
	g.nodesConsumed = 0
	g.memorySize = 4096
	g.memoryOffset = 512
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
				g.ir = append(g.ir, g.newSetInstruction(variable, 0, true))
			} else {
				variableValue, _ := strconv.Atoi(variable.Value.(*ast.NumLit).Value)
				g.ir = append(g.ir, g.newSetInstruction(variable, variableValue, false))
			}
		case "statement":

		case "anb":

		case "function":

		case "functionCall":

		case "setStatement":

		case "mathStatement":

		case "comparison":

		case "switchStatement":

		}
	}

	spew.Dump(g)
}

/*
	used for FX65
*/
func (g *Generator) newLDRInstruction(v *ast.Variable) LDR {
	instr := LDR{}

	return instr
}

/*
	used for 6XNN
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
	isntr.addr = region.Addr
	isntr.val = varValue
	return isntr
}
