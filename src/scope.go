package main

import (
	"sort"
	"strings"
)

type Var struct {
	_type *VarType // string bool int MyStruct MyClass etc...
}

type Scope struct {
	types      map[string]string
	vars       map[string]Var
	returnType *VarType
}

func (s *Scope) typeExists(name string) bool {
	_, ok := s.types[name]
	if ok {
		return true
	}
	return false
}

func (s *Scope) getVar(name string) (*Var, bool) {
	result, ok := s.vars[name]
	return &result, ok
}

func (fc *FileCompiler) declareVariable(_type *VarType, isDefine bool) {
	varName := fc.getNextToken(false, true)
	fc.checkVarNameSyntax([]byte(varName))
	_, ok := fc.getVar(varName)
	if ok {
		fc.throwAtLine("Variable name already in use: " + varName)
	}
	if fc.typeExists(varName) {
		fc.throwAtLine("Name already used as a class/struct: " + varName)
	}
	if !isDefine {
		fc.addResult("var " + varName)
		fc.expectToken("=")
		fc.addResult(" = ")
		rightType := fc.assignValue()
		if !_type.isCompatible(rightType) {
			fc.throwTypeError(_type, rightType)
		}
		fc.addResult(";")
	}
	scope := fc.scopes[fc.scopeIndex]
	scope.vars[varName] = Var{
		_type: _type,
	}
}

