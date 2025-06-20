.section .text

.globl _start

_start:
    mov x0, #1
    adr x1, msg
    ldr x2, =len
    ldr x2, [x2]
    mov w8, #64
    svc #0
    mov x0, #0
    mov w8, #93
    svc #0

msg: .ascii "hola assembler\n"
len: .quad . - msg
