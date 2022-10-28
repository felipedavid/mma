package main

import (
	"fmt"
)

const (
	ImmdMin = 0b000000
	ImmdMax = 0b111111
	AddrMin = 0b0000_0000_0000
	AddrMax = 0b1111_1111_1111
)

const (
	SymbolNone = iota
	SymbolReg
	SymbolInstr
	SymbolCmd
)

type SymbolKind uint8

// This is based on the instruction specific syntax and not on the architecture specification
const (
	InstrNoneKind     = iota
	InstrRegisterKind // <op> <reg>, <reg>, <reg>
	InstrMemoryKind   // <op> <reg>, <int>(<reg>)
	InstrImmKind      // <op> <reg>, <reg>, <int>
	InstrJumpKind
)

type InstrKind uint8

const (
	AddOp = iota
	AddiOp
	AndOp
	BeqOp
	JOp
	LwOp
	NandOp
	NorOp
	OrOp
	SltOp
	SwOp
	SubOp
)

type InstrDef struct {
	kind InstrKind
	op   uint16
}

type Symbol struct {
	name string
	kind SymbolKind

	val any
}

var Symbols = map[string]Symbol{
	// Instructions
	"add":  {"add", SymbolInstr, InstrDef{InstrRegisterKind, AddOp}},
	"addi": {"addi", SymbolInstr, InstrDef{InstrImmKind, AddiOp}},
	"and":  {"and", SymbolInstr, InstrDef{InstrRegisterKind, AndOp}},
	"beq":  {"beq", SymbolInstr, InstrDef{InstrImmKind, BeqOp}},
	"j":    {"j", SymbolInstr, InstrDef{InstrJumpKind, JOp}},
	"lw":   {"lw", SymbolInstr, InstrDef{InstrMemoryKind, LwOp}},
	"nand": {"nand", SymbolInstr, InstrDef{InstrRegisterKind, NandOp}},
	"nor":  {"nor", SymbolInstr, InstrDef{InstrRegisterKind, NorOp}},
	"or":   {"or", SymbolInstr, InstrDef{InstrRegisterKind, OrOp}},
	"slt":  {"slw", SymbolInstr, InstrDef{InstrRegisterKind, SltOp}},
	"sw":   {"sw", SymbolInstr, InstrDef{InstrMemoryKind, SwOp}},
	"sub":  {"sub", SymbolInstr, InstrDef{InstrRegisterKind, SubOp}},

	// Cmds
	".uint8":  {".uint8", SymbolCmd, cmdUint8},
	".uint16": {".uint16", SymbolCmd, cmdUint16},
	".uint32": {".uint32", SymbolCmd, cmdUint32},

	// Regs
	"$0": {"$0", SymbolReg, uint16(0)},
	"$1": {"$1", SymbolReg, uint16(1)},
	"$2": {"$2", SymbolReg, uint16(2)},
	"$3": {"$3", SymbolReg, uint16(3)},
	"$4": {"$4", SymbolReg, uint16(4)},
	"$5": {"$5", SymbolReg, uint16(5)},
	"$6": {"$6", SymbolReg, uint16(6)},
	"$7": {"$7", SymbolReg, uint16(7)},
}

func cmdUint8() {

}

func cmdUint16() {

}

func cmdUint32() {

}

func (a *Assembler) parseName() string {
	name := string(a.source[a.token.start:a.token.finish])
	//a.expectToken(TokenName)
	a.nextToken()
	return name
}

func (a *Assembler) parseSymbol() Symbol {
	name := a.parseName()
	sym, ok := a.symbols[name]
	if !ok {
		a.symbols[name] = sym
	}
	return sym
}

func (a *Assembler) isToken(kind TokenKind) bool {
	return a.token.kind == kind
}

func (a *Assembler) isTokenName(name string) bool {
	return a.token.val == name
}

func (a *Assembler) matchToken(kind TokenKind) bool {
	if a.token.kind == kind {
		a.nextToken()
		return true
	}
	return false
}

func (a *Assembler) matchName(name string) bool {
	if a.token.val == name {
		a.nextToken()
		return true
	}
	return false
}

func (a *Assembler) expectToken(kind TokenKind) bool {
	if a.isToken(kind) {
		a.nextToken()
		return true
	}
	a.parserError("Expected %s, got %s", kindStr[kind], kindStr[a.token.kind])
	return false
}

func (a *Assembler) parseNewlines() {
	for a.matchToken(TokenNewLine) {
	}
}

func (a *Assembler) parseRegister() uint16 {
	sym := a.parseSymbol()
	if sym.kind != SymbolReg {
		a.parserError("Expected register bug got '%s'", sym.name)
		return 0
	}

	val, isUint16 := sym.val.(uint16)
	if !isUint16 {
		a.parserError("Val should be uint16, not %T", val)
	}
	return val
}

func (a *Assembler) parseConst() uint16 {
	val, isUint16 := a.token.val.(int)
	if !isUint16 {
		return 0
	}
	a.nextToken()

	return uint16(val)
}

