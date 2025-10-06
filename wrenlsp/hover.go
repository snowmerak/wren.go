package wrenlsp

import (
	"fmt"
	"strings"
)

// handleHover handles the textDocument/hover request.
func (s *Server) handleHover(id interface{}, msg map[string]interface{}) map[string]interface{} {
	params, ok := msg["params"].(map[string]interface{})
	if !ok {
		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      id,
			"result":  nil,
		}
	}

	textDocument, ok := params["textDocument"].(map[string]interface{})
	if !ok {
		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      id,
			"result":  nil,
		}
	}

	uri, _ := textDocument["uri"].(string)
	position, _ := params["position"].(map[string]interface{})
	line, _ := position["line"].(float64)
	character, _ := position["character"].(float64)

	// Get the word at the cursor position
	word := s.getWordAtPosition(uri, int(line), int(character))
	if word == "" {
		return map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      id,
			"result":  nil,
		}
	}

	// Search in foreign methods
	for _, fm := range s.config.ForeignMethods {
		if fm.Method == word || fm.Class == word {
			content := fmt.Sprintf("**%s.%s**", fm.Class, fm.Method)
			if fm.Signature != "" {
				sig := fmt.Sprintf("```wren\n%s\n```", fm.Signature)
				content = fmt.Sprintf("**%s.%s**\n\n%s", fm.Class, fm.Method, sig)
			}
			if fm.Doc != "" {
				content = content + "\n\n" + fm.Doc
			}

			return map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      id,
				"result": map[string]interface{}{
					"contents": map[string]interface{}{
						"kind":  "markdown",
						"value": content,
					},
				},
			}
		}
	}

	return map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  nil,
	}
}

// getWordAtPosition extracts the word at the given position in the document.
func (s *Server) getWordAtPosition(uri string, line, character int) string {
	doc, exists := s.documents[uri]
	if !exists {
		return ""
	}

	lines := strings.Split(doc.Content, "\n")
	if line < 0 || line >= len(lines) {
		return ""
	}

	lineText := lines[line]
	if character < 0 || character > len(lineText) {
		return ""
	}

	// Find word boundaries
	start := character
	end := character

	// Go backward to find start
	for start > 0 && isWordChar(lineText[start-1]) {
		start--
	}

	// Go forward to find end
	for end < len(lineText) && isWordChar(lineText[end]) {
		end++
	}

	return lineText[start:end]
}

// isWordChar checks if a character is part of a word.
func isWordChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}
