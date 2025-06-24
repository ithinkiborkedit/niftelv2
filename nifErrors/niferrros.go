package niferrors

import "fmt"

type ErrorKind int

const (
	LexError ErrorKind = iota
	ParseError
	RuntimeError
	IOError
)

type NifError struct {
	Kind    ErrorKind
	Message string
	Line    int
	Column  int
	Token   string
	File    string
}

func (e NifError) Error() string {
	return fmt.Sprintf("[%v] %s(line %d, col %d): %s", e.Kind, e.File, e.Line, e.Column, e.Message)
}
