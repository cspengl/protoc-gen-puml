package plantuml

import "fmt"

type Datatype struct {
	Name string
}

func (d *Datatype) New(elem Element) Renderable {
	return &datatypeInstance{
		TypeName: d.Name,
		Element:  elem,
	}
}

type datatypeInstance struct {
	TypeName string
	Element
}

func (d *datatypeInstance) Render(t RenderTarget) {
	t.P(fmt.Sprintf(
		"struct \"%s\" as %s<<(D,#FF7700) %s>>",
		d.Name, d.Fullname, d.TypeName,
	))
}
