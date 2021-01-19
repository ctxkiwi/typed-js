package main

import (
	"sort"
	"strings"
)

type Var struct {
	_type *VarType // string bool int MyStruct MyClass etc...
}

type Scope struct {
	structs    map[string]string
	classes    map[string]string
	vars       map[string]Var
	returnType *VarType
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

func (c *Compile) declareVariable(_type *VarType, isDefine bool) {
	varName := c.getNextToken(false, true)
	c.checkVarNameSyntax([]byte(varName))
	_, ok := c.getVar(varName)
	if ok {
		c.throwAtLine("Variable name already in use: " + varName)
	}
	if c.typeExists(varName) {
		c.throwAtLine("Name already used as a class/struct: " + varName)
	}
	if !isDefine {
		c.addResult(c.whitespace + "var " + varName)
		c.expectToken("=")
		c.addResult(" = ")
		rightType := c.assignValue()
		if !_type.isCompatible(rightType) {
			c.throwTypeError(_type, rightType)
		}
		c.addResult(";\n")
	}
	scope := scopes[scopeIndex]
	scope.vars[varName] = Var{
		_type: _type,
	}
}

func (c *Compile) assignValue() *VarType {

	var result *VarType

	token := c.getNextToken(false, true)
	if token == "(" {
		c.addResult("(")
		result = c.assignValue()
		c.expectToken(")")
		c.addResult(")")
	} else if isNumberChar([]byte(token)[0]) {
		// Number
		number := token
		nextChar := c.readNextChar()
		if nextChar == "." {
			number += c.getNextToken(false, true)
			number += c.getNextToken(false, true)
		}
		nextChar = c.readNextChar()
		if !isNumberSyntax([]byte(number)) || nextChar == "." {
			c.throwAtLine("Invalid number")
		}
		c.addResult(number)
		_type := VarType{
			name: "number",
		}
		result = &_type
	} else if token == "null" {
		c.addResult("null")
		_type := VarType{
			name: "null",
		}
		result = &_type
	} else if token == "undefined" {
		c.addResult("undefined")
		_type := VarType{
			name: "undefined",
		}
		result = &_type
	} else if token == "new" {
		// Classes
		c.addResult("new ")
		className := c.getNextToken(false, true)
		class, classExists := c.getClass(className)
		if !classExists {
			c.throwAtLine("Unknown class: " + className)
		}
		c.addResult(className)
		_type := VarType{
			name: className,
		}
		result = &_type
		c.expectToken("(")
		c.addResult("(")

		constructorProp, hasConstructor := class.props["constructor"]
		constructorType := constructorProp.varType

		if !hasConstructor {
			c.expectToken(")")
		} else {

			// Get params
			nextChar := c.getNextToken(true, true)
			params := []*VarType{}
			for nextChar != ")" {
				np := c.assignValue()
				params = append(params, np)
				nextChar = c.getNextToken(true, false)
				if nextChar == "," {
					nextChar = c.getNextToken(false, false)
					c.addResult(",")
					nextChar = c.getNextToken(true, false)
				}
			}
			nextChar = c.getNextToken(false, false)
			c.addResult(")")

			// Check param types
			paramCount := len(params)
			if paramCount > len(constructorType.paramTypes) {
				c.throwAtLine("Too many params")
			}
			for i, pt := range constructorType.paramTypes {
				if i < paramCount {
					p := params[i]
					if !pt.isCompatible(p) {
						c.throwTypeError(pt, p)
					}
				} else {
					// Check if undefined is allowed
					if !pt.undefined {
						c.throwAtLine("Missing params")
					}
				}
			}
		}

		// c.throwAtLine("Class values not supported yet")
	} else if token == "function" {
		// Functions
		c.expectToken("(")
		c.addResult("function(")

		_type := VarType{
			name:       "func",
			paramTypes: []*VarType{},
		}
		result = &_type

		createNewScope()
		scope := getScope()

		ntoken := c.getNextToken(true, false)
		if ntoken == ")" {
			ntoken = c.getNextToken(false, false)
			c.addResult(")")
		}
		for ntoken != ")" {
			paramName := c.getNextToken(false, false)
			c.addResult(paramName)
			ptype := c.getNextType()
			ptype.paramName = paramName
			result.paramTypes = append(result.paramTypes, ptype)
			scope.vars[paramName] = Var{
				_type: ptype,
			}
			ntoken = c.getNextTokenSameLine()
			if ntoken != "," && ntoken != ")" {
				c.throwAtLine("Unexpected token: " + ntoken)
			}
			c.addResult(ntoken)
		}
		rtype := c.getNextType()
		result.returnType = rtype

		c.expectToken("{")
		c.addResult("{\n")
		scopes[scopeIndex].returnType = result.returnType
		c.compile()
		if string(c.code[c.index-1]) != "}" {
			c.throwAtLine("Expected: }")
		}
		c.addResult("}")
		popScope()

	} else if isVarNameSyntax([]byte(token)) {
		// Vars
		_var, ok := c.getVar(token)
		if ok {
			// Is variable name
			c.addResult(token)
			result = _var._type
		} else {
			c.throwAtLine("Undefined variable: " + token)
		}
	} else if token == "\"" || token == "'" {
		// String
		c.addResult(token)
		char := ""
		lastChar := ""
		for c.index <= c.maxIndex {
			lastChar = char
			charInt := c.code[c.index]
			char = string(charInt)
			c.index++
			c.col++
			c.lastTokenCol++
			c.addResult(char)
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
		c.addResult(c.whitespace + "[")

		returnType := VarType{
			name:       "array",
			assignable: false,
		}

		token = c.getNextToken(true, false)
		if token == "]" {
			token = c.getNextToken(false, false)
		} else {
			for token != "]" {
				if token == "" {
					c.throwAtLine("Unexpected end of code")
				}

				subType := c.assignValue()
				if returnType.subtype == nil {
					returnType.subtype = subType
				} else {
					if subType.name == "null" {
						returnType.subtype.nullable = true
					}
					if subType.name == "undefined" {
						returnType.subtype.undefined = true
					}
					if !returnType.subtype.isCompatible(subType) {
						c.throwTypeError(returnType.subtype, subType)
					}
				}

				token = c.getNextToken(false, false)
				if token == "," {
					c.addResult(",")
				}
			}
		}
		c.addResult(c.whitespace + "]")

		result = &returnType

	} else if token == "{" {

		c.addResult("{")
		returnType := VarType{
			name:       "object",
			props:      map[string]*Property{},
			assignable: false,
		}

		// s, _ := c.getStruct(leftType.name)
		token := c.getNextToken(false, false)
		for token != "}" {
			if token == "" {
				c.throwAtLine("Unexpected end of code")
			}
			c.checkVarNameSyntax([]byte(token))
			varName := token
			c.addResult(c.whitespace + varName)
			c.expectToken(":")
			c.addResult(":")
			// Read value
			propType := c.assignValue()
			newProp := Property{
				varType: propType,
			}
			returnType.props[varName] = &newProp

			token = c.getNextToken(false, false)
			if token == "," {
				c.addResult(",")
				token = c.getNextToken(false, false)
			}
		}

		// todo: Autofill missing fields

		//
		c.addResult(c.whitespace + "}")
		result = &returnType
	} else {
		c.throwAtLine("Setting value type '" + token + "' is not supported yet")
	}

	// Handle trailing . [ or (
	for {
		nextChar := c.readNextChar()
		if (result.name == "array") && nextChar == "[" {
			c.throwAtLine("Array accessors not ready yet")
		} else if (result.name == "func") && nextChar == "(" {
			// Function
			nextChar := c.getNextToken(false, true)
			c.addResult("(")
			nextChar = c.getNextToken(true, true)
			// Check param types
			params := []*VarType{}
			for nextChar != ")" {
				np := c.assignValue()
				params = append(params, np)
				nextChar = c.getNextToken(true, false)
				if nextChar == "," {
					nextChar = c.getNextToken(false, false)
					c.addResult(",")
					nextChar = c.getNextToken(true, false)
				}
			}
			nextChar = c.getNextToken(false, false)
			c.addResult(")")
			paramCount := len(params)
			if paramCount > len(result.paramTypes) {
				c.throwAtLine("Too many params")
			}
			for i, pt := range result.paramTypes {
				if i < paramCount {
					p := params[i]
					if !pt.isCompatible(p) {
						c.throwTypeError(pt, p)
					}
				} else {
					// Check if undefined is allowed
					if !pt.undefined {
						c.throwAtLine("Missing params")
					}
				}
			}

			// Set result to function returnType
			result = result.returnType
			result.assignable = false

		} else if nextChar == "." || nextChar == "[" {
			nextChar := c.getNextToken(false, true)
			if nextChar == "[" {
				c.throwAtLine("Dynamic properties are not allowed")
			}
			c.addResult(".")
			propName := c.getNextToken(false, true)
			if len(c.whitespace) > 0 {
				c.throwAtLine("Unexpected whitespace")
			}
			c.addResult(propName)
			// check if struct
			s, ok := c.getStruct(result.name)
			if ok {
				prop, ok := s.props[propName]
				if !ok {
					c.throwAtLine("Undefined property: " + propName + " on struct: " + result.name)
				}
				result = prop.varType
				result.assignable = true
			} else {
				// check if class
				class, ok := c.getClass(result.name)
				if ok {
					prop, ok := class.props[propName]
					if !ok {
						c.throwAtLine("Undefined property: " + propName + " on class: " + result.name)
					}
					result = c.assignValue()
					if !prop.varType.isCompatible(result) {
						c.throwTypeError(prop.varType, result)
					}
				} else {
					c.throwAtLine("Cannot load struct/class: " + result.name + " (compiler bug)")
				}
			}
		} else {
			break
		}
	}

	// Handle operators
	nextToken := c.getNextToken(true, true)
	i := sort.SearchStrings(operators, nextToken)
	if i < len(operators) && operators[i] == nextToken {
		nextToken = c.getNextToken(false, false)
		if nextToken != "++" && nextToken != "--" {
			c.addResult(" " + nextToken + " ")
		} else {
			c.addResult(nextToken)
		}
		rightType := c.assignValue()
		showError := false
		leftLower := strings.ToLower(result.name)
		rightLower := strings.ToLower(rightType.name)
		switch nextToken {
		case "+":
			if leftLower == "number" && rightLower == "number" {
			} else if leftLower == "string" && (rightLower == "bool" || rightLower == "number") {
				result = createType("string")
			} else if rightLower == "string" && (leftLower == "bool" || leftLower == "number") {
				result = createType("string")
			} else {
				showError = true
			}
		case "-", "*", "/":
			if leftLower != "number" || rightLower != "number" {
				showError = true
			}
		case "==", "===":
			if !result.isCompatible(rightType) {
				showError = true
				break
			}
			result = createType("bool")
		case "<", ">", "<=", ">=":
			if (leftLower != "number" && leftLower != "string") || (rightLower != "number" && rightLower != "string") {
				showError = true
				break
			}
			result = createType("bool")
		case "&&", "||":
			if leftLower != "bool" || rightLower != "bool" {
				showError = true
			}
		default:
			c.throwAtLine("Operator not supported yet: '" + nextToken + "'")
		}
		if showError {
			c.throwAtLine("Cannot use operator '" + nextToken + "' on type " + result.displayName() + " && " + rightType.displayName() + "")
		}
	}

	return result
}
