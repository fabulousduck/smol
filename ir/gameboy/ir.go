package ir

import "github.com/fabulousduck/smol/ast"

/*
Generator is a structure that contains information
required to generate an IR for the gameboy system
*/
type Generator struct {
}

/*
Init creates a new generator instance
that can be used to transform an AST into an
IR in gameboy format
*/
func Init() *Generator {
	g := new(Generator)
	return g
}

/*
Generate generates an IR from a given AST
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
