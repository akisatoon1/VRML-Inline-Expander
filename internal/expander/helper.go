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
	// Group Node layout
	/*
		Group {
		  children [
		    ... // content from referenced file
		  ]
		}
	*/
	const layout = `%sGroup {
  children [
%s
  ]
}`

	defPart := createDefDeclaration(node.DefName)

	childrenBlock := joinWithIndent(refLines, "    ")

	layoutResult := fmt.Sprintf(layout, defPart, childrenBlock)

	return strings.Split(layoutResult, "\n")
}

// createDefDeclaration returns the DEF part of a node declaration
// If defName is empty, returns empty string
// Otherwise, returns "DEF <defName> " with trailing space
func createDefDeclaration(defName string) string {
	if defName == "" {
		return ""
	}
	return fmt.Sprintf("DEF %s ", defName)
}

func joinWithIndent(lines []string, indent string) string {
	var builder strings.Builder
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		builder.WriteString(indent)
		builder.WriteString(line)
		builder.WriteString("\n")
	}
	return strings.TrimRight(builder.String(), "\n")
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
