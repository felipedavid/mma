package main

import "fmt"

type Instruction interface {
	encode() uint16
}

type RInstruction struct {
	op    uint16
	rs    uint16
	rt    uint16
	rd    uint16
	funct uint16
}

func (i *RInstruction) encode() uint16 {
	registerMask := uint16((1 << 3) - 1)
	opMask := uint16((1 << 4) - 1)
	functMask := uint16((1 << 3) - 1)

	op := (i.op & opMask) << 12
	rs := (i.rs & registerMask) << 9
	rt := (i.rt & registerMask) << 6
	rd := (i.rd & registerMask) << 3
	funct := i.rd & functMask

	return op | rs | rt | rd | funct
}

type IInstruction struct {
	op   uint16
	rs   uint16
	rt   uint16
	immd uint16
}

func (i *IInstruction) encode() uint16 {
	registerMask := uint16((1 << 3) - 1)
	opMask := uint16((1 << 4) - 1)
	immdMask := uint16((1 << 6) - 1)

	op := (i.op & opMask) << 12
	rs := (i.rs & registerMask) << 9
	rt := (i.rt & registerMask) << 6
	immd := i.immd & immdMask

	return op | rs | rt | immd
}

type JInstruction struct {
	op   uint16
	addr uint16
}

func (i *JInstruction) encode() uint16 {
	opMask := uint16((1 << 4) - 1)
	addrMask := uint16((1 << 12) - 1)

	op := (i.op & opMask) << 12
	addr := i.addr & addrMask

	return op | addr
}

type Assembler struct {
	Lexer
}

func newAssembler(source []byte) *Assembler {
	return &Assembler{
		*newLexer(source),
	}
}

func (a *Assembler) assemble() error {
	a.next()
	switch a.token.kind {
	case Identifier:
		fmt.Println("Assembler command")
	}
	return nil
}
