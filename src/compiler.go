package main

import (
	"strconv"
	"strings"
)

var basicTypes = []string{"bool", "number", "string", "array", "object", "func", "any", "T", "void"}
var structsEqualToClass = []string{"number", "string", "bool", "array", "object"}
var basicValues = []string{"true", "false", "undefined", "null", "[]", "{}"}
var operators = []string{"+", "-", "*", "/", "==", "===", "<", ">", "<=", ">=", "&&", "||", "++", "--"}
var operatorChars = []string{"+", "-", "*", "/", "=", "<", ">", "&", "|"}

var exportCounter = 0

var resultOrder = []*FileCompiler{}

type Compiler struct {
	readTypes bool
}

func (c *Compiler) compileCode(name string, code []byte) string {

	defaultCode := []byte("import \"@core/_imports\"\n\n")
	code = append(defaultCode, code...)

	fc := c.createNewFileCompiler(name, code)
	fc.createNewScope()

	// fmt.Println("Read only mode on")
	c.readTypes = true
	fc.compile()
	// fmt.Println("Read only mode off")
	fc.reset()
	c.readTypes = false

	fc.compile()

	resultOrder = append(resultOrder, fc)

	result := ""
	for _, ifc := range resultOrder {

		ifc.result = strings.TrimSpace(ifc.result)
		if len(ifc.result) == 0 {
			continue;
		}

		// Exports
		codeBefore := ""
		codeBefore += "\n    var " + ifc.exportVarName + " = "

		// Add imports / exports to result
		codeBefore += "(function("
		count := 0
		for _, imp := range ifc.imports {
			if count > 0 {
				codeBefore += ", "
			}
			codeBefore += imp.internalName
			count++
		}
		codeBefore += "){\n        "

		// Exports
		codeAfter := "\n"
		codeAfter += "        return {\n"
		count = 0
		for externName, internName := range ifc.exports {
			if count > 0 {
				codeAfter += ",\n"
			}
			codeAfter += "            "
			codeAfter += externName
			codeAfter += ": "
			codeAfter += internName
			count++
		}
		codeAfter += "\n        };\n    "

		// Imports input
		codeAfter += "})("
		count = 0
		for _, imp := range ifc.imports {
			if count > 0 {
				codeBefore += ", "
			}
			codeAfter += imp.fileCompiler.exportVarName
			codeAfter += "." + imp.externalName
			count++
		}
		codeAfter += ");\n"

		result += codeBefore + ifc.result + codeAfter
	}

	return result
}

func (c *Compiler) createNewFileCompiler(name string, code []byte) *FileCompiler {

	fc := FileCompiler{
		name: name,

		compiler: c,

		scopes:     []*Scope{},
		scopeIndex: -1,

		index:     0,
		col:       0,
		maxIndex:  len(code) - 1,
		line:      1,
		exitScope: false,

		imports: []*Import{},
		exports: map[string]string{},

		compiled: false,
		code:     code,
		result:   "",
	}

	fc.exportVarName = "export_" + strconv.Itoa(exportCounter)
	exportCounter++

	return &fc
}

