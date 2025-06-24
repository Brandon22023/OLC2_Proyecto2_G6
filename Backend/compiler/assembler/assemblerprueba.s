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
msg1: .ascii "HOla mundo arm"
len1: .quad . - msg1
.align 3
float_literal_2: .quad 0x4034800000000000
.align 3
float_literal_3: .quad 0x400921cac083126f
msg4: .ascii "hola"
len4: .quad . - msg4
.align 3
float_literal_5: .quad 0x4026666666666666
msg6: .ascii "ab"
len6: .quad . - msg6
.align 3
x: .quad 10
s: .ascii "Hola"
len_s: .quad . - s
.align 3
y: .quad 0x4025000000000000
k: .byte 66
len_k: .quad 1
.align 3
l: .quad 0
msg7: .ascii "valor s: "
len7: .quad . - msg7
msg8: .ascii "Hola"
len8: .quad . - msg8
.align 3
var: .quad 30
msg9: .ascii "valor de var: "
len9: .quad . - msg9
msg10: .ascii "nuevo valor de var: "
len10: .quad . - msg10
.align 3
u: .quad 10
.align 3
o: .quad 0x40091eb851eb851f
.align 3
b: .quad 0
r: .byte 90
len_r: .quad 1
.align 3
float_temp_1: .quad 0x40191eb851eb851f
msg11: .ascii "Mundo"
len11: .quad . - msg11
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
ldr x1, =msg1
ldr x2, =len1
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
ldr x1, =float_literal_2
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
mov x0, 97
ldr x1, =buffer
bl rune_to_ascii
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
ldr x1, =float_literal_3
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
ldr x1, =float_literal_5
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
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #7
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
mov x0, 65
ldr x1, =buffer
bl rune_to_ascii
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
ldr x1, =msg7
ldr x2, =len7
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg8
ldr x2, =len8
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #1
ldr x1, =msg_nl
ldr x2, =len_nl
ldr x2, [x2]
mov w8, #64
svc #0
ldr x1, =y
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
mov x0, 66
ldr x1, =buffer
bl rune_to_ascii
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
ldr x1, =msg9
ldr x2, =len9
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
ldr x1, =var
mov x0, #50
str x0, [x1]
mov x0, #1
ldr x1, =msg10
ldr x2, =len10
ldr x2, [x2]
mov w8, #64
svc #0
mov x0, #50
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
ldr x1, =u
mov x0, #20
str x0, [x1]
ldr x0, =float_temp_1
ldr x0, [x0]
ldr x1, =o
str x0, [x1]
ldr x1, =b
mov x0, #1
str x0, [x1]
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
ldr x1, =o
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
mov x0, 88
ldr x1, =buffer
bl rune_to_ascii
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


rune_to_ascii:
    mov x2, x1
    strb w0, [x2]
    mov w3, #0
    strb w3, [x2, #1]
    mov x0, #1
    ret

// Finalizar programa


// Foreign functions


//libreria estandar
// Funciones estándar aquí