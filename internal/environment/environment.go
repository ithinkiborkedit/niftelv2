package environment

import (
	"fmt"

	"github.com/ithinkiborkedit/niftelv2.git/internal/value"
)

type Environment struct {
	variables map[string]value.Value
	enclosing *Environment
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		variables: make(map[string]value.Value),
		enclosing: enclosing,
	}
}

func (e *Environment) Define(name string, val value.Value) {
	e.variables[name] = val
}

func (e *Environment) Assign(name string, val value.Value) error {

	if _, ok := e.variables[name]; ok {
		e.variables[name] = val
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.Assign(name, val)
	}

	return fmt.Errorf("undefined variable '%s'", name)
}

func (e *Environment) Get(name string) (value.Value, error) {
	val, ok := e.variables[name]
	if ok {
		return val, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	return value.Null(), fmt.Errorf("undefined variable '%s'", name)
}
