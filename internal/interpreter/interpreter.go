package interpreter

import (
	"github.com/ithinkiborkedit/niftelv2.git/internal/environment"
	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	"github.com/ithinkiborkedit/niftelv2.git/internal/value"
)

type Interpreter struct {
	env *environment.Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: environment.NewEnvironment(nil),
	}
}

func (i *Interpreter) Evaluate(expr ast.Expr) (value.Value, error) {
	switch e := expr.(type) {
	case *ast.ListExpr:
		return i.VisitLiteralExpr(e)
	case *ast.
	}
}


