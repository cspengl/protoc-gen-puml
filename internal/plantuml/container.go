package plantuml

type Container struct {
	Element
	renderables []namedRenderable
	connections []*Connection
}

type namedRenderable struct {
	Renderable
	name string
}

func NewContainer(e Element) *Container {
	return &Container{
		Element:     e,
		renderables: make([]namedRenderable, 0),
		connections: make([]*Connection, 0),
	}
}

func (d *Container) Add(name string, r Renderable) {
	d.renderables = append(d.renderables, namedRenderable{
		name:       name,
		Renderable: r,
	})
}

func (d *Container) Get(name string) (Renderable, bool) {
	for _, r := range d.renderables {
		if r.name == name {
			return r.Renderable, true
		}
	}
	return nil, false
}

func (d *Container) AddConnection(c *Connection) {
	d.connections = append(d.connections, c)
}

func (d *Container) Render(t RenderTarget) {
	for _, r := range d.renderables {
		r.Render(t)
		t.P()
	}

	for _, c := range d.connections {
		c.Render(t)
	}
}

type Package struct {
	*Container
}

func (p *Package) Render(t RenderTarget) {
	t.P("package ", p.Name, " as ", p.Fullname, " {")
	t.IncrLevel()
	t.P()
	p.Container.Render(t)
	t.P()
	t.DecrLevel()
	t.P("}")
}
