package lexer

import ntoken "go/github.com/ithinkiborkedit/niftelv2/internal/token"

type Lexer struct {
	input   string
	start   int
	current int
	line    int
	column  int
	width   int
	tokens  []ntoken.Token
}
