package value

import (
	"testing"
)

func TestValueTypeRegistry(t *testing.T) {
	BuiltinTypesInit()

	builtins := []string{"int", "float", "string", "bool", "null"}
	registry := NewValueRegistry()
	for _, name := range builtins {
		typ, ok := GetType(name)
		if !ok {
			t.Errorf("builtin type %s not found", name)
		}
		if typ == nil || typ.Name != name {
			t.Errorf("builtin type %s: wrong typeInfo returned", name)
		}

		myStruct := &TypeInfo{Name: "MyStruct", Fields: map[string]*TypeInfo{
			"foo": &TypeInfo{Name: "int"},
			"bar": &TypeInfo{Name: "string"},
		}}

		registry.RegisterType(myStruct)
	}
}
