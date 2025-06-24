package lexer

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"

	token "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"
)

type Lexer struct {
	source  string
	start   int
	current int
	line    int
	column  int
}

type TokenSource interface {
	NextToken() (token.Token, error)
}

var tokenKeyWords = map[string]token.TokenType{
	"true":     token.TokenTrue,
	"type":     token.TokenT,
	"struct":   token.TokenStruct,
	"import":   token.TokenImport,
	"as":       token.TokenAs,
	"nil":      token.TokenNil,
	"false":    token.TokenFalse,
	"if":       token.TokenIf,
	"else":     token.TokenElse,
	"for":      token.TokenFor,
	"in":       token.TokenIn,
	"var":      token.TokenVar,
	"func":     token.TokenFunc,
	"return":   token.TokenReturn,
	"while":    token.TokenWhile,
	"print":    token.TokenPrint,
	"break":    token.TokenBreak,
	"continue": token.TokenContinue,
}

func (l *Lexer) NextToken() (token.Token, error) {
	for !l.isAtEnd() {
		l.start = l.current
		tok, err := l.scanToken()
		if err != nil {
			return token.Token{}, err
		}
		// if tok.Type == token.TokenIllegal {
		// 	return tok, nil
		// }
		if tok.Type == 0 {
			continue
		}
		return token.Token{
			Type:   token.TokenEOF,
			Lexeme: "",
			Line:   l.line,
			Column: l.column,
		}, nil
	}
	// tok, matched := l.scanToken()
	// if matched {
	// 	return tok, nil
	// }
	return l.scanToken()
}

func (l *Lexer) makeToken(tt token.TokenType) token.Token {
	lexeme := l.source[l.start:l.current]
	return token.Token{
		Type:   tt,
		Lexeme: lexeme,
		Line:   l.line,
		Column: l.column,
	}
}

func New(source string) *Lexer {
	return &Lexer{
		source:  source,
		start:   0,
		current: 0,
		line:    1,
		column:  0,
		// tokens:  []token.Token{},
	}
}

// func (l *Lexer) ScanTokens() []token.Token {
// 	for !l.isAtEnd() {
// 		l.start = l.current
// 		l.scanToken()
// 	}
// 	l.tokens = append(l.tokens, token.Token{
// 		Type:   token.TokenEOF,
// 		Lexeme: "",
// 		Line:   l.line,
// 		Column: l.column,
// 	})

// 	return l.tokens
// }

func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

func (l *Lexer) advance() rune {
	if l.isAtEnd() {

		return 0

	}

	r, width := utf8.DecodeRuneInString(l.source[l.current:])
	l.current += width
	l.column++
	return r
}

// func (l *Lexer) addToken(tt token.TokenType) {
// 	lexeme := l.source[l.start:l.current]
// 	l.tokens = append(l.tokens, token.Token{
// 		Type:   tt,
// 		Lexeme: lexeme,
// 		Line:   l.line,
// 		Column: l.column,
// 	})
// }

func (l *Lexer) match(expected rune) bool {
	if l.isAtEnd() {
		return false
	}
	r, width := utf8.DecodeRuneInString(l.source[l.current:])
	if r != expected {
		return false
	}
	l.current += width
	l.column++
	return true
}

func (l *Lexer) skipLineComment() {
	for !l.isAtEnd() {
		r, _ := utf8.DecodeRuneInString(l.source[l.current:])
		if r == '\n' {
			break
		}
		l.advance()
	}
}

func (l *Lexer) skipBlockComment() {
	for !l.isAtEnd() {
		r, _ := utf8.DecodeRuneInString(l.source[l.current:])
		if r == '*' && l.peekNext() == '/' {
			l.advance()
			l.advance()
			break
		}
		if r == '\n' {
			l.line++
			l.column = 0
		}
		l.advance()
	}
}

