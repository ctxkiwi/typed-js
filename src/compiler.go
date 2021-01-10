
package main;

import(
	"fmt"
	"os"
	"io/ioutil"
)

type Compile struct {
	name string
	index int
	maxIndex int
	line int
	code []byte
	result string
}

func compileFile(file string) string{

	code, err := ioutil.ReadFile(file);
	if err != nil {
		fmt.Println("Cant read file: " + file)
		os.Exit(1)
	}

	c := Compile {
		name: file,
		index: 0,
		maxIndex: len(code) - 1,
		line: 1,
		code: code,
		result: "",
	}

	return c.compile()
}

func (c Compile) compile() string{

	for c.index <= c.maxIndex {
		c.handleNextWord();
	}

	return ""
}

func (c Compile) handleNextWord() {

	word := ""
	for c.index <= c.maxIndex {

		charInt := c.code[c.index];
		char := string(charInt)

		if isNewLine(charInt) {
			c.line++
		}

		if isWhiteSpace(charInt) && len(word) == 0 {
			c.result += char
			c.index++
			continue
		}

		if isAlphaChar(charInt) || (len(word) == 0 && char == "#") {
			word += char
			c.index++
			continue
		}

		if len(word) == 0 {
			c.throwAtLine("Unexpect token: " + char)
		}

		break;
	}
	c.throwAtLine("Test")
}

func (c Compile) throwAtLine (msg string) {

	fmt.Print("\033[31m") // Color red
	fmt.Println(msg)
	fmt.Print("\033[0m") // Color reset
	fmt.Println("Line", c.line, "in", c.name)
	fmt.Println("Line:", c.readLine(c.line))

	os.Exit(1)
}

func (c Compile) throw (msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func (c Compile) readLine(lineNr int) string {
	i := 0
	line := ""
	currentLine := 1
	for i <= c.maxIndex {

		charInt := c.code[i];
		char := string(charInt)

		if currentLine == lineNr {
			line += char
		}

		if isNewLine(charInt) {
			currentLine++
			if(currentLine > lineNr) {
				break
			}
		}

		i++;
	}
	return line
}