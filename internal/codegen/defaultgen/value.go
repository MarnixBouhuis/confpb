package defaultgen

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/marnixbouhuis/confpb/internal/codegen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/known/structpb"
)

type messageValue struct {
	Ident       protogen.GoIdent
	FillDefault bool
}

// setValueForField generates a line of code that sets a field of the message struct to a specified value.
// The value T is converted to its string / code representation.
// It returns an error of the type is not supported.
func setValueForField[T any](g *protogen.GeneratedFile, field *protogen.Field, value T) error {
	str, err := valueToString(g, value)
	if err != nil {
		return err
	}
	return setRawValueForField(g, field, str)
}

// setRawValueForField generates a line of code that sets a field of the message struct to a specific value.
// It expects the value as a raw code string. No conversion is done. The string is printed directly to the output file.
func setRawValueForField(g *protogen.GeneratedFile, field *protogen.Field, valueStr string) error {
	if codegen.NeedsPointer(field) {
		valueStr = g.QualifiedGoIdent(runtimePackage.Ident("Pointer")) + "(" + valueStr + ")"
	}

	if field.Oneof != nil {
		g.P("x.", field.Oneof.GoName, " = &", field.GoIdent, "{")
		g.P(field.GoName, ": ", valueStr, ",")
		g.P("}")
	} else {
		g.P("x.", field.GoName, " = ", valueStr)
	}

	return nil
}

// valueToString converts a go value to its code string representation.
func valueToString(g *protogen.GeneratedFile, value interface{}) (string, error) {
	switch v := value.(type) {
	case bool:
		str := "false"
		if v {
			str = "true"
		}
		return str, nil
	case int32:
		return "int32(" + strconv.FormatInt(int64(v), 10) + ")", nil
	case int64:
		return "int64(" + strconv.FormatInt(v, 10) + ")", nil
	case uint32:
		return "uint32(" + strconv.FormatUint(uint64(v), 10) + ")", nil
	case uint64:
		return "uint64(" + strconv.FormatUint(v, 10) + ")", nil
	case float32:
		return "float32(" + strconv.FormatFloat(float64(v), 'f', -1, 32) + ")", nil
	case float64:
		return "float64(" + strconv.FormatFloat(v, 'f', -1, 64) + ")", nil
	case string:
		return strconv.Quote(v), nil
	case []byte:
		bytesDecimal := make([]string, 0, len(v))
		for _, b := range v {
			bytesDecimal = append(bytesDecimal, strconv.FormatUint(uint64(b), 10))
		}
		return "[]byte{" + strings.Join(bytesDecimal, ",") + "}", nil
	case *protogen.EnumValue:
		return g.QualifiedGoIdent(v.GoIdent), nil
	case time.Time:
		newTimestampPb := g.QualifiedGoIdent(timestampPbPackage.Ident("New"))
		timeFromUnix := g.QualifiedGoIdent(timePackage.Ident("Unix"))
		timeInUnix := strconv.FormatInt(v.UnixNano(), 10)
		return newTimestampPb + "(" + timeFromUnix + "(0, " + timeInUnix + "))", nil
	case time.Duration:
		newDurationPb := g.QualifiedGoIdent(durationPbPackage.Ident("New"))
		duration := g.QualifiedGoIdent(timePackage.Ident("Duration"))
		durationNanos := strconv.FormatInt(int64(v), 10)
		return newDurationPb + "(" + duration + "(" + durationNanos + "))", nil
	case *structpb.Struct:
		return protoStructToString(g, v)
	case *structpb.Value:
		return protoStructValueToString(g, v)
	case *messageValue:
		messageType := g.QualifiedGoIdent(v.Ident)
		if !v.FillDefault {
			return "&" + messageType + "{}", nil
		}

		loaderStr := "func () *" + messageType + "{\n"
		loaderStr += "x := &" + messageType + "{}\n"
		loaderStr += "if withDefault, ok := interface{}(x).(interface{ Default() *" + messageType + " }); ok {\n"
		loaderStr += "return withDefault.Default()\n"
		loaderStr += "}\n"
		loaderStr += "return x\n"
		loaderStr += "}()"
		return loaderStr, nil
	case []bool:
		items, err := valuesToString(g, v)
		if err != nil {
			return "", err
		}
		return "[]bool{\n" + strings.Join(items, "") + "\n}", nil
	case []int32:
		items, err := valuesToString(g, v)
		if err != nil {
			return "", err
		}
		return "[]int32{\n" + strings.Join(items, "") + "\n}", nil
	case []int64:
		items, err := valuesToString(g, v)
		if err != nil {
			return "", err
		}
		return "[]int64{\n" + strings.Join(items, "") + "\n}", nil
	case []uint32:
		items, err := valuesToString(g, v)
		if err != nil {
			return "", err
		}
		return "[]uint32{\n" + strings.Join(items, "") + "\n}", nil
	case []uint64:
		items, err := valuesToString(g, v)
		if err != nil {
			return "", err
		}
		return "[]uint64{\n" + strings.Join(items, "") + "\n}", nil
	case []float32:
		items, err := valuesToString(g, v)
		if err != nil {
			return "", err
		}
		return "[]float32{\n" + strings.Join(items, "") + "\n}", nil
	case []float64:
		items, err := valuesToString(g, v)
		if err != nil {
			return "", err
		}
		return "[]float64{\n" + strings.Join(items, "") + "\n}", nil
	case []string:
		items, err := valuesToString(g, v)
		if err != nil {
			return "", err
		}
		return "[]string{\n" + strings.Join(items, "") + "\n}", nil
	case [][]byte:
		items, err := valuesToString(g, v)
		if err != nil {
			return "", err
		}
		return "[][]byte{\n" + strings.Join(items, "") + "\n}", nil
	case []*protogen.EnumValue:
		if len(v) == 0 {
			return "nil", nil // No items, we don't know the enum type
		}

		enumType := g.QualifiedGoIdent(v[0].Parent.GoIdent)
		items, err := valuesToString(g, v)
		if err != nil {
			return "", err
		}
		return "[]" + enumType + "{\n" + strings.Join(items, "") + "\n}", nil
	case []time.Time:
		items, err := valuesToString(g, v)
		if err != nil {
			return "", err
		}
		return "[]*" + g.QualifiedGoIdent(timestampPbPackage.Ident("Timestamp")) + "{\n" + strings.Join(items, "") + "\n}", nil
	case []time.Duration:
		items, err := valuesToString(g, v)
		if err != nil {
			return "", err
		}
		return "[]*" + g.QualifiedGoIdent(durationPbPackage.Ident("Duration")) + "{\n" + strings.Join(items, "") + "\n}", nil
	case []*structpb.Struct:
		items, err := valuesToString(g, v)
		if err != nil {
			return "", err
		}
		return "[]*" + g.QualifiedGoIdent(structPbPackage.Ident("Struct")) + "{\n" + strings.Join(items, "") + "\n}", nil
	case []*structpb.Value:
		items, err := valuesToString(g, v)
		if err != nil {
			return "", err
		}
		return "[]*" + g.QualifiedGoIdent(structPbPackage.Ident("Value")) + "{\n" + strings.Join(items, "") + "\n}", nil
	case []*messageValue:
		if len(v) == 0 {
			return "nil", nil // No items, we don't know the message type
		}

		messageType := g.QualifiedGoIdent(v[0].Ident)
		items, err := valuesToString(g, v)
		if err != nil {
			return "", err
		}
		return "[]*" + messageType + "{\n" + strings.Join(items, "") + "\n}", nil
	default:
		return "", fmt.Errorf("unknown value type: %T", value)
	}
}

