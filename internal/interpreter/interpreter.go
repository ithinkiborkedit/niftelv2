package interpreter

import (
	"errors"
	"fmt"

	"github.com/ithinkiborkedit/niftelv2.git/internal/controlflow"
	"github.com/ithinkiborkedit/niftelv2.git/internal/environment"
	"github.com/ithinkiborkedit/niftelv2.git/internal/function"
	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	token "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"
	"github.com/ithinkiborkedit/niftelv2.git/internal/value"
)

// Interpreter interprets and executes Niftel code.
type Interpreter struct {
	env                *environment.Environment
	envStack           []*environment.Environment
	ShouldPrintResults bool
	// Add flags, call stacks, etc. here as needed
}

// NewInterpreter returns a fresh Interpreter with a global environment.
func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: environment.NewEnvironment(nil),
	}
}

func (i *Interpreter) Eval(expr ast.Expr) (value.Value, error) {
	return i.Evaluate(expr)
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
	case *ast.StructLiteralExpr:
		return i.VisitStructLiteralExpr(e)
	case *ast.FuncExpr:
		return i.VisitFuncExpr(e)
	default:
		return value.Null(), fmt.Errorf("unknown expression type %T", expr)
	}
}

// Execute dispatches to the correct Stmt handler.
func (i *Interpreter) Execute(stmt ast.Stmt) controlflow.ExecResult {
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
	case *ast.StructStmt:
		return i.VisitStructStmt(s)
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
		return controlflow.ExecResult{Err: fmt.Errorf("unknown statement type %T", stmt)}
	}
}

// --- Expression Visitors ---

func (i *Interpreter) VisitStructLiteralExpr(expr *ast.StructLiteralExpr) (value.Value, error) {
	structName := expr.TypeName.Lexeme

	typeInfo, ok := value.GetType(structName)
	if !ok || typeInfo.Kind != value.TypeKindStruct {
		return value.Null(), fmt.Errorf("struct tpye '%s' not found", structName)
	}

	structType := &value.StructType{
		Name:   typeInfo.Name,
		Fields: make([]token.Token, 0, len(typeInfo.Fields)),
	}
	for fname := range typeInfo.Fields {
		structType.Fields = append(structType.Fields, token.Token{Lexeme: fname})
	}

	instance := &value.StructInstance{
		Type:   structType,
		Fields: make(map[string]value.Value),
	}

	for fname, exprval := range expr.Fields {
		val, err := i.Evaluate(exprval)
		if err != nil {
			return value.Null(), fmt.Errorf("error in field '%s': %w", fname, err)
		}
		instance.Fields[fname] = val
	}

	for fname := range typeInfo.Fields {
		if _, ok := instance.Fields[fname]; !ok {
			instance.Fields[fname] = value.Null()
		}
	}

	return value.Value{
		Type: value.ValueStruct,
		Data: instance,
		Meta: typeInfo,
	}, nil
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.LiteralExpr) (value.Value, error) {
	tok := expr.Value
	switch tok.Type {
	case token.TokenNumber:
		switch val := tok.Data.(type) {
		case int:
			return value.Value{Type: value.ValueInt, Data: float64(val)}, nil
		case int64:
			return value.Value{Type: value.ValueInt, Data: float64(val)}, nil
		case float64:
			return value.Value{Type: value.ValueInt, Data: val}, nil
		default:
			return value.Null(), errors.New("invalid int literal token data")

		}

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
			return value.Value{Type: value.ValueInt, Data: left.Data.(float64) + right.Data.(float64)}, nil
		}
		if left.Type == value.ValueString && right.Type == value.ValueString {
			return value.Value{Type: value.ValueString, Data: left.Data.(string) + right.Data.(string)}, nil
		}
		return value.Null(), fmt.Errorf("unsupported operand types for +: %v and %v", left.Type, right.Type)
	case token.TokenMinus:
		if left.Type == value.ValueInt && right.Type == value.ValueInt {
			return value.Value{Type: value.ValueInt, Data: left.Data.(float64) - right.Data.(float64)}, nil
		}
		return value.Null(), fmt.Errorf("unsupported operand types for -: %v and %v", left.Type, right.Type)
	case token.TokenStar:
		if left.Type == value.ValueInt && right.Type == value.ValueInt {
			return value.Value{Type: value.ValueInt, Data: left.Data.(float64) * right.Data.(float64)}, nil
		}
		return value.Null(), fmt.Errorf("unsupported operand types for *: %v and %v", left.Type, right.Type)
	case token.TokenFWDSlash:
		if left.Type == value.ValueInt && right.Type == value.ValueInt {
			if right.Data.(float64) == 0 {
				return value.Null(), fmt.Errorf("division by zero")
			}
			return value.Value{Type: value.ValueInt, Data: left.Data.(float64) / right.Data.(float64)}, nil
		}
		return value.Null(), fmt.Errorf("unsupported operand types for /: %v and %v", left.Type, right.Type)
	case token.TokenGreater:
		if left.Type == value.ValueInt && right.Type == value.ValueInt {
			return value.Value{Type: value.ValueBool, Data: left.Data.(float64) > right.Data.(float64)}, nil
		}
		return value.Null(), fmt.Errorf("unsupported operand types for >: %v and %v", left.Type, right.Type)
	case token.TokenLess:
		if left.Type == value.ValueInt && right.Type == value.ValueInt {
			return value.Value{Type: value.ValueBool, Data: left.Data.(float64) < right.Data.(float64)}, nil
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
		return value.Value{Type: value.ValueInt, Data: -right.Data.(float64)}, nil
	default:
		return value.Null(), fmt.Errorf("unsupported unary operator %v", expr.Operator.Lexeme)
	}
}

