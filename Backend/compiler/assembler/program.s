.data
heap: .space 4096
heap_ptr: .quad heap
buffer: .space 32
msg_nl: .ascii "\n"
len_nl: .quad . - msg_nl
msg_true: .ascii "true"
len_true: .quad . - msg_true
msg_falsestr: .ascii "false"
len_falsestr: .quad . - msg_falsestr
.align 3
float_100: .double 100.0
msg1: .ascii "¡Hola, mundo!"
len1: .quad . - msg1
.align 3
puntos: .quad 0
msg2: .ascii "=== Prueba de funciones simplificada ==="
len2: .quad . - msg2
.align 3
puntos_simples: .quad 0
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
fn_saludar:
mov x0, #1
ldr x1, =msg1
ldr x2, =len1
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
ret
mov x0, #1
ldr x1, =msg2
ldr x2, =len2
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
    # Salida final
mov x0, #0
mov w8, #93
svc #0
# Finalizar programa


# Foreign functions


# libreria estandar
// Funciones estándar aquí