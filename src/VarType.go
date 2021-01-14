package main

import "sort"

type VarType struct {
	name       string
	toft       string
	subtype    *VarType
	nullable   bool
	undefined  bool
	paramTypes []*VarType
	returnType *VarType
}

func (c *Compile) getTypeOfType(_type string) (string, bool) {
	i := sort.SearchStrings(basicTypes, _type)
	if i < len(basicTypes) && basicTypes[i] == _type {
		return "basic", true
	}
	_, ok := c.getStruct(_type)
	if ok {
		return "struct", true
	}
	_, ok = c.getClass(_type)
	if ok {
		return "class", true
	}
	return "", false
}
