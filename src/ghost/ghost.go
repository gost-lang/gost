package ghost

import (
	"ghostlang.org/x/ghost/interpreter"
	"ghostlang.org/x/ghost/library/modules"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/parser"
	"ghostlang.org/x/ghost/scanner"
)

type Ghost struct {
	FatalError  bool
	source      string
	Environment *object.Environment
	File        string
}

func New() *Ghost {
	ghost := &Ghost{
		Environment: object.NewEnvironment(),
	}

	ghost.RegisterEvaluator()

	return ghost
}

func (ghost *Ghost) SetDirectory(directory string) {
	ghost.Environment.SetDirectory(directory)
}

func (ghost *Ghost) GetDirectory() string {
	return ghost.Environment.GetDirectory()
}

func (ghost *Ghost) SetSource(source string) {
	ghost.source = source
}

func (ghost *Ghost) Execute() object.Object {
	scanner := scanner.New(ghost.source)
	parser := parser.New(scanner)
	program := parser.Parse()

	if len(parser.Errors()) != 0 {
		logParseErrors(parser.Errors())
		return nil
	}

	result := interpreter.Evaluate(program, ghost.Environment)

	if err, ok := result.(*object.Error); ok {
		log.Error(err.Message)

		return nil
	}

	return result
}

func (ghost *Ghost) RegisterEvaluator() {
	evaluator := interpreter.Evaluate

	object.RegisterEvaluator(evaluator)
	modules.RegisterEvaluator(evaluator)
}

func logParseErrors(errors []string) {
	for _, message := range errors {
		log.Error(message)
	}
}
