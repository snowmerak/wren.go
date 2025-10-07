package main

import (
	"testing"

	wrengo "github.com/snowmerak/gwen"
	"github.com/snowmerak/gwen/wrencli"
)

func TestStandardCLI(t *testing.T) {
	// Test that standard CLI can be created
	cli := wrencli.NewCLI(wrencli.Config{
		OnVMCreate: func() *wrengo.WrenVM {
			return wrengo.NewVMWithForeign()
		},
	})

	if cli == nil {
		t.Fatal("Failed to create CLI")
	}

	// Test version command
	err := cli.Run([]string{"wren-std", "version"})
	if err != nil {
		t.Errorf("Version command failed: %v", err)
	}

	// Test help command
	err = cli.Run([]string{"wren-std", "help"})
	if err != nil {
		t.Errorf("Help command failed: %v", err)
	}

	// Test code evaluation
	err = cli.RunCode(`System.print("Hello from wren-std!")`)
	if err != nil {
		t.Errorf("Code evaluation failed: %v", err)
	}

	// Test basic functionality
	t.Log("Standard CLI is working correctly!")
	if err != nil {
		t.Errorf("Async test failed: %v", err)
	}
}
