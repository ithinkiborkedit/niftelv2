package interpreter

import (
	"errors"
	"fmt"

	"github.com/ithinkiborkedit/niftelv2.git/internal/controlflow"
	"github.com/ithinkiborkedit/niftelv2.git/internal/environment"
	"github.com/ithinkiborkedit/niftelv2.git/internal/function"
	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	token "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"
	"github.com/ithinkiborkedit/niftelv2.git/internal/symtable"
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
	interp := &Interpreter{
		env: environment.NewEnvironment(nil),
	}
	if err := interp.RegisterBuiltInTypes(); err != nil {
		panic(fmt.Sprintf("Interpreter failed to register builtin types: %v", err))
	}

	return interp
}

func (i *Interpreter) Eval(expr ast.Expr) controlflow.ExecResult {
	return i.Evaluate(expr)
}

func (i *Interpreter) RegisterBuiltInTypes() error {
	for name, typ := range value.BuiltInTypes {
		if err := i.env.DefineType(typ); err != nil {
			fmt.Printf("SKIP DEBUG TYPE %q already defined", name)
			continue
		}
	}

	fmt.Println("DEBUG ENV TYPES after registration:[")
	for name := range i.env.SymbolTable().Types {
		fmt.Printf("%q ", name)
	}
	fmt.Println("]")
	return nil
	// for _, name := range value.BuiltInTypes {

	// 	// typ, ok := value.LookupType(name)
	// 	// if !ok {
	// 	// 	return fmt.Errorf("missing builtint type %q", name)
	// 	// }
	// 	// if err := i.env.DefineType(typ); err != nil {
	// 	// 	return fmt.Errorf("failed to register type %q: %w", name, err)
	// 	// }
	// 	// fmt.Println(name)
	// }
	// fmt.Print("DEBUG: env.Types after RegisterBuiltInTypes:[\n")
	// for name := range i.env.SymbolTable().Types {
	// 	fmt.Printf("%q ", name)
	// }
	// return nil
}

