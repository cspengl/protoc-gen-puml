package plantuml

import (
	"fmt"
)

type attributeContainer struct {
	Attributes []attribute
}

func (c *attributeContainer) AddAttribute(name, typename string) {
	c.Attributes = append(c.Attributes, attribute{
		Name: name,
		Type: typename,
	})
}

func (c *attributeContainer) DelAttribute(name string) {
	for i, a := range c.Attributes {
		if a.Name == name {
			c.Attributes = append(c.Attributes[:i], c.Attributes[i+1:]...)
		}
	}
}

func (c *attributeContainer) Render(t RenderTarget) {
	if c == nil {
		return
	}
	for _, attr := range c.Attributes {
		t.P(fmt.Sprintf("%s : %s", attr.Name, attr.Type))
	}
}

type attribute struct {
	Name, Type string
}

func NewStruct(e Element) *Struct {
	return &Struct{
		Element:            e,
		attributeContainer: new(attributeContainer),
	}
}

type Struct struct {
	Element
	*attributeContainer
}

func (s *Struct) Render(t RenderTarget) {
	t.P(fmt.Sprintf("struct \"%s\" as %s {", s.Name, s.Fullname))
	t.IncrLevel()
	s.attributeContainer.Render(t)
	t.DecrLevel()
	t.P("}")
}

type Abstract struct {
	Element
}

func (s *Abstract) Render(t RenderTarget) {
	t.P(fmt.Sprintf("abstract \"%s\" as %s", s.Name, s.Fullname))
}
