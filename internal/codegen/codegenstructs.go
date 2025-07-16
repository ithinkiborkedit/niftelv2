package codegen

type structTypes map[string]*StructTypeInfo

type StructTypeInfo struct {
	Name         string
	FieldTypes   []string
	FieldNames   []string
	FieldIndices map[string]int
	Emitted      bool
}

func NewStructTypes() structTypes {
	return make(structTypes)
}

func (st structTypes) Register(info *StructTypeInfo) {
	st[info.Name] = info
}

func (st structTypes) Lookup(name string) (*StructTypeInfo, bool) {
	info, ok := st[name]
	return info, ok
}

func (st structTypes) MarkEmitted(name string) {
	if info, ok := st[name]; ok {
		info.Emitted = true
	}
}
