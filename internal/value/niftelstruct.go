package value

import token "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"

type StructType struct {
	Name   string
	Fields []token.Token
}

type StructInstance struct {
	Type   *StructType
	Fields map[string]Value
}
