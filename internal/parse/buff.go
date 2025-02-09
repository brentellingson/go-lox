package parse

import (
	"github.com/brentellingson/go-lox/internal/token"
)

type TokenBuffer struct {
	tokens  []token.Token
	current int
}

func NewTokenBuffer(tokens []token.Token) *TokenBuffer {
	if len(tokens) == 0 || tokens[len(tokens)-1].Type != token.EOF {
		tokens = append(tokens, token.Token{Type: token.EOF})
	}
	return &TokenBuffer{tokens: tokens}
}

// Current returns the current token in the stream. All parser functions will assume that Current is the token that needs to be parsed.
func (t *TokenBuffer) Current() token.Token {
	return t.tokens[t.current]
}

// Advance moves to the next token in the stream and returns the previous current token.
func (t *TokenBuffer) Advance() token.Token {
	rval := t.tokens[t.current]
	if !t.IsAtEnd() {
		t.current++
	}
	return rval
}

// Peek returns the next token in the stream without advancing the current token.  If the current token is the last token in the stream, Peek will return an EOF token.
func (t *TokenBuffer) Peek() token.Token {
	if t.IsAtEnd() {
		return t.tokens[t.current]
	}

	return t.tokens[t.current+1]
}

// IsAtEnd returns true if the current token is the last token in the stream.
func (t *TokenBuffer) IsAtEnd() bool {
	return t.tokens[t.current].Type == token.EOF || t.current >= len(t.tokens)-1
}

// Check returns true if the current token matches any of the specified types, without advancing the current token.
func (t *TokenBuffer) Check(types ...token.TokenType) bool {
	for _, ttype := range types {
		if t.tokens[t.current].Type == ttype {
			return true
		}
	}
	return false
}

// Match returns true if the current token matches any of the specified types, and advances the current token.
func (t *TokenBuffer) Match(types ...token.TokenType) bool {
	if t.Check(types...) {
		t.Advance()
		return true
	}
	return false
}
