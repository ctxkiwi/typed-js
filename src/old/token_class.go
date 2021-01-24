package main

func (fc *FileCompiler) handleClass(isDefine bool) {

	name := fc.getNextToken(false, true)
	fc.checkVarNameSyntax([]byte(name))

	if fc.typeExists(name) {
		fc.throwAtLine("Struct/class name already in use: " + name)
	}
	_, ok := fc.getVar(name)
	if ok {
		fc.throwAtLine("Name already used as a variable: " + name)
	}

	globalName, class := createNewClass()

	scope := fc.compiler.scopes[fc.compiler.scopeIndex]
	scope.classes[name] = globalName

	token := fc.getNextToken(false, true)

	// Get fields
	if token != "{" {
		fc.throwAtLine("Unexpected token: " + token)
	}

	functionCode := map[string]string{}

	token = fc.getNextToken(false, false)
	for token != "}" {
		if token == "" {
			fc.throwAtLine("Unexpected end of code")
		}
		fc.checkVarNameSyntax([]byte(token))
		varName := token
		_, ok := class.props[varName]
		if ok {
			fc.throwAtLine("Property name '" + varName + "' already exists")
		}
		fc.expectToken(":")

		// Read type or function
		prop := Property{}
		token = fc.getNextToken(true, true)

		if varName == "constructor" && token != "function" {
			fc.throwAtLine("Constructor must be a function")
		}

		if !isDefine && token == "function" {
			fc.recordResult()
			t := fc.assignValue()
			_typeOfType, _ := fc.getTypeOfType(t.name)
			t.toft = _typeOfType
			prop.varType = t
			functionCode[varName] = fc.getRecording()
		} else {
			t := fc.getNextType()
			_typeOfType, _ := fc.getTypeOfType(t.name)
			t.toft = _typeOfType
			prop.varType = t
		}
		// Check default
		char := fc.getNextCharacterOnLine()
		if char == "=" {
			value, _ := fc.getNextValueToken()
			prop._default = value
		}

		// Store property
		class.props[varName] = &prop

		token = fc.getNextToken(false, false)
	}

	// Write result code
	if !isDefine {
		fc.addSpace()
		fc.addResult("var " + globalName + " = function(")
		constructorProp, hasConstructor := class.props["constructor"]
		if hasConstructor {
			for i, vtype := range constructorProp.varType.paramTypes {
				if i > 0 {
					fc.addResult(", ")
				}
				fc.addResult(vtype.paramName)
			}
		}
		fc.addResult(") {\n")
		extraSpace++
		fc.addSpace()
		if hasConstructor {
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
			extraSpace--
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
