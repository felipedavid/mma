package main

import (
	"fmt"
	"math"
)

const Eof = 0

const (
	TokenNone = iota
	TokenName
	TokenNumber
	TokenString
	TokenRegister
	TokenNewLine
	TokenComma
	TokenLeftParen
	TokenRightParen
)

type TokenKind uint8

var kindStr = []string{
	TokenNone:       "None",
	TokenName:       "TokenName",
	TokenNumber:     "TokenNumber",
	TokenString:     "TokenString",
	TokenRegister:   "TokenRegister",
	TokenNewLine:    "TokenNewLine",
	TokenComma:      "TokenComma",
	TokenLeftParen:  "TokenLeftParen",
	TokenRightParen: "TokenRightParen",
}

type Token struct {
	kind TokenKind
	val  any

	start  int
	finish int
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

var regNameToValue = map[string]uint16{
	"$zero": 0,
	"$0":    0,
	"$1":    1,
	"$2":    2,
	"$3":    3,
	"$4":    4,
	"$5":    5,
	"$6":    6,
	"$7":    7,
}

func (a *Assembler) lexeme() string {
	return string(a.source[a.start:a.current])
}

func (a *Assembler) nextToken() {
	a.token.kind = TokenNone
	a.token.val = nil

startOver:
	a.start = a.current
	ch := a.peek()
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
		newline := false
		for isWhiteSpace(a.peek()) {
			if a.peek() == '\n' {
				newline = true
				a.line++
			}
			a.advance()
		}
		if newline {
			a.token.kind = TokenNewLine
		} else {
			goto startOver
		}
	case ch == '$':
		a.advance()
		for isAlphaNumeric(a.peek()) {
			a.advance()
		}
		regName := a.lexeme()
		if val, isValidRegister := regNameToValue[regName]; isValidRegister {
			a.token.kind = TokenRegister
			a.token.val = val
		} else {
			a.lexError("%s is not a valid register", regName)
		}
	case ch == ',':
		a.token.kind = TokenComma
		a.advance()
	case ch == '(':
		a.token.kind = TokenLeftParen
		a.advance()
	case ch == ')':
		a.token.kind = TokenRightParen
		a.advance()
	case ch == Eof:
		return
	default:
		a.lexError("Unknown character: '%c'\n", ch)
	}
	a.token.start = a.start
	a.start = a.current
	a.token.finish = a.current
}

func (a *Assembler) scanName() {
	for isAlphaNumeric(a.peek()) || a.peek() == '_' {
		a.advance()
	}

	a.token.kind = TokenName
	a.token.val = a.lexeme()
}

func (a *Assembler) scanHexEscape() byte {
	a.advance()
	val := charToDigit[a.peek()]
	if val == 0 && a.peek() != '0' {
		a.lexError("\\x needs at least 1 hex digit")
	}
	a.advance()
	digit := charToDigit[a.advance()]
	if digit != 0 || a.advance() == '0' {
		val *= 16
		val += digit
		if val > 0xFF {
			a.lexError("\\x argument out of range")
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
	a.advance()
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
			a.lexError("Unexpected end of file within multi-line string literal")
		}
	} else {
		for !a.isAtEnd() && a.peek() != '"' {
			val := a.peek()
			if val == '\n' {
				a.lexError("String literal cannot contain newline")
				break
			} else if val == '\\' {
				a.advance()
				if a.peek() == 'x' {
					val = a.scanHexEscape()
				} else {
					val = escapeToChar[a.peek()]
					if val == 0 && a.peek() != '0' {
						a.lexError("Invalid string literal escape '\\%c'", a.peek())
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
			a.lexError("Unexpected end of file within string literal")
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
			a.lexError("Digit '%c' out of range for base %v", a.peek(), base)
			digit = 0
		}
		if val > (math.MaxInt64-int(digit))/base {
			a.lexError("Integer literal overflow")
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
		a.lexError("Expected base %v digit, got '%v'", base, a.peek())
	}

	a.token.kind = TokenNumber
	a.token.val = val
}

func (a *Assembler) lexError(fmtString string, val ...any) {
	errMsg := fmt.Sprintf(fmtString, val...)
	fmt.Printf("Lexing error on line %d: %s\n", a.line, errMsg)
}
