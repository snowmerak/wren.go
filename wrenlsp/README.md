# wrenlsp - Wren Language Server Protocol

`wrenlsp` is a library that provides LSP (Language Server Protocol) functionality for the Wren programming language.

## Features

- ✅ **Completion**: Keyword and foreign method autocompletion
- ✅ **Hover**: Hover information for symbols (planned)
- ✅ **Diagnostics**: Syntax error detection (planned)
- ✅ **Customizable**: Inject custom VM with foreign methods

## Quick Start

### Using the Standard LSP Server

The easiest way to use wrenlsp is through the standard LSP server binary:

```bash
# Build the LSP server
python build.py lsp

# Or build everything
python build.py all
```

This creates `bin/wren-lsp-std` which can be used with any LSP client.

### VS Code Integration

Create a `.vscode/settings.json` in your Wren project:

```json
{
  "wren.languageServer.path": "path/to/wren-lsp-std"
}
```

Or install a Wren VS Code extension that supports LSP.

### Vim/Neovim Integration

Add to your LSP configuration:

```lua
require'lspconfig'.configs.wren = {
  default_config = {
    cmd = {'path/to/wren-lsp-std'},
    filetypes = {'wren'},
    root_dir = function(fname)
      return vim.fn.getcwd()
    end,
  },
}

require'lspconfig'.wren.setup{}
```

## Using as a Library

If you need to customize the LSP server, you can use wrenlsp as a library:

```go
package main

import (
    "fmt"
    "os"
    
    "github.com/snowmerak/wren.go/wrenlsp"
    wrengo "github.com/snowmerak/wren.go"
)

func main() {
    config := wrenlsp.Config{
        EnableDiagnostics: true,
        OnVMCreate: func() *wrengo.WrenVM {
            // Create VM with custom foreign methods
            vm := wrengo.NewVMWithForeign()
            // Register your foreign methods here
            return vm
        },
    }
    
    server := wrenlsp.NewServer(config)
    
    // Register foreign methods for autocompletion
    server.RegisterForeignMethod(wrenlsp.ForeignMethodInfo{
        Module:    "mymodule",
        Class:     "MyClass",
        Method:    "myMethod",
        Signature: "myMethod(_,_)",
        IsStatic:  false,
        Doc:       "My custom method documentation",
    })
    
    if err := server.Serve(); err != nil {
        fmt.Fprintf(os.Stderr, "LSP server error: %v\n", err)
        os.Exit(1)
    }
}
```

## API Reference

### Config

```go
type Config struct {
    // OnVMCreate creates a custom VM instance for analysis
    OnVMCreate func() *WrenVM
    
    // ForeignMethods contains information about available foreign methods
    ForeignMethods []ForeignMethodInfo
    
    // EnableDiagnostics enables syntax error diagnostics
    EnableDiagnostics bool
}
```

### ForeignMethodInfo

```go
type ForeignMethodInfo struct {
    Module    string  // Module name
    Class     string  // Class name
    Method    string  // Method name
    Signature string  // Full signature, e.g. "method(_,_)"
    IsStatic  bool    // Whether it's a static method
    Doc       string  // Documentation text
}
```

### Server Methods

```go
// NewServer creates a new LSP server instance
func NewServer(config Config) *Server

// RegisterForeignMethod registers a foreign method for autocompletion
func (s *Server) RegisterForeignMethod(info ForeignMethodInfo)

// Serve starts the LSP server and processes requests
func (s *Server) Serve() error
```

## Supported LSP Features

### Implemented

- ✅ **initialize**: Server initialization
- ✅ **textDocument/completion**: Autocompletion for keywords and foreign methods
- ✅ **textDocument/didOpen**: Document open notification
- ✅ **textDocument/didChange**: Document change notification

### Planned

- 🚧 **textDocument/hover**: Hover information
- 🚧 **textDocument/definition**: Go to definition
- 🚧 **textDocument/publishDiagnostics**: Syntax error diagnostics
- 🚧 **textDocument/formatting**: Code formatting

## Examples

### Example 1: Basic LSP Server

```go
package main

import (
    "github.com/snowmerak/wren.go/wrenlsp"
)

func main() {
    server := wrenlsp.NewServer(wrenlsp.Config{
        EnableDiagnostics: true,
    })
    
    server.Serve()
}
```

### Example 2: Custom Foreign Methods

```go
package main

import (
    "github.com/snowmerak/wren.go/wrenlsp"
    wrengo "github.com/snowmerak/wren.go"
)

func main() {
    config := wrenlsp.Config{
        EnableDiagnostics: true,
        OnVMCreate: func() *wrengo.WrenVM {
            vm := wrengo.NewVMWithForeign()
            // Your custom VM setup
            return vm
        },
    }
    
    server := wrenlsp.NewServer(config)
    
    // Register async methods
    server.RegisterForeignMethod(wrenlsp.ForeignMethodInfo{
        Module:    "async",
        Class:     "Future",
        Method:    "new",
        Signature: "new(_)",
        IsStatic:  true,
        Doc:       "Creates a new Future that will be resolved with the result of the given fiber",
    })
    
    server.Serve()
}
```

## Testing

Run the test suite:

```bash
go test ./wrenlsp -v
```

Current test coverage includes:
- Server creation and configuration
- Foreign method registration
- Initialize request handling
- Completion request handling

## Protocol

wrenlsp implements the [Language Server Protocol](https://microsoft.github.io/language-server-protocol/) specification.

Communication uses JSON-RPC 2.0 over stdin/stdout:

```
Content-Length: 123\r\n
\r\n
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{...}}
```

## Limitations

- Currently only supports full document synchronization (not incremental)
- Hover information is not yet implemented
- No workspace symbol search yet
- Diagnostics are planned but not yet implemented

## Contributing

Contributions are welcome! Areas that need work:

1. **Hover Information**: Show type information and documentation
2. **Diagnostics**: Implement syntax error detection using Wren VM
3. **Go to Definition**: Navigate to symbol definitions
4. **Workspace Symbols**: Search for symbols across files
5. **Code Formatting**: Auto-format Wren code

## License

Same as wren.go - see LICENSE file.
