package main

import (
    "strings"
    "strconv"
    "fmt"
    "os"
)

type Parser struct {
    scanner Scanner
    symbols map[string]int
    instructions []Instruction
    nAddr int
}

func (p *Parser) Init(src []byte) {
    p.scanner.Init(src)
    p.nAddr = 0
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
        case R_INSTRUCTION:
            p.instructions = append(p.instructions, &RInstruction{lit: lit})
        case I_INSTRUCTION:
            p.instructions = append(p.instructions, &IInstruction{lit: lit})
        case J_INSTRUCTION:
            p.instructions = append(p.instructions, &JInstruction{lit: lit})
        }
    }

    for _, instr := range p.instructions {
        switch i := instr.(type) {
        case *RInstruction:
            p.parseRInstruction(i)
        case *IInstruction:
            p.parseIInstruction(i)
        case *JInstruction:
            p.parseJInstruction(i)
        }
    }
    return AssemblyFile{ Instructions: p.instructions }
}

func (p *Parser) parseRInstruction(r *RInstruction) {
    r.rs = 1
    r.rt = 2
    r.rd = 3
    r.funct = 5
}

func (p *Parser) parseIInstruction(i *IInstruction) {
    i.rs = 3
    i.rt = 0
    i.immd = 4
}

func (p *Parser) parseJInstruction(j *JInstruction) {
    label := strings.Fields(j.lit)[0]
    if val, ok := p.symbols[label]; ok { // if the label is present in the map
        j.addr = val
    } else {
        addr, err := strconv.Atoi(label)
        if err != nil {
            fmt.Printf("Label \"%v\" dosen't exist.\n", label)
            os.Exit(1)
        }
        j.addr = addr
    }
}
