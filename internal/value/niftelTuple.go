package value

import "github.com/ithinkiborkedit/niftelv2.git/internal/symtable"

type NiftelTupleType struct {
	ElementTypes []*symtable.TypeSymbol
}

type NiftelTupleValue struct {
	Type     *symtable.TypeSymbol
	Elements []Value
}

func NewTupleType(elemsType []*symtable.TypeSymbol) *NiftelTupleType {
	return &NiftelTupleType{ElementTypes: elemsType}
}

func NewTupleValue(tupleType *symtable.TypeSymbol, elements []Value) *NiftelTupleValue {
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