func (a *Assembler) parseImmd() uint16 {
	val, isUint16 := a.token.val.(int)
	if !isUint16 {
		return 0
	}
	a.expectToken(TokenNumber)
	if val < ImmdMin || val > ImmdMax {
		a.lexError("Immediate '%d' is out of range", val)
	}
	return uint16(val)
}

func (a *Assembler) parseAddress() uint16 {
	val, isUint16 := a.token.val.(int)
	if !isUint16 {
		return 0
	}
	a.expectToken(TokenNumber)
	if val < AddrMin || val > AddrMax {
		a.lexError("Address '%d' is out of range", val)
	}
	return uint16(val)
}

type Instruction struct {
	op, rd, rs, rt, immd, addr uint16
}

func (a *Assembler) encodeInstruction(instr Instruction) uint16 {
	rs := bits(instr.rs, 0, 3) << 9
	rt := bits(instr.rt, 0, 3) << 6
	rd := bits(instr.rd, 0, 3) << 3
	immd := bits(instr.immd, 0, 6)
	addr := bits(instr.addr, 0, 12)

	var op uint16
	var funct uint16

	switch instr.op {
	case AddOp:
		funct = 0
		return op | rs | rt | rd | funct
	case AddiOp:
		op = 0x8 << 12
		return op | rs | rt
	case AndOp:
		funct = 4
		return op | rs | rt | rd | funct
	case BeqOp:
		op = 4 << 12
		return op | rs | rt | immd
	case JOp:
		op = 2 << 12
		return op | addr
	case LwOp:
		op = 3 << 12
		return op | rs | rt | immd
	case NandOp:
		funct = 1
		return op | rs | rt | rd | funct
	case NorOp:
		funct = 7
		return op | rs | rt | rd | funct
	case OrOp:
		funct = 5
		return op | rs | rt | rd | funct
	case SltOp:
		funct = 2
		return op | rs | rt | rd | funct
	case SwOp:
		op = 0xA << 12
		return op | rs | rt | immd
	case SubOp:
		funct = 3
		return op | rs | rt | rd | funct
	}
	return 0
}

func (a *Assembler) assembleInstruction(encodedInstr uint16) {
	high := uint8(bits(encodedInstr, 8, 8))
	low := uint8(bits(encodedInstr, 0, 8))
	a.codeSection = append(a.codeSection, high)
	a.codeSection = append(a.codeSection, low)
}

func (a *Assembler) parseInstruction(sym Symbol) {
	instrDef, ok := sym.val.(InstrDef)
	if !ok {
		a.parserError("Symbol '%s' does not have value of the right type. Value should have value '%T'", sym.name, sym.val)
		return
	}

	var instr Instruction
	instr.op = instrDef.op

	switch instrDef.kind {
	case InstrRegisterKind:
		instr.rd = a.parseRegister()
		a.expectToken(TokenComma)
		instr.rs = a.parseRegister()
		a.expectToken(TokenComma)
		instr.rt = a.parseRegister()
	case InstrMemoryKind:
		instr.rt = a.parseRegister()
		a.expectToken(TokenComma)
		instr.immd = a.parseImmd()
		a.expectToken(TokenLeftParen)
		instr.rs = a.parseRegister()
		a.expectToken(TokenRightParen)
	case InstrImmKind:
		instr.rt = a.parseRegister()
		a.expectToken(TokenComma)
		instr.rs = a.parseRegister()
		a.expectToken(TokenComma)
		instr.immd = a.parseImmd()
	case InstrJumpKind:
		instr.addr = a.parseAddress()
	default:
		a.parserError("Could not parse instruction. The instruction class does not exist.")
		return
	}

	if !a.hasError {
		a.assembleInstruction(a.encodeInstruction(instr))
	}
}

func (a *Assembler) parseLine() {
	a.parseNewlines()
	if a.isToken(TokenName) {
		sym := a.parseSymbol()
		switch sym.kind {
		case SymbolInstr:
			a.parseInstruction(sym)
		}
	}
	a.parseNewlines()
}

func (a *Assembler) parseLines() {
	for !a.isToken(TokenNone) {
		a.parseLine()
	}
}

func (a *Assembler) parseFile() {
	for a.token.kind != TokenNone {
		a.parseLine()
	}
	a.pass++
}

func (a *Assembler) parserError(fmtString string, val ...any) {
	a.hasError = true
	errorMsg := fmt.Sprintf(fmtString, val...)
	fmt.Printf("Parsing error at line %d: %s\n", a.line, errorMsg)
}

func (a *Assembler) getCodeSectionStr() string {
	var buf string
	for i := 0; i < len(a.codeSection); i += 2 {
		buf += fmt.Sprintf("%08b%08b\n", a.codeSection[i], a.codeSection[i+1])
	}
	return buf
}

func (a *Assembler) getDebugStr() string {
	var buf string
	for i := 0; i < len(a.codeSection); i += 2 {
		buf += fmt.Sprintf("0x%04x -> %08b%08b\n", a.address, a.codeSection[i], a.codeSection[i+1])
	}
	return buf
}
