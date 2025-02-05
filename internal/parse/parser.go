package parse

import (
	"github.com/brentellingson/go-lox/internal"
	"github.com/brentellingson/go-lox/internal/ast"
	"github.com/brentellingson/go-lox/internal/token"
)

type ParseError struct {
	Token   token.Token
	Message string
}

func (e *ParseError) Error() string {
	return "Parse Error " + e.Token.String() + ": " + e.Message
}

type Parser struct {
	tokens  []token.Token
	current int
}

func NewParser(tokens []token.Token) *Parser {
	if len(tokens) == 0 || tokens[len(tokens)-1].Type != token.EOF {
		tokens = append(tokens, token.Token{Type: token.EOF})
	}
	return &Parser{tokens: tokens}
}

// Current returns the current token in the stream. All parser functions will assume that Current is the token that needs to be parsed.
func (t *Parser) Current() token.Token {
	return t.tokens[t.current]
}

// Advance moves to the next token in the stream and returns the previous current token.
func (t *Parser) Advance() token.Token {
	rval := t.tokens[t.current]
	if !t.IsAtEnd() {
		t.current++
	}
	return rval
}

// Peek returns the next token in the stream without advancing the current token.  If the current token is the last token in the stream, Peek will return an EOF token.
func (t *Parser) Peek() token.Token {
	if t.IsAtEnd() {
		return t.tokens[t.current]
	}

	return t.tokens[t.current+1]
}

// IsAtEnd returns true if the current token is the last token in the stream.
func (t *Parser) IsAtEnd() bool {
	return t.tokens[t.current].Type == token.EOF || t.current >= len(t.tokens)-1
}

// Check returns true if the current token matches any of the specified types, without advancing the current token.
func (t *Parser) Check(types ...token.TokenType) bool {
	for _, ttype := range types {
		if t.tokens[t.current].Type == ttype {
			return true
		}
	}
	return false
}

// Match returns true if the current token matches any of the specified types, and advances the current token.
func (t *Parser) Match(types ...token.TokenType) bool {
	if t.Check(types...) {
		t.Advance()
		return true
	}
	return false
}

func (p *Parser) Parse() ast.Expr {
	expr, err := p.expression()
	if err != nil {
		internal.ReportParserError(p.Current().Line, err.Error())
	}
	return expr
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.Check(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.Advance()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.Check(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.Advance()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) term() (ast.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.Check(token.PLUS, token.MINUS) {
		operator := p.Advance()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) factor() (ast.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.Check(token.STAR, token.SLASH) {
		operator := p.Advance()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.Check(token.BANG, token.MINUS) {
		operator := p.Advance()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &ast.Unary{Operator: operator, Right: right}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.Match(token.FALSE) {
		return &ast.Literal{Value: false}, nil
	}

	if p.Match(token.TRUE) {
		return &ast.Literal{Value: true}, nil
	}

	if p.Match(token.NIL) {
		return &ast.Literal{Value: nil}, nil
	}

	if p.Check(token.NUMBER, token.STRING) {
		return &ast.Literal{Value: p.Advance().Literal}, nil
	}

	if p.Match(token.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if !p.Match(token.RIGHT_PAREN) {
			return nil, &ParseError{p.Current(), "Expect ')' after expression."}
		}
		return &ast.Grouping{Expression: expr}, nil
	}

	return nil, &ParseError{p.Current(), "Expect expression."}
}
