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
	expected []token
}

func TestTokenizer_Tokenize(t *testing.T) {
	tests := []testcase{
		{
			name:     "comment",
			input:    "#comment\n",
			expected: []token{},
		},
		{
			name:     "comment not end",
			input:    "#comment",
			expected: []token{},
		},
		{
			name:     "hashtag in string",
			input:    "\"#\"",
			expected: []token{{_type: string_t, value: "#"}},
		},
		{
			name:     "whitespace",
			input:    "\r\n \t,",
			expected: []token{},
		},
		{
			name:     "whitespace in string",
			input:    "\"\r\n \t,\"",
			expected: []token{{_type: string_t, value: "\r\n \t,"}},
		},
		{
			name:     "split by space",
			input:    "1 1",
			expected: []token{{_type: number_t, value: "1"}, {_type: number_t, value: "1"}},
		},
		{
			name:     "split by spaces",
			input:    "1          1",
			expected: []token{{_type: number_t, value: "1"}, {_type: number_t, value: "1"}},
		},
		{
			name:     "split by any whitespace and comments",
			input:    "1\r\n \t,#comment\n1",
			expected: []token{{_type: number_t, value: "1"}, {_type: number_t, value: "1"}},
		},
		{
			name:     "string",
			input:    "\"string\"",
			expected: []token{{_type: string_t, value: "string"}},
		},
		{
			name:     "double quote in string",
			input:    "\"\\\"\"",
			expected: []token{{_type: string_t, value: "\""}},
		},
		{
			name:     "back slash in string",
			input:    "\"\\\\\"",
			expected: []token{{_type: string_t, value: "\\"}},
		},
		{
			name:     "sequencial strings",
			input:    "\"string1\"\"string2\"",
			expected: []token{{_type: string_t, value: "string1"}, {_type: string_t, value: "string2"}},
		},
		{
			name:     "other separated by string",
			input:    "a\"string\"b",
			expected: []token{{_type: ident_t, value: "a"}, {_type: string_t, value: "string"}, {_type: ident_t, value: "b"}},
		},
		{
			name:     "int32 base-10",
			input:    "123",
			expected: []token{{_type: number_t, value: "123"}},
		},
		{
			name:     "int32 base-10 plus",
			input:    "+123",
			expected: []token{{_type: number_t, value: "+123"}},
		},
		{
			name:     "int32 base-10 minus",
			input:    "-123",
			expected: []token{{_type: number_t, value: "-123"}},
		},
		{
			name:     "int32 base-16 small x",
			input:    "0x7B",
			expected: []token{{_type: number_t, value: "0x7B"}},
		},
		{
			name:     "int32 base-16 large X",
			input:    "0X7B",
			expected: []token{{_type: number_t, value: "0X7B"}},
		},
		{
			name:     "int32 base-16 small x plus",
			input:    "+0x7B",
			expected: []token{{_type: number_t, value: "+0x7B"}},
		},
		{
			name:     "int32 base-16 large X plus",
			input:    "+0X7B",
			expected: []token{{_type: number_t, value: "+0X7B"}},
		},
		{
			name:     "int32 base-16 small x minus",
			input:    "-0x7B",
			expected: []token{{_type: number_t, value: "-0x7B"}},
		},
		{
			name:     "int32 base-16 large X minus",
			input:    "-0X7B",
			expected: []token{{_type: number_t, value: "-0X7B"}},
		},
		{
			name:     "double",
			input:    "123.456",
			expected: []token{{_type: number_t, value: "123.456"}},
		},
		{
			name:     "double plus",
			input:    "+123.456",
			expected: []token{{_type: number_t, value: "+123.456"}},
		},
		{
			name:     "double minus",
			input:    "-123.456",
			expected: []token{{_type: number_t, value: "-123.456"}},
		},
		{
			name:     "double no fractional part",
			input:    "123.",
			expected: []token{{_type: number_t, value: "123."}},
		},
		{
			name:     "double small e",
			input:    "1.2e2",
			expected: []token{{_type: number_t, value: "1.2e2"}},
		},
		{
			name:     "double large E",
			input:    "1.2E2",
			expected: []token{{_type: number_t, value: "1.2E2"}},
		},
		{
			name:     "double e plus",
			input:    "1.2e+2",
			expected: []token{{_type: number_t, value: "1.2e+2"}},
		},
		{
			name:     "double e minus",
			input:    "1.2e-2",
			expected: []token{{_type: number_t, value: "1.2e-2"}},
		},
		{
			name:     "ident only first char letter",
			input:    "a",
			expected: []token{{_type: ident_t, value: "a"}},
		},
		{
			name:     "ident only first char upper",
			input:    "A",
			expected: []token{{_type: ident_t, value: "A"}},
		},
		{
			name:     "ident only first char underbar",
			input:    "_",
			expected: []token{{_type: ident_t, value: "_"}},
		},
		{
			name:     "ident",
			input:    "ab",
			expected: []token{{_type: ident_t, value: "ab"}},
		},
		{
			name:     "ident sign",
			input:    "a+-",
			expected: []token{{_type: ident_t, value: "a+-"}},
		},
		{
			name:     "punct .",
			input:    ".",
			expected: []token{{_type: punct_t, value: "."}},
		},
		{
			name:     "punct {",
			input:    "{",
			expected: []token{{_type: punct_t, value: "{"}},
		},
		{
			name:     "punct }",
			input:    "}",
			expected: []token{{_type: punct_t, value: "}"}},
		},
		{
			name:     "punct [",
			input:    "[",
			expected: []token{{_type: punct_t, value: "["}},
		},
		{
			name:     "punct ]",
			input:    "]",
			expected: []token{{_type: punct_t, value: "]"}},
		},
		{
			name:  "sequencial puncts",
			input: ".,{}[]",
			expected: []token{
				{_type: punct_t, value: "."},
				{_type: punct_t, value: ","},
				{_type: punct_t, value: "{"},
				{_type: punct_t, value: "}"},
				{_type: punct_t, value: "["},
				{_type: punct_t, value: "]"},
			},
		},
		{
			name:  "other separated by punct",
			input: "a.b,c{d}e[f]g",
			expected: []token{
				{_type: ident_t, value: "a"},
				{_type: punct_t, value: "."},
				{_type: ident_t, value: "b"},
				{_type: punct_t, value: ","},
				{_type: ident_t, value: "c"},
				{_type: punct_t, value: "{"},
				{_type: ident_t, value: "d"},
				{_type: punct_t, value: "}"},
				{_type: ident_t, value: "e"},
				{_type: punct_t, value: "["},
				{_type: ident_t, value: "f"},
				{_type: punct_t, value: "]"},
				{_type: ident_t, value: "g"},
			},
		},
		{
			name:  "separator . and double",
			input: "a9.0.1",
			expected: []token{
				{_type: ident_t, value: "a9"},
				{_type: punct_t, value: "."},
				{_type: number_t, value: "0.1"},
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
func assertTokensEqual(t *testing.T, got []expander.Token, expected []token) {
	t.Helper() // マークすることでスタックトレースが呼び出し元を指す

	if len(got) != len(expected) {
		t.Errorf("token count = %d, want %d\ngot:  %v\nwant: %v",
			len(got), len(expected), got, expected)
		return
	}

	for i := range got {
		gotToken, expectedToken := got[i].(token), expected[i]

		if gotToken._type != expectedToken._type {
			t.Errorf("Token[%d] type = %d, want %d", i, gotToken._type, expectedToken._type)
		}

		if gotToken.GetValue() != expectedToken.value {
			t.Errorf("Token[%d] = %q, want %q", i, gotToken.GetValue(), expectedToken.value)
		}
	}
}
