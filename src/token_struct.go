
package main

import (
)

func (c *Compile) handleStruct (isLocal bool) {

	if isLocal {

		token := c.getNextToken();
		if token != "struct" {
			c.throwAtLine("Unexpected token: " + token)
		}
	} else {
		if c.scopeIndex > 0 {
			c.throwAtLine("You must use 'local' to create a local struct")
		}
	}

	name := c.getNextToken()
	c.checkVarNameSyntax([]byte(name))

	if(c.typeExists(name)){
		c.throwAtLine("Struct name already in use: " + name)
	}

	globalName, s := createNewStruct()
	s.isLocal = isLocal

	scope := c.scopes[c.scopeIndex]
	scope.structs[name] = globalName

	// Get fields
	c.expectToken("{")

	token := c.getNextToken()
	for token != "}" {
		if token == "" {
			c.throwAtLine("Unexpected end of code")
		}
		c.checkVarNameSyntax([]byte(token))
		varName := token
		_, ok := s.props[varName]; 
		if ok {
			c.throwAtLine("Property name '"+ varName +"' already exists")
		}
		c.expectToken(":")
		// Read type
		t := c.getNextType()
		_typeOfType, ok := c.getTypeOfType(t)
		prop := Property{
			_type: t,
			_typeOfType: _typeOfType,
		}
		// Check default
		char := c.getNextCharacterOnLine()
		if char == "=" {
			value,_ := c.getNextValueToken()
			prop._default = value
		}

		// Store property
		s.props[varName] = prop

		token = c.getNextToken()
	}

}
