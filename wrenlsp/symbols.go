package wrenlsp

import "regexp"

// extractSymbols extracts class and method names from the document using regex.
func (s *Server) extractSymbols(uri string) []map[string]interface{} {
	doc, exists := s.documents[uri]
	if !exists {
		return nil
	}

	symbols := []map[string]interface{}{}
	seen := make(map[string]bool)

	// Extract class names: class ClassName { ... }
	classRegex := regexp.MustCompile(`class\s+([A-Z][a-zA-Z0-9_]*)`)
	classMatches := classRegex.FindAllStringSubmatch(doc.Content, -1)
	for _, match := range classMatches {
		if len(match) > 1 {
			className := match[1]
			if !seen[className] {
				symbols = append(symbols, map[string]interface{}{
					"label": className,
					"kind":  5, // Class
				})
				seen[className] = true
			}
		}
	}

	// Extract method names: methodName(...) { ... } or methodName { ... }
	methodRegex := regexp.MustCompile(`\s+([a-z][a-zA-Z0-9_]*)\s*(?:\([^)]*\))?\s*\{`)
	methodMatches := methodRegex.FindAllStringSubmatch(doc.Content, -1)
	for _, match := range methodMatches {
		if len(match) > 1 {
			methodName := match[1]
			// Skip keywords
			if isKeyword(methodName) {
				continue
			}
			if !seen[methodName] {
				symbols = append(symbols, map[string]interface{}{
					"label": methodName,
					"kind":  6, // Method
				})
				seen[methodName] = true
			}
		}
	}

	// Extract variable names: var varName = ...
	varRegex := regexp.MustCompile(`var\s+([a-z][a-zA-Z0-9_]*)`)
	varMatches := varRegex.FindAllStringSubmatch(doc.Content, -1)
	for _, match := range varMatches {
		if len(match) > 1 {
			varName := match[1]
			if !seen[varName] {
				symbols = append(symbols, map[string]interface{}{
					"label": varName,
					"kind":  13, // Variable
				})
				seen[varName] = true
			}
		}
	}

	return symbols
}

// isKeyword checks if a word is a Wren keyword.
func isKeyword(word string) bool {
	keywords := map[string]bool{
		"if": true, "else": true, "while": true, "for": true, "in": true,
		"return": true, "break": true, "continue": true, "import": true,
		"static": true, "is": true, "true": true, "false": true,
		"null": true, "this": true, "super": true, "class": true,
		"var": true, "construct": true, "foreign": true,
	}
	return keywords[word]
}
