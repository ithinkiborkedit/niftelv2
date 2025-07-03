package parser

import (
	"errors"
	"fmt"

	"github.com/ithinkiborkedit/niftelv2.git/internal/lexer"
	ast "github.com/ithinkiborkedit/niftelv2.git/internal/nifast"
	token "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"
)

var ErrIncomplete = errors.New("incomplete input")

type Parser struct {
	src      lexer.TokenSource
	curr     token.Token
	ahead    token.Token
	hasAhead bool
	prev     token.Token
	err      error
}

func New(src lexer.TokenSource) *Parser {
	p := &Parser{
		src: src,
	}
	p.advance()
	return p
}

func (p *Parser) Parse() ([]ast.Stmt, error) {
	var statements []ast.Stmt
	err := p.skipnewLines()
	if err != nil {
		return nil, err
	}
	for !p.isAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			statements = append(statements, stmt)
		}
		err = p.skipnewLines()
		if err != nil {
			return nil, err
		}
	}
	return statements, nil
}

func (p *Parser) match(types ...token.TokenType) (bool, error) {
	for _, tt := range types {
		if p.check(tt) {
			if err := p.advance(); err != nil {
				return false, err
			}
			return true, nil
		}
	}

	return false, nil
}

func (p *Parser) isAtEnd() bool {
	return p.curr.Type == token.TokenEOF
}

func (p *Parser) peek() (token.Token, error) {
	if !p.hasAhead {
		tok, err := p.src.NextToken()
		if err != nil {
			return token.Token{}, err
		}
		p.ahead = tok
		p.hasAhead = true
	}
	return p.ahead, nil
}

func (p *Parser) advance() error {
	p.prev = p.curr
	if p.hasAhead {
		p.curr = p.ahead
		p.hasAhead = false
	} else {
		tok, err := p.src.NextToken()
		if err != nil {
			return err
		}
		p.curr = tok
	}
	// if !p.isAtEnd() {
	// 	p.current++
	// }
	// return p.previous()
	fmt.Printf("[ADVANCE] Current: %v (type=%v, line=%v, lexeme=%q)", p.curr.Type, p.curr.Type, p.curr.Line, p.curr.Lexeme)
	return nil
}

func (p *Parser) previous() token.Token {
	return p.prev
	// return p.tokens[p.current-1]
}

func (p *Parser) consume(tt token.TokenType, message string) (token.Token, error) {
	if p.check(tt) {
		tok := p.curr
		if err := p.advance(); err != nil {
			return token.Token{}, err
		}
		return tok, nil
	}

	return token.Token{}, fmt.Errorf("[Parse error] %s. Got '%s' at line %d", message, p.curr.Lexeme, p.curr.Line)
	// 	return p.advance()
	// }
	// panic(fmt.Sprintf("[Parse error] %s. Got '%s' at line %d", message, p.peek().Lexeme, p.peek().Line))
}

func (p *Parser) check(tt token.TokenType) bool {
	return p.curr.Type == tt
}

func (p *Parser) skipnewLines() error {
	for p.check(token.TokenNewLine) {
		if err := p.advance(); err != nil {
			return err
		}
	}
	return nil
}

// func (p *Parser) skipnewLines() {
// 	for p.check(token.TokenNewLine) {
// 		p.advance()
// 	}
// }

func (p *Parser) expression() (ast.Expr, error) {
	return p.orExpr()
}

