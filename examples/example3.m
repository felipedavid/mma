.data
    n1: .word 0x5
    n2: .word 0x5
.text
    lw $1, n1
    lw $2, n2
    add $3, $1, $2
    beq $1, $2, 0x3
    sub $3, $1, $2
    sw $3, 4($0)
    j here
    sw $3, 4($0)
here:
