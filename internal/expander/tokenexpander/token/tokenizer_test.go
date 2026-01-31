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
			name:     "comment",
			input:    "#comment\n",
			expected: []string{},
		},
		{
			name:     "hashtag in string",
			input:    "\"#\"",
			expected: []string{"#"},
		},
		{
			name:     "whitespace",
			input:    "\r\n \t,",
			expected: []string{},
		},
		{
			name:     "whitespace in string",
			input:    "\"\r\n \t,\"",
			expected: []string{"\r\n \t,"},
		},
		{
			name:     "split by space",
			input:    "1 1",
			expected: []string{"1", "1"},
		},
		{
			name:     "split by spaces",
			input:    "1          1",
			expected: []string{"1", "1"},
		},
		{
			name:     "split by any whitespace and comments",
			input:    "1\r\n \t,#comment\n1",
			expected: []string{"1", "1"},
		},
		{
			name:     "string",
			input:    "\"string\"",
			expected: []string{"string"},
		},
		{
			name:     "double quote in string",
			input:    "\"\\\"\"",
			expected: []string{"\""},
		},
		{
			name:     "back slash in string",
			input:    "\"\\\\\"",
			expected: []string{"\\"},
		},
		{
			name:     "int32 base-10",
			input:    "123",
			expected: []string{"123"},
		},
		{
			name:     "int32 base-10 plus",
			input:    "+123",
			expected: []string{"+123"},
		},
		{
			name:     "int32 base-10 minus",
			input:    "-123",
			expected: []string{"-123"},
		},
		{
			name:     "int32 base-16 small x",
			input:    "0x7B",
			expected: []string{"0x7B"},
		},
		{
			name:     "int32 base-16 large X",
			input:    "0X7B",
			expected: []string{"0X7B"},
		},
		{
			name:     "int32 base-16 small x plus",
			input:    "+0x7B",
			expected: []string{"+0x7B"},
		},
		{
			name:     "int32 base-16 large X plus",
			input:    "+0X7B",
			expected: []string{"+0X7B"},
		},
		{
			name:     "int32 base-16 small x minus",
			input:    "-0x7B",
			expected: []string{"-0x7B"},
		},
		{
			name:     "int32 base-16 large X minus",
			input:    "-0X7B",
			expected: []string{"-0X7B"},
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
