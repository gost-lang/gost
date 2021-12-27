package parser

import (
	"ghostlang.org/x/ghost/ast"
	"ghostlang.org/x/ghost/token"
)

func (parser *Parser) parseExpression(precedence int) ast.ExpressionNode {
	prefix := parser.prefixParserFns[parser.currentToken.Type]

	if prefix == nil {
		return nil
	}

	leftExpression := prefix()

	for precedence < parser.nextTokenPrecedence() {
		infix := parser.infixParserFns[parser.nextToken.Type]

		if infix == nil {
			return leftExpression
		}

		parser.readToken()

		leftExpression = infix(leftExpression)
	}

	return leftExpression
}

func (parser *Parser) parseExpressionList(end token.Type) []ast.ExpressionNode {
	list := []ast.ExpressionNode{}

	if parser.nextTokenTypeIs(end) {
		parser.readToken()

		return list
	}

	parser.readToken()

	list = append(list, parser.parseExpression(LOWEST))

	for parser.nextTokenTypeIs(token.COMMA) {
		parser.readToken()
		parser.readToken()
		list = append(list, parser.parseExpression(LOWEST))
	}

	if !parser.expectNextType(end) {
		return nil
	}

	return list
}
