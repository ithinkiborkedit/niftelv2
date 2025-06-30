package environment_test

import (
	"testing"

	"github.com/ithinkiborkedit/niftelv2.git/internal/environment"
	"github.com/ithinkiborkedit/niftelv2.git/internal/symtable"
)

func TestEnvironmentDefineAndLookup(t *testing.T) {
	env := environment.NewEnvironment(nil)

	// Define variable symbol
	varSym := &symtable.VarSymbol{SymName: "foo", SymKind: symtable.SymbolVar}
	if err := env.DefineVar(varSym); err != nil {
		t.Fatalf("unexpected error defining var: %v", err)
	}

	// Lookup variable
	v, ok := env.LookupVar("foo")
	if !ok {
		t.Fatalf("expected variable foo to be found")
	}
	if v.SymName != "foo" {
		t.Errorf("expected name foo, got %s", v.SymName)
	}

	// Duplicate variable in same scope should error
	err := env.DefineVar(varSym)
	if err == nil {
		t.Errorf("expected error defining duplicate var, got nil")
	}

	// Define function symbol
	funcSym := &symtable.FuncSymbol{SymName: "bar"}
	if err := env.DefineFunc(funcSym); err != nil {
		t.Fatalf("unexpected error defining func: %v", err)
	}
	f, ok := env.LookupFunc("bar")
	if !ok {
		t.Fatalf("expected function bar to be found")
	}
	if f.SymName != "bar" {
		t.Errorf("expected name bar, got %s", f.SymName)
	}

	// Nested scope: should inherit parent symbols for lookup, but not for define
	child := environment.NewEnvironment(env)
	if _, ok := child.LookupVar("foo"); !ok {
		t.Errorf("child should see parent var foo")
	}
	childVar := &symtable.VarSymbol{SymName: "childvar", SymKind: symtable.SymbolVar}
	if err := child.DefineVar(childVar); err != nil {
		t.Errorf("unexpected error in child DefineVar: %v", err)
	}
	if _, ok := child.LookupVar("childvar"); !ok {
		t.Errorf("expected to find childvar in child env")
	}
	if _, ok := env.LookupVar("childvar"); ok {
		t.Errorf("should not find childvar in parent env")
	}
}

func TestEnvironmentDefineTypeParam(t *testing.T) {
	env := environment.NewEnvironment(nil)
	tp := &symtable.TypeParamSymbol{SymNam: "T"}
	if err := env.DefineTypeParam(tp); err != nil {
		t.Fatalf("unexpected error defining type param: %v", err)
	}
	tpsym, ok := env.LookupTypeParam("T")
	if !ok || tpsym == nil {
		t.Fatalf("failed to lookup type param")
	}
}
