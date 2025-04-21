package plantuml

import "io"

type Element struct {
	Name, Fullname string
}

type RenderTarget interface {
	io.Writer
	P(v ...any)
	IncrLevel()
	DecrLevel()
}

type Renderable interface {
	Render(RenderTarget)
}
