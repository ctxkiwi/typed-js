
package main;

func isAlphaChar (c byte) bool {

	if c >= 65 && c <= 90 {
		return true
	}

	if c >= 97 && c <= 122 {
		return true
	}

	return false
}

func isWhiteSpace (c byte) bool {
	return c <= 32
}

func isNewLine (c byte) bool {
	return c == 10
}
