package main

import (
	"strconv"
)

var allTypes = map[string]*VarType{}
var typeCount = 0

func createNewType(isClass bool) (string, *VarType) {
	name := "type_" + strconv.Itoa(typeCount)
	s := VarType{
		isClass:  isClass,
		isStruct: !isClass,
	}
	s.name = name
	s.props = map[string]*Property{}
	allTypes[name] = &s
	typeCount++
	return name, &s
}

type Property struct {
	varType  *VarType
	_default string
}
