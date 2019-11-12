package lexer

import (
	"os"
	"strings"

	"github.com/fabulousduck/smol/errors"
)

type typename map[string]string

type operatorAttributes struct {
	precedance    int
	associativity string
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
	types := map[string][]string{
		"character":         []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "_"},
		"integer":           []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
		"less_than":         []string{"<"},
		"comma":             []string{","},
		"greater_than":      []string{">"},
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
		"while_not":           []string{"whileNot"},
		"boolean_keyword":     []string{"True", "False"},
		"variable_type":       []string{"String", "Bool", "Uint32", "Uint64"},
		"print":               []string{"print"},
		"close_block":         []string{"end"},
		"set_variable":        []string{"set"},
		"equals":              []string{"eq"},
		"not_equals":          []string{"neq"},
		"less_than":           []string{"lt"},
		"greater_than":        []string{"gt"},
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

//returns the precedance and associativity of an operator
func getOperatorAttributes(operator string) operatorAttributes {
	operatorAttributeMap := map[string]operatorAttributes{
		"less_than":         {6, "left"},  // <
		"greater_than":      {6, "left"},  // >
		"exponent":          {4, "right"}, // ^
		"slash":             {3, "left"},  // /
		"star":              {3, "left"},  // *
		"plus":              {2, "left"},  // +
		"dash":              {2, "left"},  // -
		"left_parenthesis":  {1, "left"},  // (
		"right_parenthesis": {1, "left"},  // )
	}

	if val, ok := operatorAttributeMap[operator]; ok {
		return val
	}
	errors.InvalidOperatorError()
	os.Exit(65)
	//this is only here because otherwise golang will cry about no return value even though this is unreachable code
	return operatorAttributes{}
}

//DetermineStringType will determine the type of a given string
func DetermineStringType(str string) string {
	return determineType(string([]rune(str)[0]))
}
