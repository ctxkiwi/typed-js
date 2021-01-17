package main

import (
	"strconv"
)

var allStructs = map[string]*VarType{}
var structCount = 0

func createNewStruct() (string, *VarType) {
	name := "struct_" + strconv.Itoa(structCount)
	vt := VarType{}
	vt.isStruct = true
	vt.props = map[string]*Property{}
	allStructs[name] = &vt
	structCount++
	return name, &vt
}

type Property struct {
	varType  *VarType
	_default string
}
