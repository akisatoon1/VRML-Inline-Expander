package expander

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/akisatoon1/VRML-Inline-Expander/internal/parser"
	"github.com/akisatoon1/VRML-Inline-Expander/internal/reader" // TODO: readerやwriterは自前で実装する必要ある？
	"github.com/akisatoon1/VRML-Inline-Expander/internal/writer"
)

// ProcessingStack manages files currently being processed to detect circular references
type ProcessingStack struct {
	files map[string]bool
}

// newProcessingStack creates a new ProcessingStack
func newProcessingStack() *ProcessingStack {
	return &ProcessingStack{
		files: make(map[string]bool),
	}
}

// Enter adds a file to the processing stack
// Returns error if the file is already being processed (circular reference)
func (ps *ProcessingStack) Enter(path string) error {
	if ps.files[path] {
		return fmt.Errorf("circular reference detected: %s", path)
	}
	ps.files[path] = true
	return nil
}

// Exit removes a file from the processing stack
func (ps *ProcessingStack) Exit(path string) {
	delete(ps.files, path)
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
	// Initialize processing stack for circular reference detection
	stack := newProcessingStack()

	expandedLines, err := e.expandInlineNodes(inputPath, stack)
	if err != nil {
		return fmt.Errorf("failed to expand Inline nodes: %w", err)
	}

	if err := e.writer.Write(outputPath, expandedLines); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}

func (e *Expander) expandInlineNodes(path string, stack *ProcessingStack) ([]string, error) {
	lines, err := e.reader.Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %w", err)
	}

	// Circular reference guard - check if already processing this file
	if err := stack.Enter(path); err != nil {
		return nil, err
	}
	// Exit on function completion
	defer stack.Exit(path)

	nodes, err := e.parser.FindInlineNodes(lines)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Inline nodes: %w", err)
	}

	// Read referenced files for each node
	baseDir := filepath.Dir(path)
	nodesWithContent := make([]nodeWithContent, len(nodes))
	for i, node := range nodes {
		refPath := filepath.Join(baseDir, node.UrlPath)

		// Recursive call with stack
		refLines, err := e.expandInlineNodes(refPath, stack)
		if err != nil {
			return nil, fmt.Errorf("failed to read referenced file '%s': %w", node.UrlPath, err)
		}

		// remove VRML header if present
		if len(refLines) > 0 && strings.HasPrefix(strings.TrimSpace(refLines[0]), "#VRML") {
			refLines = refLines[1:]
		}

		nodesWithContent[i] = nodeWithContent{
			Node:     node,
			RefLines: refLines,
		}
	}

	expandedLines := e.expandLines(lines, nodesWithContent)
	return expandedLines, nil
}

// expandLines expands Inline nodes in the given lines (testable)
func (e *Expander) expandLines(lines []string, nodesWithContent []nodeWithContent) []string {
	// Process nodes in reverse order to avoid line number shifts
	result := lines
	for i := len(nodesWithContent) - 1; i >= 0; i-- {
		nwc := nodesWithContent[i]

		groupLines := e.buildGroupNode(nwc.Node, nwc.RefLines)
		result = e.replaceLines(result, nwc.Node.StartLine, nwc.Node.EndLine, groupLines)
	}

	return result
}

// TODO: Expanderのメンバを使っていないけど、receiverにする必要あるの？
// buildGroupNode builds a Group node with children from the referenced file content
func (e *Expander) buildGroupNode(node parser.InlineNode, refLines []string) []string {
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

// TODO: Expanderのメンバを使っていないけど、receiverにする必要あるの？
// replaceLines replaces lines from startLine to endLine with newLines
func (e *Expander) replaceLines(lines []string, startLine, endLine int, newLines []string) []string {
	var result []string

	// Add lines before the replacement
	result = append(result, lines[:startLine]...)

	// Add the new lines
	result = append(result, newLines...)

	// Add lines after the replacement
	if endLine+1 < len(lines) {
		result = append(result, lines[endLine+1:]...)
	}

	return result
}
