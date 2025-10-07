package wrenlsp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	wrengo "github.com/snowmerak/gwen"
)

// Document represents an open document.
type Document struct {
	URI     string
	Content string
	Version int
}

// Server represents an LSP server instance.
type Server struct {
	config    Config
	reader    *bufio.Reader
	writer    io.Writer
	documents map[string]*Document
}

// Config defines the configuration for the LSP server.
type Config struct {
	// OnVMCreate creates a custom VM instance for analysis
	OnVMCreate func() *wrengo.WrenVM

	// ForeignMethods contains information about available foreign methods
	ForeignMethods []ForeignMethodInfo

	// EnableDiagnostics enables syntax error diagnostics
	EnableDiagnostics bool
}

// ForeignMethodInfo describes a foreign method for autocompletion.
type ForeignMethodInfo struct {
	Module    string
	Class     string
	Method    string
	Signature string
	IsStatic  bool
	Doc       string
}

// NewServer creates a new LSP server instance.
func NewServer(config Config) *Server {
	// Set defaults
	if config.OnVMCreate == nil {
		config.OnVMCreate = func() *wrengo.WrenVM {
			return wrengo.NewVMWithForeign()
		}
	}
	if config.EnableDiagnostics {
		config.EnableDiagnostics = true
	}

	// Add built-in Wren symbols
	config.ForeignMethods = append(config.ForeignMethods, builtinSymbols...)

	return &Server{
		config:    config,
		reader:    bufio.NewReader(os.Stdin),
		writer:    os.Stdout,
		documents: make(map[string]*Document),
	}
}

// Serve starts the LSP server and processes requests.
func (s *Server) Serve() error {
	for {
		// Read message
		msg, err := s.readMessage()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("failed to read message: %w", err)
		}

		// Handle message
		response := s.handleMessage(msg)
		if response != nil {
			if err := s.writeMessage(response); err != nil {
				return fmt.Errorf("failed to write response: %w", err)
			}
		}
	}
}

// readMessage reads an LSP message from stdin.
func (s *Server) readMessage() (map[string]interface{}, error) {
	// Read headers
	headers := make(map[string]string)
	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	// Get content length
	contentLengthStr, ok := headers["Content-Length"]
	if !ok {
		return nil, fmt.Errorf("missing Content-Length header")
	}

	contentLength, err := strconv.Atoi(contentLengthStr)
	if err != nil {
		return nil, fmt.Errorf("invalid Content-Length: %w", err)
	}

	// Read content
	content := make([]byte, contentLength)
	_, err = io.ReadFull(s.reader, content)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var msg map[string]interface{}
	if err := json.Unmarshal(content, &msg); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return msg, nil
}

// writeMessage writes an LSP message to stdout.
func (s *Server) writeMessage(msg map[string]interface{}) error {
	content, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(content))

	if _, err := s.writer.Write([]byte(header)); err != nil {
		return err
	}

	if _, err := s.writer.Write(content); err != nil {
		return err
	}

	return nil
}

// handleMessage processes an LSP request and returns a response.
func (s *Server) handleMessage(msg map[string]interface{}) map[string]interface{} {
	method, ok := msg["method"].(string)
	if !ok {
		return nil
	}

	id := msg["id"]

	switch method {
	case "initialize":
		return s.handleInitialize(id, msg)
	case "initialized":
		return nil // Notification, no response needed
	case "shutdown":
		return s.handleShutdown(id)
	case "exit":
		os.Exit(0)
		return nil
	case "textDocument/didOpen":
		return s.handleDidOpen(msg)
	case "textDocument/didChange":
		return s.handleDidChange(msg)
	case "textDocument/completion":
		return s.handleCompletion(id, msg)
	case "textDocument/hover":
		return s.handleHover(id, msg)
	default:
		return s.errorResponse(id, -32601, "Method not found")
	}
}

// handleInitialize handles the initialize request.
func (s *Server) handleInitialize(id interface{}, msg map[string]interface{}) map[string]interface{} {
	capabilities := map[string]interface{}{
		"textDocumentSync": map[string]interface{}{
			"openClose": true,
			"change":    1, // Full sync
		},
		"completionProvider": map[string]interface{}{
			"triggerCharacters": []string{".", "("},
		},
		"hoverProvider": true,
	}

	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result": map[string]interface{}{
			"capabilities": capabilities,
			"serverInfo": map[string]interface{}{
				"name":    "wrenlsp",
				"version": "0.1.0",
			},
		},
	}
}

// handleShutdown handles the shutdown request.
func (s *Server) handleShutdown(id interface{}) map[string]interface{} {
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  nil,
	}
}

