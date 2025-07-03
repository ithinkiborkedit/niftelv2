package value

import (
	"errors"
	"strings"
	"sync"

	"github.com/ithinkiborkedit/niftelv2.git/internal/symtable"
)

type TypeKind int

var GenericStruct = map[string]*TypeInfo{}

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

func TupleTypeKey(elementTypes []*symtable.TypeSymbol) string {
	var parts []string
	for _, t := range elementTypes {
		parts = append(parts, t.SymName)
	}
	return "(" + strings.Join(parts, ",") + ")"
}

func GetOrRegisterTupleType(elementTypes []*symtable.TypeSymbol) *symtable.TypeSymbol {
	key := TupleTypeKey(elementTypes)
	if t, ok := GetType(key); ok {
		return t
	}
	tupleType := &symtable.TypeSymbol{
		SymName: key,
		SymKind: symtable.SymbolTypes,
		Fields:  nil,
	}
	RegisterType(key, tupleType)

	return tupleType

}

var globalRegistry = TypeRegistry{
	types: make(map[string]*TypeInfo),
}

func RegisterGenericStruct(name string, typeParams []string, fieldTemplates map[string]*TypeInfo) {
	info := &TypeInfo{
		Name:          name,
		Kind:          TypeKindStruct,
		Fields:        fieldTemplates,
		GenericParams: typeParams,
	}

	GenericStruct[name] = info
}

func InstatiateGenericStruct(name string, typeArgs []string) (*TypeInfo, error) {
	gen, ok := GenericStruct[name]
	if !ok {
		return nil, errors.New("unknown generic struct: " + name)
	}
	if len(typeArgs) != len(gen.GenericParams) {
		return nil, errors.New("wrong number of type arguments")
	}

	subst := map[string]string{}

	for i, param := range gen.GenericParams {
		subst[param] = typeArgs[i]
	}

	fields := map[string]*TypeInfo{}

	for fieldName, fieldType := range gen.Fields {
		fields[fieldName] = substituteTypeInfo(fieldType, subst)
	}
	instName := name + "[" + strings.Join(typeArgs, ",") + "]"

	info := &TypeInfo{
		Name:   instName,
		Kind:   TypeKindStruct,
		Fields: fields,
	}
	RegisterType(instName, &symtable.TypeSymbol{
		SymName: instName,
		SymKind: symtable.SymbolTypes,
	})

	return info, nil
}

func substituteTypeInfo(t *TypeInfo, subst map[string]string) *TypeInfo {
	if t == nil {
		return nil
	}
	if val, ok := subst[t.Name]; ok {
		return &TypeInfo{
			Name: val,
			Kind: TypeKindStruct,
		}

	}
	return t
}

func RegisterType(name string, sym *symtable.TypeSymbol) {
	BuiltInTypes[name] = sym
}

func GetType(name string) (*symtable.TypeSymbol, bool) {
	if sym, ok := BuiltInTypes[name]; ok {
		return sym, true
	}

	return nil, false
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

var BuiltInTypes = map[string]*symtable.TypeSymbol{}

func BuiltinTypesInit() {
	BuiltInTypes["int"] = &symtable.TypeSymbol{SymName: "int", SymKind: symtable.SymbolTypes}
	BuiltInTypes["float"] = &symtable.TypeSymbol{SymName: "float", SymKind: symtable.SymbolTypes}
	BuiltInTypes["string"] = &symtable.TypeSymbol{SymName: "string", SymKind: symtable.SymbolTypes}
	BuiltInTypes["bool"] = &symtable.TypeSymbol{SymName: "bool", SymKind: symtable.SymbolTypes}
	BuiltInTypes["null"] = &symtable.TypeSymbol{SymName: "null", SymKind: symtable.SymbolTypes}
	BuiltInTypes["tuple"] = &symtable.TypeSymbol{SymName: "tuple", SymKind: symtable.SymbolTypes}
	BuiltInTypes["list"] = &symtable.TypeSymbol{SymName: "list", SymKind: symtable.SymbolTypes}
	BuiltInTypes["struct"] = &symtable.TypeSymbol{SymName: "struct", SymKind: symtable.SymbolTypes}
	BuiltInTypes["func"] = &symtable.TypeSymbol{SymName: "func", SymKind: symtable.SymbolTypes}
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
