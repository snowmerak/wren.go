# Wren.go Examples

Complete examples demonstrating all features of wren.go in a single program.

## Quick Start

```bash
# Run all examples
go run example/main.go

# Or with go generate
cd example
go generate  # Regenerate bindings if needed
go run .
```

## What's Included

The example demonstrates **20 different features** in order:

### Basic Wren Features (1-4)

1. **Simple Print** - Hello World
2. **Variables and Math** - Numbers and arithmetic
3. **Classes and Methods** - OOP with Point class
4. **Fibers** - Coroutines/cooperative multitasking

### Foreign Functions via Code Generator (5-9)

5. **Math Operations** - `Math.multiply`
6. **String Operations** - `StringUtils.concat`
7. **Utility Functions** - `Utils.greet`
8. **Calculator Functions** - `square`, `sqrt`, `power`
9. **Complex Calculation** - Pythagorean theorem using multiple functions

### Multi-Module Usage (10-13)

10. **Loading Modules** - Define foreign classes in separate modules
11. **Circle Calculations** - `geometry` module with Circle class
12. **Rectangle Calculations** - `geometry` module with Rectangle class
13. **Cross-Module Usage** - Import and use classes from different modules

### Multiple VMs (14-15)

14. **Creating Multiple VMs** - Two VMs executing the same code independently
15. **Independent VM States** - Each VM has its own execution context

### Async Task Execution (16-20)

16. **Basic Async Task** - Submit and wait for async computations
17. **Multiple Concurrent Tasks** - Run multiple tasks in parallel
18. **Task with Timeout** - Context-based timeout handling
19. **Task Cancellation** - Cancel long-running tasks
20. **Non-blocking Future Check** - Poll for results without blocking

## Multiple VMs: How It Works

**Important:** You can create multiple VM instances simultaneously, but they share the same foreign method **registry**:

```go
// Foreign methods are registered globally (once via init())
func init() {
    RegisterWrenBindings()  // Called once for all VMs
}

// But each VM has independent execution
vm1 := wrengo.NewVMWithForeign()
vm2 := wrengo.NewVMWithForeign()
```

### Key Points

- ✅ **Registry is global**: `RegisterForeignMethod()` stores methods in a global registry
- ✅ **Execution is per-VM**: Each VM has its own:
  - Execution context (variables, classes, state)
  - Wrapper functions (up to 99 per VM)
  - Memory and garbage collector
- ✅ **Thread-safe**: Mutex-protected registry and VM data structures
- ⚠️ **300 method limit per VM**: Each VM can use up to 300 foreign methods

### Example

```go
// Both VMs can use the same foreign methods
vm1 := wrengo.NewVMWithForeign()
vm2 := wrengo.NewVMWithForeign()

// VM1 has its own state
vm1.Interpret("main", `var x = 10`)  // x = 10 in VM1

// VM2 has completely independent state
vm2.Interpret("main", `var x = 20`)  // x = 20 in VM2

// Both can call the same foreign methods
vm1.Interpret("main", `Calculator.square(5)`)  // Works in VM1
vm2.Interpret("main", `Calculator.square(7)`)  // Works in VM2
```

## File Structure

```
example/
├── main.go         # Complete example program
├── math.go         # Go functions with //wren:bind annotations
├── main_wren.go    # Auto-generated bindings (DO NOT EDIT)
└── README.md       # This file
```

## Code Generator Usage

### 1. Annotate Go Functions

See `math.go`:

```go
//go:generate go run ../wrengen/main.go -dir .

//wren:bind module=main
type Math struct{}

//wren:bind static
func (m *Math) Multiply(a, b float64) float64 {
    return a * b
}

//wren:bind module=main class=Calculator name=sqrt static
func Sqrt(x float64) float64 {
    return math.Sqrt(x)
}
```

### 2. Generate Bindings

```bash
cd example
go generate
```

This creates `main_wren.go` with all the boilerplate.

### 3. Use in Wren

```wren
class Math {
  foreign static multiply(a, b)
}

System.print(Math.multiply(7, 6))  // 42
```

## Multi-Module Example

The example also demonstrates how to use **different modules** in Wren:

### 1. Define Foreign Functions in Different Modules

In `math.go`, you can specify different modules:

```go
// Main module functions
//wren:bind module=main class=Calculator name=sqrt static
func Sqrt(x float64) float64 {
    return math.Sqrt(x)
}

// Geometry module functions
//wren:bind module=geometry class=Circle name=area static
func CircleArea(radius float64) float64 {
    return math.Pi * radius * radius
}

//wren:bind module=geometry class=Rectangle name=area static
func RectangleArea(width, height float64) float64 {
    return width * height
}
```

### 2. Define the Module in Wren

First, interpret the module definition:

```go
geometryModule := `
class Circle {
  foreign static area(radius)
  foreign static circumference(radius)
}

class Rectangle {
  foreign static area(width, height)
  foreign static perimeter(width, height)
}
`

vm.Interpret("geometry", geometryModule)
```

### 3. Import and Use in Another Module

