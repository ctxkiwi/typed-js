package main

import (
	"sort"
	"strings"
)

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

func (t *VarType) isCompatible(at *VarType) bool {
	if at == nil {
		return false
	}
	if t.name != at.name {
		nameLower := strings.ToLower(t.name)
		i := sort.SearchStrings(structsEqualToClass, nameLower)
		if !(i < len(structsEqualToClass) && structsEqualToClass[i] == nameLower && nameLower == strings.ToLower(at.name)) {
			return false
		}
	}
	if at.nullable && !t.nullable {
		return false
	}
	if at.undefined && !t.undefined {
		return false
	}
	// if one is nil && one is not nil
	if (t.subtype == nil || at.subtype == nil) && (t.subtype != nil || at.subtype != nil) {
		return false
	}
	if len(t.paramTypes) != len(at.paramTypes) {
		return false
	}
	for i, pt := range t.paramTypes {
		apt := at.paramTypes[i]
		if !pt.isCompatible(apt) {
			return false
		}
	}
	return true
}

func (t *VarType) displayName() string {

	if t == nil {
		return "???"
	}

	result := t.name

	if len(t.paramTypes) > 0 {
		result += "("
		for i, p := range t.paramTypes {
			if i > 0 {
				result += ","
			}
			result += p.displayName()
		}
		result += ")"
	}
	if t.returnType != nil {
		result += "<" + t.returnType.displayName() + ">"
	}
	if t.subtype != nil {
		result += "<" + t.subtype.displayName() + ">"
	}
	if t.nullable {
		result += "|null"
	}
	if t.undefined {
		result += "|undefined"
	}
	return result
}

func (c *Compile) throwTypeError(t *VarType, at *VarType) {
	c.throwAtLine("Types not compatible: " + t.displayName() + " <-> " + at.displayName() + "")
}
