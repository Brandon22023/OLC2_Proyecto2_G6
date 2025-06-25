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
msg1: .ascii "=== Archivo de prueba básico ==="
len1: .quad . - msg1
msg2: .ascii "==== Declaración de variables ===="
len2: .quad . - msg2
.align 3
puntos_declaracion: .quad 0
msg3: .ascii "Declaración explícita con tipo y valor"
len3: .quad . - msg3
.align 3
entero: .quad 42
.align 3
decimal: .quad 0x400921f9f01b866e
texto: .asciz "Hola, mundo!"
len_texto: .quad . - texto
.align 3
booleano: .quad 1
msg4: .ascii "entero:"
len4: .quad . - msg4
msg5: .ascii "decimal:"
len5: .quad . - msg5
msg6: .ascii "texto:"
len6: .quad . - msg6
msg7: .ascii "Hola, mundo!"
len7: .quad . - msg7
msg8: .ascii "booleano:"
len8: .quad . - msg8
.align 3
cmp_float_left_9: .quad 0x400921f9f01b866e
.align 3
cmp_float_right_10: .quad 0x4008000000000000
msg11: .ascii "OK Declaración explícita con tipo y valor: correcto"
len11: .quad . - msg11
msg12: .ascii "X Declaración explícita con tipo y valor: incorrecto"
len12: .quad . - msg12
msg13: .ascii "\nDeclaración explícita con tipo y sin valor"
len13: .quad . - msg13
.align 3
entero_sin_valor: .quad 0
.align 3
decimal_sin_valor: .quad 0x0
texto_sin_valor: .asciz ""
len_texto_sin_valor: .quad . - texto_sin_valor
.align 3
booleano_sin_valor: .quad 0
msg14: .ascii "enteroSinValor:"
len14: .quad . - msg14
msg15: .ascii "decimalSinValor:"
len15: .quad . - msg15
msg16: .ascii "textoSinValor:"
len16: .quad . - msg16
msg17: .ascii ""
len17: .quad . - msg17
msg18: .ascii "booleanoSinValor:"
len18: .quad . - msg18
.align 3
cmp_eq_float_left_19: .quad 0x0
.align 3
cmp_eq_float_right_20: .quad 0x0
msg21: .ascii "OK Declaración explícita con tipo y sin valor: correcto"
len21: .quad . - msg21
msg22: .ascii "X Declaración explícita con tipo y sin valor: incorrecto"
len22: .quad . - msg22
msg23: .ascii "\nErrores de redeclaración"
len23: .quad . - msg23
msg24: .ascii "X Errores de redeclaración: incorrecto"
len24: .quad . - msg24
msg25: .ascii "OK Errores de redeclaración: correcto"
len25: .quad . - msg25
msg26: .ascii "\n==== Asignación de variables ===="
len26: .quad . - msg26
.align 3
puntos_asignacion: .quad 0
msg27: .ascii "Asignación con tipo correcto"
len27: .quad . - msg27
.align 3
float_temp_1: .quad 0x4023cccccccccccd
msg28: .ascii "entero:"
len28: .quad . - msg28
msg29: .ascii "decimal:"
len29: .quad . - msg29
msg30: .ascii "texto:"
len30: .quad . - msg30
msg31: .ascii "Texto modificado"
len31: .quad . - msg31
msg32: .ascii "booleano:"
len32: .quad . - msg32
.align 3
cmp_eq_float_left_33: .quad 0x4023cccccccccccd
.align 3
cmp_eq_float_right_34: .quad 0x4023cccccccccccd
msg35: .ascii "OK Asignación con tipo correcto: correcto"
len35: .quad . - msg35
msg36: .ascii "X Asignación con tipo correcto: incorrecto"
len36: .quad . - msg36
msg37: .ascii "\nAsignación con tipo incorrecto"
len37: .quad . - msg37
msg38: .ascii "OK Asignación con tipo incorrecto: Se detectaron errores de tipo correctamente"
len38: .quad . - msg38
msg39: .ascii "\n==== Operaciones Aritméticas ===="
len39: .quad . - msg39
.align 3
puntos_ari: .quad 0
.align 3
resultado_suma_1: .quad 15
.align 3
resultado_suma_2: .quad 0x4030000000000000
.align 3
resultado_suma_3: .quad 0x402f000000000000
.align 3
resultado_suma_4: .quad 0x402f000000000000
msg40: .ascii "10 + 5 ="
len40: .quad . - msg40
msg41: .ascii "10.5 + 5.5 ="
len41: .quad . - msg41
msg42: .ascii "10 + 5.5 ="
len42: .quad . - msg42
msg43: .ascii "10.5 + 5 ="
len43: .quad . - msg43
.align 3
cmp_eq_float_left_44: .quad 0x4030000000000000
.align 3
cmp_eq_float_right_45: .quad 0x4030000000000000
.align 3
cmp_eq_float_left_46: .quad 0x402f000000000000
.align 3
cmp_eq_float_right_47: .quad 0x402f000000000000
.align 3
cmp_eq_float_left_48: .quad 0x402f000000000000
.align 3
cmp_eq_float_right_49: .quad 0x402f000000000000
msg50: .ascii "OK Suma: correcto"
len50: .quad . - msg50
msg51: .ascii "X Suma: incorrecto"
len51: .quad . - msg51
.align 3
resultado_resta_1: .quad 5
.align 3
resultado_resta_2: .quad 0x4014000000000000
.align 3
resultado_resta_3: .quad 0x4012000000000000
.align 3
resultado_resta_4: .quad 0x4016000000000000
msg52: .ascii "10 - 5 ="
len52: .quad . - msg52
msg53: .ascii "10.5 - 5.5 ="
len53: .quad . - msg53
msg54: .ascii "10 - 5.5 ="
len54: .quad . - msg54
msg55: .ascii "10.5 - 5 ="
len55: .quad . - msg55
.align 3
cmp_eq_float_left_56: .quad 0x4014000000000000
.align 3
cmp_eq_float_right_57: .quad 0x4014000000000000
.align 3
cmp_eq_float_left_58: .quad 0x4012000000000000
.align 3
cmp_eq_float_right_59: .quad 0x4012000000000000
.align 3
cmp_eq_float_left_60: .quad 0x4016000000000000
.align 3
cmp_eq_float_right_61: .quad 0x4016000000000000
msg62: .ascii "OK Resta: correcto"
len62: .quad . - msg62
msg63: .ascii "X Resta: incorrecto"
len63: .quad . - msg63
msg64: .ascii "\n==== Operaciones Relacionales ===="
len64: .quad . - msg64
.align 3
puntos_rel: .quad 0
.align 3
cmp_eq_float_left_65: .quad 0x4025000000000000
.align 3
cmp_eq_float_right_66: .quad 0x4025000000000000
.align 3
cmp_eq_float_left_67: .quad 0x4025000000000000
.align 3
cmp_eq_float_right_68: .quad 0x4016000000000000
msg69: .ascii "OK Relacionales: correcto"
len69: .quad . - msg69
msg70: .ascii "X Relacionales: incorrecto"
len70: .quad . - msg70
msg71: .ascii "\n==== Operaciones Lógicas ===="
len71: .quad . - msg71
.align 3
puntos_log: .quad 0
msg72: .ascii "\n==== print ===="
len72: .quad . - msg72
.align 3
puntos_print: .quad 2
.align 3
float_literal_73: .quad 0x40091eb851eb851f
msg74: .ascii "Texto de prueba"
len74: .quad . - msg74
msg75: .ascii "Hola"
len75: .quad . - msg75
msg76: .ascii "\n==== Manejo de valor nulo ===="
len76: .quad . - msg76
.align 3
puntos_nil: .quad 2
msg77: .ascii "\n==== Punto y coma opcional ===="
len77: .quad . - msg77
.align 3
puntos_semicolon: .quad 2
.align 3
a: .quad 10
.align 3
b: .quad 20
.align 3
c: .quad 30
msg78: .ascii "a ="
len78: .quad . - msg78
msg79: .ascii "b ="
len79: .quad . - msg79
msg80: .ascii "c = a + b ="
len80: .quad . - msg80
msg81: .ascii "\n=== Tabla de Resultados ==="
len81: .quad . - msg81
msg82: .ascii "| Asignación de variables"
len82: .quad . - msg82
msg83: .ascii "| Operaciones Aritméticas  |"
len83: .quad . - msg83
msg84: .ascii "| Operaciones Relacionales |"
len84: .quad . - msg84
msg85: .ascii "| Operaciones Lógicas      |"
len85: .quad . - msg85
msg86: .ascii "| print                  | "
len86: .quad . - msg86
msg87: .ascii "| Manejo de nulo           | "
len87: .quad . - msg87
msg88: .ascii "| Punto y coma opcional    |"
len88: .quad . - msg88
msg89: .ascii "| TOTAL                    |"
len89: .quad . - msg89
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
mov x0, #1
ldr x1, =msg3
ldr x2, =len3
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg4
ldr x2, =len4
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #42
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
mov x1, x1
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
ldr x1, =msg5
ldr x2, =len5
ldr x2, [x2]
mov w8, #64
svc #0
ldr x1, =decimal
ldr d0, [x1]
ldr x1, =buffer
bl float_to_ascii
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
ldr x1, =msg6
ldr x2, =len6
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg7
ldr x2, =len7
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg8
ldr x2, =len8
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, 1
bl bool_to_ascii
mov x2, x0
mov x1, x1
mov x0, #1
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, 42
mov x1, 42
cmp x0, x1
mov x2, #0
beq cmp_eq_true_1
b cmp_eq_end_1
cmp_eq_true_1:
mov x2, #1
cmp_eq_end_1:
ldr x0, =cmp_float_left_9
ldr d0, [x0]
ldr x1, =cmp_float_right_10
ldr d1, [x1]
fcmp d0, d1
mov x2, #0
bgt cmp_true_1
b cmp_end_1
cmp_true_1:
mov x2, #1
cmp_end_1:
mov x0, x2
mov x1, x2
mov x2, #0
cmp x0, #0
beq logic_end_1
cmp x1, #0
beq logic_end_1
mov x2, #1
b logic_end_1
logic_end_1:
mov x0, x2
mov x1, #1
mov x2, #0
cmp x0, #0
beq logic_end_2
cmp x1, #0
beq logic_end_2
mov x2, #1
b logic_end_2
logic_end_2:
mov x0, x2
mov x1, #0
mov x2, #0
cmp x0, #0
beq logic_end_3
cmp x1, #0
beq logic_end_3
mov x2, #1
b logic_end_3
logic_end_3:
cmp x2, #0
beq else_1
mov x0, #1
ldr x1, =msg11
ldr x2, =len11
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
else_1:
mov x0, #1
ldr x1, =msg12
ldr x2, =len12
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
ifend_1:
mov x0, #1
ldr x1, =msg13
ldr x2, =len13
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg14
ldr x2, =len14
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #0
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
mov x1, x1
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
ldr x1, =msg15
ldr x2, =len15
ldr x2, [x2]
mov w8, #64
svc #0
ldr x1, =decimal_sin_valor
ldr d0, [x1]
ldr x1, =buffer
bl float_to_ascii
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
ldr x1, =msg16
ldr x2, =len16
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg17
ldr x2, =len17
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg18
ldr x2, =len18
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, 0
bl bool_to_ascii
mov x2, x0
mov x1, x1
mov x0, #1
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, 0
mov x1, 0
cmp x0, x1
mov x2, #0
beq cmp_eq_true_2
b cmp_eq_end_2
cmp_eq_true_2:
mov x2, #1
cmp_eq_end_2:
ldr x0, =cmp_eq_float_left_19
ldr d0, [x0]
ldr x1, =cmp_eq_float_right_20
ldr d1, [x1]
fcmp d0, d1
mov x2, #0
beq cmp_eq_true_3
b cmp_eq_end_3
cmp_eq_true_3:
mov x2, #1
cmp_eq_end_3:
mov x0, x2
mov x1, x2
mov x2, #0
cmp x0, #0
beq logic_end_4
cmp x1, #0
beq logic_end_4
mov x2, #1
b logic_end_4
logic_end_4:
mov x0, x2
mov x1, #1
mov x2, #0
cmp x0, #0
beq logic_end_5
cmp x1, #0
beq logic_end_5
mov x2, #1
b logic_end_5
logic_end_5:
mov x0, x2
mov x1, #0
mov x2, #0
cmp x0, #0
beq logic_end_6
cmp x1, #0
beq logic_end_6
mov x2, #1
b logic_end_6
logic_end_6:
cmp x2, #0
beq else_2
mov x0, #1
ldr x1, =msg21
ldr x2, =len21
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
else_2:
mov x0, #1
ldr x1, =msg22
ldr x2, =len22
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
ifend_2:
mov x0, #1
ldr x1, =msg23
ldr x2, =len23
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, 42
mov x1, 100
cmp x0, x1
mov x2, #0
beq cmp_eq_true_4
b cmp_eq_end_4
cmp_eq_true_4:
mov x2, #1
cmp_eq_end_4:
cmp x2, #0
beq else_3
mov x0, #1
ldr x1, =msg24
ldr x2, =len24
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
else_3:
mov x0, #1
ldr x1, =msg25
ldr x2, =len25
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
ifend_3:
mov x0, #1
ldr x1, =msg26
ldr x2, =len26
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg27
ldr x2, =len27
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
mov x0, #99
str x0, [x1]
ldr x0, =float_temp_1
ldr x0, [x0]
ldr x1, =decimal
str x0, [x1]
mov x0, #1
mov x2, #0
cmp x0, #0
beq not_true_1
b not_end_1
not_true_1:
mov x2, #1
not_end_1:
ldr x1, =booleano
mov x0, #0
str x0, [x1]
mov x0, #1
ldr x1, =msg28
ldr x2, =len28
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #99
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
mov x1, x1
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
ldr x1, =msg29
ldr x2, =len29
ldr x2, [x2]
mov w8, #64
svc #0
ldr x1, =decimal
ldr d0, [x1]
ldr x1, =buffer
bl float_to_ascii
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
ldr x1, =msg30
ldr x2, =len30
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg31
ldr x2, =len31
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg32
ldr x2, =len32
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, x2
bl bool_to_ascii
mov x2, x0
mov x1, x1
mov x0, #1
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, 99
mov x1, 99
cmp x0, x1
mov x2, #0
beq cmp_eq_true_5
b cmp_eq_end_5
cmp_eq_true_5:
mov x2, #1
cmp_eq_end_5:
ldr x0, =cmp_eq_float_left_33
ldr d0, [x0]
ldr x1, =cmp_eq_float_right_34
ldr d1, [x1]
fcmp d0, d1
mov x2, #0
beq cmp_eq_true_6
b cmp_eq_end_6
cmp_eq_true_6:
mov x2, #1
cmp_eq_end_6:
mov x0, x2
mov x1, x2
mov x2, #0
cmp x0, #0
beq logic_end_7
cmp x1, #0
beq logic_end_7
mov x2, #1
b logic_end_7
logic_end_7:
mov x0, x2
mov x1, #1
mov x2, #0
cmp x0, #0
beq logic_end_8
cmp x1, #0
beq logic_end_8
mov x2, #1
b logic_end_8
logic_end_8:
mov x0, x2
mov x1, #0
mov x2, #0
cmp x0, #0
beq logic_end_9
cmp x1, #0
beq logic_end_9
mov x2, #1
b logic_end_9
logic_end_9:
cmp x2, #0
beq else_4
mov x0, #1
ldr x1, =msg35
ldr x2, =len35
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
else_4:
mov x0, #1
ldr x1, =msg36
ldr x2, =len36
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
ifend_4:
mov x0, #1
ldr x1, =msg37
ldr x2, =len37
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg38
ldr x2, =len38
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg39
ldr x2, =len39
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg40
ldr x2, =len40
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #15
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
mov x1, x1
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
ldr x1, =msg41
ldr x2, =len41
ldr x2, [x2]
mov w8, #64
svc #0
ldr x1, =resultado_suma_2
ldr d0, [x1]
ldr x1, =buffer
bl float_to_ascii
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
ldr x1, =msg42
ldr x2, =len42
ldr x2, [x2]
mov w8, #64
svc #0
ldr x1, =resultado_suma_3
ldr d0, [x1]
ldr x1, =buffer
bl float_to_ascii
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
ldr x1, =msg43
ldr x2, =len43
ldr x2, [x2]
mov w8, #64
svc #0
ldr x1, =resultado_suma_4
ldr d0, [x1]
ldr x1, =buffer
bl float_to_ascii
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
mov x0, 15
mov x1, 15
cmp x0, x1
mov x2, #0
beq cmp_eq_true_7
b cmp_eq_end_7
cmp_eq_true_7:
mov x2, #1
cmp_eq_end_7:
ldr x0, =cmp_eq_float_left_44
ldr d0, [x0]
ldr x1, =cmp_eq_float_right_45
ldr d1, [x1]
fcmp d0, d1
mov x2, #0
beq cmp_eq_true_8
b cmp_eq_end_8
cmp_eq_true_8:
mov x2, #1
cmp_eq_end_8:
mov x0, x2
mov x1, x2
mov x2, #0
cmp x0, #0
beq logic_end_10
cmp x1, #0
beq logic_end_10
mov x2, #1
b logic_end_10
logic_end_10:
ldr x0, =cmp_eq_float_left_46
ldr d0, [x0]
ldr x1, =cmp_eq_float_right_47
ldr d1, [x1]
fcmp d0, d1
mov x2, #0
beq cmp_eq_true_9
b cmp_eq_end_9
cmp_eq_true_9:
mov x2, #1
cmp_eq_end_9:
mov x0, x2
mov x1, x2
mov x2, #0
cmp x0, #0
beq logic_end_11
cmp x1, #0
beq logic_end_11
mov x2, #1
b logic_end_11
logic_end_11:
ldr x0, =cmp_eq_float_left_48
ldr d0, [x0]
ldr x1, =cmp_eq_float_right_49
ldr d1, [x1]
fcmp d0, d1
mov x2, #0
beq cmp_eq_true_10
b cmp_eq_end_10
cmp_eq_true_10:
mov x2, #1
cmp_eq_end_10:
mov x0, x2
mov x1, x2
mov x2, #0
cmp x0, #0
beq logic_end_12
cmp x1, #0
beq logic_end_12
mov x2, #1
b logic_end_12
logic_end_12:
cmp x2, #0
beq else_5
mov x0, #1
ldr x1, =msg50
ldr x2, =len50
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
else_5:
mov x0, #1
ldr x1, =msg51
ldr x2, =len51
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
ifend_5:
mov x0, #1
ldr x1, =msg52
ldr x2, =len52
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #5
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
mov x1, x1
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
ldr x1, =msg53
ldr x2, =len53
ldr x2, [x2]
mov w8, #64
svc #0
ldr x1, =resultado_resta_2
ldr d0, [x1]
ldr x1, =buffer
bl float_to_ascii
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
ldr x1, =msg54
ldr x2, =len54
ldr x2, [x2]
mov w8, #64
svc #0
ldr x1, =resultado_resta_3
ldr d0, [x1]
ldr x1, =buffer
bl float_to_ascii
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
ldr x1, =msg55
ldr x2, =len55
ldr x2, [x2]
mov w8, #64
svc #0
ldr x1, =resultado_resta_4
ldr d0, [x1]
ldr x1, =buffer
bl float_to_ascii
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
mov x0, 5
mov x1, 5
cmp x0, x1
mov x2, #0
beq cmp_eq_true_11
b cmp_eq_end_11
cmp_eq_true_11:
mov x2, #1
cmp_eq_end_11:
ldr x0, =cmp_eq_float_left_56
ldr d0, [x0]
ldr x1, =cmp_eq_float_right_57
ldr d1, [x1]
fcmp d0, d1
mov x2, #0
beq cmp_eq_true_12
b cmp_eq_end_12
cmp_eq_true_12:
mov x2, #1
cmp_eq_end_12:
mov x0, x2
mov x1, x2
mov x2, #0
cmp x0, #0
beq logic_end_13
cmp x1, #0
beq logic_end_13
mov x2, #1
b logic_end_13
logic_end_13:
ldr x0, =cmp_eq_float_left_58
ldr d0, [x0]
ldr x1, =cmp_eq_float_right_59
ldr d1, [x1]
fcmp d0, d1
mov x2, #0
beq cmp_eq_true_13
b cmp_eq_end_13
cmp_eq_true_13:
mov x2, #1
cmp_eq_end_13:
mov x0, x2
mov x1, x2
mov x2, #0
cmp x0, #0
beq logic_end_14
cmp x1, #0
beq logic_end_14
mov x2, #1
b logic_end_14
logic_end_14:
ldr x0, =cmp_eq_float_left_60
ldr d0, [x0]
ldr x1, =cmp_eq_float_right_61
ldr d1, [x1]
fcmp d0, d1
mov x2, #0
beq cmp_eq_true_14
b cmp_eq_end_14
cmp_eq_true_14:
mov x2, #1
cmp_eq_end_14:
mov x0, x2
mov x1, x2
mov x2, #0
cmp x0, #0
beq logic_end_15
cmp x1, #0
beq logic_end_15
mov x2, #1
b logic_end_15
logic_end_15:
cmp x2, #0
beq else_6
mov x0, #1
ldr x1, =msg62
ldr x2, =len62
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
else_6:
mov x0, #1
ldr x1, =msg63
ldr x2, =len63
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
ifend_6:
mov x0, #1
ldr x1, =msg64
ldr x2, =len64
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, 10
mov x1, 10
cmp x0, x1
mov x2, #0
beq cmp_eq_true_15
b cmp_eq_end_15
cmp_eq_true_15:
mov x2, #1
cmp_eq_end_15:
mov x0, 10
mov x1, 5
cmp x0, x1
mov x2, #0
bne cmp_eq_true_16
b cmp_eq_end_16
cmp_eq_true_16:
mov x2, #1
cmp_eq_end_16:
mov x0, x2
mov x1, x2
mov x2, #0
cmp x0, #0
beq logic_end_16
cmp x1, #0
beq logic_end_16
mov x2, #1
b logic_end_16
logic_end_16:
ldr x0, =cmp_eq_float_left_65
ldr d0, [x0]
ldr x1, =cmp_eq_float_right_66
ldr d1, [x1]
fcmp d0, d1
mov x2, #0
beq cmp_eq_true_17
b cmp_eq_end_17
cmp_eq_true_17:
mov x2, #1
cmp_eq_end_17:
mov x0, x2
mov x1, x2
mov x2, #0
cmp x0, #0
beq logic_end_17
cmp x1, #0
beq logic_end_17
mov x2, #1
b logic_end_17
logic_end_17:
ldr x0, =cmp_eq_float_left_67
ldr d0, [x0]
ldr x1, =cmp_eq_float_right_68
ldr d1, [x1]
fcmp d0, d1
mov x2, #0
bne cmp_eq_true_18
b cmp_eq_end_18
cmp_eq_true_18:
mov x2, #1
cmp_eq_end_18:
mov x0, x2
mov x1, x2
mov x2, #0
cmp x0, #0
beq logic_end_18
cmp x1, #0
beq logic_end_18
mov x2, #1
b logic_end_18
logic_end_18:
cmp x2, #0
beq else_7
mov x0, #1
ldr x1, =msg69
ldr x2, =len69
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
else_7:
mov x0, #1
ldr x1, =msg70
ldr x2, =len70
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
ifend_7:
mov x0, #1
ldr x1, =msg71
ldr x2, =len71
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg72
ldr x2, =len72
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #42
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
mov x1, x1
mov x0, #1
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
ldr x1, =float_literal_73
ldr d0, [x1]
ldr x1, =buffer
bl float_to_ascii
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
ldr x1, =msg74
ldr x2, =len74
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
mov x1, x1
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
ldr x1, =msg75
ldr x2, =len75
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #123
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
mov x1, x1
mov x0, #1
mov w8, #64
svc #0
mov x0, 1
bl bool_to_ascii
mov x2, x0
mov x1, x1
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
ldr x1, =msg76
ldr x2, =len76
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg77
ldr x2, =len77
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg78
ldr x2, =len78
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #10
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
mov x1, x1
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
ldr x1, =msg79
ldr x2, =len79
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #20
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
mov x1, x1
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
ldr x1, =msg80
ldr x2, =len80
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #30
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
mov x1, x1
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
ldr x1, =msg81
ldr x2, =len81
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg82
ldr x2, =len82
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
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
mov x1, x1
mov x0, #1
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #2
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
mov x1, x1
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
ldr x1, =msg83
ldr x2, =len83
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
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
mov x1, x1
mov x0, #1
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #2
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
mov x1, x1
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
ldr x1, =msg84
ldr x2, =len84
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
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
mov x1, x1
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
ldr x1, =msg85
ldr x2, =len85
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
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
mov x1, x1
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
ldr x1, =msg86
ldr x2, =len86
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #2
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
mov x1, x1
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
ldr x1, =msg87
ldr x2, =len87
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #2
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
mov x1, x1
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
ldr x1, =msg88
ldr x2, =len88
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #2
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
mov x1, x1
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
ldr x1, =msg89
ldr x2, =len89
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
ldr x1, =buffer
bl int_to_ascii
mov x2, x0
mov x1, x1
mov x0, #1
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


