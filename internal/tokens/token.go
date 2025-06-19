package tokens

type TokenType int

const (
	TokenIllegal TokenType = iota
	TokenEOF
	TokenIdent
	TokenNumber
	TokenString
	TokenEqual
	TokenBang
	TokenComma
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
)

var tokenTypeToString = map[TokenType]string{
	TokenIllegal:   "ILLEGAL",
	TokenEOF:       "EOF",
	TokenIdent:     "IDENTIFIER",
	TokenNumber:    "NUMBER",
	TokenString:    "STRING",
	TokenEqual:     "=",
	TokenBang:      "!",
	TokenComma:     ",",
	TokenStar:      "*",
	TokenFWDSlash:  "/",
	TokenPlus:      "+",
	TokenPlusEq:    "+=",
	TokenMinus:     "-",
	TokenMinEq:     "-=",
	TokenLParen:    "(",
	TokenRParen:    ")",
	TokenLBrace:    "{",
	TokenRBrace:    "}",
	TokenLBracket:  "[",
	TokenRBracket:  "]",
	TokenEqality:   "==",
	TokenBangEqal:  "!=",
	TokenGreater:   ">",
	TokenLess:      "<",
	TokenGreaterEq: ">=",
	TokenLessEq:    "<=",
	TokenAnd:       "&&",
	TokenOr:        "||",
	TokenColon:     ":",
	TokenDot:       ".",
	TokenTrue:      "true",
	TokenT:         "type",
	TokenStruct:    "struct",
	TokenImport:    "import",
	TokenAs:        "as",
	TokenNil:       "nil",
	TokenFalse:     "false",
	TokenIf:        "if",
	TokenElse:      "else",
	TokenFor:       "for",
	TokenIn:        "in",
	TokenVar:       "var",
	TokenFunc:      "func",
	TokenReturn:    "return",
	TokenWhile:     "while",
	TokenPrint:     "print",
}

func (tt TokenType) String() string {
	if s, ok := tokenTypeToString[tt]; ok {
		return s
	}
	return "UNKNOWN"
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
	Column  int
}
