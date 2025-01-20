package internal

type Expr interface {
	Accept(v Visitor) any
}

type Visitor interface {
	VisitBinaryExpr(expr *Binary) any
	VisitGroupingExpr(expr *Grouping) any
	VisitLiteralExpr(expr *Literal) any
	VisitUnaryExpr(expr *Unary) any
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (e *Binary) Accept(v Visitor) any {
	return v.VisitBinaryExpr(e)
}

type Grouping struct {
	Expression Expr
}

func (e *Grouping) Accept(v Visitor) any {
	return v.VisitGroupingExpr(e)
}

type Literal struct {
	Value any
}

func (e *Literal) Accept(v Visitor) any {
	return v.VisitLiteralExpr(e)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (e *Unary) Accept(v Visitor) any {
	return v.VisitUnaryExpr(e)
}
