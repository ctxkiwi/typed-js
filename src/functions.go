package main

func isAlphaChar(c byte) bool {

	if c >= 65 && c <= 90 {
		return true
	}

	if c >= 97 && c <= 122 {
		return true
	}

	return false
}

func isVarNameChar(c byte) bool {

	if c >= 65 && c <= 90 {
		return true
	}

	if c >= 97 && c <= 122 {
		return true
	}

	// Underscore
	if c == 95 {
		return true
	}

	return false
}

func isVarNameSyntax(name []byte) bool {
	for _, char := range name {
		if !isVarNameChar(char) {
			return false
		}
	}
	return true
}

func isNumberChar(c byte) bool {

	if c >= 48 && c <= 57 {
		return true
	}

	return false
}

func isWhiteSpace(c byte) bool {
	return c <= 32
}

func isNewLine(c byte) bool {
	return c == 10
}
