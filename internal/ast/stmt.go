package ast

type Stmt interface {
	Accept(v StmtVisitor) (any, error)
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt *Expression) (any, error)
	VisitPrintStmt(stmt *Print) (any, error)
}

type Print struct {
	Expression Expr
}

func (e *Print) Accept(v StmtVisitor) (any, error) {
	return v.VisitPrintStmt(e)
}

type Expression struct {
	Expression Expr
}

func (e *Expression) Accept(v StmtVisitor) (any, error) {
	return v.VisitExpressionStmt(e)
}
