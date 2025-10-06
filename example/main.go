package main

import (
	"fmt"
	"os"

	wrengo "github.com/snowmerak/wren.go"
)

func main() {
	fmt.Println("=== Wren.go Examples ===")
	fmt.Printf("Wren Version: %d\n\n", wrengo.GetVersionNumber())

	// Create a basic VM for simple examples
	basicExamples()

	// Run code generator examples with foreign functions
	codeGenExamples()

	fmt.Println("\n=== All examples completed successfully! ===")
}

func basicExamples() {
	fmt.Println("=== Basic Examples ===\n")

	vm := wrengo.NewVM()
	defer vm.Free()

	// 1. Simple Print
	fmt.Println("1. Simple Print:")
	simpleCode := `System.print("Hello from Wren!")`
	mustInterpret(vm, "main", simpleCode)

	// 2. Variables and Math
	fmt.Println("\n2. Variables and Math:")
	mathCode := `
var x = 42
var y = 3.14
System.print("x = %(x)")
System.print("y = %(y)")
System.print("x + y = %(x + y)")
`
	mustInterpret(vm, "main", mathCode)

	// 3. Classes and Methods
	fmt.Println("\n3. Classes and Methods:")
	classCode := `
class Point {
  construct new(x, y) {
    _x = x
    _y = y
  }
  
  x { _x }
  y { _y }
  
  toString { "Point(%(_x), %(_y))" }
  
  distance(other) {
    var dx = _x - other.x
    var dy = _y - other.y
    return (dx * dx + dy * dy).sqrt
  }
}

var p1 = Point.new(0, 0)
var p2 = Point.new(3, 4)
System.print("p1: %(p1.toString)")
System.print("p2: %(p2.toString)")
System.print("Distance: %(p1.distance(p2))")
`
	mustInterpret(vm, "main", classCode)

	// 4. Fibers (Coroutines)
	fmt.Println("\n4. Fibers (Coroutines):")
	fiberCode := `
var fiber = Fiber.new {
  System.print("  Fiber step 1")
  Fiber.yield()
  System.print("  Fiber step 2")
  Fiber.yield()
  System.print("  Fiber step 3")
}

System.print("  Before fiber")
fiber.call()
System.print("  Between calls")
fiber.call()
System.print("  Between calls")
fiber.call()
System.print("  After fiber")
`
	mustInterpret(vm, "main", fiberCode)
}

func codeGenExamples() {
	fmt.Println("\n=== Code Generator Examples ===\n")

	// Create VM with foreign functions (auto-registered via init())
	vm := wrengo.NewVMWithForeign()
	defer vm.Free()

	code := `
// Define foreign classes
class Math {
  foreign static multiply(a, b)
}

class StringUtils {
  foreign static concat(a, b)
}

class Utils {
  foreign static greet(name)
}

class Calculator {
  foreign static square(x)
  foreign static sqrt(x)
  foreign static power(base, exponent)
}

// Math operations
System.print("5. Math Operations:")
System.print("  Math.multiply(7, 6) = %(Math.multiply(7, 6))")

// String operations
System.print("\n6. String Operations:")
var hello = "Hello, "
var world = "Wren!"
System.print("  StringUtils.concat = %(StringUtils.concat(hello, world))")

// Utility functions
System.print("\n7. Utility Functions:")
System.print("  Utils.greet('Developer') = %(Utils.greet("Developer"))")

// Calculator functions
System.print("\n8. Calculator Functions:")
System.print("  Calculator.square(5) = %(Calculator.square(5))")
System.print("  Calculator.sqrt(16) = %(Calculator.sqrt(16))")
System.print("  Calculator.power(2, 8) = %(Calculator.power(2, 8))")

// Complex calculation - Pythagorean theorem
System.print("\n9. Complex Calculation (Pythagorean):")
var a = 3
var b = 4
var c = Calculator.sqrt(Calculator.square(a) + Calculator.square(b))
System.print("  For triangle (3, 4): c = %(c)")
`

	mustInterpret(vm, "main", code)
}

func mustInterpret(vm *wrengo.WrenVM, module, source string) {
	result, err := vm.Interpret(module, source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	if result != wrengo.ResultSuccess {
		fmt.Fprintf(os.Stderr, "Execution failed with result: %v\n", result)
		os.Exit(1)
	}
}
