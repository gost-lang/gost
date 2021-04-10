package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/environment"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateExpression(node *ast.Expression, env *environment.Environment) (object.Object, bool) {
	result, ok := Evaluate(node.Expression, env)

	if !ok {
		return result, ok
	}

	return value.NULL, ok
}