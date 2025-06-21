package environment

import (
	"fmt"
	"sync"

	"github.com/ithinkiborkedit/niftelv2.git/internal/value"
)

type Environment struct {
	variables map[string]value.Value
	enclosing *Environment
	mu        sync.RWMutex
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		variables: make(map[string]value.Value),
		enclosing: enclosing,
	}
}

func (e *Environment) Define(name string, val value.Value) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.variables[name] = val
}

func (e *Environment) Assign(name string, val value.Value) error {
	e.mu.Lock()
	defer e.mu.Unlock()

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
	e.mu.RLock()
	val, ok := e.variables[name]
	e.mu.RUnlock()
	if ok {
		return val, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	return value.Null(), fmt.Errorf("undefined variable '%s'", name)
}
