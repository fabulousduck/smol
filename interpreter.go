package losp

import (
	"fmt"
	"os"
	"strconv"
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
	case "BRK":
		fmt.Printf("\n")
		return
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
		if s.rhs.getNodeName() == "statVar" {
			rhs := s.rhs.(*statVar)
			value := i.stack.find(rhs.value)
			cast, _ := strconv.Atoi(value)
			fmt.Printf("%c", cast)
			return
		}

		cast, _ := strconv.Atoi(s.rhs.(*numLit).value)
		fmt.Printf("%c", cast)
		return

	case "INC":
		if s.rhs.getNodeName() != "statVar" {
			litIncrementError()
			os.Exit(65)
		}
		v := i.stack.find(s.rhs.(*statVar).value)
		vc, _ := strconv.Atoi(v)
		vc++
		if _, ok := i.stack[s.rhs.(*statVar).value]; ok {
			i.stack[s.rhs.(*statVar).value] = strconv.Itoa(vc)
		}
		return
	}
}

func (s *stack) set(key string, value string) {

}

func (s stack) find(key string) string {
	if val, ok := s[key]; ok {
		return val
	}
	undefinedVariableError(key)
	os.Exit(65)
	return ""
}
