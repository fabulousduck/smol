package interpreter

import (
	"fmt"
	"os"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/fabulousduck/smol/ast"
	"github.com/fabulousduck/smol/errors"
)

type stackItem struct {
	key, value string
}

type heapItem struct {
	key      string
	itemType string
	value    ast.Node
}

type stack []*stackItem

//Stacks is the global scope that hold sub scopes for varianbles
type Stacks []stack

//Heap is not really a heap since it does not hold dynamically sized types, but a good excuse to put my function decls into
type heap []*heapItem

//Heaps is a list of scopes starting from global to nested
type Heaps []heap

//Interpreter contains all data needed to Interpret an AST
type Interpreter struct {
	Stacks       Stacks
	Heaps        Heaps
	CurrentScope int
}

//NewInterpreter provides a new interpreter with empty base stack and heap
func NewInterpreter() *Interpreter {
	i := new(Interpreter)

	//scopes work on an array, where 0 is the global scope
	//all following scopes are
	i.CurrentScope = 0
	i.Stacks = Stacks{stack{}}
	i.Heaps = Heaps{heap{}}
	return i
}

//Interpret will tree walk execute an AST from left to right (topdown)
func (i Interpreter) Interpret(AST []ast.Node) {
	for j := 0; j < len(AST); j++ {
		node := AST[j]
		nodeType := node.GetNodeName()
		switch nodeType {
		case "variable":

			if node.(*ast.Variable).GetAllocationType() == "stack" {
				i.stackAlloc(len(i.Stacks)-1, node.(*ast.Variable))
				break
			}

			i.heapAlloc(len(i.Heaps)-1, node)
			break

		case "function":
			i.heapAlloc(0, node)
		case "printCall":
			i.execPrintCall(node.(*ast.PrintCall))
		case "functionCall":
			i.execFunctionCall(node.(*ast.FunctionCall))
		case "setStatement":
			i.setVariableValue(node.(*ast.SetStatement))
		case "switchStatement":
			i.execSwitchStatement(node.(*ast.SwitchStatement))
		case "freeStatement":
			i.execFreeStatement(node.(*ast.FreeStatement))
		case "directOperation":
			i.execDirectOperation(node.(*ast.DirectOperation))
		}
	}
}

func (i *Interpreter) execDirectOperation(do *ast.DirectOperation) {

	variableValue := i.Stacks.resolveVariable(do.Variable)
	cast, _ := strconv.Atoi(variableValue.value)
	if do.Operation == "++" {
		cast++
	} else {
		cast--
	}
	castBack := strconv.Itoa(cast)
	i.Stacks.set(variableValue.key, castBack)
}

func (i *Interpreter) execPrintCall(pc *ast.PrintCall) {
	printable := i.Stacks.resolveValue(pc.Printable)
	fmt.Printf(printable)
}

