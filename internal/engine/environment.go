package engine

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{values: make(map[string]any)}
}

func (e *Environment) Wrap() *Environment {
	return &Environment{enclosing: e, values: make(map[string]any)}
}

func (e *Environment) Unwrap() *Environment {
	return e.enclosing
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Assign(name string, value any) bool {
	if _, ok := e.values[name]; ok {
		e.values[name] = value
		return true
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}

	return false
}

func (e *Environment) Get(name string) (any, bool) {
	if value, ok := e.values[name]; ok {
		return value, true
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	return nil, false
}
