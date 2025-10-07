package wrencli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	wrengo "github.com/snowmerak/gwen"
)

// StartREPL starts an interactive Read-Eval-Print Loop.
func (c *CLI) StartREPL() error {
	fmt.Println("Wren REPL")
	fmt.Println("Type 'exit' or 'quit' to exit, 'help' for help")
	fmt.Println()

	// Create VM once for the entire REPL session
	vm := c.config.OnVMCreate()
	defer vm.Free()

	reader := bufio.NewReader(os.Stdin)
	var multilineBuffer strings.Builder
	inMultiline := false

	for {
		// Choose prompt
		prompt := c.config.REPLPrompt
		if inMultiline {
			prompt = c.config.REPLMultilinePrompt
		}

		// Print prompt
		fmt.Print(prompt)

		// Read line
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println()
				return nil
			}
			return fmt.Errorf("failed to read input: %w", err)
		}

		line = strings.TrimSpace(line)

		// Handle special commands
		if !inMultiline {
			switch line {
			case "exit", "quit":
				fmt.Println("Goodbye!")
				return nil
			case "help":
				c.printREPLHelp()
				continue
			case "clear":
				// Clear screen (ANSI escape code)
				fmt.Print("\033[H\033[2J")
				continue
			case "":
				continue
			}
		}

		// Check for multiline continuation
		if strings.HasSuffix(line, "{") || strings.HasSuffix(line, "(") {
			inMultiline = true
			multilineBuffer.WriteString(line)
			multilineBuffer.WriteString("\n")
			continue
		}

		// Build complete code
		var code string
		if inMultiline {
			multilineBuffer.WriteString(line)
			multilineBuffer.WriteString("\n")
			code = multilineBuffer.String()
			multilineBuffer.Reset()
			inMultiline = false
		} else {
			code = line
		}

		// Execute code
		result, err := vm.Interpret("repl", code)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		if result != wrengo.ResultSuccess {
			fmt.Printf("Execution failed with result code: %d\n", result)
		}
	}
}

// printREPLHelp prints REPL-specific help.
func (c *CLI) printREPLHelp() {
	fmt.Println(`REPL Commands:
  exit, quit    Exit the REPL
  help          Show this help message
  clear         Clear the screen

Multi-line input:
  Lines ending with { or ( will continue to the next line.
  Complete the block to execute.

Examples:
  System.print("Hello")
  
  class Point {
    construct new(x, y) {
      _x = x
      _y = y
    }
  }
  
  var p = Point.new(10, 20)`)
}
