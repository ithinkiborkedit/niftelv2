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

	for i, tok := range tokens {
		fmt.Printf("%d: %v %q", i, tok.Type, tok.Lexeme)
	}
	for i, tok := range tokens {
		if tok.Type == token.TokenIllegal {
			t.Fatalf("%d: %v %q", i, tok.Type, tok.Lexeme)
		}
	}

}

// for i, tok := range filtered {
// 	t.Logf("%03d: %-20v %q", i, tok.Type, tok.Lexeme)
// }
// if len(filtered) != len(expectedTypes) {
// 	t.Fatalf("token count mismatch: got %d, want %d", len(filtered), len(expectedTypes))
// }

// for i, tok := range filtered {

// 	for tok.Type != expectedTypes[i] {
// 		t.Errorf("token %d: got %v, want %v (lexeme: %q)", i, tok.Type, expectedTypes[i], tok.Lexeme)
// 	}

// }
