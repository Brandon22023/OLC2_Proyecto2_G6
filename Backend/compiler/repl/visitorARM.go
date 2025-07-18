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
	FuncLabels       map[string]string //las funciones se guardaran aqui en un label para mas estabilidad
	breakLabels      []string
	continueLabels   []string
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

	// Sección .data: etiquetas estándar
	v.Generator.AddData(`msg_nl: .ascii "\n"`)
	v.Generator.AddData(`len_nl: .quad . - msg_nl`)
	v.Generator.AddData(`msg_true: .ascii "true"`)
	v.Generator.AddData(`len_true: .quad . - msg_true`)
	v.Generator.AddData(`msg_falsestr: .ascii "false"`)
	v.Generator.AddData(`len_falsestr: .quad . - msg_falsestr`)
	v.Generator.AddData(`.align 3`)
	v.Generator.AddData(`float_100: .double 100.0`)

	// 1. PRIMER PASO: Registrar labels de todas las funciones
	for _, decl := range ctx.AllDeclaraciones() {
		if decl.FuncDcl() != nil {
			funcName := decl.FuncDcl().(*parser.FuncDclContext).ID().GetText()
			label := "fn_" + funcName
			v.FuncLabels[funcName] = label
		}
	}

	// 2. SEGUNDO PASO: Procesar funciones y variables globales
	for _, decl := range ctx.AllDeclaraciones() {
		if decl.FuncDcl() != nil {
			v.VisitFuncDcl(decl.FuncDcl().(*parser.FuncDclContext))
		} else if decl.FuncMain() != nil {
			// No ejecutes main aquí, lo haremos en _start
			continue
		} else {
			v.Visit(decl)
		}
	}

	// 3. Generar _start y poner el cuerpo de main ahí
	v.Generator.AddInstruction(".global _start")
	v.Generator.setLabel("_start")
	v.Generator.AddInstruction("    adr x10, heap")

	// Busca el main y genera su cuerpo aquí
	for _, decl := range ctx.AllDeclaraciones() {
		if decl.FuncMain() != nil {
			mainCtx := decl.FuncMain().(*parser.FuncMainContext)
			for _, d := range mainCtx.Block().AllDeclaraciones() {
				v.Visit(d)
			}
		}
	}

	// Termina el programa
	v.Generator.AddInstruction("    mov x0, #0")
	v.Generator.AddInstruction("    mov w8, #93")
	v.Generator.AddInstruction("    svc #0")

	// Funciones auxiliares si se usaron
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

	return nil
}
func (v *ARMVisitor) VisitTransfersentence(ctx *parser.TransfersentenceContext) interface{} {
    res := v.VisitChildren(ctx)
    if ret, ok := res.(PrintValue); ok && ret.Tipo == "__return__" {
        return ret
    }
    return res
}
func (v *ARMVisitor) VisitDeclaraciones(ctx *parser.DeclaracionesContext) interface{} {
	txt := ctx.GetText()
    //fmt.Println("[DEBUG] VisitDeclaraciones:", txt)
    // Parche: si el texto empieza con "return", llama manualmente a VisitReturnStatement
    if strings.HasPrefix(txt, "return") {
        // Crea un contexto falso o busca el hijo correcto
        // Si tienes acceso al hijo, llama a v.VisitReturnStatement(hijo)
        // Si no, usa el parche del valor:
        valor := txt[6:]
        //fmt.Println("[DEBUG] Valor de retorno (parche desde declaraciones):", valor)
        v.Generator.AddMov("x0", valor)
        return PrintValue{Tipo: "__return__"}
    }
	if ctx.Stmt() != nil {
        res := v.Visit(ctx.Stmt())
        // Propaga el return especial
        if ret, ok := res.(PrintValue); ok && ret.Tipo == "__return__" {
            return ret
        }
        return res
    }
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
	v.Generator.setLabel("fn_main")
	//fmt.Println("VisitFuncMain ejecutado")
	for _, decl := range ctx.Block().AllDeclaraciones() {
		v.Visit(decl)
	}
	v.Generator.AddInstruction("ret")
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
	//fmt.Println("Entrando a VisitPrintStatement con exprs:", ctx.GetText())
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
				fmt.Printf("⚠️ Variable no encontrada: %s\n", id)
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

			_, err := strconv.Atoi(pv.Valor)
			switch {
			case pv.Valor == "x0":
			case pv.Valor == "x2":
				v.Generator.AddInstruction("mov x0, x2")
			case err == nil:
				v.Generator.AddMov("x0", fmt.Sprintf("#%s", pv.Valor))
			default:
				v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", pv.Valor))
				v.Generator.AddInstruction("ldr x0, [x1]")
			}

			v.Generator.AddLdr("x1", "buffer")
			v.Generator.AddBl("int_to_ascii")
			v.Generator.AddInstruction("mov x2, x0")
			v.Generator.AddLdr("x1", "buffer")
			v.Generator.AddInstruction("mov x0, #1")
			v.Generator.AddInstruction("mov w8, #64")
			v.Generator.AddInstruction("svc #0")

		case "string", "cadena":
			// Etiquetas únicas seguras
			v.PrintCount++
			label := fmt.Sprintf("msg_print_%d", v.PrintCount)
			lenLabel := fmt.Sprintf("len_print_%d", v.PrintCount)

			escaped := strings.ReplaceAll(pv.Valor, `"`, `\"`)
			v.Generator.AddData(fmt.Sprintf("%s: .ascii \"%s\"", label, escaped))
			v.Generator.AddData(fmt.Sprintf("%s: .quad . - %s", lenLabel, label))

			v.Generator.AddInstruction("mov x0, #1")
			v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", label))
			v.Generator.AddInstruction(fmt.Sprintf("ldr x2, =%s", lenLabel))
			v.Generator.AddInstruction("ldr x2, [x2]")
			v.Generator.AddInstruction("mov w8, #64")
			v.Generator.AddInstruction("svc #0")

		case "bool", "booleano":
			v.UsesBoolToAscii = true
			if pv.Valor == "true" {
				v.Generator.AddMov("x0", "1")
			} else if pv.Valor == "false" {
				v.Generator.AddMov("x0", "0")
			} else {
				v.Generator.AddMov("x0", pv.Valor)
			}
			v.Generator.AddBl("bool_to_ascii")
			v.Generator.AddInstruction("mov x2, x0")
			v.Generator.AddLdr("x1", "buffer")
			v.Generator.AddInstruction("mov x0, #1")
			v.Generator.AddInstruction("mov w8, #64")
			v.Generator.AddInstruction("svc #0")

		case "caracter", "rune":
			v.UsesRuneToAscii = true
			r := []rune(pv.Valor)[0]
			v.Generator.AddMov("x0", fmt.Sprintf("%d", r))
			v.Generator.AddLdr("x1", "buffer")
			v.Generator.AddBl("rune_to_ascii")
			v.Generator.AddInstruction("mov x2, x0")
			v.Generator.AddLdr("x1", "buffer")
			v.Generator.AddInstruction("mov x0, #1")
			v.Generator.AddInstruction("mov w8, #64")
			v.Generator.AddInstruction("svc #0")

		case "float64", "decimal":
			v.UsesFloatToAscii = true
			label := ""

			if id, ok := expr.(*parser.IdContext); ok {
				varName := id.GetText()
				entry, exists := v.VarMap[varName]
				if exists {
					label = entry.Label
				}
			}

			valorStr := pv.Valor
			if label != "" {
				v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", label))
				v.Generator.AddInstruction("ldr d0, [x1]")
			} else {
				v.PrintCount++
				label = fmt.Sprintf("float_literal_%d", v.PrintCount)
				bits := math.Float64bits(toFloat(valorStr))
				v.Generator.AddData(fmt.Sprintf(".align 3\n%s: .quad 0x%x", label, bits))
				v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", label))
				v.Generator.AddInstruction("ldr d0, [x1]")
			}

			v.Generator.AddLdr("x1", "buffer")
			v.Generator.AddBl("float_to_ascii")
			v.Generator.AddInstruction("mov x2, x0")
			v.Generator.AddLdr("x1", "buffer")
			v.Generator.AddInstruction("mov x0, #1")
			v.Generator.AddInstruction("mov w8, #64")
			v.Generator.AddInstruction("svc #0")
		}
	}

	// Salto de línea al final del print
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
	if entry, ok := v.VarMap[id]; ok {

		if entry.Tipo == "float64" {
			// NO generes instrucciones aquí para floats
			return PrintValue{
				Tipo:  entry.Tipo,
				Valor: entry.Valor,
			}
		}
		// Forzar que la comparación use el valor real desde memoria
		v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", entry.Label))
		v.Generator.AddInstruction("ldr x0, [x1]")

		return PrintValue{
			Tipo:  entry.Tipo,
			Valor: "x0", // ← IMPORTANTE: se compara usando el valor cargado a x0
		}
	}
	fmt.Printf("⚠️ Variable no encontrada: %s\n", id)
	return PrintValue{Tipo: "int", Valor: "0"}
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
			switch entry.Tipo {
			case "int":
				rhs := pv.Valor

				// Cargar RHS en x2
				if _, err := strconv.Atoi(rhs); err == nil {
					v.Generator.AddMov("x2", "#"+rhs)
				} else if strings.HasPrefix(rhs, "x") {
					v.Generator.AddInstruction("mov x2, " + rhs)
				} else {
					v.Generator.AddInstruction(fmt.Sprintf("ldr x3, =%s", rhs))
					v.Generator.AddInstruction("ldr x2, [x3]")
				}

				// Cargar LHS en x0
				v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", entry.Label))
				v.Generator.AddInstruction("ldr x0, [x1]")

				// Realizar operación
				switch op {
				case "+=":
					v.Generator.AddInstruction("add x0, x0, x2")
				case "-=":
					v.Generator.AddInstruction("sub x0, x0, x2")
				}

				// Guardar el resultado
				v.Generator.AddInstruction("str x0, [x1]")

				entry.Valor = "x0"
				v.VarMap[id] = entry

			default:
				fmt.Printf("⚠️ Tipo %s no soportado aún en IMCPLICIT\n", entry.Tipo)
			}
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
	//fmt.Println("DEBUG: Entrando a AsignacionLUEGO para:", ctx.ID().GetText())
	id := ctx.ID().GetText()
	expr := ctx.Expresion()
	// Si es llamada a función
	if expr.GetRuleContext().GetRuleIndex() == parser.VlangParserRULE_funcCall {
		fmt.Println("DEBUG: Generando str x0, [x1] para", id)
		v.Visit(expr) // Esto llamará a VisitFuncCall y pondrá el valor en x0
		entry := v.VarMap[id]
		v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", entry.Label))
		v.Generator.AddInstruction("str x0, [x1]")
		entry.Valor = "x0"
		v.VarMap[id] = entry
		return nil
	}
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
    //fmt.Println("DEBUG: Entrando a VariableDeclarationImmutable para:", ctx.ID().GetText())
    id := ctx.ID().GetText()
    expr := ctx.Expresion()

    if ctx.ASSIGN() == nil {
        fmt.Println("DEBUG: No hay asignación para", id)
        return nil
    }

    // Si la variable ya existe, tratamos esto como una asignación (NO redeclaración)
    if entry, exists := v.VarMap[id]; exists {
        //fmt.Println("DEBUG: Variable ya existe, tratando como asignación:", id)
        ruleIndex := expr.GetRuleContext().GetRuleIndex()
        if ruleIndex == parser.VlangParserRULE_funcCall || ruleIndex == parser.VlangParserRULE_llamadaFuncion {
            v.Visit(expr)
            v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", entry.Label))
            v.Generator.AddInstruction("str x0, [x1]")
            entry.Valor = "x0"
            v.VarMap[id] = entry
            return nil
        }
        val := v.Visit(expr)
        if pv, ok := val.(PrintValue); ok {
            tipo := normalizeTipo(pv.Tipo)
            if tipo != entry.Tipo {
                //fmt.Printf("ERROR: Tipo incompatible en reasignación de '%s'. Esperado '%s', recibido '%s'\n", id, entry.Tipo, tipo)
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
            entry.Valor = pv.Valor
            v.VarMap[id] = entry
            return nil
        }
        return nil
    }

    // Si NO existe, es declaración nueva
    ruleIndex := expr.GetRuleContext().GetRuleIndex()
    if ruleIndex == parser.VlangParserRULE_funcCall || ruleIndex == parser.VlangParserRULE_llamadaFuncion {
        fmt.Println("DEBUG: Generando str x0, [x1] para", id)
        v.Visit(expr)
        tipo := "int" // Puedes mejorar esto si tienes el tipo real
        label := id
        v.VarMap[id] = VariableEntry{Tipo: tipo, Valor: "x0", Label: label}
        v.Generator.AddData(fmt.Sprintf(".align 3\n%s: .quad 0", label))
        v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", label))
        v.Generator.AddInstruction("str x0, [x1]")
        return nil
    }

    val := v.Visit(expr)
    if pv, ok := val.(PrintValue); ok {
        tipo := normalizeTipo(pv.Tipo)
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
	return PrintValue{Tipo: "bool", Valor: "x2"}
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
		return PrintValue{Tipo: "bool", Valor: "x2"}
	}

	// Si no es comparación válida, retorna false
	return PrintValue{Tipo: "bool", Valor: "x2"}
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
	//fmt.Println("[DEBUG] VisitOPERADORESLOGICOS: valores no son booleanos, retornando false.")
	return PrintValue{Tipo: "bool", Valor: "x2"}
}
func (v *ARMVisitor) VisitParentesisexpre(ctx *parser.ParentesisexpreContext) interface{} {
	return v.Visit(ctx.Expresion())
}

