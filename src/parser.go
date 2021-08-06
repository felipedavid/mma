package main

import (
    "fmt"
)

type Parser struct {
    scanner Scanner
    symbols map[string]int
    instructions []Instruction
    nAddr int
}

func (p *Parser) Init(src []byte) {
    p.scanner.Init(src)
    p.nAddr = 16 // address [0:15] are reserved
    p.symbols = map[string]int {
        "R0": 0, "R1": 1, "R2": 2, "R3": 3, "R4": 4, "R5": 5, "R6": 6, "R7": 7,
    }
}

func (p *Parser) Parse() AssemblyFile {
loop:
    for {
        tok, lit := p.scanner.Scan()
        switch tok {
            case EOF:
                break loop
            case LABEL:
                p.symbols[lit] = len(p.instructions)
            case ILLEGAL:
                fmt.Println(lit)
                break loop
            // implement instructions
        }
    }

    for _, instr := range p.instructions {
        switch i := instr.(type) {
        case *RInstruction:
            p.parseRInstruction(i)
        case *IInstruction:
            p.parseIInstruction(i)
        }
    }
    return AssemblyFile { Instructions: p.instructions }
}

func (p *Parser) parseRInstruction(r *RInstruction) *RInstruction {
    return r
}

func (p *Parser) parseIInstruction(i *IInstruction) *IInstruction {
    return i
}
