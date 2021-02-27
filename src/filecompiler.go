package main

import (
	"fmt"
	"os"
	"sort"
)

type FileCompiler struct {
	name string

	compiler *Compiler

	scopes     []*Scope
	scopeIndex int
	exitScope  bool

	index        int
	maxIndex     int
	line         int
	col          int
	lastTokenCol int
	whitespace   string
	indentLevel  int

	imports       []*Import
	exports       map[string]string
	exportVarName string

	compiled     bool
	code         []byte
	result       string
	resultBlocks []string
}

func (fc *FileCompiler) recordResult() {
	fc.resultBlocks = append(fc.resultBlocks, "")
}
func (fc *FileCompiler) getRecording() string {
	result := fc.resultBlocks[len(fc.resultBlocks)-1]
	fc.resultBlocks = fc.resultBlocks[:len(fc.resultBlocks)-1]
	return result
}
func (fc *FileCompiler) addResult(str string) {
	if len(fc.resultBlocks) > 0 {
		fc.resultBlocks[len(fc.resultBlocks)-1] += str
	} else {
		fc.result += str
	}
}

var extraSpace = 1

func (fc *FileCompiler) addSpace() {
	fc.addResult(fc.getSpace())
}
func (fc *FileCompiler) getSpace() string {
	i := -1 - extraSpace
	result := ""
	for i < fc.scopeIndex {
		result += "    "
		i++
	}
	return result
}

func (fc *FileCompiler) reset() {
	fc.line = 1
	fc.index = 0
	fc.col = 0
	fc.maxIndex = len(fc.code) - 1
}

func (fc *FileCompiler) compile() string {

	for fc.index <= fc.maxIndex {
		fc.handleNextWord()
		if fc.exitScope {
			fc.exitScope = false
			break
		}
	}

	if fc.exitScope {
		fc.throwAtLine("You forgot to close a certain scope, expected: \"}\"")
	}

	return fc.result
}

func (fc *FileCompiler) createNewScope() {
	s := Scope{
		types: map[string]string{},
		vars:  map[string]Var{},
	}
	fc.scopes = append(fc.scopes, &s)
	fc.scopeIndex++
}

func (fc *FileCompiler) getScope() *Scope {
	return fc.scopes[fc.scopeIndex]
}
func (fc *FileCompiler) popScope() {
	fc.scopes = fc.scopes[:len(fc.scopes)-1]
	fc.scopeIndex--
}

func (fc *FileCompiler) getNextToken(readOnly bool, sameLine bool) string {

	indexAtStart := fc.index
	if !readOnly {
		fc.lastTokenCol = fc.col
		fc.whitespace = ""
	}

	word := ""
	for fc.index <= fc.maxIndex {

		charInt := fc.code[fc.index]
		char := string(charInt)

		if isNewLine(charInt) {
			if sameLine {
				if len(word) == 0 && !readOnly {
					fc.throwAtLine("Unexpected new line")
				}
				break
			}
			fc.index++
			if !readOnly {
				fc.line++
				fc.col = 0
				fc.lastTokenCol = 0
				fc.whitespace = ""
				// lastChars := ""
				// if len(fc.result) > 1 {
				// 	lastChars = fc.result[len(fc.result)-2:]
				// 	if lastChars != "\n\n" {
				// 		fc.addResult("\n")
				// 		fc.addSpace()
				// 	}
				// }
			}
			if len(word) == 0 {
				continue
			}
			break
		}

		if isWhiteSpace(charInt) && len(word) == 0 {
			fc.index++
			if !readOnly {
				fc.col++
				fc.lastTokenCol++
				fc.whitespace += char
			}
			continue
		}

		if isVarNameChar(charInt) || (len(word) > 0 && isNumberChar(charInt)) || (len(word) == 0 && char == "#") {
			word += char
			fc.index++
			if !readOnly {
				fc.col++
			}
			continue
		}

		if len(word) == 0 {
			fc.index++
			if !readOnly {
				fc.col++
			}
			word = char
			// Comments
			if char == "/" {
				nextChar := fc.readNextChar()
				if nextChar == "/" || nextChar == "*" {
					word += nextChar
					fc.index++
					if !readOnly {
						fc.col++
					}
					break
				}
			}
			// Operators
			i := sort.SearchStrings(operatorChars, char)
			if i < len(operatorChars) && operatorChars[i] == char {
				// If operator char
				nextChar := fc.readNextChar()
				i = sort.SearchStrings(operatorChars, nextChar)
				if i < len(operatorChars) && operatorChars[i] == nextChar {
					word += nextChar
					fc.index++
					if !readOnly {
						fc.col++
					}
				}
				nextChar = fc.readNextChar()
				if nextChar == "=" {
					word += nextChar
					fc.index++
					if !readOnly {
						fc.col++
					}
				}
			}
			break
		}

		break
	}

	if readOnly {
		fc.index = indexAtStart
	}

	return word
}