func (p *Parser) orExpr() (ast.Expr, error) {
	left, err := p.andExpr()
	if err != nil {
		return nil, err
	}
	for {
		m, err := p.match(token.TokenOr)
		if err != nil {
			return nil, err
		}
		if !m {
			break
		}
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
	for {
		m, err := p.match(token.TokenAnd)
		if err != nil {
			return nil, err
		}
		if !m {
			break
		}
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
	for {
		m, err := p.match(token.TokenEqality)
		if err != nil {
			return nil, err
		}
		if !m {
			break
		}
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

func (p *Parser) parseTypeExpr() (*ast.TypeExpr, error) {
	fmt.Printf("parseTypeExpr: current token = %v\n", p.curr)
	name, err := p.consume(token.TokenIdentifier, "expected type name")
	if err != nil {
		return nil, err
	}
	typeExpr := &ast.TypeExpr{
		Name:     name,
		TypeArgs: nil,
	}
	if p.check(token.TokenLBracket) {
		_, err := p.consume(token.TokenLBracket, "expected '[' after type name for type arguments")
		if err != nil {
			return nil, err
		}
		var args []ast.TypeExpr
		for {
			arg, err := p.parseTypeExpr()
			if err != nil {
				return nil, err
			}
			args = append(args, *arg)
			ok, err := p.match(token.TokenComma)
			if err != nil {
				return nil, err
			}
			if !ok {
				break
			}
		}
		_, err = p.consume(token.TokenRBracket, "expected ']' afte type arguments")
		if err != nil {
			return nil, err
		}
		typeExpr.TypeArgs = args
	}
	return typeExpr, nil
}

func (p *Parser) comparissonExpr() (ast.Expr, error) {
	left, err := p.termExpr()
	if err != nil {
		return nil, err
	}
	for {
		m, err := p.match(token.TokenLess,
			token.TokenLessEq,
			token.TokenGreater,
			token.TokenGreaterEq)
		if err != nil {
			return nil, err
		}
		if !m {
			break
		}
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
	for {
		m, err := p.match(token.TokenPlus, token.TokenMinus)
		if err != nil {
			return nil, err
		}
		if !m {
			break
		}
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

	for {
		m, err := p.match(token.TokenStar, token.TokenFWDSlash, token.TokenPercent)
		if err != nil {
			return nil, err
		}
		if !m {
			break
		}
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
	m, err := p.match(token.TokenMinus, token.TokenBang)
	if err != nil {
		return nil, err
	}
	if m {
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
	var names []token.Token
	name, err := p.consume(token.TokenIdentifier, "expect variables name after var")
	if err != nil {
		return nil, err
	}
	names = append(names, name)
	for {
		ok, err := p.match(token.TokenComma)
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		next, err := p.consume(token.TokenIdentifier, "expected varibal name after ','")
		if err != nil {
			return nil, err
		}
		names = append(names, next)
	}
	var typ *ast.TypeExpr
	if len(names) == 1 && p.check(token.TokenColon) {

		_, err = p.consume(token.TokenColon, "expected ':' after variable name")
		if err != nil {
			return nil, err
		}
		typ, err = p.parseTypeExpr()
		if err != nil {
			return nil, err
		}

	}
	_, err = p.consume(token.TokenAssign, "expect '=' after variable(s)")
	if err != nil {
		return nil, err
	}

	init, err := p.expression()
	if err != nil {
		return nil, err
	}
	return &ast.VarStmt{
		Names: names,
		Type:  typ,
		Init:  init,
	}, nil
	// p.consume(token.TokenVar, "Expect 'var' keyword")
	// name, err := p.consume(token.TokenIdentifier, "expect variable name after var")
	// if err != nil {
	// 	return nil, err
	// }
	// _, err = p.consume(token.TokenColon, "expected ':' after variable name")
	// if err != nil {
	// 	return nil, err
	// }
	// typ, err := p.consume(token.TokenIdentifier, "expect type after variable name")
	// if err != nil {
	// 	return nil, err
	// }

	// _, err = p.consume(token.TokenAssign, "expect '=' after variable type")
	// if err != nil {
	// 	return nil, err
	// }

	// init, err := p.expression()
	// if err != nil {
	// 	return nil, err
	// }

	// err = p.skipnewLines()
	// if err != nil {
	// 	return nil, err
	// }

	// return &ast.VarStmt{
	// 	Names: []token.Token{name},
	// 	Type:  typ,
	// 	Init:  init,
	// }, nil
}

func (p *Parser) shortVarDeclaration() (ast.Stmt, error) {
	name, err := p.consume(token.TokenIdentifier, "expect variable name before ':='")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.TokenColonEqual, "Expected ':=' for short variable declaration")
	if err != nil {
		return nil, err
	}

	init, err := p.expression()
	if err != nil {
		return nil, err
	}

	err = p.skipnewLines()
	if err != nil {
		return nil, err
	}

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
		var typeArgs []*ast.TypeExpr
		if p.check(token.TokenLBracket) {

			_, err := p.consume(token.TokenLBracket, "expected '[' for generic call type arguments")
			if err != nil {
				return nil, err
			}
			for {
				typ, err := p.parseTypeExpr()
				if err != nil {
					return nil, err
				}
				typeArgs = append(typeArgs, typ)
				ok, err := p.match(token.TokenComma)
				if err != nil {
					return nil, err
				}
				if !ok {
					break
				}
			}
			_, err = p.consume(token.TokenRBracket, "expected ']' after type arguments")
			if err != nil {
				return nil, err
			}
			if !p.check(token.TokenLParen) {
				return nil, fmt.Errorf("Type arguments only allowed in function calls e.g(foo[T](args...)")
			}
		}
		ok, err := p.match(token.TokenLParen)
		if err != nil {
			return nil, err
		}
		if ok {
			expr, err = p.finishCall(expr, typeArgs)
			if err != nil {
				return nil, err
			}
			continue
		}

		ok, err = p.match(token.TokenDot)
		if err != nil {
			return nil, err
		}
		if ok {
			name, err := p.consume(token.TokenIdentifier, "Expected property name after '.'")
			if err != nil {
				return nil, err
			}
			expr = &ast.GetExpr{
				Object: expr,
				Name:   name,
			}
			continue
		}

		ok, err = p.match(token.TokenLBracket)
		if err != nil {
			return nil, err
		}
		if ok {
			index, err := p.expression()
			if err != nil {
				return nil, err
			}
			bracket, err := p.consume(token.TokenRBracket, "expected ']' after index")
			if err != nil {
				if p.curr.Type == token.TokenEOF {
					return nil, ErrIncomplete
				}
				return nil, err
			}
			expr = &ast.IndexExpr{
				Collection: expr,
				Bracket:    bracket,
				Index:      index,
			}
			continue
		}

		break
	}
	return expr, nil
}

func (p *Parser) finishCall(callee ast.Expr, typeArgs []*ast.TypeExpr) (ast.Expr, error) {
	var arguments []ast.Expr
	if !p.check(token.TokenRParen) {
		for {
			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, arg)
			ok, err := p.match(token.TokenComma)
			if err != nil {
				return nil, err
			}
			if !ok {
				break
			}
		}
	}
	paren, err := p.consume(token.TokenRParen, "expected ')' after arguments")
	if err != nil {
		if p.curr.Type == token.TokenEOF {
			return nil, ErrIncomplete
		}
		return nil, err
	}
	fmt.Printf("DEBUG CALL callee=%#v, typeArg=%#v, arguments=%#v", callee, typeArgs, arguments)
	return &ast.CallExpr{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
		TypeArgs:  typeArgs,
	}, nil
}

func (p *Parser) structLiteralExpr(typeName *ast.TypeExpr) (ast.Expr, error) {
	lbrace, err := p.consume(token.TokenLBrace, "expected '{' after struct type name")
	if err != nil {
		return nil, err
	}

	fields := make(map[string]ast.Expr)

	for {
		err := p.skipnewLines()
		if err != nil {
			return nil, err
		}
		if p.check(token.TokenRBrace) {
			break
		}

		fieldTok, err := p.consume(token.TokenIdentifier, "expected field name in struct literal")
		if err != nil {
			return nil, err
		}

		_, err = p.consume(token.TokenColon, "expected ':' after struct field name")
		if err != nil {
			return nil, err
		}

		valExpr, err := p.expression()
		if err != nil {
			return nil, err
		}
		fields[fieldTok.Lexeme] = valExpr
		if ok, _ := p.match(token.TokenComma); ok {
			continue
		}
		_ = p.skipnewLines()
		if p.check(token.TokenRBrace) {
			break
		}
	}

	_, err = p.consume(token.TokenRBrace, "expected '}' after struct literal fields")
	if err != nil {
		return nil, err
	}

	return &ast.StructLiteralExpr{
		TypeName: typeName,
		Fields:   fields,
		LBrace:   lbrace,
	}, nil

}

func (p *Parser) primaryExpr() (ast.Expr, error) {
	ok, err := p.match(token.TokenFunc)
	if err != nil {
		return nil, err
	}
	if ok {
		return p.funcExpression()
	}
	ok, err = p.match(token.TokenFalse)
	if err != nil {
		return nil, err
	}
	if ok {
		return &ast.LiteralExpr{Value: p.prev}, nil
	}

	ok, err = p.match(token.TokenTrue)
	if err != nil {
		return nil, err
	}
	if ok {
		return &ast.LiteralExpr{Value: p.prev}, nil
	}

	ok, err = p.match(token.TokenNull)
	if err != nil {
		return nil, err
	}
	if ok {
		return &ast.LiteralExpr{Value: p.prev}, nil
	}
	ok, err = p.match(token.TokenNumber)
	if err != nil {
		return nil, err
	}
	if ok {
		return &ast.LiteralExpr{Value: p.prev}, nil
	}
	ok, err = p.match(token.TokenString)
	if err != nil {
		return nil, err
	}
	if ok {
		return &ast.LiteralExpr{Value: p.prev}, nil
	}
	ok, err = p.match(token.TokenIdentifier)
	if err != nil {
		return nil, err
	}
	if ok {
		if p.check(token.TokenLBrace) {
			typeExpr, err := p.parseTypeExprFromToken(p.prev)
			if err != nil {
				return nil, err
			}
			return p.structLiteralExpr(typeExpr)
		}
		// typeExpr, err := p.parseTypeExprFromToken(p.prev)
		// if err != nil {
		// 	return nil, err
		// }
		// if p.check(token.TokenLBrace) {
		// 	return p.structLiteralExpr(typeExpr)
		// }
		return &ast.VariableExpr{Name: p.prev}, nil
	}
	// ok, err = p.match(token.TokenIdentifier)
	// if err != nil {
	// 	return nil, err
	// }
	// if ok {
	// 	return &ast.VariableExpr{Name: p.prev}, nil
	// }
	ok, err = p.match(token.TokenLParen)
	if err != nil {
		return nil, err
	}
	if ok {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.TokenRParen, "expected ')' after expression")
		if err != nil {
			if p.curr.Type == token.TokenEOF {
				return nil, ErrIncomplete
			}
			return nil, err
		}
		return expr, nil
	}
	ok, err = p.match(token.TokenLBracket)
	if err != nil {
		return nil, err
	}
	if ok {
		return p.listLiteralExpr()
	}
	ok, err = p.match(token.TokenLBrace)
	if err != nil {
		return nil, err
	}
	if ok {
		return p.dictLiteralExpr()
	}
	return nil, fmt.Errorf("unexpected token '%s' as line %d", p.curr.Lexeme, p.curr.Line)
}

func (p *Parser) parseTypeExprFromToken(tok token.Token) (*ast.TypeExpr, error) {
	typeExpr := &ast.TypeExpr{
		Name:     tok,
		TypeArgs: nil,
	}

	if p.check(token.TokenLBracket) {
		_, err := p.consume(token.TokenLBracket, "expected '[' after type name for type arguments")
		if err != nil {
			return nil, err
		}
		var args []ast.TypeExpr
		for {
			arg, err := p.parseTypeExpr()
			if err != nil {
				return nil, err
			}
			args = append(args, *arg)
			ok, err := p.match(token.TokenComma)
			if err != nil {
				return nil, err
			}
			if !ok {
				break
			}
		}
		_, err = p.consume(token.TokenRBracket, "expectd '[' after type arguments")
		if err != nil {
			return nil, err
		}
		typeExpr.TypeArgs = args
	}
	return typeExpr, nil
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
			ok, err := p.match(token.TokenComma)
			if err != nil {
				return nil, err
			}
			if !ok {
				break
			}
		}
	}
	_, err := p.consume(token.TokenRBracket, "expected ']' after list elements")
	if err != nil {
		if p.curr.Type == token.TokenEOF {
			return nil, ErrIncomplete
		}
		return nil, err
	}
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
			_, err = p.consume(token.TokenColon, "expect ':' after dictionary key")
			if err != nil {
				return nil, err
			}
			value, err := p.expression()
			if err != nil {
				return nil, err
			}
			pairs = append(pairs, [2]ast.Expr{key, value})
			ok, err := p.match(token.TokenComma)
			if err != nil {
				return nil, err
			}
			if !ok {
				break
			}
		}
	}
	_, err := p.consume(token.TokenRBrace, "Expect '}' after dictionary entries")
	if err != nil {
		if p.curr.Type == token.TokenEOF {
			return nil, ErrIncomplete
		}
		return nil, err
	}
	return &ast.DictExpr{
		Pairs: pairs,
	}, nil
}