// Implement remaining expression visitors: VisitCallExpr, VisitGetExpr, VisitListExpr, VisitDictExpr, VisitFuncExpr

// --- Statement Visitors ---

func (i *Interpreter) VisitStructStmt(stmt *ast.StructStmt) controlflow.ExecResult {
	structName := stmt.Name.Lexeme

	if value.HasType(structName) {
		return controlflow.ExecResult{Err: fmt.Errorf("struct '%s' already defined", structName)}
	}

	fields := make(map[string]*value.TypeInfo)
	for _, field := range stmt.Fields {
		if len(field.Names) != 1 {
			return controlflow.ExecResult{Err: fmt.Errorf("struct field must have exactly one name!")}
		}
		fieldName := field.Names[0].Lexeme
		fieldTypeName := field.Type.Lexeme
		fieldType, ok := value.GetType(fieldTypeName)
		if !ok {
			return controlflow.ExecResult{Err: fmt.Errorf("Uknown type '%s' for struct field '%s'", fieldTypeName, fieldName)}
		}
		fields[fieldName] = fieldType
	}
	methods := make(map[string]*value.FuncInfo)
	for _, method := range stmt.Methods {
		methods[method.Name.Lexeme] = &value.FuncInfo{
			Name: method.Name.Lexeme,
		}
	}

	typeInfo := &value.TypeInfo{
		Name:    structName,
		Kind:    value.TypeKindStruct,
		Fields:  fields,
		Methods: methods,
	}
	value.RegisterType(structName, typeInfo)
	fmt.Printf("[INFO] REGISTERED struct type: '%s'\n", stmt.Name.Lexeme)
	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
}

func (i *Interpreter) VisitVarStmt(stmt *ast.VarStmt) controlflow.ExecResult {
	val, err := i.Evaluate(stmt.Init)
	if err != nil {
		return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
	}
	if len(stmt.Names) == 1 {
		i.env.Define(stmt.Names[0].Lexeme, val)
		return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
	}

	values, ok := val.Data.([]value.Value)
	if !ok {
		return controlflow.ExecResult{Err: fmt.Errorf("cannot unpack non-tuple value to multiple variables")}
	}
	if len(values) != len(stmt.Names) {
		return controlflow.ExecResult{Err: fmt.Errorf("mismatch: %d variables but %d values returned", len(stmt.Names), len(values))}
	}
	for idx, name := range stmt.Names {
		i.env.Define(name.Lexeme, values[idx])
	}
	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
}

func (i *Interpreter) VisitShortVarStmt(stmt *ast.ShortVarStmt) controlflow.ExecResult {
	val, err := i.Evaluate(stmt.Init)
	if err != nil {
		return controlflow.ExecResult{Err: err}
	}
	i.env.Define(stmt.Name.Lexeme, val)
	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
}

func (i *Interpreter) VisitAssignStmt(stmt *ast.AssignStmt) controlflow.ExecResult {
	val, err := i.Evaluate(stmt.Value)
	if err != nil {
		return controlflow.ExecResult{Err: err}
	}
	err = i.env.Assign(stmt.Name.Lexeme, val)
	if err != nil {
		return controlflow.ExecResult{Err: err}
	}
	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) controlflow.ExecResult {
	val, err := i.Evaluate(stmt.Expr)
	if err != nil {
		return controlflow.ExecResult{Err: err}
	}
	fmt.Println(val.String())
	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
}

