package engine

type Environment struct {
	values map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{values: make(map[string]any)}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Assign(name string, value any) bool {
	if _, ok := e.values[name]; !ok {
		return false
	}
	e.values[name] = value
	return true
}

func (e *Environment) Get(name string) (any, bool) {
	value, ok := e.values[name]
	return value, ok
}
