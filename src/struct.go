
package main

import (
	"strconv"
)

var allStructs = map[string]Struct{}
var structCount = 0

func createNewStruct() (string, Struct) {
	name := strconv.Itoa(structCount)
	s := Struct{}
	s.vars = map[string]Property{}
	allStructs[name] = s
	return name, s
}

type Struct struct {
	isLocal bool
	vars map[string]Property
}

type Property struct {
	_type string
	_default string
}