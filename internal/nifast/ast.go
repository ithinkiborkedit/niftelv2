package nifast

import token "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"

type Expr interface {
	exprNode()
	Pos() (line, column int)
}

type TypeExpr struct {
	Name     token.Token
	TypeArgs []TypeExpr
}

type Stmt interface {
	stmtNode()
	Pos() (line, column int)
}

type BinaryExpr struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (*BinaryExpr) exprNode()         {}
func (e *BinaryExpr) Pos() (int, int) { return e.Operator.Line, e.Operator.Column }

type UnaryExpr struct {
	Operator token.Token
	Right    Expr
}

func (*UnaryExpr) exprNode()         {}
func (e *UnaryExpr) Pos() (int, int) { return e.Operator.Line, e.Operator.Column }

type VariableExpr struct {
	Name token.Token
}

func (*VariableExpr) exprNode()         {}
func (e *VariableExpr) Pos() (int, int) { return e.Name.Line, e.Name.Column }

type LiteralExpr struct {
	Value token.Token
}

func (*LiteralExpr) exprNode()         {}
func (e *LiteralExpr) Pos() (int, int) { return e.Value.Line, e.Value.Column }

type CallExpr struct {
	Callee    Expr
	Paren     token.Token
	Arguments []Expr
	TypeArgs  []*TypeExpr
}

func (*CallExpr) exprNode()         {}
func (e *CallExpr) Pos() (int, int) { return e.Paren.Line, e.Paren.Column }

type IndexExpr struct {
	Collection Expr
	Bracket    token.Token
	Index      Expr
}

func (*IndexExpr) exprNode()         {}
func (e *IndexExpr) Pos() (int, int) { return e.Bracket.Line, e.Bracket.Column }

type GetExpr struct {
	Object Expr
	Name   token.Token
}

func (*GetExpr) exprNode()         {}
func (e *GetExpr) Pos() (int, int) { return e.Name.Line, e.Name.Column }

type ListExpr struct {
	Elements []Expr
	LBracket token.Token
}

func (*ListExpr) exprNode()         {}
func (e *ListExpr) Pos() (int, int) { return e.LBracket.Line, e.LBracket.Column }

type DictExpr struct {
	Pairs  [][2]Expr
	LBrace token.Token
}

func (*DictExpr) exprNode()         {}
func (e *DictExpr) Pos() (int, int) { return e.LBrace.Line, e.LBrace.Column }

type StructLiteralExpr struct {
	TypeName *TypeExpr
	Fields   map[string]Expr
	LBrace   token.Token
}

func (*StructLiteralExpr) exprNode()         {}
func (e *StructLiteralExpr) Pos() (int, int) { return e.LBrace.Line, e.LBrace.Column }

type FuncExpr struct {
	Params []Param
	Body   *BlockStmt
	Func   token.Token
}

type Param struct {
	Name token.Token
	Type *TypeExpr
}

func (*FuncExpr) exprNode()         {}
func (e *FuncExpr) Pos() (int, int) { return e.Func.Line, e.Func.Column }

//STATEMENTS

type VarStmt struct {
	Names []token.Token
	Type  *TypeExpr
	Init  Expr
}

func (*VarStmt) stmtNode() {}
func (s *VarStmt) Pos() (int, int) {
	if len(s.Names) > 0 {
		return s.Names[0].Line, s.Names[0].Column
	}
	return 0, 0
}

type ShortVarStmt struct {
	Name token.Token
	Init Expr
}

func (*ShortVarStmt) stmtNode()         {}
func (s *ShortVarStmt) Pos() (int, int) { return s.Name.Line, s.Name.Column }

type AssignStmt struct {
	Name  token.Token
	Value Expr
}

func (*AssignStmt) stmtNode()         {}
func (s *AssignStmt) Pos() (int, int) { return s.Name.Line, s.Name.Column }

type PrintStmt struct {
	Expr  Expr
	Print token.Token
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
	IfToken    token.Token
}

func (*IfStmt) stmtNode()         {}
func (s *IfStmt) Pos() (int, int) { return s.IfToken.Line, s.IfToken.Column }

type WhileStmt struct {
	Conditon Expr
	Body     Stmt
	While    token.Token
}

func (*WhileStmt) stmtNode()         {}
func (s *WhileStmt) Pos() (int, int) { return s.While.Line, s.While.Column }

type ForStmt struct {
	Init     Stmt
	CondExpr Expr
	Update   Stmt
	BodyStmt Stmt
	For      token.Token
}

func (*ForStmt) stmtNode()         {}
func (s *ForStmt) Pos() (int, int) { return s.For.Line, s.For.Column }

type BlockStmt struct {
	Statements []Stmt
	LBrace     token.Token
}

func (*BlockStmt) stmtNode()         {}
func (s *BlockStmt) Pos() (int, int) { return s.LBrace.Line, s.LBrace.Column }

type FuncStmt struct {
	Name        token.Token
	Params      []Param
	TypeParams  []token.Token
	ReturnTypes []*TypeExpr
	Body        *BlockStmt
	Return      []*TypeExpr
	Func        token.Token
}

func (*FuncStmt) stmtNode()         {}
func (s *FuncStmt) Pos() (int, int) { return s.Func.Line, s.Func.Column }

type StructStmt struct {
	Name       token.Token
	TypeParams []token.Token
	Fields     []VarStmt
	Methods    []FuncStmt
	Struct     token.Token
}

func (*StructStmt) stmtNode()         {}
func (s *StructStmt) Pos() (int, int) { return s.Struct.Line, s.Struct.Column }

type ReturnStmt struct {
	Keyword token.Token
	Values  []Expr
}

func (*ReturnStmt) stmtNode()         {}
func (s *ReturnStmt) Pos() (int, int) { return s.Keyword.Line, s.Keyword.Column }

type BreakStmt struct {
	Keyword token.Token
}

func (*BreakStmt) stmtNode()         {}
func (s *BreakStmt) Pos() (int, int) { return s.Keyword.Line, s.Keyword.Column }

type ContinueStmt struct {
	Keyword token.Token
}

func (*ContinueStmt) stmtNode()         {}
func (s *ContinueStmt) Pos() (int, int) { return s.Keyword.Line, s.Keyword.Column }
