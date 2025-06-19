package lexer

import "go/token"

type Lexer struct {
	input   string
	start   int
	current int
	line    int
	column  int
	width   int
	tokens  []token.Token
}
