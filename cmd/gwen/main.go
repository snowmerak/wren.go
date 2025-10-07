package main

import (
	"fmt"
	"os"

	wrengo "github.com/snowmerak/gwen"
	_ "github.com/snowmerak/gwen/builtin" // Import builtin async functions
	"github.com/snowmerak/gwen/wrencli"
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
