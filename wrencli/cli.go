package wrencli

import (
	"fmt"
	"os"

	wrengo "github.com/snowmerak/wren.go"
)

// CLI represents a Wren command-line interface.
type CLI struct {
	config Config
}

// Config defines the configuration for the CLI.
type Config struct {
	// OnVMCreate is called to create a custom VM instance.
	// If nil, NewVMWithForeign() is used by default.
	OnVMCreate func() *wrengo.WrenVM

	// ScriptExtension is the default file extension for Wren scripts.
	// Default is ".wren".
	ScriptExtension string

	// REPLPrompt is the prompt string shown in the REPL.
	// Default is "wren> ".
	REPLPrompt string

	// REPLMultilinePrompt is the prompt string shown for continuation lines.
	// Default is "....> ".
	REPLMultilinePrompt string
}

// NewCLI creates a new CLI instance with the given configuration.
func NewCLI(config Config) *CLI {
	// Set defaults
	if config.OnVMCreate == nil {
		config.OnVMCreate = func() *wrengo.WrenVM {
			return wrengo.NewVMWithForeign()
		}
	}
	if config.ScriptExtension == "" {
		config.ScriptExtension = ".wren"
	}
	if config.REPLPrompt == "" {
		config.REPLPrompt = "wren> "
	}
	if config.REPLMultilinePrompt == "" {
		config.REPLMultilinePrompt = "....> "
	}

	return &CLI{config: config}
}

// Run executes the CLI with the given arguments.
// args[0] is expected to be the program name.
func (c *CLI) Run(args []string) error {
	if len(args) < 2 {
		return c.StartREPL()
	}

	command := args[1]

	switch command {
	case "repl":
		return c.StartREPL()
	case "run":
		if len(args) < 3 {
			return fmt.Errorf("usage: %s run <script.wren>", args[0])
		}
		return c.RunScript(args[2])
	case "eval", "-e":
		if len(args) < 3 {
			return fmt.Errorf("usage: %s eval <code>", args[0])
		}
		return c.RunCode(args[2])
	case "help", "-h", "--help":
		c.PrintHelp(args[0])
		return nil
	case "version", "-v", "--version":
		c.PrintVersion()
		return nil
	default:
		// Assume it's a script file if it has the right extension
		if len(command) > len(c.config.ScriptExtension) &&
			command[len(command)-len(c.config.ScriptExtension):] == c.config.ScriptExtension {
			return c.RunScript(command)
		}
		return fmt.Errorf("unknown command: %s\nTry '%s help' for more information.", command, args[0])
	}
}

// PrintHelp prints the help message.
func (c *CLI) PrintHelp(progName string) {
	fmt.Printf(`Wren CLI - A command-line interface for the Wren scripting language

Usage:
  %s [command] [arguments]

Commands:
  repl              Start interactive REPL
  run <script>      Run a Wren script file
  eval <code>       Evaluate Wren code directly
  help              Show this help message
  version           Show version information

Examples:
  %s repl                    # Start REPL
  %s run script.wren         # Run a script
  %s script.wren             # Run a script (shorthand)
  %s eval "System.print(42)" # Evaluate code

If no command is given, REPL is started.
`, progName, progName, progName, progName, progName)
}

// PrintVersion prints version information.
func (c *CLI) PrintVersion() {
	versionNum := wrengo.GetVersionNumber()
	major := versionNum / 1000000
	minor := (versionNum / 1000) % 1000
	patch := versionNum % 1000

	fmt.Printf("Wren CLI\n")
	fmt.Printf("Wren version: %d.%d.%d\n", major, minor, patch)
	fmt.Printf("Go bindings: github.com/snowmerak/wren.go\n")
}

// RunScript executes a Wren script from a file.
func (c *CLI) RunScript(path string) error {
	// Read the file
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Create VM
	vm := c.config.OnVMCreate()
	defer vm.Free()

	// Execute
	result, err := vm.Interpret(path, string(content))
	if err != nil {
		return fmt.Errorf("execution error: %w", err)
	}

	if result != wrengo.ResultSuccess {
		return fmt.Errorf("script failed with result code: %d", result)
	}

	return nil
}

// RunCode executes Wren code directly.
func (c *CLI) RunCode(code string) error {
	// Create VM
	vm := c.config.OnVMCreate()
	defer vm.Free()

	// Execute
	result, err := vm.Interpret("eval", code)
	if err != nil {
		return fmt.Errorf("execution error: %w", err)
	}

	if result != wrengo.ResultSuccess {
		return fmt.Errorf("code evaluation failed with result code: %d", result)
	}

	return nil
}
