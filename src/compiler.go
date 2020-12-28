
package main;

import(
	"fmt"
	"os"
	"io/ioutil"
)

func compileFile(file string) string{

	code, err := ioutil.ReadFile(file);
	if err != nil {
		fmt.Println("Cant read file: " + file)
		os.Exit(1)
	}
	return compileCode(code)
}

func compileCode(code []byte) string{

	return ""
}