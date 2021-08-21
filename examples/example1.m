.data:
    .word 0x5
    .word 0xA

.text:
    lw $2, 0($1)
    lw $3, 2($1)
    add $r4, $r2, $r3
    sw $4, 4($1)
