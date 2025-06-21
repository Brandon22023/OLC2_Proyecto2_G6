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
    # Print salida 6
mov x0, #1
adr x1, msg6
ldr x2, =len6
ldr x2, [x2]
mov w8, #64
svc #0
    # Print salida 7
mov x0, #1
adr x1, msg7
ldr x2, =len7
ldr x2, [x2]
mov w8, #64
svc #0
    # Salida final
mov x0, #0
mov w8, #93
svc #0

.section .rodata
msg1: .ascii "hola mundo\n"
len1: .quad . - msg1
msg2: .ascii "hola a todos\n"
len2: .quad . - msg2
msg3: .ascii "esto es como un string\n"
len3: .quad . - msg3
msg4: .ascii "1\n"
len4: .quad . - msg4
msg5: .ascii "1.25\n"
len5: .quad . - msg5
msg6: .ascii "True\n"
len6: .quad . - msg6
msg7: .ascii "A\n"
len7: .quad . - msg7
