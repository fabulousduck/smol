package lexer

import (
	"strings"
)

type typename map[string]string

func determineType(character string) string {

	usableChar := strings.ToLower(character)
	types := map[string][]string{
		"character":     []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "_"},
		"integer":       []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
		"left_arrow":    []string{"<"},
		"comma":         []string{","},
		"right_arrow":   []string{">"},
		"semicolon":     []string{";"},
		"left_bracket":  []string{"["},
		"right_bracket": []string{"]"},
		"double_dot":    []string{":"},
		"comment":       []string{"#"},
		"newline":       []string{"\r", "\n"},
		"tab":           []string{"\t"},
		"space":         []string{" "},
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

func getKeyword(token *Token) string {
	keywords := map[string]string{
		"DEF": "function_definition",
		"ANB": "left_not_right",
		"MEM": "variable_assignment",
		"PRI": "print_integer",
		"PRU": "print_ascii",
		"INC": "increment_value",
		"END": "close_block",
		"BRK": "print_break",
		"SET": "set_variable",
		"ADD": "addition",
		"SUB": "subtraction",
		"MUL": "multiplication",
		"DIV": "division",
		"POW": "power_of",
		"EQ":  "equals",
		"NEQ": "not_equals",
		"LT":  "less_than",
		"GT":  "greater_than",
		"SWT": "switch",
		"CAS": "case",
		"EOS": "end_of_switch",
	}

	if val, ok := keywords[token.Value]; ok {
		return val
	}

	if len(token.Value) > 1 {
		return "string"
	}
	return token.Type
}

//DetermineStringType will determine the type of a given string
func DetermineStringType(str string) string {
	return determineType(string([]rune(str)[0]))
}
