
package main;

import(
	"os"
	"fmt"
	"path/filepath"
	"strings"
	"io/ioutil"
	"sort"
)

func main(){
	// Get file from argv
	if len(os.Args) <= 1 {
		fmt.Println("Missing input file param")
		os.Exit(1)
	}
	if len(os.Args) <= 2 {
		fmt.Println("Missing output file param")
		os.Exit(1)
	}

	fn := os.Args[1];
	if !strings.HasSuffix(fn, ".tjs") {
		fmt.Println("File must be a .tjs file")
		os.Exit(1)
	}

	inputFilepath, err := filepath.Abs(fn)
	if err != nil {
		fmt.Println("Cant generate absolute filepath for input file")
		os.Exit(1)
	}

	outfn := os.Args[2];
	outputFilepath, err := filepath.Abs(outfn)
	if err != nil {
		fmt.Println("Cant generate absolute filepath for output file")
		os.Exit(1)
	}

	sort.Strings(basicTypes)

	code := compileFile(inputFilepath);

	ioutil.WriteFile(outputFilepath, []byte(code), 0644)

}
