package value

import (
	"testing"
)

func TestValueTypeRegistry(t *testing.T) {
	BuiltinTypesInit()

	if typ, ok := GetType("int"); !ok {
		t.Fatal("builtin type int not found")
	} else if typ.Name != "int" || typ.Kind != TypeKindBuiltin {
		t.Fatal("builtin type int inof incorrect")
	}

	myStruct := &TypeInfo{
		Name: "MyStruct",
		Kind: TypeKindStruct,
		Fields: map[string]*TypeInfo{
			"foo": {Name: "int", Kind: TypeKindBuiltin},
			"bar": {Name: "string", Kind: TypeKindBuiltin},
		},
		Methods:       make(map[string]*FuncInfo),
		GenericParams: nil,
	}

	RegisterType(myStruct.Name, myStruct)

	typ, ok := GetType("MyStruct")
	if !ok {
		t.Fatalf("struct type %v not found", typ.Name)
	}
	if typ.Name != "MyStruct" || typ.Kind != TypeKindStruct {
		t.Fatalf("struct type %v info incorrect", typ.Name)
	}
	if len(typ.Fields) != 2 {
		t.Fatalf("struct %v fields count mismatch", typ.Name)
	}
	if len(typ.GenericParams) == 0 || typ.GenericParams[0] != "T" {
		t.Fatalf("struct %v generic params mismatch", typ.Name)
	}

	if !HasType("int") || HasType("MyStruct") || HasType("NonExistents") {
		t.Fatal("HasType check failed")
	}

	types := ListTypes()
	found := map[string]bool{}

	for _, name := range types {
		found[name] = true
	}

	if !found["int"] || !found["MyStruct"] {
		t.Fatal("ListTypes missing registered types")
	}

}
