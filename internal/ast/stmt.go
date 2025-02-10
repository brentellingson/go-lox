package ast

import "github.com/brentellingson/go-lox/internal/token"

type Stmt interface {
	Accept(v StmtVisitor) (any, error)
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt *Expression) (any, error)
	VisitPrintStmt(stmt *Print) (any, error)
	VisitVarStmt(stmt *Var) (any, error)
	VisitBlockStmt(stmt *Block) (any, error)
	VisitIfStmt(stmt *If) (any, error)
	VisitWhileStmt(stmt *While) (any, error)
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

type Var struct {
	Name       token.Token
	Expression Expr
}

func (e *Var) Accept(v StmtVisitor) (any, error) {
	return v.VisitVarStmt(e)
}

type Block struct {
	Statements []Stmt
}

func (e *Block) Accept(v StmtVisitor) (any, error) {
	return v.VisitBlockStmt(e)
}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (e *If) Accept(v StmtVisitor) (any, error) {
	return v.VisitIfStmt(e)
}

type While struct {
	Condition Expr
	Body      Stmt
}

func (e *While) Accept(v StmtVisitor) (any, error) {
	return v.VisitWhileStmt(e)
}
