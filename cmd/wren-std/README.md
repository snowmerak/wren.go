# wren-std - Standard Wren CLI

The standard command-line interface for Wren, including built-in async support.

## Installation

```bash
go install github.com/snowmerak/wren.go/cmd/wren-std@latest
```

## Features

- ✅ **Full Wren VM** - Complete Wren language support
- ✅ **Async/Await** - Built-in asynchronous task execution
- ✅ **Interactive REPL** - Multi-line support
- ✅ **Script Execution** - Run `.wren` files
- ✅ **Direct Evaluation** - Execute code from command line

## Usage

### Interactive REPL

```bash
wren-std
# or
wren-std repl
```

```wren
wren> System.print("Hello, Wren!")
Hello, Wren!

wren> var x = 42
wren> System.print(x * 2)
84

wren> class Point {
....>   construct new(x, y) {
....>     _x = x
....>     _y = y
....>   }
....> }
```

### Run a Script

```bash
wren-std run script.wren
# or
wren-std script.wren
```

### Evaluate Code

```bash
wren-std eval "System.print(42)"
# or
wren-std -e "System.print(42)"
```

### Help and Version

```bash
wren-std help
wren-std version
```

## Built-in Features

### Async Support

The standard CLI includes built-in async/await functionality:

```wren
foreign class Async {
    foreign static await(futureId)
    foreign static isReady(futureId)
    foreign static get(futureId)
    foreign static cancel(futureId)
    foreign static getState(futureId)
    foreign static cleanup(futureId)
}

// Example: if you have a foreign method that returns a future
var futureId = SlowTask.start()

// Wait for completion
var result = Async.await(futureId)
System.print("Result: %(result)")
Async.cleanup(futureId)

// Or poll without blocking
while (!Async.isReady(futureId)) {
    System.print("Still working...")
    Fiber.yield()
}
var result = Async.get(futureId)
Async.cleanup(futureId)
```

## Example Scripts

### hello.wren

```wren
System.print("Hello, Wren!")

class Greeter {
    construct new(name) {
        _name = name
    }
    
    greet() {
        System.print("Hello, %(_name)!")
    }
}

var greeter = Greeter.new("World")
greeter.greet()
```

Run with:
```bash
wren-std hello.wren
```

### fibonacci.wren

```wren
class Fibonacci {
    static compute(n) {
        if (n <= 1) return n
        return compute(n - 1) + compute(n - 2)
    }
}

for (i in 0..10) {
    System.print("fib(%(i)) = %(Fibonacci.compute(i))")
}
```

Run with:
```bash
wren-std fibonacci.wren
```

## Building from Source

```bash
git clone https://github.com/snowmerak/wren.go.git
cd wren.go
./build_wren.sh  # or build_wren.bat on Windows
go build -o wren-std ./cmd/wren-std
```

## Creating a Custom CLI

If you need custom foreign methods, use the `wrencli` library:

```go
package main

import (
    "os"
    wrengo "github.com/snowmerak/wren.go"
    "github.com/snowmerak/wren.go/wrencli"
)

//go:generate go run github.com/snowmerak/wren.go/wrengen -dir .

//wren:bind module=main
type MyAPI struct{}

//wren:bind static
func (m *MyAPI) CustomMethod(x float64) float64 {
    return x * 2
}

func main() {
    cli := wrencli.NewCLI(wrencli.Config{
        OnVMCreate: func() *wrengo.WrenVM {
            return wrengo.NewVMWithForeign()
        },
    })
    
    cli.Run(os.Args)
}
```

See [wrencli documentation](../../wrencli/) for more details.

## License

MIT License - see repository root for details.

## Links

- [Wren Language](https://wren.io/)
- [wren.go Repository](https://github.com/snowmerak/wren.go)
- [wrencli Library](../../wrencli/)
