package bytecode

import (
	"github.com/fabulousduck/smol/ast"
)

type instruction struct {
	Type string
	LHS  string
	MHS  string
	RHS  string
}

type register struct {
	variable *ast.Variable
}

type tuple struct {
	key   string
	value string
}

type stack []*tuple

//Generator contains all the basic information needed
//to transform an AST into a chip-8 ROM
type Generator struct {
	filename                     string
	nodesConsumed, registerCount int
	memorySize, memoryOffset     int
	esp, ebp                     int //esp is the global stack pointer, and ebp is the current stack frame pointer (local variable stack pointer)
	registers                    []*register
	ir                           []*instruction
	stack                        stack
}

//NewGenerator inits the generator
func NewGenerator(filename string) *Generator {
	g := new(Generator)
	g.filename = filename
	g.nodesConsumed = 0
	g.registerCount = 16
	g.memorySize = 4096
	g.memoryOffset = 512
	return g
}

//Compile generates a ROM from an AST by converting it into a IR and then byte code
func (g *Generator) Compile(AST []ast.Node) {
	for j := 0; j < len(AST); j++ {
		node := AST[j]
		nodeType := node.GetNodeName()
		switch nodeType {
		case "variable":
			//we can do this since only ints exist in our language
			g.ir = append(g.ir, g.generateMovInstruction(node.(*ast.Variable)))
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
}

/*
	desc: 6XNN
	operation: sets V[C] to NN


*/
func (g *Generator) generateMovInstruction(variable *ast.Variable) *instruction {
	movInstruction := new(instruction)
	stackVariable := new(tuple)
	movInstruction.Type = "MOV"
	variableValue := ""

	if ast.NodeIsVariable(variable) {
		//TODO
		//resolve value from stack
	} else {
		variableValue = variable.Value.(*ast.NumLit).Value
	}

	//throw it on the stack
	stackVariable.value = variableValue
	stackVariable.key = variable.Name

	g.stack = append(g.stack, stackVariable)

	//find a suitable register for the variable

	return movInstruction
}

func (g *Generator) interpretMovInstruction(mov *instruction) {

}

/*
	findSuitableRegister is a function that will be called when
	a value needs to be moved to a register.

	it is very strict about memory usage as it will throw an error
	if more registers are being used than allowed
*/
func (g *Generator) findSuitableRegister() int {
}
