package main

import (
	"fmt"
	"os"

	"github.com/snowmerak/wren.go/wrenlsp"
)

func main() {
	config := wrenlsp.Config{
		EnableDiagnostics: true,
	}

	server := wrenlsp.NewServer(config)

	if err := server.Serve(); err != nil {
		fmt.Fprintf(os.Stderr, "LSP server error: %v\n", err)
		os.Exit(1)
	}
}
