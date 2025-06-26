package interpreter

import (
	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	"github.com/ithinkiborkedit/niftelv2.git/internal/value"
)

type ExprVisitor interface {
	VisitBinaryExpr(expr *ast.BinaryExpr) (value.Value, error)
	VisitUnaryExpr(expr *ast.UnaryExpr) (value.Value, error)
	VisitLiteralExpr(expr *ast.LiteralExpr) (value.Value, error)
	VisitVariableExpr(expr *ast.VariableExpr) (value.Value, error)
	VisitCallExpr(expr *ast.CallExpr) (value.Value, error)
	VisitIndexExpr(expr *ast.IndexExpr) (value.Value, error)
	VisitGetExpr(expr *ast.GetExpr) (value.Value, error)
	VisitListExpr(expr *ast.ListExpr) (value.Value, error)
	VisitDictExpr(expr *ast.DictExpr) (value.Value, error)
	VisitFuncExpr(expr *ast.FuncExpr) (value.Value, error)
	VisitStructLiteralExpr(expr *ast.StructLiteralExpr) (value.Value, error)
}

type StmtVisitor interface {
	VisitVarStmt(expr *ast.VarStmt) error
	VisitShortVarStmt(expr *ast.ShortVarStmt) error
	VisitAssignStmt(expr *ast.AssignStmt) error
	VisitPrintStmt(expr *ast.PrintStmt) error
	VisitExprStmt(expr *ast.ExprStmt) error
	VisitIfStmt(expr *ast.IfStmt) error
	VisitWhileStmt(expr *ast.WhileStmt) error
	VisitForStmt(expr *ast.ForStmt) error
	VisitBlockStmt(expr *ast.BlockStmt) error
	VisitFuncStmt(expr *ast.FuncStmt) error
	VisitStructStmt(expr *ast.StructStmt) error
	VisitReturnStmt(expr *ast.ReturnStmt) error
	VisitBreakStmt(expr *ast.BreakStmt) error
	VisitContinueStmt(expr *ast.ContinueStmt) error
}
