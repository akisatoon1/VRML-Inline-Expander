package reader

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Reader handles file reading operations
type Reader struct{}

// New creates a new Reader instance
func New() *Reader {
	return &Reader{}
}

// Read reads a file and returns its content as a slice of lines
func (r *Reader) Read(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return lines, nil
}

// ReadReferenced reads a referenced file with path resolution relative to baseDir
// and removes the VRML header
func (r *Reader) ReadReferenced(baseDir, relativePath string) ([]string, error) {
	// Check if it's a network path (not supported)
	if strings.HasPrefix(relativePath, "http://") ||
		strings.HasPrefix(relativePath, "https://") ||
		strings.HasPrefix(relativePath, "ftp://") {
		return nil, fmt.Errorf("network paths are not supported: %s", relativePath)
	}

	// Resolve the path
	resolvedPath := filepath.Join(baseDir, relativePath)

	// Read the file
	lines, err := r.Read(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read referenced file %s: %w", relativePath, err)
	}

	// Remove VRML header
	lines = r.removeHeader(lines)

	return lines, nil
}

// removeHeader removes VRML header lines (lines starting with #VRML)
func (r *Reader) removeHeader(lines []string) []string {
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "#VRML") {
			result = append(result, line)
		}
	}
	return result
}
