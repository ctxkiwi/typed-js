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
	for i, char := range name {
		if i > 0 && isNumberChar(char) {
			continue
		}
		if !isVarNameChar(char) {
			return false
		}
	}
	return true
}

func (fc *FileCompiler) checkVarNameSyntax(name []byte) {
	if !isVarNameSyntax(name) {
	fc.throwAtLine("Invalid variable name: " + string(name))
	}
}

func isNumberSyntax(name []byte) bool {
	hasDot := false
	for i, char := range name {
		if i > 0 && string(char) == "." {
			if hasDot {
				return false
			}
			hasDot = true
			continue
		}
		if !isNumberChar(char) {
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
