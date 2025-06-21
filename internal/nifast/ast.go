package nifast

import token "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"

type Expr interface {
	exprNode()
	Pos() (line, column int)
}

type Stmt interface {
	stmtNode()
	Pos() (line, column int)
}

type BinaryExpr struct {
	Left     Expr
	Operator token.NifToken
	Right    Expr
}

func (*BinaryExpr) exprNode()         {}
func (e *BinaryExpr) Pos() (int, int) { return e.Operator.Line, e.Operator.Column }

type UnaryExpr struct {
	Operator token.NifToken
	Right    Expr
}

func (*UnaryExpr) exprNode()         {}
func (e *UnaryExpr) Pos() (int, int) { return e.Operator.Line, e.Operator.Column }

type VariableExpr struct {
	Name token.NifToken
}

func (*VariableExpr) exprNode()         {}
func (e *VariableExpr) Pos() (int, int) { return e.Name.Line, e.Name.Column }

type LiteralExpr struct {
	Value token.NifToken
}

func (*LiteralExpr) exprNode()         {}
func (e *LiteralExpr) Pos() (int, int) { return e.Value.Line, e.Value.Column }

type CallExpr struct {
	Callee    Expr
	Paren     token.NifToken
	Arguments []Expr
}

func (*CallExpr) exprNode()         {}
func (e *CallExpr) Pos() (int, int) { return e.Paren.Line, e.Paren.Column }

type IndexExpr struct {
	Collection Expr
	Bracket    token.NifToken
	Index      Expr
}

func (*IndexExpr) exprNode()         {}
func (e *IndexExpr) Pos() (int, int) { return e.Bracket.Line, e.Bracket.Column }

type GetExpr struct {
	Object Expr
	Name   token.NifToken
}

func (*GetExpr) exprNode()         {}
func (e *GetExpr) Pos() (int, int) { return e.Name.Line, e.Name.Column }

type ListExpr struct {
	Elements []Expr
	LBracket token.NifToken
}

func (*ListExpr) exprNode()         {}
func (e *ListExpr) Pos() (int, int) { return e.LBracket.Line, e.LBracket.Column }

type DictExpr struct {
	Pairs  [][2]Expr
	LBrace token.NifToken
}

func (*DictExpr) exprNode()         {}
func (e *DictExpr) Pos() (int, int) { return e.LBrace.Line, e.LBrace.Column }

type FuncExpr struct {
	Params []Param
	Body   *BlockStmt
	Func   token.NifToken
}

type Param struct {
	Name token.NifToken
	Type token.NifToken
}

func (*FuncExpr) exprNode()         {}
func (e *FuncExpr) Pos() (int, int) { return e.Func.Line, e.Func.Column }

//STATEMENTS

type VarStmt struct {
	Name token.NifToken
	Type token.NifToken
	Init Expr
}

func (*VarStmt) stmtNode()         {}
func (s *VarStmt) Pos() (int, int) { return s.Name.Line, s.Name.Column }

type ShortVarStmt struct {
	Name token.NifToken
	Init Expr
}

func (*ShortVarStmt) stmtNode()         {}
func (s *ShortVarStmt) Pos() (int, int) { return s.Name.Line, s.Name.Column }

type AssignStmt struct {
	Name  token.NifToken
	Value Expr
}

func (*AssignStmt) stmtNode()         {}
func (s *AssignStmt) Pos() (int, int) { return s.Name.Line, s.Name.Column }

type PrintStmt struct {
	Expr  Expr
	Print token.NifToken
}

func (*PrintStmt) stmtNode()         {}
func (s *PrintStmt) Pos() (int, int) { return s.Print.Line, s.Print.Column }

type ExprStmt struct {
	Expr Expr
}

func (*ExprStmt) stmtNode()         {}
func (s *ExprStmt) Pos() (int, int) { return s.Expr.Pos() }

type IfStmt struct {
	Conditon   Expr
	ThenBranch Stmt
	ElseBranch Stmt
	IfToken    token.NifToken
}

func (*IfStmt) stmtNode()         {}
func (s *IfStmt) Pos() (int, int) { return s.IfToken.Line, s.IfToken.Column }

type WhileStmt struct {
	Conditon Expr
	Body     Stmt
	While    token.NifToken
}

func (*WhileStmt) stmtNode()         {}
func (s *WhileStmt) Pos() (int, int) { return s.While.Line, s.While.Column }

type ForStmt struct {
	Init     Stmt
	CondExpr Expr
	Update   Stmt
	BodyStmt Stmt
	For      token.NifToken
}

func (*ForStmt) stmtNode()         {}
func (s *ForStmt) Pos() (int, int) { return s.For.Line, s.For.Column }

type BlockStmt struct {
	Statements []Stmt
	LBrace     token.NifToken
}

func (*BlockStmt) stmtNode()         {}
func (s *BlockStmt) Pos() (int, int) { return s.LBrace.Line, s.LBrace.Column }

type FuncStmt struct {
	Name   token.NifToken
	Params []Param
	Body   *BlockStmt
	Func   token.NifToken
}

func (*FuncStmt) stmtNode()         {}
func (s *FuncStmt) Pos() (int, int) { return s.Func.Line, s.Func.Column }

type StructStmt struct {
	Name    token.NifToken
	Fields  []VarStmt
	Methods []FuncStmt
	Struct  token.NifToken
}

func (*StructStmt) stmtNode()         {}
func (s *StructStmt) Pos() (int, int) { return s.Struct.Line, s.Struct.Column }

type ReturnStmt struct {
	Keyword token.NifToken
	Value   Expr
}

func (*ReturnStmt) stmtNode()         {}
func (s *ReturnStmt) Pos() (int, int) { return s.Keyword.Line, s.Keyword.Column }

type BreakStmt struct {
	Keyword token.NifToken
}

func (*BreakStmt) stmtNode()         {}
func (s *BreakStmt) Pos() (int, int) { return s.Keyword.Line, s.Keyword.Column }

type ContinueStmt struct {
	Keyword token.NifToken
}

func (*ContinueStmt) stmtNode()         {}
func (s *ContinueStmt) Pos() (int, int) { return s.Keyword.Line, s.Keyword.Column }
