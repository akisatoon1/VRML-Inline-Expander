package expander

import (
	"fmt"
	"os"
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
	// Check if baseFilePath exists
	if _, err := os.Stat(baseFilePath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("baseFilePath does not exist: %s", baseFilePath)
		}
		return "", fmt.Errorf("failed to check baseFilePath: %w", err)
	}

	// Resolve relative path
	baseDir := filepath.Dir(baseFilePath)
	targetPath := filepath.Join(baseDir, relativePath)

	// Get absolute path
	absPath, err = filepath.Abs(targetPath)
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
	if err := validateReplaceLines(len(lines), startLine, endLine); err != nil {
		return nil, err
	}

	result := performLineReplacement(lines, startLine, endLine, newLines)
	return result, nil
}

// not validation
func performLineReplacement(lines []string, startLine, endLine int, newLines []string) []string {
	var result []string
	result = append(result, lines[:startLine]...)
	result = append(result, newLines...)
	result = append(result, lines[endLine+1:]...)
	return result
}

func validateReplaceLines(linesNum, startLine, endLine int) error {
	if linesNum == 0 {
		return fmt.Errorf("cannot replace lines in empty slice")
	}
	if startLine < 0 {
		return fmt.Errorf("startLine cannot be negative: %d", startLine)
	}
	if endLine < 0 {
		return fmt.Errorf("endLine cannot be negative: %d", endLine)
	}
	if startLine >= linesNum {
		return fmt.Errorf("startLine (%d) is out of range (valid range: 0-%d)", startLine, linesNum-1)
	}
	if endLine >= linesNum {
		return fmt.Errorf("endLine (%d) is out of range (valid range: 0-%d)", endLine, linesNum-1)
	}
	if endLine < startLine {
		return fmt.Errorf("startLine (%d) cannot be greater than endLine (%d)", startLine, endLine)
	}
	return nil
}
