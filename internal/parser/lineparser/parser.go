/*
TODO: 仕様を満たす

lineの配列で解析しているので、複数行にまたがるような文法を正しく認識できない。
vrmlでは空白と改行は同じ区切り文字であり、改行を空白の代わりに使っていても正しくパースする必要あり。
*/

package lineparser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/akisatoon1/VRML-Inline-Expander/internal/parser"
)

// Parser handles VRML parsing operations
type Parser struct{}

// New creates a new Parser instance
func New() *Parser {
	return &Parser{}
}

// FindInlineNodes finds all Inline nodes in the given content
func (p *Parser) FindInlineNodes(lines []string) ([]parser.InlineNode, error) {
	var nodes []parser.InlineNode

	i := 0
	for i < len(lines) {
		line := lines[i]

		// Check if this line contains "Inline"
		if strings.Contains(line, "Inline") {
			node, endLine, err := p.parseInlineNode(lines, i)
			if err != nil {
				return nil, fmt.Errorf("failed to parse Inline node at line %d: %w", i, err)
			}
			nodes = append(nodes, node)
			i = endLine + 1
		} else {
			i++
		}
	}

	return nodes, nil
}

// parseInlineNode parses an Inline node starting from the given line
func (p *Parser) parseInlineNode(lines []string, startLine int) (parser.InlineNode, int, error) {
	node := parser.InlineNode{
		StartLine: startLine,
	}

	// Extract DEF name if present
	defName := p.extractDefName(lines[startLine])
	node.DefName = defName

	// Find the url field
	urlPath, endLine, err := p.findUrlField(lines, startLine)
	if err != nil {
		return node, startLine, err
	}

	node.UrlPath = urlPath
	node.EndLine = endLine

	return node, endLine, nil
}

// extractDefName extracts the DEF name from a line containing Inline
// Example: "DEF MyModel Inline {" -> "MyModel"
// TODO: 区切り文字
func (p *Parser) extractDefName(line string) string {
	// Pattern: DEF <name> Inline
	// TODO: compile
	re := regexp.MustCompile(`DEF\s+(\w+)\s+Inline`)
	matches := re.FindStringSubmatch(line)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

// findUrlField finds the url field and returns the path and the end line of the node
// TODO: Inlineの深さは1まで
func (p *Parser) findUrlField(lines []string, startLine int) (string, int, error) {
	braceDepth := 0
	urlFound := false
	var urlPath string

	for i := startLine; i < len(lines); i++ {
		line := lines[i]

		// Count braces to track node boundaries
		braceDepth += strings.Count(line, "{")
		braceDepth -= strings.Count(line, "}")

		// Look for url field if not found yet
		if !urlFound {
			url := p.extractUrlPath(line)
			if url != "" {
				urlPath = url
				urlFound = true
			}
		}

		// If we've closed all braces, we've reached the end of the node
		if braceDepth == 0 && i > startLine {
			if !urlFound {
				return "", i, fmt.Errorf("url field not found in Inline node")
			}
			return urlPath, i, nil
		}
	}

	return "", len(lines) - 1, fmt.Errorf("unclosed Inline node")
}

// extractUrlPath extracts the URL path from a line containing url field
// Example: 'url "part.wrl"' -> "part.wrl"
// TODO: 区切り文字
func (p *Parser) extractUrlPath(line string) string {
	// Pattern: url "path"
	// TODO: compile
	re := regexp.MustCompile(`url\s+"([^"]+)"`)
	matches := re.FindStringSubmatch(line)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}
