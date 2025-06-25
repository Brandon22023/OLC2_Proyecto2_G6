package repl

import (
	parser "compiler/parser"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"
)

type VariableEntry struct {
	Tipo  string
	Label string
	Valor string
}

type ARMVisitor struct {
	parser.BaseVlangVisitor
	Generator        *ArmGenerator
	PrintCount       int
	VarMap           map[string]VariableEntry
	ScopeTrace       *ScopeTrace
	UsesIntToAscii   bool
	UsesFloatToAscii bool
	UsesBoolToAscii  bool
	UsesRuneToAscii  bool
}

type PrintValue struct {
	Tipo  string
	Valor string
}

// debido a que dicho arbol puede que le sea dificil manejarlo
// ahora de igual forma leeremos los scopes donde estan los entornos de las variables
func (v *ARMVisitor) CollectAllVariables(scope *BaseScope) {
	for name := range scope.variables {
		entry, exists := v.VarMap[name]
		if !exists {
			continue
		}

		if entry.Tipo == "string" || entry.Tipo == "rune" {
			// No se generan datos duplicados aquí para variables que ya fueron manejadas
			// en VisitVariableDeclaration.
		}
	}
	for _, child := range scope.Children() {
		v.CollectAllVariables(child)
	}
}

func (v *ARMVisitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

func (v *ARMVisitor) VisitPrograma(ctx *parser.ProgramaContext) interface{} {
	v.VarMap = make(map[string]VariableEntry)
	v.PrintCount = 0

	// Etiquetas estándar para salto de línea y booleanos
	v.Generator.AddData(`msg_nl: .ascii "\n"`)
	v.Generator.AddData(`len_nl: .quad . - msg_nl`)
	v.Generator.AddData(`msg_true: .ascii "true"`)
	v.Generator.AddData(`len_true: .quad . - msg_true`)
	v.Generator.AddData(`msg_falsestr: .ascii "false"`)
	v.Generator.AddData(`len_falsestr: .quad . - msg_falsestr`)

	// 🛠 FIX: Agrega esta constante
	v.Generator.AddData(`.align 3`)
	v.Generator.AddData(`float_100: .double 100.0`)

	for _, decl := range ctx.AllDeclaraciones() {
		v.Visit(decl)
	}

	if v.Generator != nil && v.VarMap != nil && v.ScopeTrace != nil {
		v.CollectAllVariables(v.ScopeTrace.GlobalScope)
	}
	if v.UsesIntToAscii {
		v.Generator.AddIntToAsciiFunction()
	}
	if v.UsesFloatToAscii {
		v.Generator.AddFloatToAsciiFunction()
	}
	if v.UsesBoolToAscii {
		v.Generator.AddBoolToAsciiFunction()
	}
	if v.UsesRuneToAscii {
		v.Generator.AddRuneToAsciiFunction()
	}

	v.Generator.AddInstruction("    # Salida final")
	v.Generator.AddInstruction("mov x0, #0")
	v.Generator.AddInstruction("mov w8, #93")
	v.Generator.AddInstruction("svc #0")
	return nil
}

func (v *ARMVisitor) VisitDeclaraciones(ctx *parser.DeclaracionesContext) interface{} {
	if ctx.Stmt() != nil {
		return v.Visit(ctx.Stmt())
	}
	if ctx.VarDcl() != nil {
		return v.Visit(ctx.VarDcl())
	}
	if ctx.FuncMain() != nil {
		return v.Visit(ctx.FuncMain())
	}
	if ctx.FuncDcl() != nil {
		return v.Visit(ctx.FuncDcl())
	}
	if ctx.StructDcl() != nil {
		return v.Visit(ctx.StructDcl())
	}
	return nil
}

func (v *ARMVisitor) VisitFuncMain(ctx *parser.FuncMainContext) interface{} {
	fmt.Println("VisitFuncMain ejecutado")
	for _, decl := range ctx.Block().AllDeclaraciones() {
		v.Visit(decl)
	}
	return nil
}

func (v *ARMVisitor) VisitVariableDeclaration(ctx *parser.VariableDeclarationContext) interface{} {
	id := ctx.ID().GetText()
	tipo := normalizeTipo(ctx.TIPO().GetText())
	val := v.Visit(ctx.Expresion())

	if pv, ok := val.(PrintValue); ok {
		label := id

		v.VarMap[id] = VariableEntry{
			Tipo:  tipo,
			Label: label,
			Valor: pv.Valor,
		}

		switch tipo {
		case "int":
			v.Generator.AddData(fmt.Sprintf(".align 3\n%s: .quad %d", label, toInt(pv.Valor)))
		case "bool":
			b := 0
			if pv.Valor == "true" {
				b = 1
			}
			v.Generator.AddData(fmt.Sprintf(".align 3\n%s: .quad %d", label, b))
		case "float64":
			bits := math.Float64bits(toFloat(pv.Valor))
			v.Generator.AddData(fmt.Sprintf(".align 3\n%s: .quad 0x%x", label, bits))
		case "string":
			// ✅ Cambio importante: usar .asciz para strcmp
			v.Generator.AddData(fmt.Sprintf("%s: .asciz \"%s\"", label, pv.Valor))
			v.Generator.AddData(fmt.Sprintf("len_%s: .quad . - %s", label, label))
		case "rune":
			v.Generator.AddData(fmt.Sprintf("%s: .byte %d", label, toRune(pv.Valor)))
			v.Generator.AddData(fmt.Sprintf("len_%s: .quad 1", label))
		}
	}

	return nil
}

func (v *ARMVisitor) VisitPrintStatement(ctx *parser.PrintStatementContext) interface{} {
	fmt.Println("Entrando a VisitPrintStatement con exprs:", ctx.GetText())
	exprs := ctx.AllExpresion()

	for _, expr := range exprs {
		val := v.Visit(expr)

		var pv PrintValue
		ok := false

		if pval, isPv := val.(PrintValue); isPv {
			pv = pval
			ok = true
		} else if id, isId := val.(string); isId {
			entry, exists := v.VarMap[id]
			if !exists {
				fmt.Printf("ERROR: variable no encontrada: %s\n", id)
				continue
			}
			pv = PrintValue{Tipo: entry.Tipo, Valor: entry.Valor}
			ok = true
		}

		if !ok {
			fmt.Println("ERROR: valor no reconocido en print, por lo siguiente no se imprimirá:", val)
			continue
		}

		switch pv.Tipo {
		case "int", "entero":
			v.UsesIntToAscii = true
			intVal := toInt(pv.Valor)
			if intVal < 0 {
				// Escribe '-' en buffer
				v.Generator.AddLdr("x1", "buffer")
				v.Generator.AddMov("w3", "#45") // ASCII '-'
				v.Generator.AddInstruction("strb w3, [x1], #1")
				v.Generator.AddMov("x0", fmt.Sprintf("#%d", -intVal))
				v.Generator.AddBl("int_to_ascii")
				// Ajusta para imprimir el '-' más los dígitos
				v.Generator.AddInstruction("add x2, x0, #1") // x0 = longitud de dígitos, +1 por el '-'
				v.Generator.AddLdr("x1", "buffer")
				v.Generator.AddInstruction("mov x0, #1")
				v.Generator.AddInstruction("mov w8, #64")
				v.Generator.AddInstruction("svc #0")
			} else {
				v.Generator.AddMov("x0", fmt.Sprintf("#%d", intVal))
				v.Generator.AddLdr("x1", "buffer")
				v.Generator.AddBl("int_to_ascii")
				v.Generator.AddInstruction("mov x2, x0")
				v.Generator.AddInstruction("mov x1, x1")
				v.Generator.AddInstruction("mov x0, #1")
				v.Generator.AddInstruction("mov w8, #64")
				v.Generator.AddInstruction("svc #0")
			}

		case "float64", "decimal":
			v.UsesFloatToAscii = true

			label := ""

			// Si el valor es una variable, accede a la etiqueta original
			if id, ok := expr.(*parser.IdContext); ok {
				varName := id.GetText()
				entry, exists := v.VarMap[varName]
				if exists {
					label = entry.Label
				}
			}

			valorStr := pv.Valor
			if label != "" {
				// Imprimir directamente desde la variable
				v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", label))
				v.Generator.AddInstruction("ldr d0, [x1]")
			} else {
				// Imprimir un literal flotante (positivo o negativo)
				v.PrintCount++
				label = fmt.Sprintf("float_literal_%d", v.PrintCount)
				bits := math.Float64bits(toFloat(valorStr)) // Usa el valor real, con signo
				v.Generator.AddData(fmt.Sprintf(".align 3\n%s: .quad 0x%x", label, bits))
				v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", label))
				v.Generator.AddInstruction("ldr d0, [x1]")
			}

			// Llamar a la función de conversión y luego imprimir
			v.Generator.AddLdr("x1", "buffer")        // buffer
			v.Generator.AddBl("float_to_ascii")       // x0 ← length, x1 se mueve internamente
			v.Generator.AddInstruction("mov x2, x0")  // x2 = length
			v.Generator.AddLdr("x1", "buffer")        // ❗️RELOAD x1 to buffer
			v.Generator.AddInstruction("mov x0, #1")  // stdout
			v.Generator.AddInstruction("mov w8, #64") // syscall write
			v.Generator.AddInstruction("svc #0")

		case "bool", "booleano":
			v.UsesBoolToAscii = true
			if pv.Valor == "true" || pv.Valor == "false" {
				valBool := "0"
				if pv.Valor == "true" {
					valBool = "1"
				}
				v.Generator.AddMov("x0", valBool)
			} else if pv.Valor == "x2" {
				// El resultado de una comparación está en x2
				v.Generator.AddInstruction("mov x0, x2")
			} else {
				// Por si acaso, intenta convertir el valor a entero
				v.Generator.AddMov("x0", pv.Valor)
			}
			v.Generator.AddBl("bool_to_ascii")
			v.Generator.AddInstruction("mov x2, x0")
			v.Generator.AddInstruction("mov x1, x1")
			v.Generator.AddInstruction("mov x0, #1")
			v.Generator.AddInstruction("mov w8, #64")
			v.Generator.AddInstruction("svc #0")

		case "caracter":
			v.UsesRuneToAscii = true
			r := []rune(pv.Valor)[0]
			v.Generator.AddMov("x0", fmt.Sprintf("%d", r))
			v.Generator.AddLdr("x1", "buffer")
			v.Generator.AddBl("rune_to_ascii")
			v.Generator.AddInstruction("mov x2, x0")
			v.Generator.AddInstruction("mov x1, x1")
			v.Generator.AddInstruction("mov x0, #1")
			v.Generator.AddInstruction("mov w8, #64")
			v.Generator.AddInstruction("svc #0")

		case "string", "cadena":
			v.PrintCount++
			label := fmt.Sprintf("msg%d", v.PrintCount)
			lenLabel := fmt.Sprintf("len%d", v.PrintCount)
			v.Generator.AddData(fmt.Sprintf("%s: .ascii \"%s\"", label, pv.Valor))
			v.Generator.AddData(fmt.Sprintf("%s: .quad . - %s", lenLabel, label))
			v.printSyscall(label, lenLabel)

		case "rune":
			v.UsesRuneToAscii = true
			v.Generator.AddMov("x0", fmt.Sprintf("%d", toRune(pv.Valor)))
			v.Generator.AddLdr("x1", "buffer")
			v.Generator.AddBl("rune_to_ascii")
			v.Generator.AddInstruction("mov x2, x0")
			v.Generator.AddInstruction("mov x1, x1")
			v.Generator.AddInstruction("mov x0, #1")
			v.Generator.AddInstruction("mov w8, #64")
			v.Generator.AddInstruction("svc #0")

		}
	}

	// SOLO al final se imprime salto de línea
	v.Generator.AddMov("x0", "#1")
	v.Generator.AddLdr("x1", "msg_nl")
	v.Generator.AddLdr("x2", "len_nl")
	v.Generator.AddInstruction("ldr x2, [x2]")
	v.Generator.AddMov("w8", "#64")
	v.Generator.AddSvc()

	return nil
}

func (v *ARMVisitor) VisitId(ctx *parser.IdContext) interface{} {
	id := ctx.GetText()
	if entry, exists := v.VarMap[id]; exists {
		return PrintValue{Tipo: entry.Tipo, Valor: entry.Valor}
	}
	// Si no existe, devolver como string para evitar crash, pero advertir
	fmt.Printf("⚠️  Variable no encontrada: %s\n", id)
	return id
}

func (v *ARMVisitor) VisitValorexpr(ctx *parser.ValorexprContext) interface{} {
	return v.Visit(ctx.Valor())
}

func (v *ARMVisitor) VisitValorEntero(ctx *parser.ValorEnteroContext) interface{} {
	return PrintValue{Tipo: "int", Valor: ctx.GetText()}
}

func (v *ARMVisitor) VisitValorDecimal(ctx *parser.ValorDecimalContext) interface{} {
	return PrintValue{Tipo: "float64", Valor: ctx.GetText()}
}

func (v *ARMVisitor) VisitValorCadena(ctx *parser.ValorCadenaContext) interface{} {
	text := ctx.GetText()
	if len(text) >= 2 && text[0] == '"' && text[len(text)-1] == '"' {
		text = text[1 : len(text)-1]
	}
	return PrintValue{Tipo: "string", Valor: text}
}

func (v *ARMVisitor) VisitValorBooleano(ctx *parser.ValorBooleanoContext) interface{} {
	return PrintValue{Tipo: "bool", Valor: ctx.GetText()}
}

func (v *ARMVisitor) VisitValorCaracter(ctx *parser.ValorCaracterContext) interface{} {
	text := ctx.GetText()
	if len(text) == 3 && text[0] == '\'' && text[2] == '\'' {
		text = string(text[1])
	}
	return PrintValue{Tipo: "caracter", Valor: text}
}

func (v *ARMVisitor) VisitIMCPLICIT(ctx *parser.IMCPLICITContext) interface{} {
	id := ctx.ID().GetText()
	if entry, ok := v.VarMap[id]; ok {
		op := ctx.GetOp().GetText()
		val := v.Visit(ctx.Expresion())

		if pv, ok := val.(PrintValue); ok {
			a := toFloat(entry.Valor)
			b := toFloat(pv.Valor)
			var result string

			switch op {
			case "+=":
				result = formatNumber(a + b)
			case "-=":
				result = formatNumber(a - b)
			}

			v.replaceVarData(id, result)
		}
	}
	return nil
}

func (v *ARMVisitor) replaceVarData(id string, newValue string) {
	entry := v.VarMap[id]
	entry.Valor = newValue
	v.VarMap[id] = entry

	// Evitar reemplazo directo de .data
	if entry.Tipo == "string" || entry.Tipo == "rune" {
		entry.Valor = newValue
		v.VarMap[id] = entry
		return
	}

	label := entry.Label
	for i, line := range v.Generator.Data {
		if strings.HasPrefix(line, label+":") {
			v.Generator.Data[i] = fmt.Sprintf("%s: .quad %s", label, newValue)
			break
		}
	}
}

func (v *ARMVisitor) VisitAsignacionLUEGO(ctx *parser.AsignacionLUEGOContext) interface{} {
	id := ctx.ID().GetText()
	val := v.Visit(ctx.Expresion())

	if pv, ok := val.(PrintValue); ok {
		if entry, exists := v.VarMap[id]; exists {
			tipo := normalizeTipo(pv.Tipo)

			if tipo != entry.Tipo {
				fmt.Printf("ERROR: tipo incompatible en asignación a '%s': %s vs %s\n", id, entry.Tipo, tipo)
				return nil
			}

			switch tipo {
			case "int":
				v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", entry.Label))
				v.Generator.AddInstruction(fmt.Sprintf("mov x0, #%d", toInt(pv.Valor)))
				v.Generator.AddInstruction("str x0, [x1]")
			case "bool":
				b := 0
				if pv.Valor == "true" {
					b = 1
				}
				v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", entry.Label))
				v.Generator.AddInstruction(fmt.Sprintf("mov x0, #%d", b))
				v.Generator.AddInstruction("str x0, [x1]")
			case "float64":
				bits := math.Float64bits(toFloat(pv.Valor))
				tmp := v.Generator.GenerateUniqueLabel("flt")
				v.Generator.AddData(fmt.Sprintf(".align 3\n%s: .quad 0x%x", tmp, bits))
				v.Generator.AddInstruction(fmt.Sprintf("ldr x0, =%s", tmp))
				v.Generator.AddInstruction("ldr d0, [x0]")
				v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", entry.Label))
				v.Generator.AddInstruction("str d0, [x1]")
			case "string":
				label := v.Generator.GenerateUniqueLabel("str")
				v.Generator.AddData(fmt.Sprintf("%s: .ascii \"%s\"", label, pv.Valor))
				v.Generator.AddData(fmt.Sprintf("len_%s: .quad . - %s", label, label))
				entry.Label = label
				v.VarMap[id] = entry
			case "rune":
				v.Generator.AddInstruction(fmt.Sprintf("mov x0, #%d", toRune(pv.Valor)))
				v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", entry.Label))
				v.Generator.AddInstruction("strb w0, [x1]")
			}

			entry.Valor = pv.Valor
			v.VarMap[id] = entry
		}
	}

	return nil
}

func (v *ARMVisitor) VisitSumres(ctx *parser.SumresContext) interface{} {
	l := v.Visit(ctx.Expresion(0)).(PrintValue)
	r := v.Visit(ctx.Expresion(1)).(PrintValue)
	op := ctx.GetOp().GetText()

	switch op {
	case "+":
		if l.Tipo == "cadena" && r.Tipo == "cadena" {
			return PrintValue{Tipo: "cadena", Valor: l.Valor + r.Valor}
		}
		if l.Tipo == "entero" && r.Tipo == "entero" {
			return PrintValue{Tipo: "int", Valor: fmt.Sprintf("%d", toInt(l.Valor)+toInt(r.Valor))}
		}
		return PrintValue{Tipo: "float64", Valor: formatNumber(toFloat(l.Valor) + toFloat(r.Valor))}
	case "-":
		if l.Tipo == "entero" && r.Tipo == "entero" {
			return PrintValue{Tipo: "int", Valor: fmt.Sprintf("%d", toInt(l.Valor)-toInt(r.Valor))}
		}
		return PrintValue{Tipo: "float64", Valor: formatNumber(toFloat(l.Valor) - toFloat(r.Valor))}
	}
	return nil
}

func (v *ARMVisitor) VisitMultdivmod(ctx *parser.MultdivmodContext) interface{} {
	l := v.Visit(ctx.Expresion(0)).(PrintValue)
	r := v.Visit(ctx.Expresion(1)).(PrintValue)
	op := ctx.GetOp().GetText()

	switch op {
	case "*":
		return PrintValue{Tipo: "float64", Valor: formatNumber(toFloat(l.Valor) * toFloat(r.Valor))}
	case "/":
		return PrintValue{Tipo: "float64", Valor: formatNumber(toFloat(l.Valor) / toFloat(r.Valor))}
	case "%":
		return PrintValue{Tipo: "int", Valor: fmt.Sprintf("%d", toInt(l.Valor)%toInt(r.Valor))}
	}
	return nil
}

func toFloat(valor string) float64 {
	val, err := strconv.ParseFloat(valor, 64)
	if err != nil {
		return 0.0
	}
	return val
}

func toInt(valor string) int {
	val, err := strconv.Atoi(valor)
	if err != nil {
		return 0
	}
	return val
}

func toBool(valor string) int {
	if valor == "true" {
		return 1
	}
	return 0
}

func toRune(valor string) rune {
	runes := []rune(valor)
	if len(runes) > 0 {
		return runes[0]
	}
	return 0
}

func formatNumber(n float64) string {
	if n == float64(int64(n)) {
		return fmt.Sprintf("%d", int64(n))
	}
	return fmt.Sprintf("%.6g", n) // <-- formato más corto
}

func (v *ARMVisitor) printSyscall(label, lenLabel string) {
	v.Generator.AddMov("x0", "#1")
	v.Generator.AddLdr("x1", label)
	v.Generator.AddLdr("x2", lenLabel)
	v.Generator.AddInstruction("ldr x2, [x2]")
	v.Generator.AddMov("w8", "#64")
	v.Generator.AddSvc()
}

func normalizeTipo(t string) string {
	switch t {
	case "entero":
		return "int"
	case "booleano":
		return "bool"
	case "cadena":
		return "string"
	case "caracter":
		return "rune"
	case "decimal":
		return "float64"
	}
	return t
}

func (v *ARMVisitor) VisitVariableDeclarationImmutable(ctx *parser.VariableDeclarationImmutableContext) interface{} {
	id := ctx.ID().GetText()
	fmt.Println("DEBUG: Entrando a VariableDeclarationImmutable para:", id)

	if ctx.ASSIGN() == nil {
		fmt.Println("DEBUG: No hay asignación para", id)
		return nil
	}

	val := v.Visit(ctx.Expresion())
	if pv, ok := val.(PrintValue); ok {
		tipo := normalizeTipo(pv.Tipo)

		// REASIGNACIÓN
		if entry, exists := v.VarMap[id]; exists {
			fmt.Println("DEBUG: Reasignación detectada para:", id)

			if tipo != entry.Tipo {
				fmt.Printf("ERROR: Tipo incompatible en reasignación de '%s'. Esperado '%s', recibido '%s'\n", id, entry.Tipo, tipo)
				return nil
			}

			switch tipo {
			case "int":
				intVal := toInt(pv.Valor)
				v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", entry.Label))
				v.Generator.AddInstruction(fmt.Sprintf("mov x0, #%d", intVal))
				v.Generator.AddInstruction("str x0, [x1]")
			case "bool":
				boolVal := 0
				if pv.Valor == "true" {
					boolVal = 1
				}
				v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", entry.Label))
				v.Generator.AddInstruction(fmt.Sprintf("mov x0, #%d", boolVal))
				v.Generator.AddInstruction("str x0, [x1]")
			case "float64":
				floatBits := math.Float64bits(toFloat(pv.Valor))
				tempLabel := v.Generator.GenerateUniqueLabel("float_temp")
				v.Generator.AddData(fmt.Sprintf(".align 3\n%s: .quad 0x%x", tempLabel, floatBits))
				v.Generator.AddInstruction(fmt.Sprintf("ldr x0, =%s", tempLabel))
				v.Generator.AddInstruction("ldr x0, [x0]")
				v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", entry.Label))
				v.Generator.AddInstruction("str x0, [x1]")
			case "string", "rune":
				v.replaceVarData(id, pv.Valor)
			}

			// Actualiza el valor en VarMap
			entry.Valor = pv.Valor
			v.VarMap[id] = entry
			return nil
		}

		// DECLARACIÓN
		if _, exists := v.VarMap[id]; !exists {
			fmt.Println("DEBUG: Declaración nueva para:", id)

			label := id
			v.VarMap[id] = VariableEntry{Tipo: tipo, Valor: pv.Valor, Label: label}

			switch tipo {
			case "int":
				val := toInt(pv.Valor)
				v.Generator.AddData(fmt.Sprintf(".align 3\n%s: .quad %d", label, val))
				v.Generator.AddData(fmt.Sprintf("len_%s: .quad 8", label))
			case "bool":
				val := 0
				if pv.Valor == "true" {
					val = 1
				}
				v.Generator.AddData(fmt.Sprintf(".align 3\n%s: .quad %d", label, val))
				v.Generator.AddData(fmt.Sprintf("len_%s: .quad 8", label))
			case "float64":
				floatBits := math.Float64bits(toFloat(pv.Valor))
				tempLabel := v.Generator.GenerateUniqueLabel("float_temp")
				v.Generator.AddData(fmt.Sprintf(".align 3\n%s: .quad 0x%x", tempLabel, floatBits))
				v.Generator.AddInstruction(fmt.Sprintf("ldr x0, =%s", tempLabel))
				v.Generator.AddInstruction("ldr d0, [x0]")
				v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", label))
				v.Generator.AddInstruction("str d0, [x1]")

			case "string":
				v.Generator.AddData(fmt.Sprintf("%s: .asciz \"%s\"", label, pv.Valor))
				v.Generator.AddData(fmt.Sprintf("len_%s: .quad . - %s", label, label))

			case "rune":
				v.Generator.AddData(fmt.Sprintf("%s: .byte %d", label, toRune(pv.Valor)))
				v.Generator.AddData(fmt.Sprintf("len_%s: .quad 1", label))
			}
		}
	}
	return nil
}

func (v *ARMVisitor) VisitRelacionales(ctx *parser.RelacionalesContext) interface{} {
	left := v.Visit(ctx.Expresion(0)).(PrintValue)
	right := v.Visit(ctx.Expresion(1)).(PrintValue)
	op := ctx.GetOp().GetText() // ">", "<", etc.

	tipoLeft := normalizeTipo(left.Tipo)
	tipoRight := normalizeTipo(right.Tipo)
	if tipoLeft == "int" && tipoRight == "int" {
		v.Generator.AddMov("x0", left.Valor)
		v.Generator.AddMov("x1", right.Valor)
		v.Generator.AddInstruction("cmp x0, x1")

		labelTrue := v.Generator.GenerateUniqueLabel("cmp_true")
		labelEnd := v.Generator.GenerateUniqueLabel("cmp_end")

		v.Generator.AddMov("x2", "#0") // Por defecto: false

		switch op {
		case ">":
			v.Generator.AddInstruction(fmt.Sprintf("bgt %s", labelTrue))
		case "<":
			v.Generator.AddInstruction(fmt.Sprintf("blt %s", labelTrue))
		case ">=":
			v.Generator.AddInstruction(fmt.Sprintf("bge %s", labelTrue))
		case "<=":
			v.Generator.AddInstruction(fmt.Sprintf("ble %s", labelTrue))
		}

		v.Generator.AddInstruction(fmt.Sprintf("b %s", labelEnd))
		v.Generator.setLabel(labelTrue)
		v.Generator.AddMov("x2", "#1") // true
		v.Generator.setLabel(labelEnd)

		return PrintValue{Tipo: "bool", Valor: "x2"}
	}

	if (tipoLeft == "float64" && tipoRight == "float64") ||
		(tipoLeft == "int" && tipoRight == "float64") ||
		(tipoLeft == "float64" && tipoRight == "int") {

		// Convierte ambos a float64
		var leftVal, rightVal float64
		leftVal = toFloat(left.Valor)
		rightVal = toFloat(right.Valor)

		// Carga leftVal en d0
		v.PrintCount++
		labelLeft := fmt.Sprintf("cmp_float_left_%d", v.PrintCount)
		bitsLeft := math.Float64bits(leftVal)
		v.Generator.AddData(fmt.Sprintf(".align 3\n%s: .quad 0x%x", labelLeft, bitsLeft))
		v.Generator.AddInstruction(fmt.Sprintf("ldr x0, =%s", labelLeft))
		v.Generator.AddInstruction("ldr d0, [x0]")

		// Carga rightVal en d1
		v.PrintCount++
		labelRight := fmt.Sprintf("cmp_float_right_%d", v.PrintCount)
		bitsRight := math.Float64bits(rightVal)
		v.Generator.AddData(fmt.Sprintf(".align 3\n%s: .quad 0x%x", labelRight, bitsRight))
		v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", labelRight))
		v.Generator.AddInstruction("ldr d1, [x1]")

		// Compara d0 y d1
		v.Generator.AddInstruction("fcmp d0, d1")

		labelTrue := v.Generator.GenerateUniqueLabel("cmp_true")
		labelEnd := v.Generator.GenerateUniqueLabel("cmp_end")

		v.Generator.AddMov("x2", "#0") // Por defecto: false

		switch op {
		case ">":
			v.Generator.AddInstruction(fmt.Sprintf("bgt %s", labelTrue))
		case "<":
			v.Generator.AddInstruction(fmt.Sprintf("blt %s", labelTrue))
		case ">=":
			v.Generator.AddInstruction(fmt.Sprintf("bge %s", labelTrue))
		case "<=":
			v.Generator.AddInstruction(fmt.Sprintf("ble %s", labelTrue))
		}

		v.Generator.AddInstruction(fmt.Sprintf("b %s", labelEnd))
		v.Generator.setLabel(labelTrue)
		v.Generator.AddMov("x2", "#1") // true
		v.Generator.setLabel(labelEnd)

		return PrintValue{Tipo: "bool", Valor: "x2"}
	}

	// Si no es comparación de enteros, retorna false para evitar nil
	return PrintValue{Tipo: "bool", Valor: "false"}
}

func (v *ARMVisitor) VisitIgualdad(ctx *parser.IgualdadContext) interface{} {
	left := v.Visit(ctx.Expresion(0)).(PrintValue)
	right := v.Visit(ctx.Expresion(1)).(PrintValue)
	op := ctx.GetOp().GetText() // "==" o "!="

	tipoLeft := normalizeTipo(left.Tipo)
	tipoRight := normalizeTipo(right.Tipo)

	// --- Comparación int vs int ---
	if tipoLeft == "int" && tipoRight == "int" {
		v.Generator.AddMov("x0", left.Valor)
		v.Generator.AddMov("x1", right.Valor)
		v.Generator.AddInstruction("cmp x0, x1")

		labelTrue := v.Generator.GenerateUniqueLabel("cmp_eq_true")
		labelEnd := v.Generator.GenerateUniqueLabel("cmp_eq_end")

		v.Generator.AddMov("x2", "#0") // Por defecto: false

		switch op {
		case "==":
			v.Generator.AddInstruction(fmt.Sprintf("beq %s", labelTrue))
		case "!=":
			v.Generator.AddInstruction(fmt.Sprintf("bne %s", labelTrue))
		}

		v.Generator.AddInstruction(fmt.Sprintf("b %s", labelEnd))
		v.Generator.setLabel(labelTrue)
		v.Generator.AddMov("x2", "#1") // true
		v.Generator.setLabel(labelEnd)

		return PrintValue{Tipo: "bool", Valor: "x2"}
	}

	// --- Comparación float64/int/float64 ---
	if (tipoLeft == "float64" && tipoRight == "float64") ||
		(tipoLeft == "int" && tipoRight == "float64") ||
		(tipoLeft == "float64" && tipoRight == "int") {

		// Convierte ambos a float64
		var leftVal, rightVal float64
		leftVal = toFloat(left.Valor)
		rightVal = toFloat(right.Valor)

		// Carga leftVal en d0
		v.PrintCount++
		labelLeft := fmt.Sprintf("cmp_eq_float_left_%d", v.PrintCount)
		bitsLeft := math.Float64bits(leftVal)
		v.Generator.AddData(fmt.Sprintf(".align 3\n%s: .quad 0x%x", labelLeft, bitsLeft))
		v.Generator.AddInstruction(fmt.Sprintf("ldr x0, =%s", labelLeft))
		v.Generator.AddInstruction("ldr d0, [x0]")

		// Carga rightVal en d1
		v.PrintCount++
		labelRight := fmt.Sprintf("cmp_eq_float_right_%d", v.PrintCount)
		bitsRight := math.Float64bits(rightVal)
		v.Generator.AddData(fmt.Sprintf(".align 3\n%s: .quad 0x%x", labelRight, bitsRight))
		v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", labelRight))
		v.Generator.AddInstruction("ldr d1, [x1]")

		// Compara d0 y d1
		v.Generator.AddInstruction("fcmp d0, d1")

		labelTrue := v.Generator.GenerateUniqueLabel("cmp_eq_true")
		labelEnd := v.Generator.GenerateUniqueLabel("cmp_eq_end")

		v.Generator.AddMov("x2", "#0") // Por defecto: false

		switch op {
		case "==":
			v.Generator.AddInstruction(fmt.Sprintf("beq %s", labelTrue))
		case "!=":
			v.Generator.AddInstruction(fmt.Sprintf("bne %s", labelTrue))
		}

		v.Generator.AddInstruction(fmt.Sprintf("b %s", labelEnd))
		v.Generator.setLabel(labelTrue)
		v.Generator.AddMov("x2", "#1") // true
		v.Generator.setLabel(labelEnd)

		return PrintValue{Tipo: "bool", Valor: "x2"}
	}
	if tipoLeft == "string" && tipoRight == "string" {
		var result bool
		switch op {
		case "==":
			result = left.Valor == right.Valor
		case "!=":
			result = left.Valor != right.Valor
		}
		if result {
			return PrintValue{Tipo: "bool", Valor: "true"}
		}
		return PrintValue{Tipo: "bool", Valor: "false"}
	}

	// Si no es comparación válida, retorna false
	return PrintValue{Tipo: "bool", Valor: "false"}
}

func (v *ARMVisitor) VisitOPERADORESLOGICOS(ctx *parser.OPERADORESLOGICOSContext) interface{} {
	left := v.Visit(ctx.Expresion(0)).(PrintValue)
	right := v.Visit(ctx.Expresion(1)).(PrintValue)
	op := ctx.GetOp().GetText() // "&&" o "||"

	tipoLeft := normalizeTipo(left.Tipo)
	tipoRight := normalizeTipo(right.Tipo)

	// Solo booleanos
	if tipoLeft == "bool" && tipoRight == "bool" {
		labelTrue := v.Generator.GenerateUniqueLabel("logic_true")
		labelEnd := v.Generator.GenerateUniqueLabel("logic_end")

		switch left.Valor {
		case "x2":
			v.Generator.AddInstruction("mov x0, x2")
		case "true":
			v.Generator.AddMov("x0", "#1")
		case "false":
			v.Generator.AddMov("x0", "#0")
		default:
			v.Generator.AddMov("x0", left.Valor)
		}

		switch right.Valor {
		case "x2":
			v.Generator.AddInstruction("mov x1, x2")
		case "true":
			v.Generator.AddMov("x1", "#1")
		case "false":
			v.Generator.AddMov("x1", "#0")
		default:
			v.Generator.AddMov("x1", right.Valor)
		}

		switch op {
		case "&&":
			// x2 = (x0 != 0) && (x1 != 0)
			v.Generator.AddMov("x2", "#0")
			v.Generator.AddInstruction("cmp x0, #0")
			v.Generator.AddInstruction(fmt.Sprintf("beq %s", labelEnd))
			v.Generator.AddInstruction("cmp x1, #0")
			v.Generator.AddInstruction(fmt.Sprintf("beq %s", labelEnd))
			v.Generator.AddMov("x2", "#1")
			v.Generator.AddInstruction(fmt.Sprintf("b %s", labelEnd))
		case "||":
			// x2 = (x0 != 0) || (x1 != 0)
			v.Generator.AddMov("x2", "#1")
			v.Generator.AddInstruction("cmp x0, #0")
			v.Generator.AddInstruction(fmt.Sprintf("bne %s", labelTrue))
			v.Generator.AddInstruction("cmp x1, #0")
			v.Generator.AddInstruction(fmt.Sprintf("bne %s", labelTrue))
			v.Generator.AddMov("x2", "#0")
			v.Generator.AddInstruction(fmt.Sprintf("b %s", labelEnd))
			v.Generator.setLabel(labelTrue)
			v.Generator.AddMov("x2", "#1")
		}
		v.Generator.setLabel(labelEnd)
		return PrintValue{Tipo: "bool", Valor: "x2"}
	}

	// Si no son booleanos, retorna false
	return PrintValue{Tipo: "bool", Valor: "false"}
}

func (v *ARMVisitor) VisitUnario(ctx *parser.UnarioContext) interface{} {
	val := v.Visit(ctx.Expresion()).(PrintValue)
	op := ctx.GetOp().GetText()

	switch op {
	case "!":
		// Negación lógica: solo para booleanos
		tipo := normalizeTipo(val.Tipo)
		if tipo == "bool" {
			switch val.Valor {
			case "x2":
				// El valor está en x2 (resultado de una comparación)
				v.Generator.AddInstruction("mov x0, x2")
			case "true":
				v.Generator.AddMov("x0", "#1")
			case "false":
				v.Generator.AddMov("x0", "#0")
			default:
				v.Generator.AddMov("x0", val.Valor)
			}
			// NOT lógico: x2 = (x0 == 0) ? 1 : 0
			labelTrue := v.Generator.GenerateUniqueLabel("not_true")
			labelEnd := v.Generator.GenerateUniqueLabel("not_end")
			v.Generator.AddMov("x2", "#0")
			v.Generator.AddInstruction("cmp x0, #0")
			v.Generator.AddInstruction(fmt.Sprintf("beq %s", labelTrue))
			v.Generator.AddInstruction(fmt.Sprintf("b %s", labelEnd))
			v.Generator.setLabel(labelTrue)
			v.Generator.AddMov("x2", "#1")
			v.Generator.setLabel(labelEnd)
			return PrintValue{Tipo: "bool", Valor: "x2"}
		}
		// Si no es booleano, retorna false
		return PrintValue{Tipo: "bool", Valor: "false"}

	case "-":
		// Negación aritmética: para enteros y decimales
		tipo := normalizeTipo(val.Tipo)
		if tipo == "int" || tipo == "entero" {
			return PrintValue{Tipo: "int", Valor: fmt.Sprintf("%d", -toInt(val.Valor))}
		}
		if tipo == "float64" || tipo == "decimal" {
			return PrintValue{Tipo: "float64", Valor: formatNumber(-toFloat(val.Valor))}
		}
	}

	// Si no es un caso soportado, retorna el valor original
	return val
}

func (v *ARMVisitor) VisitControlStatement(ctx *parser.ControlStatementContext) interface{} {
	return v.Visit(ctx.Sentencias_control())
}

func (v *ARMVisitor) VisitIf_context(ctx *parser.If_contextContext) interface{} {
	return v.Visit(ctx.IfDcl())
}

func (v *ARMVisitor) VisitIfDcl(ctx *parser.IfDclContext) interface{} {
	fmt.Println("Entrando a VisitIfDcl con condición:", ctx.Expresion().GetText())
	endLabel := v.Generator.GenerateUniqueLabel("ifend")
	elseIfLabels := []string{}

	// Generar etiquetas para cada else if y else
	for range ctx.AllElseIfDcl() {
		elseIfLabels = append(elseIfLabels, v.Generator.GenerateUniqueLabel("elseif"))
	}
	if ctx.ElseCondicional() != nil {
		elseIfLabels = append(elseIfLabels, v.Generator.GenerateUniqueLabel("else"))
	}

	// --- IF PRINCIPAL ---
	val := v.Visit(ctx.Expresion()).(PrintValue)
	falseLabel := ""
	if len(elseIfLabels) > 0 {
		falseLabel = elseIfLabels[0]
	} else {
		falseLabel = endLabel
	}
	v.prepareCondBranch(val, falseLabel)

	for _, decl := range ctx.AllDeclaraciones() {
		v.Visit(decl)
	}
	v.Generator.B(endLabel)

	// --- ELSE IFs ---
	for i, elseIf := range ctx.AllElseIfDcl() {
		v.Generator.setLabel(elseIfLabels[i])
		val := v.Visit(elseIf.Expresion()).(PrintValue)

		nextLabel := ""
		if i+1 < len(elseIfLabels) {
			nextLabel = elseIfLabels[i+1]
		} else {
			nextLabel = endLabel
		}
		v.prepareCondBranch(val, nextLabel)

		for _, decl := range elseIf.AllDeclaraciones() {
			v.Visit(decl)
		}
		v.Generator.B(endLabel)
	}

	// --- ELSE FINAL ---
	if ctx.ElseCondicional() != nil {
		v.Generator.setLabel(elseIfLabels[len(elseIfLabels)-1])
		for _, decl := range ctx.ElseCondicional().AllDeclaraciones() {
			v.Visit(decl)
		}
	}

	v.Generator.setLabel(endLabel)
	return nil
}

func (v *ARMVisitor) prepareCondBranch(val PrintValue, labelFalse string) {
	switch val.Valor {
	case "x2":
		v.Generator.AddInstruction("cmp x2, #0")
		v.Generator.AddInstruction(fmt.Sprintf("beq %s", labelFalse))
	case "true":
		// no branch, continúa
	case "false":
		v.Generator.AddInstruction(fmt.Sprintf("b %s", labelFalse))
	default:
		v.Generator.AddMov("x0", val.Valor)
		v.Generator.AddInstruction("cmp x0, #0")
		v.Generator.AddInstruction(fmt.Sprintf("beq %s", labelFalse))
	}
}

func (v *ARMVisitor) VisitSwitch_context(ctx *parser.Switch_contextContext) interface{} {
	return v.Visit(ctx.SwitchDcl())
}

func (v *ARMVisitor) VisitSwitchDcl(ctx *parser.SwitchDclContext) interface{} {
	switchExpr := v.Visit(ctx.Expresion())

	pv, ok := switchExpr.(PrintValue)
	if !ok {
		fmt.Println("ERROR: expresión de switch no válida")
		return nil
	}

	v.Generator.AddStringCompareFunction() // Asegura strcmp una vez

	endLabel := v.Generator.GenerateUniqueLabel("switch_end")
	caseLabels := make([]string, len(ctx.AllCaseBlock()))
	defaultLabel := ""

	for i := range ctx.AllCaseBlock() {
		caseLabels[i] = v.Generator.GenerateUniqueLabel(fmt.Sprintf("case_%d", i))
	}

	if ctx.DefaultBlock() != nil {
		defaultLabel = v.Generator.GenerateUniqueLabel("switch_default")
	}

	// === Evaluar expresión del switch (poner en x0) ===
	var switchLabel string
	if pv.Tipo == "string" {
		varName := ctx.Expresion().GetText()
		entry, exists := v.VarMap[varName]

		if exists {
			switchLabel = entry.Label
			v.Generator.AddInstruction(fmt.Sprintf("ldr x0, =%s", switchLabel))
		} else {
			// Literal directa
			v.PrintCount++
			switchLabel = fmt.Sprintf("switch_expr_str_%d", v.PrintCount)
			v.Generator.AddData(fmt.Sprintf("%s: .asciz \"%s\"", switchLabel, pv.Valor))
			v.Generator.AddInstruction(fmt.Sprintf("ldr x0, =%s", switchLabel))
		}
	} else if pv.Tipo == "int" || pv.Tipo == "bool" {
		v.Generator.AddMov("x0", pv.Valor)
	} else {
		fmt.Println("ERROR: tipo de switch no soportado:", pv.Tipo)
		return nil
	}

	// === Comparaciones para cada case ===
	for i, caseBlock := range ctx.AllCaseBlock() {
		caseExpr := v.Visit(caseBlock.Expresion())
		casePv := caseExpr.(PrintValue)

		if pv.Tipo == "string" && casePv.Tipo == "string" {
			v.PrintCount++
			label := fmt.Sprintf("switch_case_str_%d", v.PrintCount)
			v.Generator.AddData(fmt.Sprintf("%s: .asciz \"%s\"", label, casePv.Valor))
			// Re-cargar x0 = valor original del switch
			v.Generator.AddInstruction(fmt.Sprintf("ldr x0, =%s", switchLabel))
			v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", label))
			v.Generator.AddBl("strcmp")
			v.Generator.AddInstruction("cmp x0, #0")
			v.Generator.AddInstruction(fmt.Sprintf("beq %s", caseLabels[i]))

		} else if pv.Tipo == "int" && casePv.Tipo == "int" {
			v.Generator.AddMov("x1", casePv.Valor)
			v.Generator.AddInstruction("cmp x0, x1")
			v.Generator.AddInstruction(fmt.Sprintf("beq %s", caseLabels[i]))
		} else if pv.Tipo == "bool" && casePv.Tipo == "bool" {
			val := 0
			if casePv.Valor == "true" {
				val = 1
			}
			v.Generator.AddMov("x1", fmt.Sprintf("#%d", val))
			v.Generator.AddInstruction("cmp x0, x1")
			v.Generator.AddInstruction(fmt.Sprintf("beq %s", caseLabels[i]))
		} else {
			fmt.Printf("WARN: tipo no compatible en switch/case: %s vs %s\n", pv.Tipo, casePv.Tipo)
		}
	}

	// === Si ningún case coincidió, salta al default o al final ===
	if defaultLabel != "" {
		v.Generator.B(defaultLabel)
	} else {
		v.Generator.B(endLabel)
	}

	// === Cuerpos de los case ===
	for i, caseBlock := range ctx.AllCaseBlock() {
		v.Generator.setLabel(caseLabels[i])
		for _, decl := range caseBlock.AllDeclaraciones() {
			v.Visit(decl)
		}
		v.Generator.B(endLabel)
	}

	// === Cuerpo del default ===
	if ctx.DefaultBlock() != nil {
		v.Generator.setLabel(defaultLabel)
		for _, decl := range ctx.DefaultBlock().AllDeclaraciones() {
			v.Visit(decl)
		}
	}

	v.Generator.setLabel(endLabel)
	return nil
}
