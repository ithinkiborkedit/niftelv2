package symtable

type SymbolKind int

const (
	SymbolVar SymbolKind = iota
	SymbolFunc
	SymbolType
	SymbolTypeParam
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

type VarSymbol struct {
	SymName string
	SymKind SymbolKind
	Type    *TypeSymbol
	Mutable bool
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

// type TypeSymbol struct {
// 	Name string
// }

// type VarSymbol struct {
// 	Name string
// 	Kind SymbolKind
// 	Type *TypeSymbol
// }

// type SymbolTable struct {
// 	values map[string]*VarSymbol
// 	types  map[string]*TypeSymbol
// 	parent *SymbolTable
// }

// type Scope struct {
// 	Vars       map[string]*Symbol
// 	Funcs      map[string]*Symbol
// 	Types      map[string]*Symbol
// 	TypeParams map[string]*Symbol
// 	Parent     *Scope
// 	Level      int
// }
