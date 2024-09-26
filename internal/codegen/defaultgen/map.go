package defaultgen

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	confpbv1 "github.com/marnixbouhuis/confpb/pkg/gen/confpb/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
)

func processMapField(g *protogen.GeneratedFile, field *protogen.Field, defaultOption *confpbv1.Default) error {
	if !field.Desc.IsMap() {
		return errors.New("field is not a map")
	}

	d, err := getDefaultForFieldType[*confpbv1.Default_Map_](defaultOption)
	if err != nil {
		return err
	}

	keyField := field.Message.Fields[0]
	valueField := field.Message.Fields[1]

	keyType, err := fieldTypeToString(g, keyField)
	if err != nil {
		return fmt.Errorf("failed to get type for key field: %w", err)
	}

	valueType, err := fieldTypeToString(g, valueField)
	if err != nil {
		return fmt.Errorf("failed to get type for value field: %w", err)
	}

	mapStr := fmt.Sprintf("map[%s]%s{\n", keyType, valueType)
	for _, values := range d.Map.Values {
		keyStr, err := mapDefaultKeyToString(g, keyField, values)
		if err != nil {
			return err
		}

		valueStr, err := mapDefaultValueToString(g, valueField, values)
		if err != nil {
			return err
		}

		mapStr += keyStr + ": " + valueStr + ",\n"
	}
	mapStr += "}"
	return setRawValueForField(g, field, mapStr)
}

func mapDefaultKeyToString(g *protogen.GeneratedFile, keyField *protogen.Field, mapValue *confpbv1.Default_Map_Value) (string, error) {
	//nolint:exhaustive // We only check field types supported as map keys.
	switch keyField.Desc.Kind() {
	case protoreflect.Int32Kind:
		d, err := getMapKeyDefaultForFieldType[*confpbv1.Default_Map_Value_Int32Key](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.Int32Key)
	case protoreflect.Int64Kind:
		d, err := getMapKeyDefaultForFieldType[*confpbv1.Default_Map_Value_Int64Key](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.Int64Key)
	case protoreflect.Uint32Kind:
		d, err := getMapKeyDefaultForFieldType[*confpbv1.Default_Map_Value_Uint32Key](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.Uint32Key)
	case protoreflect.Uint64Kind:
		d, err := getMapKeyDefaultForFieldType[*confpbv1.Default_Map_Value_Uint64Key](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.Uint64Key)
	case protoreflect.Sint32Kind:
		d, err := getMapKeyDefaultForFieldType[*confpbv1.Default_Map_Value_Sint32Key](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.Sint32Key)
	case protoreflect.Sint64Kind:
		d, err := getMapKeyDefaultForFieldType[*confpbv1.Default_Map_Value_Sint64Key](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.Sint64Key)
	case protoreflect.Fixed32Kind:
		d, err := getMapKeyDefaultForFieldType[*confpbv1.Default_Map_Value_Fixed32Key](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.Fixed32Key)
	case protoreflect.Fixed64Kind:
		d, err := getMapKeyDefaultForFieldType[*confpbv1.Default_Map_Value_Fixed64Key](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.Fixed64Key)
	case protoreflect.Sfixed32Kind:
		d, err := getMapKeyDefaultForFieldType[*confpbv1.Default_Map_Value_Sfixed32Key](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.Sfixed32Key)
	case protoreflect.Sfixed64Kind:
		d, err := getMapKeyDefaultForFieldType[*confpbv1.Default_Map_Value_Sfixed64Key](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.Sfixed64Key)
	case protoreflect.BoolKind:
		d, err := getMapKeyDefaultForFieldType[*confpbv1.Default_Map_Value_BoolKey](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.BoolKey)
	case protoreflect.StringKind:
		d, err := getMapKeyDefaultForFieldType[*confpbv1.Default_Map_Value_StringKey](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.StringKey)
	default:
		return "", fmt.Errorf("unknown map key type: %s", keyField.Desc.Kind().String())
	}
}

