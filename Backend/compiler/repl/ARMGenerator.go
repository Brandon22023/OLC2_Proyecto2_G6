package repl

import (
	"fmt"
	"strings"
)

// Tipos de objetos para la pila
type StackObjectType int

const (
	Int StackObjectType = iota
	Float
	String
	Boolean
	Rune
	NUL
)

type StackObject struct {
	Type   StackObjectType
	Length int
	Depth  int
	Id     string
	Offset *int
}

// Simulación de una librería estándar
type StandardLibrary struct {
	Symbols map[string]string
	used    map[string]bool
}

func NewStandardLibrary() *StandardLibrary {
	return &StandardLibrary{
		Symbols: make(map[string]string),
		used:    make(map[string]bool),
	}
}
func (g *ArmGenerator) AddData(data string) {
	g.Data = append(g.Data, data)
}
func (g *ArmGenerator) AddInstruction(instr string) {
	g.Instructions = append(g.Instructions, instr)
}
func (g *ArmGenerator) AddIntToAsciiFunction() {
	g.FuncInstrucions = append(g.FuncInstrucions, `
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
`)
}

func (g *ArmGenerator) AddMov(rd, imm string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("mov %s, %s", rd, imm))
}

func (g *ArmGenerator) AddLdr(rd, label string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("ldr %s, =%s", rd, label))
}

func (g *ArmGenerator) AddBl(label string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("bl %s", label))
}

func (g *ArmGenerator) AddSvc() {
	g.Instructions = append(g.Instructions, "svc #0")
}
func (g *ArmGenerator) ToString() string {
	return g.String()
}

func (s *StandardLibrary) Use(name string) {
	s.used[name] = true
}

func (s *StandardLibrary) GetFunctionDefinitions() string {
	// Aquí deberías retornar las definiciones de funciones usadas
	return "// Funciones estándar aquí"
}

// ARMGenerator adaptado
type ArmGenerator struct {
	Instructions    []string
	StdLib          *StandardLibrary
	FuncInstrucions []string
	StackObjects    []StackObject
	Depth           int
	LabelCounter    int
	UsesIntToAscii  bool     // <-- Añade este campo para compatibilidad
	Data            []string // <-- Añade este campo para visitorARM.go
	LabelCounters   map[string]int
}

func NewArmGenerator() *ArmGenerator {
	return &ArmGenerator{
		Instructions:    []string{},
		StdLib:          NewStandardLibrary(),
		FuncInstrucions: []string{},
		StackObjects:    []StackObject{},
		Depth:           0,
		LabelCounter:    0,
		LabelCounters:   make(map[string]int),
	}
}

func (g *ArmGenerator) getLabel() string {
	label := fmt.Sprintf("L%d", g.LabelCounter)
	g.LabelCounter++
	return label
}

func (g *ArmGenerator) setLabel(label string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("%s:", label))
}

// Operaciones de pila
func (g *ArmGenerator) PushObject(obj StackObject) {
	g.StackObjects = append(g.StackObjects, obj)
}

func (g *ArmGenerator) TopObject() StackObject {
	return g.StackObjects[len(g.StackObjects)-1]
}

func (g *ArmGenerator) PopObject() StackObject {
	if len(g.StackObjects) == 0 {
		panic("No hay objetos en la pila")
	}
	obj := g.StackObjects[len(g.StackObjects)-1]
	g.StackObjects = g.StackObjects[:len(g.StackObjects)-1]
	return obj
}

