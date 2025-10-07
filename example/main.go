package main

import (
	"fmt"
	"os"

	wrengo "github.com/snowmerak/gwen"
)

func main() {
	fmt.Println("=== Wren.go Examples ===")
	fmt.Printf("Wren Version: %d\n\n", wrengo.GetVersionNumber())

	// Create a basic VM for simple examples
	basicExamples()

	// Run code generator examples with foreign functions
	codeGenExamples()

	// Demonstrate multi-module usage
	multiModuleExamples()

	// Demonstrate multiple VMs running simultaneously
	multipleVMsExample()

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
var devName = "Developer"
System.print("  Utils.greet('Developer') = %(Utils.greet(devName))")

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

func multiModuleExamples() {
	fmt.Println("\n=== Multi-Module Examples ===\n")

	// Create VM with foreign functions
	vm := wrengo.NewVMWithForeign()
	defer vm.Free()

	// First, define the geometry module with foreign classes
	geometryModule := `
// geometry.wren - Geometry module
class Circle {
  foreign static area(radius)
  foreign static circumference(radius)
}

class Rectangle {
  foreign static area(width, height)
  foreign static perimeter(width, height)
}
`

	fmt.Println("10. Loading Geometry Module:")
	mustInterpret(vm, "geometry", geometryModule)

	// Now import and use the geometry module from main
	mainCode := `
import "geometry" for Circle, Rectangle

System.print("11. Circle Calculations:")
var radius = 5
System.print("  Circle with radius %(radius):")
System.print("  - Area: %(Circle.area(radius))")
System.print("  - Circumference: %(Circle.circumference(radius))")

System.print("\n12. Rectangle Calculations:")
var width = 10
var height = 6
System.print("  Rectangle with width %(width) and height %(height):")
System.print("  - Area: %(Rectangle.area(width, height))")
System.print("  - Perimeter: %(Rectangle.perimeter(width, height))")

System.print("\n13. Mixed Module Usage:")
// You can also import from main module if needed
// Here we demonstrate using both modules
var circleArea = Circle.area(3)
var rectArea = Rectangle.area(4, 5)
System.print("  Circle area + Rectangle area = %(circleArea + rectArea))")
`

	fmt.Println("\n  Using geometry module from main:")
	mustInterpret(vm, "main", mainCode)
}

func multipleVMsExample() {
	fmt.Println("\n=== Multiple VMs Example ===\n")
	fmt.Println("14. Creating Multiple Independent VMs:")

	// IMPORTANT: All VMs share the same foreign method registry
	// RegisterForeignMethod is called globally via init()
	// But each VM has its own execution context and wrapper functions

	// Create first VM
	vm1 := wrengo.NewVMWithForeign()
	defer vm1.Free()

	// Create second VM
	vm2 := wrengo.NewVMWithForeign()
	defer vm2.Free()

	// Both VMs can use the same foreign methods
	code := `
class Calculator {
  foreign static square(x)
  foreign static sqrt(x)
}

System.print("  VM result: square(7) = %(Calculator.square(7))")
`

	fmt.Println("  VM1 executing:")
	mustInterpret(vm1, "main", code)

	fmt.Println("\n  VM2 executing (same code, independent context):")
	mustInterpret(vm2, "main", code)

	// Each VM maintains its own state
	vm1Code := `
var x = 10
System.print("  VM1: x = %(x)")
`
	vm2Code := `
var x = 20
System.print("  VM2: x = %(x)")
`

	fmt.Println("\n15. Independent VM States:")
	mustInterpret(vm1, "main", vm1Code)
	mustInterpret(vm2, "main", vm2Code)

	fmt.Println("\n  ✓ Each VM has its own independent execution context")
	fmt.Println("  ✓ Foreign methods are registered globally but executed per-VM")
	fmt.Println("  ✓ Each VM can handle up to 300 foreign methods")
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
