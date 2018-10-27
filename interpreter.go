package losp

import (
	"fmt"
	"os"
)

type stack map[string]string

type interpreter struct {
	stack stack
}

func newInterpreter() *interpreter {
	i := new(interpreter)
	i.stack = make(stack)
	return i
}

func (i interpreter) interpret(ast []node) {
	for j := 0; j < len(ast); j++ {
		node := ast[j]
		nodeType := node.getNodeName()
		switch nodeType {
		case "variable":
			v := node.(*variable)
			//we can do this since only ints exist in our language
			i.stackAlloc(v)
		case "statement":
			s := node.(*statement)
			i.execStatement(s)
		case "anb":
		case "function":
		}
	}

	// spew.Dump(i.stack)
}

func (i *interpreter) stackAlloc(v *variable) {
	// fmt.Printf("STACK ALLOCATION K : %s, V : %s\n", v.name, v.value)
	i.stack[v.name] = v.value
}

func (i *interpreter) execStatement(s *statement) {
	switch s.lhs {
	case "PRI":
		if s.rhs.getNodeName() == "statVar" {
			rhs := s.rhs.(*statVar)
			value := i.stack.find(rhs.value)
			fmt.Printf("%s", value)
			return
		}

		fmt.Printf("%s", s.rhs.(*numLit).value)
		return
	case "PRU":
		// cast, _ := strconv.Atoi(s.rhs)
		// fmt.Printf("%c", cast)
	case "INC":

	}
}

func (s stack) find(key string) string {
	if val, ok := s[key]; ok {
		return val
	}
	undefinedVariableError(key)
	os.Exit(65)
	return ""
}
