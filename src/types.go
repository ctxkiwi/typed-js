package main

import (
	"strconv"
)

var allTypes = map[string]*VarType{}

func createNewType(isClass bool, prefix string) (string, *VarType) {
	name := ""
	typeCount := 1
	for {
		name = prefix + "_" + strconv.Itoa(typeCount)
		_, exists := allTypes[name]
		if !exists {
			break
		}
		typeCount++
	}
	s := VarType{
		isClass:  isClass,
		isStruct: !isClass,
	}
	s.name = name
	s.props = map[string]*Property{}
	allTypes[name] = &s
	return name, &s
}

type Property struct {
	varType  *VarType
	_default string
}
