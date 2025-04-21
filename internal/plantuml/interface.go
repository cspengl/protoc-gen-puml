package plantuml

import "fmt"

type Interface struct {
	Element
	methods []*method
}

type method struct {
	Name                        string
	InputMessage, OutputMessage string
}

func NewInterface(e Element) *Interface {
	return &Interface{
		Element: e,
		methods: make([]*method, 0),
	}
}

func (i *Interface) Render(t RenderTarget) {
	t.P(fmt.Sprintf("interface \"%s\" as %s {", i.Name, i.Fullname))
	t.IncrLevel()
	for _, m := range i.methods {
		t.P(fmt.Sprintf("%s(%s) %s", m.Name, m.InputMessage, m.OutputMessage))
	}
	t.DecrLevel()
	t.P("}")
}

func (i *Interface) AddMethod(name, input, output string) {
	i.methods = append(i.methods, &method{
		Name:          name,
		InputMessage:  input,
		OutputMessage: output,
	})
}

func (i *Interface) DelMethod(name string) {
	for j, method := range i.methods {
		if method.Name == name {
			i.methods = append(i.methods[:j], i.methods[j+1:]...)
		}
	}
}
