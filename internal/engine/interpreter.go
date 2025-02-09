package engine

import (
	"fmt"

	"github.com/brentellingson/go-lox/internal/ast"
	"github.com/brentellingson/go-lox/internal/token"
)

type RuntimeError struct {
	token   token.Token
	message string
}

func NewRuntimeError(token token.Token, message string) *RuntimeError {
	return &RuntimeError{token: token, message: message}
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("runtime error: %v at line %v", e.message, e.token.Line)
}

type Interpreter struct {
	env *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{env: NewEnvironment()}
}

func (i *Interpreter) Interpret(stmts []ast.Stmt) (any, error) {
	var rslt any
	for _, stmt := range stmts {
		var err error
		rslt, err = i.execute(stmt)
		if err != nil {
			return nil, err
		}
	}
	return rslt, nil
}

func (i *Interpreter) execute(stmt ast.Stmt) (any, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) Evaluate(expr ast.Expr) (any, error) {
	return expr.Accept(i)
}

func (i *Interpreter) VisitExpressionStmt(stmt *ast.Expression) (any, error) {
	return i.Evaluate(stmt.Expression)
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.Print) (any, error) {
	rslt, err := i.Evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}
	fmt.Println(rslt)
	return nil, nil
}

func (i *Interpreter) VisitVarStmt(stmt *ast.Var) (any, error) {
	var rslt any
	if stmt.Expr != nil {
		var err error
		rslt, err = i.Evaluate(stmt.Expr)
		if err != nil {
			return nil, err
		}
	}
	i.env.Define(stmt.Name.Lexeme, rslt)
	return nil, nil
}

func checkNumberOperands(left, right any) (float64, float64, bool) {
	if left, ok := left.(float64); ok {
		if right, ok := right.(float64); ok {
			return left, right, true
		}
	}
	return 0, 0, false
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary) (any, error) {
	left, err := i.Evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.Evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.EQUAL_EQUAL:
		return left == right, nil
	case token.BANG_EQUAL:
		return left != right, nil
	case token.PLUS:
		if left, right, ok := checkNumberOperands(left, right); ok {
			return left + right, nil
		}
		if left, ok := left.(string); ok {
			return left + fmt.Sprintf("%v", right), nil
		}
	case token.MINUS:
		if left, right, ok := checkNumberOperands(left, right); ok {
			return left - right, nil
		}
	case token.STAR:
		if left, right, ok := checkNumberOperands(left, right); ok {
			return left * right, nil
		}
	case token.SLASH:
		if left, right, ok := checkNumberOperands(left, right); ok {
			return left / right, nil
		}
	case token.GREATER:
		if left, right, ok := checkNumberOperands(left, right); ok {
			return left > right, nil
		}
	case token.GREATER_EQUAL:
		if left, right, ok := checkNumberOperands(left, right); ok {
			return left >= right, nil
		}
	case token.LESS:
		if left, right, ok := checkNumberOperands(left, right); ok {
			return left < right, nil
		}
	case token.LESS_EQUAL:
		if left, right, ok := checkNumberOperands(left, right); ok {
			return left <= right, nil
		}
	}

	return nil, NewRuntimeError(expr.Operator, fmt.Sprintf("binary operator %v not supported for types %T, %T", expr.Operator.Type, left, right))
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.Grouping) (any, error) {
	return i.Evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.Literal) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary) (any, error) {
	right, err := i.Evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.Type {
	case token.MINUS:
		if right, ok := right.(float64); ok {
			return -right, nil
		}
	case token.BANG:
		return !isTruthy(right), nil
	}
	return nil, NewRuntimeError(expr.Operator, fmt.Sprintf("unary operator %v not supported for type %T", expr.Operator.Type, right))
}

func (i *Interpreter) VisitVariableExpr(expr *ast.Variable) (any, error) {
	v, ok := i.env.Get(expr.Name.Lexeme)
	if !ok {
		return nil, &RuntimeError{token: expr.Name, message: "undefined variable " + expr.Name.Lexeme}
	}
	return v, nil
}

func (i *Interpreter) VisitAssignExpr(expr *ast.Assign) (any, error) {
	value, err := i.Evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	if ok := i.env.Assign(expr.Name.Lexeme, value); !ok {
		return nil, &RuntimeError{token: expr.Name, message: "undefined variable " + expr.Name.Lexeme}
	}
	return value, nil
}

func isTruthy(v any) bool {
	switch v := v.(type) {
	case nil:
		return false
	case bool:
		return v
	default:
		return true
	}
}