Then use `import` to access the classes:

```wren
import "geometry" for Circle, Rectangle

System.print("Circle area: %(Circle.area(5))")
System.print("Rectangle area: %(Rectangle.area(10, 6))")
```

**Key Points:**
- Each module can have its own foreign classes
- Use `import "module_name" for Class1, Class2` to access classes
- The same VM can handle multiple modules simultaneously
- Foreign methods are registered per-module in the code generator

## Annotations Reference

| Annotation | Effect | Example |
|------------|--------|---------|
| `//wren:bind` | Basic binding | `func (m *Math) Add(a, b int)` |
| `//wren:bind static` | Static method | `func (m *Math) Multiply(a, b float64)` |
| `//wren:bind module=X` | Set module name | `//wren:bind module=game` |
| `//wren:bind class=Y` | Override class name | `//wren:bind class=Calculator` |
| `//wren:bind name=Z` | Override method name | `//wren:bind name=concat` |

Combine them:
```go
//wren:bind module=main class=Utils name=greet static
func Greet(name string) string {
    return "Hello, " + name + "!"
}
```

## Supported Types

| Go Type | Wren Type | Notes |
|---------|-----------|-------|
| `int`, `int32`, `int64` | Number | Converted to/from float64 |
| `float32`, `float64` | Number | Direct mapping |
| `string` | String | UTF-8 |
| `bool` | Bool | true/false |
| `error` (return) | Abort | Fiber aborts on error |

### Error Handling

Functions returning `(value, error)` automatically handle errors:

```go
//wren:bind
func (m *Math) Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```

In Wren, division by zero will abort with the error message.

## Expected Output

```
=== Wren.go Examples ===
Wren Version: 4000

=== Basic Examples ===

1. Simple Print:
Hello from Wren!

2. Variables and Math:
x = 42
y = 3.14
x + y = 45.14

3. Classes and Methods:
p1: Point(0, 0)
p2: Point(3, 4)
Distance: 5

4. Fibers (Coroutines):
  Before fiber
  Fiber step 1
  Between calls
  Fiber step 2
  Between calls
  Fiber step 3
  After fiber

=== Code Generator Examples ===

5. Math Operations:
  Math.multiply(7, 6) = 42

6. String Operations:
  StringUtils.concat = Hello, Wren!

7. Utility Functions:
  Utils.greet('Developer') = Hello, Developer!

8. Calculator Functions:
  Calculator.square(5) = 25
  Calculator.sqrt(16) = 4
  Calculator.power(2, 8) = 256

9. Complex Calculation (Pythagorean):
  For triangle (3, 4): c = 5

=== Multi-Module Examples ===

10. Loading Geometry Module:

  Using geometry module from main:

11. Circle Calculations:
  Circle with radius 5:
  - Area: 78.539816339745
  - Circumference: 31.415926535898

12. Rectangle Calculations:
  Rectangle with width 10 and height 6:
  - Area: 60
  - Perimeter: 32

13. Mixed Module Usage:
  Circle area + Rectangle area = 48.274333882308

=== All examples completed successfully! ===
```

## Async Examples

The example also includes **5 async task execution examples** demonstrating the Future-based async system:

### How It Works

The async system uses a **Future pattern** with a managed worker pool:

```go
// Get the global async manager
am := wrengo.GetAsyncManager()

// Submit a task
future := am.Submit(func(ctx context.Context) (interface{}, error) {
    time.Sleep(1 * time.Second)
    return "result", nil
})

// Wait for completion (blocking)
result, err := future.Wait()

// Or check without blocking
if future.IsReady() {
    result, err := future.Get()
}
```

### Key Features

- ✅ **Worker Pool**: Configurable number of goroutines (default: 4)
- ✅ **Future Pattern**: Check status, wait, or poll for results
- ✅ **Context Support**: Timeouts and cancellation via context
- ✅ **Thread-Safe**: Safe for concurrent use
- ✅ **Error Handling**: Errors captured in Future

### Examples in main.go

16. **Basic Task** - Simple async execution with Wait()
17. **Concurrent Tasks** - Submit multiple tasks at once
18. **Timeout** - Cancel task after deadline
19. **Cancellation** - Manually cancel ongoing work
20. **Polling** - Check readiness without blocking

### Future API

```go
future := am.Submit(task)
future.ID()           // Get unique ID
future.IsReady()      // Check if completed
future.State()        // Get current state (Pending/Completed/Failed/Cancelled)
future.Wait()         // Block until done
future.Get()          // Get result if ready
future.Cancel()       // Cancel execution
```

## Learn More

- [Main README](../README.md) - Project overview and installation
- [Wrengen Documentation](../wrengen/README.md) - Complete code generator reference
- [Wren Language Documentation](https://wren.io/) - Official Wren docs

## Next Steps

1. **Modify** `math.go` to add your own functions
2. **Run** `go generate` to regenerate bindings
3. **Test** in `main.go` or create your own Wren scripts
4. **Learn** more in [wrengen/README.md](../wrengen/README.md)
