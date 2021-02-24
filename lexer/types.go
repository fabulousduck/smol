package lexer

import (
	"strings"
)

type typename map[string]string

//OperatorAttributes wraps the precedance and Associativity of a operator
type OperatorAttributes struct {
	Precedance    int
	Associativity string
}

var operatorAttributeMap = map[string]OperatorAttributes{
	"left_parenthesis":  {11, "left"}, // (
	"right_parenthesis": {11, "left"}, // )
	"less_than":         {6, "left"},  // <
	"greater_than":      {6, "left"},  // >
	"exponent":          {4, "right"}, // ^
	"division":          {3, "left"},  // /
	"star":              {3, "left"},  // *
	"plus":              {2, "left"},  // +
	"dash":              {2, "left"},  // -
}

var types = map[string][]string{
	"character":         []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "_"},
	"integer":           []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
	"less_than":         []string{"<"},
	"comma":             []string{","},
	"greater_than":      []string{">"},
	"exponent":          []string{"^"},
	"left_parenthesis":  []string{"("},
	"right_parenthesis": []string{")"},
	"semicolon":         []string{";"},
	"plus":              []string{"+"},
	"double_quote":      []string{"\""},
	"star":              []string{"*"},
	"division":          []string{"/"},
	"equals":            []string{"="},
	"dash":              []string{"-"},
	"left_bracket":      []string{"["},
	"right_bracket":     []string{"]"},
	"double_dot":        []string{":"},
	"comment":           []string{"#"},
	"newline":           []string{"\r", "\n"},
	"ignoreable":        []string{"\t", " "},
}

//IsLitteral checks if a given token is a litteral type
func IsLitteral(token Token) bool {
	litteralTypes := []string{"character", "string", "integer", "string_litteral"}

	for _, litteral := range litteralTypes {
		if token.Type == litteral {
			return true
		}
	}
	return false
}

func determineType(character string) string {

	usableChar := strings.ToLower(character)

	for key, values := range types {
		if contains(usableChar, values) {
			return key
		}
	}

	return "undefined_symbol"
}

func contains(name string, list []string) bool {
	for i := 0; i < len(list); i++ {
		if string(list[i]) == name {
			return true
		}
	}
	return false
}

func containsEither(names []string, list []string) bool {
	for i := 0; i < len(names); i++ {
		if contains(names[i], list) {
			return true
		}
	}

	return false
}

func getKeyword(token *Token) string {
	keywords := map[string][]string{
		"function_definition": []string{"def"},
		"boolean_keyword":     []string{"True", "False"},
		"variable_type":       []string{"String", "Bool", "Uint32", "Uint64"},
		"print":               []string{"print"},
		"include":             []string{"include"},
		"close_block":         []string{"end"},
		"set_variable":        []string{"set"},
		"if_statement":        []string{"if"},
		"switch":              []string{"switch"},
		"case":                []string{"case"},
		"end_of_switch":       []string{"default"},
		"free":                []string{"free"},
		"plot":                []string{"plot"},
	}

	for key, values := range keywords {
		if contains(token.Value, values) {
			return key
		}
	}

	if len(token.Value) > 1 {
		return "string"
	}
	return token.Type
}

//GetOperatorAttributes returns the precedance and associativity of an operator
func GetOperatorAttributes(operator string) OperatorAttributes {
	if val, ok := operatorAttributeMap[operator]; ok {
		return val
	}

	//we assume its a variable if its not found
	//this works since functions and variables have the same associativity and precedance
	return OperatorAttributes{14, "left"}
}

//IsOperator checks if a given char is a valid operator
func IsOperator(operator string) bool {
	if _, ok := operatorAttributeMap[operator]; ok {
		return true
	}
	return false
}

//GetPrec is a simple wrapper for GetOperatorAttributes but only returns the precedance
func GetPrec(operator string) int {
	return GetOperatorAttributes(operator).Precedance
}

//HasHigherPrec is a simple check for checking if lhs has a bigger precednce than rhs
func (lhs Token) HasHigherPrec(rhs Token) bool {
	operatorA := GetOperatorAttributes(lhs.Type)
	operatorB := GetOperatorAttributes(rhs.Type)
	hasPrec := (operatorB.Associativity == "left" && operatorB.Precedance <= operatorA.Precedance) ||
		(operatorB.Associativity == "right" && operatorB.Precedance < operatorA.Precedance)

	return hasPrec
}

//DetermineStringType will determine the type of a given string
func DetermineStringType(str string) string {
	return determineType(string([]rune(str)[0]))
}
