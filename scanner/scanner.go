package scanner

import (
	"fmt"

	"ghostlang.org/x/ghost/error"
	"ghostlang.org/x/ghost/log"
	"ghostlang.org/x/ghost/token"
	"github.com/shopspring/decimal"
)

// Scanner transforms our source code into tokens.
type Scanner struct {
	source  string
	start   int
	current int
	line    int
	tokens  []token.Token
}

// keywords contains a list of all reserved keywords.
var keywords = map[string]token.Type{
	"and":      token.AND,
	"class":    token.CLASS,
	"else":     token.ELSE,
	"false":    token.FALSE,
	"for":      token.FOR,
	"function": token.FUNCTION,
	"if":       token.IF,
	"null":     token.NULL,
	"or":       token.OR,
	// "print":    token.PRINT,
	"return": token.RETURN,
	"super":  token.SUPER,
	"this":   token.THIS,
	"true":   token.TRUE,
	"while":  token.WHILE,
}

// New creates a new scanner instance.
func New(source string) Scanner {
	scanner := Scanner{source: source, line: 1}

	return scanner
}

// ScanTokens transforms the source into an array of tokens. It works its way
// through the source code, adding tokens until it runs out of characters.
// Then it appends one final "end of file" token.
func (scanner *Scanner) ScanTokens() []token.Token {
	for !scanner.isAtEnd() {
		// We are at the beginning of the next lexeme.
		scanner.start = scanner.current
		scanner.ScanToken()
	}

	scanner.tokens = append(scanner.tokens, token.Token{Type: token.EOF, Literal: nil, Line: scanner.line})

	return scanner.tokens
}

