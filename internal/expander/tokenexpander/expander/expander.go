package expander

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/akisatoon1/VRML-Inline-Expander/internal/expander/ancestor"
)

type FileExpander struct {
	tokenizer Tokenizer
	consumer  Consumer
}

func NewFileExpander(tokenizer Tokenizer, consumer Consumer) *FileExpander {
	return &FileExpander{
		tokenizer: tokenizer,
		consumer:  consumer,
	}
}

// Expands recursively a VRML file to the given writer.
func (e *FileExpander) ExpandToWriter(filePath string, writer io.Writer) error {
	// write VRML header
	writer.Write([]byte("#VRML V2.0 utf8\n"))

	ancestors := ancestor.NewEmptyAncestorSet()
	return e.expandToWriter(filePath, ancestors, writer)
}

func (e *FileExpander) expandToWriter(filePath string, ancestors ancestor.AncestorSet, writer io.Writer) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	if err := checkCircularReference(absPath, ancestors); err != nil {
		return err
	}

	tokens, err := e.tokenizeFile(filePath)
	if err != nil {
		return err
	}

	updatedAncestors := ancestors.Add(absPath)
	return e.writeAllTokens(tokens, filePath, updatedAncestors, writer)
}