// (Opcional) si usas corchetes para grouping
func (v *ARMVisitor) VisitCorchetesexpre(ctx *parser.CorchetesexpreContext) interface{} {
	return v.Visit(ctx.Expresion())
}
func (v *ARMVisitor) VisitUnario(ctx *parser.UnarioContext) interface{} {

	valRaw := v.Visit(ctx.Expresion())
	val, ok := valRaw.(PrintValue)
	if !ok {
		fmt.Println("[ERROR] VisitUnario: valor no es PrintValue o es nil:", valRaw)
		// Puedes retornar un valor por defecto o nil, según tu lógica
		return PrintValue{Tipo: "bool", Valor: "false"}
	}
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
		//fmt.Println("[DEBUG] VisitUnario valor:", val.Valor)
		//fmt.Print("[ERROR] VisitUnario: operación ! solo válida para booleanos, retornando false.\n")
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
	//fmt.Println("Entrando a VisitIfDcl con condición:", ctx.Expresion().GetText())
	endLabel := v.Generator.GenerateUniqueLabel("ifend")
	elseIfLabels := []string{}

	for range ctx.AllElseIfDcl() {
		elseIfLabels = append(elseIfLabels, v.Generator.GenerateUniqueLabel("elseif"))
	}
	if ctx.ElseCondicional() != nil {
		elseIfLabels = append(elseIfLabels, v.Generator.GenerateUniqueLabel("else"))
	}

	val := v.Visit(ctx.Expresion()).(PrintValue)
	falseLabel := ""
	if len(elseIfLabels) > 0 {
		falseLabel = elseIfLabels[0]
	} else {
		falseLabel = endLabel
	}
	v.prepareCondBranch(val, falseLabel)

	// Bloque IF
	for _, decl := range ctx.AllDeclaraciones() {
		res := v.Visit(decl)
		if str, ok := res.(string); ok && (str == "break" || str == "continue") {
			return str
		}
		if ret, ok := res.(PrintValue); ok && ret.Tipo == "__return__" {
			return ret
		}
	}
	v.Generator.B(endLabel)

	// ELSE IF
	for i, elseIf := range ctx.AllElseIfDcl() {
		v.Generator.setLabel(elseIfLabels[i])
		val := v.Visit(elseIf.Expresion()).(PrintValue)
		nextLabel := endLabel
		if i+1 < len(elseIfLabels) {
			nextLabel = elseIfLabels[i+1]
		}
		v.prepareCondBranch(val, nextLabel)

		for _, decl := range elseIf.AllDeclaraciones() {
			res := v.Visit(decl)
			if str, ok := res.(string); ok && (str == "break" || str == "continue") {
				return str
			}
		}
		v.Generator.B(endLabel)
	}

	// ELSE
	if ctx.ElseCondicional() != nil {
		v.Generator.setLabel(elseIfLabels[len(elseIfLabels)-1])
		for _, decl := range ctx.ElseCondicional().AllDeclaraciones() {
			res := v.Visit(decl)
			if str, ok := res.(string); ok && (str == "break" || str == "continue") {
				return str
			}
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

// Holi Brandon no me mates
// Aquí están las funciones de control de flujo
/*func (v *ARMVisitor) VisitBreakStatement(ctx *parser.BreakStatementContext) interface{} {
	return "break"
}

func (v *ARMVisitor) VisitContinueStatement(ctx *parser.ContinueStatementContext) interface{} {
	return "continue"
}*/

func (v *ARMVisitor) VisitExpresionStatement(ctx *parser.ExpresionStatementContext) interface{} {
	exp := ctx.Expresion()
	ruleIndex := exp.GetRuleContext().GetRuleIndex()
	switch ruleIndex {
	case parser.VlangParserRULE_funcCall,
		parser.VlangParserRULE_llamadaFuncion:
		// Solo llama a Visit, el dispatch es automático
		return v.Visit(exp)
	default:
		return v.Visit(exp)
	}
}

func (v *ARMVisitor) VisitFor_context(ctx *parser.For_contextContext) interface{} {
	//fmt.Println("[DEBUG] Entrando a VisitFor_context")
	return v.Visit(ctx.ForDcl())
}

func (v *ARMVisitor) VisitStmt(ctx *parser.StmtContext) interface{} {
    //fmt.Println("[DEBUG] VisitStmt:", ctx.GetText())
    for i := 0; i < ctx.GetChildCount(); i++ {
        child := ctx.GetChild(i)
        if node, ok := child.(antlr.ParseTree); ok {
            res := v.Visit(node)
            if ret, ok := res.(PrintValue); ok && ret.Tipo == "__return__" {
                return ret
            }
        }
    }
    return nil
}

func (v *ARMVisitor) VisitForClasico(ctx *parser.ForClasicoContext) interface{} {
	//fmt.Println("[DEBUG] ➤ Entrando a ForClásico")

	asignacion := ctx.Asignacion()
	condicion := ctx.Expresion()
	bloque := ctx.Block()

	id := asignacion.ID().GetText()
	valor := v.Visit(asignacion.Expresion())
	tipo := "int"

	if _, exists := v.VarMap[id]; !exists {
		v.VarMap[id] = VariableEntry{Tipo: tipo, Label: id, Valor: "0"}
		v.Generator.AddData(fmt.Sprintf(".align 3\n%s: .quad 0", id))
	}

	if pv, ok := valor.(PrintValue); ok {
		v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", id))
		v.Generator.AddMov("x0", pv.Valor)
		v.Generator.AddInstruction("str x0, [x1]")
		entry := v.VarMap[id]
		entry.Valor = "x0"
		v.VarMap[id] = entry
	}

	labelInicio := v.Generator.GenerateUniqueLabel("for_start")
	labelCond := v.Generator.GenerateUniqueLabel("for_cond")
	labelFin := v.Generator.GenerateUniqueLabel("for_end")
	labelContinue := v.Generator.GenerateUniqueLabel("for_continue")

	v.breakLabels = append(v.breakLabels, labelFin)
	v.continueLabels = append(v.continueLabels, labelContinue)

	v.Generator.B(labelCond)
	v.Generator.setLabel(labelInicio)

	for _, decl := range bloque.AllDeclaraciones() {
		res := v.Visit(decl)
		if str, ok := res.(string); ok {
			if str == "break" {
				v.Generator.DebugPrint("dbg_break", "→ [DEBUG] Se ejecutó un BREAK\\n")
				v.Generator.B(labelFin)
				break
			}
			if str == "continue" {
				v.Generator.DebugPrint("dbg_continue", "→ [DEBUG] Se ejecutó un CONTINUE\\n")
				v.Generator.B(labelContinue)
				break
			}
		}
	}

	v.Generator.setLabel(labelContinue)
	v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", id))
	v.Generator.AddInstruction("ldr x0, [x1]")
	v.Generator.AddInstruction("add x0, x0, #1")
	v.Generator.AddInstruction("str x0, [x1]")

	v.Generator.setLabel(labelCond)
	val := v.Visit(condicion)
	if pv, ok := val.(PrintValue); ok {
		v.prepareCondBranch(pv, labelFin)
	}
	v.Generator.B(labelInicio)

	v.Generator.setLabel(labelFin)

	v.breakLabels = v.breakLabels[:len(v.breakLabels)-1]
	v.continueLabels = v.continueLabels[:len(v.continueLabels)-1]
	delete(v.VarMap, id)

	return nil
}

// Implementación para FOR tipo while: for condicion { cuerpo }
func (v *ARMVisitor) VisitForCondicionUnica(ctx *parser.ForCondicionUnicaContext) interface{} {
	//fmt.Println("[DEBUG] Entrando a ForCondicionUnica")

	startLabel := v.Generator.GenerateUniqueLabel("for_start")
	condLabel := v.Generator.GenerateUniqueLabel("for_cond")
	endLabel := v.Generator.GenerateUniqueLabel("for_end")
	continueLabel := v.Generator.GenerateUniqueLabel("for_continue")

	// Saltar a evaluar la condición primero
	v.Generator.B(condLabel)

	// Etiqueta de inicio del cuerpo
	v.Generator.setLabel(startLabel)

	// Visitar declaraciones dentro del cuerpo del ciclo
	for _, decl := range ctx.Block().AllDeclaraciones() {
		switch {
		case decl.Stmt() != nil:
			v.Visit(decl.Stmt())
		case decl.VarDcl() != nil:
			v.Visit(decl.VarDcl())
		case decl.FuncDcl() != nil:
			v.Visit(decl.FuncDcl())
		case decl.StructDcl() != nil:
			v.Visit(decl.StructDcl())
		case decl.FuncMain() != nil:
			v.Visit(decl.FuncMain())
		default:
			v.Visit(decl)
		}
	}

	// Etiqueta opcional para continue
	v.Generator.setLabel(continueLabel)

	// Evaluar condición del ciclo
	v.Generator.setLabel(condLabel)

	val := v.Visit(ctx.Expresion())
	pv, ok := val.(PrintValue)
	if !ok {
		fmt.Println("[ERROR] Condición del for no es PrintValue")
		return nil
	}

	// Generar salto si la condición ya no se cumple
	v.prepareCondBranch(pv, endLabel)

	// Volver al cuerpo
	v.Generator.B(startLabel)

	// Fin del ciclo
	v.Generator.setLabel(endLabel)

	//fmt.Println("[DEBUG] Salida de ForCondicionUnica")
	return nil
}

func (v *ARMVisitor) VisitBreakStatement(ctx *parser.BreakStatementContext) interface{} {
	v.Generator.DebugPrint("dbg_break", "→ [DEBUG] Se ejecutó BREAK\n")
	if len(v.breakLabels) > 0 {
		v.Generator.AddInstruction("b " + v.breakLabels[len(v.breakLabels)-1])
	}
	return "break"
}

func (v *ARMVisitor) VisitContinueStatement(ctx *parser.ContinueStatementContext) interface{} {
	v.Generator.DebugPrint("dbg_continue", "→ [DEBUG] Se ejecutó CONTINUE\n")
	if len(v.continueLabels) > 0 {
		v.Generator.AddInstruction("b " + v.continueLabels[len(v.continueLabels)-1])
	}
	return "continue"
}

func (v *ARMVisitor) VisitIncremento(ctx *parser.IncrementoContext) interface{} {
	id := ctx.ID().GetText()
	entry, ok := v.VarMap[id]
	if !ok {
		fmt.Printf("⚠️ Variable no encontrada en incremento: %s\n", id)
		return nil
	}

	if entry.Tipo != "int" {
		fmt.Printf("⚠️ Incremento no soportado para tipo: %s\n", entry.Tipo)
		return nil
	}

	v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", entry.Label))
	v.Generator.AddInstruction("ldr x0, [x1]")
	v.Generator.AddInstruction("add x0, x0, #1")
	v.Generator.AddInstruction("str x0, [x1]")

	entry.Valor = "x0"
	v.VarMap[id] = entry
	return nil
}

func (v *ARMVisitor) VisitDecremento(ctx *parser.DecrementoContext) interface{} {
	id := ctx.ID().GetText()
	entry, ok := v.VarMap[id]
	if !ok {
		fmt.Printf("⚠️ Variable no encontrada en decremento: %s\n", id)
		return nil
	}

	if entry.Tipo != "int" {
		fmt.Printf("⚠️ Decremento no soportado para tipo: %s\n", entry.Tipo)
		return nil
	}

	// x = x - 1
	v.Generator.AddInstruction(fmt.Sprintf("ldr x1, =%s", entry.Label))
	v.Generator.AddInstruction("ldr x0, [x1]")
	v.Generator.AddInstruction("sub x0, x0, #1")
	v.Generator.AddInstruction("str x0, [x1]")

	entry.Valor = "x0"
	v.VarMap[id] = entry
	return nil
}

// aqui seran las funciones de las estructuras
func NewARMVisitor(generator *ArmGenerator, scope *ScopeTrace) *ARMVisitor {
	return &ARMVisitor{
		Generator:  generator,
		ScopeTrace: scope,
		FuncLabels: make(map[string]string),
	}
}

func (v *ARMVisitor) VisitFuncDcl(ctx *parser.FuncDclContext) interface{} {
	funcName := ctx.ID().GetText()
	label := "fn_" + funcName
	v.FuncLabels[funcName] = label
    //fmt.Println("[DEBUG] VisitFuncDcl:", funcName)
	// ① Definir la etiqueta de la función en el .s
	v.Generator.setLabel(label)

	// ② Extraer parámetros (hasta 4 en registros x0–x3)
	paramNames := []string{}
	if ctx.ParametrosFormales() != nil {
		for _, param := range ctx.ParametrosFormales().AllParametro() {
			paramNames = append(paramNames, param.ID().GetText())
		}
	}

	// ③ Guardar cada parámetro en memoria usando su etiqueta
	for i, name := range paramNames {
		// Sube x{i} a la etiqueta =name
		v.Generator.AddInstruction(fmt.Sprintf("str x%d, =%s", i, name))
	}

	// ④ Generar etiqueta de salida para los return
	endLabel := v.Generator.GenerateUniqueLabel("fn_end")

	// ⑤ Recorrer el cuerpo y detectar returns para saltar al endLabel
	for _, decl := range ctx.Block().AllDeclaraciones() {
		//fmt.Println("[DEBUG] VisitFuncDcl visit decl:", decl.GetText())
		res := v.Visit(decl)
		if ret, ok := res.(PrintValue); ok && ret.Tipo == "__return__" {
			//v.Generator.B(endLabel)
			v.Generator.AddInstruction("ret")
			//fmt.Println("[DEBUG] VisitFuncDcl detectó return en", funcName)
			return ret
		}
	}

	// ⑥ Marcar la etiqueta de salida y emitir el ret
	v.Generator.setLabel(endLabel)
	v.Generator.AddInstruction("ret")
	return nil
}
func (v *ARMVisitor) VisitLlamadaFuncionExpr(ctx *parser.LlamadaFuncionExprContext) interface{} {
	return v.Visit(ctx.LlamadaFuncion())
}


func (v *ARMVisitor) VisitLlamadaFuncion(ctx *parser.LlamadaFuncionContext) interface{} {
	// Solo soportamos ID LPAREN ... RPAREN por ahora
	funcName := ctx.ID().GetText()
	label, ok := v.FuncLabels[funcName]
	if !ok {
		fmt.Printf("Función '%s' no declarada\n", funcName)
		return nil
	}
	// Evalúa argumentos si tienes soporte
	// ...
	v.Generator.AddBl(label)
	return nil
}
func (v *ARMVisitor) VisitFuncCall(ctx *parser.FuncCallContext) interface{} {
	//fmt.Println("[DEBUG] VisitFuncCall:", ctx.GetText())
	funcName := ctx.ID().GetText()
	label, ok := v.FuncLabels[funcName]
	if !ok {
		//fmt.Printf("Función '%s' no declarada\n", funcName)
		return nil
	}

	v.Generator.DebugPrint("call_"+funcName, "Entrando a "+funcName+"\\n")

	// Evalúa argumentos y colócalos en x0, x1, ...
	if ctx.ParametrosReales() != nil {
		for i, expr := range ctx.ParametrosReales().AllExpresion() {
			val := v.Visit(expr)
			if pv, ok := val.(PrintValue); ok {
				v.Generator.AddMov(fmt.Sprintf("x%d", i), pv.Valor)
			}
		}
	}

	// Llama a la función
	//fmt.Println("[DEBUG] Llamando a la función:", funcName, "SU LABEL ES:", label)
	v.Generator.AddBl(label)
	v.Generator.DebugPrint("ret_"+funcName, "Saliendo de "+funcName+"\\n")
	// El valor de retorno está en x0
	return PrintValue{Tipo: "int", Valor: "x0"} // O el tipo real
}

func (v *ARMVisitor) CallFunctionByName(name string) {
	label, ok := v.FuncLabels[name]
	if !ok {
		fmt.Printf("Función '%s' no declarada\n", name)
		return
	}
	v.Generator.AddBl(label)
}
func (v *ARMVisitor) VisitReturnStatement(ctx *parser.ReturnStatementContext) interface{} {
    txt := ctx.GetText()
    //fmt.Println("[DEBUG] VisitReturnStatement:", txt)
    // Si el texto viene como "return42", sepáralo
    if strings.HasPrefix(txt, "return") && len(txt) > 6 {
           valor := txt[6:] // lo que sigue después de "return"
           //fmt.Println("[DEBUG] Valor de retorno (parche):", valor)
           v.Generator.AddMov("x0", valor)
           return PrintValue{Tipo: "__return__"}
    }
    // Si el árbol está bien, usa el hijo expresión
    if ctx.Expresion() != nil {
        val := v.Visit(ctx.Expresion())
        if pv, ok := val.(PrintValue); ok {
            //fmt.Println("[DEBUG] Valor de retorno:", pv.Valor)
            v.Generator.AddMov("x0", pv.Valor)
            return PrintValue{Tipo: "__return__"}
        }
    }
    //fmt.Println("[DEBUG] VisitReturnStatement: no se pudo obtener un valor válido, retornando __return__")
    return PrintValue{Tipo: "__return__"}
}

func (g *ArmGenerator) DebugPrint(label string, msg string) {
	dataLabel := "dbg_" + label
	lenLabel := "len_" + dataLabel
	g.AddData(fmt.Sprintf("%s: .ascii \"%s\"", dataLabel, msg))
	g.AddData(fmt.Sprintf("%s: .quad . - %s", lenLabel, dataLabel))
	g.AddInstruction("mov x0, #1") // <-- aquí ya no uses fmt.Sprintf
	g.AddInstruction(fmt.Sprintf("ldr x1, =%s", dataLabel))
	g.AddInstruction(fmt.Sprintf("ldr x2, =%s", lenLabel))
	g.AddInstruction("ldr x2, [x2]")
	g.AddInstruction("mov w8, #64")
	g.AddInstruction("svc #0")
}
