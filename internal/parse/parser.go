package parse

import (
	"errors"

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

func Parse(tokens []token.Token) ([]ast.Stmt, error) {
	parser := NewParser(tokens)
	return parser.Parse()
}

type Parser struct {
	buff *TokenBuffer
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{buff: NewTokenBuffer(tokens)}
}

func (p *Parser) Parse() ([]ast.Stmt, error) {
	var errs []error
	var stmts []ast.Stmt
	for !p.buff.IsAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			errs = append(errs, err)
			p.synchronize()
		} else {
			stmts = append(stmts, stmt)
		}
	}
	return stmts, errors.Join(errs...)
}

func (p *Parser) synchronize() {
	for !p.buff.IsAtEnd() {
		prev := p.buff.Advance()
		if prev.Type == token.SEMICOLON {
			return
		}
		if p.buff.Check(
			token.CLASS,
			token.FOR,
			token.FUN,
			token.IF,
			token.PRINT,
			token.RETURN,
			token.VAR,
			token.WHILE,
		) {
			return
		}
	}
}

func (p *Parser) declaration() (ast.Stmt, error) {
	if p.buff.Match(token.VAR) {
		return p.varStatement()
	}
	if p.buff.Match(token.LEFT_BRACE) {
		return p.blockStatement()
	}

	return p.statement()
}

func (p *Parser) varStatement() (ast.Stmt, error) {
	if !p.buff.Check(token.IDENTIFIER) {
		return nil, &ParseError{p.buff.Current(), "Expect variable name."}
	}
	name := p.buff.Advance()

	var initializer ast.Expr
	if p.buff.Match(token.EQUAL) {
		var err error
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if !p.buff.Match(token.SEMICOLON) && !p.buff.IsAtEnd() {
		return nil, &ParseError{p.buff.Current(), "Expect ';' after value."}
	}

	return &ast.Var{Name: name, Expr: initializer}, nil
}

func (p *Parser) blockStatement() (ast.Stmt, error) {
	var stmts []ast.Stmt
	for !p.buff.IsAtEnd() && !p.buff.Check(token.RIGHT_BRACE) {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	if !p.buff.Match(token.RIGHT_BRACE) {
		return nil, &ParseError{p.buff.Current(), "Expect '}' after block."}
	}
	return &ast.Block{Statements: stmts}, nil
}

func (p *Parser) statement() (ast.Stmt, error) {
	if p.buff.Match(token.PRINT) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if !p.buff.Match(token.SEMICOLON) && !p.buff.IsAtEnd() {
		return nil, &ParseError{p.buff.Current(), "Expect ';' after value."}
	}
	return &ast.Print{Expression: expr}, nil
}

func (p *Parser) expressionStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if !p.buff.Match(token.SEMICOLON) && !p.buff.IsAtEnd() {
		return nil, &ParseError{p.buff.Current(), "Expect ';' after expression."}
	}
	return &ast.Expression{Expression: expr}, nil
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}
	if p.buff.Check(token.EQUAL) {
		equals := p.buff.Advance()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}
		if variable, ok := expr.(*ast.Variable); ok {
			return &ast.Assign{Name: variable.Name, Value: value}, nil
		}
		return nil, &ParseError{equals, "Invalid assignment target."}
	}
	return expr, nil
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

	if p.buff.Check(token.IDENTIFIER) {
		return &ast.Variable{Name: p.buff.Advance()}, nil
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
