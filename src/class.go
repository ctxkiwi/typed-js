package main

import (
	"strconv"
)

var allClasses = map[string]*VarType{}
var classCount = 0

func createNewClass() (string, *VarType) {
	name := "class_" + strconv.Itoa(classCount)
	s := VarType{
		isClass: true,
	}
	s.name = name
	s.props = map[string]*Property{}
	allClasses[name] = &s
	return name, &s
}

// type Class struct {
// 	isLocal bool
// 	name    string
// 	props   map[string]*Property
// }

// type ClassProperty struct {
// 	_type string
// 	_typeOfType string
// 	_default string
// }
