package main

import (
	"fmt"
	"os"

	_ "github.com/snowmerak/wren.go/builtin" // Import builtin async functions
	wrengo "github.com/snowmerak/wren.go"
	"github.com/snowmerak/wren.go/wrencli"
)

func main() {
	// Create CLI with standard VM (includes async support)
	cli := wrencli.NewCLI(wrencli.Config{
		OnVMCreate: func() *wrengo.WrenVM {
			// Standard VM with all built-in features including async
			return wrengo.NewVMWithForeign()
		},
	})

	// Run CLI
	if err := cli.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
