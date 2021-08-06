package main

type Scanner struct {
    src      []byte
    ch       rune // current character
    offset   int  // character offset
    rdOffset int  // reading offset
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
    src_len := len(s.src)
    if s.rdOffset < src_len {
        s.offset = s.rdOffset
        s.ch = rune(s.src[s.rdOffset])
        s.rdOffset += 1
    } else {
        s.offset = src_len
    }
}

func (s *Scanner) skipWhitespace() {
    if s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
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
    s.next() // always make progress

    switch ch {
        case -1:
            tok = EOF
        case '/':
            if s.ch == '/' { // the second '/' means we have a comment
                tok = COMMENT
                lit = s.scanComment()
            }
        case 'j':
            if s.ch == ' ' {
                tok = J_INSTRUCTION
            }
        default:
            tok = ILLEGAL
            lit = s.scanComment()
    }
    return tok, lit
}

func (s *Scanner) scanLine() string {
    offs := s.offset
    for s.ch != '\n' && s.ch != '\r' && s.ch >= 0 && s.ch != ' ' {
        s.next()
    }
    return string(s.src[offs:s.offset])
}
