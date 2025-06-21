.data
heap: .space 4096
heap_ptr: .quad heap
buffer: .space 32

msg_nl: .ascii "\n"
len_nl: .quad . - msg_nl

msg1: .ascii "hola mundo\n"
len1: .quad . - msg1

msg2: .ascii "hola causa\n"
len2: .quad . - msg2

msg3: .ascii "-----------------------------\n"
len3: .quad . - msg3

.text
.global malloc
malloc:
    mov x2, x10
    add x0, x2, x0
    mov x10, x0
    mov x0, x2
    ret

.global _start
_start:
    adr x10, heap

    // imprimir msg1
    mov x0, #1
    ldr x1, =msg1
    ldr x2, =len1
    ldr x2, [x2]
    mov w8, #64
    svc #0

    // imprimir msg2
    mov x0, #1
    ldr x1, =msg2
    ldr x2, =len2
    ldr x2, [x2]
    mov w8, #64
    svc #0

    // Número: 1
    mov x0, #1
    ldr x1, =buffer
    bl int_to_ascii
    mov x2, x0      // longitud
    mov x0, #1      // stdout
    // x1 ya es puntero correcto al número
    mov w8, #64
    svc #0

    // Salto de línea
    mov x0, #1
    ldr x1, =msg_nl
    ldr x2, =len_nl
    ldr x2, [x2]
    mov w8, #64
    svc #0

    // línea divisoria
    mov x0, #1
    ldr x1, =msg3
    ldr x2, =len3
    ldr x2, [x2]
    mov w8, #64
    svc #0

    // Número: 25
    mov x0, #25
    ldr x1, =buffer
    bl int_to_ascii
    mov x2, x0
    mov x0, #1
    mov w8, #64
    svc #0

    // Salto de línea
    mov x0, #1
    ldr x1, =msg_nl
    ldr x2, =len_nl
    ldr x2, [x2]
    mov w8, #64
    svc #0

    // Número: 2025
    mov x0, #2025
    ldr x1, =buffer
    bl int_to_ascii
    mov x2, x0
    mov x0, #1
    mov w8, #64
    svc #0

    // Salto de línea
    mov x0, #1
    ldr x1, =msg_nl
    ldr x2, =len_nl
    ldr x2, [x2]
    mov w8, #64
    svc #0

    // salir
    mov x0, #0
    mov w8, #93
    svc #0

// ============================================
// x0 = número
// x1 = puntero a buffer
// devuelve:
//   x0 = longitud
//   x1 = puntero al inicio del número convertido
int_to_ascii:
    mov x2, x1      // puntero de escritura
    mov x3, #10     // divisor base 10
    mov x4, #0      // longitud
    mov x5, x0      // copia del número

    cmp x5, #0
    bne convert_loop
    mov w6, #'0'
    strb w6, [x2], #1
    mov x4, #1
    b reverse

convert_loop:
    udiv x6, x5, x3
    msub x7, x6, x3, x5 // x7 = resto
    add w7, w7, #'0'
    strb w7, [x2], #1
    mov x5, x6
    add x4, x4, #1
    cmp x5, #0
    bne convert_loop

reverse:
    sub x2, x2, x4     // x2 apunta al inicio del número al revés
    mov x7, x4
    mov x8, x1         // destino original (buffer)
rev_loop:
    cmp x7, #0
    beq done
    ldrb w9, [x2], #1
    strb w9, [x8], #1
    sub x7, x7, #1
    b rev_loop

done:
    sub x1, x8, x4     // x1 = puntero al inicio del número ya invertido
    mov x0, x4         // longitud
    ret