func (i *Interpreter) VisitExprStmt(stmt *ast.ExprStmt) controlflow.ExecResult {
	result, err := i.Evaluate(stmt.Expr)
	if err != nil {
		return controlflow.ExecResult{Err: err}
	}
	if i.ShouldPrintResults && !result.IsNull() {
		fmt.Println(result.String())
	}

	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) controlflow.ExecResult {
	cond, err := i.Evaluate(stmt.Conditon)
	if err != nil {
		return controlflow.ExecResult{Err: err}
	}
	fmt.Printf("If conditon %#v (type=%v)\n", cond, cond.Type)
	if cond.Type != value.ValueBool {
		return controlflow.ExecResult{Err: fmt.Errorf("if condition must evaluate to bool")}
	}
	if cond.Data.(bool) {
		return i.Execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.Execute(stmt.ElseBranch)
	}
	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
}

func (i *Interpreter) VisitWhileStmt(stmt *ast.WhileStmt) controlflow.ExecResult {
	for {
		cond, err := i.Evaluate(stmt.Conditon)
		if err != nil {
			return controlflow.ExecResult{Err: err}
		}
		if cond.Type != value.ValueBool {
			return controlflow.ExecResult{Err: fmt.Errorf("while condition must evaluate to bool")}
		}
		if !cond.Data.(bool) {
			break
		}
		result := i.Execute(stmt.Body)
		if result.Flow != controlflow.FlowNone || result.Err != nil {
			return result
		}
	}
	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
}

func (i *Interpreter) PushEnv(env *environment.Environment) {
	i.envStack = append(i.envStack, i.env)
	i.env = env
}

func (i *Interpreter) PopEnv() {
	n := len(i.envStack)
	if n == 0 {
		panic("env stack underflow")
	}
	i.env = i.envStack[n-1]
	i.envStack = i.envStack[:n-1]
}

func (i *Interpreter) GetEnv() *environment.Environment {
	return i.env
}

func (i *Interpreter) ExecuteBlock(block *ast.BlockStmt, env *environment.Environment) controlflow.ExecResult {
	previous := i.env
	i.env = env

	defer func() {
		i.env = previous
	}()

	for _, stmt := range block.Statements {
		result := i.Execute(stmt)
		if result.Flow != controlflow.FlowNone {
			return result
		}

	}
	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.BlockStmt) controlflow.ExecResult {
	blockEnv := environment.NewEnvironment(i.env)
	i.PushEnv(blockEnv)
	defer i.PopEnv()
	for _, s := range stmt.Statements {
		result := i.Execute(s)
		if result.Flow != controlflow.FlowNone {
			return result
		}
		if result.Flow != controlflow.FlowNone {
			return result
		}
	}

	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
}

// VisitForStmt executes a for loop.
func (i *Interpreter) VisitForStmt(stmt *ast.ForStmt) controlflow.ExecResult {
	forEnv := environment.NewEnvironment(i.env)
	i.PushEnv(forEnv)
	defer i.PopEnv()

	// Init statement
	if stmt.Init != nil {
		result := i.Execute(stmt.Init)
		if result.Flow != controlflow.FlowNone || result.Err != nil {
			return result
		}
	}

	for {
		// Condition expression
		if stmt.CondExpr != nil {
			condVal, err := i.Evaluate(stmt.CondExpr)
			if err != nil {
				return controlflow.ExecResult{Err: err}
			}
			if condVal.Type != value.ValueBool {
				return controlflow.ExecResult{Err: fmt.Errorf("for loop condition must evaluate to bool")}
			}
			if !condVal.Data.(bool) {
				break
			}
		}

		// Execute body
		result := i.Execute(stmt.BodyStmt)
		if result.Flow != controlflow.FlowNone || result.Err != nil {
			return result
		}

		// Update statement
		if stmt.Update != nil {
			result := i.Execute(stmt.Update)
			if result.Flow != controlflow.FlowNone || result.Err != nil {
				return result
			}
		}
	}

	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
}