// scanToken is responsible for scanning the current character and storing the
// correct token type for it. This is the heart of our scanner.
func (scanner *Scanner) ScanToken() {
	c := scanner.advance()

	switch c {
	case '(':
		scanner.addToken(token.LEFTPAREN)
	case ')':
		scanner.addToken(token.RIGHTPAREN)
	case '{':
		scanner.addToken(token.LEFTBRACE)
	case '}':
		scanner.addToken(token.RIGHTBRACE)
	case ',':
		scanner.addToken(token.COMMA)
	case '.':
		scanner.addToken(token.DOT)
	case '-':
		scanner.addToken(token.MINUS)
	case '+':
		scanner.addToken(token.PLUS)
	case ';':
		scanner.addToken(token.SEMICOLON)
	case '*':
		scanner.addToken(token.STAR)
	case '%':
		scanner.addToken(token.PERCENT)
	case '?':
		scanner.addToken(token.QUESTION)
	case ':':
		if scanner.match('=') {
			scanner.addToken(token.ASSIGN)
		} else {
			scanner.addToken(token.COLON)
		}
	case '!':
		if scanner.match('=') {
			scanner.addToken(token.BANGEQUAL)
		} else {
			scanner.addToken(token.BANG)
		}
	case '=':
		if scanner.match('=') {
			scanner.addToken(token.EQUALEQUAL)
		} else {
			scanner.addToken(token.EQUAL)
		}
	case '<':
		if scanner.match('=') {
			scanner.addToken(token.LESSEQUAL)
		} else {
			scanner.addToken(token.LESS)
		}
	case '>':
		if scanner.match('=') {
			scanner.addToken(token.GREATEREQUAL)
		} else {
			scanner.addToken(token.GREATER)
		}
	case '/':
		if scanner.match('/') {
			// A comment goes until the end of the line. Comments are lexemes
			// but they aren't meaningful, so we don't want to deal with them
			// and simply wish to discard them here.
			for scanner.peek() != '\n' && !scanner.isAtEnd() {
				scanner.advance()
			}
		} else {
			scanner.addToken(token.SLASH)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace
	case '\n':
		scanner.line++
	case '"':
		scanner.scanString()
	default:
		if scanner.isDigit(c) {
			scanner.scanNumber()
		} else if scanner.isAlpha(c) {
			scanner.scanIdentifier()
		} else {
			err := error.Error{
				Reason:  error.Syntax,
				Message: fmt.Sprintf("unexpected parsing error on line %d", scanner.line),
			}

			log.Error(err.Reason, err.Message)
		}
	}
}

// scanString consumes characters until it hits either the closing " or end of
// file. If we run to the end of the file without a closing ", we report an
// error.
func (scanner *Scanner) scanString() {
	for scanner.peek() != '"' && !scanner.isAtEnd() {
		if scanner.peek() == '\n' {
			scanner.line++
		}

		scanner.advance()
	}

	if scanner.isAtEnd() {
		err := error.Error{
			Reason:  error.Syntax,
			Message: fmt.Sprintf("unexpected end of string on line %d", scanner.line),
		}

		log.Error(err.Reason, err.Message)
		return
	}

	// The closing ".
	scanner.advance()

	// Extract and trim the surrounding quotes for our literal value.
	literal := scanner.source[scanner.start+1 : scanner.current-1]

	scanner.addTokenWithLiteral(token.STRING, literal)
}

// scanNumber consumes all digits for the integer part of the literal, and then
// the fractional part if we encounter a decimal point (.) followed by at least
// one digit. If we do have a fractional part, we consume all remaining digits.
func (scanner *Scanner) scanNumber() {
	for scanner.isDigit(scanner.peek()) {
		scanner.advance()
	}

	// Look for a fractional part.
	if scanner.peek() == '.' && scanner.isDigit(scanner.peekNext()) {
		// Consume the "."
		scanner.advance()

		for scanner.isDigit(scanner.peek()) {
			scanner.advance()
		}
	}

	// Look for a scientific notion part.
	if scanner.peek() == 'e' {
		// Consume the "e"
		scanner.advance()

		if scanner.peek() == '-' {
			// Consume the "-"
			scanner.advance()
		}

		if scanner.isDigit(scanner.peek()) {
			for scanner.isDigit(scanner.peek()) {
				scanner.advance()
			}
		} else {
			err := error.Error{
				Reason:  error.Syntax,
				Message: fmt.Sprintf("exponent has no digits on line %d", scanner.line),
			}

			log.Error(err.Reason, err.Message)
			return
		}
	}

	literal := scanner.source[scanner.start:scanner.current]
	number, err := decimal.NewFromString(literal)

	if err != nil {
		err := error.Error{
			Reason:  error.Syntax,
			Message: fmt.Sprintf("invalid number format on line %d", scanner.line),
		}

		log.Error(err.Reason, err.Message)
	} else {
		scanner.addTokenWithLiteral(token.NUMBER, number)
	}
}

// scanIdentifier consumes characters until it runs out of alphanumeric
// characters.
func (scanner *Scanner) scanIdentifier() {
	for scanner.isAlphaNumeric(scanner.peek()) {
		scanner.advance()
	}

	literal := scanner.source[scanner.start:scanner.current]
	keyword, ok := keywords[literal]

	if ok {
		scanner.addToken(keyword)
	} else {
		scanner.addToken(token.IDENTIFIER)
	}
}

// =============================================================================
// Helper methods

// addToken grabs the current lexeme and creates a new token for it. In this
// case, addToken is for tokens without a literal (native Go) value.
func (scanner *Scanner) addToken(tokenType token.Type) {
	scanner.addTokenWithLiteral(tokenType, nil)
}

// addTokenWithLiteral grabs the current lexeme and creates a new token of the
// passed type and literal (native Go) value.
func (scanner *Scanner) addTokenWithLiteral(tokenType token.Type, literal interface{}) {
	lexeme := scanner.source[scanner.start:scanner.current]
	scanner.tokens = append(scanner.tokens, token.Token{Type: tokenType, Lexeme: lexeme, Literal: literal, Line: scanner.line})
}

// isAtEnd tells us if we've consumed all the characters in our raw source code.
func (scanner *Scanner) isAtEnd() bool {
	return scanner.current >= len(scanner.source)
}

// isDigit tells us if the passed character is a number.
func (scanner *Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// isAlpha tells us if the passed character is an alphabetic character.
func (scanner *Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

// isAlphaNumeric tells us if the passed character is either an alphabetic or
// numeric character.
func (scanner *Scanner) isAlphaNumeric(c byte) bool {
	return scanner.isAlpha(c) || scanner.isDigit(c)
}

// advance consumes the next character in our raw source code and returns it.
func (scanner *Scanner) advance() byte {
	scanner.current++

	return scanner.source[scanner.current-1]
}

// match acts as a conditional advance, only consuming the current character if
// it's what we're looking for in "expected".
func (scanner *Scanner) match(expected byte) bool {
	if scanner.isAtEnd() {
		return false
	}

	if scanner.source[scanner.current] != expected {
		return false
	}

	scanner.current++

	return true
}

// peek looks at the current unconsumed character. We use this to lookahead,
// useful to check for multi-character tokens.
func (scanner *Scanner) peek() byte {
	if scanner.isAtEnd() {
		return 0
	}

	return scanner.source[scanner.current]
}

// peekNext operates in a similar manner to peek() however, it instead looks at
// the next upcoming character in our raw source code.
func (scanner *Scanner) peekNext() byte {
	if scanner.current+1 >= len(scanner.source) {
		return 0
	}

	return scanner.source[scanner.current+1]
}
