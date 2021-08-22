.data:
    (user_input) .word 23
.text:
    mov $r1, user_input 
    mov $r2, 1
(L1)
    beq $r1, $r0, 0x4
    sub $r1, $r1, $r2
    j L1
