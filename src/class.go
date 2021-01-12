
package main

import (
	"strconv"
)

var allClasses = map[string]Class{}
var classCount = 0

func createNewClass() (string, *Class) {
	name := strconv.Itoa(classCount)
	s := Class{}
	s.props = map[string]ClassProperty{}
	allClasses[name] = s
	return name, &s
}

type Class struct {
	isLocal bool
	props map[string]ClassProperty
}

type ClassProperty struct {
	_type string
	_typeOfType string
	_default string
}