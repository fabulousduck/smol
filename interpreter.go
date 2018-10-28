package losp

import (
	"fmt"
	"os"
	"strconv"
)

type tuple struct {
	key   string
	value string
}

type stack []*tuple

type stacks []stack

type interpreter struct {
	stacks stacks
}

func newInterpreter() *interpreter {
	i := new(interpreter)
	i.stacks = stacks{}
	baseStack := stack{}
	i.stacks = append(i.stacks, baseStack)
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
			i.stackAlloc(len(i.stacks)-1, v)
		case "statement":
			s := node.(*statement)
			i.execStatement(s)
		case "anb":
			anb := node.(*anb)
			i.execANB(anb)
		case "function":
		}
	}

	// spew.Dump(i.stack)
}

func (i *interpreter) stackAlloc(scopeLevel int, v *variable) {
	// fmt.Printf("STACK ALLOCATION K : %s, V : %s\n", v.name, v.value)
	stackTuple := new(tuple)
	stackTuple.key = v.name
	stackTuple.value = v.value
	i.stacks[scopeLevel] = append(i.stacks[scopeLevel], stackTuple)
}

func (i *interpreter) execANB(anb *anb) {
	var lhs *string
	var rhs *string

	if anb.lhs.getNodeName() == "statVar" {
		scopeLevel, index := i.stacks.find(anb.lhs.(*statVar).value)
		lhs = &i.stacks[scopeLevel][index].value
	} else {
		lhs = &anb.lhs.(*numLit).value
	}

	if anb.rhs.getNodeName() == "statVar" {
		scopeLevel, index := i.stacks.find(anb.rhs.(*statVar).value)
		rhs = &i.stacks[scopeLevel][index].value
	} else {
		rhs = &anb.rhs.(*numLit).value
	}
	scopedStack := stack{}
	i.stacks = append(i.stacks, scopedStack)
	scopeLevel := len(i.stacks)
	v, _ := strconv.Atoi(*lhs)
	n, _ := strconv.Atoi(*rhs)
	for v < n {
		// fmt.Printf(" %s < %s ", *lhs, *rhs)
		i.interpret(anb.body)
		v, _ = strconv.Atoi(*lhs)
		n, _ = strconv.Atoi(*rhs)
	}
	//GC the stacks that were used in the scoped block. ANB in this case
	i.stacks = i.stacks[scopeLevel:]
}

func (i *interpreter) execStatement(s *statement) {
	switch s.lhs {
	case "BRK":
		fmt.Printf("\n")
		return
	case "PRI":
		if s.rhs.getNodeName() == "statVar" {
			rhs := s.rhs.(*statVar)
			//scope level 0 is local block scope, and then works its way up
			scopeLevel, index := i.stacks.find(rhs.value)
			fmt.Printf("%s", i.stacks[scopeLevel][index].value)
			return
		}

		fmt.Printf("%s", s.rhs.(*numLit).value)
		return
	case "PRU":
		if s.rhs.getNodeName() == "statVar" {
			rhs := s.rhs.(*statVar)
			scopeLevel, index := i.stacks.find(rhs.value)
			cast, _ := strconv.Atoi(i.stacks[scopeLevel][index].value)
			fmt.Printf("%c", cast)
			return
		}

		cast, _ := strconv.Atoi(s.rhs.(*numLit).value)
		fmt.Printf("%c", cast)
		return

	//TODO: fucking hell make this nicer ryan good god - ryan
	case "INC":
		fmt.Println("INC")
		if s.rhs.getNodeName() != "statVar" {
			litIncrementError()
			os.Exit(65)
		}
		scopeLevel, index := i.stacks.find(s.rhs.(*statVar).value)
		vc, _ := strconv.Atoi(i.stacks[scopeLevel][index].value)
		vc++
		cast := strconv.Itoa(vc)
		i.stacks.set(scopeLevel, index, cast)
		return
	}
}

func (s stacks) set(scopeLevel int, index int, value string) {
	s[scopeLevel][index].value = value
}

func (s stacks) find(key string) (int, int) {
	//reverse stack search so we start at local scope and keep working our way up intill we find something

	for i := len(s) - 1; i > -1; i-- {
		stackIndex := s[i].stackContains(key)
		if stackIndex != -1 {
			//scopeLevel, scopedStackIndex
			return i, stackIndex
		}
	}

	undefinedVariableError(key)
	os.Exit(65)
	return -1, -1
}

func (s stack) stackContains(key string) int {
	for i := 0; i < len(s); i++ {
		if s[i].key == key {
			return i
		}
	}
	return -1
}
