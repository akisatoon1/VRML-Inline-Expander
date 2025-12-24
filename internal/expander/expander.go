package expander

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/akisatoon1/VRML-Inline-Expander/internal/parser"
	"github.com/akisatoon1/VRML-Inline-Expander/internal/reader" // TODO: readerやwriterは自前で実装する必要ある？
	"github.com/akisatoon1/VRML-Inline-Expander/internal/writer"
)

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
	if err := stack.enter(absPath); err != nil {
		return nil, err
	}
	// Exit on function completion
	defer stack.exit(absPath)

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
	result := lines

	// Process nodes in reverse order to avoid line number shifts
	for _, nwc := range slices.Backward(nodesWithContent) {
		groupLines := buildGroupNode(nwc.Node, nwc.RefLines)

		res, err := replaceLines(result, nwc.Node.StartLine, nwc.Node.EndLine, groupLines)
		if err != nil {
			return nil, fmt.Errorf("failed to replace lines for node at line %d-%d: %w", nwc.Node.StartLine, nwc.Node.EndLine, err)
		}

		result = res
	}

	return result, nil
}