// Evaluate dispatches to the correct Expr handler.
func (i *Interpreter) Evaluate(expr ast.Expr) controlflow.ExecResult {
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
		return controlflow.ExecResult{Err: fmt.Errorf("unknown expression type %T", expr)}
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

func (i *Interpreter) VisitStructLiteralExpr(expr *ast.StructLiteralExpr) controlflow.ExecResult {
	structName := expr.TypeName.Lexeme

	typeSym, ok := value.GetType(structName)
	if !ok || typeSym == nil {
		return controlflow.ExecResult{Value: value.Null(), Err: fmt.Errorf("struct type '%s' not found", structName)}
	}

	var orderedFields []token.Token

	for fname := range typeSym.Fields {
		orderedFields = append(orderedFields, token.Token{Lexeme: fname})
	}

	structType := &value.StructType{
		Name:   typeSym.SymName,
		Fields: orderedFields,
	}
	// for fname := range typeInfo.Fields {
	// 	structType.Fields = append(structType.Fields, token.Token{Lexeme: fname})
	// }

	instance := &value.StructInstance{
		Type:   structType,
		Fields: make(map[string]value.Value),
	}

	for fname, exprval := range expr.Fields {
		valRes := i.Evaluate(exprval)
		if valRes.Err != nil {
			return controlflow.ExecResult{Value: value.Null(), Err: fmt.Errorf("error in field '%s': %w", fname, valRes.Err)}
		}
		instance.Fields[fname] = valRes.Value
	}

	for fname := range typeSym.Fields {
		if _, ok := instance.Fields[fname]; !ok {
			instance.Fields[fname] = value.Null()
		}
	}

	return controlflow.ExecResult{
		Value: value.Value{
			Type: value.ValueStruct,
			Data: instance,
			// Meta: typeSym.,
		},
		Flow: controlflow.FlowNone,
	}
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.LiteralExpr) controlflow.ExecResult {
	tok := expr.Value
	switch tok.Type {
	case token.TokenNumber:
		switch val := tok.Data.(type) {
		case int:
			return controlflow.ExecResult{Value: value.Value{Type: value.ValueInt, Data: float64(val)}, Flow: controlflow.FlowNone}
		case int64:
			return controlflow.ExecResult{Value: value.Value{Type: value.ValueInt, Data: float64(val)}, Flow: controlflow.FlowNone}
		case float64:
			return controlflow.ExecResult{Value: value.Value{Type: value.ValueInt, Data: float64(val)}, Flow: controlflow.FlowNone}
		default:
			return controlflow.ExecResult{Err: errors.New("invalid int literal token data")}

		}

	case token.TokenString:
		return controlflow.ExecResult{Value: value.Value{Type: value.ValueString, Data: tok.Lexeme}, Flow: controlflow.FlowNone}
	case token.TokenBool:
		switch tok.Lexeme {
		case "true":
			return controlflow.ExecResult{Value: value.Value{Type: value.ValueBool, Data: true}, Flow: controlflow.FlowNone}
		case "false":
			return controlflow.ExecResult{Value: value.Value{Type: value.ValueBool, Data: false}, Flow: controlflow.FlowNone}
		}
		return controlflow.ExecResult{Err: errors.New("invalid bool literal token")}
	case token.TokenNull:
		return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
	default:
		return controlflow.ExecResult{Err: fmt.Errorf("unsupported literal token type %v", tok.Type)}
	}
}

func (i *Interpreter) VisitVariableExpr(expr *ast.VariableExpr) controlflow.ExecResult {
	val, err := i.env.GetVar(expr.Name.Lexeme)
	if err != nil {
		return controlflow.ExecResult{Err: fmt.Errorf("undefined variable %s", expr.Name.Lexeme)}
	}
	return controlflow.ExecResult{Value: val, Flow: controlflow.FlowNone}
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.BinaryExpr) controlflow.ExecResult {
	leftRes := i.Evaluate(expr.Left)
	if leftRes.Err != nil {
		return controlflow.ExecResult{Err: leftRes.Err}
	}
	left := leftRes.Value
	rightRes := i.Evaluate(expr.Right)
	if rightRes.Err != nil {
		return controlflow.ExecResult{Err: rightRes.Err}
	}
	right := rightRes.Value
	switch expr.Operator.Type {
	case token.TokenPlus:
		if left.Type == value.ValueInt && right.Type == value.ValueInt {
			return controlflow.ExecResult{Value: value.Value{
				Type: value.ValueInt, Data: left.Data.(float64) + right.Data.(float64)},
				Flow: controlflow.FlowNone}
		}
		if left.Type == value.ValueString && right.Type == value.ValueString {
			return controlflow.ExecResult{Value: value.Value{
				Type: value.ValueString, Data: left.Data.(string) + right.Data.(string)},
				Flow: controlflow.FlowNone}
		}
		return controlflow.ExecResult{Err: fmt.Errorf("unsupported operand types for +: %v and %v", left.Type, right.Type)}
	case token.TokenMinus:
		if left.Type == value.ValueInt && right.Type == value.ValueInt {
			return controlflow.ExecResult{Value: value.Value{
				Type: value.ValueInt, Data: left.Data.(float64) - right.Data.(float64)},
				Flow: controlflow.FlowNone}
		}
		return controlflow.ExecResult{Err: fmt.Errorf("unsupported operand types for '-': %v and %v", left.Type, right.Type)}
	case token.TokenStar:
		if left.Type == value.ValueInt && right.Type == value.ValueInt {
			return controlflow.ExecResult{Value: value.Value{
				Type: value.ValueInt, Data: left.Data.(float64) * right.Data.(float64)},
				Flow: controlflow.FlowNone}
		}
		return controlflow.ExecResult{Err: fmt.Errorf("unsupported operand types for *: %v and %v", left.Type, right.Type)}
	case token.TokenFWDSlash:
		if left.Type == value.ValueInt && right.Type == value.ValueInt {
			if right.Data.(float64) == 0 {
				return controlflow.ExecResult{Err: fmt.Errorf("division by zero")}
			}
			return controlflow.ExecResult{Value: value.Value{
				Type: value.ValueInt, Data: left.Data.(float64) / right.Data.(float64)},
				Flow: controlflow.FlowNone}
		}
		return controlflow.ExecResult{Err: fmt.Errorf("unsupported operand types for /: %v and %v", left.Type, right.Type)}
	case token.TokenGreater:
		if left.Type == value.ValueInt && right.Type == value.ValueInt {
			return controlflow.ExecResult{Value: value.Value{
				Type: value.ValueBool, Data: left.Data.(float64) > right.Data.(float64)},
				Flow: controlflow.FlowNone}
		}
		return controlflow.ExecResult{Err: fmt.Errorf("unsupported operand types for >: %v and %v", left.Type, right.Type)}
	case token.TokenLess:
		if left.Type == value.ValueInt && right.Type == value.ValueInt {
			return controlflow.ExecResult{Value: value.Value{
				Type: value.ValueBool, Data: left.Data.(float64) < right.Data.(float64)},
				Flow: controlflow.FlowNone}
		}
		return controlflow.ExecResult{Err: fmt.Errorf("unsupported operand types for <: %v and %v", left.Type, right.Type)}
	case token.TokenEqality:
		if left.Type != right.Type {
			return controlflow.ExecResult{Value: value.Value{
				Type: value.ValueBool, Data: false},
				Flow: controlflow.FlowNone}
		}
		return controlflow.ExecResult{Value: value.Value{Type: value.ValueBool, Data: left.Data == right.Data}, Flow: controlflow.FlowNone}
	case token.TokenBangEqal:
		if left.Type != right.Type {
			return controlflow.ExecResult{Value: value.Value{
				Type: value.ValueBool, Data: true},
				Flow: controlflow.FlowNone}
		}
		return controlflow.ExecResult{Value: value.Value{Type: value.ValueBool, Data: left.Data != right.Data}, Flow: controlflow.FlowNone}
	default:
		return controlflow.ExecResult{Err: fmt.Errorf("unsupported binary operator %v", expr.Operator.Lexeme)}
	}
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.UnaryExpr) controlflow.ExecResult {
	rightRes := i.Evaluate(expr.Right)
	if rightRes.Err != nil {
		return controlflow.ExecResult{Err: rightRes.Err}
	}
	right := rightRes.Value
	switch expr.Operator.Type {
	case token.TokenBang:
		if right.Type != value.ValueBool {
			return controlflow.ExecResult{Err: fmt.Errorf("operator ! requires boolean operand")}
		}
		return controlflow.ExecResult{
			Value: value.Value{Type: value.ValueBool, Data: !right.Data.(bool)},
			Flow:  controlflow.FlowNone,
		}
	case token.TokenMinus:
		if right.Type != value.ValueInt {
			return controlflow.ExecResult{Err: fmt.Errorf("operator - requires int operand")}
		}
		return controlflow.ExecResult{
			Value: value.Value{Type: value.ValueInt, Data: -right.Data.(float64)},
			Flow:  controlflow.FlowNone,
		}
	default:
		return controlflow.ExecResult{Err: fmt.Errorf("unsupported unary operator %v", expr.Operator.Lexeme)}
	}
}

// Implement remaining expression visitors: VisitCallExpr, VisitGetExpr, VisitListExpr, VisitDictExpr, VisitFuncExpr

// --- Statement Visitors ---

func (i *Interpreter) VisitStructStmt(stmt *ast.StructStmt) controlflow.ExecResult {
	structName := stmt.Name.Lexeme

	if i.env.HasLocalType(structName) {
		return controlflow.ExecResult{Err: fmt.Errorf("struct '%s' already defined", structName)}
	}

	fields := make(map[string]*symtable.TypeSymbol)
	for _, field := range stmt.Fields {
		if len(field.Names) != 1 {
			return controlflow.ExecResult{Err: fmt.Errorf("struct field must have exactly one name!")}
		}
		fieldName := field.Names[0].Lexeme
		fieldTypeName := field.Type.Lexeme
		fieldType, ok := i.env.LookupType(fieldTypeName)
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

	structSym := &symtable.TypeSymbol{
		SymName: structName,
		SymKind: symtable.SymbolTypes,
		Fields:  fields,
	}

	if err := i.env.DefineType(structSym); err != nil {
		return controlflow.ExecResult{Err: err}
	}
	fmt.Printf("[INFO] REGISTERED struct type: '%s'\n", stmt.Name.Lexeme)
	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
}

func (i *Interpreter) VisitVarStmt(stmt *ast.VarStmt) controlflow.ExecResult {
	if stmt.Type.Lexeme == "" {
		return controlflow.ExecResult{Err: fmt.Errorf("var declarations require a type")}
	}
	valRes := i.Evaluate(stmt.Init)
	if valRes.Err != nil {
		return controlflow.ExecResult{Err: valRes.Err}
	}
	val := valRes.Value
	if len(stmt.Names) == 1 {
		name := stmt.Names[0].Lexeme
		var varTypeSym *symtable.TypeSymbol = nil
		valRes := i.Evaluate(stmt.Init)
		if valRes.Err != nil {
			return controlflow.ExecResult{Err: valRes.Err}
		}
		val := valRes.Value
		// if stmt.Type.Lexeme != "" {
		// 	typeSym, ok := i.env.LookupType(stmt.Type.Lexeme)
		// 	if !ok {
		// 		return controlflow.ExecResult{Err: fmt.Errorf("unknown type '%s'", stmt.Type.Lexeme)}
		// 	}
		// 	varTypeSym = typeSym
		// }
		varSym := &symtable.VarSymbol{
			SymName: name,
			SymKind: symtable.SymbolVar,
			Type:    varTypeSym,
			Mutable: true,
		}
		if err := i.env.DefineVar(varSym); err != nil {
			return controlflow.ExecResult{Err: err}
		}
		if err := i.env.AssignVar(name, val); err != nil {
			return controlflow.ExecResult{Err: err}
		}

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
		name := name.Lexeme
		if stmt.Type.Lexeme == "" {
			return controlflow.ExecResult{Err: fmt.Errorf("var declarations require a type")}
		}
		varSym := &symtable.VarSymbol{
			SymName: name,
			SymKind: symtable.SymbolVar,
			Mutable: true,
		}
		if err := i.env.DefineVar(varSym); err != nil {
			return controlflow.ExecResult{Err: err}

		}
		if err := i.env.AssignVar(name, values[idx]); err != nil {
			return controlflow.ExecResult{Err: err}
		}
	}
	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
}

func (i *Interpreter) VisitShortVarStmt(stmt *ast.ShortVarStmt) controlflow.ExecResult {
	valRes := i.Evaluate(stmt.Init)
	if valRes.Err != nil {
		return controlflow.ExecResult{Err: valRes.Err}
	}
	val := valRes.Value
	name := stmt.Name.Lexeme

	varSym := &symtable.VarSymbol{
		SymName: name,
		SymKind: symtable.SymbolVar,
		Type:    nil,
		Mutable: true,
	}
	if err := i.env.DefineVar(varSym); err != nil {
		return controlflow.ExecResult{Err: err}
	}
	if err := i.env.AssignVar(name, val); err != nil {
		return controlflow.ExecResult{Err: err}

	}
	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
}

func (i *Interpreter) VisitAssignStmt(stmt *ast.AssignStmt) controlflow.ExecResult {
	valRes := i.Evaluate(stmt.Value)
	if valRes.Err != nil {
		return controlflow.ExecResult{Err: valRes.Err}
	}
	val := valRes.Value
	if err := i.env.AssignVar(stmt.Name.Lexeme, val); err != nil {
		return controlflow.ExecResult{Err: err}
	}
	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.PrintStmt) controlflow.ExecResult {
	valRes := i.Evaluate(stmt.Expr)
	if valRes.Err != nil {
		return controlflow.ExecResult{Err: valRes.Err}
	}
	val := valRes.Value
	fmt.Println(val.String())
	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
}

func (i *Interpreter) VisitExprStmt(stmt *ast.ExprStmt) controlflow.ExecResult {
	resultRes := i.Evaluate(stmt.Expr)
	if resultRes.Err != nil {
		return controlflow.ExecResult{Err: resultRes.Err}
	}
	result := resultRes.Value
	if i.ShouldPrintResults && !result.IsNull() {
		fmt.Println(result.String())
	}

	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
}

func (i *Interpreter) VisitIfStmt(stmt *ast.IfStmt) controlflow.ExecResult {
	condRes := i.Evaluate(stmt.Conditon)
	if condRes.Err != nil {
		return controlflow.ExecResult{Err: condRes.Err}
	}
	cond := condRes.Value
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
		condRes := i.Evaluate(stmt.Conditon)
		if condRes.Err != nil {
			return controlflow.ExecResult{Err: condRes.Err}
		}
		cond := condRes.Value
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
			condRes := i.Evaluate(stmt.CondExpr)
			if condRes.Err != nil {
				return controlflow.ExecResult{Err: condRes.Err}
			}
			condVal := condRes.Value
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
	calleeRes := i.Evaluate(expr.Callee)
	if calleeRes.Err != nil {
		return controlflow.ExecResult{Err: calleeRes.Err}
	}
	calleeVal := calleeRes.Value
	callable, ok := calleeVal.Data.(function.Callable)
	if !ok {
		return controlflow.ExecResult{Err: fmt.Errorf("attempt to call non-function value")}
	}
	args := make([]value.Value, len(expr.Arguments))
	for idx, argExpr := range expr.Arguments {
		argRes := i.Evaluate(argExpr)
		if argRes.Err != nil {
			return controlflow.ExecResult{Err: argRes.Err}
		}
		args[idx] = argRes.Value
	}

	return callable.Call(args, i)
}

func (i *Interpreter) VisitIndexExpr(expr *ast.IndexExpr) controlflow.ExecResult {
	// Evaluate the collection expression
	collectionRes := i.Evaluate(expr.Collection)
	if collectionRes.Err != nil {
		return controlflow.ExecResult{Err: collectionRes.Err}
	}
	collectionVal := collectionRes.Value

	// Evaluate the index/key expression
	indexRes := i.Evaluate(expr.Index)
	if indexRes.Err != nil {
		return controlflow.ExecResult{Err: indexRes.Err}
	}

	indexVal := indexRes.Value

	switch collectionVal.Type {
	case value.ValueList:
		list, ok := collectionVal.Data.([]value.Value)
		if !ok {
			return controlflow.ExecResult{Err: fmt.Errorf("list data is corrupted")}
		}
		idxFloat, ok := indexVal.Data.(float64)
		idx := int(idxFloat)
		if !ok || float64(idx) != idxFloat {
			return controlflow.ExecResult{Err: fmt.Errorf("list index must be integer")}
		}
		if idx < 0 || idx >= len(list) {
			return controlflow.ExecResult{Err: fmt.Errorf("list index out of range")}
		}
		return controlflow.ExecResult{Value: list[idx], Flow: controlflow.FlowNone}

	case value.ValueDict:
		dict, ok := collectionVal.Data.(map[string]value.Value)
		if !ok {
			return controlflow.ExecResult{Err: fmt.Errorf("dict data is corrupted")}
		}
		keyStr, ok := indexVal.Data.(string)
		if !ok {
			return controlflow.ExecResult{Err: fmt.Errorf("dict key must be string")}
		}
		val, exists := dict[keyStr]
		if !exists {
			return controlflow.ExecResult{Err: fmt.Errorf("dict key not found: %s", keyStr)}
		}
		return controlflow.ExecResult{Value: val, Flow: controlflow.FlowNone}

	default:
		return controlflow.ExecResult{Err: fmt.Errorf("indexing unsupported on type %v", collectionVal.Type)}
	}
}

func (i *Interpreter) VisitGetExpr(expr *ast.GetExpr) controlflow.ExecResult {
	// Evaluate object expression
	objectRes := i.Evaluate(expr.Object)
	if objectRes.Err != nil {
		return controlflow.ExecResult{Err: objectRes.Err}
	}

	objectVal := objectRes.Value

	if objectVal.Type != value.ValueStruct {
		return controlflow.ExecResult{Err: fmt.Errorf("attempt to get property on non-struct type")}
	}

	inst, ok := objectVal.Data.(*value.StructInstance)
	if !ok || inst == nil {
		return controlflow.ExecResult{Err: fmt.Errorf("struct instance is corrupt")}
	}

	fieldName := expr.Name.Lexeme

	val, exists := inst.Fields[fieldName]
	if !exists {
		return controlflow.ExecResult{Err: fmt.Errorf("struct field '%s' not found", fieldName)}
	}
	return controlflow.ExecResult{Value: val, Flow: controlflow.FlowNone}
}

func (i *Interpreter) VisitListExpr(expr *ast.ListExpr) controlflow.ExecResult {
	elements := make([]value.Value, 0, len(expr.Elements))
	for _, elementExpr := range expr.Elements {
		elemRes := i.Evaluate(elementExpr)
		if elemRes.Err != nil {
			return controlflow.ExecResult{Err: elemRes.Err}
		}
		elements = append(elements, elemRes.Value)
	}
	return controlflow.ExecResult{
		Value: value.Value{
			Type: value.ValueList,
			Data: elements,
		},
		Flow: controlflow.FlowNone,
	}
}

func (i *Interpreter) VisitDictExpr(expr *ast.DictExpr) controlflow.ExecResult {
	dict := value.NewNiftelDict()
	for _, pair := range expr.Pairs {
		keyRes := i.Evaluate(pair[0])
		if keyRes.Err != nil {
			return controlflow.ExecResult{Err: keyRes.Err}
		}
		valueRes := i.Evaluate(pair[1])
		// valueVal, err := i.Evaluate(pair[1])
		if valueRes.Err != nil {
			return controlflow.ExecResult{Err: valueRes.Err}
		}
		dict.Set(keyRes.Value, valueRes.Value)
	}
	return controlflow.ExecResult{
		Value: value.Value{
			Type: value.ValueDict,
			Data: dict,
		},
		Flow: controlflow.FlowNone,
	}
}

func (i *Interpreter) VisitFuncExpr(expr *ast.FuncExpr) controlflow.ExecResult {
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

	return controlflow.ExecResult{
		Value: value.Value{
			Type: value.ValueFunc, // Or a dedicated ValueFunc type if you have one
			Data: fn,
			// Meta: nil, // Optional: assign Func type info if available
		},
		Flow: controlflow.FlowNone,
	}
}

// VisitFuncStmt defines a function in the environment.
func (i *Interpreter) VisitFuncStmt(stmt *ast.FuncStmt) controlflow.ExecResult {
	name := stmt.Name.Lexeme

	if i.env.HasLocalFunc(name) {
		return controlflow.ExecResult{
			Err: fmt.Errorf("function '%s' already defined in this scope", stmt.Name.Lexeme),
		}
	}
	var returnTypes []*symtable.TypeSymbol
	for _, retTok := range stmt.Return {
		if retTok.Lexeme == "" {
			continue
		}
		typeSym, ok := i.env.LookupType(retTok.Lexeme)
		if !ok {
			return controlflow.ExecResult{Err: fmt.Errorf("unkown type '%s' for function '%s'", retTok.Lexeme, name)}
		}
		returnTypes = append(returnTypes, typeSym)
	}
	var params []symtable.VarSymbol

	for _, param := range stmt.Params {
		var typeSym *symtable.TypeSymbol
		if param.Type.Lexeme != "" {
			ts, ok := i.env.LookupType(param.Type.Lexeme)
			if !ok {
				fmt.Println("DEBUG: KNOWN TYPES IN INTERPRETER ENV:", func() []string {
					names := []string{}
					for name := range i.env.SymbolTable().Types {
						names = append(names, name)
					}
					return names
				}())
				return controlflow.ExecResult{Err: fmt.Errorf("unknown parameter type '%s'", param.Type.Lexeme)}
			}
			typeSym = ts
		}
		params = append(params, symtable.VarSymbol{
			SymName: param.Name.Lexeme,
			SymKind: symtable.SymbolVar,
			Type:    typeSym,
			Mutable: false,
		})

	}

	fmt.Printf("Defining function: %s\n", stmt.Name.Lexeme)
	funcSym := &symtable.FuncSymbol{
		SymName:    name,
		Params:     params,
		ReturnType: returnTypes,
		TypeParams: nil,
	}
	if err := i.env.DefineFunc(funcSym); err != nil {
		return controlflow.ExecResult{Err: err}
	}
	varSym := &symtable.VarSymbol{
		SymName: name,
		SymKind: symtable.SymbolVar,
		Mutable: false,
	}
	if err := i.env.DefineVar(varSym); err != nil {
		return controlflow.ExecResult{Err: err}
	}
	fn := function.NewUserFunc(
		stmt.Name.Lexeme,
		stmt.Params,
		stmt.Body,
		i.env,
		stmt.Func.Line,
		stmt.Func.Column)
	if err := i.env.AssignVar(name, value.Value{
		Type: value.ValueFunc,
		Data: fn,
	}); err != nil {
		return controlflow.ExecResult{Err: err}
	}

	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowNone}
}

func (i *Interpreter) VisitReturnStmt(stmt *ast.ReturnStmt) controlflow.ExecResult {
	var result value.Value

	if len(stmt.Values) == 0 {
		result = value.Null()
	} else if len(stmt.Values) == 1 {
		valRes := i.Evaluate(stmt.Values[0])
		if valRes.Err != nil {
			return controlflow.ExecResult{Err: valRes.Err}
		}
		result = valRes.Value
	} else {
		vals := make([]value.Value, len(stmt.Values))
		types := make([]*symtable.TypeSymbol, len(stmt.Values))
		for idx, expr := range stmt.Values {
			valRes := i.Evaluate(expr)
			if valRes.Err != nil {
				return controlflow.ExecResult{Err: valRes.Err}
			}
			vals[idx] = valRes.Value
			types[idx] = valRes.Value.TypeInfo()
		}
		tupleType := value.GetOrRegisterTupleType(types)
		tupleVal := value.NewTupleValue(tupleType, vals)
		result = value.Value{
			Type: value.ValueTuple,
			Data: tupleVal,
		}
	}
	return controlflow.ExecResult{Value: result, Flow: controlflow.FlowReturn}
}

// VisitBreakStmt handles break statement in loops.
func (i *Interpreter) VisitBreakStmt(stmt *ast.BreakStmt) controlflow.ExecResult {
	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowBreak}
}

// VisitContinueStmt handles continue statement in loops.
func (i *Interpreter) VisitContinueStmt(stmt *ast.ContinueStmt) controlflow.ExecResult {
	return controlflow.ExecResult{Value: value.Null(), Flow: controlflow.FlowContinue}
}
