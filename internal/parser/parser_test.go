package parser

import (
	"testing"
)

func TestFindInlineNodes(t *testing.T) {
	p := New()
	lines := []string{
		"#VRML V2.0 utf8",
		"Inline {",
		"  url \"part.wrl\"",
		"}",
	}

	nodes, err := p.FindInlineNodes(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}

	node := nodes[0]
	if node.DefName != "" {
		t.Errorf("expected DefName '', got '%s'", node.DefName)
	}
	if node.UrlPath != "part.wrl" {
		t.Errorf("expected UrlPath 'part.wrl', got '%s'", node.UrlPath)
	}
	if node.StartLine != 1 {
		t.Errorf("expected StartLine 1, got %d", node.StartLine)
	}
	if node.EndLine != 3 {
		t.Errorf("expected EndLine 3, got %d", node.EndLine)
	}
}
