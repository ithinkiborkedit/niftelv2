package function

import (
	"fmt"

	"github.com/ithinkiborkedit/niftelv2.git/internal/environment"
	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	"github.com/ithinkiborkedit/niftelv2.git/internal/runtimecontrol"
	"github.com/ithinkiborkedit/niftelv2.git/internal/value"
)

type Function struct {
	name       string
	params     []ast.Param
	body       *ast.BlockStmt
	env        *environment.Environment
	isNative   bool
	sourceLine int
	sourceCol  int
	nativeFunc func([]value.Value, InterpreterAPI) (value.Value, error)
}

type InterpreterAPI interface {
	PushEnv(*environment.Environment)
	PopEnv()
	GetEnv() *environment.Environment
	ExecuteBlock(*ast.BlockStmt, *environment.Environment) (ret value.Value, err error)
}

type Callable interface {
	Call(args []value.Value, interp InterpreterAPI) (value.Value, error)
	Arity() int
	Name() string
	IsNative() bool
	SourcePos() (line, col int)
}

func NewUserFunc(name string, params []ast.Param, body *ast.BlockStmt, env *environment.Environment, line, col int) *Function {
	return &Function{
		name:       name,
		params:     params,
		body:       body,
		env:        env,
		isNative:   false,
		sourceLine: line,
		sourceCol:  col,
	}

}

func NewNativeFunc(name string, fn func([]value.Value, InterpreterAPI) (value.Value, error)) *Function {
	return &Function{
		name:       name,
		isNative:   true,
		nativeFunc: fn,
	}
}

func (f *Function) Call(args []value.Value, interp InterpreterAPI) (value.Value, error) {
	if f.isNative {
		return f.nativeFunc(args, interp)
	}
	if len(args) != len(f.params) {
		return value.Null(), fmt.Errorf("function '%s': expected %d parameters got %d", f.name, len(f.params), len(args))
	}

	callEnv := environment.NewEnvironment(f.env)
	for i, param := range f.params {
		callEnv.Define(param.Name.Lexeme, args[i])
	}

	interp.PushEnv(callEnv)
	defer interp.PopEnv()

	var ret value.Value = value.Null()
	var err error

	defer func() {
		if r := recover(); r != nil {
			if rv, ok := r.(runtimecontrol.ReturnValue); ok {
				ret = rv.Value
				err = nil
			} else {
				panic(r)
			}
		}

	}()

	_, execErr := interp.ExecuteBlock(f.body, callEnv)
	if execErr != nil {
		err = execErr
	}
	return ret, err
}

func (f *Function) Arity() int            { return len(f.params) }
func (f *Function) Name() string          { return f.name }
func (f *Function) IsNative() bool        { return f.isNative }
func (f *Function) SourcePos() (int, int) { return f.sourceLine, f.sourceCol }
func (r runtimecontrol.ReturnValue) Error()