func (g *ArmGenerator) PushConstant(obj StackObject, value interface{}) {
	switch obj.Type {
	case Int:
		g.Mov("X0", value.(int))
		g.Push("X0")
	case Float:
		// Conversión de float64 a uint64
		floatBits := int64(value.(float64))
		for i := 0; i < 4; i++ {
			part := (floatBits >> (i * 16)) & 0xFFFF
			if i == 0 {
				g.Instructions = append(g.Instructions, fmt.Sprintf("MOVZ X0, #%d, LSL #0", part))
			} else {
				g.Instructions = append(g.Instructions, fmt.Sprintf("MOVK X0, #%d, LSL #%d", part, i*16))
			}
		}
		g.Push("X0")
	case String:
		// Simulación: solo push de dirección ficticia
		g.Push("HP")
		// Aquí deberías convertir el string a bytes y simular el push caracter por caracter
	case Boolean:
		val := 0
		if value.(bool) {
			val = 1
		}
		g.Mov("X0", val)
		g.Push("X0")
	case Rune:
		g.Mov("X0", int(value.(rune)))
		g.Push("X0")
	}
	g.PushObject(obj)
}

func (g *ArmGenerator) PopObjectTo(rd string) StackObject {
	obj := g.PopObject()
	g.Pop(rd)
	return obj
}

// Métodos para crear objetos de pila
func (g *ArmGenerator) IntObject() StackObject {
	return StackObject{Type: Int, Length: 8, Depth: g.Depth}
}
func (g *ArmGenerator) FloatObject() StackObject {
	return StackObject{Type: Float, Length: 8, Depth: g.Depth}
}
func (g *ArmGenerator) StringObject() StackObject {
	return StackObject{Type: String, Length: 8, Depth: g.Depth}
}
func (g *ArmGenerator) BooleanObject() StackObject {
	return StackObject{Type: Boolean, Length: 8, Depth: g.Depth}
}
func (g *ArmGenerator) RuneObject() StackObject {
	return StackObject{Type: Rune, Length: 8, Depth: g.Depth}
}
func (g *ArmGenerator) CloneObject(obj StackObject) StackObject {
	return StackObject{
		Type:   obj.Type,
		Length: obj.Length,
		Depth:  obj.Depth,
		Id:     obj.Id,
	}
}

// Manejo de scopes
func (g *ArmGenerator) NewScope() {
	g.Depth++
}
func (g *ArmGenerator) EndScope() int {
	byteOffset := 0
	for i := len(g.StackObjects) - 1; i >= 0; i-- {
		if g.StackObjects[i].Depth == g.Depth {
			byteOffset += g.StackObjects[i].Length
			g.StackObjects = append(g.StackObjects[:i], g.StackObjects[i+1:]...)
		} else {
			break
		}
	}
	g.Depth--
	return byteOffset
}

func (g *ArmGenerator) TagObject(id string) {
	if len(g.StackObjects) > 0 {
		g.StackObjects[len(g.StackObjects)-1].Id = id
	}
}

func (g *ArmGenerator) GetObject(id string) (int, StackObject) {
	byteOffset := 0
	for i := len(g.StackObjects) - 1; i >= 0; i-- {
		byteOffset += g.StackObjects[i].Length
		if g.StackObjects[i].Id == id {
			return byteOffset - g.StackObjects[i].Length, g.StackObjects[i]
		}
	}
	panic(fmt.Sprintf("No se encontró el objeto con ID: %s", id))
}

// Instrucciones ARM simuladas
func (g *ArmGenerator) Add(rd, rs1, rs2 string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("ADD %s, %s, %s", rd, rs1, rs2))
}
func (g *ArmGenerator) Sub(rd, rs1, rs2 string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("SUB %s, %s, %s", rd, rs1, rs2))
}
func (g *ArmGenerator) Mul(rd, rs1, rs2 string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("MUL %s, %s, %s", rd, rs1, rs2))
}
func (g *ArmGenerator) Div(rd, rs1, rs2 string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("SDIV %s, %s, %s", rd, rs1, rs2))
}
func (g *ArmGenerator) Addi(rd, rs1 string, imm int) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("ADDI %s, %s, %d", rd, rs1, imm))
}
func (g *ArmGenerator) Str(rs1, rs2 string, offset int) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("STR %s, [%s, #%d]", rs1, rs2, offset))
}
func (g *ArmGenerator) Strb(rs1, rs2 string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("STRB %s, [%s]", rs1, rs2))
}
func (g *ArmGenerator) Ldr(rd, rs1 string, offset int) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("LDR %s, [%s, #%d]", rd, rs1, offset))
}
func (g *ArmGenerator) Mov(rd string, imm int) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("MOV %s, #%d", rd, imm))
}
func (g *ArmGenerator) Push(rs string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("STR %s, [SP, #-8]!", rs))
}
func (g *ArmGenerator) Pop(rd string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("LDR %s, [SP], #8", rd))
}
func (g *ArmGenerator) Svc() {
	g.Instructions = append(g.Instructions, "SVC #0")
}

