package lexer

import (
	"testing"

	token "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"
)

func TestLexer_ScanTokens(t *testing.T) {
	source := `
	func add(a: int, b: int) -> int {
	   return a + b
	}
	var result: int = add(5,3)

	if result > 5 {
	  print("success")

	} else {
	  print("fail") 
	}
	// single line comment
	/* block
	   comment */
	`
	lex := New(source)
	tokens := lex.ScanTokens()

	var filtered []token.NifToken

	for _, tok := range tokens {

		if tok.Type == token.TokenNewLine {
			continue
		}

		filtered = append(filtered, tok)

	}

	expectedTypes := []token.TokenType{
		token.TokenIllegal,
		token.TokenEOF,
		token.TokenIdentifier,
		token.TokenNumber,
		token.TokenString,
		token.TokenBool,
		token.TokenNull,

		token.TokenAssign,
		token.TokenBang,
		token.TokenComma,
		token.TokenSemicolon,
		token.TokenStar,
		token.TokenFWDSlash,
		token.TokenPlus,
		token.TokenPlusEq,
		token.TokenMinEq,
		token.TokenMinus,
		token.TokenLParen,
		token.TokenRParen,
		token.TokenLBrace,
		token.TokenRBrace,
		token.TokenPipe,
		token.TokenLBracket,
		token.TokenRBracket,
		token.TokenRepo,
		token.TokenBranch,
		token.TokenEqality,
		token.TokenBangEqal,
		token.TokenGreater,
		token.TokenLess,
		token.TokenGreaterEq,
		token.TokenLessEq,
		token.TokenAnd,
		token.TokenOr,
		token.TokenColon,
		token.TokenDot,
		token.TokenPercent,
		token.TokenAmper,
		token.TokenNewLine,
		token.TokenArrow,

		//Keywords
		token.TokenTrue,
		token.TokenT,
		token.TokenStruct,
		token.TokenImport,
		token.TokenAs,
		token.TokenNil,
		token.TokenFalse,
		token.TokenIf,
		token.TokenElse,
		token.TokenFor,
		token.TokenIn,
		token.TokenVar,
		token.TokenFunc,
		token.TokenReturn,
		token.TokenWhile,
		token.TokenPrint,
		token.TokenBreak,
		token.TokenContinue,
	}
	for i, tok := range filtered {
		t.Logf("%03d: %-20v %q", i, tok.Type, tok.Lexeme)
	}
	if len(tokens) != len(expectedTypes) {
		t.Fatalf("token count mismatch: got %d, want %d", len(tokens), len(expectedTypes))
	}

	for i, tok := range tokens {

		for tok.Type != expectedTypes[i] {
			t.Errorf("token %d: got %v, want %v (lexeme: %q)", i, tok.Type, expectedTypes[i], tok.Lexeme)
		}

	}
}
