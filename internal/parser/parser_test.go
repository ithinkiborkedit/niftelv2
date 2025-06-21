package parser_test

import (
	"testing"

	"github.com/ithinkiborkedit/niftelv2.git/internal/lexer"
	"github.com/ithinkiborkedit/niftelv2.git/internal/parser"
)

func TestParser_Basic(t *testing.T) {
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

	l := lexer.New(source)
	tokens := l.ScanTokens()

	for i, tok := range tokens {
		t.Logf("%30d: %s\t%q", i, tok.Type, tok.Lexeme)
	}

	p := parser.New(tokens)

	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	t.Logf("AST %+v", program)

}
