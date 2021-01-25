package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func main() {
	// Get file from argv
	if len(os.Args) <= 1 {
		fmt.Println("Missing input file param")
		os.Exit(1)
	}
	if len(os.Args) <= 2 {
		fmt.Println("Missing output file param")
		os.Exit(1)
	}

	fn := os.Args[1]
	if !strings.HasSuffix(fn, ".tjs") {
		fmt.Println("File must be a .tjs file")
		os.Exit(1)
	}

	inputFilepath, err := filepath.Abs(fn)
	if err != nil {
		fmt.Println("Cant generate absolute filepath for input file")
		os.Exit(1)
	}

	outfn := os.Args[2]
	outputFilepath, err := filepath.Abs(outfn)
	if err != nil {
		fmt.Println("Cant generate absolute filepath for output file")
		os.Exit(1)
	}

	// Check if input file exists
	code, err := ioutil.ReadFile(inputFilepath)
	if err != nil {
		fmt.Println("Cant read file: " + inputFilepath)
		os.Exit(1)
	}

	// Sort global arrays
	sort.Strings(basicTypes)
	sort.Strings(basicValues)
	sort.Strings(operators)
	sort.Strings(operatorChars)
	sort.Strings(structsEqualToClass)

	// Compile
	compiler := Compiler{
		readTypes: true,
	}

	start := time.Now()
	data, err := Asset("src/core/types.tjs")
	if err != nil {
		fmt.Println("Missing assets: @core/types")
		os.Exit(1)
	}
	jscode := compiler.compileCode("src/core/types.tjs", []byte(data))

	data, _ = Asset("src/core/globals.tjs")
	jscode += compiler.compileCode("src/core/globals.tjs", []byte(data))

	jscode += "\n(function(){\n\n"
	jscode += compiler.compileCode(inputFilepath, code)
	jscode += "\n})();\n"
	elapsed := time.Since(start)

	// Result
	ioutil.WriteFile(outputFilepath, []byte(jscode), 0644)
	log.Printf("Compiled in %s", elapsed)
}
