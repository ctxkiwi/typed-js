
package main

import (
	"strconv"
)

var allStructs = map[string]Struct{}
var structCount = 0

func createNewStruct() (string, *Struct) {
	name := strconv.Itoa(structCount)
	s := Struct{}
	s.props = map[string]Property{}
	allStructs[name] = s
	return name, &s
}

type Struct struct {
	isLocal bool
	props map[string]Property
}

type Property struct {
	_type string
	_typeOfType string
	_default string
}