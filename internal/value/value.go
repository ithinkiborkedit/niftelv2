package value

import (
	"fmt"
	"reflect"
)

type ValueType int

const (
	TypeNull ValueType = iota
	TypeInt
	TypeFloat
	TypeString
	TypeBool
	TypeList
	TypeDict
	TypeStruct
)

type Value struct {
	Type ValueType
	Data interface{}
	Meta *TypeInfo
}

func (v Value) String() string {
	switch v.Type {
	case TypeNull:
		return "null"
	case TypeInt, TypeFloat, TypeBool, TypeString:
		return fmt.Sprintf("%v", v.Data)
	case TypeList:
		return fmt.Sprintf("%v", v.Data)
	case TypeDict:
		return fmt.Sprintf("%v", v.Data)
	case TypeStruct:
		return fmt.Sprintf("%v", v.Data)
	default:
		return fmt.Sprintf("<unknown value: %v>", reflect.TypeOf(v.Data))
	}
}

type TypeInfo struct {
	Name     string
	Fields   map[string]*TypeInfo
	ElemType *TypeInfo
}

func Null() Value {
	return Value{Type: TypeNull, Data: nil}
}
