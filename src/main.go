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

	sort.Strings(basicTypes)
	sort.Strings(basicValues)

	start := time.Now()
	code := compileFile(inputFilepath)
	elapsed := time.Since(start)

	ioutil.WriteFile(outputFilepath, []byte(code), 0644)

	log.Printf("Compiled in %s", elapsed)
}
