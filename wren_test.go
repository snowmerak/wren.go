package wrengo_test

import (
	"testing"
	
	"github.com/snowmerak/wren.go"
)

func TestNewVM(t *testing.T) {
	vm := wrengo.NewVM()
	if vm == nil {
		t.Fatal("Failed to create VM")
	}
	defer vm.Free()
}

func TestGetVersionNumber(t *testing.T) {
	version := wrengo.GetVersionNumber()
	if version <= 0 {
		t.Fatalf("Invalid version number: %d", version)
	}
	t.Logf("Wren version: %d", version)
}

func TestInterpretSimple(t *testing.T) {
	vm := wrengo.NewVM()
	if vm == nil {
		t.Fatal("Failed to create VM")
	}
	defer vm.Free()
	
	source := `System.print("Hello from Wren!")`
	result, err := vm.Interpret("main", source)
	
	if err != nil {
		t.Fatalf("Interpret error: %v", err)
	}
	
	if result != wrengo.ResultSuccess {
		t.Fatalf("Expected ResultSuccess, got %v", result)
	}
}

func TestInterpretWithVariable(t *testing.T) {
	vm := wrengo.NewVM()
	if vm == nil {
		t.Fatal("Failed to create VM")
	}
	defer vm.Free()
	
	source := `
var x = 42
var y = 3.14
System.print("x = %(x)")
System.print("y = %(y)")
System.print("x + y = %(x + y)")
`
	
	result, err := vm.Interpret("main", source)
	
	if err != nil {
		t.Fatalf("Interpret error: %v", err)
	}
	
	if result != wrengo.ResultSuccess {
		t.Fatalf("Expected ResultSuccess, got %v", result)
	}
}

func TestInterpretWithClass(t *testing.T) {
	vm := wrengo.NewVM()
	if vm == nil {
		t.Fatal("Failed to create VM")
	}
	defer vm.Free()
	
	source := `
class Greeter {
  construct new(name) {
    _name = name
  }
  
  greet() {
    System.print("Hello, %(_name)!")
  }
}

var greeter = Greeter.new("World")
greeter.greet()
`
	
	result, err := vm.Interpret("main", source)
	
	if err != nil {
		t.Fatalf("Interpret error: %v", err)
	}
	
	if result != wrengo.ResultSuccess {
		t.Fatalf("Expected ResultSuccess, got %v", result)
	}
}

func TestInterpretCompileError(t *testing.T) {
	vm := wrengo.NewVM()
	if vm == nil {
		t.Fatal("Failed to create VM")
	}
	defer vm.Free()
	
	source := `
var x = 
`
	
	result, err := vm.Interpret("main", source)
	
	if err != nil {
		t.Fatalf("Interpret error: %v", err)
	}
	
	if result != wrengo.ResultCompileError {
		t.Fatalf("Expected ResultCompileError, got %v", result)
	}
}

func TestInterpretWithConfig(t *testing.T) {
	config := wrengo.DefaultConfiguration()
	config.InitialHeapSize = 5 * 1024 * 1024 // 5MB
	
	vm := wrengo.NewVMWithConfig(config)
	if vm == nil {
		t.Fatal("Failed to create VM with config")
	}
	defer vm.Free()
	
	source := `System.print("VM with custom config")`
	result, err := vm.Interpret("main", source)
	
	if err != nil {
		t.Fatalf("Interpret error: %v", err)
	}
	
	if result != wrengo.ResultSuccess {
		t.Fatalf("Expected ResultSuccess, got %v", result)
	}
}

func TestCollectGarbage(t *testing.T) {
	vm := wrengo.NewVM()
	if vm == nil {
		t.Fatal("Failed to create VM")
	}
	defer vm.Free()
	
	// Create some garbage
	source := `
var list = []
for (i in 1..1000) {
  list.add(i)
}
`
	
	result, err := vm.Interpret("main", source)
	if err != nil {
		t.Fatalf("Interpret error: %v", err)
	}
	
	if result != wrengo.ResultSuccess {
		t.Fatalf("Expected ResultSuccess, got %v", result)
	}
	
	// Force garbage collection
	vm.CollectGarbage()
}
