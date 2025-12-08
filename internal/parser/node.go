package parser

// InlineNode represents an Inline node found in a VRML file
type InlineNode struct {
	DefName   string // DEF name (optional, empty string if not defined)
	UrlPath   string // URL path to the referenced file
	StartLine int    // Starting line number (0-indexed)
	EndLine   int    // Ending line number (0-indexed, inclusive)
}