func (fc *FileCompiler) getNextTokenSameLine() string {

	fc.lastTokenCol = fc.col

	word := ""
	for fc.index <= fc.maxIndex {

		charInt := fc.code[fc.index]
		char := string(charInt)

		if isNewLine(charInt) {
			if len(word) == 0 {
				fc.throwAtLine("Unexpected new line")
			}
			break
		}

		if isWhiteSpace(charInt) && len(word) == 0 {
			fc.index++
			fc.col++
			fc.lastTokenCol++
			continue
		}

		if isVarNameChar(charInt) || (len(word) == 0 && char == "#") {
			word += char
			fc.index++
			fc.col++
			continue
		}

		// If any special character
		if len(word) == 0 {
			fc.index++
			fc.col++
			return char
		}

		break
	}

	return word
}

func (fc *FileCompiler) getNextCharacterOnLine() string {

	fc.lastTokenCol = fc.col

	for fc.index <= fc.maxIndex {

		charInt := fc.code[fc.index]
		char := string(charInt)

		if isNewLine(charInt) {
			fc.index++
			fc.col = 0
			fc.lastTokenCol = 0
			fc.line++
			return char
		}

		fc.index++
		fc.col++

		if isWhiteSpace(charInt) {
			fc.lastTokenCol++
			continue
		}

		return char
	}

	fc.throwAtLine("Unexpected end of file")
	return ""
}

func (fc *FileCompiler) getNextValueToken() (string, string) {

	fc.lastTokenCol = fc.col

	vtype := ""
	word := ""
	prevChar := ""
	char := ""
	inStr := false
	hasDot := false
	endStrChar := ""

	for fc.index <= fc.maxIndex {

		prevChar = char
		charInt := fc.code[fc.index]
		char = string(charInt)

		if inStr {
			word += char
			fc.index++
			fc.col++
			fc.lastTokenCol++
			if char == endStrChar {
				break
			}
			if isNewLine(charInt) {
				if prevChar != "\\" {
					// prevent new line
					fc.throwAtLine("Unexpected new line")
				}
				fc.line++
				fc.col = 0
				fc.lastTokenCol = 0
			}
			continue
		}

		if isNewLine(charInt) {
			if len(word) == 0 {
				fc.throwAtLine("Unexpected new line")
			}
			break
		}

		if isWhiteSpace(charInt) && len(word) == 0 {
			fc.index++
			fc.col++
			fc.lastTokenCol++
			continue
		}

		// Strings
		if len(word) == 0 && (char == "\"" || char == "'") {
			word += char
			endStrChar = char
			fc.index++
			fc.col++
			vtype = "str"
			inStr = true
			continue
		}

		if isVarNameChar(charInt) {
			if vtype != "" && vtype != "word" {
				fc.throwAtLine("Unexpected char: " + char)
			}
			word += char
			fc.index++
			fc.col++
			vtype = "word"
			continue
		}

		isDot := char == "."
		if isNumberChar(charInt) || isDot {
			if vtype != "" && vtype != "num" {
				fc.throwAtLine("Unexpected char: " + char)
			}
			if isDot {
				if hasDot || len(word) == 0 {
					fc.throwAtLine("Unexpected char: " + char)
				}
				hasDot = true
			}
			word += char
			fc.index++
			fc.col++
			vtype = "num"
			continue
		}

		break
	}

	if len(word) == 0 {
		fc.throwAtLine("Missing value")
	}

	// if ends with a dot
	if string(word[len(word)-1]) == "." {
		fc.throwAtLine("Unexpected dot")
	}

	if word == "true" || word == "false" {
		return word, "bool"
	}

	if word == "[]" {
		return word, "array"
	}

	if word == "{}" {
		return word, "object"
	}

	if vtype == "str" {
		return word, "string"
	}
	if vtype == "num" {
		if hasDot {
			return word, "float"
		}
		return word, "int"
	}

	fc.throwAtLine("Unknown value: " + word)
	return "", ""
}

