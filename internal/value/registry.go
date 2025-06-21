package value

import "sync"

type TypeKind int

const (
	TypeStruct TypeKind = iota
	TypeEnum
	TypeBuiltin
)

type TypeRegistry struct {
	mu    sync.RWMutex
	types map[string]*TypeInfo
}

var globalRegistry = &TypeRegistry{
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

func BuiltinTypesInit() {
	RegisterType("int", &TypeInfo{Name: "int"})
	RegisterType("float", &TypeInfo{Name: "float"})
	RegisterType("string", &TypeInfo{Name: "string"})
	RegisterType("bool", &TypeInfo{Name: "bool"})
	RegisterType("null", &TypeInfo{Name: "null"})
}
