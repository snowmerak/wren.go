# Wren.go Runtime Features Guide

A comprehensive guide to using the wren.go runtime system and its advanced features.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Basic VM Usage](#basic-vm-usage)
3. [Foreign Functions](#foreign-functions)
4. [Code Generator (Wrengen)](#code-generator-wrengen)
5. [Asynchronous Tasks](#asynchronous-tasks)
6. [Multiple VM Instances](#multiple-vm-instances)
7. [Data Exchange with Slots](#data-exchange-with-slots)
8. [Best Practices](#best-practices)

---

## Getting Started

### Installation

```bash
go get github.com/snowmerak/wren.go
cd $(go list -f '{{.Dir}}' github.com/snowmerak/wren.go)
# Build the Wren library
./build_wren.sh  # On Linux/Mac
# or
build_wren.bat   # On Windows
```

### Basic Example

```go
package main

import (
    "fmt"
    wrengo "github.com/snowmerak/wren.go"
)

func main() {
    vm := wrengo.NewVM()
    defer vm.Free()
    
    code := `System.print("Hello from Wren!")`
    result, err := vm.Interpret("main", code)
    
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    
    if result == wrengo.ResultSuccess {
        fmt.Println("Success!")
    }
}
```

---

## Basic VM Usage

### Creating a VM

**Default Configuration:**

```go
vm := wrengo.NewVM()
defer vm.Free()  // Always free resources!
```

**Custom Configuration:**

```go
config := wrengo.DefaultConfiguration()
config.InitialHeapSize = 10 * 1024 * 1024  // 10MB
config.MinHeapSize = 1 * 1024 * 1024       // 1MB
config.HeapGrowthPercent = 50              // 50% growth

vm := wrengo.NewVMWithConfig(config)
defer vm.Free()
```

**With Foreign Methods:**

```go
// Automatically registers all foreign methods
vm := wrengo.NewVMWithForeign()
defer vm.Free()
```

### Executing Wren Code

```go
result, err := vm.Interpret("main", `
    var x = 42
    System.print("The answer is %(x)")
`)

if err != nil {
    fmt.Println("Error:", err)
}
```

### Memory Management

```go
// Manual garbage collection
vm.CollectGarbage()
```

---

## Foreign Functions

Foreign functions allow Go code to be called from Wren scripts.

### Manual Registration

```go
// Register a foreign method
wrengo.RegisterForeignMethod(
    "main",           // module name
    "Calculator",     // class name
    true,             // is static
    "add(_,_)",       // signature
    func(vm *wrengo.WrenVM) {
        // Get parameters from slots
        a := vm.GetSlotDouble(1)
        b := vm.GetSlotDouble(2)
        
        // Return result
        vm.SetSlotDouble(0, a + b)
    },
)
```

### Using in Wren

```wren
foreign class Calculator {
    foreign static add(a, b)
}

System.print(Calculator.add(10, 20))  // 30
```

---

## Code Generator (Wrengen)

Automatically generate foreign function bindings from annotated Go code.

### Basic Usage

**1. Annotate Your Go Code:**

```go
//go:generate go run github.com/snowmerak/wren.go/wrengen -dir .

//wren:bind module=main
type Math struct{}

//wren:bind static
func (m *Math) Multiply(a, b float64) float64 {
    return a * b
}

//wren:bind name=add(_,_) static
func (m *Math) Add(a, b float64) float64 {
    return a + b
}
```

**2. Generate Bindings:**

```bash
go generate
```

This creates `<package>_wren.go` with automatic bindings.

**3. Use in Wren:**

```wren
foreign class Math {
    foreign static multiply(a, b)
    foreign static add(a, b)
}

System.print(Math.multiply(7, 6))  // 42
System.print(Math.add(10, 20))     // 30
```

### Annotation Options

```go
//wren:bind                        // Use default settings
//wren:bind module=mymodule        // Specify module name
//wren:bind class=MyClass          // Specify class name
//wren:bind name=customName(_)     // Custom method name and signature
//wren:bind static                 // Mark as static method
```

### Supported Types

| Go Type | Wren Type | Get Method | Set Method |
|---------|-----------|------------|------------|
| `float64` | Num | `GetSlotDouble()` | `SetSlotDouble()` |
| `float32` | Num | `GetSlotDouble()` | `SetSlotDouble()` |
| `int`, `int64` | Num | `GetSlotDouble()` | `SetSlotDouble()` |
| `string` | String | `GetSlotString()` | `SetSlotString()` |
| `bool` | Bool | `GetSlotBool()` | `SetSlotBool()` |
| `error` | String (abort) | - | `AbortFiber()` |

### Complete Example

```go
package example

//go:generate go run github.com/snowmerak/wren.go/wrengen -dir .

//wren:bind module=geometry
type Circle struct{}

//wren:bind name=area(_) static
func (c *Circle) Area(radius float64) float64 {
    return 3.14159 * radius * radius
}

//wren:bind name=circumference(_) static
func (c *Circle) Circumference(radius float64) float64 {
    return 2 * 3.14159 * radius
}
```

```wren
// geometry.wren
foreign class Circle {
    foreign static area(radius)
    foreign static circumference(radius)
}

var r = 5
System.print("Area: %(Circle.area(r))")
System.print("Circumference: %(Circle.circumference(r))")
```

---

## Asynchronous Tasks

The async system allows non-blocking execution of long-running tasks.

### Overview

- **Future-based**: Submit tasks and get a `Future` to track results
- **Worker pool**: Background workers execute tasks concurrently
- **Context support**: Timeouts and cancellation
- **Thread-safe**: Safe for concurrent use

### Basic Async Task

**Go Side:**

```go
// Register an async foreign method
wrengo.RegisterAsyncMethod("slowComputation", func(args ...interface{}) (interface{}, error) {
    // Simulate long computation
    time.Sleep(2 * time.Second)
    return 42, nil
})
```

**Wren Side:**

```wren
foreign class SlowTask {
    foreign static start()
}

// Submit task and get future ID
var futureId = SlowTask.start()

// Wait for result (blocks until complete)
var result = Async.await(futureId)
System.print("Result: %(result)")

// Clean up
Async.cleanup(futureId)
```

### Non-Blocking Check

```wren
var futureId = SlowTask.start()

// Check without blocking
while (!Async.isReady(futureId)) {
    System.print("Still working...")
    Fiber.yield()
}

var result = Async.get(futureId)
System.print("Done: %(result)")
Async.cleanup(futureId)
```

### Task Cancellation

```wren
var futureId = LongTask.start()

// Cancel the task
Async.cancel(futureId)

// Check state
var state = Async.getState(futureId)
System.print("State: %(state)")  // "cancelled"

Async.cleanup(futureId)
```

### Multiple Concurrent Tasks

```wren
var tasks = []
for (i in 1..5) {
    tasks.add(Task.start(i))
}

// Wait for all tasks
for (futureId in tasks) {
    var result = Async.await(futureId)
    System.print("Task result: %(result)")
    Async.cleanup(futureId)
}
```

### Async API Reference

| Method | Description | Returns |
|--------|-------------|---------|
| `Async.await(id)` | Wait for future to complete | Result value or error |
| `Async.isReady(id)` | Check if future is ready | `true` or `false` |
| `Async.get(id)` | Get result (must be ready) | Result value or error |
| `Async.cancel(id)` | Cancel the task | `true` if cancelled |
| `Async.getState(id)` | Get task state | `"pending"`, `"completed"`, `"failed"`, `"cancelled"` |
| `Async.cleanup(id)` | Remove future from manager | - |

### Go Async Manager API

```go
// Get the global async manager
mgr := wrengo.GetAsyncManager()

// Submit a task
future := mgr.Submit(func(ctx context.Context) (interface{}, error) {
    // Check for cancellation
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }
    
    // Do work
    result := expensiveComputation()
    return result, nil
})

// Wait for result
result, err := future.Wait()

// Or wait with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := future.WaitContext(ctx)
```

### Complete Async Example

```go
package main

import (
    "context"
    "time"
    wrengo "github.com/snowmerak/wren.go"
)

func init() {
    // Register async foreign method
    wrengo.RegisterForeignMethod("main", "SlowTask", true, "start()", 
        func(vm *wrengo.WrenVM) {
            // Submit async task
            future := wrengo.GetAsyncManager().Submit(
                func(ctx context.Context) (interface{}, error) {
                    // Simulate work
                    for i := 0; i < 10; i++ {
                        select {
                        case <-ctx.Done():
                            return nil, ctx.Err()
                        default:
                            time.Sleep(200 * time.Millisecond)
                        }
                    }
                    return 42.0, nil
                },
            )
            
            // Store future and return ID
            id := wrengo.GetAsyncManager().StoreFuture(future)
            vm.SetSlotDouble(0, float64(id))
        },
    )
}

func main() {
    vm := wrengo.NewVMWithForeign()
    defer vm.Free()
    
    wrenScript := `
        foreign class SlowTask {
            foreign static start()
        }
        
        System.print("Starting task...")
        var id = SlowTask.start()
        
        System.print("Working on other things...")
        
        while (!Async.isReady(id)) {
            System.print("Still computing...")
            Fiber.yield()
        }
        
        var result = Async.get(id)
        System.print("Result: %(result)")
        Async.cleanup(id)
    `
    
    vm.Interpret("main", wrenScript)
}
```

---

## Multiple VM Instances

You can create multiple VM instances that execute independently.

### Key Concepts

- **Registry is global**: Foreign methods are registered once for all VMs
- **Execution is per-VM**: Each VM has its own state, variables, and memory
- **Thread-safe**: Safe to use multiple VMs concurrently
- **Wrapper limit**: Each VM can use up to 300 foreign methods

### Example

```go
// Create two independent VMs
vm1 := wrengo.NewVMWithForeign()
vm2 := wrengo.NewVMWithForeign()
defer vm1.Free()
defer vm2.Free()

// VM1 has its own state
vm1.Interpret("main", `
    var x = 10
    System.print("VM1: x = %(x)")
`)

// VM2 has completely independent state
vm2.Interpret("main", `
    var x = 20
    System.print("VM2: x = %(x)")
`)

// Both can call the same foreign methods
vm1.Interpret("main", `System.print(Math.multiply(3, 4))`)
vm2.Interpret("main", `System.print(Math.multiply(5, 6))`)
```

### Use Cases

- **Sandboxing**: Isolate script execution
- **Multi-tenancy**: Run scripts for different users
- **Testing**: Test different scenarios in parallel
- **Module isolation**: Keep module states separate

---

## Data Exchange with Slots

Slots are used to pass data between Go and Wren.

### Understanding Slots

- **Slot 0**: Return value (for Go → Wren)
- **Slot 1+**: Parameters (for Wren → Go)
- Must ensure slot exists before use

### Getting Values from Wren

```go
func myForeignMethod(vm *wrengo.WrenVM) {
    // Ensure slot exists
    vm.EnsureSlots(3)
    
    // Get parameters
    numValue := vm.GetSlotDouble(1)
    strValue := vm.GetSlotString(2)
    boolValue := vm.GetSlotBool(3)
    
    // Process...
}
```

### Returning Values to Wren

```go
func myForeignMethod(vm *wrengo.WrenVM) {
    result := 42.0
    
    // Return via slot 0
    vm.SetSlotDouble(0, result)
}
```

### Lists

```go
// Get list from Wren
vm.EnsureSlots(2)
count := vm.GetListCount(1)

for i := 0; i < count; i++ {
    vm.GetListElement(1, i, 0)
    value := vm.GetSlotDouble(0)
    // Process value...
}

// Create list in Go
vm.SetSlotNewList(0)
vm.SetSlotDouble(1, 10.0)
vm.InsertInList(0, -1, 1)  // Append 10.0
vm.SetSlotDouble(1, 20.0)
vm.InsertInList(0, -1, 1)  // Append 20.0
```

### Maps

```go
// Create map
vm.SetSlotNewMap(0)

// Set key-value pairs
vm.SetSlotString(1, "name")
vm.SetSlotString(2, "Alice")
vm.SetMapValue(0, 1, 2)

vm.SetSlotString(1, "age")
vm.SetSlotDouble(2, 30)
vm.SetMapValue(0, 1, 2)
```

### Error Handling

```go
func myForeignMethod(vm *wrengo.WrenVM) {
    if errorCondition {
        vm.SetSlotString(0, "Error message")
        vm.AbortFiber(0)
        return
    }
    
    // Normal execution...
}
```

---

## Best Practices

### 1. Always Free VMs

```go
vm := wrengo.NewVM()
defer vm.Free()  // ✅ Always use defer
```

### 2. Use Code Generator

```go
// ✅ Good: Use wrengen annotations
//wren:bind static
func (m *Math) Add(a, b float64) float64 {
    return a + b
}

// ❌ Bad: Manual registration (error-prone)
wrengo.RegisterForeignMethod("main", "Math", true, "add(_,_)", ...)
```

### 3. Check Async Task State

```wren
// ✅ Good: Always check if ready before getting
if (Async.isReady(id)) {
    var result = Async.get(id)
}

// ❌ Bad: Getting without checking (may block or error)
var result = Async.get(id)
```

### 4. Clean Up Futures

```wren
var id = Task.start()
var result = Async.await(id)
Async.cleanup(id)  // ✅ Always cleanup when done
```

### 5. Handle Errors

```go
result, err := vm.Interpret("main", code)
if err != nil {
    log.Printf("Error: %v", err)  // ✅ Handle errors
    return
}
```

### 6. Use Context for Cancellation

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

future := mgr.Submit(func(ctx context.Context) (interface{}, error) {
    select {
    case <-ctx.Done():
        return nil, ctx.Err()  // ✅ Respect cancellation
    default:
        // Do work...
    }
})
```

### 7. Limit VM Count

```go
// ✅ Good: Reuse VMs
var vmPool = make(chan *wrengo.WrenVM, 10)

// ❌ Bad: Create unlimited VMs
for i := 0; i < 10000; i++ {
    vm := wrengo.NewVM()  // Memory intensive!
}
```

### 8. Module Organization

```
project/
├── main.go              # Main application
├── math.go              # Math foreign functions
├── geometry.go          # Geometry foreign functions
├── main_wren.go         # Generated bindings (from wrengen)
└── scripts/
    ├── main.wren        # Main script
    ├── math.wren        # Math module
    └── geometry.wren    # Geometry module
```

---

## Examples

See the `example/` directory for complete working examples:

- `example/main.go` - Comprehensive feature demonstration
- `example/async_visual/` - Visual async execution demo
- `example/README.md` - Detailed example documentation

---

## Further Reading

- [Wren Language Guide](https://wren.io/)
- [Wrengen Documentation](./wrengen/README.md)
- [Example Documentation](./example/README.md)
- [Wren Embedding Guide](https://wren.io/embedding/)

---

## Support

- **Issues**: https://github.com/snowmerak/wren.go/issues
- **Discussions**: https://github.com/snowmerak/wren.go/discussions

---

**Happy Wrening! 🚀**
