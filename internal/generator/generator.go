package generator

import (
	"fmt"
	"strings"

	"github.com/cspengl/protoc-gen-puml/internal/plantuml"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	doubleType   = &plantuml.Datatype{Name: "double"}
	floatType    = &plantuml.Datatype{Name: "float"}
	int32Type    = &plantuml.Datatype{Name: "int32"}
	int64Type    = &plantuml.Datatype{Name: "int64"}
	uint32Type   = &plantuml.Datatype{Name: "uint32"}
	uint64Type   = &plantuml.Datatype{Name: "uint64"}
	sint32Type   = &plantuml.Datatype{Name: "sint32"}
	sint64Type   = &plantuml.Datatype{Name: "sint64"}
	fixed32Type  = &plantuml.Datatype{Name: "fixed32"}
	fixed64Type  = &plantuml.Datatype{Name: "fixed64"}
	sfixed32Type = &plantuml.Datatype{Name: "sfixed32"}
	sfixed64Type = &plantuml.Datatype{Name: "sfixed64"}
	boolType     = &plantuml.Datatype{Name: "bool"}
	stringType   = &plantuml.Datatype{Name: "string"}
	bytesType    = &plantuml.Datatype{Name: "bytes"}
)

type Config struct {
	DiagramTitle string
}

type PUMLGenerator struct {
	config Config
}

func NewGenerator(config Config) *PUMLGenerator {
	return &PUMLGenerator{
		config: config,
	}
}

func (g *PUMLGenerator) Generate(files ...*protogen.File) *plantuml.Diagram {
	diagram := plantuml.NewDiagram(g.config.DiagramTitle)

	for _, file := range files {
		if !file.Generate {
			continue
		}
		g.genFile(diagram, file)
	}

	return diagram
}

func (g *PUMLGenerator) genFile(d *plantuml.Diagram, file *protogen.File) {
	var container *plantuml.Container

	if file.Desc.Package() != "" {
		pkg, ok := d.Get(string(file.Desc.Package()))
		if !ok {
			pkg = &plantuml.Package{
				Container: plantuml.NewContainer(plantuml.Element{
					Name:     string(file.Desc.Package()),
					Fullname: string(file.Desc.Package()),
				}),
			}
			d.Add(string(file.Desc.Package()), pkg)
		}
		container = pkg.(*plantuml.Package).Container
	} else {
		container = d.Container
	}

	// generate messages
	for _, message := range file.Messages {
		g.genMessage(container, message)
	}

	// generate enums
	for _, enum := range file.Enums {
		g.genEnum(container, enum)
	}

	// generate services
	for _, service := range file.Services {
		g.genService(container, service)
	}
}

func (g *PUMLGenerator) genMessage(d *plantuml.Container, message *protogen.Message) {
	s := plantuml.NewStruct(elemFromDesc(message.Desc))
	if message.Desc.IsMapEntry() {
		return
	}
	for _, field := range message.Fields {
		g.genField(s, field)
	}

	for _, oneof := range message.Oneofs {
		g.genOneOf(d, s, oneof)
	}

	for _, enum := range message.Enums {
		g.genEnum(d, enum)
	}

	for _, message := range message.Messages {
		g.genMessage(d, message)
	}

	d.Add(fullName(message.Desc), s)
}

func (g *PUMLGenerator) genField(s *plantuml.Struct, field *protogen.Field) {
	var typeName string
	if field.Oneof != nil {
		return
	}
	switch field.Desc.Kind() {
	case protoreflect.MessageKind:
		typeName = string(field.Message.Desc.Name())
	case protoreflect.EnumKind:
		typeName = string(field.Desc.Enum().Name())
	default:
		typeName = typeNameFromKind(field.Desc.Kind()).Name
	}

	if field.Desc.IsList() {
		typeName = typeName + "[]"
	}
	if field.Desc.IsMap() {
		var (
			mapType    = field.Desc.Message()
			keyField   = mapType.Fields().ByNumber(field.Desc.MapKey().Number())
			valueField = mapType.Fields().ByNumber(field.Desc.MapValue().Number())
			key        = typeNameFromKind(keyField.Kind()).Name
			value      = typeNameFromKind(valueField.Kind()).Name
		)
		if keyField.Kind() == protoreflect.MessageKind {
			key = string(keyField.Message().Name())
		}
		if valueField.Kind() == protoreflect.MessageKind {
			value = string(valueField.Message().Name())
		}
		typeName = fmt.Sprintf("map<%s,%s>", key, value)
	}

	// add attribute
	s.AddAttribute(string(field.Desc.Name()), typeName)
}

