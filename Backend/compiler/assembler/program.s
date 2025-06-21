.section .text

.globl _start

_start:
    # Print salida 1
mov x0, #1
adr x1, msg1
ldr x2, =len1
ldr x2, [x2]
mov w8, #64
svc #0
    # Print salida 2
mov x0, #1
adr x1, msg2
ldr x2, =len2
ldr x2, [x2]
mov w8, #64
svc #0
    # Print salida 3
mov x0, #1
adr x1, msg3
ldr x2, =len3
ldr x2, [x2]
mov w8, #64
svc #0
    # Print salida 4
mov x0, #1
adr x1, msg4
ldr x2, =len4
ldr x2, [x2]
mov w8, #64
svc #0
    # Print salida 5
mov x0, #1
adr x1, msg5
ldr x2, =len5
ldr x2, [x2]
mov w8, #64
svc #0
    # Salida final
mov x0, #0
mov w8, #93
svc #0

.section .rodata
msg1: .ascii "1\n"
len1: .quad . - msg1
msg2: .ascii "1.25\n"
len2: .quad . - msg2
msg3: .ascii "true\n"
len3: .quad . - msg3
msg4: .ascii "hola causa\n"
len4: .quad . - msg4
msg5: .ascii "H\n"
len5: .quad . - msg5
