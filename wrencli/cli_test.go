package wrencli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	wrengo "github.com/snowmerak/wren.go"
)

func TestNewCLI(t *testing.T) {
	cli := NewCLI(Config{})
	
	if cli == nil {
		t.Fatal("NewCLI returned nil")
	}
	
	// Check defaults
	if cli.config.ScriptExtension != ".wren" {
		t.Errorf("Expected default extension .wren, got %s", cli.config.ScriptExtension)
	}
	
	if cli.config.REPLPrompt != "wren> " {
		t.Errorf("Expected default prompt 'wren> ', got %s", cli.config.REPLPrompt)
	}
	
	if cli.config.OnVMCreate == nil {
		t.Error("Expected OnVMCreate to be set")
	}
}

func TestCustomConfig(t *testing.T) {
	customCalled := false
	
	cli := NewCLI(Config{
		OnVMCreate: func() *wrengo.WrenVM {
			customCalled = true
			return wrengo.NewVM()
		},
		ScriptExtension: ".ws",
		REPLPrompt:      "custom> ",
	})
	
	// Test custom VM creation
	vm := cli.config.OnVMCreate()
	defer vm.Free()
	
	if !customCalled {
		t.Error("Custom OnVMCreate was not called")
	}
	
	if cli.config.ScriptExtension != ".ws" {
		t.Errorf("Expected extension .ws, got %s", cli.config.ScriptExtension)
	}
	
	if cli.config.REPLPrompt != "custom> " {
		t.Errorf("Expected prompt 'custom> ', got %s", cli.config.REPLPrompt)
	}
}

func TestRunCode(t *testing.T) {
	cli := NewCLI(Config{})
	
	// Test simple code
	err := cli.RunCode(`System.print("test")`)
	if err != nil {
		t.Errorf("RunCode failed: %v", err)
	}
}

func TestRunScript(t *testing.T) {
	cli := NewCLI(Config{})
	
	// Create a temporary script
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.wren")
	
	content := `System.print("Hello from script")`
	err := os.WriteFile(scriptPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test script: %v", err)
	}
	
	// Run the script
	err = cli.RunScript(scriptPath)
	if err != nil {
		t.Errorf("RunScript failed: %v", err)
	}
}

func TestRunScriptNotFound(t *testing.T) {
	cli := NewCLI(Config{})
	
	err := cli.RunScript("nonexistent.wren")
	if err == nil {
		t.Error("Expected error for nonexistent script")
	}
	
	if !strings.Contains(err.Error(), "failed to read file") {
		t.Errorf("Expected 'failed to read file' error, got: %v", err)
	}
}

func TestRunInvalidCode(t *testing.T) {
	cli := NewCLI(Config{})
	
	// Invalid Wren syntax
	err := cli.RunCode(`this is not valid wren code at all`)
	if err == nil {
		t.Error("Expected error for invalid code")
	}
}

func TestCommandParsing(t *testing.T) {
	cli := NewCLI(Config{})
	
	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "help command",
			args:        []string{"prog", "help"},
			expectError: false,
		},
		{
			name:        "version command",
			args:        []string{"prog", "version"},
			expectError: false,
		},
		{
			name:        "eval without code",
			args:        []string{"prog", "eval"},
			expectError: true,
		},
		{
			name:        "run without script",
			args:        []string{"prog", "run"},
			expectError: true,
		},
		{
			name:        "unknown command",
			args:        []string{"prog", "unknown"},
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cli.Run(tt.args)
			
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestScriptExtensionDetection(t *testing.T) {
	cli := NewCLI(Config{})
	
	// Create a temporary script
	tmpDir := t.TempDir()
	scriptPath := filepath.Join(tmpDir, "test.wren")
	
	content := `System.print("test")`
	err := os.WriteFile(scriptPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write test script: %v", err)
	}
	
	// Should detect .wren extension and run as script
	err = cli.Run([]string{"prog", scriptPath})
	if err != nil {
		t.Errorf("Failed to run script via extension detection: %v", err)
	}
}
