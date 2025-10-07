package builtin

//go:generate go run ../wrengen -dir .

import (
	"math"

	wrengo "github.com/snowmerak/wren.go"
)

// Math provides built-in mathematical utilities
//
//wren:bind module=math
type Math struct{}

// Sqrt returns the square root of x
//
//wren:bind name=sqrt(_) static
func (m *Math) Sqrt(vm *wrengo.WrenVM) error {
	x := vm.GetSlotDouble(1)
	if x < 0 {
		vm.SetSlotString(0, "Error: sqrt of negative number")
		return nil
	}
	result := math.Sqrt(x)
	vm.SetSlotDouble(0, result)
	return nil
}

// Pow returns x^y, the base-x exponential of y
//
//wren:bind name=pow(_,_) static
func (m *Math) Pow(vm *wrengo.WrenVM) error {
	x := vm.GetSlotDouble(1)
	y := vm.GetSlotDouble(2)
	result := math.Pow(x, y)
	vm.SetSlotDouble(0, result)
	return nil
}

// Sin returns the sine of the radian argument x
//
//wren:bind name=sin(_) static
func (m *Math) Sin(vm *wrengo.WrenVM) error {
	x := vm.GetSlotDouble(1)
	result := math.Sin(x)
	vm.SetSlotDouble(0, result)
	return nil
}

// Cos returns the cosine of the radian argument x
//
//wren:bind name=cos(_) static
func (m *Math) Cos(vm *wrengo.WrenVM) error {
	x := vm.GetSlotDouble(1)
	result := math.Cos(x)
	vm.SetSlotDouble(0, result)
	return nil
}

// Abs returns the absolute value of x
//
//wren:bind name=abs(_) static
func (m *Math) Abs(vm *wrengo.WrenVM) error {
	x := vm.GetSlotDouble(1)
	result := math.Abs(x)
	vm.SetSlotDouble(0, result)
	return nil
}

// Max returns the larger of x or y
//
//wren:bind name=max(_,_) static
func (m *Math) Max(vm *wrengo.WrenVM) error {
	x := vm.GetSlotDouble(1)
	y := vm.GetSlotDouble(2)
	result := math.Max(x, y)
	vm.SetSlotDouble(0, result)
	return nil
}

// Min returns the smaller of x or y
//
//wren:bind name=min(_,_) static
func (m *Math) Min(vm *wrengo.WrenVM) error {
	x := vm.GetSlotDouble(1)
	y := vm.GetSlotDouble(2)
	result := math.Min(x, y)
	vm.SetSlotDouble(0, result)
	return nil
}

// Pi returns the mathematical constant π
//
//wren:bind name=pi static
func (m *Math) Pi(vm *wrengo.WrenVM) error {
	vm.SetSlotDouble(0, math.Pi)
	return nil
}