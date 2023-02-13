package main

import "fmt"

const (
	InvalidToken = iota
	Dot
	Comma
	Colon
	DollarSign
	Identifier
	Number
	Eof
)

type TokenKind int

type Token struct {
	kind   TokenKind
	lexeme string
	line   int
}

type Lexer struct {
	source []byte

	offset int
	token  Token
	line   int
}

func newLexer(source []byte) *Lexer {
	return &Lexer{
		source: source,
		line:   1,
	}
}

func (l *Lexer) newToken(kind TokenKind, startOffset int) Token {
	return Token{kind: kind, lexeme: string(l.source[startOffset:l.offset])}
}

func (l *Lexer) peek() byte {
	if l.offset >= len(l.source) {
		return 0
	}
	return l.source[l.offset]
}

func (l *Lexer) next() bool {
REPEAT:
	ch := l.peek()
	start := l.offset
	l.offset++
	switch {
	case ch == 0:
		l.token = Token{kind: Eof}
	case ch == '.':
		l.token = l.newToken(Dot, start)
	case ch == ',':
		l.token = l.newToken(Comma, start)
	case ch == ':':
		l.token = l.newToken(Colon, start)
	case ch == '$':
		l.token = l.newToken(DollarSign, start)
	case ch == '\n':
		l.line++
		goto REPEAT
	case ch == ' ':
		for isSpace(l.peek()) {
			l.offset++
		}
		goto REPEAT
	case isAlpha(ch) || ch == '_':
		ch = l.peek()
		for isAlphaNumeric(ch) || ch == '_' {
			l.offset++
			ch = l.peek()
		}
		l.token = l.newToken(Identifier, start)
	case isDigit(ch):
		for isDigit(l.peek()) {
			l.offset++
		}
		l.token = l.newToken(Number, start)
	case ch == '#':
		for l.peek() != '\n' {
			l.offset++
		}
		start = l.offset
		goto REPEAT
	default:
		l.token = l.newToken(InvalidToken, start)
		return false
	}
	return true
}

func (l *Lexer) logCurrentToken() {
	var tokenKindStr = []string{
		InvalidToken: "InvalidToken",
		Dot:          "Dot",
		Comma:        "Comma",
		DollarSign:   "DollarSign",
		Colon:        "Colon",
		Identifier:   "Identifier",
		Number:       "Number",
		Eof:          "Eof",
	}

	fmt.Printf("line: %d, kind: %s, lexeme: \"%s\"\n", l.line, tokenKindStr[l.token.kind], string(l.token.lexeme))
}
