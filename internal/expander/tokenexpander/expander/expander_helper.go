package expander

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/akisatoon1/VRML-Inline-Expander/internal/expander/ancestor"
)

// checkCircularReference checks if the file path exists in ancestors
func checkCircularReference(absPath string, ancestors ancestor.AncestorSet) error {
	if ancestors.Contain(absPath) {
		return fmt.Errorf("circular reference detected: %s", absPath)
	}
	return nil
}

// tokenizeFile opens a file and returns its tokens
func (e *FileExpander) tokenizeFile(filePath string) ([]Token, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	tokens, err := e.tokenizer.Tokenize(file)
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize file: %w", err)
	}

	return tokens, nil
}

// Write all tokens. Find Inline Node in tokens, expand Inline Node and write expanded tokens.
func (e *FileExpander) writeAllTokens(tokens []Token, currentFilePath string, ancestors ancestor.AncestorSet, writer io.Writer) error {
	tokenIndex := 0
	for tokenIndex < len(tokens) {
		inlineNodeFound, inline, tokensConsumed, err := e.consumer.ConsumeInlineNode(tokens, tokenIndex)
		if err != nil {
			return fmt.Errorf("failed to consume inline node: %w", err)
		}

		if inlineNodeFound {
			if err := e.writeExpandedInline(inline, currentFilePath, ancestors, writer); err != nil {
				return err
			}
			tokenIndex += tokensConsumed
		} else {
			if err := writeTokenAndSpace(tokens[tokenIndex], writer); err != nil {
				return err
			}
			tokenIndex++
		}
	}

	return nil
}

// Write expanded Inline Node.
func (e *FileExpander) writeExpandedInline(inline InlineNode, currentFilePath string, ancestors ancestor.AncestorSet, writer io.Writer) error {
	if err := writeGroupNodeStart(writer); err != nil {
		return fmt.Errorf("failed to write group start: %w", err)
	}

	filePath := pathFromWD(currentFilePath, inline.UrlPath())
	if err := e.expandToWriter(filePath, ancestors, writer); err != nil {
		return fmt.Errorf("failed to expand referenced file %s: %w", inline.UrlPath(), err)
	}

	if err := writeGroupNodeEnd(writer); err != nil {
		return fmt.Errorf("failed to write group end: %w", err)
	}

	return nil
}

// Writes the beginning boilerplate of a Group node
func writeGroupNodeStart(writer io.Writer) error {
	if _, err := writer.Write([]byte("Group { children [")); err != nil {
		return err
	}

	return nil
}

// Writes the ending boilerplate of a Group node
func writeGroupNodeEnd(writer io.Writer) error {
	if _, err := writer.Write([]byte("] }")); err != nil {
		return err
	}

	return nil
}

// Resolve referenced file path from current working directory
func pathFromWD(currentFilePath, relativePath string) string {
	baseDir := filepath.Dir(currentFilePath)
	return filepath.Join(baseDir, relativePath)
}

// Write a token and a space
func writeTokenAndSpace(token Token, writer io.Writer) error {
	if _, err := writer.Write([]byte(token.GetValue())); err != nil {
		return err
	}

	if _, err := writer.Write([]byte(" ")); err != nil {
		return err
	}
	return nil
}
