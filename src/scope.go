package main

type Var struct {
	_typeOfType string // basic,struct,class
	_type       string // string bool int MyStruct MyClass etc...
	nullable    bool
}

type Scope struct {
	structs map[string]string
	classes map[string]string
	vars    map[string]Var
}

func (s *Scope) typeExists(name string) bool {
	_, ok := s.structs[name]
	if ok {
		return true
	}
	_, ok = s.classes[name]
	if ok {
		return true
	}
	return false
}

func (s *Scope) hasStruct(name string) bool {
	for _, str := range s.structs {
		if str == name {
			return true
		}
	}
	return false
}

func (s *Scope) hasClass(name string) bool {
	for _, str := range s.classes {
		if str == name {
			return true
		}
	}
	return false
}

func (s *Scope) getVar(name string) (*Var, bool) {
	result, ok := s.vars[name]
	return &result, ok
}

func (c *Compile) declareVariable(_type string, _typeOfType string) {
	varName := c.getNextTokenSameLine()
	c.checkVarNameSyntax([]byte(varName))
	_, ok := c.getVar(varName)
	if ok {
		c.throwAtLine("Variable name already in use: " + varName)
	}
	c.result += "var " + varName
	c.expectToken("=")
	c.result += " = "
	c.assignValue(_type, _typeOfType)
	// value, _ := c.getNextValueToken()
	// c.result += value
}

func (c *Compile) assignValue(_type string, _typeOfType string) {

	token := c.getNextTokenSameLine()
	if token == "(" {
		c.result += "("
		c.assignValue(_type, _typeOfType)
		c.expectToken(")")
		c.result += ")"
		return
	}
	if isVarNameSyntax([]byte(token)) {
		_, ok := c.getVar(token)
		if ok {
			// Is variable name
			c.result += token
			// ...
		}
	}

	if _typeOfType == "basic" {

		if _type == "string" {
			if token != "\"" && token != "'" {
				c.throwAtLine("Unexpected token: " + token)
			}
			char := ""
			lastChar := ""
			for c.index <= c.maxIndex {
				lastChar = char
				charInt := c.code[c.index]
				char = string(charInt)
				c.index++
				c.result += char
				if isNewLine(charInt) {
					if lastChar != "\\" {
						c.throwAtLine("Unexpected newline")
					}
					c.line++
					c.col = 0
					c.lastTokenCol = 0
				}
				if char == token && lastChar != "\\" {
					break
				}
			}
			if c.index > c.maxIndex {
				c.throwAtLine("You forgot to close a string somewhere")
			}
		}

	} else if _typeOfType == "struct" {
		if token != "{" {
			c.throwAtLine("Expected token: {")
		}
		c.result += "{\n"

		s, _ := c.getStruct(_type)
		token := c.getNextToken()
		for token != "}" {
			if token == "" {
				c.throwAtLine("Unexpected end of code")
			}
			c.checkVarNameSyntax([]byte(token))
			varName := token
			prop, ok := s.props[varName]
			if !ok {
				c.throwAtLine("Unknown property '" + varName + "' in struct '" + _type + "'")
			}
			c.expectToken(":")
			// Read value
			c.assignValue(prop._type, prop._typeOfType)

			token = c.getNextToken()
		}

	} else if _typeOfType == "class" {

	}
}
