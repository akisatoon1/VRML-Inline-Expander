package consumer

import (
	"fmt"

	"github.com/akisatoon1/VRML-Inline-Expander/internal/expander/tokenexpander/expander"
)

var (
	errMissingOpenBrace      = fmt.Errorf("expected '{' after Inline keyword")
	errMissingUrlField       = fmt.Errorf("Inline node must have a url field")
	errMissingUrlValue       = fmt.Errorf("expected url value after url keyword")
	errMissingCloseBrace     = fmt.Errorf("expected '}' to close Inline node")
	errUnexpectedEndOfTokens = fmt.Errorf("unexpected end of tokens while parsing Inline node")
	errIndexOutOfRange       = fmt.Errorf("index is out of range of tokens")
)

// token values want to check
const (
	keywordInline     = "Inline"
	keywordUrl        = "url"
	keywordOpenBrace  = "{"
	keywordCloseBrace = "}"
)

// implements expander.InlineNode interface
type inlineNode struct {
	urlPath string
}

func (n *inlineNode) UrlPath() string {
	return n.urlPath
}

// implements expander.Consumer interface
type InlineConsumer struct{}

func NewInlineConsumer() *InlineConsumer {
	return &InlineConsumer{}
}

// Consume tokens to find Inline node and url value in the node.
func (c *InlineConsumer) ConsumeInlineNode(tokens []expander.Token, startIndex int) (bool, expander.InlineNode, int, error) {
	if !isValidPosition(tokens, startIndex) {
		return false, nil, 0, errIndexOutOfRange
	}

	currentPos := startIndex

	// check "Inline" token
	ok, err := isExpectedToken(tokens, currentPos, keywordInline)
	if err != nil {
		return false, nil, 0, err
	}
	if !ok {
		// if not "Inline", end consuming.
		return false, nil, 0, nil
	}
	currentPos++

	// check "{" token
	if ok, err = isExpectedToken(tokens, currentPos, keywordOpenBrace); err != nil {
		return false, nil, 0, err
	}
	if !ok {
		return false, nil, 0, errMissingOpenBrace
	}
	currentPos++

	// check "url" token
	if ok, err = isExpectedToken(tokens, currentPos, keywordUrl); err != nil {
		return false, nil, 0, err
	}
	if !ok {
		return false, nil, 0, errMissingUrlField
	}
	currentPos++

	// check value of url field
	url, err := expectString(tokens, currentPos)
	if err != nil {
		return false, nil, 0, err
	}
	currentPos++

	// check "}" token
	if ok, err = isExpectedToken(tokens, currentPos, keywordCloseBrace); err != nil {
		return false, nil, 0, err
	}
	if !ok {
		return false, nil, 0, errMissingCloseBrace
	}

	node := &inlineNode{urlPath: url}
	consumed := currentPos - startIndex + 1
	return true, node, consumed, nil
}

// prevent indexOutOfRange
func isValidPosition(tokens []expander.Token, position int) bool {
	return 0 <= position && position < len(tokens)
}

// check if the token value is expected value, return bool about match or not.
func isExpectedToken(tokens []expander.Token, position int, expectedValue string) (bool, error) {
	if !isValidPosition(tokens, position) {
		return false, errUnexpectedEndOfTokens
	}
	if tokens[position].GetValue() != expectedValue {
		return false, nil
	}
	return true, nil
}

// expect tokens[position] is string and return its string value.
func expectString(tokens []expander.Token, position int) (string, error) {
	if !isValidPosition(tokens, position) {
		return "", errUnexpectedEndOfTokens
	}
	if !tokens[position].IsString() {
		return "", errMissingUrlValue
	}

	rawStr := tokens[position].GetValue()
	strValue := rawStr[1 : len(rawStr)-1] // remove surrounding quotes
	return strValue, nil
}