func (g *PUMLGenerator) genOneOf(d *plantuml.Container, s *plantuml.Struct, oneof *protogen.Oneof) {
	// add attribute to message
	s.AddAttribute(string(oneof.Desc.Name()), string(oneof.Desc.Name()))

	// add abstract class for oneof type
	d.Add(string(oneof.Desc.FullName()), &plantuml.Abstract{
		Element: elemFromDesc(oneof.Desc),
	})
	// add connection between oneof abstract type and message field
	d.AddConnection(&plantuml.Connection{
		Source: fmt.Sprintf("%s::%s", fullName(oneof.Parent.Desc), oneof.Desc.Name()),
		Target: fullName(oneof.Desc),
	})

	// iterate over fields of oneof
	for _, oneOfMessage := range oneof.Fields {
		var sourceName string
		// add field type
		switch oneOfMessage.Desc.Kind() {
		case protoreflect.MessageKind:
			sourceName = fullName(oneOfMessage.Message.Desc)
		case protoreflect.EnumKind:
			sourceName = fullName(oneOfMessage.Enum.Desc)
		default:
			datatype := typeNameFromKind(oneOfMessage.Desc.Kind())
			sourceName = string(fullName(oneOfMessage.Desc)) + "_" + datatype.Name
			d.Add(sourceName, datatype.New(
				plantuml.Element{
					Name:     string(oneOfMessage.Desc.Name()),
					Fullname: sourceName,
				},
			))
		}

		// add connection to abstract type
		d.AddConnection(&plantuml.Connection{
			Source:   sourceName,
			Target:   fullName(oneof.Desc),
			Relation: plantuml.Extension,
		})
	}
}

func (g *PUMLGenerator) genEnum(d *plantuml.Container, enum *protogen.Enum) {
	e := plantuml.NewEnum(elemFromDesc(enum.Desc))
	for _, val := range enum.Values {
		e.AddValue(string(val.Desc.Name()), int(val.Desc.Number()))
	}
	d.Add(fullName(enum.Desc), e)
}

func (g *PUMLGenerator) genService(d *plantuml.Container, service *protogen.Service) {
	iface := plantuml.NewInterface(elemFromDesc(service.Desc))
	for _, method := range service.Methods {
		iface.AddMethod(
			string(method.Desc.Name()),
			string(method.Input.Desc.Name()),
			string(method.Output.Desc.Name()),
		)
	}
	d.Add(fullName(service.Desc), iface)
}

func typeNameFromKind(k protoreflect.Kind) *plantuml.Datatype {
	switch k {
	case protoreflect.DoubleKind:
		return doubleType
	case protoreflect.FloatKind:
		return floatType
	case protoreflect.Int32Kind:
		return int32Type
	case protoreflect.Int64Kind:
		return int64Type
	case protoreflect.Uint32Kind:
		return uint32Type
	case protoreflect.Uint64Kind:
		return uint64Type
	case protoreflect.Sint32Kind:
		return sint32Type
	case protoreflect.Sint64Kind:
		return sint64Type
	case protoreflect.Fixed32Kind:
		return fixed32Type
	case protoreflect.Fixed64Kind:
		return fixed64Type
	case protoreflect.Sfixed32Kind:
		return sfixed32Type
	case protoreflect.Sfixed64Kind:
		return sfixed64Type
	case protoreflect.BoolKind:
		return boolType
	case protoreflect.StringKind:
		return stringType
	case protoreflect.BytesKind:
		return bytesType
	default:
		return &plantuml.Datatype{Name: k.String()}
	}
}

func elemFromDesc(d protoreflect.Descriptor) plantuml.Element {
	return plantuml.Element{
		Name:     string(d.Name()),
		Fullname: fullName(d),
	}
}

func fullName(d protoreflect.Descriptor) string {
	return strings.ReplaceAll(string(d.FullName()), ".", "_")
}
