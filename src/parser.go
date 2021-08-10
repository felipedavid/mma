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
    no_commas := strings.ReplaceAll(r.lit, ",", " ")
    lits := strings.Fields(no_commas)

    // Parsing funct
    switch lits[0] {
        case "add":
            r.funct = 0
        case "and":
            r.funct = 4
        case "nand":
            r.funct = 1
        case "nor":
            r.funct = 7
        case "or":
            r.funct = 5
        case "slt":
            r.funct = 2
        case "sub":
            r.funct = 3
        default:
            fmt.Println("Invalid instruction:", lits[0])
            os.Exit(1)
    }

    // Parsing registers
    var reg_i int
    for n, reg := range lits[1:] {
        if reg[0] == '$' {
            if val, ok := p.symbols[reg[1:]]; ok {
                reg_i = val
            } else {
                reg_n, err := strconv.Atoi(reg[1:])
                if err != nil {
                    fmt.Println("Invalid register:", reg)
                    os.Exit(1)
                }
                reg_i = reg_n
            }

            switch n {
            case 0:
                r.rd = reg_i
            case 1:
                r.rs = reg_i
            case 2:
                r.rt = reg_i
            }
        }
    }
}

func (p *Parser) parseIInstruction(i *IInstruction) {
    clean_lit := strings.ReplaceAll(i.lit, ",", " ")
    clean_lit = strings.ReplaceAll(clean_lit, "(", " ")
    clean_lit = strings.ReplaceAll(clean_lit, ")", " ")
    lits := strings.Fields(clean_lit)
    fmt.Println(lits)

    switch lits[0] {
        case "addi":
            i.op = 8
        case "beq":
            i.op = 4
        case "lw":
            i.op = 3
        case "sw":
            i.op = 10
        default:
            fmt.Println("Invalid instruction:", lits[0])
            os.Exit(1)
    }

    if i.op == 3 || i.op == 10 { // if "lw" or "sw"
        var reg_i int
        for n, field := range lits[1:4] {
            if field[0] == '$' {
                if val, ok := p.symbols[field[1:]]; ok {
                    reg_i = val
                } else {
                    reg_n, err := strconv.Atoi(field[1:])
                    if err != nil {
                        fmt.Println("Invalid register:", field)
                        os.Exit(1)
                    }
                    reg_i = reg_n
                }

                switch n {
                case 0:
                    i.rt = reg_i
                case 2:
                    i.rs = reg_i
                }
            } else if n == 1 {
                immd, err := strconv.Atoi(field)
                if err != nil {
                    fmt.Println("Invalid immediate:", immd)
                    os.Exit(1)
                }
                i.immd = immd
            }
        }
    } else if i.op == 8 || i.op == 4 { // if "addi" or "beq"
        var reg_i int
        for n, field := range lits[1:4] {
            if field[0] == '$' {
                if val, ok := p.symbols[field[1:]]; ok {
                    reg_i = val
                } else {
                    reg_n, err := strconv.Atoi(field[1:])
                    if err != nil {
                        fmt.Println("Invalid register:", field)
                        os.Exit(1)
                    }
                    reg_i = reg_n
                }

                switch n {
                case 0:
                    i.rs = reg_i
                case 1:
                    i.rt = reg_i
                }
            } else if n == 2 {
                immd, err := strconv.Atoi(field[:])
                if err != nil {
                    fmt.Println("Invalid immediate:", field[:])
                    fmt.Println("Only valid immediates are base 10 integers.:")
                    os.Exit(1)
                }
                i.immd = immd
            }
        }
    } else {
        fmt.Println("Invalid I Instruction:", lits[0])
        os.Exit(1)
    }
}

func (p *Parser) parseJInstruction(j *JInstruction) {
    fmt.Println(p.symbols)
    fmt.Println(j.lit)
    label := strings.Fields(j.lit)[1]
    if val, ok := p.symbols[label]; ok { // if the label is present in the map
        j.addr = val * 2
    } else {
        addr, err := strconv.Atoi(label)
        if err != nil {
            fmt.Printf("Label \"%v\" dosen't exist.\n", label)
            os.Exit(1)
        }
        j.addr = addr
    }
}
