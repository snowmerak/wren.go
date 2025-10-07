package main

import (
	"fmt"
	"os"

	_ "github.com/snowmerak/gwen/builtin"
	"github.com/snowmerak/gwen/wrenlsp"
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
