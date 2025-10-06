package main

import (
	"fmt"
	"os"

	wrengo "github.com/snowmerak/wren.go"
)

func main() {
	// Print Wren version
	fmt.Printf("Wren Version: %d\n", wrengo.GetVersionNumber())

	// Create a new VM
	vm := wrengo.NewVM()
	defer vm.Free()

	// Simple example
	fmt.Println("\n=== Simple Print ===")
	simpleCode := `System.print("Hello from Wren!")`
	if result, err := vm.Interpret("main", simpleCode); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	} else if result != wrengo.ResultSuccess {
		fmt.Fprintf(os.Stderr, "Execution failed with result: %v\n", result)
		os.Exit(1)
	}

	// Class example
	fmt.Println("\n=== Class Example ===")
	classCode := `
class Point {
  construct new(x, y) {
    _x = x
    _y = y
  }
  
  x { _x }
  y { _y }
  
  toString {
    return "Point(%(_x), %(_y))"
  }
  
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

	if result, err := vm.Interpret("main", classCode); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	} else if result != wrengo.ResultSuccess {
		fmt.Fprintf(os.Stderr, "Execution failed with result: %v\n", result)
		os.Exit(1)
	}

	// Fiber example
	fmt.Println("\n=== Fiber Example ===")
	fiberCode := `
var fiber = Fiber.new {
  System.print("Fiber step 1")
  Fiber.yield()
  System.print("Fiber step 2")
  Fiber.yield()
  System.print("Fiber step 3")
}

System.print("Before fiber")
fiber.call()
System.print("Between calls")
fiber.call()
System.print("Between calls")
fiber.call()
System.print("After fiber")
`

	if result, err := vm.Interpret("main", fiberCode); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	} else if result != wrengo.ResultSuccess {
		fmt.Fprintf(os.Stderr, "Execution failed with result: %v\n", result)
		os.Exit(1)
	}

	fmt.Println("\n=== All examples completed successfully! ===")
}
