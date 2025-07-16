package tokentoval

import (
	"fmt"
	"strconv"

	tokens "github.com/ithinkiborkedit/niftelv2.git/internal/niftokens"
	"github.com/ithinkiborkedit/niftelv2.git/internal/value"
)

func Convert(tok tokens.Token) (value.Value, error) {
	switch tok.Type {
	case tokens.TokenNumber:
		fval, err := strconv.ParseFloat(tok.Lexeme, 64)
		if err != nil {
			return value.Null(), fmt.Errorf("invalid numbner token: %v", err)
		}
		return value.Value{
			Type: value.ValueFloat,
			Data: fval,
		}, nil
	case tokens.TokenFloat:
		fval, err := strconv.ParseFloat(tok.Lexeme, 64)
		if err != nil {
			return value.Null(), fmt.Errorf("invalid float token %v", err)
		}
		return value.Value{
			Type: value.ValueFloat,
			Data: fval,
		}, nil
	case tokens.TokenString:
		str := tok.Lexeme
		if len(str) >= 2 && (str[0] == '"' || str[0] == '\'') && str[len(str)-1] == str[0] {
			str = str[1 : len(str)-1]
		}
		return value.Value{
			Type: value.ValueString,
			Data: str,
		}, nil
	case tokens.TokenTrue, tokens.TokenFalse:
		return value.Value{
			Type: value.ValueBool,
			Data: tok.Type == tokens.TokenTrue,
		}, nil
	case tokens.TokenNull:
		return value.Null(), nil
	default:
		return value.Null(), fmt.Errorf("unsported token type for conversion: %v", tok.Type)
	}
}
