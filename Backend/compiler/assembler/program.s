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
.align 3
puntos: .quad 0
msg_print_1: .ascii "=== Prueba Básica Simplificada ==="
len_print_1: .quad . - msg_print_1
.align 3
entero: .quad 42
.align 3
decimal: .quad 0x40091eb851eb851f
texto: .asciz "Hola, mundo!"
len_texto: .quad . - texto
.align 3
booleano: .quad 1
msg_print_2: .ascii "OK entero"
len_print_2: .quad . - msg_print_2
.align 3
cmp_float_left_3: .quad 0x40091eb851eb851f
.align 3
cmp_float_right_4: .quad 0x4008000000000000
msg_print_5: .ascii "OK decimal"
len_print_5: .quad . - msg_print_5
msg_print_6: .ascii "OK texto"
len_print_6: .quad . - msg_print_6
msg_print_7: .ascii "OK booleano"
len_print_7: .quad . - msg_print_7
msg_print_8: .ascii "OK asignación entero"
len_print_8: .quad . - msg_print_8
.align 3
suma: .quad 15
msg_print_9: .ascii "OK suma"
len_print_9: .quad . - msg_print_9
.align 3
resta: .quad 5
msg_print_10: .ascii "OK resta"
len_print_10: .quad . - msg_print_10
msg_print_11: .ascii "OK igualdad"
len_print_11: .quad . - msg_print_11
msg_print_12: .ascii "OK lógica AND"
len_print_12: .quad . - msg_print_12
msg_print_13: .ascii "OK lógica OR negada"
len_print_13: .quad . - msg_print_13
msg_print_14: .ascii "Texto de prueba"
len_print_14: .quad . - msg_print_14
msg_print_15: .ascii "Puntos obtenidos: "
len_print_15: .quad . - msg_print_15
msg_print_16: .ascii " / 10"
len_print_16: .quad . - msg_print_16
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
mov x0, #1
ldr x1, =msg_print_1
ldr x2, =len_print_1
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
ldr x1, =entero
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
ldr x1, =msg_print_2
ldr x2, =len_print_2
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x2, #1
ldr x1, =puntos
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
B ifend_1
ifend_1:
ldr x0, =cmp_float_left_3
ldr d0, [x0]
ldr x1, =cmp_float_right_4
ldr d1, [x1]
fcmp d0, d1
mov x2, #0
bgt cmp_true_1
b cmp_end_1
cmp_true_1:
mov x2, #1
cmp_end_1:
cmp x2, #0
beq ifend_2
mov x0, #1
ldr x1, =msg_print_5
ldr x2, =len_print_5
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x2, #1
ldr x1, =puntos
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
B ifend_2
ifend_2:
ldr x1, =texto
ldr x0, [x1]
cmp x2, #0
beq ifend_3
mov x0, #1
ldr x1, =msg_print_6
ldr x2, =len_print_6
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x2, #1
ldr x1, =puntos
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
B ifend_3
ifend_3:
ldr x1, =booleano
ldr x0, [x1]
mov x0, x0
cmp x0, #0
beq ifend_4
mov x0, #1
ldr x1, =msg_print_7
ldr x2, =len_print_7
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x2, #1
ldr x1, =puntos
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
B ifend_4
ifend_4:
ldr x1, =entero
mov x0, #99
str x0, [x1]
ldr x1, =entero
ldr x0, [x1]
mov x0, x0
mov x1, 99
cmp x0, x1
mov x2, #0
beq cmp_eq_true_2
b cmp_eq_end_2
cmp_eq_true_2:
mov x2, #1
cmp_eq_end_2:
cmp x2, #0
beq ifend_5
mov x0, #1
ldr x1, =msg_print_8
ldr x2, =len_print_8
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x2, #1
ldr x1, =puntos
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
B ifend_5
ifend_5:
ldr x1, =suma
ldr x0, [x1]
mov x0, x0
mov x1, 15
cmp x0, x1
mov x2, #0
beq cmp_eq_true_3
b cmp_eq_end_3
cmp_eq_true_3:
mov x2, #1
cmp_eq_end_3:
cmp x2, #0
beq ifend_6
mov x0, #1
ldr x1, =msg_print_9
ldr x2, =len_print_9
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x2, #1
ldr x1, =puntos
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
B ifend_6
ifend_6:
ldr x1, =resta
ldr x0, [x1]
mov x0, x0
mov x1, 5
cmp x0, x1
mov x2, #0
beq cmp_eq_true_4
b cmp_eq_end_4
cmp_eq_true_4:
mov x2, #1
cmp_eq_end_4:
cmp x2, #0
beq ifend_7
mov x0, #1
ldr x1, =msg_print_10
ldr x2, =len_print_10
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x2, #1
ldr x1, =puntos
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
B ifend_7
ifend_7:
mov x0, 10
mov x1, 10
cmp x0, x1
mov x2, #0
beq cmp_eq_true_5
b cmp_eq_end_5
cmp_eq_true_5:
mov x2, #1
cmp_eq_end_5:
cmp x2, #0
beq ifend_8
mov x0, #1
ldr x1, =msg_print_11
ldr x2, =len_print_11
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x2, #1
ldr x1, =puntos
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
B ifend_8
ifend_8:
mov x0, #1
mov x1, #1
mov x2, #0
cmp x0, #0
beq logic_end_1
cmp x1, #0
beq logic_end_1
mov x2, #1
b logic_end_1
logic_end_1:
cmp x2, #0
beq ifend_9
mov x0, #1
ldr x1, =msg_print_12
ldr x2, =len_print_12
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x2, #1
ldr x1, =puntos
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
B ifend_9
ifend_9:
mov x0, #0
mov x1, #0
mov x2, #1
cmp x0, #0
bne logic_true_2
cmp x1, #0
bne logic_true_2
mov x2, #0
b logic_end_2
logic_true_2:
mov x2, #1
logic_end_2:
mov x0, x2
mov x2, #0
cmp x0, #0
beq not_true_1
b not_end_1
not_true_1:
mov x2, #1
not_end_1:
cmp x2, #0
beq ifend_10
mov x0, #1
ldr x1, =msg_print_13
ldr x2, =len_print_13
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x2, #1
ldr x1, =puntos
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
B ifend_10
ifend_10:
mov x0, #42
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
mov x0, #1
ldr x1, =msg_print_14
ldr x2, =len_print_14
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, 1
bl bool_to_ascii
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
mov x0, #1
ldr x1, =msg_print_15
ldr x2, =len_print_15
ldr x2, [x2]
mov w8, #64
svc #0
ldr x1, =puntos
ldr x0, [x1]
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
ldr x1, =buffer
mov x0, #1
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_print_16
ldr x2, =len_print_16
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
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


bool_to_ascii:
    cmp x0, #1
    beq bool_true
    adr x2, msg_falsestr
    ldr x0, =len_falsestr
    ldr x0, [x0]
    mov x1, x2
    ret
bool_true:
    adr x2, msg_true
    ldr x0, =len_true
    ldr x0, [x0]
    mov x1, x2
    ret

# Finalizar programa


# Foreign functions


# libreria estandar
// Funciones estándar aquí