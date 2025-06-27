package value

type NiftelTupleType struct {
	ElementTypes []*TypeInfo
}

type NiftelTupleValue struct {
	Type     *TypeInfo
	Elements []Value
}

func NewTupleType(elemsType []*TypeInfo) *NiftelTupleType {
	return &NiftelTupleType{ElementTypes: elemsType}
}

func NewTupleValue(tupleType *TypeInfo, elements []Value) *NiftelTupleValue {
	return &NiftelTupleValue{
		Type:     tupleType,
		Elements: elements,
	}
}

func (tt *NiftelTupleType) TypeInfo() *TypeInfo {
	return &TypeInfo{
		Name: TupleTypeKey(tt.ElementTypes),
		Kind: TypeKindTuple,
	}
}