func (l *Lexer) errorf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	panic(fmt.Sprintf("Lexer error at line %d, col %d: %s", l.line, l.column, msg))
}

func (l *Lexer) peekNext() rune {
	if l.current >= len(l.source) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.source[l.current:])
	return r
}

func (l *Lexer) string(quote rune) (token.Token, error) {
	for !l.isAtEnd() {
		r, _ := utf8.DecodeRuneInString(l.source[l.current:])
		if r == quote {
			l.advance()
			lexeme := l.source[l.start:l.current]
			return token.Token{
				Type:   token.TokenString,
				Lexeme: lexeme,
				Line:   l.line,
				Column: l.column,
			}, nil
		}
		if r == '\n' {
			l.line++
			l.column = 0
		}

		if r == '\\' {
			l.advance()
		}
		l.advance()
		// r, _ := utf8.DecodeRuneInString(l.source[l.current:])
		// if r == quote {
		// 	l.advance()
		// 	lexeme := l.source[l.start:l.current]
		// 	l.tokens = append(l.tokens, token.Token{
		// 		Type:   token.TokenString,
		// 		Lexeme: lexeme,
		// 		Line:   l.line,
		// 		Column: l.column,
		// 	})
		// 	return
		// }
		// if r == '\n' {
		// 	l.line++
		// 	l.column = 0
		// }
		// if r == '\\' {
		// 	l.advance()
		// }
		// l.advance()
	}
	return token.Token{}, fmt.Errorf("unterminated string literal")
}

func (l *Lexer) number() (token.Token, error) {
	for !l.isAtEnd() {
		r, _ := utf8.DecodeRuneInString(l.source[l.current:])
		if !unicode.IsDigit(r) && r != '.' {
			break
		}
		l.advance()
	}
	lexeme := l.source[l.start:l.current]

	if _, err := strconv.ParseFloat(lexeme, 64); err != nil {
		return token.Token{}, fmt.Errorf("invalid number literal: %s", lexeme)
	}
	// 	r, _ := utf8.DecodeRuneInString(l.source[l.current:])
	// 	if !(unicode.IsDigit(r) || r == '.') {
	// 		break
	// 	}
	// 	l.advance()
	// }

	// lexeme := l.source[l.start:l.current]
	// val, err := strconv.ParseFloat(lexeme, 64)
	// if err != nil {
	// 	l.errorf("invalid number literal: %s", lexeme)
	// 	val = 0.0
	// }

	return token.Token{
		Type:   token.TokenNumber,
		Lexeme: lexeme,
		// Data:   val,
		Line:   l.line,
		Column: l.column,
	}, nil
}

func (l *Lexer) identifier() (token.Token, error) {
	for !l.isAtEnd() {
		r, _ := utf8.DecodeRuneInString(l.source[l.current:])
		if !unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			break
		}
		l.advance()
	}
	lexeme := l.source[l.start:l.current]

	tokType := token.TokenIdentifier
	if kwType, ok := tokenKeyWords[lexeme]; ok {
		tokType = kwType
	}

	return token.Token{
		Type:   tokType,
		Lexeme: lexeme,
		Line:   l.line,
		Column: l.column,
	}, nil
	// 	r, _ := utf8.DecodeRuneInString(l.source[l.current:])
	// 	if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
	// 		break
	// 	}
	// 	l.advance()
	// }
	// lexeme := l.source[l.start:l.current]
	// tt := lookupIdentifier(lexeme)
	// l.tokens = append(l.tokens, token.Token{
	// 	Type:   tt,
	// 	Lexeme: lexeme,
	// 	Line:   l.line,
	// 	Column: l.column,
	// })
}

