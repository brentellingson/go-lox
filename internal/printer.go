package internal

import (
	"fmt"
	"strings"
)

func PrintAst(expr Expr) string {
	result, err := expr.Accept(&AstPrinter{})
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return result.(string)
}

type AstPrinter struct{}

func (p *AstPrinter) VisitBinaryExpr(expr *Binary) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitGroupingExpr(expr *Grouping) (any, error) {
	return p.parenthesize("group", expr.Expression)
}

func (p *AstPrinter) VisitLiteralExpr(expr *Literal) (any, error) {
	return fmt.Sprintf("%#v", expr.Value), nil
}

func (p *AstPrinter) VisitUnaryExpr(expr *Unary) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr) (any, error) {
	var b strings.Builder
	b.WriteRune('(')
	b.WriteString(name)
	for _, e := range exprs {
		b.WriteRune(' ')
		v, err := e.Accept(p)
		if err != nil {
			return nil, err
		}
		b.WriteString(v.(string))
	}
	b.WriteRune(')')
	return b.String(), nil
}
