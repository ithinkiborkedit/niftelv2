package controlflow

import "github.com/ithinkiborkedit/niftelv2.git/internal/value"

type ControlFlow int

const (
	FlowNormal ControlFlow = iota
	FlowNone
	FlowReturn
	FlowBreak
	FlowContinue
)

type ExecResult struct {
	Value value.Value
	Flow  ControlFlow
	Err   error
}
