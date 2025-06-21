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
    # Print salida 8
mov x0, #1
adr x1, msg8
ldr x2, =len8
ldr x2, [x2]
mov w8, #64
svc #0
    # Print salida 9
mov x0, #1
adr x1, msg9
ldr x2, =len9
ldr x2, [x2]
mov w8, #64
svc #0
    # Print salida 10
mov x0, #1
adr x1, msg10
ldr x2, =len10
ldr x2, [x2]
mov w8, #64
svc #0
    # Print salida 11
mov x0, #1
adr x1, msg11
ldr x2, =len11
ldr x2, [x2]
mov w8, #64
svc #0
    # Print salida 12
mov x0, #1
adr x1, msg12
ldr x2, =len12
ldr x2, [x2]
mov w8, #64
svc #0
    # Print salida 13
mov x0, #1
adr x1, msg13
ldr x2, =len13
ldr x2, [x2]
mov w8, #64
svc #0
    # Print salida 14
mov x0, #1
adr x1, msg14
ldr x2, =len14
ldr x2, [x2]
mov w8, #64
svc #0
    # Print salida 15
mov x0, #1
adr x1, msg15
ldr x2, =len15
ldr x2, [x2]
mov w8, #64
svc #0
    # Print salida 16
mov x0, #1
adr x1, msg16
ldr x2, =len16
ldr x2, [x2]
mov w8, #64
svc #0
    # Print salida 17
mov x0, #1
adr x1, msg17
ldr x2, =len17
ldr x2, [x2]
mov w8, #64
svc #0
    # Print salida 18
mov x0, #1
adr x1, msg18
ldr x2, =len18
ldr x2, [x2]
mov w8, #64
svc #0
    # Print salida 19
mov x0, #1
adr x1, msg19
ldr x2, =len19
ldr x2, [x2]
mov w8, #64
svc #0
    # Salida final
mov x0, #0
mov w8, #93
svc #0

    # Salida final
mov x0, #0
mov w8, #93
svc #0

.section .rodata
msg1: .ascii "12\n"
len1: .quad . - msg1
msg2: .ascii "Hola arm\n"
len2: .quad . - msg2
msg3: .ascii "1.25\n"
len3: .quad . - msg3
msg4: .ascii "true\n"
len4: .quad . - msg4
msg5: .ascii "H\n"
len5: .quad . - msg5
msg_puy: .ascii "100\n"
len_puy: .quad . - msg_puy
msg6: .ascii "valor de puy:  100\n"
len6: .quad . - msg6
msg_i: .ascii "15\n"
len_i: .quad . - msg_i
msg_j: .ascii "como estas\n"
len_j: .quad . - msg_j
msg_k: .ascii "B\n"
len_k: .quad . - msg_k
msg7: .ascii "15 como estas B\n"
len7: .quad . - msg7
msg_defect: .ascii "0\n"
len_defect: .quad . - msg_defect
msg_defect1: .ascii "false\n"
len_defect1: .quad . - msg_defect1
msg8: .ascii "valor por defecto de un int:  0\n"
len8: .quad . - msg8
msg9: .ascii "valor por defecto de un bool:  false\n"
len9: .quad . - msg9
msg_var: .ascii "30\n"
len_var: .quad . - msg_var
msg10: .ascii "valor de var:  30\n"
len10: .quad . - msg10
msg11: .ascii "GANO GUATE VAMOSS 10.5\n"
len11: .quad . - msg11
msg12: .ascii "1\n"
len12: .quad . - msg12
msg13: .ascii "0.333333\n"
len13: .quad . - msg13
msg14: .ascii "13\n"
len14: .quad . - msg14
msg15: .ascii "concatenacion\n"
len15: .quad . - msg15
msg16: .ascii "ahora el valor de I es:  15\n"
len16: .quad . - msg16
msg17: .ascii "ahora el valor de I es:  15\n"
len17: .quad . - msg17
msg18: .ascii "valor del booleno cambiado porque guate gano GUATE ES CLAVE:  false\n"
len18: .quad . - msg18
msg19: .ascii "valor del booleno cambiado porque guate gano GUATE ES CLAVE:  0\n"
len19: .quad . - msg19
