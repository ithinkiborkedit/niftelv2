package codegen

import (
	"fmt"
	"strings"

	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
)

type Codegen struct {
	builder strings.Builder
}

func NewCodeGen() *Codegen {
	return &Codegen{}
}

func (c *Codegen) GenerateLLVM(stmts []ast.Stmt) (string, error) {
	c.emitPreamble()
	c.emitMainFunction(stmts)
	return c.builder.String(), nil
}

func (c *Codegen) emitPreamble() {
	c.builder.WriteString(`declare i32 @print(i8*,...)
	@print.str = constant [4 x i8] c"%\d\0A\00"

	`)
}

func (c *Codegen) emitMainFunction(stmts []ast.Stmt) {
	c.builder.WriteString("define i32 @main(){\n")
	c.builder.WriteString("entry:\n")

	for _, stmt := range stmts {
		c.emitStmt(stmt)
	}

	c.builder.WriteString(" ret i32 0\n")
	c.builder.WriteString("}\n")
}

func (c *Codegen) emitStmt(s ast.Stmt) {
	switch stmt := s.(type) {
	case *ast.PrintStmt:
		c.emitPrint(stmt)
	default:
		fmt.Printf("warning: unsupported statement type: %T\n", stmt)
	}
}

func (c *Codegen) emitPrint(s *ast.PrintStmt) {
	lit, ok := s.Expr.(*ast.LiteralExpr)
	if !ok {
		fmt.Println("warning: print only supports literal ints now")
	}

	val := lit.Value.Lexeme
	c.builder.WriteString(fmt.Sprintf("call i32 (i8*,...) @printf(i8* getelementptr ([4 x i8]* @print.str,i32 0,i32 0 %s)\n", val))
}
