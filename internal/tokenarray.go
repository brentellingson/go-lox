package internal

type TokenArray struct {
	tokens  []Token
	current int
}

// TokenStream interface. Expects that final token to be EOF. Will never advance current beyond the final EOF token.
func NewTokenArray(tokens []Token) *TokenArray {
	if len(tokens) == 0 || tokens[len(tokens)-1].Type != EOF {
		tokens = append(tokens, Token{Type: EOF})
	}
	return &TokenArray{tokens: tokens}
}

// Current returns the current token in the stream. All parser functions will assume that Current is the token that needs to be parsed.
func (t *TokenArray) Current() Token {
	return t.tokens[t.current]
}

// Advance moves to the next token in the stream and returns the previous current token.
func (t *TokenArray) Advance() Token {
	rval := t.tokens[t.current]
	if !t.IsAtEnd() {
		t.current++
	}
	return rval
}

// Peek returns the next token in the stream without advancing the current token.  If the current token is the last token in the stream, Peek will return an EOF token.
func (t *TokenArray) Peek() Token {
	if t.IsAtEnd() {
		return t.tokens[t.current]
	}

	return t.tokens[t.current+1]
}

// IsAtEnd returns true if the current token is the last token in the stream.
func (t *TokenArray) IsAtEnd() bool {
	return t.tokens[t.current].Type == EOF || t.current >= len(t.tokens)-1
}

// Check returns true if the current token matches any of the specified types, without advancing the current token.
func (t *TokenArray) Check(types ...TokenType) bool {
	for _, ttype := range types {
		if t.tokens[t.current].Type == ttype {
			return true
		}
	}
	return false
}

// Match returns true if the current token matches any of the specified types, and advances the current token.
func (t *TokenArray) Match(types ...TokenType) bool {
	if t.Check(types...) {
		t.Advance()
		return true
	}
	return false
}
