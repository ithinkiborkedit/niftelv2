package runtimecontrol

import (
	"fmt"

	"github.com/ithinkiborkedit/niftelv2.git/internal/value"
)

type ReturnValue struct {
	Value value.Value
}

type BreakSignal struct{}
type ContinueSignal struct{}

func (b BreakSignal) Error() string { return "break" }

func (c ContinueSignal) Error() string { return "continue" }

func (r ReturnValue) Error() string { return fmt.Sprintf("return: '%v'", r.Value) }
