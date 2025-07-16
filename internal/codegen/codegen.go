package codegen

import (
	"fmt"
	"strings"

	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	tokens "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"
	"github.com/ithinkiborkedit/niftelv2.git/internal/tokentoval"
)

type Codegen struct {
	builder      strings.Builder
	strings      map[string]string
	nextStrIndex int
}

func NewCodeGen() *Codegen {
	return &Codegen{
		strings: make(map[string]string),
	}
}

func (c *Codegen) GenerateLLVM(stmts []ast.Stmt) (string, error) {
	c.emitPreamble()
	for _, stmt := range stmts {
		c.collectStringsStmt(stmt)
	}
	c.emitStringConstants()
	c.emitMainFunction(stmts)
	return c.builder.String(), nil
}

func (c *Codegen) registerStringLiteral(s string) string {
	if name, ok := c.strings[s]; ok {
		return name
	}
	name := fmt.Sprintf("@.str%d", c.nextStrIndex)
	c.strings[s] = name
	c.nextStrIndex++
	return name
}

func (c *Codegen) emitStringConstants() {
	for s, name := range c.strings {
		length := len(s) + 1
		escaped := escapeStringForLLVM(s)
		c.builder.WriteString(fmt.Sprintf("%s = private constant [%d x i8] c\"%s\\00\"\n", name, length, escaped))
	}
}

func (c *Codegen) collectStringsStmt(s ast.Stmt) {
	switch stmt := s.(type) {
	case *ast.PrintStmt:
		c.collectStringsExpr(stmt.Expr)
	}
}

func (c *Codegen) collectStringsExpr(e ast.Expr) {
	switch expr := e.(type) {
	case *ast.LiteralExpr:
		if expr.Value.Type == tokens.TokenString {
			c.registerStringLiteral(expr.Value.Lexeme)
		}
	}
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
		c.builder.WriteString("\n")
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
	val, err := tokentoval.Convert(lit.Value)
	if err != nil {
		fmt.Printf("error converting token to vale %v\n", err)
		return
	}
	llvmLiteral := c.emitValueLiteral(val)
	var formatName string
	switch lit.Value.Type {
	case tokens.TokenNumber:
		formatName = "@print_int_format"
	case tokens.TokenString:
		formatName = "@print_str_format"
	default:
		fmt.Printf("usupported literal type in print %v", lit)
		return
	}
	c.builder.WriteString(fmt.Sprintf(
		"call i32 (i8*,...) @printf(i8* getelementptr ([4 x i8], [4 x i8]* %s, i32 0, i32 0), %s)",
		formatName, llvmLiteral))
}