func lookupIdentifier(lexeme string) token.TokenType {
	keywords := map[string]token.TokenType{

		"true":   token.TokenTrue,
		"type":   token.TokenT,
		"struct": token.TokenStruct,
		"import": token.TokenImport,
		"as":     token.TokenAs,
		"nil":    token.TokenNil,
		"false":  token.TokenFalse,
		"if":     token.TokenIf,
		"else":   token.TokenElse,
		"for":    token.TokenFor,
		"in":     token.TokenIn,
		"var":    token.TokenVar,
		"func":   token.TokenFunc,
		"return": token.TokenReturn,
		"while":  token.TokenWhile,
	}

	if tt, ok := keywords[lexeme]; ok {
		return tt
	}
	return token.TokenIdentifier
}

func (l *Lexer) scanToken() (token.Token, error) {
	ch := l.advance()
	switch ch {
	case '(':
		l.makeToken(token.TokenLParen)
	case ')':
		l.makeToken(token.TokenRParen)
	case '{':
		l.makeToken(token.TokenLBrace)
	case '}':
		l.makeToken(token.TokenRBrace)
	case '[':
		l.makeToken(token.TokenLBracket)
	case ']':
		l.makeToken(token.TokenRBracket)
	case ',':
		l.makeToken(token.TokenComma)
	case '.':
		l.makeToken(token.TokenDot)
	case ';':
		l.makeToken(token.TokenSemicolon)
	case ':':
		if l.match('=') {
			l.makeToken(token.TokenColonEqual)
		} else {
			l.makeToken(token.TokenColon)
		}
	case '+':
		l.makeToken(token.TokenPlus)
	case '-':
		if l.match('>') {
			l.makeToken(token.TokenArrow)
		} else {
			l.makeToken(token.TokenBang)
		}
	case '*':
		l.makeToken(token.TokenStar)
	case '/':
		if l.match('/') {
			l.skipLineComment()
		} else if l.match('*') {
			l.skipBlockComment()
		} else {
			l.makeToken(token.TokenFWDSlash)
		}
	case '%':
		l.makeToken(token.TokenPercent)
	case '=':
		if l.match('=') {
			l.makeToken(token.TokenEqality)
		} else {
			l.makeToken(token.TokenAssign)
		}
	case '!':
		if l.match('=') {
			l.makeToken(token.TokenBangEqal)
		} else {
			l.makeToken(token.TokenBang)
		}
	case '<':
		if l.match('=') {
			l.makeToken(token.TokenLessEq)
		} else {
			l.makeToken(token.TokenLess)
		}
	case '>':
		if l.match('=') {
			l.makeToken(token.TokenGreaterEq)
		} else {
			l.makeToken(token.TokenGreater)
		}
	case '&':
		if l.match('&') {
			l.makeToken(token.TokenAnd)
		} else {
			l.errorf("unexpected character: %q", ch)
		}
	case '|':
		if l.match('|') {
			l.makeToken(token.TokenOr)
		} else {
			l.errorf("unexpected character: %q", ch)
		}
	case '"', '\'':
		l.string(ch)
	case '\n':
		tok := l.makeToken(token.TokenNewLine)
		l.line++
		l.column = 0
		return tok, nil
	case ' ', '\r', '\t':
		return token.Token{}, nil
	default:
		if unicode.IsDigit(ch) {
			l.number()
		} else if unicode.IsLetter(ch) || ch == '_' {
			l.identifier()
		} else {
			l.errorf("unexpected character: '%q'", ch)
		}
	}
	return token.Token{}, fmt.Errorf("unexpected character %q", ch)
}

// func (l *Lexer) ScanTokens() []niftokens.NifToken {
// 	for !l.isAtEnd() {
// 		l.start = l.current
// 		l.ScanToken()
// 	}
// 	l.tokens = append(l.tokens, niftokens.NifToken{
// 		Type:niftokens.TokenEOF,
// 		Line:l.line,
// 		Column: l.column,
// 	})

// 	return l.tokens
// }

