.data:
    .word 78
    .word 51

(a) .word 11
(b) .word 22
(c) .word 0

.text:
    // Referenciando a memória principal usando endereços
    mov $r1, [0]
    mov $r2, [2]
    add $r3, $r2, $r1
    mov [0x4], $r3

    // Referenciando valores na memória principal através de labels
    mov $r1, a
    mov $r2, b
    add $r3, $r2, $r1
    mov c, $r3

    // Imediatos para registradores
    mov $r1, 35
    mov $r2, 66
    add $r3, $r2, $r1

