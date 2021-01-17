package main

func (c *Compile) handleStruct(isLocal bool, isDefine bool) {

	if !isLocal {
		if scopeIndex > 0 {
			c.throwAtLine("You must use 'local' to create a local struct")
		}
	}

	name := c.getNextToken(false, true)
	c.checkVarNameSyntax([]byte(name))

	if c.typeExists(name) {
		c.throwAtLine("Struct/class name already in use: " + name)
	}
	_, ok := c.getVar(name)
	if ok {
		c.throwAtLine("Name already used as a variable: " + name)
	}

	globalName, s := createNewStruct()
	s.isLocal = isLocal

	scope := scopes[scopeIndex]
	scope.structs[name] = globalName

	token := c.getNextToken(false, true)
	var class *Class = nil
	if isDefine && token == "," {
		c.expectToken("class")
		className := c.getNextToken(false, true)
		c.checkVarNameSyntax([]byte(className))
		if c.typeExists(className) {
			c.throwAtLine("Struct/class name already in use: " + className)
		}
		globalName, class = createNewClass()
		scope.classes[className] = globalName
		token = c.getNextTokenSameLine()
	}

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
		_, ok := s.props[varName]
		if ok {
			c.throwAtLine("Property name '" + varName + "' already exists")
		}
		c.expectToken(":")
		// Read type
		t := c.getNextType()
		_typeOfType, ok := c.getTypeOfType(t.name)
		t.toft = _typeOfType
		prop := Property{
			varType: t,
		}
		// Check default
		char := c.getNextCharacterOnLine()
		if char == "=" {
			value, _ := c.getNextValueToken()
			prop._default = value
		}

		// Store property
		s.props[varName] = &prop
		if class != nil {
			class.props[varName] = &prop
		}

		token = c.getNextToken(false, false)
	}

}
