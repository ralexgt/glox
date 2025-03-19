package main

import (
	"fmt"

	"github.com/ralexgt/glox/expr"
	"github.com/ralexgt/glox/token"
)

func main() {
	// Expression is a hard coded instance of a tree
	// Used to test before actually implementing the parser
	expression := expr.Binary[string]{
		Left: expr.Unary[string]{
			Operator: token.NewToken(token.TokenType_Minus, "-", nil, 1),
			Right: expr.Literal[string]{
				Value: 123,
			},
		},
		Operator: token.NewToken(token.TokenType_Star, "*", nil, 1),
		Right: expr.Grouping[string]{
			Expression: expr.Literal[string]{
				Value: 45.7,
			},
		},
	}
	printer := expr.AstPrinter{}
	fmt.Println(printer.Print(expression))

}
