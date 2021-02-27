package main

func (fc *FileCompiler) handleIf() {

	if !fc.compiler.readTypes {
		fc.addResult("if ")
	}
	token := "else"
	nextToken := "if"

	for token == "else" {

		if nextToken == "if" {
			fc.expectToken("(")
			if !fc.compiler.readTypes {
				fc.addResult("(")
			}

			if fc.compiler.readTypes {
				fc.skipValue()
			} else {
				fc.assignValue()
			}

			fc.expectToken(")")
			if !fc.compiler.readTypes {
				fc.addResult(")")
			}
		}

		fc.expectToken("{")
		fc.createNewScope()
		// extraSpace++
		if !fc.compiler.readTypes {
			fc.addResult("{\n")
			fc.addSpace()
		}

		//////////////////

		fc.compile()

		//////////////////

		// extraSpace--
		fc.popScope()
		if !fc.compiler.readTypes {
			fc.addResult("\n")
			fc.addSpace()
			fc.addResult("}")
		}

		if nextToken != "if" {
			break
		}

		token = fc.getNextToken(true, false)
		if token == "else" {
			token = fc.getNextToken(false, false)
			if !fc.compiler.readTypes {
				fc.addResult(" else ")
			}
			nextToken = fc.getNextToken(true, true)
			if nextToken == "if" {
				nextToken = fc.getNextToken(false, true)
				if !fc.compiler.readTypes {
					fc.addResult("if ")
				}
			}
		}
	}

}
