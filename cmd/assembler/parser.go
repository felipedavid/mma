package main

import "fmt"

type Parser struct {
	l *Lexer
}

func newParser(fileName string, src []byte) *Parser {
	newL := newLexer(fileName, src)
	return &Parser{
		l: newL,
	}
}

func (p *Parser) isToken(kind TokenKind) bool {
	return p.l.Token.kind == kind
}

func (p *Parser) isTokenIdentifier(ident string) bool {
	return p.isToken(TokenIdentifier) && string(p.l.Token.strVal) == ident
}

func (p *Parser) matchToken(kind TokenKind) bool {
	if p.isToken(kind) {
		p.l.NextToken()
		return true
	}
	return false
}

func (p *Parser) expectToken(kind TokenKind) bool {
	if p.isToken(kind) {
		p.l.NextToken()
		return true
	}
	p.error("Expected token %s, got %s", TokenKindToString[kind], TokenKindToString[p.l.Token.kind])
	return false
}

func (p *Parser) error(fmtString string, val ...any) {
	fmt.Printf(fmtString, val...)
}
