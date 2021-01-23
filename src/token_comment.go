package main

func (fc *FileCompiler) handleComment() {

	token := string(fc.code[fc.index-1])

	if token == "/" {
		for fc.index <= fc.maxIndex {

			charInt := fc.code[fc.index]
			// char := string(charInt)
			fc.index++

			if isNewLine(charInt) {
				return
			}
		}
		return
	}
	if token == "*" {
		for fc.index <= fc.maxIndex {

			charInt := fc.code[fc.index]
			char := string(charInt)
			fc.index++
			fc.col++

			if isNewLine(charInt) {
				fc.line++
				fc.col = 0
				fc.lastTokenCol = 0
			}

			if char == "*" {
				nextToken := fc.readNextChar()
				if nextToken == "/" {
					fc.index++
					fc.col++
					return
				}
			}
		}
		return
	}

	fc.throwAtLine("Unknown token: " + token)
}
