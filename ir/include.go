package ir

import "github.com/fabulousduck/smol/ast"

func (g *Generator) createIncludeStatement(includeStatement *ast.IncludeStatement) {
	//dont bother with empty includes
	//maybe this should be caught at AST stage
	if includeStatement.Name == "" {
		return
	}

}
