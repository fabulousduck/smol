package ir

/*
compressMemoryLayout relocates all variables next to the opcodes to reduce the size of the rom
*/
func (g *Generator) compressMemoryLayout() {
	variablesReplaced := 0

	//make sure the game does not start reading variable space
	g.WrapCodeInLoop()

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
