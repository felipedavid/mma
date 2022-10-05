package main

import (
	"log"
	"math"
	"os"
)

// charToDigit maps a ASCII character to its actual value on 2-16 numeric bases
var charToDigit = [256]uint8{
	'0': 0,
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'a': 10,
	'b': 11,
	'c': 12,
	'd': 13,
	'e': 14,
	'f': 15,
	'A': 10,
	'B': 11,
	'C': 12,
	'D': 13,
	'E': 14,
	'F': 15,
}

const (
	None = iota
	Identifier
	Integer
	End
)

type TokenKind uint

// Token just assign a meaning to a sequence of characters. It's easier to work with tokens then every time we get
// a sequence of characters make some sense out of it
type Token struct {
	kind   TokenKind
	intVal int64
}

// Lexer takes a byte stream and transform it in a Token stream
type Lexer struct {
	fileName string
	src      []byte // src is our raw source code read from a file
	Token    Token

	lineNumber int // lineNumber just tells us which line we are scanning at a given time
	lineStart  int // lineStart tells where is the start of the line specified by lineNumber
	start      int // start specifies the start of a lexeme
	current    int // current specified the current character we are checking

	errLogger *log.Logger
}

// NewLexer just returns a initialized instance of Lexer
func NewLexer(source []byte) *Lexer {
	lex := &Lexer{
		src:        source,
		lineNumber: 1,
		start:      0,
		current:    0,
		errLogger:  log.New(os.Stdout, "SYNTAX_ERROR\t", log.LstdFlags),
	}

	lex.NextToken()
	return lex
}

// NextToken consumes a token
func (l *Lexer) NextToken() {
	if l.start >= len(l.src) {
		l.Token.kind = End
		return
	}

StartOver:
	switch l.src[l.current] {
	case ' ':
		fallthrough
	case '\t':
		fallthrough
	case '\n':
		fallthrough
	case '\r':
		fallthrough
	case '\v':
		for isSpace(l.src[l.current]) {
			if l.src[l.current] == '\n' {
				l.lineStart = l.current + 1
				l.lineNumber++
			}
			l.current++
		}
		// Yes, gotos. This is actually a good use case for them
		goto StartOver
	default:
		if isDigit(l.src[l.current]) {
			l.Token.kind, l.Token.intVal = l.scanInt()
		}
	}

	l.start = l.current
}

// scanInt parses integers from base 2, 8, 10 and 16
func (l *Lexer) scanInt() (TokenKind, int64) {
	base := int64(10)
	startDigits := l.current
	// Check for binary/hex/octal notations
	if l.src[l.current] == '0' {
		l.current++
		if toLower(l.src[l.current]) == 'x' {
			l.current++
			base = 16
			startDigits = l.current
		} else if toLower(l.src[l.current]) == 'b' {
			l.current++
			base = 2
			startDigits = l.current
		} else if isDigit(l.src[l.current]) {
			base = 8
			startDigits = l.current
		}
	}

	// Actually parse the number based adjusting for the base
	var val int64
	for l.current < len(l.src) {
		if l.src[l.current] == '_' {
			l.current++
			continue
		}
		digit := int64(charToDigit[l.src[l.current]])
		if digit == 0 && l.src[l.current] != '0' {
			break
		}
		if digit >= base {
			l.error("Digit '%v' out of range for base %v", l.src[l.current], base)
			digit = 0
		}
		if val > (math.MaxInt64-digit)/base {
			l.error("Integer literal overflow")
			for isDigit(l.src[l.current]) {
				l.current++
			}
			val = 0
			break
		}
		val = val*base + digit
		l.current++
	}
	if l.current == startDigits {
		l.error("Expected base %v digit, got '%v'", base, l.src[l.current])
	}

	return Integer, val
}

// error logs errors while lexing
func (l *Lexer) error(fmtString string, val ...any) {
	l.errLogger.Fatalf(fmtString, val...)
}