// Operaciones con flotantes
func (g *ArmGenerator) Scvtf(rd, rs string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("SCVTF %s, %s", rd, rs))
}
func (g *ArmGenerator) Frintm(rd, rs string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("FRINTM %s, %s", rd, rs))
}
func (g *ArmGenerator) Fmov(rd, rs string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("FMOV %s, %s", rd, rs))
}
func (g *ArmGenerator) Fadd(rd, rs1, rs2 string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("FADD %s, %s, %s", rd, rs1, rs2))
}
func (g *ArmGenerator) Fsub(rd, rs1, rs2 string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("FSUB %s, %s, %s", rd, rs1, rs2))
}
func (g *ArmGenerator) Fmul(rd, rs1, rs2 string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("FMUL %s, %s, %s", rd, rs1, rs2))
}
func (g *ArmGenerator) Fdiv(rd, rs1, rs2 string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("FDIV %s, %s, %s", rd, rs1, rs2))
}

// Comparaciones y saltos
func (g *ArmGenerator) Cmp(rs1, rs2 string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("CMP %s, %s", rs1, rs2))
}
func (g *ArmGenerator) Fcmp(rs1, rs2 string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("FCMP %s, %s", rs1, rs2))
}
func (g *ArmGenerator) Beq(label string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("BEQ %s", label))
}
func (g *ArmGenerator) Bne(label string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("BNE %s", label))
}
func (g *ArmGenerator) Bgt(label string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("BGT %s", label))
}
func (g *ArmGenerator) Blt(label string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("BLT %s", label))
}
func (g *ArmGenerator) Bge(label string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("BGE %s", label))
}
func (g *ArmGenerator) Ble(label string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("BLE %s", label))
}
func (g *ArmGenerator) B(label string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("B %s", label))
}
func (g *ArmGenerator) Cbz(rs, label string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("CBZ %s, %s", rs, label))
}
func (g *ArmGenerator) Neg(rd, rs string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("NEG %s, %s", rd, rs))
}
func (g *ArmGenerator) Fneg(rd, rs string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("FNEG %s, %s", rd, rs))
}
func (g *ArmGenerator) Br(label string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("BR %s", label))
}
func (g *ArmGenerator) Ldrb(rd, rs1 string, offset int) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("LDRB %s, [%s, #%d]", rd, rs1, offset))
}

// Finalizar programa
func (g *ArmGenerator) EndProgram() {
	g.Mov("X0", 0)
	g.Mov("X8", 93)
	g.Svc()
}

func (g *ArmGenerator) Emit(instruction string) {
	g.Instructions = append(g.Instructions, instruction)
}

func (g *ArmGenerator) PrintInteger(rs string) {
	g.StdLib.Use("print_integer")
	g.Instructions = append(g.Instructions, fmt.Sprintf("MOV X0, %s", rs))
	g.Instructions = append(g.Instructions, "BL print_integer")
}

func (g *ArmGenerator) PrintFloat() {
	g.StdLib.Use("print_integer")
	g.StdLib.Use("print_double")
	g.Instructions = append(g.Instructions, "BL print_double")
}

func (g *ArmGenerator) Comment(comment string) {
	g.Instructions = append(g.Instructions, fmt.Sprintf("// %s", comment))
}

func (g *ArmGenerator) PrintString(rs string) {
	g.StdLib.Use("print_string")
	g.Instructions = append(g.Instructions, fmt.Sprintf("MOV X0, %s", rs))
	g.Instructions = append(g.Instructions, "BL print_string")
}

