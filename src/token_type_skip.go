package main

func (fc *FileCompiler) handleTypeSkip(isLocal bool, isDefine bool, isStruct bool, isClass bool) {

	if (isStruct && isClass) || (!isStruct && !isClass) {
		fc.throwAtLine("Compiler bug, new type must either be struct or class, not both")
	}

	if !isLocal {
		if fc.scopeIndex > 0 {
			fc.throwAtLine("You must use 'local' to create a local struct")
		}
	}

	name := fc.getNextToken(false, true)
	fc.checkVarNameSyntax([]byte(name))

	if fc.typeExists(name) {
		fc.throwAtLine("Type name already in use: " + name)
	}

	var class *VarType = nil
	var s *VarType = nil
	var globalName string

	scope := fc.scopes[fc.scopeIndex]
	if isStruct {
		globalName, s = createNewType(false, name)
		s.isLocal = isLocal
		scope.types[name] = globalName
	}

	if isClass {
		globalName, class = createNewType(true, name)
		class.isLocal = isLocal
		scope.types[name] = globalName
	}

	token := fc.getNextToken(false, true)
	if isDefine && isStruct && token == "," {
		fc.expectToken("class")
		className := fc.getNextToken(false, true)
		fc.checkVarNameSyntax([]byte(className))
		if fc.typeExists(className) {
			fc.throwAtLine("Struct/class name already in use: " + className)
		}
		globalName, class = createNewType(true, className)
		scope.types[className] = globalName
		token = fc.getNextTokenSameLine()
	}

	// Get fields
	if token != "{" {
		fc.throwAtLine("Unexpected token: " + token)
	}

	// Class functions
	// functionCode := map[string]string{}

	token = fc.getNextToken(false, false)
	for token != "}" {
		if token == "" {
			fc.throwAtLine("Unexpected end of code")
		}
		fc.checkVarNameSyntax([]byte(token))
		varName := token

		if s != nil {
			_, ok := s.props[varName]
			if ok {
				fc.throwAtLine("Property name '" + varName + "' already exists")
			}
		}
		if class != nil {
			_, ok := class.props[varName]
			if ok {
				fc.throwAtLine("Property name '" + varName + "' already exists")
			}
		}

		fc.expectToken(":")

		// Read type
		prop := Property{}
		token = fc.getNextToken(true, true)

		if isClass && varName == "constructor" && token != "function" {
			fc.throwAtLine("Constructor must be a function")
		}

		if isClass && !isDefine && token == "function" {
			fc.getNextToken(false, true)
			t := fc.skipFunction()
			// _typeOfType, _ := fc.getTypeOfType(t.name)
			// t.toft = _typeOfType
			prop.varType = t
		} else {
			t := fc.getNextType()
			// _typeOfType, _ := fc.getTypeOfType(t.name)
			// t.toft = _typeOfType
			prop.varType = t
		}

		// Check default
		char := fc.getNextToken(true, true)
		if char == "=" {
			char = fc.getNextToken(false, true)
			fc.skipValue()
		}

		// Store property
		if s != nil {
			s.props[varName] = &prop
		}
		if class != nil {
			class.props[varName] = &prop
		}

		token = fc.getNextToken(false, false)
	}
}
