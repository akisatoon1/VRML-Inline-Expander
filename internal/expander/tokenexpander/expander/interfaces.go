package expander

import "io"

type Tokenizer interface {
	Tokenize(r io.Reader) ([]Token, error)
}

type Token interface {
	IsPunct() bool
	IsIdent() bool
	IsNumber() bool
	IsString() bool
	GetValue() string
}

type InlineNode interface {
	UrlPath() string
}

type Consumer interface {
	ConsumeInlineNode(tokens []Token, startIndex int) (isFound bool, inline InlineNode, consumeNum int, err error)
}