// handleDidOpen handles the textDocument/didOpen notification.
func (s *Server) handleDidOpen(msg map[string]interface{}) map[string]interface{} {
	params, ok := msg["params"].(map[string]interface{})
	if !ok {
		return nil
	}

	textDocument, ok := params["textDocument"].(map[string]interface{})
	if !ok {
		return nil
	}

	uri, _ := textDocument["uri"].(string)
	content, _ := textDocument["text"].(string)
	version, _ := textDocument["version"].(float64)

	// Store document
	s.documents[uri] = &Document{
		URI:     uri,
		Content: content,
		Version: int(version),
	}

	// Analyze and send diagnostics
	s.publishDiagnostics(uri, content)

	return nil
}

// handleDidChange handles the textDocument/didChange notification.
func (s *Server) handleDidChange(msg map[string]interface{}) map[string]interface{} {
	params, ok := msg["params"].(map[string]interface{})
	if !ok {
		return nil
	}

	textDocument, ok := params["textDocument"].(map[string]interface{})
	if !ok {
		return nil
	}

	uri, _ := textDocument["uri"].(string)
	version, _ := textDocument["version"].(float64)

	contentChanges, ok := params["contentChanges"].([]interface{})
	if !ok || len(contentChanges) == 0 {
		return nil
	}

	// Full document sync - get the new content
	change, ok := contentChanges[0].(map[string]interface{})
	if !ok {
		return nil
	}

	content, _ := change["text"].(string)

	// Update document
	if doc, exists := s.documents[uri]; exists {
		doc.Content = content
		doc.Version = int(version)
	} else {
		s.documents[uri] = &Document{
			URI:     uri,
			Content: content,
			Version: int(version),
		}
	}

	// Analyze and send diagnostics
	s.publishDiagnostics(uri, content)

	return nil
}

// handleCompletion handles the textDocument/completion request.
func (s *Server) handleCompletion(id interface{}, msg map[string]interface{}) map[string]interface{} {
	completionItems := []map[string]interface{}{}

	// Add foreign methods
	for _, fm := range s.config.ForeignMethods {
		item := map[string]interface{}{
			"label":  fm.Method,
			"kind":   2, // Method
			"detail": fmt.Sprintf("%s.%s", fm.Class, fm.Method),
		}
		if fm.Doc != "" {
			item["documentation"] = fm.Doc
		}
		completionItems = append(completionItems, item)
	}

	// Add Wren keywords
	keywords := []string{
		"class", "var", "if", "else", "while", "for", "in",
		"return", "break", "continue", "import", "static",
		"is", "true", "false", "null", "this", "super",
		"construct", "foreign",
	}
	for _, kw := range keywords {
		completionItems = append(completionItems, map[string]interface{}{
			"label": kw,
			"kind":  14, // Keyword
		})
	}

	// Add symbols from current document
	params, ok := msg["params"].(map[string]interface{})
	if ok {
		textDocument, ok := params["textDocument"].(map[string]interface{})
		if ok {
			uri, _ := textDocument["uri"].(string)
			symbols := s.extractSymbols(uri)
			for _, symbol := range symbols {
				completionItems = append(completionItems, symbol)
			}
		}
	}

	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  completionItems,
	}
}

// errorResponse creates an error response.
func (s *Server) errorResponse(id interface{}, code int, message string) map[string]interface{} {
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}
}

// RegisterForeignMethod registers a foreign method for autocompletion.
func (s *Server) RegisterForeignMethod(info ForeignMethodInfo) {
	s.config.ForeignMethods = append(s.config.ForeignMethods, info)
}

// publishDiagnostics analyzes the document and publishes diagnostics.
func (s *Server) publishDiagnostics(uri, content string) {
	if !s.config.EnableDiagnostics {
		return
	}

	diagnostics := []map[string]interface{}{}

	// Create a VM to check syntax
	vm := s.config.OnVMCreate()
	if vm == nil {
		return
	}
	defer vm.Free()

	// Capture errors using Interpret result
	result, _ := vm.Interpret("main", content)

	// If there's a compile error, create a diagnostic
	// Note: We can't get exact line/column from current API,
	// so we'll create a general diagnostic
	if result == wrengo.ResultCompileError {
		diagnostic := map[string]interface{}{
			"range": map[string]interface{}{
				"start": map[string]interface{}{
					"line":      0,
					"character": 0,
				},
				"end": map[string]interface{}{
					"line":      0,
					"character": 1,
				},
			},
			"severity": 1, // Error
			"source":   "wrenlsp",
			"message":  "Syntax error in Wren code",
		}
		diagnostics = append(diagnostics, diagnostic)
	}

	// Send diagnostics notification
	notification := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "textDocument/publishDiagnostics",
		"params": map[string]interface{}{
			"uri":         uri,
			"diagnostics": diagnostics,
		},
	}

	s.writeMessage(notification)
}
