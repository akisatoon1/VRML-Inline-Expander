/*
	Test Tokenizer.Tokenize method
*/

package token

import (
	"strings"
	"testing"

	"github.com/akisatoon1/VRML-Inline-Expander/internal/expander/tokenexpander/expander"
)

type testcase struct {
	name     string
	input    string
	expected []string
}

func TestTokenizer_Tokenize(t *testing.T) {
	tests := []testcase{
		{
			name:  "dummy case",
			input: "DEF MyShape Shape { }",
			expected: []string{
				"DEF", "MyShape", "Shape", "{", "}",
			},
		},
	}

	testAllCases(t, tests)
}

func testAllCases(t *testing.T, tests []testcase) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewTokenizer()
			reader := strings.NewReader(tt.input)

			got, err := tokenizer.Tokenize(reader)
			if err != nil {
				t.Fatalf("Tokenize() error = %v", err)
			}

			assertTokensEqual(t, got, tt.expected)
		})
	}
}

// assertTokensEqual asserts that tokens match expected values
func assertTokensEqual(t *testing.T, got []expander.Token, expected []string) {
	t.Helper() // マークすることでスタックトレースが呼び出し元を指す

	gotStrings := tokensToStrings(got)

	if len(gotStrings) != len(expected) {
		t.Errorf("token count = %d, want %d\ngot:  %v\nwant: %v",
			len(gotStrings), len(expected), gotStrings, expected)
		return
	}

	// TODO: タイプを比較していないので、あとでトークンごと直接比較
	for i := range gotStrings {
		if gotStrings[i] != expected[i] {
			t.Errorf("Token[%d] = %q, want %q", i, gotStrings[i], expected[i])
		}
	}
}

// Helper function to convert tokens to strings for debugging
func tokensToStrings(tokens []expander.Token) []string {
	result := make([]string, len(tokens))
	for i, token := range tokens {
		result[i] = token.GetValue()
	}
	return result
}
