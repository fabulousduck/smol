package interpreter

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/fabulousduck/smol/ast"
	"github.com/fabulousduck/smol/errors"
)

type tuple struct {
	key   string
	value string
}

type stack []*tuple

//Stacks is the global scope that hold sub scopes for varianbles
type Stacks []stack

//Heap is not really a heap since it does not hold dynamically sized types, but a good excuse to put my function decls into
type Heap []*ast.Function

//Interpreter contains all data needed to Interpret an AST
type Interpreter struct {
	Stacks Stacks
	Heap   Heap
}

//NewInterpreter provides a new interpreter with empty base stack and heap
func NewInterpreter() *Interpreter {
	i := new(Interpreter)
	i.Stacks = Stacks{}
	i.Heap = Heap{}
	baseStack := stack{}
	i.Stacks = append(i.Stacks, baseStack)
	return i
}

//Interpret will tree walk execute an AST from left to right (topdown)
func (i Interpreter) Interpret(AST []ast.Node) {
	for j := 0; j < len(AST); j++ {
		node := AST[j]
		nodeType := node.GetNodeName()
		switch nodeType {
		case "variable":
			//we can do this since only ints exist in our language
			i.stackAlloc(len(i.Stacks)-1, node.(*ast.Variable))
		case "statement":
			i.execStatement(node.(*ast.Statement))
		case "anb":
			i.execANB(node.(*ast.Anb))
		case "function":
			i.execFunctionDecl(node.(*ast.Function))
		case "functionCall":
			i.execFunctionCall(node.(*ast.FunctionCall))
		case "setStatement":
			i.setVariableValue(node.(*ast.SetStatement))
		case "mathStatement":
			i.execMathStatement(node.(*ast.MathStatement))
		case "comparison":
			i.execComparison(node.(*ast.Comparison))
		case "switchStatement":
			i.execSwitchStatement(node.(*ast.SwitchStatement))
		}
	}
}

func (i *Interpreter) execSwitchStatement(ss *ast.SwitchStatement) {
	var defaultCase []ast.Node
	var caseMatchValue string
	var matchValue string

	if ast.NodeIsVariable(ss.MatchValue) {
		matchValue = i.Stacks.resolveVariable(ss.MatchValue).value
	} else {
		matchValue = ss.MatchValue.(*ast.NumLit).Value
	}

	for _, switchCase := range ss.Cases {
		if switchCase.GetNodeName() != "switchCase" && switchCase.GetNodeName() != "end_of_switch" {
			errors.UnknownSwitchNode()
			os.Exit(65)
		}
		if switchCase.GetNodeName() == "end_of_switch" {
			defaultCase = switchCase.(*ast.Eos).Body
			continue
		}

		currentCase := switchCase.(*ast.SwitchCase)
		caseMatchValue = i.Stacks.resolveValue(currentCase.MatchValue)

		if matchValue == caseMatchValue {
			i.Interpret(currentCase.Body)
			return
		}
	}

	if defaultCase != nil {
		i.Interpret(defaultCase)
	}

	return
}

func (i *Interpreter) execComparison(cm *ast.Comparison) {

	clhs := 0
	crhs := 0
	beforeScopeLevel := len(i.Stacks)

	clhs, _ = strconv.Atoi(i.Stacks.resolveValue(cm.LHS))
	crhs, _ = strconv.Atoi(i.Stacks.resolveValue(cm.RHS))

	//create a stack for the block inside the comparisons body
	i.Stacks = append(i.Stacks, stack{})

	// do static analysis on same variable comparisons
	switch cm.Operator {
	case "LT":
		if clhs < crhs {
			i.Interpret(cm.Body)
		}
	case "GT":
		if clhs > crhs {
			i.Interpret(cm.Body)
		}
	case "EQ":
		if clhs == crhs {
			i.Interpret(cm.Body)
		}
	case "NEQ":
		if clhs != crhs {
			i.Interpret(cm.Body)
		}
	}

	i.Stacks = i.Stacks[:beforeScopeLevel]
	return
}

func (i *Interpreter) execMathStatement(ms *ast.MathStatement) {
	operator := ms.LHS

	if ms.MHS.GetNodeName() != "statVar" {
		errors.MathInvalidReceiverError()
	}

	receiverVariable := i.Stacks.resolveVariable(ms.MHS)
	combinatorValue := i.Stacks.resolveValue(ms.RHS)

	i.Stacks.set(receiverVariable.key, evalMathExpression(operator, receiverVariable.value, combinatorValue))

}

func evalMathExpression(expressionType string, LHS string, RHS string) string {
	clhs, _ := strconv.Atoi(LHS)
	crhs, _ := strconv.Atoi(RHS)
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
	return RHS
}

