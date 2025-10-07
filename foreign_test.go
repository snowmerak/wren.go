package wrengo_test

import (
	"testing"

	"github.com/snowmerak/gwen"
)

func TestForeignMethod(t *testing.T) {
	// Register a simple foreign method
	wrengo.RegisterForeignMethod("main", "Math", true, "add(_,_)", func(vm *wrengo.WrenVM) {
		a := vm.GetSlotDouble(1)
		b := vm.GetSlotDouble(2)
		result := a + b
		vm.SetSlotDouble(0, result)
	})

	vm := wrengo.NewVMWithForeign()
	defer vm.Free()

	code := `
class Math {
  foreign static add(a, b)
}

System.print(Math.add(3, 5))
`

	result, err := vm.Interpret("main", code)
	if err != nil {
		t.Fatalf("Interpret error: %v", err)
	}

	if result != wrengo.ResultSuccess {
		t.Fatalf("Expected ResultSuccess, got %v", result)
	}
}

func TestForeignMethodString(t *testing.T) {
	wrengo.RegisterForeignMethod("main", "Greeter", true, "greet(_)", func(vm *wrengo.WrenVM) {
		name := vm.GetSlotString(1)
		greeting := "Hello, " + name + "!"
		vm.SetSlotString(0, greeting)
	})

	vm := wrengo.NewVMWithForeign()
	defer vm.Free()

	code := `
class Greeter {
  foreign static greet(name)
}

System.print(Greeter.greet("World"))
`

	result, err := vm.Interpret("main", code)
	if err != nil {
		t.Fatalf("Interpret error: %v", err)
	}

	if result != wrengo.ResultSuccess {
		t.Fatalf("Expected ResultSuccess, got %v", result)
	}
}

func TestForeignMethodList(t *testing.T) {
	wrengo.RegisterForeignMethod("main", "ListUtils", true, "sum(_)", func(vm *wrengo.WrenVM) {
		// Get list from slot 1
		count := vm.GetListCount(1)
		sum := 0.0

		// Iterate through list
		for i := 0; i < count; i++ {
			vm.GetListElement(1, i, 2) // Get element i into slot 2
			sum += vm.GetSlotDouble(2)
		}

		// Return sum in slot 0
		vm.SetSlotDouble(0, sum)
	})

	vm := wrengo.NewVMWithForeign()
	defer vm.Free()

	code := `
class ListUtils {
  foreign static sum(list)
}

var numbers = [1, 2, 3, 4, 5]
System.print(ListUtils.sum(numbers))
`

	result, err := vm.Interpret("main", code)
	if err != nil {
		t.Fatalf("Interpret error: %v", err)
	}

	if result != wrengo.ResultSuccess {
		t.Fatalf("Expected ResultSuccess, got %v", result)
	}
}

func TestForeignMethodMap(t *testing.T) {
	wrengo.RegisterForeignMethod("main", "MapUtils", true, "hasKey(_,_)", func(vm *wrengo.WrenVM) {
		// Slot 1: map, Slot 2: key
		hasKey := vm.GetMapContainsKey(1, 2)
		vm.SetSlotBool(0, hasKey)
	})

	vm := wrengo.NewVMWithForeign()
	defer vm.Free()

	code := `
class MapUtils {
  foreign static hasKey(map, key)
}

var dict = {"name": "Wren", "type": "language"}
System.print(MapUtils.hasKey(dict, "name"))
System.print(MapUtils.hasKey(dict, "missing"))
`

	result, err := vm.Interpret("main", code)
	if err != nil {
		t.Fatalf("Interpret error: %v", err)
	}

	if result != wrengo.ResultSuccess {
		t.Fatalf("Expected ResultSuccess, got %v", result)
	}
}

func TestMultipleVMs(t *testing.T) {
	// Register different methods for different modules
	wrengo.RegisterForeignMethod("vm1", "Math", true, "add(_,_)", func(vm *wrengo.WrenVM) {
		a := vm.GetSlotDouble(1)
		b := vm.GetSlotDouble(2)
		vm.SetSlotDouble(0, a+b)
	})

	wrengo.RegisterForeignMethod("vm2", "Math", true, "multiply(_,_)", func(vm *wrengo.WrenVM) {
		a := vm.GetSlotDouble(1)
		b := vm.GetSlotDouble(2)
		vm.SetSlotDouble(0, a*b)
	})

	// Create two VMs
	vm1 := wrengo.NewVMWithForeign()
	defer vm1.Free()

	vm2 := wrengo.NewVMWithForeign()
	defer vm2.Free()

	// Test VM1 with add
	code1 := `
class Math {
  foreign static add(a, b)
}

System.print(Math.add(10, 20))
`

	result1, err := vm1.Interpret("vm1", code1)
	if err != nil {
		t.Fatalf("VM1 Interpret error: %v", err)
	}
	if result1 != wrengo.ResultSuccess {
		t.Fatalf("VM1: Expected ResultSuccess, got %v", result1)
	}

	// Test VM2 with multiply
	code2 := `
class Math {
  foreign static multiply(a, b)
}

System.print(Math.multiply(7, 6))
`

	result2, err := vm2.Interpret("vm2", code2)
	if err != nil {
		t.Fatalf("VM2 Interpret error: %v", err)
	}
	if result2 != wrengo.ResultSuccess {
		t.Fatalf("VM2: Expected ResultSuccess, got %v", result2)
	}
}