float_to_ascii:
    // ✅ Guardar copia de d0 antes de convertirlo
    fmov d4, d0               // ← Copia para mantener el original

    // Convertir parte entera de float en d0
    fcvtzu x0, d0             // x0 = int(d0)
    mov x2, x1                // x2 = buffer
    mov x3, #10
    mov x4, #0
    mov x5, x0                // copia para loop
    cmp x5, #0
    bne float_int_loop
    mov w6, #'0'
    strb w6, [x2], #1
    mov x4, #1
    b float_int_done

float_int_loop:
    udiv x6, x5, x3
    msub x7, x6, x3, x5
    add w7, w7, #'0'
    strb w7, [x2], #1
    mov x5, x6
    add x4, x4, #1
    cmp x5, #0
    bne float_int_loop

float_int_done:
    sub x2, x2, x4
    mov x5, x4
    mov x6, #0

float_reverse_loop:
    cmp x6, x5
    bge float_reverse_done
    add x7, x2, x6
    sub x8, x2, x6
    add x8, x8, x5
    sub x8, x8, #1
    ldrb w9, [x7]
    ldrb w10, [x8]
    strb w10, [x7]
    strb w9, [x8]
    add x6, x6, #1
    cmp x6, x5, lsr #1
    blt float_reverse_loop

float_reverse_done:
    add x1, x2, x5
    mov w6, #'.'
    strb w6, [x1], #1

    // ✅ Calcular parte decimal usando d4 (no d0, ya destruido)
    scvtf d1, x0
    fsub d2, d4, d1
    ldr x7, =float_100
    ldr d3, [x7]
    fmul d2, d2, d3
    fcvtzu x0, d2

    // Convertir dos dígitos decimales
    mov x3, #10
    udiv x4, x0, x3
    msub x5, x4, x3, x0
    add w4, w4, #'0'
    add w5, w5, #'0'
    strb w4, [x1], #1
    strb w5, [x1], #1

    // Fin de cadena
    mov w6, #0
    strb w6, [x1]

    sub x0, x1, x2      // x0 = longitud
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