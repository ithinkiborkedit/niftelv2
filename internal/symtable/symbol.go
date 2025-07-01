package symtable

import "fmt"

type SymbolKind int

const (
	SymbolVar SymbolKind = iota
	SymbolFuncs
	SymbolTypes
	SymbolTypeParams
)

type Symbol interface {
	Name() string
	Kind() SymbolKind
	String() string
}

type TypeSymbol struct {
	SymName    string
	SymKind    SymbolKind
	Fields     map[string]*TypeSymbol
	TypeParams []string
}

type TypeParamSymbol struct {
	SymNam string
}

type VarSymbol struct {
	SymName string
	SymKind SymbolKind
	Type    *TypeSymbol
	Mutable bool
}

type FuncSymbol struct {
	SymName    string
	Params     []VarSymbol
	ReturnType []*TypeSymbol
	TypeParams []string
}

func (v *VarSymbol) Name() string {
	return v.SymName
}

func (v *VarSymbol) Kind() SymbolKind {
	return v.SymKind
}

func (v *VarSymbol) String() string {
	return v.SymName
}

func (t *TypeSymbol) Name() string {
	return t.SymName
}

func (t *TypeSymbol) Kind() SymbolKind {
	return t.SymKind
}

func (t *TypeSymbol) String() string {
	return t.SymName
}

func (f *FuncSymbol) Name() string {
	return f.SymName
}

func (f *FuncSymbol) Kind() SymbolKind {
	return SymbolFuncs
}

func (f *FuncSymbol) String() string {
	return fmt.Sprintf("func %s", f.SymName)
}

func (tp *TypeParamSymbol) Name() string {
	return tp.SymNam
}

func (tp *TypeParamSymbol) Kind() SymbolKind {
	return SymbolTypeParams
}

func (tp *TypeParamSymbol) String() string {
	return fmt.Sprintf("typeparam %s", tp.SymNam)
}
