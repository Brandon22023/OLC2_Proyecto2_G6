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
.align 3
numero: .quad 0
msg3: .ascii "Número obtenido:"
len3: .quad . - msg3
msg4: .ascii "OK obtener_numero"
len4: .quad . - msg4
.text
.global malloc
malloc:
    mov x2, x10
    add x0, x2, x0
    mov x10, x0
    mov x0, x2
    ret
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
fn_end_1:
ret
fn_obtener_numero:
fn_end_2:
ret
.global _start
_start:
    adr x10, heap
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
bl fn_saludar
bl fn_obtener_numero
mov x0, #1
ldr x1, =msg3
ldr x2, =len3
ldr x2, [x2]
mov w8, #64
svc #0
ldr x1, =numero
ldr x0, [x1]
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
ldr x1, =buffer
mov x0, #1
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
ldr x1, =numero
ldr x0, [x1]
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
ldr x1, =buffer
mov x0, #1
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
ldr x1, =numero
ldr x0, [x1]
mov x0, x0
mov x1, 42
cmp x0, x1
mov x2, #0
beq cmp_eq_true_1
b cmp_eq_end_1
cmp_eq_true_1:
mov x2, #1
cmp_eq_end_1:
cmp x2, #0
beq ifend_1
mov x0, #1
ldr x1, =msg4
ldr x2, =len4
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
ldr x1, =puntos_simples
ldr x0, [x1]
mov x2, #1
add x0, x0, x2
str x0, [x1]
B ifend_1
ifend_1:
    mov x0, #0
    mov w8, #93
    svc #0

int_to_ascii:
    mov x2, x1          // x2 = buffer (puntero de escritura)
    mov x3, #10         // divisor
    mov x4, #0          // contador de dígitos
    mov x5, x0          // copia del número
    cmp x5, #0
    bne int_to_ascii_loop
    mov w6, #'0'
    strb w6, [x2], #1
    mov x4, #1
    b int_to_ascii_done
int_to_ascii_loop:
    udiv x6, x5, x3
    msub x7, x6, x3, x5 // x7 = x5 - x6*x3 (resto)
    add w7, w7, #'0'
    strb w7, [x2], #1
    mov x5, x6
    add x4, x4, #1
    cmp x5, #0
    bne int_to_ascii_loop
int_to_ascii_done:
    sub x2, x2, x4      // x2 = inicio de los dígitos
    mov x5, x4          // x5 = longitud
    mov x6, #0          // i = 0
reverse_inplace_loop:
    cmp x6, x5
    bge reverse_inplace_done
    add x7, x2, x6      // &buffer[i]
    sub x8, x2, x6
    add x8, x8, x5
    sub x8, x8, #1      // &buffer[len-1-i]
    ldrb w9, [x7]
    ldrb w10, [x8]
    strb w10, [x7]
    strb w9, [x8]
    add x6, x6, #1
    cmp x6, x5, lsr #1
    blt reverse_inplace_loop
reverse_inplace_done:
    mov x1, x2          // x1 = puntero al inicio
    mov x0, x5          // x0 = longitud
    ret

# Finalizar programa


# Foreign functions


# libreria estandar
// Funciones estándar aquí