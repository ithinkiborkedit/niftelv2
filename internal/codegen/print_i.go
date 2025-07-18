package codegen

import (
	"fmt"

	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	"github.com/ithinkiborkedit/niftelv2.git/internal/tokentoval"
	"github.com/ithinkiborkedit/niftelv2.git/internal/value"
)

type printableLiteralExp struct {
	lit *ast.LiteralExpr
}

type printableVariableExpr struct {
	varExpr *ast.VariableExpr
}

type Printable interface {
	EmitPrint(c *Codegen) error
}

func (p *printableLiteralExp) EmitPrint(c *Codegen) error {
	val, err := tokentoval.Convert(p.lit.Value)
	if err != nil {
		fmt.Printf("error converting token to value %v\n", err)
		return err
	}
	llvmLiteral := c.emitValueLiteral(val)
	var formatName string
	switch val.Type {

	case value.ValueInt:
		formatName = "@print_int_format"
	case value.ValueFloat:
		formatName = "@print_float_format"
	case value.ValueString:
		formatName = "@print_str_format"
	default:
		fmt.Printf("usupported literal type in print %v", val.Type)
	}
	c.builder.WriteString(fmt.Sprintf(
		"call i32 (i8*,...) @printf(i8* getelementptr ([4 x i8], [4 x i8]* %s, i32 0, i32 0), %s)",
		formatName, llvmLiteral))
	return nil
}

func (p *printableVariableExpr) EmitPrint(c *Codegen) error {
	llvmVarInfo, ok := c.lookupVariable(p.varExpr.Name.Lexeme)
	if !ok {
		return fmt.Errorf("variable lookup failed: %v", llvmVarInfo)
	}

	var formatName string
	switch llvmVarInfo.LLVMType {
	case "i64":
		formatName = "@print_int_format"
	case "double":
		formatName = "@print_float_format"
	case "i8*":
		formatName = "@print_str_format"
	default:
		return fmt.Errorf("unsupported variable type %s for print", llvmVarInfo.LLVMType)
	}

	c.builder.WriteString(fmt.Sprintf(
		"call i32 (i8*,...) @printf(i8* getelementptr ([4 x i8], [4 x i8]* %s, i32 0, i32 0), %s %s)",
		formatName, llvmVarInfo.LLVMType, llvmVarInfo.LLVMName))

	return nil

}
