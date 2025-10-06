# wren.go

Go bindings for the [Wren scripting language](https://wren.io/).

Wren is a small, fast, class-based concurrent scripting language designed for embedding in applications.

## Features

- ✅ VM creation and management
- ✅ Execute Wren scripts
- ✅ Memory management with automatic garbage collection
- ✅ Custom configuration support
- ✅ Foreign function/class bindings
- ✅ Complete slot API for data exchange
- ✅ **Code generator for automatic bindings**
- ✅ Version information

## Prerequisites

- Go 1.25.1 or later
- CGO enabled
- C compiler (GCC/MinGW on Windows)
- Git (for submodules)

## Installation

### As a Library (Recommended)

When using as a dependency in your project:

```bash
go get github.com/snowmerak/wren.go
```

**Important**: Before building your project, you need to build the Wren static library once:

```bash
cd $GOPATH/src/github.com/snowmerak/wren.go  # or wherever go get placed it
# On Linux/Mac
./build_wren.sh
# On Windows
build_wren.bat
```

Or if using Go modules and your module cache:

```bash
cd $(go list -f '{{.Dir}}' github.com/snowmerak/wren.go)
# Run the appropriate build script
```

### For Development

1. Clone the repository with submodules:

```bash
git clone --recursive https://github.com/snowmerak/wren.go.git
cd wren.go
```

2. Build the Wren static library:

```bash
# On Linux/Mac
./build_wren.sh

# On Windows
build_wren.bat
```

3. Build the Go package:

```bash
go build
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/snowmerak/wren.go"
)

func main() {
    // Create a new Wren VM
    vm := wrengo.NewVM()
    defer vm.Free()
    
    // Execute Wren code
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

### With Custom Configuration

```go
config := wrengo.DefaultConfiguration()
config.InitialHeapSize = 5 * 1024 * 1024 // 5MB

vm := wrengo.NewVMWithConfig(config)
defer vm.Free()
```

## API

### VM Management

- `NewVM()` - Create a new VM with default configuration
- `NewVMWithConfig(config)` - Create a new VM with custom configuration
- `Free()` - Dispose of VM resources
- `CollectGarbage()` - Manually trigger garbage collection

### Script Execution

- `Interpret(module, source)` - Execute Wren source code

### Configuration

- `DefaultConfiguration()` - Get default configuration
- `InitialHeapSize` - Initial heap size in bytes (default: 10MB)
- `MinHeapSize` - Minimum heap size in bytes (default: 1MB)
- `HeapGrowthPercent` - Heap growth percentage (default: 50)

### Utilities

- `GetVersionNumber()` - Get Wren version number

## Code Generator (Wrengen)

Automatically generate Wren bindings from annotated Go code!

```go
//go:generate go run github.com/snowmerak/wren.go/wrengen -dir .

//wren:bind module=main
type Math struct{}

//wren:bind static
func (m *Math) Multiply(a, b float64) float64 {
    return a * b
}
```

Run `go generate` and use in Wren:

```wren
System.print(Math.multiply(7, 6))  // 42
```

See [wrengen/README.md](./wrengen/README.md) for full documentation.

## Examples

See the [example](./example) directory for more examples:

```bash
go run example/main.go
```

## Testing

```bash
go test -v
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Credits

- [Wren Programming Language](https://wren.io/) by Bob Nystrom
- This binding by snowmerak
