package repl

type FrameElement struct {
    Name   string
    Offset int
}

type FrameVisitor struct {
    Frame      []FrameElement
    LocalOffset int
    BaseOffset  int
}

func NewFrameVisitor(baseOffset int) *FrameVisitor {
    return &FrameVisitor{
        Frame:      []FrameElement{},
        LocalOffset: 0,
        BaseOffset:  baseOffset,
    }
}

// VisitVarDcl simula el método VisitVarDcl del visitor
func (v *FrameVisitor) VisitVarDcl(ctx *VarDclContext) interface{} {
    name := ctx.ID().GetText()
    v.Frame = append(v.Frame, FrameElement{
        Name:   name,
        Offset: v.BaseOffset + v.LocalOffset,
    })
    v.LocalOffset += 1
    return nil
}

// VisitBlockStmt simula el método VisitBlockStmt del visitor
func (v *FrameVisitor) VisitBlockStmt(ctx *BlockStmtContext) interface{} {
    for _, dcl := range ctx.Dcl() {
        v.VisitVarDcl(dcl)
    }
    return nil
}

// VisitIfStmt simula el método VisitIfStmt del visitor
func (v *FrameVisitor) VisitIfStmt(ctx *IfStmtContext) interface{} {
    stmts := ctx.Stmt()
    if len(stmts) > 0 {
        v.VisitBlockStmt(stmts[0])
    }
    if len(stmts) > 1 {
        v.VisitBlockStmt(stmts[1])
    }
    return nil
}

// VisitForNormalStmt simula el método VisitForNormalStmt del visitor
func (v *FrameVisitor) VisitForNormalStmt(ctx *ForNormalStmtContext) interface{} {
    if ctx.VarDcl() != nil {
        v.VisitVarDcl(ctx.VarDcl())
    }
    v.VisitBlockStmt(ctx.Stmt())
    return nil
}

// VisitForStmt simula el método VisitForStmt del visitor
func (v *FrameVisitor) VisitForStmt(ctx *ForStmtContext) interface{} {
    v.VisitBlockStmt(ctx.Stmt())
    return nil
}

// ---
// Los siguientes tipos y métodos son placeholders para simular el contexto de ANTLR en Go.
// Debes reemplazarlos por los generados por tu parser ANTLR en Go.

type VarDclContext struct{}
func (c *VarDclContext) ID() *IDContext { return &IDContext{} }

type IDContext struct{}
func (c *IDContext) GetText() string { return "variable" }

type BlockStmtContext struct{}
func (c *BlockStmtContext) Dcl() []*VarDclContext { return []*VarDclContext{} }

// ...existing code...

type IfStmtContext struct{}
func (c *IfStmtContext) Stmt() []*BlockStmtContext { return []*BlockStmtContext{} }

// ...existing code...
type ForNormalStmtContext struct{}
func (c *ForNormalStmtContext) VarDcl() *VarDclContext { return &VarDclContext{} }
func (c *ForNormalStmtContext) Stmt() *BlockStmtContext { return &BlockStmtContext{} }

type ForStmtContext struct{}
func (c *ForStmtContext) Stmt() *BlockStmtContext { return &BlockStmtContext{} }