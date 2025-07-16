package codegen

import (
	"fmt"

	"github.com/ithinkiborkedit/niftelv2.git/internal/value"
)

func (c *Codegen) llvmTypeForValueType(vt value.ValueType) string {
	switch vt {
	case value.ValueInt:
		return "i64"
	case value.ValueFloat:
		return "double"
	case value.ValueBool:
		return "i1"
	case value.ValueString:
		return "i8*"
	default:
		return "void"
	}
}

func (c *Codegen) emitValueLiteral(v value.Value) string {
	llvmType := c.llvmTypeForValueType(v.Type)
	switch v.Type {
	case value.ValueInt:
		intVal := int64(v.Data.(float64))
		return fmt.Sprintf("%s %d", llvmType, intVal)
	case value.ValueFloat:
		floatVal := v.Data.(float64)
		return fmt.Sprintf("%s %f", llvmType, floatVal)
	case value.ValueBool:
		boolVal := 0
		if v.Data.(bool) {
			boolVal = 1
		}
		return fmt.Sprintf("%s %d", llvmType, boolVal)
	case value.ValueString:
		strName := c.registerStringLiteral(v.Data.(string))
		length := len(v.Data.(string)) + 1
		return fmt.Sprintf("i8* getelementptr ([%d x i8], [%d x i8]* %s, i32 0, i32 0)", length, length, strName)
	default:
		return "void"
	}
}
