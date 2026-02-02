package token

import (
	"fmt"
	"io"

	"github.com/akisatoon1/VRML-Inline-Expander/internal/expander/tokenexpander/expander"
	"github.com/alecthomas/participle/v2/lexer"
)

var vrmlLexer = lexer.MustStateful(lexer.Rules{
	"Root": {
		{Name: "Comment", Pattern: `#[^\r\n]*`, Action: nil},
		{Name: "Whitespace", Pattern: `[\r\n \t,]+`, Action: nil},
		{Name: "String", Pattern: `"(\\["\\]|[^"\\])*"`, Action: nil},
		{Name: "DoubleWithExp", Pattern: `[+\-]?((([0-9]*\.[0-9]+)|([0-9]+(\.)?))([eE][+\-]?[0-9]+))`, Action: nil},
		{Name: "DoubleWithNoExp", Pattern: `[+\-]?(([0-9]*\.[0-9]+)|([0-9]+\.))`, Action: nil},
		{Name: "Int", Pattern: `[+\-]?((0[xX][0-9a-fA-F]+)|([0-9]+))`, Action: nil},
		{Name: "Punct", Pattern: `[.{}\[\]]`, Action: nil},
		{Name: "Ident", Pattern: `[^\x00-\x20\x22\x23\x27\x2b\x2c\x2d\x2e\x30-\x39\x5b\x5c\x5d\x7b\x7d\x7f][^\x00-\x20\x22\x23\x27\x2c\x2e\x5b\x5c\x5d\x7b\x7d\x7f]*`, Action: nil},
	},
})

// Implement Tokenizer interface in expander package
type Tokenizer struct{}

// Creates a new Tokenizer instance
func NewTokenizer() *Tokenizer {
	return &Tokenizer{}
}

// Implement Tokenizer interface signature
func (t *Tokenizer) Tokenize(r io.Reader) ([]expander.Token, error) {
	// 入力からLexerのインスタンスを作成
	lex, err := vrmlLexer.Lex("", r)
	if err != nil {
		return nil, err
	}

	// トークンを一括取得
	lexerTokens, err := lexer.ConsumeAll(lex)
	if err != nil {
		return nil, err
	}

	// participle.lexer.Tokenをexpander.Tokenに変換
	externalTokens := make([]expander.Token, 0, len(lexerTokens))
	for _, lexerToken := range lexerTokens {
		// EOFは除外
		if lexerToken.Type == lexer.EOF {
			continue
		}

		tokenT, err := convertLexerTokenType(lexerToken.Type)
		if err != nil {
			return nil, err
		}

		externalToken := newToken(tokenT, lexerToken.Value)

		if !externalToken.IsWhitespace() {
			externalTokens = append(externalTokens, externalToken)
		}
	}

	return externalTokens, nil
}

// Convert lexer.TokenType to int type assigned to token._type
func convertLexerTokenType(lexerType lexer.TokenType) (tokenType, error) {
	switch lexerType {
	case vrmlLexer.Symbols()["Punct"]:
		return punct_t, nil
	case vrmlLexer.Symbols()["Ident"]:
		return ident_t, nil
	case vrmlLexer.Symbols()["Int"], vrmlLexer.Symbols()["DoubleWithExp"], vrmlLexer.Symbols()["DoubleWithNoExp"]:
		return number_t, nil
	case vrmlLexer.Symbols()["String"]:
		return string_t, nil
	case vrmlLexer.Symbols()["Whitespace"], vrmlLexer.Symbols()["Comment"]:
		return whitespace_t, nil
	default:
		return -1, fmt.Errorf("invalid lexer token type: %d", lexerType)
	}
}
