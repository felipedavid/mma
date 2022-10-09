package main

import "fmt"

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

type Op uint8

type InstrDef struct {
	kind InstrKind
	op   Op
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
