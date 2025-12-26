package expander

import "github.com/akisatoon1/VRML-Inline-Expander/internal/parser"

type inlineNodesFinder interface {
	FindInlineNodes(lines []string) ([]parser.InlineNode, error)
}
