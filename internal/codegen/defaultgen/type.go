package defaultgen

import (
	"errors"
	"fmt"

	"github.com/marnixbouhuis/confpb/internal/codegen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// fieldTypeToString converts a protobuf field descriptor to its go code type string representation.
func fieldTypeToString(g *protogen.GeneratedFile, field *protogen.Field) (string, error) {
	var typeString string
	if codegen.NeedsPointer(field) {
		typeString += "*"
	}

	if field.Desc.IsList() {
		typeString += "[]"
	}

	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		return typeString + "bool", nil
	case protoreflect.EnumKind:
		return typeString + g.QualifiedGoIdent(field.Enum.GoIdent), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return typeString + "int32", nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return typeString + "int64", nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return typeString + "uint32", nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return typeString + "uint64", nil
	case protoreflect.FloatKind:
		return typeString + "float32", nil
	case protoreflect.DoubleKind:
		return typeString + "float64", nil
	case protoreflect.StringKind:
		return typeString + "string", nil
	case protoreflect.BytesKind:
		return typeString + "[]byte", nil
	case protoreflect.MessageKind:
		return typeString + "*" + g.QualifiedGoIdent(field.Message.GoIdent), nil
	case protoreflect.GroupKind:
		// No support needed since we the minimum proto version that we support is proto3.
		return "", errors.New("groups are not supported")
	default:
		return "", fmt.Errorf("unknown field type: %T", field.Desc.Kind())
	}
}
