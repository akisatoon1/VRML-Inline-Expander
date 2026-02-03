package tokenexpander

import (
	"fmt"
	"os"

	"github.com/akisatoon1/VRML-Inline-Expander/internal/expander/tokenexpander/consumer"
	"github.com/akisatoon1/VRML-Inline-Expander/internal/expander/tokenexpander/expander"
	"github.com/akisatoon1/VRML-Inline-Expander/internal/expander/tokenexpander/token"
)

// Expander implements the InlineExpander interface using token-based expansion
type Expander struct {
	fileExpander *expander.FileExpander
}

// New creates a new Expander instance
func New() *Expander {
	return &Expander{
		fileExpander: expander.NewFileExpander(token.NewTokenizer(), consumer.NewInlineConsumer()),
	}
}

// Expand expands all Inline nodes in a VRML file and writes the result to outputPath
func (e *Expander) Expand(inputPath, outputPath string) error {
	// Create output file
	outputFile, err := e.createOutputFile(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Expand the input file
	if err := e.fileExpander.ExpandToWriter(inputPath, outputFile); err != nil {
		return fmt.Errorf("failed to expand file: %w", err)
	}

	return nil
}

// createOutputFile creates the output file for writing
func (e *Expander) createOutputFile(outputPath string) (*os.File, error) {
	file, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create output file: %w", err)
	}
	return file, nil
}
