package interpreter

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/value"
)

func evaluateFunction(node *ast.Function, env *object.Environment) (object.Object, bool) {
	function := &object.UserFunction{Env: env, Body: node.Body, Parameters: node.Parameters}

	if node.Name.Lexeme != "" {
		env.Set(node.Name.Lexeme, function)

		return value.NULL, true
	}

	return function, true
}