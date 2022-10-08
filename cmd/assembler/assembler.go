package main

import (
	"fmt"
	"math"
)

const (
	None = iota
	TokenName
	TokenNumber
	TokenString
)

type TokenKind uint8

type Token struct {
	kind TokenKind
	val  any
}

type Assembler struct {
	filename string
	source   []byte

	token Token

	start   int
	current int
	line    int
}

func newAssembler(filename string, source []byte) *Assembler {
	asm := &Assembler{
		filename: filename,
		source:   source,
	}
	asm.nextToken()
	return asm
}

func (a *Assembler) isAtEnd() bool {
	return a.current >= len(a.source)
}

func (a *Assembler) peek() byte {
	if a.isAtEnd() {
		return 0
	}
	curr := a.source[a.current]
	return curr
}

func (a *Assembler) advance() byte {
	ch := a.peek()
	a.current++
	return ch
}

func (a *Assembler) nextToken() {
startOver:
	a.start = a.current
	ch := a.advance()
	switch {
	case ch >= 'A' && ch <= 'Z':
		fallthrough
	case ch >= 'a' && ch <= 'z':
		a.scanName()
	case ch >= '0' && ch <= '9':
		a.scanNumber()
	case ch == '"':
		a.scanString()
	case isWhiteSpace(ch):
		for isWhiteSpace(a.peek()) {
			a.advance()
		}
		goto startOver
	default:
		a.error("Unknown character: %c\n", ch)
	}
	a.start = a.current
}

func (a *Assembler) scanName() {
	for isAlphaNumeric(a.peek()) || a.peek() == '_' {
		a.advance()
	}

	a.token.kind = TokenName
	a.token.val = string(a.source[a.start:a.current])
}

func (a *Assembler) scanHexEscape() byte {
	a.advance()
	val := charToDigit[a.peek()]
	if val == 0 && a.peek() != '0' {
		a.error("\\x needs at least 1 hex digit")
	}
	a.advance()
	digit := charToDigit[a.advance()]
	if digit != 0 || a.advance() == '0' {
		val *= 16
		val += digit
		if val > 0xFF {
			a.error("\\x argument out of range")
			val = 0xFF
		}
		a.advance()
	}
	return val
}

var escapeToChar = [256]byte{
	'\'': '\'',
	'"':  '"',
	'\\': '\\',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'v':  '\v',
	'b':  '\b',
	'a':  '\a',
	'0':  0,
}

func (a *Assembler) scanString() {
	buf := make([]byte, 0, 16)
	if a.source[0] == '"' && a.source[1] == '"' {
		a.advance()
		for !a.isAtEnd() {
			if a.source[0] == '"' && a.source[1] == '"' && a.source[2] == '"' {
				a.current += 3
				break
			}
			if a.peek() != '\r' {
				buf = append(buf, a.peek())
			}
			if a.peek() == '\n' {
				a.line++
			}
			a.advance()
		}
		if !a.isAtEnd() {
			a.error("Unexpected end of file within multi-line string literal")
		}
	} else {
		for !a.isAtEnd() && a.peek() != '"' {
			val := a.peek()
			if val == '\n' {
				a.error("String literal cannot contain newline")
				break
			} else if val == '\\' {
				a.advance()
				if a.peek() == 'x' {
					val = a.scanHexEscape()
				} else {
					val = escapeToChar[a.peek()]
					if val == 0 && a.peek() != '0' {
						a.error("Invalid string literal escape '\\%c'", a.peek())
					}
					a.advance()
				}
			} else {
				a.advance()
			}
			buf = append(buf, val)
		}
		if !a.isAtEnd() {
			a.advance()
		} else {
			a.error("Unexpected end of file within string literal")
		}
	}

	a.token.kind = TokenString
	a.token.val = string(buf)
}

var charToDigit = [256]byte{
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

func (a *Assembler) scanNumber() {
	base := 10
	startDigits := a.current
	// Check for binary/hex/octal notations
	if a.peek() == '0' {
		a.advance()
		if toLower(a.peek()) == 'x' { // Hex
			a.advance()
			base = 16
			startDigits = a.current
		} else if toLower(a.peek()) == 'b' { // Binary
			a.advance()
			base = 2
			startDigits = a.current
		} else if isDigit(a.peek()) { // Octal
			base = 8
			startDigits = a.current
		}
	}

	// Actually parse the number based adjusting for the base
	var val int
	for !a.isAtEnd() {
		if a.peek() == '_' {
			a.advance()
			continue
		}
		digit := charToDigit[a.peek()]
		if digit == 0 && a.peek() != '0' {
			break
		}
		if int(digit) >= base {
			a.error("Digit '%c' out of range for base %v", a.peek(), base)
			digit = 0
		}
		if val > (math.MaxInt64-int(digit))/base {
			a.error("Integer literal overflow")
			for isDigit(a.peek()) {
				a.advance()
			}
			val = 0
			break
		}
		val = val*base + int(digit)
		a.advance()
	}
	if a.current == startDigits {
		a.error("Expected base %v digit, got '%v'", base, a.peek())
	}

	a.token.kind = TokenNumber
	a.token.val = val
}

func (a *Assembler) error(fmtString string, val ...any) {
	errMsg := fmt.Sprintf(fmtString, val...)
	fmt.Printf("[line %d] %s", a.line, errMsg)
}
