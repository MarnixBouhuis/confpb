package defaultgen

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	confpbv1 "github.com/marnixbouhuis/confpb/pkg/gen/confpb/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
)

func processSingleField(g *protogen.GeneratedFile, field *protogen.Field, defaultOption *confpbv1.Default) error {
	kind := field.Desc.Kind()
	switch {
	case kind == protoreflect.BoolKind:
		d, err := getDefaultForFieldType[*confpbv1.Default_Bool](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.Bool)
	case kind == protoreflect.Int32Kind:
		d, err := getDefaultForFieldType[*confpbv1.Default_Int32](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.Int32)
	case kind == protoreflect.Sint32Kind:
		d, err := getDefaultForFieldType[*confpbv1.Default_Sint32](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.Sint32)
	case kind == protoreflect.Sfixed32Kind:
		d, err := getDefaultForFieldType[*confpbv1.Default_Sfixed32](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.Sfixed32)
	case kind == protoreflect.Int64Kind:
		d, err := getDefaultForFieldType[*confpbv1.Default_Int64](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.Int64)
	case kind == protoreflect.Sint64Kind:
		d, err := getDefaultForFieldType[*confpbv1.Default_Sint64](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.Sint64)
	case kind == protoreflect.Sfixed64Kind:
		d, err := getDefaultForFieldType[*confpbv1.Default_Sfixed64](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.Sfixed64)
	case kind == protoreflect.Uint32Kind:
		d, err := getDefaultForFieldType[*confpbv1.Default_Uint32](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.Uint32)
	case kind == protoreflect.Fixed32Kind:
		d, err := getDefaultForFieldType[*confpbv1.Default_Fixed32](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.Fixed32)
	case kind == protoreflect.Uint64Kind:
		d, err := getDefaultForFieldType[*confpbv1.Default_Uint64](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.Uint64)
	case kind == protoreflect.Fixed64Kind:
		d, err := getDefaultForFieldType[*confpbv1.Default_Fixed64](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.Fixed64)
	case kind == protoreflect.FloatKind:
		d, err := getDefaultForFieldType[*confpbv1.Default_Float](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.Float)
	case kind == protoreflect.DoubleKind:
		d, err := getDefaultForFieldType[*confpbv1.Default_Double](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.Double)
	case kind == protoreflect.StringKind:
		d, err := getDefaultForFieldType[*confpbv1.Default_String_](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.String_)
	case kind == protoreflect.BytesKind:
		d, err := getDefaultForFieldType[*confpbv1.Default_Bytes](defaultOption)
		if err != nil {
			return err
		}
		bytes, err := base64.StdEncoding.DecodeString(d.Bytes)
		if err != nil {
			return fmt.Errorf("failed to base64 decode default value \"%s\": %w", d.Bytes, err)
		}
		return setValueForField(g, field, bytes)
	case kind == protoreflect.EnumKind:
		d, err := getDefaultForFieldType[*confpbv1.Default_Enum](defaultOption)
		if err != nil {
			return err
		}

		var chosenOption *protogen.EnumValue
		for _, enumOption := range field.Enum.Values {
			if string(enumOption.Desc.Name()) == d.Enum {
				chosenOption = enumOption
				break
			}
		}

		if chosenOption == nil {
			return fmt.Errorf(
				"default value \"%s\" is not valid for this enum, value must be one of: [%s]",
				d.Enum,
				strings.Join(getValidOptionsForEnum(field.Enum), ", "),
			)
		}

		return setValueForField(g, field, chosenOption)
	case kind == protoreflect.MessageKind && field.Message.Desc.FullName() == "google.protobuf.Timestamp":
		d, err := getDefaultForFieldType[*confpbv1.Default_Timestamp](defaultOption)
		if err != nil {
			return err
		}
		timestamp, err := time.Parse(time.RFC3339, d.Timestamp)
		if err != nil {
			return fmt.Errorf("default timestamp \"%s\" is invalid, unable to parse value as RFC3339 time string: %w", d.Timestamp, err)
		}
		return setValueForField(g, field, timestamp)
	case kind == protoreflect.MessageKind && field.Message.Desc.FullName() == "google.protobuf.Duration":
		d, err := getDefaultForFieldType[*confpbv1.Default_Duration](defaultOption)
		if err != nil {
			return err
		}
		duration, err := time.ParseDuration(d.Duration)
		if err != nil {
			return fmt.Errorf("default duration \"%s\" is invalid, unable to parse value: %w", d.Duration, err)
		}
		return setValueForField(g, field, duration)
	case kind == protoreflect.MessageKind && field.Message.Desc.FullName() == "google.protobuf.Struct":
		d, err := getDefaultForFieldType[*confpbv1.Default_Struct](defaultOption)
		if err != nil {
			return err
		}

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(d.Struct), &data); err != nil {
			return fmt.Errorf("invalid value \"%s\", unable to JSON decode object: %w", d.Struct, err)
		}

		s, err := structpb.NewStruct(data)
		if err != nil {
			return fmt.Errorf("invalid value \"%s\", unable to convert decoded data into struct: %w", d.Struct, err)
		}
		return setValueForField(g, field, s)
	case kind == protoreflect.MessageKind && field.Message.Desc.FullName() == "google.protobuf.Value":
		d, err := getDefaultForFieldType[*confpbv1.Default_Value](defaultOption)
		if err != nil {
			return err
		}

		var data interface{}
		if err := json.Unmarshal([]byte(d.Value), &data); err != nil {
			return fmt.Errorf("invalid value \"%s\", unable to JSON decode value: %w", d.Value, err)
		}

		v, err := structpb.NewValue(data)
		if err != nil {
			return fmt.Errorf("invalid value \"%s\", unable to convert decoded data into value: %w", d.Value, err)
		}
		return setValueForField(g, field, v)
	case kind == protoreflect.MessageKind:
		d, err := getDefaultForFieldType[*confpbv1.Default_Message_](defaultOption)
		if err != nil {
			return err
		}

		if !d.Message.FillDefaults {
			// No need to fill message, return empty message struct
			return setValueForField(g, field, &messageValue{
				Ident:       field.Message.GoIdent,
				FillDefault: false,
			})
		}

		// Make sure message has no circular dependency in its child messages. This will cause infinite loops when
		// loading defaults at runtime.
		if err := fieldHasCircularDefaultMessageDependency(field); err != nil {
			return err
		}

		return setValueForField(g, field, &messageValue{
			Ident:       field.Message.GoIdent,
			FillDefault: true,
		})
	default:
		return fmt.Errorf("unknown field type: %s", kind.String())
	}
}
