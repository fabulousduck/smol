package lexer

type typename map[string]string

func determineType(character string) string {
	chars := []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
		"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "_",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
		"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	}

	symbols := typename{
		"<": "LEFT_ARROW",
		",": "COMMA",
		">": "RIGHT_ARROW",
		";": "SEMI_COLON",
		"[": "LEFT_BRACKET",
		"]": "RIGHT_BRACKET",
		":": "DOUBLE_DOT",
	}

	if character == "#" {
		return "COMMENT"
	}

	escapeChars := typename{
		"\r": "WIN_NEWLINE",
		"\n": "UNIX_NEWLINE",
		"\t": "TAB",
	}

	numbers := []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	}

	if contains(character, chars) {
		return "CHAR"
	}

	if val, ok := symbols[character]; ok {
		return val
	}

	if val, ok := escapeChars[character]; ok {
		return val
	}

	if contains(character, numbers) {
		return "NUMB"
	}

	if character == " " {
		return "SPACE"
	}

	return "UDEF"
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
