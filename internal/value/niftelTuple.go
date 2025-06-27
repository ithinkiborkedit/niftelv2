package value

type NiftelTupleType struct {
	ElementTypes []*TypeInfo
}

type NiftelTupleValue struct {
	Type     *NiftelTupleType
	Elements []Value
}

func NewTupleType(elemsType []*TypeInfo) *NiftelTupleType {
	return &NiftelTupleType{ElementTypes: elemsType}
}

func NewTupleValue(tupleType *NiftelTupleType, elements []Value) *NiftelTupleValue {
	return &NiftelTupleValue{
		Type:     tupleType,
		Elements: elements,
	}
}
