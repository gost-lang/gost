package ast

import "ghostlang.org/x/ghost/token"

type Assign struct {
	ExpressionNode
	Token token.Token
	Value ExpressionNode
}