func mapDefaultValueToString(g *protogen.GeneratedFile, valueField *protogen.Field, mapValue *confpbv1.Default_Map_Value) (string, error) {
	kind := valueField.Desc.Kind()
	switch {
	case kind == protoreflect.DoubleKind:
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_DoubleValue](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.DoubleValue)
	case kind == protoreflect.FloatKind:
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_FloatValue](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.FloatValue)
	case kind == protoreflect.Int32Kind:
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_Int32Value](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.Int32Value)
	case kind == protoreflect.Int64Kind:
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_Int64Value](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.Int64Value)
	case kind == protoreflect.Uint32Kind:
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_Uint32Value](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.Uint32Value)
	case kind == protoreflect.Uint64Kind:
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_Uint64Value](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.Uint64Value)
	case kind == protoreflect.Sint32Kind:
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_Sint32Value](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.Sint32Value)
	case kind == protoreflect.Sint64Kind:
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_Sint64Value](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.Sint64Value)
	case kind == protoreflect.Fixed32Kind:
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_Fixed32Value](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.Fixed32Value)
	case kind == protoreflect.Fixed64Kind:
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_Fixed64Value](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.Fixed64Value)
	case kind == protoreflect.Sfixed32Kind:
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_Sfixed32Value](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.Sfixed32Value)
	case kind == protoreflect.Sfixed64Kind:
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_Sfixed64Value](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.Sfixed64Value)
	case kind == protoreflect.BoolKind:
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_BoolValue](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.BoolValue)
	case kind == protoreflect.StringKind:
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_StringValue](mapValue)
		if err != nil {
			return "", err
		}
		return valueToString(g, d.StringValue)
	case kind == protoreflect.BytesKind:
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_BytesValue](mapValue)
		if err != nil {
			return "", err
		}
		bytes, err := base64.StdEncoding.DecodeString(d.BytesValue)
		if err != nil {
			return "", fmt.Errorf("failed to base64 decode map default value \"%s\": %w", d.BytesValue, err)
		}
		return valueToString(g, bytes)
	case kind == protoreflect.EnumKind:
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_EnumValue](mapValue)
		if err != nil {
			return "", err
		}

		var chosenOption *protogen.EnumValue
		for _, enumOption := range valueField.Enum.Values {
			if string(enumOption.Desc.Name()) == d.EnumValue {
				chosenOption = enumOption
				break
			}
		}

		if chosenOption == nil {
			return "", fmt.Errorf(
				"default value \"%s\" is not valid for this enum, value must be one of: [%s]",
				d.EnumValue,
				strings.Join(getValidOptionsForEnum(valueField.Enum), ", "),
			)
		}

		return valueToString(g, chosenOption)
	case kind == protoreflect.MessageKind && valueField.Message.Desc.FullName() == "google.protobuf.Timestamp":
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_TimestampValue](mapValue)
		if err != nil {
			return "", err
		}
		timestamp, err := time.Parse(time.RFC3339, d.TimestampValue)
		if err != nil {
			return "", fmt.Errorf("default timestamp \"%s\" is invalid, unable to parse value as RFC3339 time string: %w", d.TimestampValue, err)
		}
		return valueToString(g, timestamp)
	case kind == protoreflect.MessageKind && valueField.Message.Desc.FullName() == "google.protobuf.Duration":
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_DurationValue](mapValue)
		if err != nil {
			return "", err
		}
		duration, err := time.ParseDuration(d.DurationValue)
		if err != nil {
			return "", fmt.Errorf("default duration \"%s\" is invalid, unable to parse value: %w", d.DurationValue, err)
		}
		return valueToString(g, duration)
	case kind == protoreflect.MessageKind && valueField.Message.Desc.FullName() == "google.protobuf.Struct":
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_StructValue](mapValue)
		if err != nil {
			return "", err
		}

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(d.StructValue), &data); err != nil {
			return "", fmt.Errorf("invalid value \"%s\", unable to JSON decode object: %w", d.StructValue, err)
		}

		s, err := structpb.NewStruct(data)
		if err != nil {
			return "", fmt.Errorf("invalid value \"%s\", unable to convert decoded data into struct: %w", d.StructValue, err)
		}
		return valueToString(g, s)
	case kind == protoreflect.MessageKind && valueField.Message.Desc.FullName() == "google.protobuf.Value":
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_ValueValue](mapValue)
		if err != nil {
			return "", err
		}

		var data interface{}
		if err := json.Unmarshal([]byte(d.ValueValue), &data); err != nil {
			return "", fmt.Errorf("invalid value \"%s\", unable to JSON decode value: %w", d.ValueValue, err)
		}

		v, err := structpb.NewValue(data)
		if err != nil {
			return "", fmt.Errorf("invalid value \"%s\", unable to convert decoded data into value: %w", d.ValueValue, err)
		}
		return valueToString(g, v)
	case kind == protoreflect.MessageKind:
		d, err := getMapValueDefaultForFieldType[*confpbv1.Default_Map_Value_MessageValue](mapValue)
		if err != nil {
			return "", err
		}

		if !d.MessageValue.FillDefaults {
			// No need to fill message, return empty message struct
			return valueToString(g, &messageValue{
				Ident:       valueField.Message.GoIdent,
				FillDefault: false,
			})
		}

		// Make sure message has no circular dependency in its child messages. This will cause infinite loops when
		// loading defaults at runtime.
		if err := fieldHasCircularDefaultMessageDependency(valueField); err != nil {
			return "", err
		}

		return valueToString(g, &messageValue{
			Ident:       valueField.Message.GoIdent,
			FillDefault: true,
		})
	default:
		return "", fmt.Errorf("unknown map value type: %s", valueField.Desc.Kind().String())
	}
}
