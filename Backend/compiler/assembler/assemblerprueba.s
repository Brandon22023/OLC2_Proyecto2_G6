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
.align 3
puntos_entornos: .quad 0
.align 3
a: .quad 10
msg_print_1: .ascii "a ="
len_print_1: .quad . - msg_print_1
msg_print_2: .ascii "OK a = 10"
len_print_2: .quad . - msg_print_2
.align 3
b: .quad 10
msg_print_3: .ascii "b ="
len_print_3: .quad . - msg_print_3
msg_print_4: .ascii "OK b = 20"
len_print_4: .quad . - msg_print_4
.align 3
c: .quad 10
.align 3
d: .quad 10
msg_print_5: .ascii "c ="
len_print_5: .quad . - msg_print_5
msg_print_6: .ascii "d ="
len_print_6: .quad . - msg_print_6
msg_print_7: .ascii "OK c = 30"
len_print_7: .quad . - msg_print_7
.align 3
puntos_if: .quad 0
msg_print_8: .ascii "OK true"
len_print_8: .quad . - msg_print_8
msg_print_9: .ascii "OK 1 == 1"
len_print_9: .quad . - msg_print_9
msg_print_10: .ascii "OK 2 > 1"
len_print_10: .quad . - msg_print_10
.align 3
puntos_while: .quad 0
.align 3
i: .quad 0
.align 3
suma1: .quad 0
msg_print_11: .ascii "OK suma1 == 10"
len_print_11: .quad . - msg_print_11
msg_print_12: .ascii "OK i == 5"
len_print_12: .quad . - msg_print_12
.align 3
j: .quad 3
.align 3
k: .quad 0
.align 3
puntos_for: .quad 0
.align 3
suma2: .quad 0
.align 3
x: .quad 0
msg_print_13: .ascii "OK suma2 == 10"
len_print_13: .quad . - msg_print_13
.align 3
y: .quad 0
.align 3
z: .quad 0
.align 3
puntos_case: .quad 0
.align 3
dia: .quad 1
msg_print_14: .ascii "Lunes"
len_print_14: .quad . - msg_print_14
msg_print_15: .ascii "Martes"
len_print_15: .quad . - msg_print_15
msg_print_16: .ascii "Miércoles"
len_print_16: .quad . - msg_print_16
msg_print_17: .ascii "Jueves"
len_print_17: .quad . - msg_print_17
msg_print_18: .ascii "Viernes"
len_print_18: .quad . - msg_print_18
msg_print_19: .ascii "Sábado"
len_print_19: .quad . - msg_print_19
msg_print_20: .ascii "Domingo"
len_print_20: .quad . - msg_print_20
msg_print_21: .ascii "Día inválido"
len_print_21: .quad . - msg_print_21
msg_print_22: .ascii "Puntos totales:"
len_print_22: .quad . - msg_print_22
msg_print_23: .ascii "/ 26"
len_print_23: .quad . - msg_print_23
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
ldr x1, =a
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
ldr x1, =a
ldr x0, [x1]
mov x0, x0
mov x1, 10
cmp x0, x1
mov x2, #0
beq cmp_eq_true_1
b cmp_eq_end_1
cmp_eq_true_1:
mov x2, #1
cmp_eq_end_1:
cmp x2, #0
beq ifend_1
mov x2, #1
ldr x1, =puntos_entornos
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
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
B ifend_1
ifend_1:
ldr x1, =b
mov x0, #20
str x0, [x1]
mov x0, #1
ldr x1, =msg_print_3
ldr x2, =len_print_3
ldr x2, [x2]
mov w8, #64
svc #0
ldr x1, =b
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
ldr x1, =b
ldr x0, [x1]
mov x0, x0
mov x1, 20
cmp x0, x1
mov x2, #0
beq cmp_eq_true_2
b cmp_eq_end_2
cmp_eq_true_2:
mov x2, #1
cmp_eq_end_2:
cmp x2, #0
beq ifend_2
mov x2, #1
ldr x1, =puntos_entornos
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
mov x0, #1
ldr x1, =msg_print_4
ldr x2, =len_print_4
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
B ifend_2
ifend_2:
ldr x1, =c
mov x0, #30
str x0, [x1]
mov x0, #1
ldr x1, =msg_print_5
ldr x2, =len_print_5
ldr x2, [x2]
mov w8, #64
svc #0
ldr x1, =c
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
mov x0, #1
ldr x1, =msg_print_6
ldr x2, =len_print_6
ldr x2, [x2]
mov w8, #64
svc #0
ldr x1, =d
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
ldr x1, =c
ldr x0, [x1]
mov x0, x0
mov x1, 30
cmp x0, x1
mov x2, #0
beq cmp_eq_true_3
b cmp_eq_end_3
cmp_eq_true_3:
mov x2, #1
cmp_eq_end_3:
cmp x2, #0
beq ifend_3
mov x2, #1
ldr x1, =puntos_entornos
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
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
B ifend_3
ifend_3:
mov x2, #1
ldr x1, =puntos_if
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
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
B ifend_4
ifend_4:
mov x0, 1
mov x1, 1
cmp x0, x1
mov x2, #0
beq cmp_eq_true_4
b cmp_eq_end_4
cmp_eq_true_4:
mov x2, #1
cmp_eq_end_4:
cmp x2, #0
beq ifend_5
mov x2, #1
ldr x1, =puntos_if
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
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
B ifend_5
ifend_5:
mov x0, 2
mov x1, 1
cmp x0, x1
mov x2, #0
bgt cmp_true_1
b cmp_end_1
cmp_true_1:
mov x2, #1
cmp_end_1:
cmp x2, #0
beq ifend_6
mov x2, #1
ldr x1, =puntos_if
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
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
B ifend_6
ifend_6:
B for_cond_1
for_start_1:
ldr x1, =i
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
ldr x1, =i
ldr x0, [x1]
mov x2, x0
ldr x1, =suma1
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
mov x2, #1
ldr x1, =i
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
for_continue_1:
for_cond_1:
ldr x1, =i
ldr x0, [x1]
mov x0, x0
mov x1, 5
cmp x0, x1
mov x2, #0
blt cmp_true_2
b cmp_end_2
cmp_true_2:
mov x2, #1
cmp_end_2:
cmp x2, #0
beq for_end_1
B for_start_1
for_end_1:
ldr x1, =suma1
ldr x0, [x1]
mov x0, x0
mov x1, 10
cmp x0, x1
mov x2, #0
beq cmp_eq_true_5
b cmp_eq_end_5
cmp_eq_true_5:
mov x2, #1
cmp_eq_end_5:
cmp x2, #0
beq ifend_7
mov x2, #1
ldr x1, =puntos_while
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
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
B ifend_7
ifend_7:
ldr x1, =i
ldr x0, [x1]
mov x0, x0
mov x1, 5
cmp x0, x1
mov x2, #0
beq cmp_eq_true_6
b cmp_eq_end_6
cmp_eq_true_6:
mov x2, #1
cmp_eq_end_6:
cmp x2, #0
beq ifend_8
mov x2, #1
ldr x1, =puntos_while
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
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
B ifend_8
ifend_8:
B for_cond_2
for_start_2:
ldr x1, =j
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
mov x2, #1
ldr x1, =j
ldr x0, [x1]
sub x0, x0, x2
str x0, [x1]
for_continue_2:
for_cond_2:
ldr x1, =j
ldr x0, [x1]
mov x0, x0
mov x1, 0
cmp x0, x1
mov x2, #0
bgt cmp_true_3
b cmp_end_3
cmp_true_3:
mov x2, #1
cmp_end_3:
cmp x2, #0
beq for_end_2
B for_start_2
for_end_2:
B for_cond_3
for_start_3:
ldr x1, =k
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
mov x2, #2
ldr x1, =k
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
for_continue_3:
for_cond_3:
ldr x1, =k
ldr x0, [x1]
mov x0, x0
mov x1, 10
cmp x0, x1
mov x2, #0
ble cmp_true_4
b cmp_end_4
cmp_true_4:
mov x2, #1
cmp_end_4:
cmp x2, #0
beq for_end_3
B for_start_3
for_end_3:
ldr x1, =x
mov x0, 0
str x0, [x1]
B for_cond_4
for_start_4:
ldr x1, =x
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
ldr x1, =x
ldr x0, [x1]
mov x2, x0
ldr x1, =suma2
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
for_continue_4:
ldr x1, =x
ldr x0, [x1]
add x0, x0, #1
str x0, [x1]
for_cond_4:
ldr x1, =x
ldr x0, [x1]
mov x0, x0
mov x1, 5
cmp x0, x1
mov x2, #0
blt cmp_true_5
b cmp_end_5
cmp_true_5:
mov x2, #1
cmp_end_5:
cmp x2, #0
beq for_end_4
B for_start_4
for_end_4:
ldr x1, =suma2
ldr x0, [x1]
mov x0, x0
mov x1, 10
cmp x0, x1
mov x2, #0
beq cmp_eq_true_7
b cmp_eq_end_7
cmp_eq_true_7:
mov x2, #1
cmp_eq_end_7:
cmp x2, #0
beq ifend_9
mov x2, #1
ldr x1, =puntos_for
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
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
B ifend_9
ifend_9:
ldr x1, =y
mov x0, 0
str x0, [x1]
B for_cond_5
for_start_5:
ldr x1, =y
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
for_continue_5:
ldr x1, =y
ldr x0, [x1]
add x0, x0, #1
str x0, [x1]
for_cond_5:
ldr x1, =y
ldr x0, [x1]
mov x0, x0
mov x1, 3
cmp x0, x1
mov x2, #0
blt cmp_true_6
b cmp_end_6
cmp_true_6:
mov x2, #1
cmp_end_6:
cmp x2, #0
beq for_end_5
B for_start_5
for_end_5:
ldr x1, =z
mov x0, 0
str x0, [x1]
B for_cond_6
for_start_6:
ldr x1, =z
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
for_continue_6:
ldr x1, =z
ldr x0, [x1]
add x0, x0, #1
str x0, [x1]
for_cond_6:
ldr x1, =z
ldr x0, [x1]
mov x0, x0
mov x1, 2
cmp x0, x1
mov x2, #0
blt cmp_true_7
b cmp_end_7
cmp_true_7:
mov x2, #1
cmp_end_7:
cmp x2, #0
beq for_end_6
B for_start_6
for_end_6:
mov x2, #2
ldr x1, =puntos_for
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
ldr x1, =dia
ldr x0, [x1]
mov x0, x0
mov x1, 1
cmp x0, x1
beq case_0_1
mov x1, 2
cmp x0, x1
beq case_1_1
mov x1, 3
cmp x0, x1
beq case_2_1
mov x1, 4
cmp x0, x1
beq case_3_1
mov x1, 5
cmp x0, x1
beq case_4_1
mov x1, 6
cmp x0, x1
beq case_5_1
mov x1, 7
cmp x0, x1
beq case_6_1
B switch_default_1
case_0_1:
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
mov x2, #1
ldr x1, =puntos_case
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
B switch_end_1
case_1_1:
mov x0, #1
ldr x1, =msg_print_15
ldr x2, =len_print_15
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
B switch_end_1
case_2_1:
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
mov x2, #1
ldr x1, =puntos_case
ldr x0, [x1]
add x0, x0, x2
str x0, [x1]
B switch_end_1
case_3_1:
mov x0, #1
ldr x1, =msg_print_17
ldr x2, =len_print_17
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
B switch_end_1
case_4_1:
mov x0, #1
ldr x1, =msg_print_18
ldr x2, =len_print_18
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
B switch_end_1
case_5_1:
mov x0, #1
ldr x1, =msg_print_19
ldr x2, =len_print_19
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
B switch_end_1
case_6_1:
mov x0, #1
ldr x1, =msg_print_20
ldr x2, =len_print_20
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
B switch_end_1
switch_default_1:
mov x0, #1
ldr x1, =msg_print_21
ldr x2, =len_print_21
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
switch_end_1:
ldr x1, =puntos_entornos
ldr x0, [x1]
ldr x1, =puntos_if
ldr x0, [x1]
ldr x1, =puntos_while
ldr x0, [x1]
ldr x1, =puntos_for
ldr x0, [x1]
ldr x1, =puntos_case
ldr x0, [x1]
mov x0, #1
ldr x1, =msg_print_22
ldr x2, =len_print_22
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
ldr x1, =msg_print_23
ldr x2, =len_print_23
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

strcmp:
    // Guardar los punteros en registros temporales
    mov x9, x0   // x9 = ptr1
    mov x10, x1  // x10 = ptr2

.loop:
    ldrb w2, [x9], #1
    ldrb w3, [x10], #1
    cmp w2, w3
    bne .noteq
    cmp w2, #0
    bne .loop
    mov x0, #0   // iguales
    ret

.noteq:
    sub x0, x2, x3
    ret


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