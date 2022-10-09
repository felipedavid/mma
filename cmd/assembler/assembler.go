package main

type Assembler struct {
	filename string
	source   []byte

	token Token

	start   int
	current int
	line    int
	pass    int
	address int

	symbols map[string]Symbol
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
