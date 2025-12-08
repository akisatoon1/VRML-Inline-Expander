package writer

import (
	"fmt"
	"os"
	"strings"
)

// Writer handles file writing operations
type Writer struct{}

// New creates a new Writer instance
func New() *Writer {
	return &Writer{}
}

// Write writes the content to a file
func (w *Writer) Write(filePath string, lines []string) error {
	// Join lines with newline
	content := strings.Join(lines, "\n")

	// Add final newline if content is not empty
	if len(content) > 0 {
		content += "\n"
	}

	// Write to file
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
