package main

import (
	"fmt"
	"os"
	"sort"
)

var basicTypes = []string{"bool", "number", "string", "array", "object", "func", "any", "T", "void"}
var structsEqualToClass = []string{"number", "string", "bool", "array", "object"}
var basicValues = []string{"true", "false", "undefined", "null", "[]", "{}"}
var operators = []string{"+", "-", "*", "/", "==", "===", "<", ">", "<=", ">=", "&&", "||", "++", "--"}
var operatorChars = []string{"+", "-", "*", "/", "=", "<", ">", "&", "|"}

var scopes []*Scope
var scopeIndex = -1

type Compile struct {
	name string

	index        int
	maxIndex     int
	line         int
	col          int
	lastTokenCol int
	whitespace   string
	exitScope    bool

	code   []byte
	result string
}

func compileCode(name string, code []byte) string {

	c := Compile{
		name: name,

		index:     0,
		maxIndex:  len(code) - 1,
		line:      1,
		exitScope: false,

		code:   code,
		result: "",
	}

	if len(scopes) == 0 {
		createNewScope()
	}

	return c.compile()
}

func (c *Compile) compile() string {

	for c.index <= c.maxIndex {
		c.handleNextWord()
		if c.exitScope {
			c.exitScope = false
			break
		}
	}

	return c.result
}

func createNewScope() {
	s := Scope{
		structs: map[string]string{},
		classes: map[string]string{},
		vars:    map[string]Var{},
	}
	scopes = append(scopes, &s)
	scopeIndex++
}
func getScope() *Scope {
	return scopes[scopeIndex]
}
func popScope() {
	scopes = scopes[:len(scopes)-1]
	scopeIndex--
}

func (c *Compile) getNextToken(readOnly bool, sameLine bool) string {

	indexAtStart := c.index
	if !readOnly {
		c.lastTokenCol = c.col
		c.whitespace = ""
	}

	word := ""
	for c.index <= c.maxIndex {

		charInt := c.code[c.index]
		char := string(charInt)

		if isNewLine(charInt) {
			if sameLine {
				if len(word) == 0 && !readOnly {
					c.throwAtLine("Unexpected new line")
				}
				break
			}
			c.index++
			if !readOnly {
				c.line++
				c.col = 0
				c.lastTokenCol = 0
				c.whitespace = ""
				lastChars := ""
				if len(c.result) > 1 {
					lastChars = c.result[len(c.result)-2:]
					if lastChars != "\n\n" {
						c.result += "\n"
					}
				}
			}
			if len(word) == 0 {
				continue
			}
			break
		}

		if isWhiteSpace(charInt) && len(word) == 0 {
			c.index++
			if !readOnly {
				c.col++
				c.lastTokenCol++
				c.whitespace += char
			}
			continue
		}

		if isVarNameChar(charInt) || (len(word) > 0 && isNumberChar(charInt)) || (len(word) == 0 && char == "#") {
			word += char
			c.index++
			if !readOnly {
				c.col++
			}
			continue
		}

		if len(word) == 0 {
			c.index++
			if !readOnly {
				c.col++
			}
			word = char
			// Comments
			if char == "/" {
				nextChar := c.readNextChar()
				if nextChar == "/" || nextChar == "*" {
					word += nextChar
					c.index++
					if !readOnly {
						c.col++
					}
					break
				}
			}
			// Operators
			i := sort.SearchStrings(operatorChars, char)
			if i < len(operatorChars) && operatorChars[i] == char {
				// If operator char
				nextChar := c.readNextChar()
				i = sort.SearchStrings(operatorChars, nextChar)
				if i < len(operatorChars) && operatorChars[i] == nextChar {
					word += nextChar
					c.index++
					if !readOnly {
						c.col++
					}
				}
				nextChar = c.readNextChar()
				if nextChar == "=" {
					word += nextChar
					c.index++
					if !readOnly {
						c.col++
					}
				}
			}
			break
		}

		break
	}

	if readOnly {
		c.index = indexAtStart
	}

	return word
}

func (c *Compile) getNextTokenSameLine() string {

	c.lastTokenCol = c.col

	word := ""
	for c.index <= c.maxIndex {

		charInt := c.code[c.index]
		char := string(charInt)

		if isNewLine(charInt) {
			if len(word) == 0 {
				c.throwAtLine("Unexpected new line")
			}
			break
		}

		if isWhiteSpace(charInt) && len(word) == 0 {
			c.index++
			c.col++
			c.lastTokenCol++
			continue
		}

		if isVarNameChar(charInt) || (len(word) == 0 && char == "#") {
			word += char
			c.index++
			c.col++
			continue
		}

		// If any special character
		if len(word) == 0 {
			c.index++
			c.col++
			return char
		}

		break
	}

	return word
}

