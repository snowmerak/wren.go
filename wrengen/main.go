package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var (
		dir    string
		output string
	)

	flag.StringVar(&dir, "dir", ".", "Directory to scan for Go files")
	flag.StringVar(&output, "output", "", "Output file name (default: <package>_wren.go)")
	flag.Parse()

	if err := run(dir, output); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(dir, output string) error {
	// Parse all Go files in the directory
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, func(fi os.FileInfo) bool {
		// Skip generated files and test files
		name := fi.Name()
		return !strings.HasSuffix(name, "_test.go") &&
			!strings.HasSuffix(name, "_wren.go") &&
			strings.HasSuffix(name, ".go")
	}, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse directory: %w", err)
	}

	// Process each package
	for pkgName, pkg := range pkgs {
		bindings := extractBindings(fset, pkg)
		if len(bindings) == 0 {
			continue
		}

		// Generate code
		code := generateCode(pkgName, bindings)

		// Determine output file name
		outFile := output
		if outFile == "" {
			outFile = filepath.Join(dir, pkgName+"_wren.go")
		}

		// Write to file
		if err := os.WriteFile(outFile, []byte(code), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}

		fmt.Printf("Generated: %s\n", outFile)
	}

	return nil
}

// Binding represents a Go function/method that should be bound to Wren
type Binding struct {
	TypeName   string   // Empty for top-level functions
	MethodName string   // Go method/function name
	WrenName   string   // Name in Wren (default: lowercase first letter)
	Module     string   // Wren module name
	ClassName  string   // Wren class name (default: TypeName or inferred)
	IsStatic   bool     // Whether it's a static method
	Params     []Param  // Parameter types
	Results    []Result // Return types
	RecvType   string   // Receiver type for methods
}

type Param struct {
	Name string
	Type string
}

type Result struct {
	Type string
}

func extractBindings(fset *token.FileSet, pkg *ast.Package) []*Binding {
	var bindings []*Binding
	defaultModule := "main"

	for _, file := range pkg.Files {
		// First pass: find type declarations with //wren:bind to set default module
		for _, decl := range file.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
				if genDecl.Doc != nil {
					for _, comment := range genDecl.Doc.List {
						if binding := parseWrenBindComment(comment.Text); binding != nil {
							if binding.Module != "" {
								defaultModule = binding.Module
							}
						}
					}
				}
			}
		}

		// Second pass: extract function/method bindings
		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				if funcDecl.Doc == nil {
					continue
				}

				var bindingTag *BindingTag
				for _, comment := range funcDecl.Doc.List {
					if tag := parseWrenBindComment(comment.Text); tag != nil {
						bindingTag = tag
						break
					}
				}

				if bindingTag == nil {
					continue
				}

				binding := &Binding{
					MethodName: funcDecl.Name.Name,
					WrenName:   bindingTag.Name,
					Module:     bindingTag.Module,
					ClassName:  bindingTag.Class,
					IsStatic:   bindingTag.Static,
				}

				if binding.Module == "" {
					binding.Module = defaultModule
				}

				// Extract receiver type for methods
				if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
					recvType := funcDecl.Recv.List[0].Type
					binding.RecvType = typeToString(recvType)
					if binding.ClassName == "" {
						binding.ClassName = extractTypeName(binding.RecvType)
					}
				}

				// Use function name if no class specified
				if binding.ClassName == "" && funcDecl.Recv == nil {
					// Top-level function - infer class from function name or use "Lib"
					binding.ClassName = "Lib"
					binding.IsStatic = true
				}

				// Set Wren name (camelCase by default)
				if binding.WrenName == "" {
					binding.WrenName = toWrenMethodName(funcDecl.Name.Name)
				}

				// Extract parameters
				if funcDecl.Type.Params != nil {
					for _, field := range funcDecl.Type.Params.List {
						paramType := typeToString(field.Type)
						for _, name := range field.Names {
							binding.Params = append(binding.Params, Param{
								Name: name.Name,
								Type: paramType,
							})
						}
					}
				}

				// Extract results
				if funcDecl.Type.Results != nil {
					for _, field := range funcDecl.Type.Results.List {
						resultType := typeToString(field.Type)
						binding.Results = append(binding.Results, Result{
							Type: resultType,
						})
					}
				}

				bindings = append(bindings, binding)
			}
		}
	}

	return bindings
}

type BindingTag struct {
	Module string
	Class  string
	Name   string
	Static bool
}

