package lexer

import (
	"fmt"
	"testing"

	token "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"
)

func TestLexer_ScanTokens(t *testing.T) {
	source := `
	func add(a: int, b: int) -> int {
	   return a + b
	}
	var result: int = add(5,3)

	if result 5 {
	  print("success")

	} else {
	  print("fail") 
	}

	struct Person {
	  a: int
	  b: int
	}

	// single line comment
	/* block
	   comment */
	`
	lex := New(source)
	var tokens []token.Token
	for {
		tok, err := lex.NextToken()
		if err != nil {
			t.Fatalf("lexer error %v", err)
		}

		fmt.Printf("%d: %v %q \n", len(tokens), tok.Type, tok.Lexeme)

		if tok.Type == token.TokenIllegal {
			t.Fatalf("%d: ILLEGAL TOKEN %q", len(tokens), tok.Lexeme)
		}

		tokens = append(tokens, tok)

		if tok.Type == token.TokenEOF {
			break
		}
	}

}

func TestLexer_StructTokens(t *testing.T) {
	source := `
struct Person {
    a: int
    b: int
}
`
	lex := New(source)
	for {
		tok, err := lex.NextToken()
		if err != nil {
			t.Fatalf("lexer error %v", err)
		}
		fmt.Printf("%v (%q) line=%d col=%d\n", tok.Type, tok.Lexeme, tok.Line, tok.Column)
		if tok.Type == token.TokenEOF {
			break
		}
	}
}
