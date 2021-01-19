package main

import (
	"strconv"
)

var allClasses = map[string]Class{}
var classCount = 0

func createNewClass() (string, *Class) {
	name := "class_" + strconv.Itoa(classCount)
	s := Class{}
	s.props = map[string]*Property{}
	allClasses[name] = s
	return name, &s
}

type Class struct {
	isLocal bool
	props   map[string]*Property
}

// type ClassProperty struct {
// 	_type string
// 	_typeOfType string
// 	_default string
// }
