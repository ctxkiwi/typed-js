package main

func (fc *FileCompiler) handleType(isLocal bool, isDefine bool, isStruct bool, isClass bool) {

	name := fc.getNextToken(false, true)

	_, ok := fc.getVar(name)
	if ok {
		fc.throwAtLine("Name already used as a variable: " + name)
	}

	_type, ok := fc.getType(name)
	if !ok {
		fc.throwAtLine("Compiler bug, type not found: " + name)
	}

	token := fc.getNextToken(false, true)
	if isDefine && token == "," {
		fc.expectToken("class")
		fc.getNextToken(false, true)
		fc.getNextTokenSameLine()
	}

	// token = {
	// Get fields
	functionCode := map[string]string{}
	token = fc.getNextToken(false, false)
	for token != "}" {

		varName := token
		fc.expectToken(":")

		// Read type
		token = fc.getNextToken(true, true)

		if isClass && !isDefine && token == "function" {
			fc.recordResult()
			fc.assignValue()
			functionCode[varName] = fc.getRecording()
		} else {
			fc.getNextType()
		}

		// Check default
		char := fc.getNextCharacterOnLine()
		if char == "=" {
			value, _ := fc.getNextValueToken()
			_type.props[varName]._default = value
		}

		token = fc.getNextToken(false, false)
	}

	// Write result code
	if isClass && !isDefine {
		globalName := _type.name
		fc.addSpace()
		fc.addResult("var " + globalName + " = function(")
		constructorProp, hasConstructor := _type.props["constructor"]
		if hasConstructor {
			for i, vtype := range constructorProp.varType.paramTypes {
				if i > 0 {
					fc.addResult(", ")
				}
				fc.addResult(vtype.paramName)
			}
		}
		if hasConstructor {
			fc.addResult(") {\n")
			extraSpace++
			fc.addSpace()
			fc.addResult("this.constructor = " + functionCode["constructor"] + ";\n")
			fc.addSpace()
			fc.addResult("this.constructor(")
			for i, vtype := range constructorProp.varType.paramTypes {
				if i > 0 {
					fc.addResult(", ")
				}
				fc.addResult(vtype.paramName)
			}
			fc.addResult(");\n")
			extraSpace--
			fc.addSpace()
		} else {
			fc.addResult(") {")
		}
		fc.addResult("};\n")
		fc.addSpace()

		for funcName, fcode := range functionCode {
			if funcName == "constructor" {
				continue
			}
			fc.addResult(globalName + ".prototype." + funcName + " = " + fcode + ";\n")
			fc.addSpace()
		}
	}

}
