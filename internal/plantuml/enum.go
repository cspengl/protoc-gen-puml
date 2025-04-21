package plantuml

import (
	"strings"
)

type Enum struct {
	Element
	values []enumvalue
}

type enumvalue struct {
	Name  string
	Value int
}

func NewEnum(e Element) *Enum {
	return &Enum{
		Element: e,
		values:  make([]enumvalue, 0),
	}
}

func (e *Enum) AddValue(name string, value int) {
	e.values = append(e.values, enumvalue{
		Name:  name,
		Value: value,
	})
}

func (e *Enum) DelValue(name string) {
	for i, val := range e.values {
		if val.Name == name {
			e.values = append(e.values[:i], e.values[i+1:]...)
		}
	}
}

func (e *Enum) Render(t RenderTarget) {
	t.P("enum \"", e.Name, "\" as ", e.Fullname, "{")
	t.IncrLevel()
	for _, val := range e.values {
		t.P(strings.ToUpper(val.Name), " (", val.Value, ")")
	}
	t.DecrLevel()
	t.P("}")
}
