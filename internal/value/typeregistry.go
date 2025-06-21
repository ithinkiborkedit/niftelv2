package value

import (
	"errors"
	"sync"
)

type TypeKind int

const (
	TypeKindStruct TypeKind = iota
	TypeKindEnum
	TypeKindBuiltin
)

type TypeInfo struct {
	Name          string
	Kind          TypeKind
	Fields        map[string]*TypeInfo
	Methods       map[string]*FuncInfo
	GenericParams []string
}

type FuncInfo struct {
	Name       string
	Params     []*TypeInfo
	ReturnType *TypeInfo
}

type TypeRegistry struct {
	mu    sync.RWMutex
	types map[string]*TypeInfo
}

var globalRegistry = TypeRegistry{
	types: make(map[string]*TypeInfo),
}

func RegisterType(name string, info *TypeInfo) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.types[name] = info
}

func GetType(name string) (*TypeInfo, bool) {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	t, ok := globalRegistry.types[name]
	return t, ok
}

func HasType(name string) bool {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	_, ok := globalRegistry.types[name]
	return ok
}

func ListTypes() []string {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	names := make([]string, 0, len(globalRegistry.types))
	for name := range globalRegistry.types {
		names = append(names, name)
	}

	return names
}

func BuiltinTypesInit() {
	RegisterType("int", &TypeInfo{Name: "int", Kind: TypeKindBuiltin})
	RegisterType("float", &TypeInfo{Name: "float", Kind: TypeKindBuiltin})
	RegisterType("string", &TypeInfo{Name: "string", Kind: TypeKindBuiltin})
	RegisterType("bool", &TypeInfo{Name: "bool", Kind: TypeKindBuiltin})
	RegisterType("null", &TypeInfo{Name: "null", Kind: TypeKindBuiltin})
}

func RegisterExampleStruct() {
	myStruct := &TypeInfo{
		Name: "Mystruct",
		Kind: TypeKindStruct,
		Fields: map[string]*TypeInfo{
			"foo": {Name: "int", Kind: TypeKindBuiltin},
			"bar": {Name: "string", Kind: TypeKindBuiltin},
		},
		Methods:       make(map[string]*FuncInfo),
		GenericParams: nil,
	}

	RegisterType(myStruct.Name, myStruct)
}

func (t *TypeInfo) FieldByName(name string) (*TypeInfo, error) {
	if t.Kind != TypeKindStruct {
		return nil, errors.New("type is not a struct")
	}
	field, ok := t.Fields[name]
	if !ok {
		return nil, errors.New("field not found")
	}
	return field, nil
}
