package internal

type ParseError struct {
	Token   Token
	Message string
}

func (e *ParseError) Error() string {
	return "Parse Error " + e.Token.String() + ": " + e.Message
}

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	if len(tokens) == 0 || tokens[len(tokens)-1].Type != EOF {
		tokens = append(tokens, Token{Type: EOF})
	}
	return &Parser{tokens: tokens}
}

// Current returns the current token in the stream. All parser functions will assume that Current is the token that needs to be parsed.
func (t *Parser) Current() Token {
	return t.tokens[t.current]
}

// Advance moves to the next token in the stream and returns the previous current token.
func (t *Parser) Advance() Token {
	rval := t.tokens[t.current]
	if !t.IsAtEnd() {
		t.current++
	}
	return rval
}

// Peek returns the next token in the stream without advancing the current token.  If the current token is the last token in the stream, Peek will return an EOF token.
func (t *Parser) Peek() Token {
	if t.IsAtEnd() {
		return t.tokens[t.current]
	}

	return t.tokens[t.current+1]
}

// IsAtEnd returns true if the current token is the last token in the stream.
func (t *Parser) IsAtEnd() bool {
	return t.tokens[t.current].Type == EOF || t.current >= len(t.tokens)-1
}

// Check returns true if the current token matches any of the specified types, without advancing the current token.
func (t *Parser) Check(types ...TokenType) bool {
	for _, ttype := range types {
		if t.tokens[t.current].Type == ttype {
			return true
		}
	}
	return false
}

// Match returns true if the current token matches any of the specified types, and advances the current token.
func (t *Parser) Match(types ...TokenType) bool {
	if t.Check(types...) {
		t.Advance()
		return true
	}
	return false
}

func (p *Parser) Parse() Expr {
	expr, err := p.expression()
	if err != nil {
		Error(p.Current().Line, err.Error())
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

	for p.Check(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.Advance()
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

	for p.Check(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.Advance()
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

	for p.Check(PLUS, MINUS) {
		operator := p.Advance()
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

	for p.Check(STAR, SLASH) {
		operator := p.Advance()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.Check(BANG, MINUS) {
		operator := p.Advance()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &Unary{Operator: operator, Right: right}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.Match(FALSE) {
		return &Literal{Value: false}, nil
	}

	if p.Match(TRUE) {
		return &Literal{Value: true}, nil
	}

	if p.Match(NIL) {
		return &Literal{Value: nil}, nil
	}

	if p.Check(NUMBER, STRING) {
		return &Literal{Value: p.Advance().Literal}, nil
	}

	if p.Match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if !p.Match(RIGHT_PAREN) {
			return nil, &ParseError{p.Current(), "Expect ')' after expression."}
		}
		return &Grouping{Expression: expr}, nil
	}

	return nil, &ParseError{p.Current(), "Expect expression."}
}
