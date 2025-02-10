package internal

import (
	"fmt"
	"strings"

	"github.com/brentellingson/go-lox/internal/ast"
)

func PrintAst(expr ast.Expr) string {
	result, err := expr.Accept(&AstPrinter{})
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return result.(string)
}

type AstPrinter struct{}

func (p *AstPrinter) VisitBinaryExpr(expr *ast.Binary) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) VisitGroupingExpr(expr *ast.Grouping) (any, error) {
	return p.parenthesize("group", expr.Expression)
}

func (p *AstPrinter) VisitLiteralExpr(expr *ast.Literal) (any, error) {
	return fmt.Sprintf("%#v", expr.Value), nil
}

func (p *AstPrinter) VisitUnaryExpr(expr *ast.Unary) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p *AstPrinter) VisitVariableExpr(expr *ast.Variable) (any, error) {
	return expr.Name.Lexeme, nil
}

func (p *AstPrinter) VisitAssignExpr(expr *ast.Assign) (any, error) {
	return p.parenthesize("set! "+expr.Name.Lexeme, expr.Value)
}

func (p *AstPrinter) VisitLogicalExpr(expr *ast.Logical) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) parenthesize(name string, exprs ...ast.Expr) (any, error) {
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