func parseWrenBindComment(comment string) *BindingTag {
	comment = strings.TrimSpace(comment)
	if !strings.HasPrefix(comment, "//wren:bind") {
		return nil
	}

	tag := &BindingTag{}
	parts := strings.Fields(comment)

	for _, part := range parts[1:] { // Skip "//wren:bind"
		if strings.Contains(part, "=") {
			kv := strings.SplitN(part, "=", 2)
			switch kv[0] {
			case "module":
				tag.Module = kv[1]
			case "class":
				tag.Class = kv[1]
			case "name":
				tag.Name = kv[1]
			}
		} else if part == "static" {
			tag.Static = true
		}
	}

	return tag
}

func typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeToString(t.X)
	case *ast.ArrayType:
		return "[]" + typeToString(t.Elt)
	case *ast.MapType:
		return "map[" + typeToString(t.Key) + "]" + typeToString(t.Value)
	case *ast.SelectorExpr:
		return typeToString(t.X) + "." + t.Sel.Name
	default:
		return "unknown"
	}
}

func extractTypeName(typeStr string) string {
	// Remove pointer and package prefixes
	typeStr = strings.TrimPrefix(typeStr, "*")
	if idx := strings.LastIndex(typeStr, "."); idx != -1 {
		typeStr = typeStr[idx+1:]
	}
	return typeStr
}

func toWrenMethodName(goName string) string {
	// Convert Go method name to Wren-style camelCase
	if len(goName) == 0 {
		return goName
	}
	// Keep first letter lowercase unless it's an acronym
	runes := []rune(goName)
	runes[0] = []rune(strings.ToLower(string(runes[0])))[0]
	return string(runes)
}

func generateCode(pkgName string, bindings []*Binding) string {
	var sb strings.Builder

	// Header
	sb.WriteString("// Code generated by wrengen. DO NOT EDIT.\n\n")
	sb.WriteString(fmt.Sprintf("package %s\n\n", pkgName))
	sb.WriteString("import wrengo \"github.com/snowmerak/gwen\"\n\n")

	// Init function
	sb.WriteString("func init() {\n")
	sb.WriteString("\tRegisterWrenBindings()\n")
	sb.WriteString("}\n\n")

	// Registration function
	sb.WriteString("func RegisterWrenBindings() {\n")

	for _, binding := range bindings {
		sb.WriteString(generateBinding(binding))
	}

	sb.WriteString("}\n")

	return sb.String()
}

func generateBinding(b *Binding) string {
	var sb strings.Builder

	// Comment
	sb.WriteString(fmt.Sprintf("\t// %s.%s\n", b.ClassName, b.WrenName))

	// Generate signature
	signature := generateSignature(b)

	// Register foreign method
	sb.WriteString(fmt.Sprintf("\twrengo.RegisterForeignMethod(%q, %q, %v, %q, func(vm *wrengo.WrenVM) {\n",
		b.Module, b.ClassName, b.IsStatic, signature))

	// Extract parameters from slots
	slotIndex := 1
	for _, param := range b.Params {
		sb.WriteString(generateParamExtraction(param, slotIndex))
		slotIndex++
	}

	// Call Go function
	sb.WriteString(generateFunctionCall(b))

	// Handle results
	if len(b.Results) > 0 {
		sb.WriteString(generateResultHandling(b))
	}

	sb.WriteString("\t})\n\n")

	return sb.String()
}

func generateSignature(b *Binding) string {
	// Generate parameter placeholders
	var params []string
	for range b.Params {
		params = append(params, "_")
	}

	paramStr := ""
	if len(params) > 0 {
		paramStr = "(" + strings.Join(params, ",") + ")"
	}

	// Don't include "static" in the signature, just return the method name with params
	return b.WrenName + paramStr
}

func generateParamExtraction(p Param, slotIndex int) string {
	switch {
	case isNumericType(p.Type):
		return fmt.Sprintf("\t\t%s := %s(vm.GetSlotDouble(%d))\n", p.Name, p.Type, slotIndex)
	case p.Type == "string":
		return fmt.Sprintf("\t\t%s := vm.GetSlotString(%d)\n", p.Name, slotIndex)
	case p.Type == "bool":
		return fmt.Sprintf("\t\t%s := vm.GetSlotBool(%d)\n", p.Name, slotIndex)
	case p.Type == "float64":
		return fmt.Sprintf("\t\t%s := vm.GetSlotDouble(%d)\n", p.Name, slotIndex)
	case p.Type == "float32":
		return fmt.Sprintf("\t\t%s := float32(vm.GetSlotDouble(%d))\n", p.Name, slotIndex)
	default:
		return fmt.Sprintf("\t\t// TODO: Extract %s (%s) from slot %d\n", p.Name, p.Type, slotIndex)
	}
}

