package main

import (
	"sort"
)

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

// func (s *Scope) hasStruct(name string) bool {
// 	for _, str := range s.structs {
// 		if str == name {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (s *Scope) hasClass(name string) bool {
// 	for _, str := range s.classes {
// 		if str == name {
// 			return true
// 		}
// 	}
// 	return false
// }

func (s *Scope) getVar(name string) (*Var, bool) {
	result, ok := s.vars[name]
	return &result, ok
}

func (c *Compile) declareVariable(_type *VarType) {
	varName := c.getNextTokenSameLine()
	c.checkVarNameSyntax([]byte(varName))
	_, ok := c.getVar(varName)
	if ok {
		c.throwAtLine("Variable name already in use: " + varName)
	}
	c.result += c.whitespace + "var " + varName
	c.expectToken("=")
	c.result += " = "
	rightType := c.assignValue(_type)
	c.result += ";"
	if !_type.isCompatible(rightType) {
		c.throwTypeError(_type, rightType)
	}
	// value, _ := c.getNextValueToken()
	// c.result += value
}

func (c *Compile) assignValue(leftType *VarType) *VarType {

	var result *VarType

	token := c.getNextTokenSameLine()
	if token == "(" {
		c.result += "("
		result = c.assignValue(leftType)
		c.expectToken(")")
		c.result += ")"
	} else if isVarNameSyntax([]byte(token)) {
		// Vars
		_, ok := c.getVar(token)
		if ok {
			// Is variable name
			c.result += token
			// ...
		}
	} else if token == "null" {
		c.result += "null"
		_type := VarType{
			name: "null",
		}
		result = &_type
	} else if token == "undefined" {
		c.result += "undefined"
		_type := VarType{
			name: "undefined",
		}
		result = &_type
	} else if token == "\"" || token == "'" {
		// String
		c.result += token
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

		_type := VarType{
			name: "string",
		}
		result = &_type

	} else if token == "[" {
		// Array
		c.result += c.whitespace + "["
		token = c.getNextToken(true, false)
		for token != "]" {
			if token == "" {
				c.throwAtLine("Unexpected end of code")
			}

			valueType := c.assignValue(leftType.subtype)
			if !leftType.subtype.isCompatible(valueType) {
				c.throwTypeError(leftType.subtype, valueType)
			}

			token = c.getNextToken(false, false)
			if token == "," {
				c.result += ","
			}
		}
		c.result += c.whitespace + "]"

		result = leftType

	} else if token == "{" {

		c.result += "{"

		s, _ := c.getStruct(leftType.name)
		token := c.getNextToken(false, false)
		for token != "}" {
			if token == "" {
				c.throwAtLine("Unexpected end of code")
			}
			c.checkVarNameSyntax([]byte(token))
			varName := token
			c.result += c.whitespace + varName
			prop, ok := s.props[varName]
			if !ok {
				c.throwAtLine("Unknown property '" + varName + "' in struct '" + leftType.name + "'")
			}
			c.expectToken(":")
			c.result += ":"
			// Read value
			c.assignValue(prop.varType)

			token = c.getNextToken(false, false)
			if token == "," {
				c.result += ","
				token = c.getNextToken(false, false)
			}
		}

		// todo: Autofill missing fields

		//
		c.result += c.whitespace + "}"
		result = leftType

	} else if token == "new" {
		c.throwAtLine("Class values not supported yet")
	} else {
		c.throwAtLine("Setting value type '" + token + "' is not supported yet")
	}

	// Handle operators
	nextToken := c.getNextToken(true, true)
	i := sort.SearchStrings(operators, nextToken)
	if i < len(operators) && operators[i] == nextToken {
		nextToken = c.getNextToken(false, false)
		rightType := c.assignValue(result)
		c.result += nextToken
		showError := false
		switch nextToken {
		case "+":
		case "-", "*", "/":
		case "==", "===":
		case "<=", ">=":
		case "&&", "||":
		default:
			c.throwAtLine("Operator not supported yet: '" + nextToken + "'")
		}
		if showError {
			c.throwAtLine("Cannot use operator '" + nextToken + "' on type " + result.displayName() + " && " + rightType.displayName() + "")
		}
	}

	return result
}