func (c *Compile) getNextCharacterOnLine() string {

	c.lastTokenCol = c.col

	for c.index <= c.maxIndex {

		charInt := c.code[c.index]
		char := string(charInt)

		if isNewLine(charInt) {
			c.index++
			c.col = 0
			c.lastTokenCol = 0
			c.line++
			return char
		}

		c.index++
		c.col++

		if isWhiteSpace(charInt) {
			c.lastTokenCol++
			continue
		}

		return char
	}

	c.throwAtLine("Unexpected end of file")
	return ""
}

func (c *Compile) getNextValueToken() (string, string) {

	c.lastTokenCol = c.col

	vtype := ""
	word := ""
	prevChar := ""
	char := ""
	inStr := false
	hasDot := false
	endStrChar := ""

	for c.index <= c.maxIndex {

		prevChar = char
		charInt := c.code[c.index]
		char = string(charInt)

		if inStr {
			word += char
			c.index++
			c.col++
			c.lastTokenCol++
			if char == endStrChar {
				break
			}
			if isNewLine(charInt) {
				if prevChar != "\\" {
					// prevent new line
					c.throwAtLine("Unexpected new line")
				}
				c.line++
				c.col = 0
				c.lastTokenCol = 0
			}
			continue
		}

		if isNewLine(charInt) {
			if len(word) == 0 {
				c.throwAtLine("Unexpected new line")
			}
			break
		}

		if isWhiteSpace(charInt) && len(word) == 0 {
			c.index++
			c.col++
			c.lastTokenCol++
			continue
		}

		// Strings
		if len(word) == 0 && (char == "\"" || char == "'") {
			word += char
			endStrChar = char
			c.index++
			c.col++
			vtype = "str"
			inStr = true
			continue
		}

		if isVarNameChar(charInt) {
			if vtype != "" && vtype != "word" {
				c.throwAtLine("Unexpected char: " + char)
			}
			word += char
			c.index++
			c.col++
			vtype = "word"
			continue
		}

		isDot := char == "."
		if isNumberChar(charInt) || isDot {
			if vtype != "" && vtype != "num" {
				c.throwAtLine("Unexpected char: " + char)
			}
			if isDot {
				if hasDot || len(word) == 0 {
					c.throwAtLine("Unexpected char: " + char)
				}
				hasDot = true
			}
			word += char
			c.index++
			c.col++
			vtype = "num"
			continue
		}

		break
	}

	if len(word) == 0 {
		c.throwAtLine("Missing value")
	}

	// if ends with a dot
	if string(word[len(word)-1]) == "." {
		c.throwAtLine("Unexpected dot")
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

	c.throwAtLine("Unknown value: " + word)
	return "", ""
}

func (c *Compile) getNextType() *VarType {
	token := c.getNextToken(false, true)
	result := VarType{
		name: token,
	}
	i := sort.SearchStrings(basicTypes, token)
	if i < len(basicTypes) && basicTypes[i] == token {
		if token == "array" || token == "object" {
			c.expectToken("<")
			result.subtype = c.getNextType()
			c.expectToken(">")
		}
		if token == "func" {
			c.expectToken("(")
			currentIndex := c.index
			currentCol := c.col
			ntoken := c.getNextTokenSameLine()
			if ntoken != ")" {
				c.col = currentCol
				c.index = currentIndex
			}
			for ntoken != ")" {
				ptype := c.getNextType()
				result.paramTypes = append(result.paramTypes, ptype)
				ntoken = c.getNextTokenSameLine()
				if ntoken != "," && ntoken != ")" {
					c.throwAtLine("Unexpected token: " + ntoken)
				}
			}
			rtype := c.getNextType()
			result.returnType = rtype
		}
	} else {
		_, foundStruct := c.getStruct(token)
		_, foundClass := c.getStruct(token) // todo: fix
		if foundStruct {
			result.toft = "struct"
		} else if foundClass {
			result.toft = "class"
		} else {
			c.throwAtLine("Unknown type: " + token)
		}
		nchar := c.readNextChar()
		if nchar == "<" {
			c.expectToken("<")
			result.subtype = c.getNextType()
			c.expectToken(">")
		}
	}
	nchar := c.readNextChar()
	for nchar == "|" {
		c.expectToken("|")
		ptype := c.getNextTokenSameLine()
		if ptype == "null" {
			result.nullable = true
		} else if ptype == "undefined" {
			result.undefined = true
		} else {
			c.throwAtLine("Expected null or undefined")
		}
		nchar = c.readNextChar()
	}

	return &result
}

func (c *Compile) readNextChar() string {
	if c.index == c.maxIndex {
		return ""
	}
	return string(c.code[c.index])
}

func (c *Compile) expectToken(token string) {
	ntoken := c.getNextTokenSameLine()
	if ntoken != token {
		c.throwAtLine("Expected: " + token)
	}
}

func (c *Compile) handleNextWord() {

	token := c.getNextToken(false, false)

	if token == "#" {
		word := c.getNextToken(false, false)
		c.handleMacro(word)
		return
	}
	if token == "}" && scopeIndex > 0 {
		c.exitScope = true
		return
	}

	isDefine := token == "define"
	if isDefine {
		token = c.getNextTokenSameLine()
	}

	if token == "struct" || token == "local" {
		isLocal := token == "local"
		if isLocal {
			token = c.getNextTokenSameLine()
			if token != "struct" {
				c.throwAtLine("Unexpected token: " + token)
			}
		}
		c.handleStruct(isLocal, isDefine)
		return
	}

	if token == "return" {
		c.result += "return "
		rtype := c.assignValue()
		c.result += ";"
		scope := scopes[scopeIndex]
		if scope.returnType == nil {
			c.throwAtLine("Unexpected return statement")
		}
		if !scope.returnType.isCompatible(rtype) {
			c.throwTypeError(scope.returnType, rtype)
		}
		return
	}

	if token == "class" {
		c.handleClass(isDefine)
		return
	}

	if token == "include" {
		c.handleInclude()
		return
	}

	if len(token) == 0 {
		return
	}

	if token == "import" {
		c.handleImport()
		return
	}

	if token == "//" || token == "/*" {
		c.handleComment()
		return
	}

	// Check if type/struct/class
	if c.isValidType(token) {
		c.index -= len(token)
		c.col -= len(token)
		if c.col < 0 {
			c.col = 0
		}
		_type := c.getNextType()
		c.declareVariable(_type, isDefine)
		return
	}

	_, isVar := c.getVar(token)
	if isVar || token == "(" || token == "[" {
		c.index -= len(token)
		c.col -= len(token)
		vt := c.assignValue()
		if vt.assignable {
			// Check for = sign
			nextToken := c.getNextToken(true, false)
			if nextToken == "=" {
				nextToken = c.getNextToken(false, false)
				c.result += "="
				assignType := c.assignValue()
				if !vt.isCompatible(assignType) {
					c.throwTypeError(vt, assignType)
				}
			}
		} else {
			nextToken := c.getNextToken(true, false)
			if nextToken == "=" {
				c.throwAtLine("Cannot assign a value to this")
			}
		}
		nextChar := c.getNextToken(true, false)
		if nextChar == ";" {
			nextChar = c.getNextToken(false, false)
		}
		c.result += ";"
		// Todo: check for missing props if left is struct
		return
	}

	// Unknown
	if isVarNameSyntax([]byte(token)) {
		c.throwAtLine("Unknown variable: " + token)
	}

	c.throwAtLine("Unexpected token: " + token)
}

func (c *Compile) getStruct(name string) (*VarType, bool) {
	var sci = scopeIndex
	for sci >= 0 {
		scope := scopes[sci]
		realName, ok := scope.structs[name]
		if ok {
			result, ok := allStructs[realName]
			return result, ok
		}
		sci--
	}
	return nil, false
}

func (c *Compile) getClass(name string) (*Class, bool) {
	var sci = scopeIndex
	for sci >= 0 {
		scope := scopes[sci]
		realName, ok := scope.classes[name]
		if ok {
			result, ok := allClasses[realName]
			return &result, ok
		}
		sci--
	}
	return nil, false
}

func (c *Compile) getVar(name string) (*Var, bool) {
	var sci = scopeIndex
	for sci >= 0 {
		scope := scopes[sci]
		result, ok := scope.vars[name]
		if ok {
			return &result, ok
		}
		sci--
	}
	return nil, false
}

func (c *Compile) typeExists(name string) bool {
	var sci = scopeIndex
	for sci >= 0 {
		scope := scopes[sci]
		if scope.typeExists(name) {
			return true
		}
		sci--
	}
	return false
}

func (c *Compile) isValidType(token string) bool {
	i := sort.SearchStrings(basicTypes, token)
	if i < len(basicTypes) && basicTypes[i] == token {
		return true
	}
	return c.typeExists(token)
}

func (c *Compile) throwAtLine(msg string) {

	fmt.Print("\033[31m") // Color red
	fmt.Println(msg)
	fmt.Print("\033[0m") // Color reset
	fmt.Println("Line", c.line, "col", c.col, "in", c.name)
	fmt.Println(c.readLine(c.line))
	fmt.Print("\033[31m") // Color red
	i := 0
	mark := ""
	for i < c.lastTokenCol {
		mark += " "
		i++
	}
	mark += "^"
	fmt.Println(mark)
	fmt.Print("\033[0m") // Color reset

	os.Exit(1)
}

func (c *Compile) throw(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func (c *Compile) readLine(lineNr int) string {
	i := 0
	line := ""
	currentLine := 1
	for i <= c.maxIndex {

		charInt := c.code[i]
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