func (p *Parser) assignmentStatement() (ast.Stmt, error) {
	name, err := p.consume(token.TokenIdentifier, "expect variable name before '='")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.TokenAssign, "expect '=' for assignment statement")
	if err != nil {
		return nil, err
	}

	value, err := p.expression()
	if err != nil {
		return nil, err
	}

	err = p.skipnewLines()
	if err != nil {
		return nil, err
	}

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
	err = p.skipnewLines()
	if err != nil {
		return nil, err
	}
	return &ast.PrintStmt{
		Expr: expr,
	}, nil
}

func (p *Parser) returnStatement() (ast.Stmt, error) {
	keyword := p.previous()
	var values []ast.Expr

	// Only parse a value if not immediately at newline, '}', or EOF.
	if !p.check(token.TokenNewLine) && !p.check(token.TokenRBrace) && !p.isAtEnd() {
		first, err := p.expression()
		if err != nil {
			return nil, err
		}
		values = append(values, first)

		// Now, keep parsing comma-separated values
		for p.check(token.TokenComma) {
			_, err := p.match(token.TokenComma)
			if err != nil {
				return nil, err
			}
			next, err := p.expression()
			if err != nil {
				return nil, err
			}
			values = append(values, next)
		}
	}

	// After all return values, skip optional newlines.
	if err := p.skipnewLines(); err != nil {
		return nil, err
	}

	fmt.Printf("[RETURN] After parsing expr, current=%v (type=%v, lexeme=%q)\n", p.curr.Type, p.curr.Type, p.curr.Lexeme)
	return &ast.ReturnStmt{
		Keyword: keyword,
		Values:  values,
	}, nil
}