func (fc *FileCompiler) assignValue() *VarType {

	var result *VarType

	token := fc.getNextToken(false, true)
	if token == "(" {
		fc.addResult("(")
		result = fc.assignValue()
		fc.expectToken(")")
		fc.addResult(")")
	} else if isNumberChar([]byte(token)[0]) {
		// Number
		number := token
		nextChar := fc.readNextChar()
		if nextChar == "." {
			number += fc.getNextToken(false, true)
			number += fc.getNextToken(false, true)
		}
		nextChar = fc.readNextChar()
		if !isNumberSyntax([]byte(number)) || nextChar == "." {
			fc.throwAtLine("Invalid number")
		}
		fc.addResult(number)
		_type := VarType{
			name: "number",
		}
		result = &_type
	} else if token == "null" {
		fc.addResult("null")
		_type := VarType{
			name: "null",
		}
		result = &_type
	} else if token == "undefined" {
		fc.addResult("undefined")
		_type := VarType{
			name: "undefined",
		}
		result = &_type
	} else if token == "new" {
		// Classes
		fc.addResult("new ")
		className := fc.getNextToken(false, true)
		_type, typeExists := fc.getType(className)
		if !typeExists || !_type.isClass {
			fc.throwAtLine("Unknown class: " + className)
		}
		fc.addResult(_type.name)
		rtype := VarType{
			name:    className,
			isClass: true,
		}
		result = &rtype
		fc.expectToken("(")
		fc.addResult("(")

		constructorProp, hasConstructor := _type.props["constructor"]
		constructorType := constructorProp.varType

		if !hasConstructor {
			fc.expectToken(")")
		} else {

			// Get params
			nextChar := fc.getNextToken(true, true)
			params := []*VarType{}
			for nextChar != ")" {
				np := fc.assignValue()
				params = append(params, np)
				nextChar = fc.getNextToken(true, false)
				if nextChar == "," {
					nextChar = fc.getNextToken(false, false)
					fc.addResult(",")
					nextChar = fc.getNextToken(true, false)
				}
			}
			nextChar = fc.getNextToken(false, false)
			fc.addResult(")")

			// Check param types
			paramCount := len(params)
			if paramCount > len(constructorType.paramTypes) {
				fc.throwAtLine("Too many params")
			}
			for i, pt := range constructorType.paramTypes {
				if i < paramCount {
					p := params[i]
					if !pt.isCompatible(p) {
						fc.throwTypeError(pt, p)
					}
				} else {
					// Check if undefined is allowed
					if !pt.undefined {
						fc.throwAtLine("Missing params")
					}
				}
			}
		}

		//fc.throwAtLine("Class values not supported yet")
	} else if token == "function" {
		// Functions
		fc.expectToken("(")
		fc.addResult("function(")

		_type := VarType{
			name:       "func",
			paramTypes: []*VarType{},
		}
		result = &_type

		fc.createNewScope()
		scope := fc.getScope()

		ntoken := fc.getNextToken(true, false)
		if ntoken == ")" {
			ntoken = fc.getNextToken(false, false)
			fc.addResult(")")
		}
		for ntoken != ")" {
			paramName := fc.getNextToken(false, false)
			fc.addResult(paramName)
			ptype := fc.getNextType()
			ptype.paramName = paramName
			result.paramTypes = append(result.paramTypes, ptype)
			scope.vars[paramName] = Var{
				_type: ptype,
			}
			ntoken = fc.getNextTokenSameLine()
			if ntoken != "," && ntoken != ")" {
				fc.throwAtLine("Unexpected token: " + ntoken)
			}
			fc.addResult(ntoken)
		}
		rtype := fc.getNextType()
		result.returnType = rtype

		fc.expectToken("{")
		fc.addResult("{")
		fc.scopes[fc.scopeIndex].returnType = result.returnType
		fc.compile()
		if string(fc.code[fc.index-1]) != "}" {
			fc.throwAtLine("Expected: }")
		}
		fc.addResult("}")
		fc.popScope()

	} else if isVarNameSyntax([]byte(token)) {
		// Vars
		_var, ok := fc.getVar(token)
		if !ok {
			fc.throwAtLine("Undefined variable: " + token)
		}
		// Is variable name
		fc.addResult(token)
		result = _var._type
		result.assignable = true
	} else if token == "\"" || token == "'" {
		// String
		fc.addResult(token)
		char := ""
		lastChar := ""
		for fc.index <= fc.maxIndex {
			lastChar = char
			charInt := fc.code[fc.index]
			char = string(charInt)
			fc.index++
			fc.col++
			fc.lastTokenCol++
			fc.addResult(char)
			if isNewLine(charInt) {
				if lastChar != "\\" {
					fc.throwAtLine("Unexpected newline")
				}
				fc.line++
				fc.col = 0
				fc.lastTokenCol = 0
			}
			if char == token && lastChar != "\\" {
				break
			}
		}
		if fc.index > fc.maxIndex {
			fc.throwAtLine("You forgot to close a string somewhere")
		}

		_type := VarType{
			name: "string",
		}
		result = &_type

	} else if token == "[" {
		// Array
		fc.addResult("[")

		returnType := VarType{
			name:       "array",
			assignable: false,
		}

		token = fc.getNextToken(true, false)
		if token == "]" {
			token = fc.getNextToken(false, false)
		} else {
			for token != "]" {
				if token == "" {
					fc.throwAtLine("Unexpected end of code")
				}

				subType := fc.assignValue()
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
						fc.throwTypeError(returnType.subtype, subType)
					}
				}

				token = fc.getNextToken(false, false)
				if token == "," {
					fc.addResult(",")
				}
			}
		}
		fc.addResult("]")

		result = &returnType

	} else if token == "{" {

		fc.addResult("{")
		returnType := VarType{
			name:       "object",
			props:      map[string]*Property{},
			assignable: false,
		}

		// s, _ := fc.getStruct(leftType.name)
		extraSpace++
		token := fc.getNextToken(false, false)
		for token != "}" {
			if token == "" {
				fc.throwAtLine("Unexpected end of code")
			}
			fc.checkVarNameSyntax([]byte(token))
			varName := token
			fc.addResult(varName)
			fc.expectToken(":")
			fc.addResult(":")
			// Read value
			propType := fc.assignValue()
			newProp := Property{
				varType: propType,
			}
			returnType.props[varName] = &newProp

			token = fc.getNextToken(true, false)
			if token == "}" {
				extraSpace--
			}
			token = fc.getNextToken(false, false)
			if token == "," {
				fc.addResult(",")
				token = fc.getNextToken(false, false)
			}
		}

		// todo: Autofill missing fields

		//
		fc.addResult("}")
		result = &returnType
	} else {
		fc.throwAtLine("Setting value type '" + token + "' is not supported yet")
	}

	// Handle trailing . [ or (
	for {
		nextChar := fc.readNextChar()
		if (result.name == "array") && nextChar == "[" {
			fc.throwAtLine("Array accessors not ready yet")
		} else if (result.name == "func") && nextChar == "(" {
			// Function
			nextChar := fc.getNextToken(false, true)
			fc.addResult("(")
			nextChar = fc.getNextToken(true, true)
			// Check param types
			params := []*VarType{}
			for nextChar != ")" {
				np := fc.assignValue()
				params = append(params, np)
				nextChar = fc.getNextToken(true, false)
				if nextChar == "," {
					nextChar = fc.getNextToken(false, false)
					fc.addResult(",")
					nextChar = fc.getNextToken(true, false)
				}
			}
			nextChar = fc.getNextToken(false, false)
			fc.addResult(")")
			paramCount := len(params)
			if paramCount > len(result.paramTypes) {
				fc.throwAtLine("Too many params")
			}
			for i, pt := range result.paramTypes {
				if i < paramCount {
					p := params[i]
					if !pt.isCompatible(p) {
						fc.throwTypeError(pt, p)
					}
				} else {
					// Check if undefined is allowed
					if !pt.undefined {
						fc.throwAtLine("Missing params")
					}
				}
			}

			// Set result to function returnType
			result = result.returnType
			result.assignable = false

		} else if nextChar == "." || nextChar == "[" {
			nextChar := fc.getNextToken(false, true)
			if nextChar == "[" {
				fc.throwAtLine("Dynamic properties are not allowed")
			}
			fc.addResult(".")
			propName := fc.getNextToken(false, true)
			if len(fc.whitespace) > 0 {
				fc.throwAtLine("Unexpected whitespace")
			}
			fc.addResult(propName)
			// check if struct
			_type, ok := fc.getType(result.name)
			if ok {
				prop, ok := _type.props[propName]
				if !ok {
					fc.throwAtLine("Undefined property: " + propName + " on: " + result.name)
				}
				result = prop.varType
				result.assignable = true
			} else {
				fc.throwAtLine("Cannot load struct/class: " + result.name + " (compiler bug)")
			}
		} else {
			break
		}
	}

	// Handle operators
	nextToken := fc.getNextToken(true, true)
	i := sort.SearchStrings(operators, nextToken)
	if i < len(operators) && operators[i] == nextToken {
		nextToken = fc.getNextToken(false, false)
		if nextToken != "++" && nextToken != "--" {
			fc.addResult(" " + nextToken + " ")
		} else {
			fc.addResult(nextToken)
		}
		rightType := fc.assignValue()
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
			fc.throwAtLine("Operator not supported yet: '" + nextToken + "'")
		}
		if showError {
			fc.throwAtLine("Cannot use operator '" + nextToken + "' on type " + result.displayName() + " && " + rightType.displayName() + "")
		}
	}

	return result
}
