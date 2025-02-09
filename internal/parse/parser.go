package parse

import (
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

func Parse(tokens []token.Token) (ast.Expr, error) {
	parser := NewParser(tokens)
	return parser.Parse()
}

type Parser struct {
	buff *TokenBuffer
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{buff: NewTokenBuffer(tokens)}
}

func (p *Parser) Parse() (ast.Expr, error) {
	expr, err := p.expression()
	return expr, err
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.buff.Check(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.buff.Advance()
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

	for p.buff.Check(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.buff.Advance()
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

	for p.buff.Check(token.PLUS, token.MINUS) {
		operator := p.buff.Advance()
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

	for p.buff.Check(token.STAR, token.SLASH) {
		operator := p.buff.Advance()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.buff.Check(token.BANG, token.MINUS) {
		operator := p.buff.Advance()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &ast.Unary{Operator: operator, Right: right}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.buff.Match(token.FALSE) {
		return &ast.Literal{Value: false}, nil
	}

	if p.buff.Match(token.TRUE) {
		return &ast.Literal{Value: true}, nil
	}

	if p.buff.Match(token.NIL) {
		return &ast.Literal{Value: nil}, nil
	}

	if p.buff.Check(token.NUMBER, token.STRING) {
		return &ast.Literal{Value: p.buff.Advance().Literal}, nil
	}

	if p.buff.Match(token.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if !p.buff.Match(token.RIGHT_PAREN) {
			return nil, &ParseError{p.buff.Current(), "Expect ')' after expression."}
		}
		return &ast.Grouping{Expression: expr}, nil
	}

	return nil, &ParseError{p.buff.Current(), "Expect expression."}
}
