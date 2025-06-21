package repl

import (
	parser "compiler/parser"
	"fmt"

	"github.com/antlr4-go/antlr/v4"
)

// Visitor personalizado para recorrer el árbol de sintaxis

// Constructor del visitor
type ARMVisitor struct {
	parser.BaseVlangVisitor
	ScopeTrace  *ScopeTrace
	//CallStack   *CallStack
	Console     *Console
	StructNames []string
	Generator  *ARMGenerator 
	PrintCount  int
}

type PrintValue struct {
    Tipo  string // "entero", "decimal", "cadena", "booleano", "caracter"
    Valor string
}
var _ parser.VlangVisitor = &ARMVisitor{}


func (v *ARMVisitor) Visit(tree antlr.ParseTree) interface{} {
    return tree.Accept(v)
}

func (v *ARMVisitor) VisitPrintStatement(ctx *parser.PrintStatementContext) interface{} {
    for _, expr := range ctx.AllExpresion() {
        v.PrintCount++
        label := fmt.Sprintf("msg%d", v.PrintCount)
        lenLabel := fmt.Sprintf("len%d", v.PrintCount)

        val := v.Visit(expr)
        pv, ok := val.(PrintValue)
        if !ok {
            continue
        }
		strVal := pv.Valor
        switch pv.Tipo {
        case "booleano":
            // Imprime "True" o "False" (puedes cambiarlo a minúsculas si prefieres)
            if strVal == "true" {
                strVal = "True"
            } else {
                strVal = "False"
            }
        case "caracter":
            // Ya está bien como string
        case "entero":
            // Ya está bien como string
        case "decimal":
            // Ya está bien como string
        case "cadena":
            // Ya está bien como string
        }
         





        // Asegura salto de línea
        if len(strVal) == 0 || strVal[len(strVal)-1] != '\n' {
            strVal += "\\n"
        }

        v.Generator.AddInstruction(fmt.Sprintf("    # Print salida %d", v.PrintCount))
        v.Generator.AddInstruction("mov x0, #1")
        v.Generator.AddInstruction(fmt.Sprintf("adr x1, %s", label))
        v.Generator.AddInstruction(fmt.Sprintf("ldr x2, =%s", lenLabel))
        v.Generator.AddInstruction("ldr x2, [x2]")
        v.Generator.AddInstruction("mov w8, #64")
        v.Generator.AddInstruction("svc #0")

        // Guarda los datos para la sección .rodata
        v.Generator.AddData(fmt.Sprintf("%s: .ascii \"%s\"", label, strVal))
        v.Generator.AddData(fmt.Sprintf("%s: .quad . - %s", lenLabel, label))
    }
    return nil
}


func (v *ARMVisitor) VisitValorexpr(ctx *parser.ValorexprContext) interface{} {
    return v.Visit(ctx.Valor())
}

func (v *ARMVisitor) VisitPrograma(ctx *parser.ProgramaContext) interface{} {
    for _, decl := range ctx.AllDeclaraciones() {
        v.Visit(decl)
    }
    // Solo una vez al final
    v.Generator.AddInstruction("    # Salida final")
    v.Generator.AddInstruction("mov x0, #0")
    v.Generator.AddInstruction("mov w8, #93")
    v.Generator.AddInstruction("svc #0")
    return nil
}

func (v *ARMVisitor) VisitDeclaraciones(ctx *parser.DeclaracionesContext) interface{} {
    fmt.Println("VisitDeclaraciones")
    if ctx.Stmt() != nil {
        return v.Visit(ctx.Stmt())
    }
    return nil
}

func (v *ARMVisitor) VisitStmt(ctx *parser.StmtContext) interface{} {
    fmt.Println("VisitStmt")
    return v.VisitChildren(ctx)
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