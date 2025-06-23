package interpreter

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ithinkiborkedit/niftelv2.git/internal/environment"
	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	token "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"
	"github.com/ithinkiborkedit/niftelv2.git/internal/value"
)

// Interpreter interprets and executes Niftel code.
type Interpreter struct {
	env *environment.Environment
	// Add flags, call stacks, etc. here as needed
}

// NewInterpreter returns a fresh Interpreter with a global environment.
func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: environment.NewEnvironment(nil),
	}
}

// Evaluate dispatches to the correct Expr handler.
func (i *Interpreter) Evaluate(expr ast.Expr) (value.Value, error) {
	switch e := expr.(type) {
	case *ast.LiteralExpr:
		return i.VisitLiteralExpr(e)
	case *ast.VariableExpr:
		return i.VisitVariableExpr(e)
	case *ast.BinaryExpr:
		return i.VisitBinaryExpr(e)
	case *ast.UnaryExpr:
		return i.VisitUnaryExpr(e)
	case *ast.CallExpr:
		return i.VisitCallExpr(e)
	case *ast.GetExpr:
		return i.VisitGetExpr(e)
	case *ast.ListExpr:
		return i.VisitListExpr(e)
	case *ast.DictExpr:
		return i.VisitDictExpr(e)
	case *ast.FuncExpr:
		return i.VisitFuncExpr(e)
	default:
		return value.Null(), fmt.Errorf("unknown expression type %T", expr)
	}
}

// Execute dispatches to the correct Stmt handler.
func (i *Interpreter) Execute(stmt ast.Stmt) error {
	switch s := stmt.(type) {
	case *ast.VarStmt:
		return i.VisitVarStmt(s)
	case *ast.ShortVarStmt:
		return i.VisitShortVarStmt(s)
	case *ast.AssignStmt:
		return i.VisitAssignStmt(s)
	case *ast.PrintStmt:
		return i.VisitPrintStmt(s)
	case *ast.ExprStmt:
		return i.VisitExprStmt(s)
	case *ast.IfStmt:
		return i.VisitIfStmt(s)
	case *ast.WhileStmt:
		return i.VisitWhileStmt(s)
	case *ast.ForStmt:
		return i.VisitForStmt(s)
	case *ast.FuncStmt:
		return i.VisitFuncStmt(s)
	case *ast.ReturnStmt:
		return i.VisitReturnStmt(s)
	case *ast.BreakStmt:
		return i.VisitBreakStmt(s)
	case *ast.ContinueStmt:
		return i.VisitContinueStmt(s)
	case *ast.BlockStmt:
		return i.VisitBlockStmt(s)
	default:
		return fmt.Errorf("unknown statement type %T", stmt)
	}
}

// --- Expression Visitors ---

func (i *Interpreter) VisitLiteralExpr(expr *ast.LiteralExpr) (value.Value, error) {
	tok := expr.Value
	switch tok.Type {
	case token.TokenNumber:
		intVal, err := strconv.Atoi(tok.Lexeme)
		if err != nil {
			return value.Null(), errors.New("invalid int literal token data")
		}
		return value.Value{Type: value.ValueInt, Data: intVal}, nil

	case token.TokenString:
		return value.Value{Type: value.ValueString, Data: tok.Lexeme}, nil
	case token.TokenBool:
		switch tok.Lexeme {
		case "true":
			return value.Value{Type: value.ValueBool, Data: true}, nil
		case "false":
			return value.Value{Type: value.ValueBool, Data: false}, nil
		}
		return value.Null(), errors.New("invalid bool literal token")
	case token.TokenNull:
		return value.Null(), nil
	default:
		return value.Null(), fmt.Errorf("unsupported literal token type %v", tok.Type)
	}
}