func (p *Parser) ifStatement() (ast.Stmt, error) {
	cond, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.TokenLBrace, "expect '{' after if condition")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.blockStatement()
	if err != nil {
		return nil, err
	}

	var elseBranch ast.Stmt
	ok, err := p.match(token.TokenElse)
	if err != nil {
		return nil, err
	}
	if ok {
		_, err = p.consume(token.TokenLBrace, "expect '{' after else")
		if err != nil {
			return nil, err
		}
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

	err := p.skipnewLines()
	if err != nil {
		return nil, err
	}

	for {
		// err := p.skipnewLines()
		// if err != nil {
		// 	return nil, err
		// }
		fmt.Printf("[BLOCK] Next token: type=%v lexeme=%q line=%d\n", p.curr.Type, p.curr.Lexeme, p.curr.Line)
		if p.check(token.TokenRBrace) || p.isAtEnd() {
			break
		}
		// if p.check(token.TokenRBrace) || p.isAtEnd() {
		// 	break
		// }
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			statements = append(statements, stmt)
		}
		err = p.skipnewLines()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(token.TokenRBrace, "expected '}' after block")
	if err != nil {
		if p.curr.Type == token.TokenEOF {
			return nil, ErrIncomplete
		}
		return nil, err
	}

	return &ast.BlockStmt{
		Statements: statements,
		LBrace:     lbrace,
	}, nil

}