func (i *Interpreter) setVariableValue(ss *ast.SetStatement) {
	if ss.MHS.GetNodeName() != "statVar" {
		errors.LitAssignError()
		os.Exit(65)
	}

	receiverVariable := i.Stacks.resolveVariable(ss.MHS)
	rhs := i.Stacks.resolveValue(ss.RHS)

	i.Stacks.set(receiverVariable.key, rhs)

}

func (i *Interpreter) execFunctionCall(fc *ast.FunctionCall) {
	functionDecl := i.Heap.resolveFunction(fc.Name)

	if len(fc.Args) != len(functionDecl.Params) {
		errors.IncorrectFunctionParamCountError(functionDecl.Name, len(fc.Args), len(functionDecl.Params))
		os.Exit(65)
		return
	}

	beforeScopeLevel := len(i.Stacks)
	scopedStack := stack{}

	for paramListIndex, param := range functionDecl.Params {

		if ast.NodeIsVariable(fc.Args[paramListIndex]) {
			scopedStack = append(scopedStack, i.Stacks.resolveVariable(fc.Args[paramListIndex]))
		} else {
			scopedStack = append(scopedStack, &tuple{key: param, value: fc.Args[paramListIndex].(*ast.NumLit).Value})
		}
	}

	//left off here refactoring

	i.Stacks = append(i.Stacks, scopedStack)
	i.Interpret(functionDecl.Body)
	i.Stacks = i.Stacks[:beforeScopeLevel]
}

func (i *Interpreter) execFunctionDecl(f *ast.Function) {
	i.Heap = append(i.Heap, f)
}

func (i *Interpreter) stackAlloc(scopeLevel int, v *ast.Variable) {
	stackTuple := new(tuple)
	stackTuple.key = v.Name

	stackTuple.value = i.Stacks.resolveValue(v.Value)

	i.Stacks[scopeLevel] = append(i.Stacks[scopeLevel], stackTuple)
}

func (i *Interpreter) execANB(anb *ast.Anb) {
	LHS := i.Stacks.resolvePtrValue(anb.LHS)
	RHS := i.Stacks.resolvePtrValue(anb.RHS)

	i.Stacks = append(i.Stacks, stack{})
	scopeLevel := len(i.Stacks)
	v, _ := strconv.Atoi(*LHS)
	n, _ := strconv.Atoi(*RHS)
	for v != n {

		i.Interpret(anb.Body)
		//check if it works without htis
		v, _ = strconv.Atoi(*LHS)
		n, _ = strconv.Atoi(*RHS)
	}
	//GC the Stacks that were used in the scoped block. ANB in this case
	i.Stacks = i.Stacks[scopeLevel:]
}

func (i *Interpreter) execStatement(s *ast.Statement) {
	switch s.LHS {
	case "BRK":
		fmt.Printf("\n")
		return
	case "PRI":
		fmt.Printf("%s", i.Stacks.resolveValue(s.RHS))
		return
	case "PRU":
		cast, _ := strconv.Atoi(i.Stacks.resolveValue(s.RHS))
		fmt.Printf("%c", cast)
		return
	case "INC":
		if s.RHS.GetNodeName() != "statVar" {
			errors.LitIncrementError()
			os.Exit(65)
		}

		variableValue := i.Stacks.resolveVariable(s.RHS)
		cast, _ := strconv.Atoi(variableValue.value)
		cast++
		castBack := strconv.Itoa(cast)
		i.Stacks.set(variableValue.key, castBack)
		return
	}
}

func (s Stacks) set(name string, value string) {
	scopeLevel, index := s.find(name)
	s[scopeLevel][index].value = value
}

func (s Stacks) find(key string) (int, int) {
	//reverse stack search so we start at local scope and keep working our way up intill we find something

	for i := len(s) - 1; i > -1; i-- {
		stackIndex := s[i].stackContains(key)
		if stackIndex != -1 {
			//scopeLevel, scopedStackIndex
			return i, stackIndex
		}
	}

	errors.UndefinedVariableError(key)
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

func (h Heap) find(name string) int {
	for i := 0; i < len(h); i++ {
		if h[i].Name == name {
			return i
		}
	}
	errors.UndefinedFunctionReferenceError(name)
	os.Exit(65)
	return -1
}

func (h Heap) resolveFunction(name string) *ast.Function {
	return h[h.find(name)]
}

func (s *Stacks) get(scopeLevel int, index int) *tuple {
	return (*s)[scopeLevel][index]
}

func (s *Stacks) resolveVariable(node ast.Node) *tuple {
	return s.get(s.find(node.(*ast.StatVar).Value))
}

func (s *Stacks) resolvePtrValue(node ast.Node) *string {
	if ast.NodeIsVariable(node) {
		return &s.resolveVariable(node).value
	}
	return &node.(*ast.NumLit).Value
}

func (s *Stacks) resolveValue(node ast.Node) string {
	if ast.NodeIsVariable(node) {
		return s.resolveVariable(node).value
	}
	return node.(*ast.NumLit).Value

}