func (i *Interpreter) VisitVariableExpr(expr *ast.VariableExpr) (value.Value, error) {
	val, err := i.env.Get(expr.Name.Lexeme)
	if err != nil {
		return value.Null(), fmt.Errorf("undefined variable %s", expr.Name.Lexeme)
	}
	return val, nil
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.BinaryExpr) (value.Value, error) {
	left, err := i.Evaluate(expr.Left)
	if err != nil {
		return value.Null(), err
	}
	right, err := i.Evaluate(expr.Right)
	if err != nil {
		return value.Null(), err
	}

	switch expr.Operator.Type {
	case token.TokenPlus:
		if left.Type == value.ValueInt && right.Type == value.ValueInt {
			return value.Value{Type: value.ValueInt, Data: left.Data.(int64) + right.Data.(int64)}, nil
		}
		if left.Type == value.ValueString && right.Type == value.ValueString {
			return value.Value{Type: value.ValueString, Data: left.Data.(string) + right.Data.(string)}, nil
		}
		return value.Null(), fmt.Errorf("unsupported operand types for +: %v and %v", left.Type, right.Type)
	case token.TokenMinus:
		if left.Type == value.ValueInt && right.Type == value.ValueInt {
			return value.Value{Type: value.ValueInt, Data: left.Data.(int64) - right.Data.(int64)}, nil
		}
		return value.Null(), fmt.Errorf("unsupported operand types for -: %v and %v", left.Type, right.Type)
	case token.TokenStar:
		if left.Type == value.ValueInt && right.Type == value.ValueInt {
			return value.Value{Type: value.ValueInt, Data: left.Data.(int64) * right.Data.(int64)}, nil
		}
		return value.Null(), fmt.Errorf("unsupported operand types for *: %v and %v", left.Type, right.Type)
	case token.TokenFWDSlash:
		if left.Type == value.ValueInt && right.Type == value.ValueInt {
			if right.Data.(int64) == 0 {
				return value.Null(), fmt.Errorf("division by zero")
			}
			return value.Value{Type: value.ValueInt, Data: left.Data.(int64) / right.Data.(int64)}, nil
		}
		return value.Null(), fmt.Errorf("unsupported operand types for /: %v and %v", left.Type, right.Type)
	case token.TokenGreater:
		if left.Type == value.ValueInt && right.Type == value.ValueInt {
			return value.Value{Type: value.ValueBool, Data: left.Data.(int64) > right.Data.(int64)}, nil
		}
		return value.Null(), fmt.Errorf("unsupported operand types for >: %v and %v", left.Type, right.Type)
	case token.TokenLess:
		if left.Type == value.ValueInt && right.Type == value.ValueInt {
			return value.Value{Type: value.ValueBool, Data: left.Data.(int64) < right.Data.(int64)}, nil
		}
		return value.Null(), fmt.Errorf("unsupported operand types for <: %v and %v", left.Type, right.Type)
	case token.TokenEqality:
		if left.Type != right.Type {
			return value.Value{Type: value.ValueBool, Data: false}, nil
		}
		return value.Value{Type: value.ValueBool, Data: left.Data == right.Data}, nil
	case token.TokenBangEqal:
		if left.Type != right.Type {
			return value.Value{Type: value.ValueBool, Data: true}, nil
		}
		return value.Value{Type: value.ValueBool, Data: left.Data != right.Data}, nil
	default:
		return value.Null(), fmt.Errorf("unsupported binary operator %v", expr.Operator.Lexeme)
	}
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.UnaryExpr) (value.Value, error) {
	right, err := i.Evaluate(expr.Right)
	if err != nil {
		return value.Null(), err
	}
	switch expr.Operator.Type {
	case token.TokenBang:
		if right.Type != value.ValueBool {
			return value.Null(), fmt.Errorf("operator ! requires boolean operand")
		}
		return value.Value{Type: value.ValueBool, Data: !right.Data.(bool)}, nil
	case token.TokenMinus:
		if right.Type != value.ValueInt {
			return value.Null(), fmt.Errorf("operator - requires int operand")
		}
		return value.Value{Type: value.ValueInt, Data: -right.Data.(int64)}, nil
	default:
		return value.Null(), fmt.Errorf("unsupported unary operator %v", expr.Operator.Lexeme)
	}
}

// Implement remaining expression visitors: VisitCallExpr, VisitGetExpr, VisitListExpr, VisitDictExpr, VisitFuncExpr

// --- Statement Visitors ---

func (i *Interpreter) VisitVarStmt(stmt *ast.VarStmt) error {
	val, err := i.Evaluate(stmt.Init)
	if err != nil {
		return err
	}
	i.env.Define(stmt.Name.Lexeme, val)
	return nil
}

func (i *Interpreter) VisitShortVarStmt(stmt *ast.ShortVarStmt) error {
	val, err := i.Evaluate(stmt.Init)
	if err != nil {
		return err
	}
	i.env.Define(stmt.Name.Lexeme, val)
	return nil
}

func (i *Interpreter) VisitAssignStmt(stmt *ast.AssignStmt) error {
	val, err := i.Evaluate(stmt.Value)
	if err != nil {
		return err
	}
	return i.env.Assign(stmt.Name.Lexeme, val)
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) error {
	val, err := i.Evaluate(stmt.Expr)
	if err != nil {
		return err
	}
	fmt.Println(val.String())
	return nil
}