func (p *Parser) continueStatement() (ast.Stmt, error) {
	keyword := p.previous()
	err := p.skipnewLines()
	if err != nil {
		return nil, err
	}
	return &ast.ContinueStmt{
		Keyword: keyword,
	}, nil
}

func (p *Parser) expressionStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	err = p.skipnewLines()
	if err != nil {
		return nil, err
	}
	return &ast.ExprStmt{
		Expr: expr,
	}, nil
}

func (p *Parser) breakStatement() (ast.Stmt, error) {
	keyword := p.previous()
	err := p.skipnewLines()
	if err != nil {
		return nil, err
	}
	return &ast.BreakStmt{
		Keyword: keyword,
	}, nil
}

func (p *Parser) whileStatement() (ast.Stmt, error) {

	cond, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.TokenLBrace, "expect '{' after while condition")
	if err != nil {
		return nil, err
	}
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
	tok, err := p.peek()
	if err != nil {
		return false
	}
	p.ahead = tok
	p.hasAhead = true

	return p.ahead.Type == tt
}

func (p *Parser) funcDeclaration() (ast.Stmt, error) {
	funcTok := p.previous()

	name, err := p.consume(token.TokenIdentifier, "expected function name after 'func'")
	if err != nil {
		return nil, err
	}
	var typeParams []token.Token
	if p.check(token.TokenLBracket) {
		_, err := p.consume(token.TokenLBracket, "expected '[' after function name for type params")
		if err != nil {
			return nil, err
		}
		for {
			param, err := p.consume(token.TokenIdentifier, "expected generic type params")
			if err != nil {
				return nil, err
			}
			typeParams = append(typeParams, param)
			ok, err := p.match(token.TokenComma)
			if err != nil {
				return nil, err
			}
			if !ok {
				break
			}
		}
		_, err = p.consume(token.TokenRBracket, "expeded ']' after generic type parameters")
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.TokenLParen, "expect '(' after function name")
	if err != nil {
		return nil, err
	}

	var params []ast.Param
	if !p.check(token.TokenRParen) {
		for {
			paramName, err := p.consume(token.TokenIdentifier, "expected parameter name")
			if err != nil {
				return nil, err
			}
			_, err = p.consume(token.TokenColon, "expected colon after parameter name")
			if err != nil {
				return nil, err
			}
			paramType, err := p.parseTypeExpr()
			if err != nil {
				return nil, err
			}

			params = append(params, ast.Param{Name: paramName, Type: paramType})

			ok, err := p.match(token.TokenComma)
			if err != nil {
				return nil, err
			}
			if !ok {
				break
			}
		}
	}
	_, err = p.consume(token.TokenRParen, "expect ')' after parameters")
	if err != nil {
		if p.curr.Type == token.TokenEOF {
			return nil, ErrIncomplete
		}
		return nil, err
	}

	var returnTypes []*ast.TypeExpr
	ok, err := p.match(token.TokenArrow)
	if err != nil {
		return nil, err
	}
	if ok {
		for p.check(token.TokenNewLine) {
			if err := p.advance(); err != nil {
				return nil, err
			}
		}
		if p.check(token.TokenIdentifier) {
			typ, err := p.parseTypeExpr()
			if err != nil {
				return nil, err
			}
			returnTypes = append(returnTypes, typ)
		} else if p.check(token.TokenLParen) {
			_, err := p.consume(token.TokenLParen, "expected '(' after -> for multiple return types")
			if err != nil {
				return nil, err
			}
			for {
				typ, err := p.parseTypeExpr()
				if err != nil {
					return nil, err
				}
				returnTypes = append(returnTypes, typ)
				ok, err := p.match(token.TokenComma)
				if err != nil {
					return nil, err
				}
				if !ok {
					break
				}
			}
			_, err = p.consume(token.TokenRParen, "expected ')' after multple return types")
			if err != nil {
				return nil, err
			}
		}
		for p.check(token.TokenNewLine) {
			if err := p.advance(); err != nil {
				return nil, err
			}
		}
		_, err = p.consume(token.TokenLBrace, "expected '{' before function body")
		if err != nil {
			return nil, err
		}
		body, err := p.blockStatement()
		if err != nil {
			return nil, err
		}
		return &ast.FuncStmt{
			Func:       funcTok,
			Name:       name,
			Params:     params,
			TypeParams: typeParams,
			Body:       body,
			Return:     returnTypes,
		}, nil
	}

	_, err = p.consume(token.TokenLBrace, "expect '{' before function body")
	if err != nil {
		return nil, err
	}
	body, err := p.blockStatement()
	if err != nil {
		return nil, err
	}

	return &ast.FuncStmt{
		Func:       funcTok,
		Name:       name,
		Params:     params,
		TypeParams: typeParams,
		Body:       body,
		Return:     returnTypes,
	}, nil
}