func (i *Interpreter) execFreeStatement(r *ast.FreeStatement) {
	if !ast.NodeIsVariable(r.Variable) {
		errors.LitteralFree()
		os.Exit(65)
	}

	//find the variable
	scopeLevel, stackIndex := i.Stacks.find(r.Variable.(*ast.StatVar).Value)

	//remove the variable
	i.Stacks[scopeLevel] = append(i.Stacks[scopeLevel][:stackIndex], i.Stacks[scopeLevel][stackIndex+1:]...)
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
	//indexing this to 0 means we only look at top level scope.
	//this disallows functions in functions since we wont find them this way
	functionDecl := i.Heaps[0].resolveFunction(fc.Name)

	if len(fc.Args) != len(functionDecl.Params) {
		errors.IncorrectFunctionParamCountError(functionDecl.Name, len(fc.Args), len(functionDecl.Params))
		os.Exit(65)
		return
	}

	beforeScopeLevel := len(i.Stacks)
	scopedStack := stack{}
	scopedHeap := heap{}

	for paramListIndex, param := range functionDecl.Params {
		currentArg := fc.Args[paramListIndex]

		if ast.NodeIsVariable(currentArg) {
			if currentArg.(*ast.Variable).GetAllocationType() == "stack" {
				scopedStack = append(scopedStack, i.Stacks.resolveVariable(currentArg))
				continue
			}
			scopedHeap = append(scopedHeap, i.Heaps.resolveVariable(currentArg))
		} else {
			nodeName := currentArg.GetNodeName()
			switch nodeName {
			case "BoolLit":
				boolLit := currentArg.(*ast.BoolLit)
				scopedStack = append(scopedStack, &stackItem{key: param, value: boolLit.Value})
				break
			case "StringLit":
				scopedHeap = append(scopedHeap, &heapItem{key: param, itemType: "stringLit", value: currentArg})
				break
			case "NumLit":
				numLit := currentArg.(*ast.NumLit)
				scopedStack = append(scopedStack, &stackItem{key: param, value: numLit.Value})
				break
			default:
				//TODO: error: attempted stack allocation for unknown litteral
			}
		}
	}

	//left off here refactoring

	i.Stacks = append(i.Stacks, scopedStack)
	i.Interpret(functionDecl.Body)
	i.Stacks = i.Stacks[:beforeScopeLevel]
}
func (i *Interpreter) stackAlloc(scopeLevel int, v *ast.Variable) {
	stackItem := new(stackItem)
	spew.Dump(v)
	stackItem.key = v.Name

	stackItem.value = i.Stacks.resolveValue(v.Value)

	i.Stacks[scopeLevel] = append(i.Stacks[scopeLevel], stackItem)
}

func (i *Interpreter) heapAlloc(scopeLevel int, node ast.Node) {
	heapItem := new(heapItem)
	nodeType := node.GetNodeName()

	switch nodeType {
	case "Function":
		heapItem.key = node.(*ast.Function).Name
		break
	case "StringLit":
		heapItem.key = node.(*ast.Variable).Name
		break
	default:
		//TODO: error message: invalid heap allocation
	}

	heapItem.value = node
	heapItem.itemType = node.GetNodeName()
	i.Heaps[scopeLevel] = append(i.Heaps[scopeLevel], heapItem)
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

func (s *Stacks) get(scopeLevel int, index int) *stackItem {
	return (*s)[scopeLevel][index]
}

func (s *Stacks) resolveVariable(node ast.Node) *stackItem {
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

	switch node.GetNodeName() {
	case "numLit":
		return node.(*ast.NumLit).Value
	case "boolLit":
		return node.(*ast.BoolLit).Value
	case "stringLit":
		return node.(*ast.StringLit).Value
	}

	errors.UnresolvableVariableValueError()
	os.Exit(65)
	//something janky to keep go happy
	//this is literally unreachable code but it doesnt detect it
	return node.(*ast.NumLit).Value
}

func (h *Heaps) resolveVariable(node ast.Node) *heapItem {
	return h.get(h.find(node.(*ast.StatVar).Value))
}

func (h *Heaps) get(scopeLevel int, index int) *heapItem {
	return (*h)[scopeLevel][index]
}

func (h heap) heapContains(name string) int {
	for i := 0; i < len(h); i++ {
		//There are two types on the heap. String and function.
		//If we add more than that, this function should be made more generic
		if h[i].itemType == "Function" {
			if h[i].value.(*ast.Function).Name == name {
				return i
			}
		} else if h[i].itemType == "StringLit" {
			if h[i].value.(*ast.Variable).Name == name {
				return i
			}
		} else {
			continue
		}
	}
	errors.UndefinedFunctionOrStringReferenceError(name)
	os.Exit(65)
	return -1
}

func (h heap) resolveFunction(name string) *ast.Function {
	return h[h.heapContains(name)].value.(*ast.Function)
}

func (h Heaps) find(key string) (int, int) {
	//reverse stack search so we start at local scope and keep working our way up intill we find something

	for i := len(h) - 1; i > -1; i-- {
		heapIndex := h[i].heapContains(key)
		if heapIndex != -1 {
			//scopeLevel, scopedStackIndex
			return i, heapIndex
		}
	}

	errors.UndefinedVariableError(key)
	os.Exit(65)
	return -1, -1
}
