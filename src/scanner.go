package main

type Scanner struct {
    src      []byte
    ch       rune
    offset   int
    rdOffset int
}

const bom = 0xFEFF
func (s *Scanner) Init(src []byte) {
    s.src = src
    s.ch = ' '
    s.offset = 0
    s.rdOffset = 0

    s.next()
    if s.ch == bom {
        s.next()
    }
}

func (s *Scanner) next() {
    if s.rdOffset < len(s.src) {
        s.offset = s.rdOffset
        s.ch = rune(s.src[s.rdOffset])
        s.rdOffset += 1
    } else {
        s.offset = len(s.src)
        s.ch = -1 // EOF
    }
}

func (s *Scanner) skipWhitespace() {
    for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
        s.next()
    }
}

func (s *Scanner) scanComment() string {
    s.next()
    offs := s.offset
    for s.ch != '\n' && s.ch >= 0 {
        s.next()
    }
    return string(s.src[offs:s.offset])
}

func (s *Scanner) Scan() (tok Token, lit string) {
    s.skipWhitespace()
    ch := s.ch
    ch_str := string(ch)
    s.next()

    // Yeah, yeah, I know, this is kinda dumb. First time using Go. I know no better.
    // @todo Rewrite this crap - fdavid 11/08/2021
    switch ch {
    case -1:
        tok = EOF
        return
    case '.':
        if s.ch == 't' {
            tok = LABEL
            lit = ".text"
        } else if s.ch == 'd' {
            tok = LABEL
            lit = ".data"
        } else if s.ch == 'w' {
            tok = WORD
            lit = ch_str + s.scanLine()
        }
    case '/':
        if s.ch == '/' {
            tok = COMMENT
            lit = ch_str + s.scanComment()
        }
    case 'j':
        if s.ch == ' ' {
            tok = J_INSTRUCTION
            lit = ch_str + s.scanLine()
        }
    case 'l':
        if s.ch == 'w' {
            tok = I_INSTRUCTION
            lit = ch_str + s.scanLine()
        }
    case 's':
        if s.ch == 'w' {
            tok = I_INSTRUCTION
            lit = ch_str + s.scanLine()
        } else if s.ch == 'u' || s.ch == 'l' {
            tok = R_INSTRUCTION
            lit = ch_str + s.scanLine()
        }
    case 'a':
        if (s.ch == 'n' || s.ch == 'd') && s.src[s.offset+1] == 'd' {
            tok = R_INSTRUCTION
            lit = ch_str + s.scanLine()
        } else if s.ch == 'd' && s.src[s.offset+1] == 'd' && s.src[s.offset+2] == 'i' {
            tok = I_INSTRUCTION
            lit = ch_str + s.scanLine()
        }
    case 'b':
        if s.ch == 'e' && s.src[s.offset+1] == 'q' {
            tok = I_INSTRUCTION
            lit = ch_str + s.scanLine()
        }
    case '(':
		tok = LABEL
		lit = ch_str + s.scanLine()
    case 'm':
        if s.ch == 'o' && s.src[s.offset+1] == 'v' {
            tok = PSEUDO_INSTRUCTION
            lit = ch_str + s.scanLine()
        }
    default:
		lit = ch_str + s.scanLine()
        tok = ILLEGAL
    }
    return
}

func (s *Scanner) scanLine() string {
    offs := s.offset
    for s.ch != '\n' && s.ch != '\r' && s.ch >= 0 {// && s.ch != ' ' {
        s.next()
    }
    return string(s.src[offs:s.offset])
}

func (s *Scanner) scanLabel() string {
	offs := s.offset
	for {
		ch := s.ch
		if ch == '\n' || ch == '\r' || ch < 0 {
			break
		}
		s.next()
		if ch == ')' {
			break
		}
	}
	return string(s.src[offs:s.offset-1])
}
