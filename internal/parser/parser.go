package parser

import (
	"fmt"

	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	token "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"
)

type Parser struct {
	tokens  []token.NifToken
	current int
}

func New(tokens []token.NifToken) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() ([]ast.Stmt, error) {
	var statements []ast.Stmt
	for !p.isAtEnd() {
		p.skipnewLines()
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}
	return statements, nil
}

func (p *Parser) match(types ...token.TokenType) bool {
	for _, tt := range types {
		if p.check(tt) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.TokenEOF
}

func (p *Parser) peek() token.NifToken {
	return p.tokens[p.current]
}

func (p *Parser) advance() token.NifToken {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) previous() token.NifToken {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tt token.TokenType, message string) token.NifToken {
	if p.check(tt) {
		return p.advance()
	}
	panic(fmt.Sprintf("[Parse error] %s. Got '%s' at line %d", message, p.peek().Lexeme, p.peek().Line))
}

func (p *Parser) check(tt token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == tt
}

func (p *Parser) skipnewLines() {
	for p.check(token.TokenNewLine) {
		p.advance()
	}
}

func (p *Parser) consumeStatementEnd(message string) {
	if p.check(token.TokenNewLine) {
		p.advance()
	}
	if p.check(token.TokenRBrace) || p.isAtEnd() {
		return
	}

	panic(fmt.Sprintf("[Parse error] %s  got '%s' instead at line %d", message, p.peek().Lexeme, p.peek().Line))
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.orExpr()
}

func (p *Parser) orExpr() (ast.Expr, error) {
	left, err := p.andExpr()
	if err != nil {
		return nil, err
	}
	for p.match(token.TokenOr) {
		operator := p.previous()
		right, err := p.andExpr()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}
	return left, nil
}

func (p *Parser) andExpr() (ast.Expr, error) {
	left, err := p.equalityExpr()
	if err != nil {
		return nil, err
	}
	for p.match(token.TokenAnd) {
		operator := p.previous()
		right, err := p.equalityExpr()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}
	return left, nil

}

func (p *Parser) equalityExpr() (ast.Expr, error) {
	left, err := p.comparissonExpr()
	if err != nil {
		return nil, err
	}
	for p.match(token.TokenEqality, token.TokenBangEqal) {
		operator := p.previous()
		right, err := p.comparissonExpr()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}
	return left, nil
}

func (p *Parser) comparissonExpr() (ast.Expr, error) {
	left, err := p.termExpr()
	if err != nil {
		return nil, err
	}

	for p.match(
		token.TokenLess,
		token.TokenLessEq,
		token.TokenGreater,
		token.TokenGreaterEq,
	) {
		operator := p.previous()
		right, err := p.termExpr()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) termExpr() (ast.Expr, error) {
	left, err := p.factorExpr()
	if err != nil {
		return nil, err
	}

	for p.match(token.TokenPlus, token.TokenMinus) {
		operator := p.previous()
		right, err := p.factorExpr()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) factorExpr() (ast.Expr, error) {
	left, err := p.UnaryExpr()
	if err != nil {
		return nil, err
	}

	for p.match(token.TokenStar, token.TokenFWDSlash, token.TokenPercent) {
		operator := p.previous()
		right, err := p.UnaryExpr()
		if err != nil {
			return nil, err
		}
		left = &ast.BinaryExpr{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left, nil
}

func (p *Parser) UnaryExpr() (ast.Expr, error) {
	if p.match(token.TokenMinus, token.TokenBang) {
		operator := p.previous()
		right, err := p.UnaryExpr()
		if err != nil {
			return nil, err
		}
		return &ast.UnaryExpr{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.CallExpr()
}

func (p *Parser) varDeclaration() (ast.Stmt, error) {
	// p.consume(token.TokenVar, "Expect 'var' keyword")
	name := p.consume(token.TokenIdentifier, "expect variable name after var")

	typ := p.consume(token.TokenIdentifier, "expect type after variable name")

	p.consume(token.TokenAssign, "expect '=' after variable type")

	init, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.consumeStatementEnd("expect end of variable declaration")

	return &ast.VarStmt{
		Name: name,
		Type: typ,
		Init: init,
	}, nil
}

func (p *Parser) shortVarDeclaration() (ast.Stmt, error) {
	name := p.consume(token.TokenIdentifier, "expect variable name before ':='")

	p.consume(token.TokenColonEqual, "Expected ':=' for short variable declaration")

	init, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.consumeStatementEnd("expected end of short variable declaration.")

	return &ast.ShortVarStmt{
		Name: name,
		Init: init,
	}, nil
}

func (p *Parser) CallExpr() (ast.Expr, error) {
	expr, err := p.primaryExpr()
	if err != nil {
		return nil, err
	}
	for {
		if p.match(token.TokenLParen) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if p.match(token.TokenDot) {
			name := p.consume(token.TokenIdentifier, "Expected property name '.'")
			expr = &ast.GetExpr{
				Object: expr,
				Name:   name,
			}
		} else if p.match(token.TokenLBracket) {
			index, err := p.expression()
			if err != nil {
				return nil, err
			}
			braket := p.consume(token.TokenRBracket, "expected ']' after index")
			expr = &ast.IndexExpr{
				Collection: expr,
				Bracket:    braket,
				Index:      index,
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee ast.Expr) (ast.Expr, error) {
	var arguments []ast.Expr
	if !p.check(token.TokenRParen) {
		for {
			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)
			if !p.match(token.TokenComma) {
				break
			}
		}
	}
	paren := p.consume(token.TokenRParen, "expected ')' after arguments")
	return &ast.CallExpr{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}, nil
}

func (p *Parser) primaryExpr() (ast.Expr, error) {
	if p.match(token.TokenFalse) {
		return &ast.LiteralExpr{Value: p.previous()}, nil
	}
	if p.match(token.TokenTrue) {
		return &ast.LiteralExpr{Value: p.previous()}, nil
	}
	if p.match(token.TokenNull) {
		return &ast.LiteralExpr{Value: p.previous()}, nil
	}
	if p.match(token.TokenNumber) {
		return &ast.LiteralExpr{Value: p.previous()}, nil
	}
	if p.match(token.TokenString) {
		return &ast.LiteralExpr{Value: p.previous()}, nil
	}
	if p.match(token.TokenIdentifier) {
		return &ast.LiteralExpr{Value: p.previous()}, nil
	}
	if p.match(token.TokenLParen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		p.consume(token.TokenRParen, "expected ')' after expression")
		return expr, nil
	}
	if p.match(token.TokenLBracket) {
		return p.listLiteralExpr()
	}
	if p.match(token.TokenLBrace) {
		return p.dictLiteralExpr()
	}

	return nil, fmt.Errorf("unexpected token '%s' as line %d", p.peek().Lexeme, p.peek().Line)
}

func (p *Parser) listLiteralExpr() (ast.Expr, error) {
	var elements []ast.Expr
	if !p.check(token.TokenRBracket) {
		for {
			elem, err := p.expression()
			if err != nil {
				return nil, err
			}
			elements = append(elements, elem)
			if !p.match(token.TokenComma) {
				break
			}
		}
	}
	p.consume(token.TokenRBracket, "expected ']' after list elements")
	return &ast.ListExpr{
		Elements: elements,
	}, nil
}

func (p *Parser) dictLiteralExpr() (ast.Expr, error) {
	var pairs [][2]ast.Expr
	if !p.check(token.TokenRBrace) {
		for {
			key, err := p.expression()
			if err != nil {
				return nil, err
			}
			p.consume(token.TokenColon, "expect ':' after dictionary key")
			value, err := p.expression()
			if err != nil {
				return nil, err
			}
			pairs = append(pairs, [2]ast.Expr{key, value})
			if !p.match(token.TokenComma) {
				break
			}
		}
	}
	p.consume(token.TokenRBrace, "Expect '}' after dictionary entries")
	return &ast.DictExpr{
		Pairs: pairs,
	}, nil
}

func (p *Parser) assignmentStatement() (ast.Stmt, error) {
	name := p.consume(token.TokenIdentifier, "expect variable name before '='")

	p.consume(token.TokenAssign, "expect '=' for assignment statement")

	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.consumeStatementEnd("expected end of an assignment statement")

	return &ast.AssignStmt{
		Name:  name,
		Value: value,
	}, nil
}

func (p *Parser) printStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consumeStatementEnd("expect end of a print statement")
	return &ast.PrintStmt{
		Expr: expr,
	}, nil
}

func (p *Parser) returnStatement() (ast.Stmt, error) {
	keyword := p.previous()

	var value ast.Expr

	if !p.check(token.TokenNewLine) && !p.check(token.TokenRBrace) && !p.isAtEnd() {
		val, err := p.expression()
		if err != nil {
			return nil, err
		}
		value = val
	}
	p.consumeStatementEnd("expect end of a print statement")

	return &ast.ReturnStmt{
		Keyword: keyword,
		Value:   value,
	}, nil
}

func (p *Parser) ifStatement() (ast.Stmt, error) {
	cond, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.consume(token.TokenLBrace, "expect '{' after if condition")

	thenBranch, err := p.blockStatement()
	if err != nil {
		return nil, err
	}

	var elseBranch ast.Stmt
	if p.match(token.TokenElse) {
		p.consume(token.TokenLBrace, "expect '{' after else")
		elseBranch, err = p.blockStatement()
		if err != nil {
			return nil, err
		}

	}

	return &ast.IfStmt{
		Conditon:   cond,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

func (p *Parser) blockStatement() (*ast.BlockStmt, error) {
	lbrace := p.previous()
	var statements []ast.Stmt

	for {
		p.skipnewLines()

		if p.check(token.TokenRBrace) || p.isAtEnd() {
			break
		}
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			statements = append(statements, stmt)
		}
	}

	p.consume(token.TokenRBrace, "expected '}' after block")
	return &ast.BlockStmt{
		Statements: statements,
		LBrace:     lbrace,
	}, nil
}

func (p *Parser) continueStatement() (ast.Stmt, error) {
	keyword := p.previous()
	p.consumeStatementEnd("expect end of continue statement")
	return &ast.ContinueStmt{
		Keyword: keyword,
	}, nil
}

func (p *Parser) expressionStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.consumeStatementEnd("expected end of expression statement")
	return &ast.ExprStmt{
		Expr: expr,
	}, nil
}

func (p *Parser) breakStatement() (ast.Stmt, error) {
	keyword := p.previous()
	p.consumeStatementEnd("expected end of break statement")
	return &ast.BreakStmt{
		Keyword: keyword,
	}, nil
}

func (p *Parser) whileStatement() (ast.Stmt, error) {

	cond, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.consume(token.TokenLBrace, "expect '{' after while condition")
	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}
	return &ast.WhileStmt{
		Conditon: cond,
		Body:     body,
	}, nil
}

func (p *Parser) checkNext(tt token.TokenType) bool {
	if p.current+1 >= len(p.tokens) {
		return false
	}

	return p.tokens[p.current+1].Type == tt
}

func (p *Parser) funcDeclaration() (ast.Stmt, error) {
	funcTok := p.previous()

	name := p.consume(token.TokenIdentifier, "expected function name after 'func'")
	p.consume(token.TokenLParen, "expect '(' after function name")

	var params []ast.Param
	if !p.check(token.TokenRParen) {
		for {
			paramName := p.consume(token.TokenIdentifier, "expected parameter name")
			p.consume(token.TokenColon, "expected colon after parameter name")
			paramType := p.consume(token.TokenIdentifier, "expected ':' after parameter name")

			params = append(params, ast.Param{Name: paramName, Type: paramType})

			if !p.match(token.TokenComma) {
				break
			}
		}
	}
	p.consume(token.TokenRParen, "expect ')' after parameters")

	var returnType token.NifToken
	if p.match(token.TokenArrow) {
		returnType = p.consume(token.TokenIdentifier, "expect return type after '->'")
	}

	p.consume(token.TokenLBrace, "expect '{' before function body")
	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}

	return &ast.FuncStmt{
		Func:   funcTok,
		Name:   name,
		Params: params,
		Body:   body,
		Return: returnType,
	}, nil
}

// func (p *Parser) blockStatement() {}

func (p *Parser) funcExpression() (ast.Expr, error) {
	funcTok := p.previous()

	p.consume(token.TokenLParen, "expect '(' after func in function literal")

	var params []ast.Param

	if !p.check(token.TokenRParen) {
		for {
			paramName := p.consume(token.TokenIdentifier, "expected parameter name")
			p.consume(token.TokenColon, "expected ':' after parameter name")
			paramType := p.consume(token.TokenIdentifier, "expected parameter type")

			params = append(params, ast.Param{Name: paramName, Type: paramType})
			if !p.match(token.TokenComma) {
				break
			}
		}
	}
	p.consume(token.TokenRParen, "expect ')' after parameter list")

	p.consume(token.TokenLBrace, "expect '{' before function body in function literal")
	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}

	return &ast.FuncExpr{
		Func:   funcTok,
		Params: params,
		Body:   body,
	}, nil
}

func (p *Parser) forStatement() (ast.Stmt, error) {
	forTok := p.previous()
	var init ast.Stmt
	var err error
	if !p.check(token.TokenSemicolon) {
		if p.check(token.TokenVar) {
			init, err = p.varDeclaration()
		} else if p.check(token.TokenIdentifier) && p.checkNext(token.TokenColonEqual) {
			init, err = p.shortVarDeclaration()
		} else if p.check(token.TokenIdentifier) && p.checkNext(token.TokenAssign) {
			init, err = p.assignmentStatement()
		} else {
			init = nil
		}

		if err != nil {
			return nil, err
		}
	}
	p.consume(token.TokenSemicolon, "expected ';' after for init")

	var cond ast.Expr
	if !p.check(token.TokenSemicolon) {
		cond, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	p.consume(token.TokenSemicolon, "expected ';' after for init")

	var post ast.Stmt

	if !p.check(token.TokenLBrace) {
		if p.check(token.TokenIdentifier) && p.checkNext(token.TokenColonEqual) {
			post, err = p.shortVarDeclaration()
		} else if p.check(token.TokenIdentifier) && p.checkNext(token.TokenAssign) {
			post, err = p.assignmentStatement()
		}
		if err != nil {
			return nil, err
		}
	}

	p.consume(token.TokenLBrace, "expected '{' after for clause")
	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}

	return &ast.ForStmt{
		Init:     init,
		CondExpr: cond,
		Update:   post,
		BodyStmt: body,
		For:      forTok,
	}, nil
}

func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(token.TokenVar) {
		return p.varDeclaration()
	}
	if p.check(token.TokenIdentifier) {
		return p.shortVarDeclaration()
	}
	if p.check(token.TokenIdentifier) && p.checkNext(token.TokenAssign) {
		return p.assignmentStatement()
	}
	if p.match(token.TokenPrint) {
		return p.printStatement()
	}
	if p.match(token.TokenIf) {
		return p.ifStatement()
	}
	if p.match(token.TokenWhile) {
		return p.whileStatement()
	}
	if p.match(token.TokenFor) {
		return p.forStatement()
	}
	if p.match(token.TokenFunc) {
		return p.funcDeclaration()
	}
	if p.match(token.TokenReturn) {
		return p.returnStatement()
	}
	if p.match(token.TokenBreak) {
		return p.breakStatement()
	}
	if p.match(token.TokenContinue) {
		return p.continueStatement()
	}
	if p.match(token.TokenLBrace) {
		return p.expressionStatement()
	}

	return p.expressionStatement()
}
