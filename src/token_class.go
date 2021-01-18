package main

func (c *Compile) handleClass(isDefine bool) {

	name := c.getNextToken(false, true)
	c.checkVarNameSyntax([]byte(name))

	if c.typeExists(name) {
		c.throwAtLine("Struct/class name already in use: " + name)
	}
	_, ok := c.getVar(name)
	if ok {
		c.throwAtLine("Name already used as a variable: " + name)
	}

	globalName, class := createNewClass()

	scope := scopes[scopeIndex]
	scope.classes[name] = globalName

	token := c.getNextToken(false, true)

	// Get fields
	if token != "{" {
		c.throwAtLine("Unexpected token: " + token)
	}

	token = c.getNextToken(false, false)
	for token != "}" {
		if token == "" {
			c.throwAtLine("Unexpected end of code")
		}
		c.checkVarNameSyntax([]byte(token))
		varName := token
		_, ok := class.props[varName]
		if ok {
			c.throwAtLine("Property name '" + varName + "' already exists")
		}
		c.expectToken(":")

		// Read type or function
		prop := Property{}
		token = c.getNextToken(true, true)
		if !isDefine && token == "function" {
			t := c.assignValue()
			_typeOfType, _ := c.getTypeOfType(t.name)
			t.toft = _typeOfType
			prop.varType = t
		} else {
			t := c.getNextType()
			_typeOfType, _ := c.getTypeOfType(t.name)
			t.toft = _typeOfType
			prop.varType = t
		}
		// Check default
		char := c.getNextCharacterOnLine()
		if char == "=" {
			value, _ := c.getNextValueToken()
			prop._default = value
		}

		// Store property
		class.props[varName] = &prop

		token = c.getNextToken(false, false)
	}

}
