package smol

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

type tuple struct {
	key   string
	value string
}

type stack []*tuple

type stacks []stack

type heap []*function //not really a heap since its not dynamically sized types, but a good excuse to put my function decls into

type interpreter struct {
	stacks stacks
	heap   heap
}

func newInterpreter() *interpreter {
	i := new(interpreter)
	i.stacks = stacks{}
	i.heap = heap{}
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
			function := node.(*function)
			i.execFunctionDecl(function)
		case "functionCall":
			fc := node.(*functionCall)
			i.execFunctionCall(fc)
		case "setStatement":
			ss := node.(*setStatement)
			i.setVariableValue(ss)
		case "mathStatement":
			ms := node.(*mathStatement)
			i.execMathStatement(ms)
		case "comparison":
			cm := node.(*comparison)
			i.execComparison(cm)
		}
	}
}

func (i *interpreter) execComparison(cm *comparison) {

	clhs := 0
	crhs := 0
	beforeScopeLevel := len(i.stacks)
	scopedStack := stack{}

	if cm.lhs.getNodeName() == "statVar" {
		scopeLevel, index := i.stacks.find(cm.lhs.(*statVar).value)
		clhs, _ = strconv.Atoi(i.stacks[scopeLevel][index].value)
	} else {
		clhs, _ = strconv.Atoi(cm.lhs.(*numLit).value)
	}

	if cm.rhs.getNodeName() == "statVar" {
		scopeLevel, index := i.stacks.find(cm.rhs.(*statVar).value)
		crhs, _ = strconv.Atoi(i.stacks[scopeLevel][index].value)
	} else {
		crhs, _ = strconv.Atoi(cm.rhs.(*numLit).value)
	}

	// do static analysis on same variable comparisons
	switch cm.operator {
	case "LT":
		if clhs < crhs {
			i.stacks = append(i.stacks, scopedStack)
			i.interpret(cm.body)
		}
	case "GT":
		if clhs > crhs {
			i.stacks = append(i.stacks, scopedStack)
			i.interpret(cm.body)
		}
	case "EQ":
		if clhs == crhs {
			i.stacks = append(i.stacks, scopedStack)
			i.interpret(cm.body)
		}
	case "NEQ":
		if clhs != crhs {
			i.stacks = append(i.stacks, scopedStack)
			i.interpret(cm.body)
		}
	}
	i.stacks = i.stacks[:beforeScopeLevel]
	return
}

func (i *interpreter) execMathStatement(ms *mathStatement) {
	operator := ms.lhs
	if ms.mhs.getNodeName() != "statVar" {
		additionInvalidReceiverError()
	}

	receiverVariableName := ms.mhs.(*statVar).value
	receiverVariableScopeLevel, receiverVariableIndex := i.stacks.find(receiverVariableName)
	receiverVariableValue := i.stacks[receiverVariableScopeLevel][receiverVariableIndex].value
	result := ""
	if ms.rhs.getNodeName() == "statVar" {
		scopeLevel, index := i.stacks.find(ms.rhs.(*statVar).value)
		rhs := i.stacks[scopeLevel][index].value
		result = evalMathExpression(operator, receiverVariableValue, rhs)
	} else {
		rhs := ms.rhs.(*numLit).value
		result = evalMathExpression(operator, receiverVariableValue, rhs)
	}

	i.stacks.set(receiverVariableScopeLevel, receiverVariableIndex, result)

}

func evalMathExpression(expressionType string, lhs string, rhs string) string {
	clhs, _ := strconv.Atoi(lhs)
	crhs, _ := strconv.Atoi(rhs)
	switch expressionType {
	case "ADD":
		return strconv.Itoa(clhs + crhs)
	case "SUB":
		return strconv.Itoa(clhs - crhs)
	case "MUL":
		return strconv.Itoa(clhs * crhs)
	case "DIV":
		return strconv.Itoa(clhs / crhs)
	case "POW":
		return strconv.Itoa(int(math.Pow(float64(clhs), float64(crhs))))
	}
	//not sure what to return here
	//TODO: figure above out and apply accordingly
	return rhs
}

func (i *interpreter) setVariableValue(ss *setStatement) {
	if ss.mhs.getNodeName() != "statVar" {
		litAssignError()
		os.Exit(65)
	}

	scopeLevel, index := i.stacks.find(ss.mhs.(*statVar).value)
	if ss.rhs.getNodeName() == "statVar" {
		rhsScopeLevel, rhsIndex := i.stacks.find(ss.rhs.(*statVar).value)
		i.stacks[scopeLevel][index].value = i.stacks[rhsScopeLevel][rhsIndex].value
		return
	}
	i.stacks[scopeLevel][index].value = ss.rhs.(*numLit).value

}

func (i *interpreter) execFunctionCall(fc *functionCall) {
	functionDecl := i.heap[i.heap.find(fc.name)]
	if len(fc.args) != len(functionDecl.params) {

		incorrectFunctionParamCountError(functionDecl.name, len(fc.args), len(functionDecl.params))
		os.Exit(65)
		return
	}
	beforeScopeLevel := len(i.stacks)
	scopedStack := stack{}
	for j := 0; j < len(functionDecl.params); j++ {
		if determineStringType(fc.args[j]) == "CHAR" {
			scopeLevel, index := i.stacks.find(fc.args[j])
			value := i.stacks[scopeLevel][index].value
			scopedStack = append(scopedStack, &tuple{key: functionDecl.params[j], value: value})
			continue
		}
		scopedStack = append(scopedStack, &tuple{key: functionDecl.params[j], value: fc.args[j]})
	}
	i.stacks = append(i.stacks, scopedStack)
	i.interpret(functionDecl.body)
	i.stacks = i.stacks[:beforeScopeLevel]
}

func (i *interpreter) execFunctionDecl(f *function) {
	i.heap = append(i.heap, f)
}

func (i *interpreter) stackAlloc(scopeLevel int, v *variable) {
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
		if s.rhs.getNodeName() != "statVar" {
			litIncrementError()
			os.Exit(65)
		}

		scopeLevel, index := i.stacks.find(s.rhs.(*statVar).value)
		vc, _ := strconv.Atoi(i.stacks[scopeLevel][index].value)
		vc++
		i.stacks.set(scopeLevel, index, strconv.Itoa(vc))
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

func (h heap) find(name string) int {
	for i := 0; i < len(h); i++ {
		if h[i].name == name {
			return i
		}
	}
	undefinedFunctionReferenceError(name)
	os.Exit(65)
	return -1
}
