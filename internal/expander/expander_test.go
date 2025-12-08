package expander

import (
	"testing"

	"github.com/akisatoon1/VRML-Inline-Expander/internal/parser"
)

func TestExpandLines(t *testing.T) {
	e := New()

	lines := []string{
		"#VRML V2.0 utf8",
		"Inline {",
		"  url \"part.wrl\"",
		"}",
		"Shape { geometry Sphere {} }",
	}

	nodesWithContent := []nodeWithContent{
		{
			Node: parser.InlineNode{
				DefName:   "",
				UrlPath:   "part.wrl",
				StartLine: 1,
				EndLine:   3,
			},
			RefLines: []string{
				"Shape { geometry Box {} }",
			},
		},
	}

	result := e.expandLines(lines, nodesWithContent)

	expected := []string{
		"#VRML V2.0 utf8",
		"Group {",
		"  children [",
		"    Shape { geometry Box {} }",
		"  ]",
		"}",
		"Shape { geometry Sphere {} }",
	}

	if len(result) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(result))
	}

	for i, line := range expected {
		if result[i] != line {
			t.Errorf("line %d: expected %q, got %q", i, line, result[i])
		}
	}
}
