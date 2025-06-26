package value

import (
	"fmt"
	"hash/fnv"
	"math"
	"reflect"
	"strings"
)

type ValueType int

const (
	ValueNull ValueType = iota
	ValueInt
	ValueFloat
	ValueString
	ValueBool
	ValueList
	ValueDict
	ValueStruct
	ValueFunc
)

type Value struct {
	Type ValueType
	Data interface{}
	Meta *TypeInfo
}

func (v Value) String() string {
	switch v.Type {
	case ValueNull:
		return "null"
	case ValueInt, ValueFloat, ValueBool, ValueString:
		return fmt.Sprintf("%v", v.Data)
	case ValueList:
		return formatList(v.Data)
	case ValueDict:
		if dict, ok := v.Data.(*NiftelDict); ok {
			items := []string{}
			for _, bucket := range dict.buckets {
				for _, entry := range bucket {
					items = append(items, fmt.Sprintf("%v: %v", entry.Key.String(), entry.Value.String()))
				}
			}
			return "{" + strings.Join(items, ", ") + "}"
		}
		return "<dict>"
	case ValueStruct:
		inst, ok := v.Data.(*StructInstance)
		if !ok {
			return "<struct-corrupt>"
		}
		fields := []string{}
		for _, fnameTok := range inst.Type.Fields {
			fname := fnameTok.Lexeme
			fields = append(fields, fmt.Sprintf("%s: %v", fname, inst.Fields[fname].String()))
		}
		return fmt.Sprintf("%s{%s}", inst.Type.Name, strings.Join(fields, ", "))
	default:
		return fmt.Sprintf("<unknown value: %v>", reflect.TypeOf(v.Data))
	}
}

func formatList(data interface{}) string {
	lst, ok := data.([]Value)
	if !ok {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteString("[")
	for i, v := range lst {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v.String())
	}
	sb.WriteString("]")
	return sb.String()
}

// func formatDict(data interface{}) string {
// 	dict, ok := data.(map[string]Value)
// }

func Null() Value {
	return Value{Type: ValueNull, Data: nil}
}

func (v Value) IsNull() bool {
	return v.Type == ValueNull
}

func (v Value) Equals(other Value) bool {
	if v.Type != other.Type {
		return false
	}

	switch v.Type {
	case ValueInt, ValueBool, ValueString, ValueFloat:
		return v.Data == other.Data
	default:
		return false
	}
}

func (v Value) Hash() uint64 {
	switch v.Type {
	case ValueInt:
		return uint64(v.Data.(float64))
	case ValueBool:
		if v.Data.(bool) {
			return 1
		}
		return 0
	case ValueString:
		h := fnv.New64()
		h.Write([]byte(v.Data.(string)))
		return h.Sum64()
	case ValueFloat:
		bits := math.Float64bits(v.Data.(float64))
		return bits
	default:
		return 0
	}
}
