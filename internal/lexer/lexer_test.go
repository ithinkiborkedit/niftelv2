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

		// if tok.Type == 0 {
		// 	continue
		// }

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
