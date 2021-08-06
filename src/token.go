package main

type Token int

const (
    ILLEGAL Token = iota
    EOF
    COMMENT
    LABEL
    R_INSTRUCTION
    I_INSTRUCTION
    J_INSTRUCTION
)