// func (p *Parser) blockStatement() {}

func (p *Parser) funcExpression() (ast.Expr, error) {
	funcTok := p.previous()
	_, err := p.consume(token.TokenLParen, "expect '(' after func in function literal")
	if err != nil {
		return nil, err
	}

	var params []ast.Param

	if !p.check(token.TokenRParen) {
		for {
			paramName, err := p.consume(token.TokenIdentifier, "expected parameter name")
			if err != nil {
				return nil, err
			}
			_, err = p.consume(token.TokenColon, "expected ':' after parameter name")
			if err != nil {
				return nil, err
			}
			paramType, err := p.parseTypeExpr()
			if err != nil {
				return nil, err
			}

			params = append(params, ast.Param{Name: paramName, Type: paramType})
			ok, err := p.match(token.TokenComma)
			if err != nil {
				return nil, err
			}
			if !ok {
				break
			}
		}
	}

	_, err = p.consume(token.TokenRParen, "expect ')' after parameter list")
	if err != nil {
		if p.curr.Type == token.TokenEOF {
			return nil, ErrIncomplete
		}
		return nil, err
	}

	_, err = p.consume(token.TokenLBrace, "expect '{' before function body in function literal")
	if err != nil {
		return nil, err
	}
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

func (p *Parser) structDeclartion() (ast.Stmt, error) {

	fmt.Printf("[AFTER SKIP NEW LINES] struct body: token %v lexem='%s' line='%d'\n", p.curr.Type, p.curr.Lexeme, p.curr.Line)
	structTok := p.previous()

	name, err := p.consume(token.TokenIdentifier, "expected a struct name after 'struct'")
	if err != nil {
		return nil, err
	}
	var typeParams []token.Token
	if p.check(token.TokenLBracket) {
		_, err := p.consume(token.TokenLBracket, "expected '[' after struct name for type parameters")
		if err != nil {
			return nil, err
		}
		for {
			param, err := p.consume(token.TokenIdentifier, "expected type parameter")
			if err != nil {
				return nil, err
			}
			typeParams = append(typeParams, param)
			ok, err := p.match(token.TokenComma)
			if err != nil {
				return nil, err
			}
			if !ok {
				break
			}
		}
		_, err = p.consume(token.TokenRBracket, "expected ']' after type paramters")
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.TokenLBrace, "expected '{' after struct name")
	fmt.Printf("[AFTER OPENING '{' ] struct body: token %v lexem='%s' line='%d'\n", p.curr.Type, p.curr.Lexeme, p.curr.Line)
	if err != nil {
		return nil, err
	}

	for p.check(token.TokenNewLine) {
		if err := p.advance(); err != nil {
			return nil, err
		}
	}
	fmt.Printf("[Trying to check for \n after '{' ] struct body: token %v lexem='%s' line='%d'\n", p.curr.Type, p.curr.Lexeme, p.curr.Line)
	// err = p.skipnewLines()
	// if err != nil {
	// 	return nil, err
	// }

	var fields []ast.VarStmt
	var methods []ast.FuncStmt

	for !p.check(token.TokenRBrace) && !p.isAtEnd() {
		err := p.skipnewLines()
		if err != nil {
			return nil, err
		}
		fmt.Printf("[PARSERLOOP] struct body: type=%d  lexem='%s' line='%d'\n", p.curr.Type, p.curr.Lexeme, p.curr.Line)
		// Allow and skip any number of blank lines or newlines
		err = p.skipnewLines()
		if err != nil {
			return nil, err
		}

		// After skipping, if at end of struct, break
		if p.check(token.TokenRBrace) {
			break
		}

		// Parse methods (functions)
		ok, err := p.match(token.TokenFunc)
		if err != nil {
			return nil, err
		}
		if ok {
			method, err := p.funcDeclaration()
			if err != nil {
				return nil, err
			}
			fn, ok := method.(*ast.FuncStmt)
			if !ok {
				return nil, fmt.Errorf("expected function statement in struct body")
			}
			methods = append(methods, *fn)
			err = p.skipnewLines()
			if err != nil {
				return nil, err
			}
			continue
		}

		// Parse fields
		if p.check(token.TokenIdentifier) {
			fieldName, err := p.consume(token.TokenIdentifier, "expected field name in struct")
			if err != nil {
				return nil, err
			}
			_, err = p.consume(token.TokenColon, "expected ':' after struct field name")
			if err != nil {
				return nil, err
			}
			fieldType, err := p.parseTypeExpr()
			if err != nil {
				return nil, err
			}
			fields = append(fields, ast.VarStmt{
				Names: []token.Token{fieldName},
				Type:  fieldType,
			})
			err = p.skipnewLines()
			if err != nil {
				return nil, err
			}
			continue
		}

		// Anything else: error
		return nil, fmt.Errorf("unexpected token in struct body: '%s' at line %d", p.curr.Lexeme, p.curr.Line)
	}

	err = p.skipnewLines()
	if err != nil {
		return nil, err
	}

	// Must see a closing brace
	_, err = p.consume(token.TokenRBrace, "expected '}' after struct body")
	if err != nil {
		return nil, err
	}

	return &ast.StructStmt{
		Name:    name,
		Fields:  fields,
		Methods: methods,
		Struct:  structTok,
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
	_, err = p.consume(token.TokenSemicolon, "expected ';' after for init")
	if err != nil {
		return nil, err
	}

	var cond ast.Expr
	if !p.check(token.TokenSemicolon) {
		cond, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.TokenSemicolon, "expected ';' after for init")
	if err != nil {
		return nil, err
	}

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

	_, err = p.consume(token.TokenLBrace, "expected '{' after for clause")
	if err != nil {
		return nil, err
	}
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
	if p.check(token.TokenRBrace) || p.isAtEnd() {
		return nil, nil
	}
	fmt.Printf("[STATMENT] At: %v Current: %v (type=%v, line=%v, lexeme=%q)", p.curr.Data, p.curr.Type, p.curr.Type, p.curr.Line, p.curr.Lexeme)
	ok, err := p.match(token.TokenVar)
	if err != nil {
		return nil, err
	}
	if ok {
		return p.varDeclaration()
	}

	ok, err = p.match(token.TokenStruct)
	if err != nil {
		return nil, err
	}
	if ok {
		return p.structDeclartion()
	}

	if p.check(token.TokenIdentifier) && p.checkNext(token.TokenColonEqual) {
		return p.shortVarDeclaration()
	}

	if p.check(token.TokenIdentifier) && p.checkNext(token.TokenAssign) {
		return p.assignmentStatement()
	}
	ok, err = p.match(token.TokenPrint)
	if err != nil {
		return nil, err
	}
	if ok {
		return p.printStatement()
	}

	ok, err = p.match(token.TokenIf)
	if err != nil {
		return nil, err
	}
	if ok {
		return p.ifStatement()
	}

	ok, err = p.match(token.TokenWhile)
	if err != nil {
		return nil, err
	}
	if ok {
		return p.whileStatement()
	}

	ok, err = p.match(token.TokenFor)
	if err != nil {
		return nil, err
	}
	if ok {
		return p.forStatement()
	}
	ok, err = p.match(token.TokenFunc)
	if err != nil {
		return nil, err
	}
	if ok {
		return p.funcDeclaration()
	}

	ok, err = p.match(token.TokenReturn)
	if err != nil {
		return nil, err
	}
	if ok {
		return p.returnStatement()
	}

	ok, err = p.match(token.TokenBreak)
	if err != nil {
		return nil, err
	}
	if ok {
		return p.breakStatement()
	}
	ok, err = p.match(token.TokenContinue)
	if err != nil {
		return nil, err
	}
	if ok {
		return p.continueStatement()
	}
	ok, err = p.match(token.TokenLBrace)
	if err != nil {
		return nil, err
	}
	if ok {
		return p.expressionStatement()
	}

	return p.expressionStatement()
}
