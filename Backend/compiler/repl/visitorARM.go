package repl

import (
	parser "compiler/parser"
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

type VariableEntry struct {
	Tipo  string
	Label string
	Valor string
}

type ARMVisitor struct {
	parser.BaseVlangVisitor
	Generator  *ArmGenerator
	PrintCount int
	VarMap     map[string]VariableEntry
	ScopeTrace *ScopeTrace
	UsesIntToAscii bool
}

type PrintValue struct {
	Tipo  string
	Valor string
}
//debido a que dicho arbol puede que le sea dificil manejarlo
//ahora de igual forma leeremos los scopes donde estan los entornos de las variables
func (v *ARMVisitor) CollectAllVariables(scope *BaseScope) {
    for name, variable := range scope.variables {
        valStr := fmt.Sprint(variable.Value.Value())
        v.VarMap[name] = VariableEntry{
            Tipo:  variable.Type,
            Label: "msg_" + name,
            Valor: valStr,
        }
        v.Generator.AddData(fmt.Sprintf("msg_%s: .ascii \"%s\\n\"", name, valStr))
        v.Generator.AddData(fmt.Sprintf("len_%s: .quad . - msg_%s", name, name))
    }
    // Recursivamente para los hijos
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
    v.Generator.AddData(`msg_true: .ascii "true\n"`)
    v.Generator.AddData(`len_true: .quad . - msg_true`)
    v.Generator.AddData(`msg_false: .ascii "false\n"`)
    v.Generator.AddData(`len_false: .quad . - msg_false`)

    for _, decl := range ctx.AllDeclaraciones() {
        v.Visit(decl)
    }
	
    if v.Generator != nil && v.VarMap != nil && v.ScopeTrace != nil {
        v.CollectAllVariables(v.ScopeTrace.GlobalScope)
    }
	if v.UsesIntToAscii {
		v.Generator.AddIntToAsciiFunction()
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
	tipo := "int" // valor por defecto si no hay tipo

	if ctx.TIPO() != nil {
		tipo = ctx.TIPO().GetText()
	}

	var valor string
	if ctx.Expresion() != nil {
		val := v.Visit(ctx.Expresion())
		if pv, ok := val.(PrintValue); ok {
			valor = pv.Valor
		}
	} else {
		switch tipo {
		case "int":
			valor = "0"
		case "string":
			valor = ""
		case "bool":
			valor = "false"
		case "rune":
			valor = "?"
		case "float64":
			valor = "0.0"
		default:
			valor = "(undef)"
		}
	}

	label := fmt.Sprintf("msg_%s", id)
	v.VarMap[id] = VariableEntry{Tipo: tipo, Label: label, Valor: valor}
	v.Generator.AddData(fmt.Sprintf("%s: .ascii \"%s\\n\"", label, valor))
	v.Generator.AddData(fmt.Sprintf("len_%s: .quad . - %s", id, label))
	return nil
}

func (v *ARMVisitor) VisitPrintStatement(ctx *parser.PrintStatementContext) interface{} {
    for _, expr := range ctx.AllExpresion() {
        val := v.Visit(expr)

        // Si es un identificador (variable)
        if id, ok := val.(string); ok {
            if entry, exists := v.VarMap[id]; exists {
                switch entry.Tipo {
                case "int":
                    v.UsesIntToAscii = true
    				v.Generator.AddInstruction(fmt.Sprintf("mov x0, #%s", entry.Valor))
					v.Generator.AddInstruction("ldr x1, =buffer")
					v.Generator.AddInstruction("bl int_to_ascii")
					v.Generator.AddInstruction("mov x2, x0") // longitud
					v.Generator.AddInstruction("mov x1, x1") // puntero al string
					v.Generator.AddInstruction("mov x0, #1") // stdout
					v.Generator.AddInstruction("mov w8, #64")
					v.Generator.AddInstruction("svc #0")
                    // Imprimir salto de línea
                    v.Generator.AddInstruction("mov x0, #1")
                    v.Generator.AddInstruction("ldr x1, =msg_nl")
                    v.Generator.AddInstruction("ldr x2, =len_nl")
                    v.Generator.AddInstruction("ldr x2, [x2]")
                    v.Generator.AddInstruction("mov w8, #64")
                    v.Generator.AddInstruction("svc #0")
                case "bool":
                    boolLabel := "msg_true"
                    lenLabel := "len_true"
                    if entry.Valor == "false" {
                        boolLabel = "msg_false"
                        lenLabel = "len_false"
                    }
                    v.printSyscall(boolLabel, lenLabel)
                default:
                    v.printSyscall(entry.Label, fmt.Sprintf("len_%s", id))
                }
            }
        } else if pv, ok := val.(PrintValue); ok {
            switch pv.Tipo {
            case "entero":
                v.UsesIntToAscii = true
				v.Generator.AddMov("x0", "#"+pv.Valor)
				v.Generator.AddLdr("x1", "buffer")
				v.Generator.AddBl("int_to_ascii")
				v.Generator.AddInstruction("mov x2, x0")
				v.Generator.AddInstruction("mov x1, x1")
				v.Generator.AddInstruction("mov x0, #1")
				v.Generator.AddInstruction("mov w8, #64")
				v.Generator.AddInstruction("svc #0")
                // salto de línea
                v.Generator.AddMov("x0", "#1")
                v.Generator.AddLdr("x1", "msg_nl")
                v.Generator.AddLdr("x2", "len_nl")
                v.Generator.AddInstruction("ldr x2, [x2]")
                v.Generator.AddMov("w8", "#64")
                v.Generator.AddSvc()
            case "booleano":
                boolLabel := "msg_true"
                lenLabel := "len_true"
                if pv.Valor == "false" {
                    boolLabel = "msg_false"
                    lenLabel = "len_false"
                }
                v.printSyscall(boolLabel, lenLabel)
            case "cadena", "caracter":
                v.PrintCount++
                label := fmt.Sprintf("msg%d", v.PrintCount)
                lenLabel := fmt.Sprintf("len%d", v.PrintCount)
                v.Generator.AddData(fmt.Sprintf("%s: .ascii \"%s\\n\"", label, pv.Valor))
                v.Generator.AddData(fmt.Sprintf("%s: .quad . - %s", lenLabel, label))
                v.printSyscall(label, lenLabel)
				fmt.Println("Generando etiqueta:", label)
            }
        }
    }
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
	return PrintValue{Tipo: "booleano", Valor: ctx.GetText()}
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
	label := fmt.Sprintf("msg_%s", id)
	lenLabel := fmt.Sprintf("len_%s", id)

	// Actualizar también el valor en VarMap
	entry := v.VarMap[id]
	entry.Valor = newValue
	v.VarMap[id] = entry

	newData := []string{}
	skip := false

	for _, line := range v.Generator.Data {
		if skip {
			skip = false
			continue
		}
		if len(line) > 0 && line[:len(label)+1] == label+":" {
			skip = true
			continue
		}
		newData = append(newData, line)
	}
	v.Generator.Data = newData

	v.Generator.AddData(fmt.Sprintf("%s: .ascii \"%s\\n\"", label, newValue))
	v.Generator.AddData(fmt.Sprintf("%s: .quad . - %s", lenLabel, label))
}

func (v *ARMVisitor) VisitAsignacionLUEGO(ctx *parser.AsignacionLUEGOContext) interface{} {
	id := ctx.ID().GetText()
	tipo := ctx.TIPO().GetText()
	val := v.Visit(ctx.Expresion())

	if pv, ok := val.(PrintValue); ok {
		if entry, exists := v.VarMap[id]; exists {
			entry.Tipo = tipo
			v.replaceVarData(id, pv.Valor) // <-- Esto actualiza también el valor internamente
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
		return PrintValue{Tipo: "float64", Valor: formatNumber(toFloat(l.Valor) + toFloat(r.Valor))}
	case "-":
		return PrintValue{Tipo: "float64", Valor: formatNumber(toFloat(l.Valor) + toFloat(r.Valor))}
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

func toFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}

func toInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
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
