package main

import (
	"path/filepath"
	"strings"
)

var allImports map[string]*FileCompiler

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

	filepath, _ := filepath.Abs(dir + from)

	if !fileExists(filepath) {
		fc.throwAtLine("File not found: " + from + " (" + filepath + ")")
	}

	if fc.compiler.readTypes {
	}

	fc.throwAtLine("Import feature ready yet")
}
