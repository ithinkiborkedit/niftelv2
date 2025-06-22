package environment

import (
	"testing"

	"github.com/ithinkiborkedit/niftelv2.git/internal/value"
)

func TestEnvironemtn_DefineAndGet(t *testing.T) {
	env := NewEnvironment(nil)

	v := value.Value{Type: value.ValueInt, Data: 42}
	env.Define("x", v)

	got, err := env.Get("x")
	if err != nil {
		t.Fatalf("unexpectd error: %v", err)
	}
	if got.Type != value.ValueInt || got.Data != 42 {
		t.Fatalf("expected value 42, got %v", got)
	}
}

func TestEnvironment_EnclosingScopes(t *testing.T) {
	global := NewEnvironment(nil)
	global.Define("a", value.Value{Type: value.ValueInt, Data: 1})

	local := NewEnvironment(nil)
	local.Define("b", value.Value{Type: value.ValueInt, Data: 2})

	if _, err := global.Get("a"); err != nil {
		t.Fatalf("expected to find 'a' in the global scope")
	}

	if _, err := local.Get("b"); err != nil {
		t.Fatalf("expected to find 'a' in the global scope")
	}

	if _, err := global.Get("b"); err == nil {
		t.Fatalf("expected to find undefined 'b' in the global scope")
	}

	if _, err := local.Get("a"); err == nil {
		t.Fatalf("expected to find undefined for 'a' in the local scope")
	}

}

func TestEnvironment_AssignEnclosing(t *testing.T) {
	global := NewEnvironment(nil)
	global.Define("x", value.Value{Type: value.ValueInt, Data: 1})

	local := NewEnvironment(global)

	err := local.Assign("x", value.Value{Type: value.ValueInt, Data: 30})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := global.Get("x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val, ok := got.Data.(int); !ok || val != 30 {
		t.Fatalf("expected 30 after assignment in enclosing, got %v", got.Data)
	}

}

func TestEnvironment_UndefinedVariable(t *testing.T) {
	env := NewEnvironment(nil)
	_, err := env.Get("missing")
	if err == nil {
		t.Fatalf("expected error for missing variable")
	}

	err = env.Assign("missing", value.Null())
	if err == nil {
		t.Fatalf("expected error for assigning missing variable")
	}
}
