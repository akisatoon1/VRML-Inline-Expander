package main

import (
	"fmt"
	"os"

	"github.com/akisatoon1/VRML-Inline-Expander/internal/expander/lineexpander"
)

// InlineExpander is an interface for expanding Inline nodes in VRML files
type InlineExpander interface {
	Expand(inputPath, outputPath string) error
}

func main() {
	// Check command line arguments
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <inputfile> <outputfile>\n", os.Args[0])
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	// Create expander and execute
	var exp InlineExpander = lineexpander.New()
	err := exp.Expand(inputPath, outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully expanded %s to %s\n", inputPath, outputPath)
}
