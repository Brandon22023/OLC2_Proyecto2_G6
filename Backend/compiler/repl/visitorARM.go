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
			v.Generator.AddData(fmt.Sprintf("%s: .ascii \"%s\"", label, pv.Valor))
			v.Generator.AddData(fmt.Sprintf("len_%s: .quad . - %s", label, label))
		case "rune":
			v.Generator.AddData(fmt.Sprintf("%s: .byte %d", label, toRune(pv.Valor)))
			v.Generator.AddData(fmt.Sprintf("len_%s: .quad 1", label))
		}
	}

	return nil
}

func (v *ARMVisitor) VisitPrintStatement(ctx *parser.PrintStatementContext) interface{} {
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
			fmt.Println("ERROR: valor no reconocido en print")
			continue
		}

		switch pv.Tipo {
		case "int", "entero":
			v.UsesIntToAscii = true
			v.Generator.AddMov("x0", "#"+pv.Valor)
			v.Generator.AddLdr("x1", "buffer")
			v.Generator.AddBl("int_to_ascii")
			v.Generator.AddInstruction("mov x2, x0")
			v.Generator.AddInstruction("mov x1, x1")
			v.Generator.AddInstruction("mov x0, #1")
			v.Generator.AddInstruction("mov w8, #64")
			v.Generator.AddInstruction("svc #0")

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

			if label != "" {
				// Imprimir directamente desde la variable
				v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", label))
				v.Generator.AddInstruction("ldr d0, [x1]")
			} else {
				// Imprimir un literal flotante
				v.PrintCount++
				label = fmt.Sprintf("float_literal_%d", v.PrintCount)
				bits := math.Float64bits(toFloat(pv.Valor))
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
			valBool := "0"
			if pv.Valor == "true" {
				valBool = "1"
			}
			v.Generator.AddMov("x0", valBool)
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
	return ctx.GetText()
}

func (v *ARMVisitor) VisitValorexpr(ctx *parser.ValorexprContext) interface{} {
	return v.Visit(ctx.Valor())
}

func (v *ARMVisitor) VisitValorEntero(ctx *parser.ValorEnteroContext) interface{} {
	return PrintValue{Tipo: "entero", Valor: ctx.GetText()}
}

func (v *ARMVisitor) VisitValorDecimal(ctx *parser.ValorDecimalContext) interface{} {
	return PrintValue{Tipo: "decimal", Valor: ctx.GetText()}
}

func (v *ARMVisitor) VisitValorCadena(ctx *parser.ValorCadenaContext) interface{} {
	text := ctx.GetText()
	if len(text) >= 2 && text[0] == '"' && text[len(text)-1] == '"' {
		text = text[1 : len(text)-1]
	}
	return PrintValue{Tipo: "cadena", Valor: text}
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
			v.VarMap[id] = VariableEntry{Tipo: normalizeTipo(tipo), Valor: pv.Valor, Label: label}

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
				v.Generator.AddInstruction("ldr d0, [x0]") // ✅ CORREGIDO
				v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", label))
				v.Generator.AddInstruction("str d0, [x1]") // ✅ CORREGIDO

			case "string", "rune":
				v.Generator.AddData(fmt.Sprintf("%s: .ascii \"%s\"", label, pv.Valor))
				v.Generator.AddData(fmt.Sprintf("len_%s: .quad . - %s", label, label))
			}
		}
	}
	return nil
}
