.data
    a: .word 5
    b: .word 10    
    c: .word 0
.text
    lw $2, a
    lw $3, b
    add $4, $2, $3
    sw $4, c
