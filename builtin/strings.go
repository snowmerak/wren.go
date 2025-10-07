package builtin

//go:generate go run ../wrengen -dir .

import (
	"strings"

	wrengo "github.com/snowmerak/gwen"
)

// Strings provides built-in string manipulation utilities
//
//wren:bind module=strings
type Strings struct{}

// ToUpper converts string to uppercase
//
//wren:bind name=upper(_) static
func (s *Strings) ToUpper(vm *wrengo.WrenVM) error {
	str := vm.GetSlotString(1)
	result := strings.ToUpper(str)
	vm.SetSlotString(0, result)
	return nil
}

// ToLower converts string to lowercase
//
//wren:bind name=lower(_) static
func (s *Strings) ToLower(vm *wrengo.WrenVM) error {
	str := vm.GetSlotString(1)
	result := strings.ToLower(str)
	vm.SetSlotString(0, result)
	return nil
}

// Trim removes leading and trailing whitespace
//
//wren:bind name=trim(_) static
func (s *Strings) Trim(vm *wrengo.WrenVM) error {
	str := vm.GetSlotString(1)
	result := strings.TrimSpace(str)
	vm.SetSlotString(0, result)
	return nil
}

// Contains checks if string contains substring
//
//wren:bind name=contains(_,_) static
func (s *Strings) Contains(vm *wrengo.WrenVM) error {
	str := vm.GetSlotString(1)
	substr := vm.GetSlotString(2)
	result := strings.Contains(str, substr)
	vm.SetSlotBool(0, result)
	return nil
}

// Split splits string by delimiter
//
//wren:bind name=split(_,_) static
func (s *Strings) Split(vm *wrengo.WrenVM) error {
	str := vm.GetSlotString(1)
	delimiter := vm.GetSlotString(2)
	parts := strings.Split(str, delimiter)

	// Return as a comma-separated string for simplicity
	// In a real implementation, you might want to return a Wren List
	result := strings.Join(parts, ",")
	vm.SetSlotString(0, result)
	return nil
}

// Join joins array elements with delimiter
//
//wren:bind name=join(_,_) static
func (s *Strings) Join(vm *wrengo.WrenVM) error {
	elements := vm.GetSlotString(1) // Comma-separated elements
	delimiter := vm.GetSlotString(2)

	parts := strings.Split(elements, ",")
	result := strings.Join(parts, delimiter)
	vm.SetSlotString(0, result)
	return nil
}
