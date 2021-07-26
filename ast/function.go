package ast

import (
	"ghostlang.org/x/ghost/token"
)

type Function struct {
	Token      token.Token
	Name       token.Token
	Parameters []Identifier
	Defaults   map[string]ExpressionNode
	Body       []StatementNode
}
