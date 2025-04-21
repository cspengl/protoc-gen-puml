package plantuml

type Diagram struct {
	*Container
}

func NewDiagram(title string) *Diagram {
	return &Diagram{
		Container: NewContainer(Element{
			Name:     title,
			Fullname: title,
		}),
	}
}

func (d *Diagram) Render(t RenderTarget) {
	t.P("@startuml ", d.Name)
	t.P("hide empty methods")
	t.P("hide empty fields")
	t.P()
	d.Container.Render(t)
	t.P()
	t.P("@enduml")
}
