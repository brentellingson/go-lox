package repl

import (
	"github.com/brentellingson/go-lox/internal/ast"
	"github.com/brentellingson/go-lox/internal/token"
)

type Interpreter interface {
	Interpret(statements ast.Expr) (any, error)
}

type Repl struct {
	Scan        func(source string) ([]token.Token, error)
	Parse       func(tokens []token.Token) (ast.Expr, error)
	Interpreter Interpreter
}

func NewRepl(scan func(source string) ([]token.Token, error), parse func(tokens []token.Token) (ast.Expr, error), interpreter Interpreter) *Repl {
	return &Repl{
		Scan:        scan,
		Parse:       parse,
		Interpreter: interpreter,
	}
}

func (r *Repl) Run(source string) (any, error) {
	tokens, err := r.Scan(source)
	if err != nil {
		return nil, err
	}

	expr, err := r.Parse(tokens)
	if err != nil {
		return nil, err
	}

	return r.Interpreter.Interpret(expr)
}
