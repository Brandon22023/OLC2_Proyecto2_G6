# Manual Técnico del Intérprete VLangCherry

## Universidad de San Carlos de Guatemala

**Facultad de Ingeniería**  
**Escuela de Ingeniería en Ciencias y Sistemas**  
**Curso: Organización de Lenguajes y Compiladores 2**  
**Vacaciones del Primer Semestre 2025**  
**Grupo: 6**

---

## Tabla de Contenidos

1. [Introducción](#introducción)
2. [Tecnologías Utilizadas](#tecnologías-utilizadas)
3. [Resumen General del Proyecto](#resumen-general-del-proyecto)
4. [Arquitectura VLangCherry](#-arquitectura-vlangcherry)
    - [Diagrama de Flujo de Procesamiento](#-diagrama-de-flujo-de-procesamiento)
    - [Componentes de la Arquitectura](#-componentes-de-la-arquitectura)
    - [Ventajas de la Arquitectura](#-ventajas-de-la-arquitectura)
    - [Flujo de Datos Detallado](#-flujo-de-datos-detallado)
5. [Diseño e Implementación de la Gramática](#diseño-e-implementación-de-la-gramática)
6. [Estructura y Funcionamiento del Visitor `ReplVisitor`](#estructura-y-funcionamiento-del-visitor-replvisitor)
7. [Backend de Generación ARM: visitorARM y ARMGenerator](#6-backend-de-generación-arm-visitorarm-y-armgenerator)
8. [Conclusiones](#conclusiones)
---

## Introducción

Este documento técnico describe detalladamente la implementación del intérprete para el lenguaje VLangCherry, un lenguaje con sintaxis inspirada en Go y diseñado con fines académicos para poner en práctica los principios de diseño de compiladores e intérpretes. El énfasis se encuentra en la implementación del recorrido del árbol de sintaxis abstracta (AST) mediante el patrón Visitor en el lenguaje Go, con ayuda de la herramienta ANTLR4.

---

## Tecnologías Utilizadas

* **ANTLR4**: Para generar automáticamente analizadores léxicos y sintácticos a partir de una gramática formal.
* **Go (Golang)**: Utilizado para implementar la lógica del intérprete, incluyendo recorrido de AST, manejo de entornos y ejecución.
* **Paquetes personalizados**: como `value` para representación de tipos, y `ScopeTrace` para el control de entornos.

---

## Resumen General del Proyecto

El proyecto consiste en la creación de un lenguaje llamado VLangCherry y su respectivo intérprete. El lenguaje permite manejar:

* Variables mutables e inmutables
* Tipos primitivos (`int`, `float64`, `string`, `bool`, `rune`)
* Tipos compuestos (`struct`, `slice`)
* Estructuras de control (`if`, `for`, `while`, `switch`, `break`, `continue`)
* Declaración y llamado de funciones
* Múltiples entornos y ámbito de variables

---
# 🍒 Arquitectura VLangCherry

## 📊 Diagrama de Flujo de Procesamiento

A continuación se presenta un diagrama de alto nivel que ilustra el flujo de procesamiento en **VLangCherry**, desde el código fuente hasta la ejecución interpretada o la generación de código ARM64:

---

## 🧩 Componentes de la Arquitectura

### 🔄 **Flujo Principal**

| Etapa | Componente | Descripción |
|-------|------------|-------------|
| **1** | 📄 **Código Fuente** | Archivo `.vlc` con sintaxis VLangCherry |
| **2** | 🔍 **ANTLR4** | Parser que genera tokens y estructura sintáctica |
| **3** | 🌳 **AST** | Representación abstracta del código en forma de árbol |
| **4** | 🎯 **Visitor Pattern** | Patrón de diseño para recorrer y procesar el AST |

### 🚀 **Modos de Ejecución**

#### ⚡ **Modo Intérprete (ReplVisitor)**
- ✅ **Ejecución directa** del AST
- 🏃‍♂️ **Respuesta inmediata** sin compilación
- 🧪 **Ideal para prototipado** y testing rápido
- 💡 **Debugging interactivo** línea por línea

#### 🔧 **Modo Compilación (ARMVisitor)**
- 🎯 **Traducción a ARM64** assembly
- ⚙️ **Optimización** de rendimiento
- 📦 **Generación de ejecutables** nativos
- 🏗️ **Compatibilidad** con hardware ARM

---

## 🛠️ **Componentes Especializados**

### 🎯 **ReplVisitor**
```
┌─────────────────────┐
│    ReplVisitor      │
├─────────────────────┤
│ • Evaluación directa│
│ • Manejo de estado  │
│ • Variables runtime │
│ • Funciones built-in│
└─────────────────────┘
```

### 🔧 **ARMVisitor** 
```
┌─────────────────────┐
│     ARMVisitor      │
├─────────────────────┤
│ • Mapeo AST→ARM64   │
│ • Gestión registros │
│ • Control de flujo  │
│ • Llamadas sistema  │
└─────────────────────┘
```

### ⚙️ **ArmGenerator**
```
┌─────────────────────┐
│    ArmGenerator     │
├─────────────────────┤
│ • Instrucciones ARM │
│ • Manejo de pila    │
│ • Gestión memoria   │
│ • Optimizaciones    │
└─────────────────────┘
```

---

## 🎯 **Ventajas de la Arquitectura**

### 🔄 **Flexibilidad**
- **Dual-mode**: Interpretación **Y** compilación desde el mismo AST
- **Modularidad**: Componentes intercambiables y extensibles
- **Escalabilidad**: Fácil adición de nuevos backends

### ⚡ **Rendimiento**
- **Modo intérprete**: Desarrollo rápido
- **Modo compilado**: Ejecución optimizada
- **ARM64 nativo**: Máximo rendimiento en hardware compatible

### 🛠️ **Extensibilidad**
- **Nuevos visitors**: Para diferentes arquitecturas
- **Backends adicionales**: LLVM, x86, RISC-V, etc.
- **Optimizaciones**: A nivel de AST o código generado

---

## 🔄 **Flujo de Datos Detallado**

```
📄 Código VLangCherry
        ⬇️
🔍 Tokenización (ANTLR4)
        ⬇️
🌳 Construcción AST
        ⬇️
    🎯 Decisión de Modo
       /              \
      /                \
⚡ ReplVisitor      🔧 ARMVisitor
     ⬇️                  ⬇️
🚀 Ejecución         ⚙️ ArmGenerator
   Inmediata            ⬇️
                    📝 ARM64 Code
                        ⬇️
                    🔨 Ensamblado
                        ⬇️
                    💻 Ejecución
```

--- 
## Diseño e Implementación de la Gramática

### Introducción a la gramática ANTLR

ANTLR (ANother Tool for Language Recognition) permite construir lenguajes formales mediante una gramática LL(\*). Para VLangCherry, se ha definido un conjunto de reglas que expresan tanto la estructura de alto nivel del lenguaje (como funciones, declaraciones y sentencias) como expresiones aritméticas y lógicas.

### Axioma Principal

```antlr
programa : (declaraciones)* EOF ;
```

Esta es la regla de inicio. Indica que un programa está compuesto por cero o más declaraciones y debe terminar en el fin de archivo (`EOF`).

### Declaraciones

```antlr
declaraciones
    : varDcl
    | funcDcl
    | funcMain
    | structDcl
    | stmt
    ;
```

Cada `declaraciones` puede ser una declaración de variable (`varDcl`), función (`funcDcl`), función principal (`funcMain`), definición de estructura (`structDcl`) o una sentencia ejecutable (`stmt`).

Esta modularidad permite que la gramática sea flexible, facilitando el parsing de bloques complejos y programas ordenados.

### Función Principal

```antlr
funcMain
    : 'fn' 'main' LPAREN RPAREN block
    ;
```

Define el punto de entrada del programa. Una función `main` no recibe parámetros y contiene un bloque de instrucciones entre llaves `{}`.

### Declaración de Variables

```antlr
varDcl
    : 'mut'? ID TIPO (ASSIGN expresion)?
    | 'mut'? ID ASSIGN expresion
    ;
```

La declaración de variables puede contener un tipo explícito o inferirlo con `:=` (mediante `ASSIGN`). Soporta mutabilidad opcional con la palabra clave `mut`.

Ejemplos válidos:

* `mut x int = 10`
* `y := 5.4`

### Tipos

Los tipos válidos están definidos léxicamente. La gramática reconoce `int`, `float64`, `string`, `bool`, `rune`, así como nombres de estructuras.

```antlr
TIPO : 'int' | 'float64' | 'string' | 'bool' | 'rune' ;
```

### Bloques y Sentencias

```antlr
block : LCURLY declaraciones* RCURLY ;
stmt  : (ifDcl | forDcl | whileDcl | switchDcl | printStmt | returnStmt | breakStmt | continueStmt | asignacion | expresionStmt) ';'? ;
```

Un `block` es un grupo de declaraciones entre llaves. Las `stmt` representan las unidades ejecutables.

### Expresiones

Las expresiones están altamente estructuradas para respetar precedencia:

```antlr
expresion
    : expresion op=('*' | '/' | '%') expresion
    | expresion op=('+' | '-') expresion
    | expresion op=('<' | '<=' | '>' | '>=') expresion
    | expresion op=('==' | '!=') expresion
    | ID
    | literal
    | '(' expresion ')'
    ;
```

ANTLR resolverá las ambigüedades mediante precedencia y asociación si se definen correctamente las reglas. En VLangCherry, se usan distintas reglas separadas (`Sumres`, `Relacionales`, etc.) para mantener claridad en el visitor.

### Literales y Tipos Primitivos

```antlr
literal
    : INT
    | FLOAT
    | STRING
    | BOOL
    | CHAR
    ;
```

Los literales soportados corresponden a los tipos primitivos. Se pueden extender fácilmente para admitir nuevos tipos o formas sintácticas.

### Estructuras y Atributos

```antlr
structDcl
    : 'struct' ID LCURLY atributos RCURLY
    ;
atributos
    : (TIPO ID)+
    ;
```

Las estructuras (`struct`) permiten definir tipos compuestos personalizados. Luego pueden ser instanciadas y manipuladas mediante el acceso por punto (`obj.atributo`).

### Acceso por punto

```antlr
expresion : ID DOT ID ;
```

Este patrón permite acceder a atributos dentro de una estructura:

```vch
p := Persona{ Nombre: "Juan", Edad: 25 }
print(p.Nombre)
```

El visitor lo interpreta en `VisitExpdotexp1`.

### Funciones

```antlr
funcDcl
    : 'fn' ID LPAREN parametrosFormales? RPAREN (TIPO)? block
    ;
```

Permite definir funciones con nombre, parámetros opcionales y tipo de retorno. El contenido se ejecuta como un `block` en un nuevo scope.

### Llamadas a Funciones

```antlr
expresion : ID LPAREN listaArgumentos? RPAREN ;
```

La gramática también soporta llamadas a funciones. El visitor evalúa los argumentos y busca el nombre de función en el mapa `functions`.

### Control de Flujo

```antlr
ifDcl
    : 'if' expresion block ('else' block)?
    ;
forDcl
    : 'for' ID 'in' expresion block
    ;
switchDcl
    : 'switch' expresion caseBlock+
    ;
```

Estas construcciones permiten evaluar condiciones, iterar sobre slices o rangos, y seleccionar bloques de ejecución.

### Transferencias

```antlr
breakStmt     : 'break' ;
continueStmt  : 'continue' ;
returnStmt    : 'return' expresion? ;
```

Estas reglas se interpretan por el visitor y propagan control hacia funciones o ciclos superiores.

---

La gramática de VLangCherry es lo suficientemente expresiva para representar estructuras modernas de control, funciones, expresiones complejas, y tipos de datos compuestos, manteniendo la simplicidad necesaria para un proyecto académico. A partir de esta gramática, se generan los analizadores léxico y sintáctico que construyen el árbol de sintaxis abstracta (AST), que será recorrido y ejecutado por `ReplVisitor`.

En la siguiente sección, entraremos en detalle sobre cómo este AST es recorrido, interpretado y evaluado mediante código Go implementado en el visitor `ReplVisitor`.


## 1. Función Principal: `funcMain`

### Gramática

```antlr
funcMain
    : 'fn' 'main' LPAREN RPAREN block
    ;
```

### Ejemplo VLangCherry

```vch
fn main() {
  mut x int = 5
  print(x)
}
```

### Explicación

El nodo `funcMain` es clave, pues representa el punto de entrada del programa. Durante la ejecución, el método `VisitPrograma` lo ubica y ejecuta su contenido. Se genera un nuevo scope llamado `fn_main` en `ScopeTrace`, el cual aísla las variables de `main()`.

### AST Esperado

```
(Programa
  (FuncMain fn main ( )
    (Block {
      (Declaraciones var x int = 5)
      (PrintStatement print(x))
    })
)
```

### Lógica en Visitor

```go
v.ScopeTrace.PushScope("fn_main")
// ...ejecutar declaraciones del bloque...
v.ScopeTrace.PopScope()
```

Esto garantiza encapsulamiento y limpieza del entorno.

---

## 2. Declaraciones de Variables: `varDcl`

### Gramática

```antlr
varDcl
    : 'mut'? ID TIPO (ASSIGN expresion)?
    | 'mut'? ID ASSIGN expresion
    ;
```

### Ejemplos VLangCherry

```vch
mut edad int = 25
nombre := "Juan"
```

### Explicación

Se permiten dos formas de declaración:

* Explícita: `mut x int = 5`
* Implícita: `x := 10`

En el AST se distinguirá mediante el tipo y si hay presencia de `TIPO`.

### Visitor: `VisitVariableDeclaration`

Valida el tipo, asigna el valor correcto con coerción y registra en el entorno:

```go
switch tipo {
 case "int": value.NewIntValue(...)
 case "float64": ...
}
```

Se utiliza `strconv.Atoi` o `ParseFloat` para convertir los valores.

---

## 3. Declaración de Structs: `structDcl`

### Gramática

```antlr
structDcl
    : 'struct' ID LCURLY atributos RCURLY
    ;
atributos
    : (TIPO ID)+
    ;
```

### Ejemplo VLangCherry

```vch
struct Persona {
  string nombre
  int edad
}
```

### Explicación

La gramática permite definir tipos personalizados. Al parsear `struct Persona`, el visitor guarda su definición en `map[string]map[string]string`, lo que facilita el acceso posterior.

### Visitor: `VisitStructDcl`

```go
v.structs["Persona"] = map[string]string{
  "nombre": "string",
  "edad": "int",
}
```

Este diccionario sirve como plantilla para instancias.

---

## 4. Expresiones: Operadores y Precedencia

### Gramática

ANTLR separa expresiones por precedencia:

```antlr
expresion : expresion op=('*'|'/') expresion   # MultDiv
          | expresion op=('+'|'-') expresion   # SumRes
          | expresion op=('=='|'!=') expresion # Igualdad
          | ...
```

### Visitor: `VisitSumres`, `VisitMultdivmod`

Cada regla se traduce en una función que:

* Evalúa ambos operandos
* Detecta si son strings, floats o ints
* Aplica la operación correspondiente

### Ejemplo VLangCherry

```vch
x := 5 + 3.2
saludo := "Hola, " + nombre
```

---

## 5. Acceso a Struct: `ID.DOT.ID`

### Gramática

```antlr
expresion : ID DOT ID ;
```

### Ejemplo

```vch
print(p.nombre)
```

### Visitor: `VisitExpdotexp1`

```go
structInstance := variable.Value.(*value.StructInstance)
attrVal := structInstance.Attributes["nombre"]
```

Este acceso valida existencia del atributo y su tipo.

---

## 6. Control de Flujo

### IF / ELSE

```antlr
ifDcl : 'if' expresion block ('else' block)? ;
```

### FOR

```antlr
forDcl : 'for' ID 'in' expresion block ;
```

### SWITCH

```antlr
switchDcl : 'switch' expresion caseBlock+ ;
```

### Visitor: `VisitIf_context`, `VisitFor_context`, `VisitSwitch_context`

Cada uno evalúa la condición y decide si ejecutar el bloque correspondiente. `VisitFor_context` soporta ciclos tipo `for-each` sobre slices o rangos.

### Ejemplo VLangCherry

```vch
for i in 0..10 {
  print(i)
}
```

---

## 7. Funciones

### Declaración

```antlr
funcDcl : 'fn' ID LPAREN parametrosFormales? RPAREN (TIPO)? block ;
```

### Llamado

```antlr
expresion : ID LPAREN listaArgumentos? RPAREN ;
```

### Visitor: `VisitFuncDcl`, `VisitFuncCall`

Las funciones se almacenan en `map[string]*StoredFunction`, con:

* Parámetros
* Tipo de retorno
* Cuerpo (`block`)

Durante la llamada se crea un nuevo scope con parámetros evaluados.

### Ejemplo

```vch
fn saludar(nombre string) {
  print("Hola, ", nombre)
}
saludar("Carlos")
```

---

## 8. Errores Semánticos

Todos los errores semánticos se almacenan en `v.SemanticErrors`.

### Ejemplos manejados

* Variable no declarada
* Reasignación de tipo diferente
* División por cero
* Atributo inexistente en struct

### Código

```go
v.SemanticErrors.NewSemanticError(ctx.GetStart(), "mensaje")
v.HasSemanticError = true
```

Esto evita múltiples errores en cascada.

---

## 9. Mejoras y Extensiones Futuras

A medida que se avance en el desarrollo, estas funcionalidades pueden ser añadidas:

* **Funciones anidadas**
* **Soporte para funciones recursivas**
* **Overloading de funciones**
* **Operador ternario** (`cond ? a : b`)
* **Soporte para `match` o `pattern matching` como en Rust**
* **Manejo explícito de errores como `try-catch`**
* **Tipado fuerte en structs y validación en profundidad**
* **Slices multidimensionales y matrices**


## Estructura y Funcionamiento del Visitor `ReplVisitor`

El corazón del intérprete VLangCherry es el `ReplVisitor`. Este struct implementa la interfaz `VlangVisitor` generada por ANTLR4, lo que le permite recorrer el AST y ejecutar el código.

```go
package repl

import (
	parser "compiler/parser"
	"compiler/value"
	"fmt"

	"strconv"
	"strings"

	"[github.com/antlr4-go/antlr/v4](https://github.com/antlr4-go/antlr/v4)"
)

// Visitor personalizado para recorrer el árbol de sintaxis
type ReplVisitor struct {
	parser.BaseVlangVisitor // Embebemos el visitor generado por ANTLR
	ScopeTrace              *ScopeTrace
	IfScope                 []*BaseScope
	inForLoop               bool                         // Bandera para rastrear si estamos en un bucle for
	functions               map[string]*StoredFunction   // Mapa para almacenar funciones definidas
	structs                 map[string]map[string]string // Mapa structs -> mapa (atributos)
	SemanticErrors          *ErrorTable                  // o []*Error si prefieres
	HasSemanticError        bool
}

type StoredFunction struct {
	ParamNames []string
	ParamTypes []string
	ReturnType string
	Block      parser.IBlockContext
}

var _ parser.VlangVisitor = &ReplVisitor{} // <-- Esto asegura la interfaz
```

El `ReplVisitor` contiene varios campos esenciales para su operación:

* `ScopeTrace`: Un gestor de ámbitos que permite crear, entrar y salir de ámbitos de variables, asegurando la correcta resolución de identificadores.
* `inForLoop`: Una bandera booleana para controlar el comportamiento de sentencias `break` y `continue` dentro de bucles `for`.
* `functions`: Un mapa que almacena las funciones definidas en el código, indexadas por su nombre.
* `structs`: Un mapa anidado que almacena las definiciones de los `structs`, con el nombre del `struct` como clave principal y un mapa de sus atributos (nombre del atributo a tipo) como valor.
* `SemanticErrors`: Una tabla para registrar los errores semánticos encontrados durante la ejecución.
* `HasSemanticError`: Una bandera para indicar si se ha encontrado algún error semántico, lo que a menudo detiene la ejecución posterior para evitar cascadas de errores.

### Inicialización del Visitor

El constructor `NewReplVisitor` inicializa el `ScopeTrace` y los mapas necesarios, preparando el visitor para el recorrido del AST.

```go
// Constructor del visitor
func NewReplVisitor() *ReplVisitor {
	scopeTrace := NewScopeTrace()
	return &ReplVisitor{
		ScopeTrace:       scopeTrace,
		inForLoop:        false,
		functions:        make(map[string]*StoredFunction),
		structs:          make(map[string]map[string]string),
		SemanticErrors:   NewErrorTable(),
		HasSemanticError: false,
	}
}
```

### Manejo de Tipos Primitivos

El método `VisitValor` es crucial para interpretar los literales de tipos primitivos (enteros, flotantes, cadenas y booleanos) del AST. Convierte el texto del token en un tipo concreto definido por `value.IVOR`, que es una interfaz para los valores del lenguaje VLangCherry. Esto permite una manipulación unificada de los valores.

```go
func (v *ReplVisitor) VisitValor(ctx *parser.ValorContext) interface{} {
    if ctx.INT() != nil {
        val, _ := strconv.Atoi(ctx.INT().GetText())
        return value.NewIntValue(val)
    }
    if ctx.FLOAT() != nil {
        val, _ := strconv.ParseFloat(ctx.FLOAT().GetText(), 64)
        return value.NewFloatValue(val)
    }
    if ctx.STRING() != nil {
        texto := strings.Trim(ctx.STRING().GetText(), "\"")
        return value.NewStringValue(texto)
    }
    if ctx.BOOL() != nil {
        return value.NewBoolValue(ctx.BOOL().GetText() == "true")
    }
    return nil // En caso de que no sea ninguno de los anteriores (ej. nil literal)
}
```

### Declaración y Asignación de Variables

Los métodos `VisitDeclaracion` y `VisitAsignacion` gestionan la creación y modificación de variables en el ámbito actual. `VisitDeclaracion` registra una nueva variable con su tipo, mientras que `VisitAsignacion` actualiza el valor de una variable existente, realizando verificaciones de tipo para asegurar la consistencia.

```go
// VisitDeclaracion es para 'var x int;' o 'var x Persona;'
func (v *ReplVisitor) VisitDeclaracion(ctx *parser.DeclaracionContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	id := ctx.ID().GetText()
	tipo := ctx.Tipo().GetText()

	// Comprobación semántica: Evitar redeclaración en el mismo scope
	if v.ScopeTrace.GetCurrentScope().HasVariable(id) {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Variable '%s' ya declarada en este ámbito", id))
		v.HasSemanticError = true
		return nil
	}

	// Manejo de tipos de structs
	if _, isStruct := v.structs[tipo]; isStruct {
		v.ScopeTrace.AddVariable(id, tipo, nil) // Inicializar en nil para structs
	} else {
		// Asignar valor por defecto según el tipo primitivo
		var defaultValue interface{}
		switch tipo {
		case "int":
			defaultValue = value.NewIntValue(0)
		case "float":
			defaultValue = value.NewFloatValue(0.0)
		case "string":
			defaultValue = value.NewStringValue("")
		case "bool":
			defaultValue = value.NewBoolValue(false)
		default:
			// Error: Tipo no reconocido
			v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Tipo '%s' no reconocido", tipo))
			v.HasSemanticError = true
			return nil
		}
		v.ScopeTrace.AddVariable(id, tipo, defaultValue)
	}
	return nil
}

// VisitAsignacion es para 'x = 10;' o 'p.nombre = "Ana";'
func (v *ReplVisitor) VisitAsignacion(ctx *parser.AsignacionContext) interface{} {
	if v.HasSemanticError {
		return nil
	}

	id := ctx.ID().GetText()
	val := v.Visit(ctx.Expresion())

	if val == nil { // Esto puede ocurrir si hay un error en la expresión
		return nil
	}

	variable, found := v.ScopeTrace.LookupVariable(id)
	if !found {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Variable '%s' no declarada", id))
		v.HasSemanticError = true
		return nil
	}

	// Comprobación de tipos en la asignación
	expectedType := variable.Type
	receivedType := val.(value.IVOR).GetType()

	// Convertir tipos de VLangCherry a strings para comparación
	var receivedTypeStr string
	switch receivedType {
	case value.IntValueType:
		receivedTypeStr = "int"
	case value.FloatValueType:
		receivedTypeStr = "float"
	case value.StringValueType:
		receivedTypeStr = "string"
	case value.BoolValueType:
		receivedTypeStr = "bool"
	case value.NilType: // Permitir asignar nil a cualquier tipo
		receivedTypeStr = "nil"
	case value.StructType:
		// Para structs, el tipo es el nombre del struct
		if structVal, ok := val.(*value.StructValue); ok {
			receivedTypeStr = structVal.StructName
		}
	default:
		receivedTypeStr = "desconocido"
	}

	// Si el tipo esperado es un struct, verificar que el valor asignado sea una instancia de ese struct
	if _, ok := v.structs[expectedType]; ok { // expectedType es un struct
		if receivedType != value.StructType {
			v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Tipo incompatible: Se esperaba '%s' pero se recibió '%s'", expectedType, receivedTypeStr))
			v.HasSemanticError = true
			return nil
		}
		if receivedStructVal, ok := val.(*value.StructValue); ok && receivedStructVal.StructName != expectedType {
			v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Tipo de struct incompatible: Se esperaba '%s' pero se recibió '%s'", expectedType, receivedStructVal.StructName))
			v.HasSemanticError = true
			return nil
		}
	} else if expectedType != receivedTypeStr && receivedType != value.NilType { // Para tipos primitivos, y no es nil
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Tipo incompatible: Se esperaba '%s' pero se recibió '%s'", expectedType, receivedTypeStr))
		v.HasSemanticError = true
		return nil
	}

	v.ScopeTrace.UpdateVariable(id, val.(value.IVOR))
	return nil
}
```

### Manejo de Expresiones Aisladas

`VisitExpresionStmt` es un método simple pero efectivo para ejecutar expresiones que no forman parte de una asignación o un retorno. Esto es particularmente útil para llamadas a funciones que tienen efectos secundarios (ej. `fmt.Println("Hola")`) pero no devuelven un valor que deba ser capturado.

```go
func (v *ReplVisitor) VisitExpresionStmt(ctx *parser.ExpresionStmtContext) interface{} {
    return v.Visit(ctx.Expresion())
}
```

### Operaciones con Atributos de Structs

Los `structs` son una característica importante de VLangCherry. Los métodos `VisitStructAccess` y `VisitStructAssign` permiten acceder y modificar los atributos de una instancia de `struct`.

`VisitStructAccess` resuelve el valor de un atributo de un `struct` dado su nombre y el nombre del atributo. Realiza verificaciones semánticas para asegurar que el `struct` y el atributo existan.

```go
// VisitStructAccess maneja el acceso a atributos de un struct: p.nombre
func (v *ReplVisitor) VisitStructAccess(ctx *parser.StructAccessContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	varName := ctx.ID(0).GetText() // Nombre de la instancia del struct (ej. p)
	attrName := ctx.ID(1).GetText() // Nombre del atributo (ej. nombre)

	// Buscar la variable que representa la instancia del struct
	variable, found := v.ScopeTrace.LookupVariable(varName)
	if !found {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Variable '%s' no declarada", varName))
		v.HasSemanticError = true
		return nil
	}

	// Asegurarse de que la variable sea un struct
	structInstance, ok := variable.Value.(*value.StructValue)
	if !ok {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("La variable '%s' no es una instancia de struct", varName))
		v.HasSemanticError = true
		return nil
	}

	// Verificar si el atributo existe en la definición del struct
	structDef, exists := v.structs[structInstance.StructName]
	if !exists {
		// Esto no debería pasar si el structInstance se creó correctamente, pero es una buena verificación de seguridad
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Definición de struct '%s' no encontrada", structInstance.StructName))
		v.HasSemanticError = true
		return nil
	}

	if _, ok := structDef[attrName]; !ok {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("El atributo '%s' no existe en el struct '%s'", attrName, structInstance.StructName))
		v.HasSemanticError = true
		return nil
	}

	// Obtener el valor del atributo
	val, exists := structInstance.Attributes[attrName]
	if !exists {
		// Esto no debería ocurrir si los atributos se inicializan correctamente
		return value.NewNilValue() // O un valor por defecto si es apropiado
	}
	return val
}
```

`VisitStructAssign` permite modificar el valor de un atributo específico dentro de una instancia de `struct`. Incluye la validación de tipos para asegurar que el valor asignado coincida con el tipo esperado del atributo.

```go
// VisitStructAssign maneja la asignación a atributos de un struct: p.nombre = "Juan"
func (v *ReplVisitor) VisitStructAssign(ctx *parser.StructAssignContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	varName := ctx.ID(0).GetText() // Nombre de la instancia del struct (ej. p)
	attrName := ctx.ID(1).GetText() // Nombre del atributo (ej. nombre)

	// Buscar la variable que representa la instancia del struct
	variable, found := v.ScopeTrace.LookupVariable(varName)
	if !found {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Variable '%s' no declarada", varName))
		v.HasSemanticError = true
		return nil
	}

	// Asegurarse de que la variable sea un struct
	structInstance, ok := variable.Value.(*value.StructValue)
	if !ok {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("La variable '%s' no es una instancia de struct", varName))
		v.HasSemanticError = true
		return nil
	}

	// Verificar si el atributo existe en la definición del struct
	structDef, exists := v.structs[structInstance.StructName]
	if !exists {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Definición de struct '%s' no encontrada", structInstance.StructName))
		v.HasSemanticError = true
		return nil
	}

	expectedAttrType, ok := structDef[attrName]
	if !ok {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("El atributo '%s' no existe en el struct '%s'", attrName, structInstance.StructName))
		v.HasSemanticError = true
		return nil
	}

	newValue := v.Visit(ctx.Expresion())
	if newValue == nil { // Podría ser nil si la expresión tuvo un error semántico
		return nil
	}

	// Comprobación de tipos en la asignación del atributo
	receivedType := newValue.(value.IVOR).GetType()
	var receivedTypeStr string
	switch receivedType {
	case value.IntValueType:
		receivedTypeStr = "int"
	case value.FloatValueType:
		receivedTypeStr = "float"
	case value.StringValueType:
		receivedTypeStr = "string"
	case value.BoolValueType:
		receivedTypeStr = "bool"
	case value.NilType:
		receivedTypeStr = "nil"
	case value.StructType: // Si se asigna otro struct
		if structVal, ok := newValue.(*value.StructValue); ok {
			receivedTypeStr = structVal.StructName
		}
	default:
		receivedTypeStr = "desconocido"
	}

	if expectedAttrType != receivedTypeStr && receivedType != value.NilType {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Tipo incompatible para el atributo '%s': Se esperaba '%s' pero se recibió '%s'", attrName, expectedAttrType, receivedTypeStr))
		v.HasSemanticError = true
		return nil
	}

	// Actualizar el valor del atributo en la instancia del struct
	structInstance.Attributes[attrName] = newValue
	v.ScopeTrace.UpdateVariable(varName, structInstance) // Asegurarse de que el ScopeTrace sepa que el struct ha sido modificado

	return nil
}
```

### Declaración e Inicialización Directa de Structs

`VisitStructDirectInitDeclaration` permite la declaración e inicialización de una instancia de `struct` en una sola línea, similar a `p := Persona{nombre: "Ana"}`. Este método primero verifica que el `struct` esté definido, luego crea una instancia y asigna los valores a sus atributos.

```go
func (v *ReplVisitor) VisitStructDirectInitDeclaration(ctx *parser.StructDirectInitDeclarationContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	varName := ctx.ID(0).GetText()    // p
	structName := ctx.ID(1).GetText() // Persona

	// Verificamos que el struct esté definido
	def, exists := v.structs[structName]
	if !exists {
		fmt.Printf("SEMANTICO: Struct '%s' no está definido\n", structName)
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Struct '%s' no está definido", structName))
		v.HasSemanticError = true
		return nil
	}

	// Comprobación semántica: Evitar redeclaración en el mismo scope
	if v.ScopeTrace.GetCurrentScope().HasVariable(varName) {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Variable '%s' ya declarada en este ámbito", varName))
		v.HasSemanticError = true
		return nil
	}

	// Crear la instancia del struct
	attributes := make(map[string]interface{})
	for k, vType := range def {
		// Inicializar atributos a sus valores por defecto (o nil para structs anidados)
		var defaultValue interface{}
		switch vType {
		case "int":
			defaultValue = value.NewIntValue(0)
		case "float":
			defaultValue = value.NewFloatValue(0.0)
		case "string":
			defaultValue = value.NewStringValue("")
		case "bool":
			defaultValue = value.NewBoolValue(false)
		default:
			// Si es un tipo struct anidado, se inicializa a nil
			if _, isNestedStruct := v.structs[vType]; isNestedStruct {
				defaultValue = value.NewNilValue() // Representa un struct no inicializado
			} else {
				v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Tipo de atributo desconocido '%s' en struct '%s'", vType, structName))
				v.HasSemanticError = true
				return nil
			}
		}
		attributes[k] = defaultValue
	}

	// Asignar valores a los atributos especificados en la inicialización
	for _, assign := range ctx.ListaAsignaciones().AllAsignacionStruct() {
		attr := assign.ID().GetText()
		val := v.Visit(assign.Expresion())
		if val == nil { // Podría ser nil si la expresión tuvo un error semántico
			return nil
		}

		expectedAttrType, ok := def[attr]
		if !ok {
			fmt.Printf("SEMANTICO: El atributo '%s' no existe en el struct '%s'\n", attr, structName)
			v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("El atributo '%s' no existe en el struct '%s'", attr, structName))
			v.HasSemanticError = true
			return nil
		}

		// Comprobación de tipos para los valores asignados a los atributos
		receivedType := val.(value.IVOR).GetType()
		var receivedTypeStr string
		switch receivedType {
		case value.IntValueType:
			receivedTypeStr = "int"
		case value.FloatValueType:
			receivedTypeStr = "float"
		case value.StringValueType:
			receivedTypeStr = "string"
		case value.BoolValueType:
			receivedTypeStr = "bool"
		case value.NilType:
			receivedTypeStr = "nil"
		case value.StructType:
			if structVal, ok := val.(*value.StructValue); ok {
				receivedTypeStr = structVal.StructName
			}
		default:
			receivedTypeStr = "desconocido"
		}

		if expectedAttrType != receivedTypeStr && receivedType != value.NilType {
			v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Tipo incompatible para el atributo '%s': Se esperaba '%s' pero se recibió '%s'", attr, expectedAttrType, receivedTypeStr))
			v.HasSemanticError = true
			return nil
		}

		attributes[attr] = val
	}

	// Crear la instancia final del struct
	newStructInstance := value.NewStructValue(structName, attributes)

	// Agregar la variable al ámbito actual
	v.ScopeTrace.AddVariable(varName, structName, newStructInstance)

	return nil
}
```

### Operaciones Aritméticas

Las operaciones aritméticas son fundamentales. El `ReplVisitor` tiene métodos dedicados para cada tipo de operación, como `VisitSuma`, `VisitResta`, `VisitMultiplicacion`, `VisitDivision` y `VisitModulo`. Estos métodos recuperan los valores de los operandos, realizan la operación y manejan la conversión de tipos (por ejemplo, `int` a `float` si es necesario) y la detección de errores (como división por cero o tipos incompatibles).

```go
// VisitSuma maneja la operación de suma: expr + expr
func (v *ReplVisitor) VisitSuma(ctx *parser.SumaContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	left := v.Visit(ctx.Expresion(0)).(value.IVOR)
	right := v.Visit(ctx.Expresion(1)).(value.IVOR)

	result, err := left.Add(right)
	if err != nil {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), err.Error())
		v.HasSemanticError = true
		return nil
	}
	return result
}

// VisitResta maneja la operación de resta: expr - expr
func (v *ReplVisitor) VisitResta(ctx *parser.RestaContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	left := v.Visit(ctx.Expresion(0)).(value.IVOR)
	right := v.Visit(ctx.Expresion(1)).(value.IVOR)

	result, err := left.Subtract(right)
	if err != nil {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), err.Error())
		v.HasSemanticError = true
		return nil
	}
	return result
}

// VisitMultiplicacion maneja la operación de multiplicación: expr * expr
func (v *ReplVisitor) VisitMultiplicacion(ctx *parser.MultiplicacionContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	left := v.Visit(ctx.Expresion(0)).(value.IVOR)
	right := v.Visit(ctx.Expresion(1)).(value.IVOR)

	result, err := left.Multiply(right)
	if err != nil {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), err.Error())
		v.HasSemanticError = true
		return nil
	}
	return result
}

// VisitDivision maneja la operación de división: expr / expr
func (v *ReplVisitor) VisitDivision(ctx *parser.DivisionContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	left := v.Visit(ctx.Expresion(0)).(value.IVOR)
	right := v.Visit(ctx.Expresion(1)).(value.IVOR)

	// Manejo de división por cero
	if (right.GetType() == value.IntValueType && right.GetValue().(int) == 0) ||
		(right.GetType() == value.FloatValueType && right.GetValue().(float64) == 0.0) {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), "División por cero")
		v.HasSemanticError = true
		return nil
	}

	result, err := left.Divide(right)
	if err != nil {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), err.Error())
		v.HasSemanticError = true
		return nil
	}
	return result
}

// VisitModulo maneja la operación de módulo: expr % expr
func (v *ReplVisitor) VisitModulo(ctx *parser.ModuloContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	left := v.Visit(ctx.Expresion(0)).(value.IVOR)
	right := v.Visit(ctx.Expresion(1)).(value.IVOR)

	// El módulo solo aplica a enteros
	if left.GetType() != value.IntValueType || right.GetType() != value.IntValueType {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), "Operación de módulo solo es válida para enteros")
		v.HasSemanticError = true
		return nil
	}

	// Manejo de módulo por cero
	if right.GetValue().(int) == 0 {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), "Módulo por cero")
		v.HasSemanticError = true
		return nil
	}

	result, err := left.Modulo(right)
	if err != nil {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), err.Error())
		v.HasSemanticError = true
		return nil
	}
	return result
}
```

### Operaciones Relacionales

Las operaciones relacionales (`==`, `!=`, `<`, `<=`, `>`, `>=`) se utilizan para comparar valores. Los métodos como `VisitIgualacion`, `VisitDiferenciacion`, `VisitMenorQue`, etc., evalúan las expresiones a ambos lados del operador y devuelven un valor booleano. Se implementan verificaciones de compatibilidad de tipos para asegurar comparaciones significativas.

```go
// VisitIgualacion maneja la operación de igualdad: expr == expr
func (v *ReplVisitor) VisitIgualacion(ctx *parser.IgualacionContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	left := v.Visit(ctx.Expresion(0)).(value.IVOR)
	right := v.Visit(ctx.Expresion(1)).(value.IVOR)

	result, err := left.Equals(right)
	if err != nil {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), err.Error())
		v.HasSemanticError = true
		return nil
	}
	return result
}

// VisitDiferenciacion maneja la operación de desigualdad: expr != expr
func (v *ReplVisitor) VisitDiferenciacion(ctx *parser.DiferenciacionContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	left := v.Visit(ctx.Expresion(0)).(value.IVOR)
	right := v.Visit(ctx.Expresion(1)).(value.IVOR)

	result, err := left.NotEquals(right)
	if err != nil {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), err.Error())
		v.HasSemanticError = true
		return nil
	}
	return result
}

// VisitMenorQue maneja la operación de menor que: expr < expr
func (v *ReplVisitor) VisitMenorQue(ctx *parser.MenorQueContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	left := v.Visit(ctx.Expresion(0)).(value.IVOR)
	right := v.Visit(ctx.Expresion(1)).(value.IVOR)

	result, err := left.LessThan(right)
	if err != nil {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), err.Error())
		v.HasSemanticError = true
		return nil
	}
	return result
}

// VisitMenorIgual maneja la operación de menor o igual que: expr <= expr
func (v *ReplVisitor) VisitMenorIgual(ctx *parser.MenorIgualContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	left := v.Visit(ctx.Expresion(0)).(value.IVOR)
	right := v.Visit(ctx.Expresion(1)).(value.IVOR)

	result, err := left.LessThanOrEquals(right)
	if err != nil {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), err.Error())
		v.HasSemanticError = true
		return nil
	}
	return result
}

// VisitMayorQue maneja la operación de mayor que: expr > expr
func (v *ReplVisitor) VisitMayorQue(ctx *parser.MayorQueContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	left := v.Visit(ctx.Expresion(0)).(value.IVOR)
	right := v.Visit(ctx.Expresion(1)).(value.IVOR)

	result, err := left.GreaterThan(right)
	if err != nil {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), err.Error())
		v.HasSemanticError = true
		return nil
	}
	return result
}

// VisitMayorIgual maneja la operación de mayor o igual que: expr >= expr
func (v *ReplVisitor) VisitMayorIgual(ctx *parser.MayorIgualContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	left := v.Visit(ctx.Expresion(0)).(value.IVOR)
	right := v.Visit(ctx.Expresion(1)).(value.IVOR)

	result, err := left.GreaterThanOrEquals(right)
	if err != nil {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), err.Error())
		v.HasSemanticError = true
		return nil
	}
	return result
}
```

### Operaciones Lógicas

Las operaciones lógicas (`&&`, `||`, `!`) son esenciales para la toma de decisiones. `VisitAnd`, `VisitOr` y `VisitNot` implementan la lógica booleana, asegurando que los operandos sean de tipo booleano y devolviendo el resultado correcto.

```go
// VisitAnd maneja la operación lógica AND: expr && expr
func (v *ReplVisitor) VisitAnd(ctx *parser.AndContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	left := v.Visit(ctx.Expresion(0)).(value.IVOR)
	right := v.Visit(ctx.Expresion(1)).(value.IVOR)

	// Ambos operandos deben ser booleanos
	if left.GetType() != value.BoolValueType || right.GetType() != value.BoolValueType {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), "Operación AND solo es válida para booleanos")
		v.HasSemanticError = true
		return nil
	}

	leftVal := left.GetValue().(bool)
	rightVal := right.GetValue().(bool)
	return value.NewBoolValue(leftVal && rightVal)
}

// VisitOr maneja la operación lógica OR: expr || expr
func (v *ReplVisitor) VisitOr(ctx *parser.OrContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	left := v.Visit(ctx.Expresion(0)).(value.IVOR)
	right := v.Visit(ctx.Expresion(1)).(value.IVOR)

	// Ambos operandos deben ser booleanos
	if left.GetType() != value.BoolValueType || right.GetType() != value.BoolValueType {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), "Operación OR solo es válida para booleanos")
		v.HasSemanticError = true
		return nil
	}

	leftVal := left.GetValue().(bool)
	rightVal := right.GetValue().(bool)
	return value.NewBoolValue(leftVal || rightVal)
}

// VisitNot maneja la operación lógica NOT: !expr
func (v *ReplVisitor) VisitNot(ctx *parser.NotContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	expr := v.Visit(ctx.Expresion()).(value.IVOR)

	// El operando debe ser booleano
	if expr.GetType() != value.BoolValueType {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), "Operación NOT solo es válida para booleanos")
		v.HasSemanticError = true
		return nil
	}

	exprVal := expr.GetValue().(bool)
	return value.NewBoolValue(!exprVal)
}
```

### Control de Flujo: Condicionales (`if-else`)

Los condicionales `if-else` son implementados por `VisitIfStmt`, `VisitElseIfStmt` y `VisitElseStmt`. Estos métodos evalúan la expresión de la condición y, si es verdadera, ejecutan el bloque de código asociado. Es crucial el manejo de ámbitos aquí: se crea un nuevo ámbito para cada bloque `if` o `else` para que las variables declaradas dentro de ellos sean locales.

```go
// VisitIfStmt maneja la sentencia if
func (v *ReplVisitor) VisitIfStmt(ctx *parser.IfStmtContext) interface{} {
	if v.HasSemanticError {
		return nil
	}

	condition := v.Visit(ctx.Expresion()).(value.IVOR)
	if condition == nil {
		return nil // Error en la condición
	}

	if condition.GetType() != value.BoolValueType {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), "La condición del if debe ser booleana")
		v.HasSemanticError = true
		return nil
	}

	if condition.GetValue().(bool) {
		v.ScopeTrace.PushScope() // Entrar a un nuevo ámbito para el bloque if
		v.Visit(ctx.Block())
		v.ScopeTrace.PopScope() // Salir del ámbito
	} else if ctx.ElseIfStmt() != nil {
		v.Visit(ctx.ElseIfStmt())
	} else if ctx.ElseStmt() != nil {
		v.Visit(ctx.ElseStmt())
	}
	return nil
}

// VisitElseIfStmt maneja la sentencia else if
func (v *ReplVisitor) VisitElseIfStmt(ctx *parser.ElseIfStmtContext) interface{} {
	if v.HasSemanticError {
		return nil
	}

	condition := v.Visit(ctx.Expresion()).(value.IVOR)
	if condition == nil {
		return nil // Error en la condición
	}

	if condition.GetType() != value.BoolValueType {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), "La condición del else if debe ser booleana")
		v.HasSemanticError = true
		return nil
	}

	if condition.GetValue().(bool) {
		v.ScopeTrace.PushScope() // Entrar a un nuevo ámbito
		v.Visit(ctx.Block())
		v.ScopeTrace.PopScope() // Salir del ámbito
	} else if ctx.ElseIfStmt() != nil {
		v.Visit(ctx.ElseIfStmt())
	} else if ctx.ElseStmt() != nil {
		v.Visit(ctx.ElseStmt())
	}
	return nil
}

// VisitElseStmt maneja la sentencia else
func (v *ReplVisitor) VisitElseStmt(ctx *parser.ElseStmtContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	v.ScopeTrace.PushScope() // Entrar a un nuevo ámbito
	v.Visit(ctx.Block())
	v.ScopeTrace.PopScope() // Salir del ámbito
	return nil
}
```

### Control de Flujo: Bucles (`for`)

Los bucles `for` son manejados por `VisitForStmt`. Este método implementa la lógica de inicialización, condición y post-sentencia de un bucle `for` tradicional. También gestiona las sentencias `break` y `continue` mediante la bandera `inForLoop` y valores de retorno especiales.

```go
// VisitForStmt maneja la sentencia for
func (v *ReplVisitor) VisitForStmt(ctx *parser.ForStmtContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	v.ScopeTrace.PushScope() // Entrar a un nuevo ámbito para el for

	// Parte 1: Inicialización (opcional)
	if ctx.ForInit() != nil {
		v.Visit(ctx.ForInit())
	}

	v.inForLoop = true // Indicamos que estamos dentro de un bucle

	// Bucle principal
	for {
		// Parte 2: Condición (opcional)
		if ctx.ForCondition() != nil {
			condition := v.Visit(ctx.ForCondition()).(value.IVOR)
			if condition == nil { // Error en la condición
				v.inForLoop = false
				v.ScopeTrace.PopScope()
				return nil
			}
			if condition.GetType() != value.BoolValueType {
				v.SemanticErrors.NewSemanticError(ctx.GetStart(), "La condición del for debe ser booleana")
				v.HasSemanticError = true
				v.inForLoop = false
				v.ScopeTrace.PopScope()
				return nil
			}
			if !condition.GetValue().(bool) {
				break // Salir del bucle si la condición es falsa
			}
		}

		// Ejecutar el bloque del bucle
		result := v.Visit(ctx.Block())

		if result != nil {
			// Manejo de sentencias 'break' y 'continue'
			if controlFlow, ok := result.(*ControlFlow); ok {
				if controlFlow.Type == Break {
					break // Salir completamente del for
				} else if controlFlow.Type == Continue {
					// Continuar a la siguiente iteración (ejecutar post-statement)
				}
			} else if _, isReturn := result.(*ReturnValue); isReturn {
				// Si hay un return dentro del for, debe propagarse hacia arriba
				v.inForLoop = false
				v.ScopeTrace.PopScope()
				return result
			}
		}

		// Parte 3: Post-sentencia (opcional)
		if ctx.ForPost() != nil {
			v.Visit(ctx.ForPost())
		}
	}

	v.inForLoop = false
	v.ScopeTrace.PopScope() // Salir del ámbito del for
	return nil
}

// VisitBreakStmt maneja la sentencia break
func (v *ReplVisitor) VisitBreakStmt(ctx *parser.BreakStmtContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	if !v.inForLoop {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), "La sentencia 'break' debe estar dentro de un bucle 'for'")
		v.HasSemanticError = true
		return nil
	}
	return &ControlFlow{Type: Break} // Devolver una señal para romper el bucle
}

// VisitContinueStmt maneja la sentencia continue
func (v *ReplVisitor) VisitContinueStmt(ctx *parser.ContinueStmtContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	if !v.inForLoop {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), "La sentencia 'continue' debe estar dentro de un bucle 'for'")
		v.HasSemanticError = true
		return nil
	}
	return &ControlFlow{Type: Continue} // Devolver una señal para continuar la siguiente iteración
}

// Definimos structs para controlar el flujo (break, continue, return)
type ControlFlowType int

const (
	Normal ControlFlowType = iota
	Break
	Continue
	Return
)

type ControlFlow struct {
	Type ControlFlowType
}

type ReturnValue struct {
	Value interface{}
}
```

### Declaración y Llamada a Funciones

La implementación de funciones es una parte compleja pero vital.

* `VisitDeclaracionFuncion` registra la función en el mapa `functions`, almacenando su nombre, parámetros, tipo de retorno y el bloque de código.
* `VisitLlamadaFuncion` es el que maneja la ejecución. Cuando se llama una función, se crea un *nuevo ámbito* para esa llamada, se evalúan los argumentos y se asignan a los parámetros de la función. Luego, el visitor se adentra en el bloque de la función. Si la función retorna un valor, este es capturado y devuelto.

```go
// VisitDeclaracionFuncion maneja la declaración de funciones
func (v *ReplVisitor) VisitDeclaracionFuncion(ctx *parser.DeclaracionFuncionContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	funcName := ctx.ID().GetText()

	// Comprobación semántica: Evitar redeclaración de función
	if _, exists := v.functions[funcName]; exists {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Función '%s' ya declarada", funcName))
		v.HasSemanticError = true
		return nil
	}

	var paramNames []string
	var paramTypes []string
	if ctx.Parametros() != nil {
		for _, param := range ctx.Parametros().AllParametro() {
			paramNames = append(paramNames, param.ID().GetText())
			paramTypes = append(paramTypes, param.Tipo().GetText())
		}
	}

	returnType := "void" // Por defecto si no se especifica
	if ctx.Tipo() != nil {
		returnType = ctx.Tipo().GetText()
	}

	v.functions[funcName] = &StoredFunction{
		ParamNames: paramNames,
		ParamTypes: paramTypes,
		ReturnType: returnType,
		Block:      ctx.Block(),
	}
	return nil
}

// VisitLlamadaFuncion maneja la llamada a funciones
func (v *ReplVisitor) VisitLlamadaFuncion(ctx *parser.LlamadaFuncionContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	funcName := ctx.ID().GetText()

	// Manejo de la función print
	if funcName == "print" {
		return v.VisitPrintFunction(ctx)
	}

	storedFunc, exists := v.functions[funcName]
	if !exists {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Función '%s' no declarada", funcName))
		v.HasSemanticError = true
		return nil
	}

	// Obtener los argumentos pasados
	var args []value.IVOR
	if ctx.Argumentos() != nil {
		for _, exprCtx := range ctx.Argumentos().AllExpresion() {
			argVal := v.Visit(exprCtx)
			if argVal == nil { // Si hubo un error al evaluar un argumento
				return nil
			}
			args = append(args, argVal.(value.IVOR))
		}
	}

	// Validar número de argumentos
	if len(args) != len(storedFunc.ParamNames) {
		v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Número incorrecto de argumentos para la función '%s'. Se esperaban %d, se recibieron %d", funcName, len(storedFunc.ParamNames), len(args)))
		v.HasSemanticError = true
		return nil
	}

	// Validar tipos de argumentos y crear un nuevo ámbito para la función
	v.ScopeTrace.PushScope() // Nuevo ámbito para la función
	for i, paramName := range storedFunc.ParamNames {
		expectedType := storedFunc.ParamTypes[i]
		receivedArg := args[i]
		receivedType := receivedArg.GetType()

		var receivedTypeStr string
		switch receivedType {
		case value.IntValueType:
			receivedTypeStr = "int"
		case value.FloatValueType:
			receivedTypeStr = "float"
		case value.StringValueType:
			receivedTypeStr = "string"
		case value.BoolValueType:
			receivedTypeStr = "bool"
		case value.NilType:
			receivedTypeStr = "nil"
		case value.StructType:
			if structVal, ok := receivedArg.(*value.StructValue); ok {
				receivedTypeStr = structVal.StructName
			}
		default:
			receivedTypeStr = "desconocido"
		}

		if expectedType != receivedTypeStr && receivedType != value.NilType {
			v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Tipo incompatible para el parámetro '%s' en la función '%s': Se esperaba '%s' pero se recibió '%s'", paramName, funcName, expectedType, receivedTypeStr))
			v.HasSemanticError = true
			v.ScopeTrace.PopScope() // Salir del ámbito antes de regresar
			return nil
		}
		v.ScopeTrace.AddVariable(paramName, expectedType, receivedArg)
	}

	// Ejecutar el bloque de la función
	returnValue := v.Visit(storedFunc.Block)

	v.ScopeTrace.PopScope() // Salir del ámbito de la función

	// Manejar el valor de retorno de la función
	if returnVal, ok := returnValue.(*ReturnValue); ok {
		// Validar el tipo de retorno si la función no es 'void'
		if storedFunc.ReturnType != "void" && returnVal.Value == nil {
			v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("La función '%s' de tipo '%s' no devuelve un valor", funcName, storedFunc.ReturnType))
			v.HasSemanticError = true
			return nil
		}

		if returnVal.Value != nil {
			returnedIVOR, isIVOR := returnVal.Value.(value.IVOR)
			if !isIVOR {
				// Esto no debería pasar si todo va bien con value.IVOR
				v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Valor de retorno inválido de la función '%s'", funcName))
				v.HasSemanticError = true
				return nil
			}

			var returnedTypeStr string
			switch returnedIVOR.GetType() {
			case value.IntValueType:
				returnedTypeStr = "int"
			case value.FloatValueType:
				returnedTypeStr = "float"
			case value.StringValueType:
				returnedTypeStr = "string"
			case value.BoolValueType:
				returnedTypeStr = "bool"
			case value.NilType:
				returnedTypeStr = "nil"
			case value.StructType:
				if structVal, ok := returnedIVOR.(*value.StructValue); ok {
					returnedTypeStr = structVal.StructName
				}
			default:
				returnedTypeStr = "desconocido"
			}

			if storedFunc.ReturnType != returnedTypeStr && returnedIVOR.GetType() != value.NilType {
				v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Tipo de retorno incompatible para la función '%s': Se esperaba '%s' pero se recibió '%s'", funcName, storedFunc.ReturnType, returnedTypeStr))
				v.HasSemanticError = true
				return nil
			}
			return returnedIVOR // Retornar el valor que devolvió la función
		} else {
			// Si la función debía retornar algo pero no lo hizo (o retornó nil)
			if storedFunc.ReturnType != "void" {
				v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("La función '%s' de tipo '%s' no devolvió un valor válido", funcName, storedFunc.ReturnType))
				v.HasSemanticError = true
				return nil
			}
		}
	} else {
		// Si no hay ReturnValue, significa que la función no tuvo una sentencia 'return' o fue void
		if storedFunc.ReturnType != "void" {
			v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("La función '%s' de tipo '%s' no tiene una sentencia de retorno", funcName, storedFunc.ReturnType))
			v.HasSemanticError = true
			return nil
		}
	}

	return value.NewNilValue() // Si no se devuelve nada explícitamente, se retorna nil
}

// VisitReturnStmt maneja la sentencia return
func (v *ReplVisitor) VisitReturnStmt(ctx *parser.ReturnStmtContext) interface{} {
	if v.HasSemanticError {
		return nil
	}
	if ctx.Expresion() != nil {
		val := v.Visit(ctx.Expresion()).(value.IVOR)
		return &ReturnValue{Value: val} // Empaquetar el valor de retorno
	}
	return &ReturnValue{Value: value.NewNilValue()} // Retorno sin valor (void)
}
```

### Función `print`

La función `print` es una función incorporada esencial para la salida de datos. `VisitPrintFunction` toma los argumentos pasados a `print`, los convierte a su representación de cadena y los imprime en la consola.

```go
// VisitPrintFunction maneja la función print (implementación interna)
func (v *ReplVisitor) VisitPrintFunction(ctx *parser.LlamadaFuncionContext) interface{} {
	if v.HasSemanticError {
		return nil
	}

	var toPrint []string
	if ctx.Argumentos() != nil {
		for _, exprCtx := range ctx.Argumentos().AllExpresion() {
			val := v.Visit(exprCtx)
			if val == nil {
				return nil // Error en el argumento
			}
			if ivorVal, ok := val.(value.IVOR); ok {
				toPrint = append(toPrint, ivorVal.String())
			} else {
				// Esto es un caso de error si algo que no es IVOR se pasó
				v.SemanticErrors.NewSemanticError(ctx.GetStart(), fmt.Sprintf("Tipo de argumento inválido para print: %T", val))
				v.HasSemanticError = true
				return nil
			}
		}
	}
	fmt.Println(strings.Join(toPrint, " "))
	return value.NewNilValue() // La función print no devuelve nada
}
```

**Otros Métodos Importantes en `ReplVisitor` (Mencionados y Explicados Brevemente):**

* **`VisitProg`**: El punto de entrada principal para el recorrido del AST. Recorre las sentencias del programa una por una.
* **`VisitBloque`**: Gestiona la entrada y salida de ámbitos para bloques de código (ej. dentro de funciones, `if`s, `for`s).
* **`VisitDeclaracionStruct`**: Registra la estructura de un `struct` (sus atributos y tipos) en el mapa `structs`.
* **`VisitAccesoID`**: Busca el valor de un identificador (variable) en el `ScopeTrace`.

Estas funciones, combinadas con las vistas anteriormente, completan el recorrido del AST generado por ANTLR y ejecutan de manera estructurada todo el lenguaje VLangCherry. El diseño modular permite agregar nuevas construcciones fácilmente, manteniendo la coherencia del sistema y facilitando la depuración y extensión del lenguaje.

---
## 6. Backend de Generación ARM: visitorARM y ARMGenerator

### Introducción

Además del intérprete, VLangCherry cuenta con un backend de **generación de código ARM64**, permitiendo traducir programas escritos en el lenguaje a instrucciones ensamblador ARM. Esto posibilita la ejecución nativa en arquitecturas ARM reales o emuladas, y demuestra la separación entre la lógica de interpretación y la de compilación.

---

### Arquitectura General

El backend ARM está compuesto principalmente por dos módulos:

- **visitorARM.go (`ARMVisitor`)**: Visitor que recorre el AST generado por ANTLR y traduce cada nodo a instrucciones ARM.
- **ARMGenerator.go (`ArmGenerator`)**: Encapsula la lógica para emitir instrucciones ARM, gestionar etiquetas, variables, pila y la sección de datos.

La interacción entre ambos módulos sigue el patrón:  
**AST → ARMVisitor → ArmGenerator → Código ARM**

---

### Principales Responsabilidades

#### ARMVisitor

- Traduce cada construcción del AST (declaraciones, expresiones, control de flujo, funciones, etc.) a instrucciones ARM.
- Maneja el mapeo de variables a registros, offsets o labels ARM.
- Controla el flujo de ejecución mediante la generación de etiquetas y saltos (`B`, `BEQ`, `BNE`, etc.).
- Propaga correctamente las sentencias de control (`return`, `break`, `continue`) usando saltos y etiquetas.
- Gestiona el valor de retorno de funciones a través del registro `x0` y saltos al final de la función.

#### ArmGenerator

- Provee métodos para emitir instrucciones ARM (`MOV`, `ADD`, `CMP`, `B`, etc.).
- Genera etiquetas únicas para saltos y bloques de código.
- Maneja la pila y la memoria local (push/pop, offsets).
- Define la sección de datos (strings, variables globales, etc.).
- Implementa funciones auxiliares como `strcmp`, conversión de tipos y rutinas de impresión.

---

### Traducción de Sentencias y Expresiones

- **Variables**: Se asignan a registros o a posiciones relativas en la pila.
- **Expresiones aritméticas y lógicas**: Se traducen a instrucciones ARM nativas, respetando la precedencia y los tipos.
- **Control de flujo**: Se generan etiquetas y saltos condicionales para `if`, `for`, `switch`, etc.
- **Funciones**: Los parámetros y valores de retorno se manejan mediante los registros estándar (`x0`, `x1`, ...). El visitor asegura que el valor de retorno esté en `x0` y salta al final de la función.
- **Return**: Coloca el valor en `x0` y genera un salto incondicional al label de fin de función.

---

### Ejemplo de Traducción

Supongamos el siguiente código VLangCherry:

```vch
fn suma(a int, b int) int {
  return a + b
}
```

### Manejo de Control de Flujo y Return
If/Else: Se generan etiquetas para las ramas y saltos condicionales según el resultado de la comparación.
For: Se crean etiquetas para el inicio, condición, cuerpo, incremento y fin del ciclo.
Switch: Se compara la expresión principal con cada caso y se salta al bloque correspondiente.
Return: El visitor coloca el valor en x0 y salta al final de la función usando un label único. Se utiliza una pila de labels de retorno para soportar funciones anidadas o recursivas.

## Ventajas y Limitaciones
### Ventajas:

Separación de lógica: El visitor ARM se encarga de la traducción, mientras que el generador abstrae la emisión de instrucciones.
Extensible: Es sencillo agregar nuevas construcciones del lenguaje o instrucciones ARM.
Portabilidad: El código generado puede ejecutarse en hardware ARM real o emuladores.

### Limitaciones:

Actualmente solo se soporta ARM64.
No se implementan optimizaciones avanzadas.
El manejo de errores en tiempo de ejecución es limitado.

## Ejemplo de Flujo de Traducción
### Parsing: ANTLR genera el AST a partir del código fuente.
### Visita: ARMVisitor recorre el AST y utiliza ArmGenerator para emitir instrucciones ARM.
### Salida: El código ARM generado se puede ensamblar y ejecutar en una arquitectura ARM64.


## Conclusiones

La implementación de VLangCherry representa un proyecto integral que abarca tanto la interpretación como la compilación a bajo nivel, integrando conceptos avanzados de diseño de lenguajes, estructuras de datos, teoría de compiladores y arquitectura de computadoras. A lo largo del desarrollo, se han consolidado aprendizajes clave y se han superado retos técnicos que enriquecen la experiencia y el valor académico del sistema.

### Consolidación de la Arquitectura

El uso de ANTLR4 para la generación del parser y el patrón Visitor para el recorrido del AST han permitido desacoplar la lógica de análisis sintáctico de la semántica y la ejecución. Esto facilita la extensión del lenguaje, la incorporación de nuevas construcciones y la experimentación con diferentes paradigmas de ejecución. El sistema de manejo de ámbitos (`ScopeTrace`) ha sido esencial para la correcta gestión de variables, funciones y estructuras, permitiendo la implementación de reglas de visibilidad, recursión y encapsulamiento.

### Impacto del Backend ARM

Uno de los aportes más significativos de VLangCherry es la integración de un backend de generación de código ARM64. Este componente transforma el lenguaje de un simple intérprete académico a una herramienta capaz de producir código ensamblador ejecutable en arquitecturas ARM reales o emuladas, como Raspberry Pi, servidores ARM o entornos virtualizados. La generación de código ARM implica desafíos adicionales, como el manejo explícito de registros, la gestión de la pila, la asignación eficiente de variables y la traducción de estructuras de control de alto nivel (if, for, switch) a instrucciones de bajo nivel.

El visitor `ARMVisitor` y el generador `ArmGenerator` trabajan en conjunto para recorrer el AST y emitir instrucciones ARM optimizadas. El visitor se encarga de traducir cada construcción del lenguaje fuente a una secuencia de instrucciones, mientras que el generador abstrae la complejidad de la sintaxis ARM, permitiendo la reutilización de rutinas como la comparación de cadenas (`strcmp`), la conversión de tipos y la impresión de valores. Esta separación de responsabilidades facilita la depuración, el mantenimiento y la futura extensión del sistema.

### Desafíos Técnicos y Soluciones

Durante el desarrollo del backend ARM, se enfrentaron varios retos técnicos:

- **Manejo de Control de Flujo:** Traducir sentencias como `if`, `for`, `switch`, `break`, `continue` y `return` a saltos y etiquetas ARM requiere una comprensión profunda del flujo de ejecución a bajo nivel. Se implementó una pila de labels de retorno para soportar funciones anidadas y recursivas, asegurando que cada `return` salte correctamente al final de la función.
- **Gestión de la Pila y Registros:** La asignación de variables a registros o posiciones en la pila fue fundamental para evitar colisiones y garantizar la integridad de los datos durante la ejecución. Se diseñaron métodos específicos para el manejo de push/pop y offsets en la pila, así como para la gestión de variables locales y globales.
- **Traducción de Tipos y Operaciones:** Se implementaron rutinas para la conversión y comparación de tipos, permitiendo que operaciones aritméticas, lógicas y de comparación funcionen correctamente en ARM, respetando la semántica del lenguaje fuente.
- **Sección de Datos:** El generador ARM administra una sección de datos donde se almacenan cadenas, variables globales y otros literales, facilitando el acceso eficiente durante la ejecución.

### Ventajas de la Arquitectura Modular

El diseño modular de VLangCherry, con componentes claramente diferenciados para la interpretación (`ReplVisitor`) y la compilación (`ARMVisitor` y `ArmGenerator`), ofrece múltiples ventajas:

- **Extensibilidad:** Es sencillo agregar nuevas construcciones del lenguaje o instrucciones ARM, así como adaptar el backend para otras arquitecturas en el futuro (por ejemplo, x86 o RISC-V).
- **Portabilidad:** El código ARM generado puede ejecutarse en una amplia variedad de dispositivos, desde sistemas embebidos hasta servidores de alto rendimiento.
- **Separación de responsabilidades:** La lógica de interpretación y compilación está desacoplada, permitiendo evolucionar ambos componentes de manera independiente.
- **Facilidad de depuración y pruebas:** La generación de código intermedio y la posibilidad de inspeccionar el ensamblador ARM facilitan la identificación de errores y la validación del comportamiento del lenguaje.

### Aprendizajes y Aplicaciones Futuras

El desarrollo de VLangCherry ha consolidado conocimientos clave en:

- Diseño de gramáticas formales y generación de analizadores sintácticos con ANTLR4.
- Implementación del patrón Visitor para el recorrido y evaluación de árboles de sintaxis abstracta.
- Traducción de alto nivel a bajo nivel, enfrentando los retos de la arquitectura ARM64.
- Manejo de ámbitos, tipos, errores semánticos y control de flujo en un lenguaje propio.
- Optimización de recursos y eficiencia en la generación de código ensamblador.

Como líneas de trabajo futuro se identifican:

- **Optimización del código ARM generado:** Implementar técnicas como eliminación de código muerto, propagación de constantes y mejor asignación de registros.
- **Soporte para más arquitecturas:** Adaptar el backend para generar código para x86, RISC-V u otras plataformas populares.
- **Mejoras en el manejo de errores en tiempo de ejecución:** Incluir rutinas de verificación y manejo de excepciones en el código ensamblador.
- **Generación de código intermedio:** Introducir una capa de representación intermedia (IR) para facilitar optimizaciones y análisis estático antes de la traducción final a ARM.
- **Herramientas de depuración:** Desarrollar utilidades para visualizar el flujo de ejecución y el estado de la pila/registros durante la ejecución del código ARM.
- **Integración con toolchains ARM:** Automatizar el ensamblado y la ejecución del código generado en dispositivos ARM reales o emulados.

### Reflexión Final

VLangCherry es mucho más que un intérprete académico: es una plataforma de experimentación y aprendizaje en el diseño de lenguajes, capaz de producir código eficiente y portable para arquitecturas modernas. La integración de un backend ARM64 demuestra la madurez del proyecto y su potencial para aplicaciones reales, tanto en el ámbito educativo como en el desarrollo de software de sistemas. Este trabajo sienta las bases para futuras investigaciones y desarrollos en compiladores, intérpretes y generación de código para arquitecturas heterogéneas, y constituye un aporte valioso para la formación de ingenieros en ciencias de la computación y sistemas.

