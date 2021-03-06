package main

func (fc *FileCompiler) handleExport() {

	fc.expectToken("{")

	token := fc.getNextToken(false, false)
	for token != "}" {
		if token == "" {
			fc.throwAtLine("Unexpected end of code")
		}
		fc.checkVarNameSyntax([]byte(token))
		varName := token

		fc.expectToken(":")
		token = fc.getNextToken(false, false)
		fc.checkVarNameSyntax([]byte(token))
		valueName := token

		_, isVar := fc.getVar(valueName)
		t, isType := fc.getType(valueName)
		if !isVar && !isType {
			fc.throwAtLine("Unknown var/type: " + valueName)
		}
		if isType {
			valueName = t.name
		}

		fc.exports[varName] = valueName

		token = fc.getNextToken(false, false)
		if token == "}" {
			break
		}
		if token != "," {
			fc.throwAtLine("Unexpected token: " + token)
		}

		token = fc.getNextToken(false, false)
	}

	token = fc.getNextToken(true, false)
	if token == ";" {
		token = fc.getNextToken(false, false)
	}

}
