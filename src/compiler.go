package main

import (
	"fmt"
	"os"
	"sort"
)

var basicTypes = []string{"bool", "int", "float", "string", "array", "object", "func"}
var basicValues = []string{"true", "false", "undefined", "null", "[]", "{}"}

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

	code   []byte
	result string
}

func compileCode(name string, code []byte) string {

	c := Compile{
		name: name,

		index:    0,
		maxIndex: len(code) - 1,
		line:     1,

		code:   code,
		result: "",
	}

	if len(scopes) == 0 {
		c.createNewScope()
	}

	return c.compile()
}

func (c *Compile) compile() string {

	for c.index <= c.maxIndex {
		c.handleNextWord()
	}

	return c.result
}

func (c *Compile) createNewScope() {
	s := Scope{
		structs: map[string]string{},
		classes: map[string]string{},
	}
	scopes = append(scopes, &s)
	scopeIndex++
}

func (c *Compile) getNextToken() string {

	c.lastTokenCol = c.col
	c.whitespace = ""

	word := ""
	for c.index <= c.maxIndex {

		charInt := c.code[c.index]
		char := string(charInt)

		if isNewLine(charInt) {
			c.index++
			c.line++
			c.col = 0
			c.lastTokenCol = 0
			c.whitespace = "\n"
			if len(word) == 0 {
				continue
			}
			break
		}

		if isWhiteSpace(charInt) && len(word) == 0 {
			c.index++
			c.col++
			c.lastTokenCol++
			c.whitespace += char
			continue
		}

		if isVarNameChar(charInt) || (len(word) == 0 && char == "#") {
			word += char
			c.index++
			c.col++
			continue
		}

		if len(word) == 0 {
			c.index++
			c.col++
			return char
		}

		break
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

func (c *Compile) getNextType() string {
	result := ""
	token := c.getNextTokenSameLine()
	i := sort.SearchStrings(basicTypes, token)
	if i < len(basicTypes) && basicTypes[i] == token {
		result += token
		if token == "array" || token == "object" {
			c.expectToken("<")
			result += ":"
			subtype := c.getNextType()
			result += subtype
		}
		if token == "func" {
			c.expectToken("(")
			result += ":"
			currentIndex := c.index
			ntoken := c.getNextToken()
			c.index = currentIndex
			for ntoken != ")" {
				ptype := c.getNextType()
				result += "|" + ptype
				ntoken = c.getNextToken()
				if ntoken != "," && ntoken != ")" {
					c.throwAtLine("Unexpected token: " + ntoken)
				}
			}
			result += ":"
			result += c.getNextType()
		}
		return result
	}
	if !c.typeExists(token) {
		c.throwAtLine("Unknown type: " + token)
	}
	return token
}

func (c *Compile) expectToken(token string) {
	ntoken := c.getNextTokenSameLine()
	if ntoken != token {
		c.throwAtLine("Expected: " + token)
	}
}

func (c *Compile) handleNextWord() {

	token := c.getNextToken()

	if token == "#" {
		word := c.getNextToken()
		c.handleMacro(word)
		return
	}

	if token == "struct" || token == "local" {
		c.handleStruct(token == "local")
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

	if token == "/" {
		c.handleComment()
		return
	}

	_typeOfType, ok := c.getTypeOfType(token)
	if ok {
		c.declareVariable(token, _typeOfType)
		return
	}

	_, ok = c.getVar(token)
	if ok {
		c.throwAtLine("Variables not ready yet")
	}

	// Unknown
	c.col = c.lastTokenCol

	if isVarNameSyntax([]byte(token)) {
		c.throwAtLine("Unknown variable/function/struct: " + token)
	}

	c.throwAtLine("Unknown token: " + token)
}

func (c *Compile) getTypeOfType(_type string) (string, bool) {
	i := sort.SearchStrings(basicTypes, _type)
	if i < len(basicTypes) && basicTypes[i] == _type {
		return "basic", true
	}
	_, ok := c.getStruct(_type)
	if ok {
		return "struct", true
	}
	_, ok = c.getClass(_type)
	if ok {
		return "class", true
	}
	return "", false
}

func (c *Compile) checkVarNameSyntax(name []byte) {
	if !isVarNameSyntax(name) {
		c.col = c.lastTokenCol
		c.throwAtLine("Invalid variable name: " + string(name))
	}
}

func (c *Compile) getStruct(name string) (*Struct, bool) {
	var sci = scopeIndex
	for sci >= 0 {
		scope := scopes[sci]
		realName, ok := scope.structs[name]
		if ok {
			result, ok := allStructs[realName]
			return &result, ok
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

func (c *Compile) throwAtLine(msg string) {

	fmt.Print("\033[31m") // Color red
	fmt.Println(msg)
	fmt.Print("\033[0m") // Color reset
	fmt.Println("Line", c.line, "col", c.col, "in", c.name)
	fmt.Println(c.readLine(c.line))
	fmt.Print("\033[31m") // Color red
	i := 0
	mark := ""
	for i < c.col {
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
