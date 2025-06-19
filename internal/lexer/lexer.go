package lexer

import (
	"unicode"

	"github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"
)

type Lexer struct {
	input    string
	start    int
	current  int
	line     int
	column   int
	width    int
	tokens   []niftokens.NifToken
	keywords map[string]niftokens.TokenType
}

func New(input string) *Lexer {
	return &Lexer{
		input:  input,
		line:   1,
		column: 1,
		keywords: map[string]niftokens.TokenType{
			"if":     niftokens.TokenIf,
			"else":   niftokens.TokenElse,
			"while":  niftokens.TokenWhile,
			"for":    niftokens.TokenFor,
			"func":   niftokens.TokenFunc,
			"return": niftokens.TokenReturn,
			"var":    niftokens.TokenVar,
			"true":   niftokens.TokenTrue,
			"false":  niftokens.TokenFalse,
			"nil":    niftokens.TokenNil,
			"print":  niftokens.TokenPrint,
		},
	}
}

func (l *Lexer) ScanTokens() []niftokens.NifToken {
	for !l.isAtEnd() {
		l.start = l.current
		l.ScanToken()
	}
	l.tokens = append(l.tokens, niftokens.NifToken{
		Type:   niftokens.TokenEOF,
		Line:   l.line,
		Column: l.column,
	})

	return l.tokens
}

func (l *Lexer) ScanToken() {
	ch := l.advance()
	switch ch {
	case '+':
		l.addToken(niftokens.TokenPlus)
	case '-':
		l.addToken(niftokens.TokenMinus)
	case '*':
		l.addToken(niftokens.TokenStar)
	case '/':
		if l.match('/') {
			l.skipLineComment()
		} else {
			l.addToken(niftokens.TokenFWDSlash)
		}
	case '%':
		l.addToken(niftokens.TokenPercent)
	case '=':
		if l.match('=') {
			l.addToken(niftokens.TokenEqality)
		} else {
			l.addToken(niftokens.TokenEqual)
		}
	case '!':
		if l.match('=') {
			l.addToken(niftokens.TokenBangEqal)
		} else {
			l.addToken(niftokens.TokenBang)
		}
	case '<':
		if l.match('=') {
			l.addToken(niftokens.TokenLessEq)
		} else {
			l.addToken(niftokens.TokenLess)
		}
	case '>':
		if l.match('=') {
			l.addToken(niftokens.TokenGreaterEq)
		} else {
			l.addToken(niftokens.TokenGreater)
		}
	case '&':
		if l.match('&') {
			l.addToken(niftokens.TokenAnd)
		} else {
			l.addToken(niftokens.TokenAmper)
		}
	case '|':
		if l.match('|') {
			l.addToken(niftokens.TokenOr)
		} else {
			l.errorf("unexpected character '%q' ", ch)
		}
	case ',':
		l.addToken(niftokens.TokenComma)
	case ':':
		l.addToken(niftokens.TokenColon)
	case '.':
		l.addToken(niftokens.TokenDot)
	case '(':
		l.addToken(niftokens.TokenLParen)
	case ')':
		l.addToken(niftokens.TokenRParen)
	case '{':
		l.addToken(niftokens.TokenLBrace)
	case '}':
		l.addToken(niftokens.TokenRBrace)
	case '[':
		l.addToken(niftokens.TokenLBracket)
	case ']':
		l.addToken(niftokens.TokenRBracket)
	case '"':
		l.string()
	case '\'':
		l.string()
	case ' ', '\r', '\t':

	case '\n':
		l.line++
		l.column = 0
	default:
		if unicode.IsDigit(ch) {
			l.number()
		} else if unicode.IsLetter(ch) || ch == '-' {
			l.identifier()
		} else {
			l.errorf("unexpected character: '%q'", ch)
		}
	}
}

func (l *Lexer) addToken(tt niftokens.TokenType) {}
