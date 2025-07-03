package typeenv

import (
	"fmt"

	"github.com/ithinkiborkedit/niftelv2.git/internal/symtable"
)

type TypeEnv struct {
	Params map[string]*symtable.TypeSymbol
	Parent *TypeEnv
}

func NewTypeEnv(parent *TypeEnv) *TypeEnv {
	return &TypeEnv{
		Params: make(map[string]*symtable.TypeSymbol),
		Parent: parent,
	}
}

func (e *TypeEnv) DefineTypeParam(name string, sym *symtable.TypeSymbol) error {
	if _, exists := e.Params[name]; exists {
		return fmt.Errorf("type parameter '%s' already defined in this scope", name)
	}
	e.Params[name] = sym
	return nil
}

func (e *TypeEnv) LookupTypeParam(name string) (*symtable.TypeSymbol, bool) {
	for curr := e; curr != nil; curr = curr.Parent {
		if sym, ok := curr.Params[name]; ok {
			return sym, true
		}
	}
	return nil, false
}

func (e *TypeEnv) HasLocalTypeparam(name string) bool {
	_, ok := e.Params[name]
	return ok
}

func (e *TypeEnv) Dump() {
	for curr := e; curr != nil; curr = curr.Parent {
		for k := range curr.Params {
			fmt.Printf("TypeParam: %s\n", k)
		}
	}
}
