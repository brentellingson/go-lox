package main

import (
	"fmt"

	"github.com/brentellingson/go-lox/internal"
)

func main() {
	expr := internal.Binary{
		Left: &internal.Unary{
			Operator: internal.Token{Type: internal.MINUS, Lexeme: "-", Literal: nil, Line: 1},
			Right:    &internal.Literal{Value: 123},
		},
		Operator: internal.Token{Type: internal.STAR, Lexeme: "*", Literal: nil, Line: 1},
		Right:    &internal.Grouping{Expression: &internal.Literal{Value: 45.67}},
	}

	fmt.Println(internal.PrintAst(&expr))
}
