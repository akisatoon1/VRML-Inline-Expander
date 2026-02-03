package consumer

import (
	"testing"

	"github.com/akisatoon1/VRML-Inline-Expander/internal/expander/tokenexpander/expander"
)

const (
	punctT = iota
	identT
	numberT
	stringT
)

// testToken implements expander.Token interface for testing
type testToken struct {
	_type int
	value string
}

func (t *testToken) IsPunct() bool {
	return t._type == punctT
}

func (t *testToken) IsIdent() bool {
	return t._type == identT
}

func (t *testToken) IsNumber() bool {
	return t._type == numberT
}

func (t *testToken) IsString() bool {
	return t._type == stringT
}

func (t *testToken) GetValue() string {
	return t.value
}

func newToken(_type int, value string) expander.Token {
	return &testToken{_type: _type, value: value}
}

func TestConsumeInlineNode(t *testing.T) {
	tests := []struct {
		name           string
		tokens         []expander.Token
		startIndex     int
		expectFound    bool
		expectUrlPath  string
		expectConsumed int
		expectIsError  bool
	}{
		{
			name: "inline node with url specified",
			tokens: []expander.Token{
				newToken(identT, "Inline"),
				newToken(punctT, "{"),
				newToken(identT, "url"),
				newToken(stringT, `"file.wrl"`),
				newToken(punctT, "}"),
			},
			startIndex:     0,
			expectFound:    true,
			expectUrlPath:  "file.wrl",
			expectConsumed: 5,
			expectIsError:  false,
		},
		{
			name: "inline node at start of token sequence",
			tokens: []expander.Token{
				newToken(identT, "Inline"),
				newToken(punctT, "{"),
				newToken(identT, "url"),
				newToken(stringT, `"first.wrl"`),
				newToken(punctT, "}"),
				newToken(identT, "Shape"),
				newToken(punctT, "{"),
				newToken(punctT, "}"),
			},
			startIndex:     0,
			expectFound:    true,
			expectUrlPath:  "first.wrl",
			expectConsumed: 5,
			expectIsError:  false,
		},
		{
			name: "inline node in middle of token sequence",
			tokens: []expander.Token{
				newToken(identT, "DEF"),
				newToken(identT, "MyNode"),
				newToken(identT, "Inline"),
				newToken(punctT, "{"),
				newToken(identT, "url"),
				newToken(stringT, `"middle.wrl"`),
				newToken(punctT, "}"),
			},
			startIndex:     2,
			expectFound:    true,
			expectUrlPath:  "middle.wrl",
			expectConsumed: 5,
			expectIsError:  false,
		},
		{
			name: "inline node at end of token sequence",
			tokens: []expander.Token{
				newToken(identT, "Shape"),
				newToken(punctT, "{"),
				newToken(punctT, "}"),
				newToken(identT, "Inline"),
				newToken(punctT, "{"),
				newToken(identT, "url"),
				newToken(stringT, `"last.wrl"`),
				newToken(punctT, "}"),
			},
			startIndex:     3,
			expectFound:    true,
			expectUrlPath:  "last.wrl",
			expectConsumed: 5,
			expectIsError:  false,
		},
		{
			name: "inline node without url",
			tokens: []expander.Token{
				newToken(identT, "Inline"),
				newToken(punctT, "{"),
				newToken(punctT, "}"),
			},
			startIndex:     0,
			expectFound:    false,
			expectUrlPath:  "",
			expectConsumed: 0,
			expectIsError:  true,
		},
		{
			name: "inline node without opening brace",
			tokens: []expander.Token{
				newToken(identT, "Inline"),
				newToken(identT, "url"),
				newToken(stringT, `"file.wrl"`),
				newToken(punctT, "}"),
			},
			startIndex:     0,
			expectFound:    false,
			expectUrlPath:  "",
			expectConsumed: 0,
			expectIsError:  true,
		},
		{
			name: "inline node without closing brace",
			tokens: []expander.Token{
				newToken(identT, "Inline"),
				newToken(punctT, "{"),
				newToken(identT, "url"),
				newToken(stringT, `"file.wrl"`),
				newToken(identT, "Shape"),
			},
			startIndex:     0,
			expectFound:    false,
			expectUrlPath:  "",
			expectConsumed: 0,
			expectIsError:  true,
		},
		{
			name: "inline node with non-url field",
			tokens: []expander.Token{
				newToken(identT, "Inline"),
				newToken(punctT, "{"),
				newToken(identT, "fake"),
				newToken(numberT, "0"),
				newToken(numberT, "0"),
				newToken(numberT, "0"),
				newToken(identT, "url"),
				newToken(stringT, `"model.wrl"`),
				newToken(punctT, "}"),
			},
			startIndex:     0,
			expectFound:    false,
			expectUrlPath:  "",
			expectConsumed: 0,
			expectIsError:  true,
		},
		{
			name: "shape node",
			tokens: []expander.Token{
				newToken(identT, "Shape"),
				newToken(punctT, "{"),
				newToken(identT, "geometry"),
				newToken(identT, "Box"),
				newToken(punctT, "{"),
				newToken(punctT, "}"),
				newToken(punctT, "}"),
			},
			startIndex:     0,
			expectFound:    false,
			expectUrlPath:  "",
			expectConsumed: 0,
			expectIsError:  false,
		},
		{
			name: "startIndex out of range",
			tokens: []expander.Token{
				newToken(identT, "Inline"),
				newToken(punctT, "{"),
				newToken(identT, "url"),
				newToken(stringT, `"file.wrl"`),
				newToken(punctT, "}"),
			},
			startIndex:     10,
			expectFound:    false,
			expectUrlPath:  "",
			expectConsumed: 0,
			expectIsError:  true,
		},
		{
			name: "startIndex does not point to inline",
			tokens: []expander.Token{
				newToken(identT, "Inline"),
				newToken(punctT, "{"),
				newToken(identT, "url"),
				newToken(stringT, `"file.wrl"`),
				newToken(punctT, "}"),
			},
			startIndex:     1,
			expectFound:    false,
			expectUrlPath:  "",
			expectConsumed: 0,
			expectIsError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumer := NewInlineConsumer()
			found, inlineNode, consumed, err := consumer.ConsumeInlineNode(tt.tokens, tt.startIndex)

			if found != tt.expectFound {
				t.Errorf("expected found=%v, got %v", tt.expectFound, found)
			}

			if (err == nil) == tt.expectIsError {
				t.Errorf("expected isError=%v, got %v", tt.expectIsError, err)
			}

			if found && inlineNode != nil {
				if inlineNode.UrlPath() != tt.expectUrlPath {
					t.Errorf("expected urlPath=%q, got %q", tt.expectUrlPath, inlineNode.UrlPath())
				}
			}

			if consumed != tt.expectConsumed {
				t.Errorf("expected consumed=%d, got %d", tt.expectConsumed, consumed)
			}
		})
	}
}
