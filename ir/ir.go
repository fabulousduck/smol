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

/*
LDR MR to R2
*/
type LDR struct {
	MR, R2 int
}

func (l LDR) GetInstructionName() string {
	return "LDR"
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

//Generator contains all the basic information needed
//to transform an AST into a chip-8 ROM
type Generator struct {
	filename                 string
	nodesConsumed            int
	memorySize, memoryOffset int
	ir                       []instruction
	memTable                 *memtable.MemTable
	regTable                 *registertable.RegisterTable
}

//NewGenerator inits the generator
func NewGenerator(filename string) *Generator {
	g := new(Generator)
	g.memTable = new(memtable.MemTable)
	g.regTable = new(registertable.RegisterTable)
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
				g.ir = append(g.ir, g.newMovInstruction(variable))
			} else {
				variableValue, _ := strconv.Atoi(variable.Value.(*ast.NumLit).Value)
				g.ir = append(g.ir, g.newSetInstruction(variable, variableValue))
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

	spew.Dump(g.ir)
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

func (g *Generator) newSetInstruction(v *ast.Variable, varValue int) SET {
	isntr := SET{}
	region := g.memTable.Put(v.Name, varValue)
	isntr.addr = region.Addr
	isntr.val = varValue
	return isntr
}
