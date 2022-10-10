package main

import (
	"fmt"
)

const (
	SymbolNone = iota
	SymbolReg
	SymbolInstr
	SymbolCmd
)

type SymbolKind uint8

const (
	InstrRKind = iota
	InstrIKind
	InstrJKind
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
	"add":  {"add", SymbolInstr, InstrDef{InstrRKind, AddOp}},
	"addi": {"addi", SymbolInstr, InstrDef{InstrIKind, AddiOp}},
	"and":  {"and", SymbolInstr, InstrDef{InstrRKind, AndOp}},
	"beq":  {"beq", SymbolInstr, InstrDef{InstrIKind, BeqOp}},
	"j":    {"j", SymbolInstr, InstrDef{InstrJKind, JOp}},
	"lw":   {"lw", SymbolInstr, InstrDef{InstrIKind, LwOp}},
	"nand": {"nand", SymbolInstr, InstrDef{InstrRKind, NandOp}},
	"nor":  {"nor", SymbolInstr, InstrDef{InstrRKind, NorOp}},
	"or":   {"or", SymbolInstr, InstrDef{InstrRKind, OrOp}},
	"slt":  {"slw", SymbolInstr, InstrDef{InstrRKind, SltOp}},
	"sw":   {"sw", SymbolInstr, InstrDef{InstrRKind, SwOp}},
	"sub":  {"sub", SymbolInstr, InstrDef{InstrRKind, SubOp}},

	// Cmds
	".uint8":  {".uint8", SymbolCmd, cmdUint8},
	".uint16": {".uint16", SymbolCmd, cmdUint16},
	".uint32": {".uint32", SymbolCmd, cmdUint32},
}

func cmdUint8() {

}

func cmdUint16() {

}

func cmdUint32() {

}

func (a *Assembler) parseName() string {
	name := a.token.val
	a.expectToken(TokenName)
	if val, ok := name.(string); ok {
		return val
	}
	return ""
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
	a.parserError("Expected %s, got %s")
	return false
}

func (a *Assembler) parseNewlines() {
	for a.matchToken(TokenNewLine) {
	}
}

func (a *Assembler) parseRegister() uint16 {
	if a.isToken(TokenRegister) {
		if val, isUint16 := a.token.val.(uint16); isUint16 {
			a.nextToken()
			return val
		} else {
			a.parserError("Val should be uint16, not %T", val)
		}
	}
	return 0
}

func (a *Assembler) parseConst() uint16 {
	return 0
}

func (a *Assembler) parseAddress() uint16 {
	return 0
}

type Instruction struct {
	op, rd, rs, rt, immd, addr uint16
}

func (a *Assembler) assemble(instr Instruction) {
	rd := bits(instr.rd, 0, 3)
	rs := bits(instr.rs, 0, 3)
	rt := bits(instr.rt, 0, 3)
	immd := bits(instr.immd, 0, 6)
	//addr := bits(instr.immd, 0, 12)

	var binInstr uint16
	var op uint16
	var funct uint16

	switch instr.op {
	case LwOp:
		op = 3
		binInstr |= op << 12
		binInstr |= rs << 9
		binInstr |= rt << 6
		binInstr |= immd
	case AddOp:
		op = 0
		funct = 0
		binInstr |= op << 12
		binInstr |= rs << 9
		binInstr |= rt << 6
		binInstr |= rd << 3
		binInstr |= funct
		fmt.Printf("-> %016b", binInstr)
	}
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
	case InstrRKind:
		instr.rd = a.parseRegister()
		a.expectToken(TokenComma)
		instr.rs = a.parseRegister()
		a.expectToken(TokenComma)
		instr.rt = a.parseRegister()
	case InstrIKind:
		instr.rt = a.parseRegister()
		a.expectToken(TokenComma)
		instr.immd = a.parseConst()
		a.expectToken(TokenLeftParen)
		instr.rs = a.parseRegister()
	case InstrJKind:
		instr.addr = a.parseAddress()
	default:
		a.parserError("Could not parse instruction. The instruction class does not exist.")
		return
	}
	a.assemble(instr)
}

func (a *Assembler) parseLine() {
	a.parseNewlines()
	if a.isToken(TokenName) {
		sym := a.parseSymbol()
		switch sym.kind {
		case SymbolInstr:
			fmt.Printf("Instruction: %s\n", sym.name)
		}
	}
	a.parseNewlines()
}

func (a *Assembler) parseFile() {
	for a.token.kind != TokenNone {
		a.parseLine()
	}
	a.pass++
}

func (a *Assembler) parserError(fmtString string, val ...any) {
	errorMsg := fmt.Sprintf(fmtString, val...)
	fmt.Printf("Parsing error at line %d: %s\n", a.line, errorMsg)
}
