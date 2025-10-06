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

The example demonstrates **9 different features** in order:

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

=== All examples completed successfully! ===
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
