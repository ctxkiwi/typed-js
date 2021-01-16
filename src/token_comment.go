package main

func (c *Compile) handleComment() {

	token := string(c.code[c.index-1])

	if token == "/" {
		for c.index <= c.maxIndex {

			charInt := c.code[c.index]
			// char := string(charInt)
			c.index++

			if isNewLine(charInt) {
				return
			}
		}
		return
	}
	if token == "*" {
		for c.index <= c.maxIndex {

			charInt := c.code[c.index]
			char := string(charInt)
			c.index++
			c.col++

			if isNewLine(charInt) {
				c.line++
				c.col = 0
				c.lastTokenCol = 0
			}

			if char == "*" {
				nextToken := c.readNextChar()
				if nextToken == "/" {
					c.index++
					c.col++
					return
				}
			}
		}
		return
	}

	c.throwAtLine("Unknown token: " + token)
}
