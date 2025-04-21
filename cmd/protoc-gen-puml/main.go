// protoc-gen-puml is a plugin for the Google protocol buffer compiler to
// generate Plantuml diagrams. Install it by building this program and making
// it accessible within your PATH with the name:
//
//	protoc-gen-puml
//
// The 'puml' suffix becomes part of the argument for protoc such that it can
// be invoked as:
//
//	protoc --puml_out=. path/to/file.proto
//
// This generates a Plantuml diagram for the protocol buffer defined by your
// file.proto. With that input, the output will be written to:
//
// ./diagram.pb.puml
package main

import (
	"flag"

	"github.com/cspengl/protoc-gen-puml/internal/generator"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
)

var (
	config generator.Config

	flags flag.FlagSet

	outputFile = flags.String("output", "diagram.pb.puml", "output file")
)

func main() {
	flags.StringVar(&config.DiagramTitle, "title", "diagram", "title of the diagram")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(runPlugin)
}

func runPlugin(p *protogen.Plugin) error {
	p.SupportedEditionsMinimum = descriptorpb.Edition_EDITION_PROTO3
	p.SupportedEditionsMaximum = descriptorpb.Edition_EDITION_PROTO3

	// generate diagram
	diagram := generator.NewGenerator(config).Generate(p.Files...)

	// render to file
	indentWriter := &generator.IndentWriter{
		Indent:  "\t",
		Wrapped: p.NewGeneratedFile(*outputFile, ""),
	}
	diagram.Render(indentWriter)
	return nil
}
