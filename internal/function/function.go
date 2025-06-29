package function

import (
	"fmt"

	"github.com/ithinkiborkedit/niftelv2.git/internal/controlflow"
	"github.com/ithinkiborkedit/niftelv2.git/internal/environment"
	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
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
	nativeFunc func([]value.Value, InterpreterAPI) controlflow.ExecResult
}

type InterpreterAPI interface {
	PushEnv(*environment.Environment)
	PopEnv()
	GetEnv() *environment.Environment
	ExecuteBlock(*ast.BlockStmt, *environment.Environment) controlflow.ExecResult
	// ExecuteBlock(*ast.BlockStmt, *environment.Environment) (ret value.Value, err error)
}

type Callable interface {
	Call(args []value.Value, interp InterpreterAPI) controlflow.ExecResult
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

func NewNativeFunc(name string, fn func([]value.Value, InterpreterAPI) controlflow.ExecResult) *Function {
	return &Function{
		name:       name,
		isNative:   true,
		nativeFunc: fn,
	}
}

func (f *Function) Call(args []value.Value, interp InterpreterAPI) controlflow.ExecResult {
	fmt.Printf("CALLING Function: %v, args: %v", f.name, f.params)
	if f.isNative {
		return f.nativeFunc(args, interp)
	}
	if len(args) != len(f.params) {
		return controlflow.ExecResult{Err: fmt.Errorf("function '%s': expected %d parameters got %d", f.name, len(f.params), len(args))}
	}

	callEnv := environment.NewEnvironment(f.env)
	for i, param := range f.params {
		callEnv.Define(param.Name.Lexeme, args[i])
	}

	interp.PushEnv(callEnv)
	defer interp.PopEnv()
	// fmt.Printf("RETURNING from function.Call: %#v, err: %v\n", controlflow.ExecResult{Value: })
	execResult := interp.ExecuteBlock(f.body, callEnv)
	return execResult

	// var ret value.Value = value.Null()
	// var err error

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Printf("Function CALLER caught panic: %#v\n", r)
	// 		if rv, ok := r.(runtimecontrol.ReturnValue); ok {
	// 			fmt.Printf("Function returning %v\n", rv.Value)
	// 			ret = rv.Value
	// 			err = nil
	// 		} else {
	// 			panic(r)
	// 		}
	// 	}

	// }()
	// ret, err = interp.ExecuteBlock(f.body, callEnv)
	// _, execErr := interp.ExecuteBlock(f.body, callEnv)
	// if execErr != nil {
	// 	return value.Null(), execErr
}

// return value.Null(), nil

func (f *Function) Arity() int            { return len(f.params) }
func (f *Function) Name() string          { return f.name }
func (f *Function) IsNative() bool        { return f.isNative }
func (f *Function) SourcePos() (int, int) { return f.sourceLine, f.sourceCol }
