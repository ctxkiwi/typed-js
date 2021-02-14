package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var allFileCompilers = map[string]*FileCompiler{}

type Import struct {
	fileCompiler *FileCompiler
	internalName string
	externalName string
}

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
	fromBytes := []byte(from)
	isAsset := false
	fullpath := ""
	if string(fromBytes[0]) == "@" {
		isAsset = true
		fromBytes = fromBytes[1:]
		fullpath = from          // with @
		from = string(fromBytes) // without @
	} else {

		dir := filepath.Dir(fc.name) + "/"

		fp, err := filepath.Abs(dir + from)
		if err != nil {
			fc.throwAtLine("Invalid filepath: " + from)
		}
		fullpath = fp
	}

	if fc.compiler.readTypes {

		nfc, exists := allFileCompilers[fullpath]
		if !exists {

			code := []byte{}
			if isAsset {

				_code, err := Asset("src/" + from)
				if err != nil {
					fc.throwAtLine("Cannot find import: " + from)
				}
				code = []byte(_code)

			} else {
				if !fileExists(fullpath) {
					fc.throwAtLine("File not found: " + from + " (" + fullpath + ")")
				}

				_code, err := ioutil.ReadFile(fullpath)
				if err != nil {
					fmt.Println("Cant read file: " + fullpath)
					os.Exit(1)
				}
				code = _code

				defaultCode := []byte("import \"@core/_imports\"\n\n")
				code = append(defaultCode, code...)
			}

			nfc = fc.compiler.createNewFileCompiler(fullpath, code)
			nfc.createNewScope()

			scope := fc.getScope()
			nscope := nfc.getScope()

			for varName, varType := range scope.vars {
				nscope.vars[varName] = varType
			}

			allFileCompilers[fullpath] = nfc

			nfc.compile()
		}
		// Check if imports exist
		for typeName, typeAlias := range imports {

			_, exists := nfc.exports[typeName]
			if !exists {
				fc.throwAtLine("Class/Struct " + typeName + " not found or not exported in: " + from)
			}

			if fc.typeExists(typeAlias) {
				fc.throwAtLine("Cannot import, name already in use: " + typeAlias)
			}

			scope := fc.scopes[fc.scopeIndex]
			_type, _ := nfc.getType(typeName)
			scope.types[typeAlias] = _type.name

			fc.imports = append(fc.imports, &Import{
				fileCompiler: nfc,
				internalName: typeAlias,
				externalName: typeName,
			})
		}

		//
		return
	} else {
		nfc, _ := allFileCompilers[fullpath]
		if !nfc.compiled {
			nfc.line = 1
			nfc.index = 0
			nfc.col = 0
			nfc.compile()

			nfc.compiled = true

			resultOrder = append(resultOrder, nfc)
		}

		scope := fc.getScope()
		nscope := nfc.getScope()

		for varName, varType := range nscope.vars {
			scope.vars[varName] = varType
		}
	}

}
