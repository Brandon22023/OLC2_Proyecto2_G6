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
}
var _ parser.VlangVisitor = &ARMVisitor{}


func (v *ARMVisitor) Visit(tree antlr.ParseTree) interface{} {
    return tree.Accept(v)
}

func (v *ARMVisitor) VisitPrintStatement(ctx *parser.PrintStatementContext) interface{} {
    for _, expr := range ctx.AllExpresion() {
        val := v.Visit(expr)
        strVal, ok := val.(string)
        if !ok {
            continue
        }
		//hace que siempre de un salto de linea 
		if len(strVal) == 0 || strVal[len(strVal)-1] != '\n' {
            strVal += "\\n"
        }

        v.Generator.AddInstruction("mov x0, #1") // stdout
		v.Generator.AddInstruction("adr x1, msg") // dirección del mensaje
		v.Generator.AddInstruction("ldr x2, =len") // dirección de la longitud
		v.Generator.AddInstruction("ldr x2, [x2]") // cargar el valor de la longitud
		v.Generator.AddInstruction("mov w8, #64") // syscall write
		v.Generator.AddInstruction("svc #0")      // llamada al sistema

		v.Generator.AddInstruction("mov x0, #0")  // código de salida
		v.Generator.AddInstruction("mov w8, #93") // syscall exit
		v.Generator.AddInstruction("svc #0")

		v.Generator.AddData("msg: .ascii \"" + strVal + "\"")
		v.Generator.AddData("len: .quad . - msg")
    }
    return nil
}

func (v *ARMVisitor) VisitValorCadena(ctx *parser.ValorCadenaContext) interface{} {
    // Elimina las comillas del string
	fmt.Println("Entrando a VisitValorCadena")
    text := ctx.GetText()
    if len(text) >= 2 && text[0] == '"' && text[len(text)-1] == '"' {
        return text[1 : len(text)-1]
    }
    return text
}

func (v *ARMVisitor) VisitValorexpr(ctx *parser.ValorexprContext) interface{} {
    return v.Visit(ctx.Valor())
}

func (v *ARMVisitor) VisitPrograma(ctx *parser.ProgramaContext) interface{} {
    fmt.Println("VisitPrograma")
    for _, decl := range ctx.AllDeclaraciones() {
        v.Visit(decl)
    }
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