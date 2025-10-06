package wrenlsp

import (
	"testing"
)

func TestServerCreation(t *testing.T) {
	config := Config{
		EnableDiagnostics: true,
	}

	server := NewServer(config)
	if server == nil {
		t.Fatal("Failed to create server")
	}

	if !server.config.EnableDiagnostics {
		t.Error("EnableDiagnostics should be true")
	}

	if server.config.OnVMCreate == nil {
		t.Error("OnVMCreate should have default value")
	}
}

func TestForeignMethodRegistration(t *testing.T) {
	server := NewServer(Config{})

	info := ForeignMethodInfo{
		Module:    "test",
		Class:     "TestClass",
		Method:    "testMethod",
		Signature: "testMethod()",
		IsStatic:  false,
		Doc:       "Test method documentation",
	}

	server.RegisterForeignMethod(info)

	if len(server.config.ForeignMethods) != 1 {
		t.Errorf("Expected 1 foreign method, got %d", len(server.config.ForeignMethods))
	}

	if server.config.ForeignMethods[0].Method != "testMethod" {
		t.Errorf("Expected method name 'testMethod', got '%s'", server.config.ForeignMethods[0].Method)
	}
}

func TestInitializeResponse(t *testing.T) {
	server := NewServer(Config{})

	msg := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]interface{}{
			"processId": 12345,
			"capabilities": map[string]interface{}{},
		},
	}

	response := server.handleInitialize(1, msg)

	if response["jsonrpc"] != "2.0" {
		t.Error("Response should have jsonrpc 2.0")
	}

	if response["id"] != 1 {
		t.Error("Response should have id 1")
	}

	result, ok := response["result"].(map[string]interface{})
	if !ok {
		t.Fatal("Response should have result")
	}

	capabilities, ok := result["capabilities"].(map[string]interface{})
	if !ok {
		t.Fatal("Result should have capabilities")
	}

	if capabilities["completionProvider"] == nil {
		t.Error("Capabilities should include completionProvider")
	}

	if capabilities["hoverProvider"] != true {
		t.Error("Capabilities should include hoverProvider")
	}
}

func TestCompletionResponse(t *testing.T) {
	server := NewServer(Config{})

	// Register a foreign method
	info := ForeignMethodInfo{
		Module:    "test",
		Class:     "TestClass",
		Method:    "testMethod",
		Signature: "testMethod(_)",
		Doc:       "Test documentation",
	}
	server.RegisterForeignMethod(info)

	msg := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      2,
		"method":  "textDocument/completion",
		"params": map[string]interface{}{
			"textDocument": map[string]interface{}{
				"uri": "file:///test.wren",
			},
			"position": map[string]interface{}{
				"line":      0,
				"character": 5,
			},
		},
	}

	response := server.handleCompletion(2, msg)

	if response["id"] != 2 {
		t.Error("Response should have id 2")
	}

	result, ok := response["result"].([]map[string]interface{})
	if !ok {
		t.Fatal("Result should be an array of completion items")
	}

	// Should have at least the foreign method and some keywords
	if len(result) < 2 {
		t.Errorf("Expected at least 2 completion items, got %d", len(result))
	}

	// Check if foreign method is included
	foundForeignMethod := false
	for _, item := range result {
		if item["label"] == "testMethod" {
			foundForeignMethod = true
			break
		}
	}

	if !foundForeignMethod {
		t.Error("Completion should include registered foreign method")
	}
}
