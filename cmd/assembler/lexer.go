package main

import (
	"fmt"
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

const (
	None = iota
	TokenIdentifier
	TokenInteger
	TokenString
	TokenDivOp
	End
)

var TokenKindToString = []string{
	None:            "None",
	TokenIdentifier: "TokenIdentifier",
	TokenInteger:    "TokenInteger",
	TokenString:     "TokenString",
	TokenDivOp:      "TokenDivOp",
	End:             "End",
}

type TokenKind uint

// Token just assign a meaning to a sequence of characters. It's easier to work with tokens then every time we get
// a sequence of characters make some sense out of it
type Token struct {
	kind   TokenKind
	intVal int64
	strVal []byte
}

func (t *Token) GetValue() any {
	switch t.kind {
	case TokenInteger:
		return t.intVal
	case TokenIdentifier:
		fallthrough
	case TokenString:
		return string(t.strVal)
	}
	return nil
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
func NewLexer(fileName string, source []byte) *Lexer {
	lex := &Lexer{
		fileName:   fileName,
		src:        source,
		lineNumber: 1,
		start:      0,
		current:    0,
		errLogger:  log.New(os.Stdout, "[SYNTAX_ERROR] ", 0),
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
		l.start = l.current
		// Yes, gotos. This is actually a good use case for them
		goto StartOver
	case '0':
		fallthrough
	case '1':
		fallthrough
	case '2':
		fallthrough
	case '3':
		fallthrough
	case '4':
		fallthrough
	case '5':
		fallthrough
	case '6':
		fallthrough
	case '7':
		fallthrough
	case '8':
		fallthrough
	case '9':
		l.scanInt()
	case '"':
		l.scanString()
	case '\'':
		l.scanChar()
	case 'a':
		fallthrough
	case 'b':
		fallthrough
	case 'c':
		fallthrough
	case 'd':
		fallthrough
	case 'e':
		fallthrough
	case 'f':
		fallthrough
	case 'g':
		fallthrough
	case 'h':
		fallthrough
	case 'i':
		fallthrough
	case 'j':
		fallthrough
	case 'k':
		fallthrough
	case 'l':
		fallthrough
	case 'm':
		fallthrough
	case 'n':
		fallthrough
	case 'o':
		fallthrough
	case 'p':
		fallthrough
	case 'q':
		fallthrough
	case 'r':
		fallthrough
	case 's':
		fallthrough
	case 't':
		fallthrough
	case 'u':
		fallthrough
	case 'v':
		fallthrough
	case 'w':
		fallthrough
	case 'x':
		fallthrough
	case 'y':
		fallthrough
	case 'z':
		fallthrough
	case 'A':
		fallthrough
	case 'B':
		fallthrough
	case 'C':
		fallthrough
	case 'D':
		fallthrough
	case 'E':
		fallthrough
	case 'F':
		fallthrough
	case 'G':
		fallthrough
	case 'H':
		fallthrough
	case 'I':
		fallthrough
	case 'J':
		fallthrough
	case 'K':
		fallthrough
	case 'L':
		fallthrough
	case 'M':
		fallthrough
	case 'N':
		fallthrough
	case 'O':
		fallthrough
	case 'P':
		fallthrough
	case 'Q':
		fallthrough
	case 'R':
		fallthrough
	case 'S':
		fallthrough
	case 'T':
		fallthrough
	case 'U':
		fallthrough
	case 'V':
		fallthrough
	case 'W':
		fallthrough
	case 'X':
		fallthrough
	case 'Y':
		fallthrough
	case fallthrough
	case 'g':
		fallthrough
	case 'h':
		fallthrough
	casefallthrough
	case 'g':
		fallthrough
	case 'h':
		fallthrough
	case'Z':
		l.current++
		for isAlphaNumeric(l.src[l.current]) || l.src[l.current] == '_' {
			l.current++
		}
		l.Token.strVal = l.src[l.start:l.current]
		l.Token.kind = TokenIdentifier
	case '/':
		l.current++
		if l.src[l.current] == '/' { // Check if it is a one line comment
			for l.current < len(l.src) && l.src[l.current] != '\n' {
				l.current++
			}
			l.start = l.current
			goto StartOver
		} else if l.src[l.current] == '*' { // Check if it is a multi-line comment
			for l.current < len(l.src) {
				if l.src[l.current] == '*' && l.src[l.current+1] == '/' {
					l.current += 2
					l.start = l.current
					goto StartOver
				}
				l.current++
			}
			l.error("Unterminated multi-line comment")
		} else {
			l.Token.kind = TokenDivOp
		}
	default:
	}

	l.start = l.current
}

func (l *Lexer) scanString() {
	l.current++
	buf := make([]byte, 0, 16)
	if l.src[0] == '"' && l.src[1] == '"' {
		l.current++
		for l.current < len(l.src) {
			if l.src[0] == '"' && l.src[1] == '"' && l.src[2] == '"' {
				l.current += 3
				break
			}
			if l.src[l.current] != '\r' {
				buf = append(buf, l.src[l.current])
			}
			if l.src[l.current] == '\n' {
				l.lineNumber++
			}
			l.current++
		}
		if l.current < len(l.src) {
			l.error("Unexpected end of file within multi-line string literal")
		}
	} else {
		for l.current < len(l.src) && l.src[l.current] != '"' {
			val := l.src[l.current]
			if val == '\n' {
				l.error("String literal cannot contain newline")
				break
			} else if val == '\\' {
				l.current++
				if l.src[l.current] == 'x' {
					val = byte(l.scanHexEscape())
				} else {
					val = escapeToChar[l.src[l.current]]
					if val == 0 && l.src[l.current] != '0' {
						l.error("Invalid string literal escape '\\%c'", l.src[l.current])
					}
					l.current++
				}
			} else {
				l.current++
			}
			buf = append(buf, val)
		}
		if l.current < len(l.src) {
			l.current++
		} else {
			l.error("Unexpected end of file within string literal")
		}
	}

	l.Token.kind = TokenString
	l.Token.strVal = buf
}

func (l *Lexer) scanChar() {
	l.current++
	val := byte(0)
	if l.src[l.current] == '\'' {
		l.error("Char literal cannot be empty")
		l.current++
	} else if l.src[l.current] == '\n' {
		l.error("Char literal cannot contain newline")
	} else if l.src[l.current] == '\\' {
		l.current++
		if l.src[l.current] == 'x' {
			val = byte(l.scanHexEscape())
		} else {
			val = escapeToChar[l.src[l.current]]
			if val == 0 && l.src[l.current] != '0' {
				l.error("Invalid char literal escape '\\%c'", l.src[l.current])
			}
			l.current++
		}
	} else {
		val = l.src[l.current]
		l.current++
	}
	if l.src[l.current] != '\'' {
		l.error("Expected closing char quote, got '%c'", l.src[l.current])
	} else {
		l.current++
	}
	l.Token.kind = TokenInteger
	l.Token.intVal = int64(val)
}

func (l *Lexer) scanHexEscape() int {
	l.current++
	val := int(charToDigit[l.src[l.current]])
	if val == 0 && l.src[l.current] != '0' {
		l.error("\\x needs at least 1 hex digit")
	}
	l.current++
	digit := charToDigit[l.src[l.current]]
	if digit != 0 || l.src[l.current] == '0' {
		val *= 16
		val += int(digit)
		if val > 0xFF {
			l.error("\\x argument out of range")
			val = 0xFF
		}
		l.current++
	}
	return val
}

// scanInt parses integers from base 2, 8, 10 and 16
func (l *Lexer) scanInt() {
	l.Token.kind = TokenInteger

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
			l.error("Digit '%c' out of range for base %v", l.src[l.current], base)
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

	l.Token.intVal = val
}

// error logs errors while lexing
func (l *Lexer) error(fmtString string, val ...any) {
	// Get the whole lexeme to print with the error message
	lexemeEnd := l.current
	for lexemeEnd < len(l.src) && !isSpace(l.src[lexemeEnd]) {
		lexemeEnd++
	}

	errorMsg := fmt.Sprintf(fmtString, val...)
	l.errLogger.Printf("[Line %d] [Lexeme: %s] %s", l.lineNumber, l.src[l.start:lexemeEnd], errorMsg)
}
