# Wrencli - Wren Command-Line Interface Library

A flexible and customizable CLI library for the Wren scripting language.

## Overview

`wrencli` provides a reusable command-line interface that can be customized with your own Wren VM configuration and foreign methods. It's designed to be used as a library, allowing you to create your own Wren CLI with custom capabilities.

## Features

- ✅ **Interactive REPL** - Read-Eval-Print Loop with multi-line support
- ✅ **Script Execution** - Run Wren scripts from files
- ✅ **Code Evaluation** - Execute Wren code directly from command line
- ✅ **Customizable VM** - Inject your own VM with custom foreign methods
- ✅ **Flexible Configuration** - Customize prompts, extensions, and behavior

## Installation

```bash
go get github.com/snowmerak/wren.go/wrencli
```

## Quick Start

### Using the Standard CLI

```go
package main

import (
    "os"
    "github.com/snowmerak/wren.go/wrencli"
)

func main() {
    // Create CLI with default configuration (includes async support)
    cli := wrencli.NewCLI(wrencli.Config{})
    
    // Run with command-line arguments
    if err := cli.Run(os.Args); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

### Creating a Custom CLI

```go
package main

import (
    "os"
    wrengo "github.com/snowmerak/wren.go"
    "github.com/snowmerak/wren.go/wrencli"
)

//go:generate go run github.com/snowmerak/wren.go/wrengen -dir .

//wren:bind module=myapp
type MyAPI struct{}

//wren:bind static
func (m *MyAPI) Hello(name string) string {
    return "Hello, " + name + "!"
}

func main() {
    cli := wrencli.NewCLI(wrencli.Config{
        OnVMCreate: func() *wrengo.WrenVM {
            // Create VM with your custom foreign methods
            // (they're registered automatically via init())
            return wrengo.NewVMWithForeign()
        },
        REPLPrompt: "myapp> ",
    })
    
    if err := cli.Run(os.Args); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

## Configuration

### Config Options

```go
type Config struct {
    // OnVMCreate creates a custom VM instance
    // Default: NewVMWithForeign()
    OnVMCreate func() *wrengo.WrenVM

    // ScriptExtension is the file extension for Wren scripts
    // Default: ".wren"
    ScriptExtension string

    // REPLPrompt is the prompt shown in REPL
    // Default: "wren> "
    REPLPrompt string

    // REPLMultilinePrompt is the continuation prompt
    // Default: "....> "
    REPLMultilinePrompt string
}
```

## Usage

### Commands

```bash
# Start REPL (default if no command given)
your-cli
your-cli repl

# Run a script file
your-cli run script.wren
your-cli script.wren  # Shorthand

# Evaluate code directly
your-cli eval "System.print(42)"
your-cli -e "System.print(42)"

# Show help
your-cli help

# Show version
your-cli version
```

### REPL Features

```wren
wren> System.print("Hello!")
Hello!

wren> var x = 42
wren> System.print(x)
42

wren> class Point {
....>   construct new(x, y) {
....>     _x = x
....>     _y = y
....>   }
....> }

wren> var p = Point.new(10, 20)
```

**REPL Commands:**
- `exit` or `quit` - Exit REPL
- `help` - Show help
- `clear` - Clear screen

**Multi-line Input:**
- Lines ending with `{` or `(` continue to the next line
- Complete the block to execute

## Examples

### Example 1: Database CLI

```go
package main

import (
    "database/sql"
    "os"
    wrengo "github.com/snowmerak/wren.go"
    "github.com/snowmerak/wren.go/wrencli"
)

//go:generate go run github.com/snowmerak/wren.go/wrengen -dir .

var db *sql.DB

//wren:bind module=main
type Database struct{}

//wren:bind name=query(_) static
func (d *Database) Query(query string) error {
    rows, err := db.Query(query)
    if err != nil {
        return err
    }
    defer rows.Close()
    // Process rows...
    return nil
}

func main() {
    // Initialize database
    var err error
    db, err = sql.Open("sqlite3", "test.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // Create CLI
    cli := wrencli.NewCLI(wrencli.Config{
        OnVMCreate: func() *wrengo.WrenVM {
            return wrengo.NewVMWithForeign()
        },
        REPLPrompt: "db> ",
    })

    cli.Run(os.Args)
}
```

### Example 2: Game Scripting CLI

```go
package main

import (
    "os"
    wrengo "github.com/snowmerak/wren.go"
    "github.com/snowmerak/wren.go/wrencli"
)

//go:generate go run github.com/snowmerak/wren.go/wrengen -dir .

//wren:bind module=game
type Player struct{}

//wren:bind name=spawn(_,_) static
func (p *Player) Spawn(x, y float64) {
    // Spawn player at position
}

//wren:bind module=game
type Enemy struct{}

//wren:bind name=spawn(_,_,_) static
func (e *Enemy) Spawn(x, y float64, health float64) {
    // Spawn enemy
}

func main() {
    cli := wrencli.NewCLI(wrencli.Config{
        OnVMCreate: func() *wrengo.WrenVM {
            return wrengo.NewVMWithForeign()
        },
        REPLPrompt: "game> ",
    })

    cli.Run(os.Args)
}
```

## Building Your Own CLI

1. **Create your project:**
```bash
mkdir my-wren-cli
cd my-wren-cli
go mod init my-wren-cli
```

2. **Add dependencies:**
```bash
go get github.com/snowmerak/wren.go
go get github.com/snowmerak/wren.go/wrencli
```

3. **Create main.go:**
```go
package main

import (
    "fmt"
    "os"
    wrengo "github.com/snowmerak/wren.go"
    "github.com/snowmerak/wren.go/wrencli"
)

//go:generate go run github.com/snowmerak/wren.go/wrengen -dir .

//wren:bind module=main
type MyAPI struct{}

//wren:bind static
func (m *MyAPI) Greet(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

func main() {
    cli := wrencli.NewCLI(wrencli.Config{
        OnVMCreate: func() *wrengo.WrenVM {
            return wrengo.NewVMWithForeign()
        },
        REPLPrompt: "my-cli> ",
    })

    if err := cli.Run(os.Args); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

4. **Generate bindings and build:**
```bash
go generate
go build -o my-wren-cli
```

5. **Use your CLI:**
```bash
./my-wren-cli repl
my-cli> MyAPI.greet("World")
Hello, World!
```

## API Reference

### CLI Methods

```go
// NewCLI creates a new CLI instance
func NewCLI(config Config) *CLI

// Run executes the CLI with arguments
func (c *CLI) Run(args []string) error

// StartREPL starts the interactive REPL
func (c *CLI) StartREPL() error

// RunScript executes a Wren script file
func (c *CLI) RunScript(path string) error

// RunCode evaluates Wren code directly
func (c *CLI) RunCode(code string) error

// PrintHelp prints the help message
func (c *CLI) PrintHelp(progName string)

// PrintVersion prints version information
func (c *CLI) PrintVersion()
```

## License

MIT License - see main repository for details.

## See Also

- [wren.go](../) - Go bindings for Wren
- [wrengen](../wrengen/) - Code generator for foreign methods
- [Standard CLI](../cmd/wren-std/) - Standard CLI implementation
