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
	Generator  *ARMGenerator
	PrintCount int
	VarMap     map[string]VariableEntry
}

type PrintValue struct {
	Tipo  string
	Valor string
}

func (v *ARMVisitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.Accept(v)
}

func (v *ARMVisitor) VisitPrograma(ctx *parser.ProgramaContext) interface{} {
	v.VarMap = make(map[string]VariableEntry)
	for _, decl := range ctx.AllDeclaraciones() {
		v.Visit(decl)
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
	var outputLine string

	for i, expr := range ctx.AllExpresion() {
		val := v.Visit(expr)

		if id, ok := val.(string); ok {
			if entry, exists := v.VarMap[id]; exists {
				outputLine += entry.Valor
			}
		} else if pv, ok := val.(PrintValue); ok {
			outputLine += pv.Valor
		}

		if i < len(ctx.AllExpresion())-1 {
			outputLine += " " // espacio entre argumentos
		}
	}

	v.PrintCount++
	label := fmt.Sprintf("msg%d", v.PrintCount)
	lenLabel := fmt.Sprintf("len%d", v.PrintCount)

	v.Generator.AddData(fmt.Sprintf("%s: .ascii \"%s\\n\"", label, outputLine))
	v.Generator.AddData(fmt.Sprintf("%s: .quad . - %s", lenLabel, label))
	v.Generator.AddInstruction(fmt.Sprintf("    # Print salida %d", v.PrintCount))
	v.printSyscall(label, lenLabel)

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
	v.Generator.AddInstruction("mov x0, #1")
	v.Generator.AddInstruction(fmt.Sprintf("adr x1, %s", label))
	v.Generator.AddInstruction(fmt.Sprintf("ldr x2, =%s", lenLabel))
	v.Generator.AddInstruction("ldr x2, [x2]")
	v.Generator.AddInstruction("mov w8, #64")
	v.Generator.AddInstruction("svc #0")
}
