package interpreter

import (
	"github.com/ithinkiborkedit/niftelv2.git/internal/controlflow"
	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
)

type ExprVisitor interface {
	VisitBinaryExpr(expr *ast.BinaryExpr) controlflow.ExecResult
	VisitUnaryExpr(expr *ast.UnaryExpr) controlflow.ExecResult
	VisitLiteralExpr(expr *ast.LiteralExpr) controlflow.ExecResult
	VisitVariableExpr(expr *ast.VariableExpr) controlflow.ExecResult
	VisitCallExpr(expr *ast.CallExpr) controlflow.ExecResult
	VisitIndexExpr(expr *ast.IndexExpr) controlflow.ExecResult
	VisitGetExpr(expr *ast.GetExpr) controlflow.ExecResult
	VisitListExpr(expr *ast.ListExpr) controlflow.ExecResult
	VisitDictExpr(expr *ast.DictExpr) controlflow.ExecResult
	VisitFuncExpr(expr *ast.FuncExpr) controlflow.ExecResult
	VisitStructLiteralExpr(expr *ast.StructLiteralExpr) controlflow.ExecResult
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
