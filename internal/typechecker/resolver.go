package typechecker

import (
	"fmt"

	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"

	"github.com/ithinkiborkedit/niftelv2.git/internal/symtable"
)

func ResolveTypeExpr(expr *ast.TypeExpr, symbols *symtable.SymbolTable, genericParams map[string]*symtable.TypeSymbol) (*symtable.TypeSymbol, error) {
	if expr == nil {
		return nil, nil
	}
	if param, ok := genericParams[expr.Name.Lexeme]; ok {
		return param, nil
	}
	typ, ok := symbols.ResolveType(expr.Name.Lexeme)
	if !ok {
		return nil, fmt.Errorf("unknown type '%s' ", expr.Name.Lexeme)
	}

	return typ, nil
}
