package codegen

import (
	"fmt"
	"strings"

	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	tokens "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"
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

func (c *Codegen) emitStringLiteral(s string) string {
	name := fmt.Sprintf("@.str%d", len(s))
	length := len(s) + 1

	escaped := escapeStringForLLVM(s)

	c.builder.WriteString(fmt.Sprintf("%s = private constant [%d x i8] c\"%s\\00\"\n", name, length, escaped))

	return name
}

func escapeStringForLLVM(s string) string {
	var sb strings.Builder

	for _, r := range s {
		if r == '\\' || r == '"' {
			sb.WriteRune('\\')
			sb.WriteRune(r)
		} else if r < 32 || r > 126 {
			sb.WriteString(fmt.Sprintf("\\%02X", r))
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func (c *Codegen) emitPreamble() {
	c.builder.WriteString(`declare i32 @printf(i8*,...)
	@print_str_format = constant [4 x i8] c"%s\0A\00"
	@print_int_format = constant [4 x i8] c"%d\0A\00"
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
		fmt.Println("warning: print only supports strings and ints")
		return
	}
	switch lit.Value.Type {
	case tokens.TokenNumber:
		val := lit.Value.Lexeme
		c.builder.WriteString(fmt.Sprintf("call i32 (i8*,...) @printf(i8* getelementptr ([4 x i8], [4 x i8]* @print_int_format, i32 0, i32 0), i32 %s)\n", val))
	case tokens.TokenString:
		strName := c.emitStringLiteral(lit.Value.Lexeme)
		length := len(lit.Value.Lexeme) + 1
		c.builder.WriteString(fmt.Sprintf("call i32 (i8*,...) @printf(i8* getelementptr ([4 x i8], [4 x i8]* @print_str_format, i32 0, i32 0), i8* getelmtptr ([%d x i8], [%d x i8]* %s, i32 0, i32 0))\n", length, length, strName))
	default:
		fmt.Println("warning print only works with string or ints")

	}
	// lit, ok := s.Expr.(*ast.LiteralExpr)

	// if !ok {
	// 	fmt.Println("warning: print only supports literal ints now")
	// }

	// val := lit.Value.Lexeme
	// c.builder.WriteString(fmt.Sprintf("call i32 (i8*,...) @printf(i8* getelementptr ([4 x i8], [4 x i8]* @print.str, i32 0, i32 0), i32 %s)\n", val))
}
