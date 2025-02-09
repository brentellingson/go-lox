package ast

import "github.com/brentellingson/go-lox/internal/token"

type Expr interface {
	Accept(v ExprVisitor) (any, error)
}

type ExprVisitor interface {
	VisitBinaryExpr(expr *Binary) (any, error)
	VisitGroupingExpr(expr *Grouping) (any, error)
	VisitLiteralExpr(expr *Literal) (any, error)
	VisitUnaryExpr(expr *Unary) (any, error)
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (e *Binary) Accept(v ExprVisitor) (any, error) {
	return v.VisitBinaryExpr(e)
}

type Grouping struct {
	Expression Expr
}

func (e *Grouping) Accept(v ExprVisitor) (any, error) {
	return v.VisitGroupingExpr(e)
}

type Literal struct {
	Value any
}

func (e *Literal) Accept(v ExprVisitor) (any, error) {
	return v.VisitLiteralExpr(e)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (e *Unary) Accept(v ExprVisitor) (any, error) {
	return v.VisitUnaryExpr(e)
}
