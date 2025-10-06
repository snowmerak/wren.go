# wren.go

Go bindings for the [Wren scripting language](https://wren.io/).

Wren is a small, fast, class-based concurrent scripting language designed for embedding in applications.

## Features

- ✅ VM creation and management
- ✅ Execute Wren scripts
- ✅ Memory management with automatic garbage collection
- ✅ Custom configuration support
- ✅ Version information

## Prerequisites

- Go 1.25.1 or later
- CGO enabled
- C compiler (GCC/MinGW on Windows)
- Git (for submodules)

## Installation

1. Clone the repository with submodules:

```bash
git clone --recursive https://github.com/snowmerak/wren.go.git
cd wren.go
```

2. Build the Wren static library:

```bash
# On Windows with TDM-GCC
C:\path\to\gcc.exe -c -I deps/wren/src/include -I deps/wren/src/vm -I deps/wren/src/optional -std=c99 -O2 deps/wren/src/vm/*.c deps/wren/src/optional/*.c
C:\path\to\ar.exe rcs libwren.a *.o
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
