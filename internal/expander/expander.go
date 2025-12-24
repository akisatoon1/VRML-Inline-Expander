package expander

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/akisatoon1/VRML-Inline-Expander/internal/parser"
	"github.com/akisatoon1/VRML-Inline-Expander/internal/reader" // TODO: readerやwriterは自前で実装する必要ある？
	"github.com/akisatoon1/VRML-Inline-Expander/internal/writer"
)

// processingStack manages files currently being processed to detect circular references
type processingStack struct {
	files map[string]bool
}

// newProcessingStack creates a new processingStack
func newProcessingStack() *processingStack {
	return &processingStack{
		files: make(map[string]bool),
	}
}

// Enter adds a file to the processing stack
// Returns error if the file is already being processed (circular reference)
func (ps *processingStack) Enter(path string) error {
	if ps == nil {
		return fmt.Errorf("processing stack is nil")
	}
	if ps.files[path] {
		return fmt.Errorf("circular reference detected: %s", path)
	}
	ps.files[path] = true
	return nil
}

// Exit removes a file from the processing stack
func (ps *processingStack) Exit(path string) {
	if ps == nil || ps.files == nil {
		return // TODO: 異常終了するべき
	}
	delete(ps.files, path)
}

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

// Expander handles the expansion of Inline nodes
type Expander struct {
	parser *parser.Parser
	reader *reader.Reader
	writer *writer.Writer
}

// New creates a new Expander instance
func New() *Expander {
	return &Expander{
		parser: parser.New(),
		reader: reader.New(),
		writer: writer.New(),
	}
}

// nodeWithContent represents an Inline node with its referenced file content
type nodeWithContent struct {
	Node     parser.InlineNode
	RefLines []string
}

// Expand expands Inline nodes in the input file and writes the result to the output file
func (e *Expander) Expand(inputPath, outputPath string) error {
	if e == nil {
		return fmt.Errorf("expander is nil")
	}
	if e.parser == nil || e.reader == nil || e.writer == nil {
		return fmt.Errorf("expander not properly initialized, use New()")
	}

	// Initialize processing stack for circular reference detection
	stack := newProcessingStack()

	inputAbsPath, err := filepath.Abs(inputPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for input file: %w", err)
	}

	expandedLines, err := e.expandInlineNodes(inputAbsPath, stack)
	if err != nil {
		return fmt.Errorf("failed to expand Inline nodes: %w", err)
	}

	if err := e.writer.Write(outputPath, expandedLines); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

func (e *Expander) expandInlineNodes(absPath string, stack *processingStack) ([]string, error) {
	// Circular reference guard - check if already processing this file
	if err := stack.Enter(absPath); err != nil {
		return nil, err
	}
	// Exit on function completion
	defer stack.Exit(absPath)

	lines, err := e.reader.Read(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %w", err)
	}

	nodes, err := e.parser.FindInlineNodes(lines)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Inline nodes: %w", err)
	}

	// Read referenced files for each node
	nodesWithContent := make([]nodeWithContent, len(nodes))
	for i, node := range nodes {
		nwc, err := e.expandInlineNode(absPath, node, stack)
		if err != nil {
			return nil, err
		}
		nodesWithContent[i] = nwc
	}

	expandedLines, err := e.expandLines(lines, nodesWithContent)
	if err != nil {
		return nil, fmt.Errorf("failed to expand lines: %w", err)
	}
	return expandedLines, nil
}

// expandInlineNode expands a single Inline node
func (e *Expander) expandInlineNode(basePath string, node parser.InlineNode, stack *processingStack) (nodeWithContent, error) {
	refAbsPath, err := resolveAbsolutePath(basePath, node.UrlPath)
	if err != nil {
		return nodeWithContent{}, fmt.Errorf("failed to resolve path for referenced file '%s': %w", node.UrlPath, err)
	}

	// Recursive call with stack
	refLines, err := e.expandInlineNodes(refAbsPath, stack)
	if err != nil {
		return nodeWithContent{}, fmt.Errorf("failed to expand referenced file '%s' from '%s': %w", node.UrlPath, basePath, err)
	}

	// remove VRML header if present
	refLines = removeVRMLHeader(refLines)

	return nodeWithContent{
		Node:     node,
		RefLines: refLines,
	}, nil
}

// expandLines expands Inline nodes in the given lines (testable)
func (e *Expander) expandLines(lines []string, nodesWithContent []nodeWithContent) ([]string, error) {
	// Process nodes in reverse order to avoid line number shifts
	result := lines
	for i := len(nodesWithContent) - 1; i >= 0; i-- {
		nwc := nodesWithContent[i]

		groupLines := buildGroupNode(nwc.Node, nwc.RefLines)
		var err error
		result, err = replaceLines(result, nwc.Node.StartLine, nwc.Node.EndLine, groupLines)
		if err != nil {
			return nil, fmt.Errorf("failed to replace lines for node at line %d-%d: %w", nwc.Node.StartLine, nwc.Node.EndLine, err)
		}
	}

	return result, nil
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