// func (l *Lexer) ScanToken() {
// 	ch := l.advance()
// 	switch ch {
// 	case '+':
// 		l.addToken(niftokens.TokenPlus)
// 	case '-':
// 		l.addToken(niftokens.TokenMinus)
// 	case '*':
// 		l.addToken(niftokens.TokenStar)
// 	case '/':
// 		if l.match('/') {
// 			l.skipLineComment()
// 		} else {
// 			l.addToken(niftokens.TokenFWDSlash)
// 		}
// 	case '%':
// 		l.addToken(niftokens.TokenPercent)
// 	case '=':
// 		if l.match('=') {
// 			l.addToken(niftokens.TokenEqality)
// 		} else {
// 			l.addToken(niftokens.TokenEqual)
// 		}
// 	case '!':
// 		if l.match('=') {
// 			l.addToken(niftokens.TokenBangEqal)
// 		} else {
// 			l.addToken(niftokens.TokenBang)
// 		}
// 	case '<':
// 		if l.match('=') {
// 			l.addToken(niftokens.TokenLessEq)
// 		} else {
// 			l.addToken(niftokens.TokenLess)
// 		}
// 	case '>':
// 		if l.match('=') {
// 			l.addToken(niftokens.TokenGreaterEq)
// 		} else {
// 			l.addToken(niftokens.TokenGreater)
// 		}
// 	case '&':
// 		if l.match('&') {
// 			l.addToken(niftokens.TokenAnd)
// 		} else {
// 			l.addToken(niftokens.TokenAmper)
// 		}
// 	case '|':
// 		if l.match('|') {
// 			l.addToken(niftokens.TokenOr)
// 		} else {
// 			l.errorf("unexpected character '%q' ", ch)
// 		}
// 	case ',':
// 		l.addToken(niftokens.TokenComma)
// 	case ':':
// 		l.addToken(niftokens.TokenColon)
// 	case '.':
// 		l.addToken(niftokens.TokenDot)
// 	case '(':
// 		l.groupdepth++
// 		l.addToken(niftokens.TokenLParen)
// 	case ')':
// 		if l.groupdepth > 0 {
// 			l.groupdepth--
// 		}
// 		l.addToken(niftokens.TokenRParen)
// 	case '{':
// 		l.groupdepth++
// 		l.addToken(niftokens.TokenLBrace)
// 	case '}':
// 		if l.groupdepth > 0 {
// 			l.groupdepth--
// 		}
// 		l.addToken(niftokens.TokenRBrace)
// 	case '[':
// 		l.groupdepth++
// 		l.addToken(niftokens.TokenLBracket)
// 	case ']':
// 		if l.groupdepth > 0 {
// 			l.groupdepth--
// 		}
// 		l.addToken(niftokens.TokenRBracket)
// 	case '"':
// 		l.string()
// 	case '\'':
// 		l.string()
// 	case ' ', '\r', '\t':

// 	default:
// 		if unicode.IsDigit(ch) {
// 			l.number()
// 		} else if unicode.IsLetter(ch) || ch == '-' {
// 			l.identifier()
// 		} else {
// 			l.errorf("unexpected character: '%q'", ch)
// 		}
// 	}
// }

// func (l *Lexer) addToken(tt niftokens.TokenType) {
// 	lexeme := l.input[l.start:l.current]
// 	tok := niftokens.NifToken{
// 		Type:   tt,
// 		Lexeme: lexeme,
// 		Line:   l.line,
// 		Column: l.column,
// 	}

// 	l.tokens = append(l.tokens, tok)
// }

// func (l *Lexer) addTokenLiteral(tt niftokens.TokenType, literal interface{}) {
// 	lexeme := l.input[l.start:l.current]
// 	tok := niftokens.NifToken{
// 		Type:    tt,
// 		Lexeme:  lexeme,
// 		Literal: literal,
// 		Line:    l.line,
// 		Column:  l.column,
// 	}

// 	l.tokens = append(l.tokens, tok)
// }

// func (l *Lexer) isAtEnd() bool {
// 	return l.current >= len(l.input)
// }