func (i *Interpreter) VisitCallExpr(expr *ast.CallExpr) controlflow.ExecResult {
	// Evaluate the callee expression (should be a function)
	calleeVal, err := i.Evaluate(expr.Callee)
	if err != nil {
		return controlflow.ExecResult{Err: err}
	}

	callable, ok := calleeVal.Data.(function.Callable)
	if !ok {
		return controlflow.ExecResult{Err: fmt.Errorf("attempt to call non-function value")}
	}
	args := make([]value.Value, len(expr.Arguments))
	for idx, argExpr := range expr.Arguments {
		args[idx], err = i.Evaluate(argExpr)
		if err != nil {
			return controlflow.ExecResult{Err: err}
		}
	}

	fmt.Printf("CALL EXPR DEBUG: %v, %v", args, i.env)

	return callable.Call(args, i)
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
		idxFloat, ok := indexVal.Data.(float64)
		idx := int(idxFloat)
		if !ok || float64(idx) != idxFloat {
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

	inst, ok := objectVal.Data.(*value.StructInstance)
	if !ok || inst == nil {
		return value.Null(), fmt.Errorf("struct instance is corrupt")
	}

	// Access struct field by name (string from token)
	fieldName := expr.Name.Lexeme

	// Assume struct data stored as map[string]value.Value
	// fields, ok := objectVal.Data.(map[string]value.Value)
	// if !ok {
	// 	return value.Null(), fmt.Errorf("struct data corrupted")
	// }

	val, exists := inst.Fields[fieldName]
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
	dict := value.NewNiftelDict()
	for _, pair := range expr.Pairs {
		keyVal, err := i.Evaluate(pair[0])
		if err != nil {
			return value.Null(), err
		}
		valueVal, err := i.Evaluate(pair[1])
		if err != nil {
			return value.Null(), err
		}
		dict.Set(keyVal, valueVal)
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

	fn := function.NewUserFunc(
		"<anonymous>",
		expr.Params,
		expr.Body,
		i.env,
		expr.Func.Line,
		expr.Func.Column)

	return value.Value{
		Type: value.ValueFunc, // Or a dedicated ValueFunc type if you have one
		Data: fn,
		// Meta: nil, // Optional: assign Func type info if available
	}, nil
}

// VisitFuncStmt defines a function in the environment.
func (i *Interpreter) VisitFuncStmt(stmt *ast.FuncStmt) controlflow.ExecResult {
	fmt.Printf("Defining function: %s\n", stmt.Name.Lexeme)
	fn := function.NewUserFunc(
		stmt.Name.Lexeme,
		stmt.Params,
		stmt.Body,
		i.env,
		stmt.Func.Line,
		stmt.Func.Column)
	i.env.Define(stmt.Name.Lexeme, value.Value{
		Type: value.ValueFunc,
		Data: fn,
	})
	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
}

func (i *Interpreter) VisitReturnStmt(stmt *ast.ReturnStmt) controlflow.ExecResult {
	var result value.Value
	var err error

	if len(stmt.Values) == 0 {
		result = value.Null()
	} else if len(stmt.Values) == 1 {
		result, err = i.Evaluate(stmt.Values[0])
		if err != nil {
			return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
		}
	} else {
		vals := make([]value.Value, len(stmt.Values))
		types := make([]*value.TypeInfo, len(stmt.Values))
		for idx, expr := range stmt.Values {
			val, err := i.Evaluate(expr)
			if err != nil {
				return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
			}
			vals[idx] = val
			types[idx] = val.TypeInfo()
		}
		tupleType := value.GetOrRegisterTupleType(types)
		tupleVal := value.NewTupleValue(tupleType, vals)
		result = value.Value{
			Type: value.ValueTuple,
			Data: tupleVal,
		}
	}
	return controlflow.ExecResult{Value: result, Flow: controlflow.FlowReturn}
	// panic(runtimecontrol.ReturnValue{Value: result})

	// for _, expr := range stmt.Values {
	// 	val, err := i.Evaluate(expr)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	results = append(results, val)
	// }
	// if len(results) == 0 {
	// 	panic(runtimecontrol.ReturnValue{Value: value.Null()})
	// } else if len(results) == 1 {
	// 	panic(runtimecontrol.ReturnValue{Value: results[0]})
	// } else {
	// 	panic(runtimecontrol.ReturnValue{Value: value.Value{
	// 		Type: value.ValueList,
	// 		Data: results,
	// 	}})
	// }
	// var retVal value.Value
	// var err error
	// if stmt.Value != nil {
	// 	retVal, err = i.Evaluate(stmt.Value)
	// 	if err != nil {
	// 		return err
	// 	}
	// } else {
	// 	retVal = value.Null()
	// }
	// panic with a special ReturnValue to unwind execution
	// panic(runtimecontrol.ReturnValue{Value: retVal})
}

// VisitBreakStmt handles break statement in loops.
func (i *Interpreter) VisitBreakStmt(stmt *ast.BreakStmt) controlflow.ExecResult {
	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowBreak}
}

// VisitContinueStmt handles continue statement in loops.
func (i *Interpreter) VisitContinueStmt(stmt *ast.ContinueStmt) controlflow.ExecResult {
	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowContinue}
}
