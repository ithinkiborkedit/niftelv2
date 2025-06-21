package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	token "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"
)

// type Lexer struct {
// 	input      string
// 	start      int
// 	current    int
// 	line       int
// 	column     int
// 	width      int
// 	tokens     []niftokens.NifToken
// 	keywords   map[string]niftokens.TokenType
// 	groupdepth int
// }

type Lexer struct {
	source  string
	start   int
	current int
	line    int
	column  int
	tokens  []token.NifToken
}

func New(source string) *Lexer {
	return &Lexer{
		source:  source,
		start:   0,
		current: 0,
		line:    1,
		column:  0,
		tokens:  []token.NifToken{},
	}
}

func (l *Lexer) ScanTokens() []token.NifToken {
	for !l.isAtEnd() {
		l.start = l.current
		l.scanToken()
	}
	l.tokens = append(l.tokens, token.NifToken{
		Type:   token.TokenEOF,
		Lexeme: "",
		Line:   l.line,
		Column: l.column,
	})

	return l.tokens
}

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

func (l *Lexer) addToken(tt token.TokenType) {
	lexeme := l.source[l.start:l.current]
	l.tokens = append(l.tokens, token.NifToken{
		Type:   tt,
		Lexeme: lexeme,
		Line:   l.line,
		Column: l.column,
	})
}

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

func (l *Lexer) string(quote rune) {
	for !l.isAtEnd() {
		r, _ := utf8.DecodeRuneInString(l.source[l.current:])
		if r == quote {
			l.advance()
			lexeme := l.source[l.start:l.current]
			l.tokens = append(l.tokens, token.NifToken{
				Type:   token.TokenString,
				Lexeme: lexeme,
				Line:   l.line,
				Column: l.column,
			})
			return
		}
		if r == '\n' {
			l.line++
			l.column = 0
		}
		if r == '\\' {
			l.advance()
		}
		l.advance()
	}
	l.errorf("unterminated string literal")
}

func (l *Lexer) number() {
	for !l.isAtEnd() {
		r, _ := utf8.DecodeRuneInString(l.source[l.current:])
		if !(unicode.IsDigit(r) || r == '.') {
			break
		}
		l.advance()
	}

	lexeme := l.source[l.start:l.current]
	l.tokens = append(l.tokens, token.NifToken{
		Type:   token.TokenNumber,
		Lexeme: lexeme,
		Line:   l.line,
		Column: l.column,
	})
}

func (l *Lexer) identifier() {
	for !l.isAtEnd() {
		r, _ := utf8.DecodeRuneInString(l.source[l.current:])
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_') {
			break
		}
		l.advance()
	}
	lexeme := l.source[l.start:l.current]
	tt := lookupIdentifier(lexeme)
	l.tokens = append(l.tokens, token.NifToken{
		Type:   tt,
		Lexeme: lexeme,
		Line:   l.line,
		Column: l.column,
	})
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

func (l *Lexer) scanToken() {
	ch := l.advance()
	switch ch {
	case '(':
		l.addToken(token.TokenLParen)
	case ')':
		l.addToken(token.TokenRParen)
	case '{':
		l.addToken(token.TokenLBrace)
	case '}':
		l.addToken(token.TokenRBrace)
	case '[':
		l.addToken(token.TokenLBracket)
	case ']':
		l.addToken(token.TokenRBracket)
	case ',':
		l.addToken(token.TokenComma)
	case '.':
		l.addToken(token.TokenDot)
	case ';':
		l.addToken(token.TokenSemicolon)
	case ':':
		if l.match('=') {
			l.addToken(token.TokenColonEqual)
		}
		l.addToken(token.TokenColon)
	case '+':
		l.addToken(token.TokenPlus)
	case '-':
		if l.match('>') {
			l.addToken(token.TokenArrow)
		} else {
			l.addToken(token.TokenBang)
		}
	case '*':
		l.addToken(token.TokenStar)
	case '/':
		if l.match('/') {
			l.skipLineComment()
		} else if l.match('*') {
			l.skipBlockComment()
		} else {
			l.addToken(token.TokenFWDSlash)
		}
	case '%':
		l.addToken(token.TokenPercent)
	case '=':
		if l.match('=') {
			l.addToken(token.TokenEqality)
		} else {
			l.addToken(token.TokenAssign)
		}
	case '!':
		if l.match('=') {
			l.addToken(token.TokenBangEqal)
		} else {
			l.addToken(token.TokenBang)
		}
	case '<':
		if l.match('=') {
			l.addToken(token.TokenLessEq)
		} else {
			l.addToken(token.TokenLess)
		}
	case '>':
		if l.match('=') {
			l.addToken(token.TokenGreaterEq)
		} else {
			l.addToken(token.TokenGreater)
		}
	case '&':
		if l.match('&') {
			l.addToken(token.TokenAnd)
		} else {
			l.errorf("unexpected character: %q", ch)
		}
	case '|':
		if l.match('|') {
			l.addToken(token.TokenOr)
		} else {
			l.errorf("unexpected character: %q", ch)
		}
	case '"', '\'':
		l.string(ch)
	case '\n':
		l.addToken(token.TokenNewLine)
		l.line++
		l.column = 0
	case ' ', '\r', '\t':

	default:
		if unicode.IsDigit(ch) {
			l.number()
		} else if unicode.IsLetter(ch) || ch == '_' {
			l.identifier()
		} else {
			l.errorf("unexpected character: '%q'", ch)
		}
	}
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
