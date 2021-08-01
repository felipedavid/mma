lw $1, 0($0)
lw $2, 2($0)
add $3, $1, $2
here:
beq $1, $2, 0x03
sub $3, $1, $2
sw $3, 4($0)
j here
sw $3, 4($0)
