package expander

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/akisatoon1/VRML-Inline-Expander/internal/parser"
)

// removeVRMLHeader removes the VRML header line if present
func removeVRMLHeader(lines []string) []string {
	if len(lines) > 0 && strings.HasPrefix(strings.TrimSpace(lines[0]), "#VRML") {
		return lines[1:]
	}
	return lines
}

// resolveAbsolutePath resolves a relative path to an absolute path based on the base file path
// baseFilePath: the absolute path of the file containing the reference
// relativePath: the relative path to resolve
func resolveAbsolutePath(baseFilePath, relativePath string) (absPath string, err error) {
	baseDir := filepath.Dir(baseFilePath)
	resolvedPath := filepath.Join(baseDir, relativePath)
	absPath, err = filepath.Abs(resolvedPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for '%s' (base: %s): %w", relativePath, baseFilePath, err)
	}
	return absPath, nil
}

// buildGroupNode builds a Group node with children from the referenced file content
func buildGroupNode(node parser.InlineNode, refLines []string) []string {
	var result []string

	// Build the Group node header
	if node.DefName != "" {
		result = append(result, fmt.Sprintf("DEF %s Group {", node.DefName))
	} else {
		result = append(result, "Group {")
	}

	// Add children field
	result = append(result, "  children [")

	// Add referenced content with indentation
	for _, line := range refLines {
		// Skip empty lines at the beginning and end
		if strings.TrimSpace(line) == "" {
			continue
		}
		result = append(result, "    "+line)
	}

	// Close children and Group
	result = append(result, "  ]")
	result = append(result, "}")

	return result
}

// replaceLines replaces lines in the closed interval [startLine, endLine] with newLines.
// Indices are 0-based and both inclusive.
//
// Example: replaceLines(["A", "B", "C", "D", "E"], 1, 2, ["X", "Y"]) => ["A", "X", "Y", "D", "E"]
func replaceLines(lines []string, startLine, endLine int, newLines []string) ([]string, error) {
	// Validate input
	if len(lines) == 0 {
		return nil, fmt.Errorf("cannot replace lines in empty slice")
	}
	if startLine < 0 {
		return nil, fmt.Errorf("startLine cannot be negative: %d", startLine)
	}
	if endLine < 0 {
		return nil, fmt.Errorf("endLine cannot be negative: %d", endLine)
	}
	if startLine >= len(lines) {
		return nil, fmt.Errorf("startLine (%d) is out of range (valid range: 0-%d)", startLine, len(lines)-1)
	}
	if endLine >= len(lines) {
		return nil, fmt.Errorf("endLine (%d) is out of range (valid range: 0-%d)", endLine, len(lines)-1)
	}
	if endLine < startLine {
		return nil, fmt.Errorf("startLine (%d) cannot be greater than endLine (%d)", startLine, endLine)
	}

	var result []string

	// Add lines before the replacement
	result = append(result, lines[:startLine]...)

	// Add the new lines
	result = append(result, newLines...)

	// Add lines after the replacement
	if endLine+1 < len(lines) {
		result = append(result, lines[endLine+1:]...)
	}

	return result, nil
}
