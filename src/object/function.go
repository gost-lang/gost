package object

import (
	"io"

	"ghostlang.org/x/ghost/ast"
)

const FUNCTION = "FUNCTION"

// Function objects consist of a user-generated function.
type Function struct {
	Parameters  []*ast.Identifier
	Body        *ast.Block
	Defaults    map[string]ast.ExpressionNode
	Environment *Environment
}

// String represents the function object's value as a string.
func (function *Function) String() string {
	return "function"
}

// Type returns the function object type.
func (function *Function) Type() Type {
	return FUNCTION
}

// Method defines the set of methods available on function objects.
func (function *Function) Method(method string, args []Object) (Object, bool) {
	return nil, false
}

func (function *Function) Evaluate(args []Object, writer io.Writer) (Object, bool) {
	env := NewEnvironment()

	if writer != nil {
		env.SetWriter(writer)
	}

	result, ok := evaluator(function.Body, env)

	return result, ok
}
