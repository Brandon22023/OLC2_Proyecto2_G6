package main

import (
	"fmt"
	"os"
	"path/filepath"

	"compiler/errors"
	compiler "compiler/parser"
	repl "compiler/repl"
	"compiler/symbols"

	"github.com/antlr4-go/antlr/v4"
)

func main() {
	// Leer código fuente desde archivo
	fmt.Println("===> Leyendo archivo de entrada...")
	inputCode, err := readStdin()
	if err != nil {
		fmt.Println("Error leyendo entrada:", err)
		return
	}

	// === Análisis léxico
	lexicalErrorListener := errors.NewLexicalErrorListener()
	lexer := compiler.NewVlangLexer(antlr.NewInputStream(inputCode))
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(lexicalErrorListener)

	// === Tokens
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// === Parser
	parser := compiler.NewVlangParser(stream)
	parser.BuildParseTrees = true
	syntaxErrorListener := errors.NewSyntaxErrorListener(lexicalErrorListener.ErrorTable)
	parser.RemoveErrorListeners()
	parser.SetErrorHandler(errors.NewCustomErrorStrategy())
	parser.AddErrorListener(syntaxErrorListener)

	fmt.Println("===> Se terminó el análisis sintáctico")
	arbolito := parser.Programa()

	// === Imprimir árbol sintáctico (opcional)
	PrintVerticalTree(arbolito, parser.RuleNames)

	// === Recorre con visitor normal (para tabla de símbolos)
	visitor := repl.NewReplVisitor()
	visitor.Visit(arbolito)

	fmt.Println("==== Variables en los entornos ====")
	visitor.ScopeTrace.GlobalScope.PrintScopeVariables(0)

	// === Crear tabla de símbolos
	tabla := symbols.NewSymbolTable()
	visitor.ScopeTrace.GlobalScope.CollectSymbols(tabla)

	fmt.Println("==== Tabla de símbolos ====")
	for _, sym := range tabla.Symbols {
		fmt.Printf("ID: %s | Tipo: %s | Dato: %s | Ámbito: %s | Línea: %d | Columna: %d\n",
			sym.ID, sym.SymType, sym.DataType, sym.Scope, sym.Line, sym.Column)
	}

	err = tabla.ToHTML("tabla_simbolos.html")
	if err != nil {
		fmt.Println("Error generando HTML de tabla de símbolos:", err)
	} else {
		fmt.Println("Tabla de símbolos exportada a tabla_simbolos.html")
	}

	// === Generación de código Assembler
	fmt.Println("===> Generando código ARM assembler...")
	generator := &repl.ARMGenerator{}
	armVisitor := &repl.ARMVisitor{
		Generator: generator,
	}
	armVisitor.Visit(arbolito)

	asmCode := generator.ToString()
	assemblerDir := "assembler"
	os.MkdirAll(assemblerDir, 0755)
	asmPath := filepath.Join(assemblerDir, "assemblerprueba.s")
	err = os.WriteFile(asmPath, []byte(asmCode), 0644)
	if err != nil {
		fmt.Println("Error al escribir el archivo Assembly:", err)
	} else {
		fmt.Println("Archivo Assembly generado correctamente en:", asmPath)
	}
}

/*
	func readStdin() (string, error) {
		input, err := os.ReadFile("/dev/stdin")
		return string(input), err
	}
*/
func readStdin() (string, error) {
	//input, err := os.ReadFile("/home/brandon/Escritorio/OLC2_Proyecto1_202300813/Backend/compiler/arhivoP.vch")
	//input, err := os.ReadFile("/home/brandon/Escritorio/OLC2_Proyecto1_202300813/Backend/compiler/basicas.vch")
	//input, err := os.ReadFile("/home/pablo/Escritorio/OLC2_Proyecto1_202300813/Backend/compiler/arhivoP.vch")
	input, err := os.ReadFile("/home/vboxuser/Documents/OLC2_Proyecto2_G6/Backend/compiler/arhivoP.vch")
	return string(input), err
}

// Funciones para visualizar nuestro árbol
func PrintVerticalTree(node antlr.Tree, ruleNames []string) {
	printVerticalNode(node, ruleNames, "", true)
}