func (fc *FileCompiler) getNextType() *VarType {
	token := fc.getNextToken(false, true)
	result := &VarType{
		name: token,
	}
	i := sort.SearchStrings(basicTypes, token)
	if i < len(basicTypes) && basicTypes[i] == token {
		if token == "array" || token == "object" {
			fc.expectToken("<")
			result.subtype = fc.getNextType()
			fc.expectToken(">")
		}
		if token == "func" {
			fc.expectToken("(")
			currentIndex := fc.index
			currentCol := fc.col
			ntoken := fc.getNextTokenSameLine()
			if ntoken != ")" {
				fc.col = currentCol
				fc.index = currentIndex
			}
			for ntoken != ")" {
				ptype := fc.getNextType()
				result.paramTypes = append(result.paramTypes, ptype)
				ntoken = fc.getNextTokenSameLine()
				if ntoken != "," && ntoken != ")" {
					fc.throwAtLine("Unexpected token: " + ntoken)
				}
			}
			rtype := fc.getNextType()
			result.returnType = rtype
		}
	} else {
		_type, foundType := fc.getType(token)
		if !foundType {
			fc.throwAtLine("Unknown type: " + token)
		}
		result = _type
		nchar := fc.readNextChar()
		if nchar == "<" {
			fc.expectToken("<")
			result.subtype = fc.getNextType()
			fc.expectToken(">")
		}
	}
	nchar := fc.readNextChar()
	for nchar == "|" {
		fc.expectToken("|")
		ptype := fc.getNextTokenSameLine()
		if ptype == "null" {
			result.nullable = true
		} else if ptype == "undefined" {
			result.undefined = true
		} else {
			fc.throwAtLine("Expected null or undefined")
		}
		nchar = fc.readNextChar()
	}

	return result
}

func (fc *FileCompiler) readNextChar() string {
	if fc.index == fc.maxIndex {
		return ""
	}
	return string(fc.code[fc.index])
}

func (fc *FileCompiler) expectToken(token string) {
	ntoken := fc.getNextTokenSameLine()
	if ntoken != token {
		fc.throwAtLine("Expected: " + token)
	}
}

func (fc *FileCompiler) handleNextWord() {

	token := fc.getNextToken(false, false)

	if token == "" {
		// End of file
		return
	}

	if token == "#" {
		word := fc.getNextToken(false, false)
		fc.handleMacro(word)
		return
	}
	if token == "}" && fc.scopeIndex > 0 {
		fc.exitScope = true
		return
	}

	isDefine := token == "define"
	if isDefine {
		token = fc.getNextTokenSameLine()
	}

	if token == "return" {
		fc.addResult("return ")
		scope := fc.getScope()
		if scope.returnType == nil {
			fc.throwAtLine("Unexpected return statement")
		}
		if scope.returnType.name == "void" {
			return
		}
		if fc.compiler.readTypes {
			fc.skipValue()
			return
		}
		rtype := fc.assignValue()
		fc.addResult(";")
		if !scope.returnType.isCompatible(rtype) {
			fc.throwTypeError(scope.returnType, rtype)
		}
		return
	}

	if token == "while" {
		fc.throwAtLine("Feature not ready yet")
		return
	}
	if token == "for" {
		fc.throwAtLine("Feature not ready yet")
		return
	}
	if token == "if" {
		fc.handleIf()
		return
	}
	if token == "switch" {
		fc.throwAtLine("Feature not ready yet")
		return
	}
	if token == "function" {
		fc.throwAtLine("Global functions are currently not allowed")
		return
	}
	if token == "var" || token == "let" {
		fc.throwAtLine("Declaring variables with var/let is not allowed, use a type instead")
		return
	}

	if token == "struct" || token == "local" {
		isLocal := token == "local"
		if isLocal {
			token = fc.getNextTokenSameLine()
			if token != "struct" {
				fc.throwAtLine("Unexpected token: " + token)
			}
		}
		if fc.compiler.readTypes {
			fc.handleTypeSkip(isLocal, isDefine, true, false)
			return
		}
		fc.handleType(isLocal, isDefine, true, false)
		return
	}

	if token == "class" {
		if fc.compiler.readTypes {
			fc.handleTypeSkip(false, isDefine, false, true)
			return
		}
		fc.handleType(false, isDefine, false, true)
		return
	}

	if token == "export" {
		fc.handleExport()
		return
	}

	if len(token) == 0 {
		return
	}

	if token == "import" {
		fc.handleImport()
		return
	}

	if token == "//" || token == "/*" {
		fc.handleComment()
		return
	}

	// Check if type/struct/class
	if fc.isValidType(token) {
		fc.index -= len(token)
		fc.col -= len(token)
		if fc.col < 0 {
			fc.col = 0
		}
		_type := fc.getNextType()
		if fc.compiler.readTypes {
			fc.getNextToken(false, true)
			if !isDefine {
				fc.expectToken("=")
				fc.skipValue()
				token := fc.getNextToken(true, false)
				if token == ";" {
					token = fc.getNextToken(false, false)
				}
			}
			return
		}
		fc.declareVariable(_type, isDefine)
		return
	}

	if fc.compiler.readTypes && (isVarNameSyntax([]byte(token)) || token == "(" || token == "[") {
		fc.index -= len(token)
		fc.col -= len(token)

		fc.skipValue()
		nextToken := fc.getNextToken(true, false)
		if nextToken == "=" {
			fc.getNextToken(false, false)
			fc.skipValue()
		}
		nextChar := fc.getNextToken(true, false)
		if nextChar == ";" {
			fc.getNextToken(false, false)
		}
		return
	}

	_, isVar := fc.getVar(token)
	if isVar || token == "(" || token == "[" {
		fc.index -= len(token)
		fc.col -= len(token)

		vt := fc.assignValue()
		if vt.assignable {
			// Check for = sign
			nextToken := fc.getNextToken(true, false)
			if nextToken == "=" {
				nextToken = fc.getNextToken(false, false)
				fc.addResult(" = ")
				assignType := fc.assignValue()
				if !vt.isCompatible(assignType) {
					fc.throwTypeError(vt, assignType)
				}
			}
		} else {
			nextToken := fc.getNextToken(true, false)
			if nextToken == "=" {
				fc.throwAtLine("Cannot assign a value to this")
			}
		}
		nextChar := fc.getNextToken(true, false)
		if nextChar == ";" {
			nextChar = fc.getNextToken(false, false)
		}
		fc.addResult(";")
		// Todo: check for missing props if left is struct
		return
	}

	// Unknown
	if isVarNameSyntax([]byte(token)) {
		fc.throwAtLine("Unknown variable: " + token)
	}

	fc.throwAtLine("Unexpected token: " + token)
}

