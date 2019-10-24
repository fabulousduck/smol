package lexer

import (
	"strings"
)

type typename map[string]string

func determineType(character string) string {

	usableChar := strings.ToLower(character)
	types := map[string][]string{
		"character":         []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "_"},
		"integer":           []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"},
		"left_arrow":        []string{"<"},
		"comma":             []string{","},
		"right_arrow":       []string{">"},
		"left_parenthesis":  []string{"("},
		"right_parenthesis": []string{")"},
		"semicolon":         []string{";"},
		"plus":              []string{"+"},
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

//DetermineStringType will determine the type of a given string
func DetermineStringType(str string) string {
	return determineType(string([]rune(str)[0]))
}