// func (l *Lexer) advance() rune {
// 	if l.isAtEnd() {
// 		l.width = 0
// 		return 0
// 	}

// 	ch, w := utf8.DecodeRuneInString(l.input[l.current:])
// 	l.current += w
// 	l.column += w
// 	l.width = w
// 	return ch
// }

// func (l *Lexer) match(expected rune) bool {
// 	if l.isAtEnd() {
// 		return false
// 	}
// 	ch, w := utf8.DecodeRuneInString(l.input[l.current:])
// 	if ch != expected {
// 		return false
// 	}

// 	l.current += w
// 	l.column += w
// 	return true
// }
// func (l *Lexer) number() {
// 	for unicode.IsDigit(l.peek()) {
// 		l.advance()

// 	}
// 	if l.peek() == '.' && unicode.IsDigit(l.peekNext()) {
// 		l.advance()
// 		for unicode.IsDigit(l.peek()) {
// 			l.advance()
// 		}
// 	}
// 	val := l.input[l.start:l.current]
// 	var number interface{}
// 	var err error
// 	if containsDot(val) {
// 		number, err = parseFloat(val)
// 	} else {
// 		number, err = parseInt(val)
// 	}
// 	if err != nil {
// 		l.errorf("invalid number %s", val)
// 		return
// 	}
// 	l.addTokenLiteral(niftokens.TokenNumber, number)
// }
// func (l *Lexer) peek() rune {
// 	if l.isAtEnd() {
// 		return 0
// 	}

// 	ch, _ := utf8.DecodeRuneInString(l.input[l.current:])
// 	return ch
// }

// func (l *Lexer) peekNext() rune {
// 	if l.current+1 >= len(l.input) {
// 		return 0
// 	}

// 	ch, _ := utf8.DecodeRuneInString(l.input[l.current+1:])
// 	return ch
// }

// func (l *Lexer) skipLineComment() {
// 	for {
// 		if l.isAtEnd() {
// 			return
// 		}
// 		ch, _ := utf8.DecodeRuneInString(l.input[l.current:])
// 		if ch == '\n' {
// 			return
// 		}
// 		l.advance()
// 	}
// }

// func (l *Lexer) identifier() {
// 	for unicode.IsLetter(l.peek()) || unicode.IsDigit(l.peek()) || l.peek() == '_' {
// 		l.advance()
// 	}
// 	lexeme := l.input[l.start:l.current]
// 	if tt, ok := l.keywords[lexeme]; ok {
// 		l.addToken(tt)
// 	} else {
// 		l.addToken(niftokens.TokenIdent)
// 	}
// }

// func (l *Lexer) string() {
// 	quote := l.input[l.start]
// 	l.advance()
// 	for {
// 		if l.isAtEnd() {
// 			l.errorf("unterminated string")
// 			return
// 		}
// 		ch := l.peek()
// 		if ch == rune(quote) {

// 			break

// 		}
// 		if ch == '\n' {
// 			l.line++
// 			l.column = 0
// 		}
// 		l.advance()
// 	}
// 	l.advance()
// 	val := l.input[l.start+1 : l.current-1]
// 	l.addTokenLiteral(niftokens.TokenString, val)
// }

// func (l *Lexer) errorf(format string, args ...interface{}) {
// 	msg := fmt.Sprintf("Line %d, Col %d: ", l.line, l.column) + fmt.Sprintf(format, args...)
// 	panic(msg)
// }

// func containsDot(s string) bool {
// 	for _, c := range s {
// 		if c == '.' {
// 			return true
// 		}
// 	}

// 	return false
// }

// func parseFloat(s string) (float64, error) {
// 	var v float64
// 	_, err := fmt.Sscanf(s, "%f", &v)
// 	return v, err
// }

// func parseInt(s string) (int64, error) {
// 	var v int64
// 	_, err := fmt.Sscanf(s, "%d", &v)
// 	return v, err
// }
