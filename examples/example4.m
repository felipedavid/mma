.text
    lw $2, 0($1)
start:
    lw $3, 2($1)
    j start
    add $4, $2, $3
    sw $4, 4($1)