func (g *ArmGenerator) PrintChar(ch rune) {
	label := fmt.Sprintf("char_%d", ch)
	if _, ok := g.StdLib.Symbols[label]; !ok {
		val := string(ch)
		if ch == '\n' {
			val = "\\n"
		}
		g.StdLib.Symbols[label] = fmt.Sprintf("%s: .ascii \"%s\"", label, val)
	}
	g.Instructions = append(g.Instructions, fmt.Sprintf(`
MOV X0, #1
ADR X1, %s
MOV X2, #1
MOV X8, #64
SVC #0`, label))
}

func (g *ArmGenerator) ConcatStrings() {
	g.StdLib.Use("string_concat")
	//str2 := g.PopObjectTo("X1")
	//str1 := g.PopObjectTo("X0")
	g.Comment("Ensure stack alignment")
	g.Instructions = append(g.Instructions, "MOV X9, SP")
	g.Instructions = append(g.Instructions, "AND X9, X9, #-16")
	g.Instructions = append(g.Instructions, "MOV SP, X9")
	g.Comment("Call string_concat (str1 + str2)")
	g.Instructions = append(g.Instructions, "BL string_concat")
	g.Push("X0")
	g.PushObject(g.StringObject())
}

func (g *ArmGenerator) String() string {
	var sb strings.Builder
	sb.WriteString(".data\n")
	sb.WriteString("heap: .space 4096\n")
	sb.WriteString("heap_ptr: .quad heap\n")
	sb.WriteString("buffer: .space 32\n")
	for _, data := range g.Data {
		sb.WriteString(data + "\n")
	}
	sb.WriteString(".text\n")
	sb.WriteString(".global malloc\n")
	sb.WriteString(`malloc:
    mov x2, x10
    add x0, x2, x0
    mov x10, x0
    mov x0, x2
    ret
`)
	sb.WriteString(".global _start\n")
	sb.WriteString("_start:\n")
	sb.WriteString("    adr x10, heap\n")
	// Instrucciones principales
	for _, instr := range g.Instructions {
		sb.WriteString(instr + "\n")
	}
	// FUNCIONES AUXILIARES (como int_to_ascii) AQUÍ
	for _, f := range g.FuncInstrucions {
		sb.WriteString(f + "\n")
	}
	sb.WriteString("// Finalizar programa\n")
	sb.WriteString("\n\n// Foreign functions\n")
	sb.WriteString("\n\n//libreria estandar\n")
	sb.WriteString(g.StdLib.GetFunctionDefinitions())
	return sb.String()
}

// StringToBytesArray convierte un string en un slice de bytes y agrega un terminador nulo al final.
func StringToBytesArray(str string) []byte {
	resultado := []byte(str)
	resultado = append(resultado, 0) // Agregar el terminador nulo al final
	return resultado
}

func (g *ArmGenerator) GenerateUniqueLabel(base string) string {
	count := g.LabelCounters[base] + 1
	g.LabelCounters[base] = count
	return fmt.Sprintf("%s_%d", base, count)
}

func (g *ArmGenerator) DeclareQuadVariable(name string, value int) {
	g.Data = append(g.Data,
		fmt.Sprintf(".align 3\n%s: .quad %d", name, value))
}

func (g *ArmGenerator) AddFloatToAsciiFunction() {
	g.FuncInstrucions = append(g.FuncInstrucions, `
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

`)
}

func (g *ArmGenerator) AddRuneToAsciiFunction() {
	g.FuncInstrucions = append(g.FuncInstrucions, `
rune_to_ascii:
    mov x2, x1
    strb w0, [x2]
    mov w3, #0
    strb w3, [x2, #1]
    mov x0, #1
    ret
`)
}

func (g *ArmGenerator) AddBoolToAsciiFunction() {
	g.FuncInstrucions = append(g.FuncInstrucions, `
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
`)
}
