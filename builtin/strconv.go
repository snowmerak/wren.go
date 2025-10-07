package builtin

//go:generate go run ../wrengen -dir .

import (
	"fmt"
	"strconv"

	wrengo "github.com/snowmerak/gwen"
)

// StrConv provides built-in string conversion utilities
//
//wren:bind module=strconv
type StrConv struct{}

// Atoi converts string to integer
//
//wren:bind name=atoi(_) static
func (sc *StrConv) Atoi(vm *wrengo.WrenVM) error {
	str := vm.GetSlotString(1)
	if num, err := strconv.Atoi(str); err == nil {
		vm.SetSlotDouble(0, float64(num))
	} else {
		vm.SetSlotString(0, fmt.Sprintf("Error: %s", err.Error()))
	}
	return nil
}

// ParseFloat converts string to float
//
//wren:bind name=parseFloat(_) static
func (sc *StrConv) ParseFloat(vm *wrengo.WrenVM) error {
	str := vm.GetSlotString(1)
	if num, err := strconv.ParseFloat(str, 64); err == nil {
		vm.SetSlotDouble(0, num)
	} else {
		vm.SetSlotString(0, fmt.Sprintf("Error: %s", err.Error()))
	}
	return nil
}

// Itoa converts integer to string
//
//wren:bind name=itoa(_) static
func (sc *StrConv) Itoa(vm *wrengo.WrenVM) error {
	num := vm.GetSlotDouble(1)
	result := strconv.Itoa(int(num))
	vm.SetSlotString(0, result)
	return nil
}

// FormatFloat converts float to string with precision
//
//wren:bind name=formatFloat(_,_) static
func (sc *StrConv) FormatFloat(vm *wrengo.WrenVM) error {
	num := vm.GetSlotDouble(1)
	precision := int(vm.GetSlotDouble(2))
	result := strconv.FormatFloat(num, 'f', precision, 64)
	vm.SetSlotString(0, result)
	return nil
}

// ParseBool converts string to boolean ("true"/"false", "1"/"0", etc.)
//
//wren:bind name=parseBool(_) static
func (sc *StrConv) ParseBool(vm *wrengo.WrenVM) error {
	str := vm.GetSlotString(1)
	if b, err := strconv.ParseBool(str); err == nil {
		vm.SetSlotBool(0, b)
	} else {
		vm.SetSlotString(0, fmt.Sprintf("Error: %s", err.Error()))
	}
	return nil
}

// FormatBool converts boolean to string
//
//wren:bind name=formatBool(_) static
func (sc *StrConv) FormatBool(vm *wrengo.WrenVM) error {
	b := vm.GetSlotBool(1)
	result := strconv.FormatBool(b)
	vm.SetSlotString(0, result)
	return nil
}
