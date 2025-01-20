package internal

import (
	"fmt"
	"strings"
)

func PrintAst(expr Expr) string {
	return expr.Accept(&AstPrinter{}).(string)
}

type AstPrinter struct{}

func (p *AstPrinter) VisitBinaryExpr(expr *Binary) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitGroupingExpr(expr *Grouping) any {
	return p.parenthesize("group", expr.Expression)
}

func (p *AstPrinter) VisitLiteralExpr(expr *Literal) any {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

func (p *AstPrinter) VisitUnaryExpr(expr *Unary) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var b strings.Builder
	b.WriteRune('(')
	b.WriteString(name)
	for _, e := range exprs {
		b.WriteRune(' ')
		b.WriteString(e.Accept(p).(string))
	}
	b.WriteRune(')')
	return b.String()
}
