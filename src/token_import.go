package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var allFileCompilers = map[string]*FileCompiler{}

func (fc *FileCompiler) handleImport() {

	imports := map[string]string{}

	token := fc.getNextToken(false, true)
	if token != "'" && token != "\"" {
		for token != "from" {
			typeName := ""
			alias := ""
			imp := strings.Split(token, ":")
			typeName = imp[0]
			if len(imp) > 1 {
				alias = imp[1]
			} else {
				alias = typeName
			}

			imports[typeName] = alias

			token = fc.getNextToken(false, true)
			if token == "" {
				fc.throwAtLine("Missing \"from\" filepath")
			}
		}
		token = fc.getNextToken(false, true)
	}
	if token != "'" && token != "\"" {
		fc.throwAtLine("Expected a string containing a filepath")
	}

	from := strings.Trim(fc.skipString(token), token) + ".tjs"
	dir := filepath.Dir(fc.name) + "/"

	filepath, err := filepath.Abs(dir + from)
	if err != nil {
		fc.throwAtLine("Invalid filepath: " + from)
	}

	if fc.compiler.readTypes {
		if !fileExists(filepath) {
			fc.throwAtLine("File not found: " + from + " (" + filepath + ")")
		}

		nfc, exists := allFileCompilers[filepath]
		if !exists {

			code, err := ioutil.ReadFile(filepath)
			if err != nil {
				fmt.Println("Cant read file: " + filepath)
				os.Exit(1)
			}
			nfc = &FileCompiler{
				name: filepath,

				compiler: fc.compiler,

				scopes:     []*Scope{},
				scopeIndex: -1,

				index:     0,
				col:       0,
				maxIndex:  len(code) - 1,
				line:      1,
				exitScope: false,

				compiled: false,
				code:     code,
				result:   "",
			}

			nfc.createNewScope()

			allFileCompilers[filepath] = nfc

			nfc.compile()
		}
		// Check if imports exist
		for typeName, typeAlias := range imports {

			if !nfc.typeExists(typeName) {
				fc.throwAtLine("Class/Struct " + typeName + " not found in: " + from)
			}

			if fc.typeExists(typeAlias) {
				fc.throwAtLine("Cannot import, name already in use: " + typeAlias)
			}

			scope := fc.scopes[fc.scopeIndex]
			_type, _ := nfc.getType(typeName)
			scope.types[typeAlias] = _type.name
		}

		//
		return
	} else {
		nfc, _ := allFileCompilers[filepath]
		if !nfc.compiled {
			nfc.line = 1
			nfc.index = 0
			nfc.col = 0
			nfc.compile()

			fc.result += nfc.result
		}
	}

}
