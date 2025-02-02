package internal

import "fmt"

type RuntimeError struct {
	token   Token
	message string
}

func NewRuntimeError(token Token, message string) *RuntimeError {
	return &RuntimeError{token: token, message: message}
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("runtime error: %v at line %v", e.message, e.token.Line)
}

type Interpreter struct{}

func (i *Interpreter) Interpret(expr Expr) {
	result, err := expr.Accept(&Interpreter{})
	if err != nil {
		ReportRuntimeError(err)
		return
	}
	fmt.Println(result)
}

func (i *Interpreter) Evaluate(expr Expr) (any, error) {
	return expr.Accept(i)
}

func checkNumberOperands(left, right any) (float64, float64, bool) {
	if left, ok := left.(float64); ok {
		if right, ok := right.(float64); ok {
			return left, right, true
		}
	}
	return 0, 0, false
}

func (i *Interpreter) VisitBinaryExpr(expr *Binary) (any, error) {
	left, err := i.Evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.Evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case EQUAL_EQUAL:
		return left == right, nil
	case BANG_EQUAL:
		return left != right, nil
	case PLUS:
		if left, right, ok := checkNumberOperands(left, right); ok {
			return left + right, nil
		}
		if left, ok := left.(string); ok {
			return left + fmt.Sprintf("%v", right), nil
		}
	case MINUS:
		if left, right, ok := checkNumberOperands(left, right); ok {
			return left - right, nil
		}
	case STAR:
		if left, right, ok := checkNumberOperands(left, right); ok {
			return left * right, nil
		}
	case SLASH:
		if left, right, ok := checkNumberOperands(left, right); ok {
			return left / right, nil
		}
	case GREATER:
		if left, right, ok := checkNumberOperands(left, right); ok {
			return left > right, nil
		}
	case GREATER_EQUAL:
		if left, right, ok := checkNumberOperands(left, right); ok {
			return left >= right, nil
		}
	case LESS:
		if left, right, ok := checkNumberOperands(left, right); ok {
			return left < right, nil
		}
	case LESS_EQUAL:
		if left, right, ok := checkNumberOperands(left, right); ok {
			return left <= right, nil
		}
	}

	return nil, NewRuntimeError(expr.Operator, fmt.Sprintf("binary operator %v not supported for types %T, %T", expr.Operator.Type, left, right))
}

func (i *Interpreter) VisitGroupingExpr(expr *Grouping) (any, error) {
	return i.Evaluate(expr.Expression)
}

func (i *Interpreter) VisitLiteralExpr(expr *Literal) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitUnaryExpr(expr *Unary) (any, error) {
	right, err := i.Evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.Type {
	case MINUS:
		if right, ok := right.(float64); ok {
			return -right, nil
		}
	case BANG:
		return !isTruthy(right), nil
	}
	return nil, NewRuntimeError(expr.Operator, fmt.Sprintf("unary operator %v not supported for type %T", expr.Operator.Type, right))
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
