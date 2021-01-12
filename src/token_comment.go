package main

func (c *Compile) handleComment() {

	token := c.getNextToken()

	if token == "/" {

		return
	}
	if token == "*" {
		return
	}

	c.throwAtLine("Unknown token: " + token)
}
