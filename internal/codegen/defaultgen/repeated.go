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

func processRepeatedField(g *protogen.GeneratedFile, field *protogen.Field, defaultOption *confpbv1.Default) error {
	kind := field.Desc.Kind()
	switch {
	case kind == protoreflect.BoolKind:
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedBool_](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.RepeatedBool.Values)
	case kind == protoreflect.Int32Kind:
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedInt32_](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.RepeatedInt32.Values)
	case kind == protoreflect.Sint32Kind:
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedSint32_](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.RepeatedSint32.Values)
	case kind == protoreflect.Sfixed32Kind:
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedSfixed32_](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.RepeatedSfixed32.Values)
	case kind == protoreflect.Int64Kind:
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedInt64_](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.RepeatedInt64.Values)
	case kind == protoreflect.Sint64Kind:
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedSint64_](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.RepeatedSint64.Values)
	case kind == protoreflect.Sfixed64Kind:
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedSfixed64_](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.RepeatedSfixed64.Values)
	case kind == protoreflect.Uint32Kind:
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedUint32_](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.RepeatedUint32.Values)
	case kind == protoreflect.Fixed32Kind:
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedFixed32_](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.RepeatedFixed32.Values)
	case kind == protoreflect.Uint64Kind:
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedUint64_](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.RepeatedUint64.Values)
	case kind == protoreflect.Fixed64Kind:
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedFixed64_](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.RepeatedFixed64.Values)
	case kind == protoreflect.FloatKind:
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedFloat_](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.RepeatedFloat.Values)
	case kind == protoreflect.DoubleKind:
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedDouble_](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.RepeatedDouble.Values)
	case kind == protoreflect.StringKind:
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedString_](defaultOption)
		if err != nil {
			return err
		}
		return setValueForField(g, field, d.RepeatedString.Values)
	case kind == protoreflect.BytesKind:
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedBytes_](defaultOption)
		if err != nil {
			return err
		}

		decodedItems := make([][]byte, 0, len(d.RepeatedBytes.Values))
		for _, str := range d.RepeatedBytes.Values {
			bytes, err := base64.StdEncoding.DecodeString(str)
			if err != nil {
				return fmt.Errorf("failed to base64 decode default value \"%s\": %w", str, err)
			}
			decodedItems = append(decodedItems, bytes)
		}

		return setValueForField(g, field, decodedItems)
	case kind == protoreflect.EnumKind:
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedEnum_](defaultOption)
		if err != nil {
			return err
		}

		validOptions := make(map[string]*protogen.EnumValue)
		for _, option := range field.Enum.Values {
			validOptions[string(option.Desc.Name())] = option
		}

		items := make([]*protogen.EnumValue, 0, len(d.RepeatedEnum.Values))
		for _, str := range d.RepeatedEnum.Values {
			value, isValid := validOptions[str]
			if !isValid {
				return fmt.Errorf(
					"default value \"%s\" is not valid for this enum, value must be one of: [%s]",
					str,
					strings.Join(getValidOptionsForEnum(field.Enum), ", "),
				)
			}
			items = append(items, value)
		}
		return setValueForField(g, field, items)
	case kind == protoreflect.MessageKind && field.Message.Desc.FullName() == "google.protobuf.Timestamp":
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedTimestamp_](defaultOption)
		if err != nil {
			return err
		}

		timestamps := make([]time.Time, 0, len(d.RepeatedTimestamp.Values))
		for _, timestampStr := range d.RepeatedTimestamp.Values {
			timestamp, err := time.Parse(time.RFC3339, timestampStr)
			if err != nil {
				return fmt.Errorf("default timestamp \"%s\" is invalid, unable to parse value as RFC3339 time string: %w", timestampStr, err)
			}
			timestamps = append(timestamps, timestamp)
		}
		return setValueForField(g, field, timestamps)
	case kind == protoreflect.MessageKind && field.Message.Desc.FullName() == "google.protobuf.Duration":
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedDuration_](defaultOption)
		if err != nil {
			return err
		}

		durations := make([]time.Duration, 0, len(d.RepeatedDuration.Values))
		for _, durationStr := range d.RepeatedDuration.Values {
			duration, err := time.ParseDuration(durationStr)
			if err != nil {
				return fmt.Errorf("default duration \"%s\" is invalid, unable to parse value: %w", durationStr, err)
			}
			durations = append(durations, duration)
		}
		return setValueForField(g, field, durations)
	case kind == protoreflect.MessageKind && field.Message.Desc.FullName() == "google.protobuf.Struct":
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedStruct_](defaultOption)
		if err != nil {
			return err
		}

		structs := make([]*structpb.Struct, 0, len(d.RepeatedStruct.Values))
		for _, structStr := range d.RepeatedStruct.Values {
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(structStr), &data); err != nil {
				return fmt.Errorf("invalid value \"%s\", unable to JSON decode object: %w", structStr, err)
			}

			s, err := structpb.NewStruct(data)
			if err != nil {
				return fmt.Errorf("invalid value \"%s\", unable to convert decoded data into struct: %w", structStr, err)
			}
			structs = append(structs, s)
		}
		return setValueForField(g, field, structs)
	case kind == protoreflect.MessageKind && field.Message.Desc.FullName() == "google.protobuf.Value":
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedValue_](defaultOption)
		if err != nil {
			return err
		}

		values := make([]*structpb.Value, 0, len(d.RepeatedValue.Values))
		for _, valueStr := range d.RepeatedValue.Values {
			var data interface{}
			if err := json.Unmarshal([]byte(valueStr), &data); err != nil {
				return fmt.Errorf("invalid value \"%s\", unable to JSON decode value: %w", valueStr, err)
			}

			v, err := structpb.NewValue(data)
			if err != nil {
				return fmt.Errorf("invalid value \"%s\", unable to convert decoded data into value: %w", valueStr, err)
			}
			values = append(values, v)
		}
		return setValueForField(g, field, values)
	case kind == protoreflect.MessageKind:
		d, err := getDefaultForFieldType[*confpbv1.Default_RepeatedMessage_](defaultOption)
		if err != nil {
			return err
		}

		values := make([]*messageValue, 0, len(d.RepeatedMessage.Values))
		for _, value := range d.RepeatedMessage.Values {
			values = append(values, &messageValue{
				Ident:       field.Message.GoIdent,
				FillDefault: value.FillDefaults,
			})
		}
		return setValueForField(g, field, values)
	default:
		return fmt.Errorf("unknown field type: %s", kind.String())
	}
}