func valuesToString[T any](g *protogen.GeneratedFile, values []T) ([]string, error) {
	items := make([]string, 0, len(values))
	for i, value := range values {
		str, err := valueToString(g, value)
		if err != nil {
			return nil, fmt.Errorf("failed to convert item (index=%d, value=\"%v\") to string: %w", i, value, err)
		}
		items = append(items, str+", \n")
	}
	return items, nil
}

func protoStructToString(g *protogen.GeneratedFile, s *structpb.Struct) (string, error) {
	// Get all fields of the struct
	keys := make([]string, 0, len(s.Fields))
	for k := range s.Fields {
		keys = append(keys, k)
	}

	// Sort keys so code generation is stable between compiled versions of the code generator
	slices.Sort(keys)

	out := "&" + g.QualifiedGoIdent(structPbPackage.Ident("Struct")) + "{\n"
	out += "Fields: map[string]*" + g.QualifiedGoIdent(structPbPackage.Ident("Value")) + "{\n"
	for _, key := range keys {
		valueStr, err := protoStructValueToString(g, s.Fields[key])
		if err != nil {
			return "", err
		}
		out += strconv.Quote(key) + ": " + valueStr + ",\n"
	}
	out += "},\n"
	out += "}"
	return out, nil
}

func protoStructValueToString(g *protogen.GeneratedFile, value *structpb.Value) (string, error) {
	switch kind := value.Kind.(type) {
	case *structpb.Value_NullValue:
		return g.QualifiedGoIdent(structPbPackage.Ident("NewNullValue")) + "()", nil
	case *structpb.Value_NumberValue:
		numberStr := strconv.FormatFloat(kind.NumberValue, 'f', -1, 64)
		return g.QualifiedGoIdent(structPbPackage.Ident("NewNumberValue")) + "(" + numberStr + ")", nil
	case *structpb.Value_StringValue:
		return g.QualifiedGoIdent(structPbPackage.Ident("NewStringValue")) + "(" + strconv.Quote(kind.StringValue) + ")", nil
	case *structpb.Value_BoolValue:
		if kind.BoolValue {
			return g.QualifiedGoIdent(structPbPackage.Ident("NewBoolValue")) + "(true)", nil
		}
		return g.QualifiedGoIdent(structPbPackage.Ident("NewBoolValue")) + "(false)", nil
	case *structpb.Value_StructValue:
		structStr, err := protoStructToString(g, kind.StructValue)
		if err != nil {
			return "", err
		}
		return g.QualifiedGoIdent(structPbPackage.Ident("NewStructValue")) + "(" + structStr + ")", nil
	case *structpb.Value_ListValue:
		newListValue := g.QualifiedGoIdent(structPbPackage.Ident("NewListValue"))
		listValueStr := newListValue + "(&" + g.QualifiedGoIdent(structPbPackage.Ident("ListValue")) + "{\n"
		listValueStr += "Values: []*" + g.QualifiedGoIdent(structPbPackage.Ident("Value")) + "{\n"

		for _, item := range kind.ListValue.Values {
			valueStr, err := protoStructValueToString(g, item)
			if err != nil {
				return "", err
			}
			listValueStr += valueStr + ",\n"
		}

		listValueStr += "},\n"
		listValueStr += "})"

		return listValueStr, nil
	default:
		return "", fmt.Errorf("unknown struct value type: %T", kind)
	}
}