func generateFunctionCall(b *Binding) string {
	var sb strings.Builder

	// Build parameter list
	var params []string
	for _, p := range b.Params {
		params = append(params, p.Name)
	}

	// Generate call
	if b.RecvType != "" {
		// Method call - need to create receiver
		receiverVar := "receiver"
		sb.WriteString(fmt.Sprintf("\t\t%s := &%s{}\n", receiverVar, extractTypeName(b.RecvType)))

		if len(b.Results) > 0 {
			resultVars := make([]string, len(b.Results))
			for i := range b.Results {
				if i == 0 {
					resultVars[i] = "result"
				} else {
					resultVars[i] = fmt.Sprintf("result%d", i)
				}
			}
			sb.WriteString(fmt.Sprintf("\t\t%s := %s.%s(%s)\n",
				strings.Join(resultVars, ", "), receiverVar, b.MethodName, strings.Join(params, ", ")))
		} else {
			sb.WriteString(fmt.Sprintf("\t\t%s.%s(%s)\n", receiverVar, b.MethodName, strings.Join(params, ", ")))
		}
	} else {
		// Function call
		if len(b.Results) > 0 {
			resultVars := make([]string, len(b.Results))
			for i := range b.Results {
				if i == 0 {
					resultVars[i] = "result"
				} else {
					resultVars[i] = fmt.Sprintf("result%d", i)
				}
			}
			sb.WriteString(fmt.Sprintf("\t\t%s := %s(%s)\n",
				strings.Join(resultVars, ", "), b.MethodName, strings.Join(params, ", ")))
		} else {
			sb.WriteString(fmt.Sprintf("\t\t%s(%s)\n", b.MethodName, strings.Join(params, ", ")))
		}
	}

	return sb.String()
}

func generateResultHandling(b *Binding) string {
	var sb strings.Builder

	// Handle error if last result is error
	if len(b.Results) > 0 && b.Results[len(b.Results)-1].Type == "error" {
		errorVar := "result"
		if len(b.Results) > 1 {
			errorVar = fmt.Sprintf("result%d", len(b.Results)-1)
		}

		sb.WriteString(fmt.Sprintf("\t\tif %s != nil {\n", errorVar))
		sb.WriteString(fmt.Sprintf("\t\t\tvm.SetSlotString(0, %s.Error())\n", errorVar))
		sb.WriteString("\t\t\tvm.AbortFiber(0)\n")
		sb.WriteString("\t\t\treturn\n")
		sb.WriteString("\t\t}\n")

		// Set non-error result
		if len(b.Results) > 1 {
			sb.WriteString(generateResultSetter(b.Results[0], "result"))
		}
	} else if len(b.Results) > 0 {
		// Just set the result
		sb.WriteString(generateResultSetter(b.Results[0], "result"))
	}

	return sb.String()
}

func generateResultSetter(r Result, varName string) string {
	switch {
	case isNumericType(r.Type):
		return fmt.Sprintf("\t\tvm.SetSlotDouble(0, float64(%s))\n", varName)
	case r.Type == "string":
		return fmt.Sprintf("\t\tvm.SetSlotString(0, %s)\n", varName)
	case r.Type == "bool":
		return fmt.Sprintf("\t\tvm.SetSlotBool(0, %s)\n", varName)
	case r.Type == "float64":
		return fmt.Sprintf("\t\tvm.SetSlotDouble(0, %s)\n", varName)
	case r.Type == "float32":
		return fmt.Sprintf("\t\tvm.SetSlotDouble(0, float64(%s))\n", varName)
	default:
		return fmt.Sprintf("\t\t// TODO: Set result %s (%s)\n", varName, r.Type)
	}
}

func isNumericType(t string) bool {
	return t == "int" || t == "int8" || t == "int16" || t == "int32" || t == "int64" ||
		t == "uint" || t == "uint8" || t == "uint16" || t == "uint32" || t == "uint64" ||
		t == "float32" || t == "float64"
}
