package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type VarType struct {
	name       string
	toft       string
	subtype    *VarType
	props      map[string]*Property
	nullable   bool
	undefined  bool
	paramTypes []*VarType
	returnType *VarType
	assignable bool
	paramName  string
	//
	isStruct bool
	isClass  bool
	isLocal  bool
}

func (fc *FileCompiler) getTypeOfType(_type string) (string, bool) {
	i := sort.SearchStrings(basicTypes, _type)
	if i < len(basicTypes) && basicTypes[i] == _type {
		return "basic", true
	}
	_, ok := fc.getStruct(_type)
	if ok {
		return "struct", true
	}
	_, ok = fc.getClass(_type)
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
		if i < len(structsEqualToClass) && structsEqualToClass[i] == nameLower && nameLower == strings.ToLower(at.name) {
			return true
		}
		// Check props
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
	if t.subtype != nil && !t.subtype.isCompatible(at.subtype) {
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

func (fc *FileCompiler) throwTypeError(t *VarType, at *VarType) {
fc.throwAtLine("Types not compatible: " + t.displayName() + " <-> " + at.displayName() + "")
}

func createType(name string) *VarType {
	result := VarType{}
	switch name {
	case "string":
		result.name = "string"
	case "bool":
		result.name = "bool"
	default:
		fmt.Println("Unknown type: " + name)
		os.Exit(1)
	}
	return &result
}
