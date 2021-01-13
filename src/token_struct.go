package main

func (c *Compile) handleStruct(isLocal bool, isDefine bool) {

	if !isLocal {
		if scopeIndex > 0 {
			c.throwAtLine("You must use 'local' to create a local struct")
		}
	}

	name := c.getNextTokenSameLine()
	c.checkVarNameSyntax([]byte(name))

	if c.typeExists(name) {
		c.throwAtLine("Struct/class name already in use: " + name)
	}

	globalName, s := createNewStruct()
	s.isLocal = isLocal

	scope := scopes[scopeIndex]
	scope.structs[name] = globalName

	token := c.getNextTokenSameLine()
	var class *Class = nil
	if isDefine && token == "," {
		c.expectToken("class")
		className := c.getNextTokenSameLine()
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

	token = c.getNextToken()
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
		_typeOfType, ok := c.getTypeOfType(t)
		prop := Property{
			_type:       t,
			_typeOfType: _typeOfType,
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

		token = c.getNextToken()
	}

}