func (fc *FileCompiler) getType(name string) (*VarType, bool) {
	var sci = fc.scopeIndex
	for sci >= 0 {
		scope := fc.scopes[sci]
		realName, ok := scope.types[name]
		if ok {
			result, ok := allTypes[realName]
			return result, ok
		}
		sci--
	}
	return nil, false
}

func (fc *FileCompiler) getVar(name string) (*Var, bool) {
	var sci = fc.scopeIndex
	for sci >= 0 {
		scope := fc.scopes[sci]
		result, ok := scope.vars[name]
		if ok {
			return &result, ok
		}
		sci--
	}
	return nil, false
}

func (fc *FileCompiler) typeExists(name string) bool {
	var sci = fc.scopeIndex
	for sci >= 0 {
		scope := fc.scopes[sci]
		if scope.typeExists(name) {
			return true
		}
		sci--
	}
	return false
}

func (fc *FileCompiler) isValidType(token string) bool {
	i := sort.SearchStrings(basicTypes, token)
	if i < len(basicTypes) && basicTypes[i] == token {
		return true
	}
	return fc.typeExists(token)
}

func (fc *FileCompiler) throwAtLine(msg string) {

	fmt.Print("\033[31m") // Color red
	fmt.Println(msg)
	fmt.Print("\033[0m") // Color reset
	fmt.Println("Line", fc.line, "col", fc.col, "in", fc.name)
	fmt.Println(fc.readLine(fc.line))
	fmt.Print("\033[31m") // Color red
	i := 0
	mark := ""
	for i < fc.lastTokenCol {
		mark += " "
		i++
	}
	mark += "^"
	fmt.Println(mark)
	fmt.Print("\033[0m") // Color reset

	panic("---")
	// os.Exit(1)
}

func (fc *FileCompiler) throw(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func (fc *FileCompiler) readLine(lineNr int) string {
	i := 0
	line := ""
	currentLine := 1
	for i <= fc.maxIndex {

		charInt := fc.code[i]
		char := string(charInt)
		isLF := isNewLine(charInt)

		if currentLine == lineNr && !isLF {
			line += char
		}

		if isLF {
			currentLine++
			if currentLine > lineNr {
				break
			}
		}

		i++
	}
	return line
}