func printVerticalNode(node antlr.Tree, ruleNames []string, prefix string, isLast bool) {
	connector := "+-- "
	if !isLast {
		connector = "|-- "
	}

	var label string
	switch n := node.(type) {
	case antlr.RuleNode:
		label = ruleNames[n.GetRuleContext().GetRuleIndex()]
	case antlr.TerminalNode:
		label = fmt.Sprintf("\"%s\"", n.GetText())
	default:
		label = fmt.Sprintf("%T", n)
	}

	fmt.Printf("%s%s%s\n", prefix, connector, label)

	childCount := node.GetChildCount()
	for i := 0; i < childCount; i++ {
		child := node.GetChild(i)
		newPrefix := prefix
		if isLast {
			newPrefix += "    "
		} else {
			newPrefix += "|   "
		}
		printVerticalNode(child, ruleNames, newPrefix, i == childCount-1)
	}
}

/*package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"compiler/cst"
	compiler "compiler/parser"
	repl "compiler/repl"

	"compiler/symbols"

	//"main/cst"
	"compiler/errors"
	// "main/repl"

	"github.com/antlr4-go/antlr/v4"
	"github.com/gofiber/fiber/v2"
)

func main() {
	var absRuta string
    app := fiber.New()

    app.Post("/analizar", func(c *fiber.Ctx) error {
        inputCode := string(c.Body())

		var buf bytes.Buffer
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

        // 1. Análisis Léxico
        lexicalErrorListener := errors.NewLexicalErrorListener()
        lexer := compiler.NewVlangLexer(antlr.NewInputStream(inputCode))
        lexer.RemoveErrorListeners()
        lexer.AddErrorListener(lexicalErrorListener)

        // 2. Tokens
        stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

        // 3. Parser + errores sintácticos
        parser := compiler.NewVlangParser(stream)
        parser.BuildParseTrees = true
        syntaxErrorListener := errors.NewSyntaxErrorListener(lexicalErrorListener.ErrorTable)
        parser.RemoveErrorListeners()
        parser.SetErrorHandler(errors.NewCustomErrorStrategy())
        parser.AddErrorListener(syntaxErrorListener)

        arbolito := parser.Programa()
		// Imprime el árbol sintáctico para depuración
		//PrintVerticalTree(arbolito, parser.RuleNames)
		//----------------arm

		// 1. Crear el generador y el visitor ARM
		generator := &repl.ARMGenerator{}
		armVisitor := &repl.ARMVisitor{
		Generator: generator,
			// Puedes inicializar otros campos si lo necesitas
		}

		// 2. Recorrer el árbol con el visitor ARM
		armVisitor.Visit(arbolito)

		asmCode := generator.ToString()
		assemblerDir := "assembler"
		os.MkdirAll(assemblerDir, 0755) // Crea la carpeta si no existe
		asmPath := filepath.Join(assemblerDir, "program.s")
		err := os.WriteFile(asmPath, []byte(asmCode), 0644)
		if err != nil {
			fmt.Println("Error al escribir el archivo Assembly:", err)
		} else {
			fmt.Println("Archivo Assembly generado correctamente en:", asmPath)
		}
		//----------------------
        visitor := repl.NewReplVisitor()

        visitor.Visit(arbolito)

		//fmt.Println("====aqui imprimire las variables de todos los entornos=========")
		//visitor.ScopeTrace.GlobalScope.PrintScopeVariables(0)

        // Cierra y recupera la salida
		w.Close()
		os.Stdout = old
		buf.ReadFrom(r)
		output := buf.String()
		// 1. Crear la tabla y llenarla
		tabla := symbols.NewSymbolTable()
		// Supón que tienes acceso al scope global o raíz:
		if visitor != nil && visitor.ScopeTrace != nil {
			visitor.ScopeTrace.GlobalScope.CollectSymbols(tabla)
		}

		// 2. Generar el HTML
		ruta := filepath.Join("reportes", "tabla_simbolos.html")
		// Crear la carpeta si no existe
		os.MkdirAll("reportes", 0755)
		err = tabla.ToHTML(ruta)
		if err != nil {
			fmt.Println("Error al crear el HTML:", err)
		} else {
			absRuta, _ := filepath.Abs(ruta)
			fmt.Println("Tabla de símbolos generada correctamente en:", absRuta)

		}
		fmt.Println("errores llegados: ", lexicalErrorListener.ErrorTable.Errors, syntaxErrorListener.ErrorTable.Errors)
		allErrors := append(lexicalErrorListener.ErrorTable.Errors, syntaxErrorListener.ErrorTable.Errors...)
		allErrors = append(allErrors, visitor.SemanticErrors.Errors...)
		allErrors = errors.RemoveDuplicateErrors(allErrors)

		errorTable := &repl.ErrorTable{Errors: allErrors}
		rutaErrores := filepath.Join("reportes", "errores.html")

		os.MkdirAll("reportes", 0755)
		_ = errors.SaveErrorsHTML(errorTable, rutaErrores)



		// Si no hay salida, muestra un mensaje por defecto
		if output == "" {
			output = "No hubo salida del programa."
		}
		os.WriteFile("ultimo_codigo.vch", []byte(inputCode), 0644)



        return c.SendString(output)
    })

	app.Get("/reporte-simbolos", func(c *fiber.Ctx) error {
			ruta := filepath.Join("reportes", "tabla_simbolos.html")
			fmt.Println("Ruta del reporte que se envia al frontend:", absRuta)
			absRuta, err := filepath.Abs(ruta)
			if err != nil {
				return c.Status(500).SendString("No se pudo obtener la ruta del reporte")
			}

    	    // Abre el HTML en el navegador del servidor
    		go symbols.OpenHTML(absRuta)
			return c.SendString(absRuta)
	})

	app.Get("/reporte-cst", func(c *fiber.Ctx) error {
		// Lee el último código analizado (puedes guardarlo en una variable global si lo necesitas)
		input, err := os.ReadFile("ultimo_codigo.vch") // O usa el input que desees
		if err != nil {
			return c.Status(500).SendString("No se pudo leer el código fuente")
		}
		svg := cst.CstReport(string(input))
		ruta := filepath.Join("reportes", "arbol_cst.svg")
		os.MkdirAll("reportes", 0755)
		err = cst.SaveCSTSVG(svg, ruta)
		if err != nil {
			return c.Status(500).SendString("No se pudo guardar el SVG")
		}
		absRuta, _ := filepath.Abs(ruta)
		go symbols.OpenHTML(absRuta) // Abre el SVG en el navegador
		fmt.Println("Ruta del CST que se envia al frontend:", absRuta)
		return c.SendString(absRuta)
	})
	app.Get("/reporte-errores", func(c *fiber.Ctx) error {
		ruta := filepath.Join("reportes", "errores.html")
		absRuta, err := filepath.Abs(ruta)
		if err != nil {
			return c.Status(500).SendString("No se pudo obtener la ruta del reporte de errores")
		}
		go symbols.OpenHTML(absRuta)
		return c.SendString(absRuta)
	})

    app.Listen(":3000")
}

/*
	func readStdin() (string, error) {
		input, err := os.ReadFile("/dev/stdin")
		return string(input), err
	}

func readStdin() (string, error) {
	//input, err := os.ReadFile("/home/brandon/Escritorio/OLC2_Proyecto1_202300813/Backend/compiler/arhivoP.vch")
	input, err := os.ReadFile("/home/brandon/Escritorio/OLC2_Proyecto1_202300813/Backend/compiler/basicas.vch")
	//input, err := os.ReadFile("/home/pablo/Escritorio/OLC2_Proyecto1_202300813/Backend/compiler/arhivoP.vch")
	//input, err := os.ReadFile("/home/vboxuser/Documents/OLC2_Proyecto1_202300813/Backend/compiler/arhivoP.vch")
	return string(input), err
}

// Funciones para visualizar nuestro arbol
func PrintVerticalTree(node antlr.Tree, ruleNames []string) {
	printVerticalNode(node, ruleNames, "", true)
}

func printVerticalNode(node antlr.Tree, ruleNames []string, prefix string, isLast bool) {
	connector := "+-- "
	if !isLast {
		connector = "|-- "
	}

	var label string
	switch n := node.(type) {
	case antlr.RuleNode:
		label = ruleNames[n.GetRuleContext().GetRuleIndex()]
	case antlr.TerminalNode:
		label = fmt.Sprintf("\"%s\"", n.GetText())
	default:
		label = fmt.Sprintf("%T", n)
	}

	fmt.Printf("%s%s%s\n", prefix, connector, label)

	// Actualizar el prefijo para los hijos
	childCount := node.GetChildCount()
	for i := 0; i < childCount; i++ {
		child := node.GetChild(i)
		newPrefix := prefix
		if isLast {
			newPrefix += "    "
		} else {
			newPrefix += "|   "
		}
		printVerticalNode(child, ruleNames, newPrefix, i == childCount-1)
	}
}*/
