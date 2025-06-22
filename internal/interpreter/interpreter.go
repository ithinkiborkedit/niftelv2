package interpreter

import (
	"fmt"
	"errors"

	"github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	"github.com/ithinkiborkedit/niftelv2.git/internal/environment"
	"github.com/ithinkiborkedit/niftelv2.git/internal/token"
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
func (i *Interpreter) Evaluate(expr nifast.Expr) (value.Value, error) {
	switch e := expr.(type) {
	case *nifast.LiteralExpr:
		return i.VisitLiteralExpr(e)
	case *nifast.VariableExpr:
		return i.VisitVariableExpr(e)
	case *nifast.BinaryExpr:
		return i.VisitBinaryExpr(e)
	case *nifast.UnaryExpr:
		return i.VisitUnaryExpr(e)
	case *nifast.CallExpr:
		return i.VisitCallExpr(e)
	case *nifast.GetExpr:
		return i.VisitGetExpr(e)
	case *nifast.ListExpr:
		return i.VisitListExpr(e)
	case *nifast.DictExpr:
		return i.VisitDictExpr(e)
	case *nifast.FuncExpr:
		return i.VisitFuncExpr(e)
	default:
		return value.Null(), fmt.Errorf("unknown expression type %T", expr)
	}
}

// Execute dispatches to the correct Stmt handler.
func (i *Interpreter) Execute(stmt nifast.Stmt) error {
	switch s := stmt.(type) {
	case *nifast.VarStmt:
		return i.VisitVarStmt(s)
	case *nifast.ShortVarStmt:
		return i.VisitShortVarStmt(s)
	case *nifast.AssignStmt:
		return i.VisitAssignStmt(s)
	case *nifast.PrintStmt:
		return i.VisitPrintStmt(s)
	case *nifast.ExprStmt:
		return i.VisitExprStmt(s)
	case *nifast.IfStmt:
		return i.VisitIfStmt(s)
	case *nifast.WhileStmt:
		return i.VisitWhileStmt(s)
	case *nifast.ForStmt:
		return i.VisitForStmt(s)
	case *nifast.FuncStmt:
		return i.VisitFuncStmt(s)
	case *nifast.ReturnStmt:
		return i.VisitReturnStmt(s)
	case *nifast.BreakStmt:
		return i.VisitBreakStmt(s)
	case *nifast.ContinueStmt:
		return i.VisitContinueStmt(s)
	case *nifast.BlockStmt:
		return i.VisitBlockStmt(s)
	default:
		return fmt.Errorf("unknown statement type %T", stmt)
	}
}

// --- Expression Visitors ---

func (i *Interpreter) VisitLiteralExpr(expr *nifast.LiteralExpr) (value.Value, error) {
	tok := expr.Value
	switch tok.Type {
	case value.ValueInt:
		if val, ok := tok.Data.(int64); ok {
			return value.Value{Type: value.ValueInt, Data: val}, nil
		}
		return value.Null(), errors.New("invalid int literal token data")
	case value.ValueString:
		return value.Value{Type: value.ValueString, Data: tok.Lexeme}, nil
	case value.ValueBool:
		if tok.Lexeme == "true" {
			return value.Value{Type: value.ValueBool, Data: true}, nil
		} else if tok.Lexeme == "false" {
			return value.Value{Type: value.ValueBool, Data: false}, nil
		}
		return value.Null(), errors.New("invalid bool literal token")
	case value.ValueNull:
		return value.Null(), nil
	default:
		return value.Null(), fmt.Errorf("unsupported literal token type %v", tok.Type)
	}
}

func (i *Interpreter) VisitVariableExpr(expr *nifast.VariableExpr) (value.Value, error) {
	val, err := i.env.Get(expr.Name.Lexeme)
	if err != nil {
		return value.Null(), fmt.Errorf("undefined variable %s", expr.Name.Lexeme)
	}
	return val, nil
}

func (i *Interpreter) VisitBinaryExpr(expr *nifast.BinaryExpr) (value.Value, error) {
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

func (i *Interpreter) VisitUnaryExpr(expr *nifast.UnaryExpr) (value.Value, error) {
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

func (i *Interpreter) VisitVarStmt(stmt *nifast.VarStmt) error {
	val, err := i.Evaluate(stmt.Init)
	if err != nil {
		return err
	}
	i.env.Define(stmt.Name.Lexeme, val)
	return nil
}

func (i *Interpreter) VisitShortVarStmt(stmt *nifast.ShortVarStmt) error {
	val, err := i.Evaluate(stmt.Init)
	if err != nil {
		return err
	}
	i.env.Define(stmt.Name.Lexeme, val)
	return nil
}

func (i *Interpreter) VisitAssignStmt(stmt *nifast.AssignStmt) error {
	val, err := i.Evaluate(stmt.Value)
	if err != nil {
		return err
	}
	return i.env.Assign(stmt.Name.Lexeme, val)
}

func (i *Interpreter) VisitPrintStmt(stmt *nifast.PrintStmt) error {
	val, err := i.Evaluate(stmt.Expr)
	if err != nil {
		return err
	}
	fmt.Println(val.String())
	return nil
}

func (i *Interpreter) VisitExprStmt(stmt *nifast.ExprStmt) error {
	_, err := i.Evaluate(stmt.Expr)
	return err
}

func (i *Interpreter) VisitIfStmt(stmt *nifast.IfStmt) error {
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

func (i *Interpreter) VisitWhileStmt(stmt *nifast.WhileStmt) error {
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

func (i *Interpreter) VisitBlockStmt(stmt *nifast.BlockStmt) error {
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

// TODO: Implement VisitForStmt, VisitFuncStmt, VisitReturnStmt, VisitBreakStmt, VisitContinueStmt
	}
}


