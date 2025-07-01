package value

import (
	"errors"
	"strings"
	"sync"
)

type TypeKind int

const (
	TypeKindStruct TypeKind = iota
	TypeKindEnum
	TypeKindBuiltin
	TypeKindTuple
	TypeKindList
	TypeKindInt
	TypeKindFloat
	TypeKindString
	TypeKindBool
	TypeKindNull
	TypeKindFunc
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

var tupleTypes = map[string]*NiftelTupleType{}

func TupleTypeKey(elementTypes []*TypeInfo) string {
	var parts []string
	for _, t := range elementTypes {
		parts = append(parts, t.Name)
	}
	return "(" + strings.Join(parts, ",") + ")"
}

func GetOrRegisterTupleType(elementTypes []*TypeInfo) *TypeInfo {
	key := TupleTypeKey(elementTypes)
	if t, ok := GetType(key); ok {
		return t
	}
	tupleType := &TypeInfo{
		Name:          key,
		Kind:          TypeKindTuple,
		Fields:        nil,
		Methods:       nil,
		GenericParams: nil,
	}
	RegisterType(key, tupleType)
	// t = NewTupleType(elementTypes)
	// tupleTypes[key] = t
	return tupleType

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

func ResetTypes() {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	globalRegistry.types = make(map[string]*TypeInfo)
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
	RegisterType("int", &TypeInfo{Name: "int", Kind: TypeKindInt})
	RegisterType("float", &TypeInfo{Name: "float", Kind: TypeKindFloat})
	RegisterType("string", &TypeInfo{Name: "string", Kind: TypeKindString})
	RegisterType("bool", &TypeInfo{Name: "bool", Kind: TypeKindBool})
	RegisterType("null", &TypeInfo{Name: "null", Kind: TypeKindNull})
	RegisterType("tuple", &TypeInfo{Name: "tuple", Kind: TypeKindTuple})
	RegisterType("list", &TypeInfo{Name: "list", Kind: TypeKindList})
	RegisterType("struct", &TypeInfo{Name: "struct", Kind: TypeKindStruct})
	RegisterType("func", &TypeInfo{Name: "fun", Kind: TypeKindFunc})
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
