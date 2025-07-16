package codegen

type structTypes map[string]*StructTypeInfo

type StructTypeInfo struct {
	Name         string
	FieldTypes   []string
	FieldNames   []string
	FieldIndices map[string]int
	Emitted      bool
}
