package main

type Assembler struct {
	filename string
	source   []byte

	token Token

	start   int
	current int
	line    int
	pass    int
	address uint16

	symbols map[string]Symbol

	dataSection []byte
	codeSection []byte

	hasError bool
}

func newAssembler(filename string, source []byte) *Assembler {
	asm := &Assembler{
		filename: filename,
		source:   source,

		start:   0,
		current: 0,
		line:    1,

		symbols: Symbols,
	}
	asm.nextToken()
	return asm
}

func (a *Assembler) resetState(src []byte) {
	a.source = src
	a.start = 0
	a.current = 0
	a.line = 1
	a.codeSection = make([]byte, 0)
	a.hasError = false
	a.nextToken()
}
