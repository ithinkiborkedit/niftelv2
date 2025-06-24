package value

import (
	"fmt"
	"reflect"
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
		return fmt.Sprintf("%v", v.Data)
	case ValueDict:
		return fmt.Sprintf("%v", v.Data)
	case ValueStruct:
		return fmt.Sprintf("Struct{%v}", v.Data)
	default:
		return fmt.Sprintf("<unknown value: %v>", reflect.TypeOf(v.Data))
	}
}

func Null() Value {
	return Value{Type: ValueNull, Data: nil}
}
