package token

type tokenType int

const (
	punct_t tokenType = iota
	ident_t
	number_t
	string_t
	whitespace_t
)

type token struct {
	_type tokenType
	value string
}

func newToken(tokenType tokenType, value string) token {
	return token{
		_type: tokenType,
		value: value,
	}
}

func (t token) IsPunct() bool {
	return t._type == punct_t
}

func (t token) IsIdent() bool {
	return t._type == ident_t
}

func (t token) IsNumber() bool {
	return t._type == number_t
}

func (t token) IsString() bool {
	return t._type == string_t
}

func (t token) IsWhitespace() bool {
	return t._type == whitespace_t
}

func (t token) GetValue() string {
	return t.value
}
