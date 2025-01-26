package internal

type TokenStream interface {
	// Current returns the current token in the stream. All parser functions will assume that Current is the token that needs to be parsed.
	Current() Token

	// Advance moves to the next token in the stream and returns the previous current token.
	Advance() Token

	// Peek returns the next token in the stream without advancing the current token.
	Peek() Token

	// IsAtEnd returns true if the current token is the last token in the stream.
	IsAtEnd() bool

	// Check returns true if the current token matches any of the specified types, without advancing the current token.
	Check(types ...TokenType) bool

	// Match returns true if the current token matches any of the specified types, and advances the current token.
	Match(types ...TokenType) bool
}

type ParseError struct {
	Token   Token
	Message string
}

func (e *ParseError) Error() string {
	return "Parse Error " + e.Token.String() + ": " + e.Message
}

type Parser struct {
	stream TokenStream
}

func NewParser(stream TokenStream) *Parser {
	return &Parser{stream: stream}
}

func (p *Parser) Parse() Expr {
	expr, err := p.expression()
	if err != nil {
		Error(p.stream.Current().Line, err.Error())
	}
	return expr
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.stream.Check(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.stream.Advance()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.stream.Check(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.stream.Advance()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.stream.Check(PLUS, MINUS) {
		operator := p.stream.Advance()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.stream.Check(STAR, SLASH) {
		operator := p.stream.Advance()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.stream.Check(BANG, MINUS) {
		operator := p.stream.Advance()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &Unary{Operator: operator, Right: right}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.stream.Match(FALSE) {
		return &Literal{Value: false}, nil
	}

	if p.stream.Match(TRUE) {
		return &Literal{Value: true}, nil
	}

	if p.stream.Match(NIL) {
		return &Literal{Value: nil}, nil
	}

	if p.stream.Check(NUMBER, STRING) {
		return &Literal{Value: p.stream.Advance().Literal}, nil
	}

	if p.stream.Match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if !p.stream.Match(RIGHT_PAREN) {
			return nil, &ParseError{p.stream.Current(), "Expect ')' after expression."}
		}
		return &Grouping{Expression: expr}, nil
	}

	return nil, &ParseError{p.stream.Current(), "Expect expression."}
}
