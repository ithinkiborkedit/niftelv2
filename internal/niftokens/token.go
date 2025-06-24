package niftokens

type TokenType int

const (
	TokenIllegal TokenType = iota
	TokenEOF
	TokenIdentifier
	TokenNumber
	TokenString
	TokenBool
	TokenNull

	TokenAssign
	TokenBang
	TokenComma
	TokenSemicolon
	TokenStar
	TokenFWDSlash
	TokenPlus
	TokenPlusEq
	TokenMinEq
	TokenMinus
	TokenLParen
	TokenRParen
	TokenLBrace
	TokenRBrace
	TokenPipe
	TokenLBracket
	TokenRBracket
	TokenRepo
	TokenBranch
	TokenEqality
	TokenBangEqal
	TokenGreater
	TokenLess
	TokenGreaterEq
	TokenLessEq
	TokenAnd
	TokenOr
	TokenColon
	TokenDot
	TokenPercent
	TokenAmper
	TokenNewLine
	TokenArrow
	TokenColonEqual

	//Keywords
	TokenTrue
	TokenT
	TokenStruct
	TokenImport
	TokenAs
	TokenNil
	TokenFalse
	TokenIf
	TokenElse
	TokenFor
	TokenIn
	TokenVar
	TokenFunc
	TokenReturn
	TokenWhile
	TokenPrint
	TokenBreak
	TokenContinue
)

var tokenTypeToString = map[TokenType]string{
	TokenIllegal:    "ILLEGAL",
	TokenEOF:        "EOF",
	TokenIdentifier: "IDENTIFIER",
	TokenNumber:     "NUMBER",
	TokenString:     "STRING",
	TokenAssign:     "=",
	TokenBang:       "!",
	TokenComma:      ",",
	TokenStar:       "*",
	TokenFWDSlash:   "/",
	TokenPercent:    "%",
	TokenPlus:       "+",
	TokenPlusEq:     "+=",
	TokenMinus:      "-",
	TokenArrow:      "->",
	TokenMinEq:      "-=",
	TokenLParen:     "(",
	TokenRParen:     ")",
	TokenLBrace:     "{",
	TokenRBrace:     "}",
	TokenLBracket:   "[",
	TokenRBracket:   "]",
	TokenEqality:    "==",
	TokenColonEqual: ":=",
	TokenBangEqal:   "!=",
	TokenGreater:    ">",
	TokenLess:       "<",
	TokenGreaterEq:  ">=",
	TokenLessEq:     "<=",
	TokenAmper:      "&",
	TokenAnd:        "&&",
	TokenPipe:       "|",
	TokenOr:         "||",
	TokenColon:      ":",
	TokenSemicolon:  ";",
	TokenDot:        ".",
	TokenTrue:       "true",
	TokenT:          "type",
	TokenStruct:     "struct",
	TokenImport:     "import",
	TokenAs:         "as",
	TokenNil:        "nil",
	TokenFalse:      "false",
	TokenIf:         "if",
	TokenElse:       "else",
	TokenFor:        "for",
	TokenIn:         "in",
	TokenVar:        "var",
	TokenFunc:       "func",
	TokenReturn:     "return",
	TokenWhile:      "while",
	TokenPrint:      "print",
	TokenBreak:      "break",
	TokenContinue:   "continue",
	TokenNewLine:    "\n",
}

var tokenKeyWords = map[string]TokenType{
	"true":     TokenTrue,
	"type":     TokenT,
	"struct":   TokenStruct,
	"import":   TokenImport,
	"as":       TokenAs,
	"nil":      TokenNil,
	"false":    TokenFalse,
	"if":       TokenIf,
	"else":     TokenElse,
	"for":      TokenFor,
	"in":       TokenIn,
	"var":      TokenVar,
	"func":     TokenFunc,
	"return":   TokenReturn,
	"while":    TokenWhile,
	"print":    TokenPrint,
	"break":    TokenBreak,
	"continue": TokenContinue,
}

func (tt TokenType) String() string {
	if s, ok := tokenTypeToString[tt]; ok {
		return s
	}
	return "UNKNOWN"
}

type Token struct {
	Type   TokenType
	Lexeme string
	Data   interface{}
	Line   int
	Column int
}
