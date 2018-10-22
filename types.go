package losp

type typename map[string]string

func determineType(character string) string {
	chars := []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
		"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}

	symbols := typename{
		"<": "LEFT_ARROW",
		",": "COMMA",
		">": "RIGHT_ARROW",
		";": "SEMI_COLON",
		"[": "LEFT_BRACE",
		"]": "RIGHT_BRACE",
	}

	escapeChars := typename{
		"\n": "NEWLINE",
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
