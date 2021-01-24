package main

import (
	"sort"
)

func (fc *FileCompiler) skipFunction() *VarType {
	fc.expectToken("(")

	_type := VarType{
		name:       "func",
		paramTypes: []*VarType{},
	}

	ntoken := fc.getNextToken(true, false)
	if ntoken == ")" {
		ntoken = fc.getNextToken(false, false)
	}
	for ntoken != ")" {

		paramName := fc.getNextToken(false, false)
		ptype := fc.getNextType()
		ptype.paramName = paramName
		_type.paramTypes = append(_type.paramTypes, ptype)

		ntoken = fc.getNextTokenSameLine()
		if ntoken != "," && ntoken != ")" {
			fc.throwAtLine("Unexpected token: " + ntoken)
		}

	}
	_type.returnType = fc.getNextType()

	fc.expectToken("{")
	fc.skipScope("}")

	return &_type
}

func (fc *FileCompiler) skipScope(endChar string) {

	level := 0
	char := ""

	for fc.index <= fc.maxIndex {

		charInt := fc.code[fc.index]
		char = string(charInt)

		fc.index++
		fc.col++
		fc.lastTokenCol++

		if isNewLine(charInt) {
			fc.line++
		}

		if level == 0 && char == endChar {
			return
		}

		if char == "\"" || char == "'" {
			fc.skipString(char)
		} else if char == "(" {
			level++
		} else if char == ")" {
			level--
		} else if char == "{" {
			level++
		} else if char == "}" {
			level--
		}
	}

	fc.throwAtLine("Unexpected end of code, expected: " + endChar)
}

func (fc *FileCompiler) skipString(endStrChar string) {

	char := ""
	prevChar := ""

	for fc.index <= fc.maxIndex {

		prevChar = char

		charInt := fc.code[fc.index]
		char = string(charInt)

		if isNewLine(charInt) {
			fc.line++
		}

		fc.index++
		fc.col++
		fc.lastTokenCol++

		if char == endStrChar && prevChar != "\\" {
			return
		}
	}

	fc.throwAtLine("Unexpected end of code, expected: " + endStrChar)
}

func (fc *FileCompiler) skipValue() {

	// allowDot := true
	token := fc.getNextToken(false, true)
	if token == "(" {
		fc.skipScope(")")
	} else if token == "[" {
		fc.skipScope("]")
	} else if token == "{" {
		fc.skipScope("}")
	} else if token == "function" {
		fc.skipFunction()
	} else if token == "new" {
		fc.getNextToken(false, true) // skip classname
	} else if token == "\"" || token == "'" {
		fc.skipString(token)
	}

	for {
		nextChar := fc.readNextChar()
		if nextChar == "[" {
			nextChar = fc.getNextToken(false, true)
			fc.skipScope("]")
		} else if nextChar == "(" {
			nextChar = fc.getNextToken(false, true)
			fc.skipScope(")")
		} else if nextChar == "." {
			nextChar = fc.getNextToken(false, true) // dot
			nextChar = fc.getNextToken(false, true) // word
		} else {
			break
		}
	}

	// Handle operators
	nextToken := fc.getNextToken(true, true)
	i := sort.SearchStrings(operators, nextToken)
	if i < len(operators) && operators[i] == nextToken {
		nextToken = fc.getNextToken(false, false) // skip operator
		fc.skipValue()
	}

}