func (i *Interpreter) VisitExprStmt(stmt *ast.ExprStmt) error {
	_, err := i.Evaluate(stmt.Expr)
	return err
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) error {
	cond, err := i.Evaluate(stmt.Conditon)
	if err != nil {
		return err
	}
	if cond.Type != value.ValueBool {
		return fmt.Errorf("if condition must evaluate to bool")
	}
	if cond.Data.(bool) {
		return i.Execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.Execute(stmt.ElseBranch)
	}
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *ast.WhileStmt) error {
	for {
		cond, err := i.Evaluate(stmt.Conditon)
		if err != nil {
			return err
		}
		if cond.Type != value.ValueBool {
			return fmt.Errorf("while condition must evaluate to bool")
		}
		if !cond.Data.(bool) {
			break
		}
		if err := i.Execute(stmt.Body); err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) error {
	previousEnv := i.env
	i.env = environment.NewEnvironment(previousEnv)
	defer func() { i.env = previousEnv }()

	for _, s := range stmt.Statements {
		if err := i.Execute(s); err != nil {
			return err
		}
	}
	return nil
}

// VisitForStmt executes a for loop.
func (i *Interpreter) VisitForStmt(stmt *ast.ForStmt) error {
	previousEnv := i.env
	i.env = environment.NewEnvironment(previousEnv)
	defer func() { i.env = previousEnv }()

	// Init statement
	if stmt.Init != nil {
		if err := i.Execute(stmt.Init); err != nil {
			return err
		}
	}

	for {
		// Condition expression
		if stmt.CondExpr != nil {
			condVal, err := i.Evaluate(stmt.CondExpr)
			if err != nil {
				return err
			}
			if condVal.Type != value.ValueBool {
				return fmt.Errorf("for loop condition must evaluate to bool")
			}
			if !condVal.Data.(bool) {
				break
			}
		}

		// Execute body
		err := i.Execute(stmt.BodyStmt)
		if err != nil {
			// Implement break/continue handling here if you add flags or custom errors
			return err
		}

		// Update statement
		if stmt.Update != nil {
			if err := i.Execute(stmt.Update); err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *Interpreter) VisitCallExpr(expr *ast.CallExpr) (value.Value, error) {
	// Evaluate the callee expression (should be a function)
	calleeVal, err := i.Evaluate(expr.Callee)
	if err != nil {
		return value.Null(), err
	}

	// Evaluate all argument expressions
	args := make([]value.Value, len(expr.Arguments))
	for idx, argExpr := range expr.Arguments {
		args[idx], err = i.Evaluate(argExpr)
		if err != nil {
			return value.Null(), err
		}
	}

	// Check if callee is a function value (assumed to have type ValueStruct or special FuncInfo)
	// You may need a function wrapper or closure struct; adjust as per your implementation
	// funcMeta := calleeVal.Meta
	// if funcMeta == nil || funcMeta.Kind != value.TypeFunc {
	// 	return value.Null(), fmt.Errorf("attempt to call non-function value")
	// }

	// TODO: Implement function invocation logic (bytecode, native call, etc.)
	// For now, assume native Go function stored in calleeVal.Data as func([]value.Value) (value.Value, error)
	if nativeFunc, ok := calleeVal.Data.(func([]value.Value) (value.Value, error)); ok {
		return nativeFunc(args)
	}

	return value.Null(), fmt.Errorf("function call not implemented for this callee type")
}

func (i *Interpreter) VisitIndexExpr(expr *ast.IndexExpr) (value.Value, error) {
	// Evaluate the collection expression
	collectionVal, err := i.Evaluate(expr.Collection)
	if err != nil {
		return value.Null(), err
	}

	// Evaluate the index/key expression
	indexVal, err := i.Evaluate(expr.Index)
	if err != nil {
		return value.Null(), err
	}

	switch collectionVal.Type {
	case value.ValueList:
		list, ok := collectionVal.Data.([]value.Value)
		if !ok {
			return value.Null(), fmt.Errorf("list data is corrupted")
		}
		idx, ok := indexVal.Data.(int)
		if !ok {
			return value.Null(), fmt.Errorf("list index must be integer")
		}
		if idx < 0 || idx >= len(list) {
			return value.Null(), fmt.Errorf("list index out of range")
		}
		return list[idx], nil

	case value.ValueDict:
		dict, ok := collectionVal.Data.(map[string]value.Value)
		if !ok {
			return value.Null(), fmt.Errorf("dict data is corrupted")
		}
		keyStr, ok := indexVal.Data.(string)
		if !ok {
			return value.Null(), fmt.Errorf("dict key must be string")
		}
		val, exists := dict[keyStr]
		if !exists {
			return value.Null(), fmt.Errorf("dict key not found: %s", keyStr)
		}
		return val, nil

	default:
		return value.Null(), fmt.Errorf("indexing unsupported on type %v", collectionVal.Type)
	}
}

func (i *Interpreter) VisitGetExpr(expr *ast.GetExpr) (value.Value, error) {
	// Evaluate object expression
	objectVal, err := i.Evaluate(expr.Object)
	if err != nil {
		return value.Null(), err
	}

	// Only structs support property access
	if objectVal.Type != value.ValueStruct {
		return value.Null(), fmt.Errorf("attempt to get property on non-struct type")
	}

	// Access struct field by name (string from token)
	fieldName := expr.Name.Lexeme

	// Assume struct data stored as map[string]value.Value
	fields, ok := objectVal.Data.(map[string]value.Value)
	if !ok {
		return value.Null(), fmt.Errorf("struct data corrupted")
	}

	val, exists := fields[fieldName]
	if !exists {
		return value.Null(), fmt.Errorf("struct field '%s' not found", fieldName)
	}
	return val, nil
}

func (i *Interpreter) VisitListExpr(expr *ast.ListExpr) (value.Value, error) {
	elements := make([]value.Value, 0, len(expr.Elements))
	for _, elementExpr := range expr.Elements {
		val, err := i.Evaluate(elementExpr)
		if err != nil {
			return value.Null(), err
		}
		elements = append(elements, val)
	}
	return value.Value{
		Type: value.ValueList,
		Data: elements,
	}, nil
}

func (i *Interpreter) VisitDictExpr(expr *ast.DictExpr) (value.Value, error) {
	dict := make(map[string]value.Value, len(expr.Pairs))
	for _, pair := range expr.Pairs {
		keyVal, err := i.Evaluate(pair[0])
		if err != nil {
			return value.Null(), err
		}
		keyStr, ok := keyVal.Data.(string)
		if !ok {
			return value.Null(), fmt.Errorf("dictionary key must be a string")
		}

		valueVal, err := i.Evaluate(pair[1])
		if err != nil {
			return value.Null(), err
		}
		dict[keyStr] = valueVal
	}
	return value.Value{
		Type: value.ValueDict,
		Data: dict,
	}, nil
}

func (i *Interpreter) VisitFuncExpr(expr *ast.FuncExpr) (value.Value, error) {
	// Create a callable function closure value from this FuncExpr AST node

	// For now, pack the FuncExpr itself as data and keep meta for type info
	// Actual call logic to be implemented in VisitCallExpr

	return value.Value{
		Type: value.ValueStruct, // Or a dedicated ValueFunc type if you have one
		Data: expr,
		Meta: nil, // Optional: assign Func type info if available
	}, nil
}

// VisitFuncStmt defines a function in the environment.
func (i *Interpreter) VisitFuncStmt(stmt *ast.FuncStmt) error {
	// Wrap the function expression in a callable Value
	fnValue := value.Value{
		Type: value.ValueStruct, // or define a specific function value type
		Data: stmt,              // Keep the AST node to execute on call
		Meta: nil,
	}
	i.env.Define(stmt.Name.Lexeme, fnValue)
	return nil
}

// VisitReturnStmt returns a value from a function.
// Here, we use panic as control flow to unwind the call stack for returns.
func (i *Interpreter) VisitReturnStmt(stmt *ast.ReturnStmt) error {
	var retVal value.Value
	var err error
	if stmt.Value != nil {
		retVal, err = i.Evaluate(stmt.Value)
		if err != nil {
			return err
		}
	} else {
		retVal = value.Null()
	}
	// panic with a special ReturnValue to unwind execution
	panic(ReturnValue{Value: retVal})
}

// VisitBreakStmt handles break statement in loops.
func (i *Interpreter) VisitBreakStmt(stmt *ast.BreakStmt) error {
	panic(BreakSignal{})
}

// VisitContinueStmt handles continue statement in loops.
func (i *Interpreter) VisitContinueStmt(stmt *ast.ContinueStmt) error {
	panic(ContinueSignal{})
}

// ReturnValue is used internally to signal return from function.
type ReturnValue struct {
	Value value.Value
}

func (r ReturnValue) Error() string { return "function return" }

// BreakSignal signals a break in loops.
type BreakSignal struct{}

func (b BreakSignal) Error() string { return "break" }

// ContinueSignal signals a continue in loops.
type ContinueSignal struct{}

func (c ContinueSignal) Error() string { return "continue" }
