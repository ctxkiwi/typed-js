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

			if char == "*" {
				nextToken := c.getNextToken(false, false)
				if nextToken == "/" {
					return
				}
			}
		}
		return
	}

	c.throwAtLine("Unknown token: " + token)
}
