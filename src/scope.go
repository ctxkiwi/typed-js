
package main

type Var struct {
	_typeOfType string // basic,struct,class
	_type string // string bool int MyStruct MyClass etc...
	nullable bool
}

type Scope struct {
	structs map[string]string
	classes map[string]string
	vars map[string]Var
}

func (s *Scope) typeExists (name string) bool {
	_, ok := s.structs[name]
	if ok { return true }
	_, ok = s.classes[name]
	if ok { return true }
	return false
}

func (s *Scope) hasStruct (name string) bool {
	for _, str := range s.structs {
        if str == name {
            return true
        }
    }
    return false
}

func (s *Scope) hasClass (name string) bool {
	for _, str := range s.classes {
        if str == name {
            return true
        }
    }
    return false
}