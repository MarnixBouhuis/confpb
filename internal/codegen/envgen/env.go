package envgen

import (
	"errors"
	"fmt"
	"slices"

	"github.com/marnixbouhuis/confpb/internal/codegen"
	confpbv1 "github.com/marnixbouhuis/confpb/pkg/gen/confpb/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	runtimePackage = protogen.GoImportPath("github.com/marnixbouhuis/confpb/pkg/runtime")
	scanPackage    = protogen.GoImportPath("github.com/marnixbouhuis/confpb/pkg/runtime/scan")
	fmtPackage     = protogen.GoImportPath("fmt")
)

type envFieldDescriptor struct {
	field  *protogen.Field
	envKey string
}

var _ codegen.FileGeneratorFunc = GenerateFile

func GenerateFile(plugin *protogen.Plugin, file *protogen.File) error {
	fileName := file.GeneratedFilenamePrefix + ".envpb.go"
	g := codegen.NewFileWithBoilerplate(plugin, file, fileName)
	//nolint:wrapcheck // IterateMessages already supplies all the message info
	return codegen.IterateMessages(file.Messages, func(message *protogen.Message) error {
		return processMessage(g, message)
	})
}

func processMessage(g *protogen.GeneratedFile, message *protogen.Message) error {
	oneofFieldGroups := make(map[string][]*envFieldDescriptor)
	normalFields := make([]*envFieldDescriptor, 0, len(message.Fields))
	for _, field := range message.Fields {
		envKey, isStr := proto.GetExtension(field.Desc.Options(), confpbv1.E_Env).(string)
		if !isStr || envKey == "" {
			// Field does not have env option set, no need to generate code
			continue
		}

		if field.Desc.IsWeak() {
			return fmt.Errorf("field \"%s\" is invalid, weak fields are not supported", field.Desc.FullName())
		}

		if field.Desc.IsMap() {
			// We can't cleanly map environment variables to (possible nested) maps.
			// In most cases a list can also be used to model the data. Lists map better to environment variables.
			return fmt.Errorf("field \"%s\" is invalid, maps are not supported", field.Desc.FullName())
		}

		if field.Oneof != nil {
			oneof, hasOneof := oneofFieldGroups[field.Oneof.GoName]
			if !hasOneof {
				oneof = make([]*envFieldDescriptor, 0, 1)
			}
			oneof = append(oneof, &envFieldDescriptor{
				field:  field,
				envKey: envKey,
			})
			oneofFieldGroups[field.Oneof.GoName] = oneof
			continue
		}

		// Normal field, not part of an oneof group
		normalFields = append(normalFields, &envFieldDescriptor{
			field:  field,
			envKey: envKey,
		})
	}

	g.P("func ", message.GoIdent, "FromEnv() (*", message.GoIdent, ", error) {")
	g.P("m, _, err := ", message.GoIdent, "FromEnvWithPrefix(\"\")")
	g.P("return m, err")
	g.P("}")
	g.P()

	g.P("func ", message.GoIdent, "FromEnvWithPrefix(prefix string) (x *", message.GoIdent, ", fieldsPresent bool, err error) {")
	g.P("if prefix != \"\" {")
	g.P("prefix = prefix + \"_\"")
	g.P("}")
	g.P("")

	// Check if there are any environment variables starting with the prefix
	// We do this to prevent infinite loops while constructing recursively embedded messages.
	g.P("if hasKey := ", runtimePackage.Ident("HasEnvKeyWithPrefix"), "(prefix); !hasKey {")
	g.P("return nil, false, nil")
	g.P("}")
	g.P("")

	g.P("x = &", message.GoIdent, "{}")

	if len(oneofFieldGroups) == 0 && len(normalFields) == 0 {
		g.P("// No fields to scan, message has no fields with env option set")
	}

	if err := processOneOfGroups(g, oneofFieldGroups); err != nil {
		return err
	}

	if err := processNormalFields(g, normalFields); err != nil {
		return err
	}

	g.P("return x, fieldsPresent, nil")
	g.P("}")
	g.P()

	// Add compile time check to make sure XXXFromEnvWithPrefix follows the Scanner type
	g.P("var _ ", scanPackage.Ident("Scanner"), "[*", message.GoIdent, "] = ", message.GoIdent, "FromEnvWithPrefix")

	return nil
}

// processOneOfGroups generates environment scanning code for fields that are part of an oneof group.
func processOneOfGroups(g *protogen.GeneratedFile, oneofFieldGroups map[string][]*envFieldDescriptor) error {
	if len(oneofFieldGroups) == 0 {
		return nil
	}

	// To implement oneof fields we keep track if another field in the same oneof group has already been set.
	// We do this using a map with as key the oneof group name and as value the field that was set.
	// Every time we set a field that is part of an oneof group we first check if another field from this group has not
	// been set yet.
	g.P("oneofs := make(map[string]string)")

	// Get all oneof group names and sort the keys, this way code generation order is stable
	oneofGroupNames := make([]string, 0, len(oneofFieldGroups))
	for oneofGroupName := range oneofFieldGroups {
		oneofGroupNames = append(oneofGroupNames, oneofGroupName)
	}
	slices.Sort(oneofGroupNames)

	for _, oneofGroupName := range oneofGroupNames {
		for _, f := range oneofFieldGroups[oneofGroupName] {
			// Start a new block scope, this way we don't have to worry about conflicting variable names
			g.P("{")

			if err := generateFieldEnvScanner(g, f); err != nil {
				return fmt.Errorf("failed to process field (\"%s\"): %w", f.field.Desc.FullName(), err)
			}
			g.P("if err != nil {")
			g.P("return nil, false, ", fmtPackage.Ident("Errorf"), "(\"error scanning field \\\"", f.field.GoName, "\\\": %w\", err)")
			g.P("}")

			g.P("if hasResult {")
			g.P("fieldsPresent = true")

			// Make sure field is not already set by another field in this oneof group
			g.P("if conflictFieldName, alreadySet := oneofs[\"", oneofGroupName, "\"]; alreadySet {")
			g.P("return nil, false, ", fmtPackage.Ident("Errorf"), "(\"could not set field \\\"", f.field.GoName, "\\\", field is part of a oneof group and another group item is already set (\\\"%s\\\")\", conflictFieldName)")
			g.P("}")
			g.P("oneofs[\"", oneofGroupName, "\"] = \"", f.field.GoName, "\"")

			g.P("x.", f.field.Oneof.GoName, " = &", f.field.GoIdent, "{")
			g.P(f.field.GoName, ": result,") // No need to check for field presence since oneof fields can not have field presence
			g.P("}")

			// End for hasResult check
			g.P("}")

			// End for block
			g.P("}")
		}
	}

	return nil
}

// processNormalFields generates scanning env variable data for normal fields (not part of an oneof group).
func processNormalFields(g *protogen.GeneratedFile, fields []*envFieldDescriptor) error {
	for _, f := range fields {
		// Start a new block scope, this way we don't have to worry about conflicting variable names
		g.P("{")

		if err := generateFieldEnvScanner(g, f); err != nil {
			return fmt.Errorf("failed to process field (\"%s\"): %w", f.field.Desc.FullName(), err)
		}
		g.P("if err != nil {")
		g.P("return nil, false, ", fmtPackage.Ident("Errorf"), "(\"error scanning field \\\"", f.field.GoName, "\\\": %w\", err)")
		g.P("}")
		g.P("if hasResult {")
		g.P("fieldsPresent = true")
		if codegen.NeedsPointer(f.field) {
			g.P("x.", f.field.GoName, " = &result")
		} else {
			g.P("x.", f.field.GoName, " = result")
		}
		g.P("}")

		g.P("}")
	}
	return nil
}

func generateFieldEnvScanner(g *protogen.GeneratedFile, f *envFieldDescriptor) error {
	var scanner protogen.GoIdent
	switch f.field.Desc.Kind() {
	case protoreflect.BoolKind:
		scanner = scanPackage.Ident("Bool")
	case protoreflect.EnumKind:
		g.P("scanner := ", scanPackage.Ident("NewEnumScanner"), "[", f.field.Enum.GoIdent, "](", f.field.Enum.GoIdent, "_value)")
		scanner = protogen.GoIdent{GoName: "scanner", GoImportPath: f.field.GoIdent.GoImportPath}
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		scanner = scanPackage.Ident("Int32")
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		scanner = scanPackage.Ident("Int64")
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		scanner = scanPackage.Ident("Uint32")
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		scanner = scanPackage.Ident("Uint64")
	case protoreflect.FloatKind:
		scanner = scanPackage.Ident("Float")
	case protoreflect.DoubleKind:
		scanner = scanPackage.Ident("Double")
	case protoreflect.StringKind:
		scanner = scanPackage.Ident("String")
	case protoreflect.BytesKind:
		scanner = scanPackage.Ident("Bytes")
	case protoreflect.MessageKind:
		switch f.field.Message.Desc.FullName() {
		case "google.protobuf.Timestamp":
			scanner = scanPackage.Ident("Timestamp")
		case "google.protobuf.Duration":
			scanner = scanPackage.Ident("Duration")
		case "google.protobuf.Struct":
			scanner = scanPackage.Ident("Struct")
		case "google.protobuf.Value":
			scanner = scanPackage.Ident("Value")
		default:
			scanner = protogen.GoIdent{
				GoName:       f.field.Message.GoIdent.GoName + "FromEnvWithPrefix",
				GoImportPath: f.field.Message.GoIdent.GoImportPath,
			}
		}
	case protoreflect.GroupKind:
		// No support needed since we the minimum proto version that we support is proto3.
		return errors.New("groups are not supported")
	default:
		return fmt.Errorf("unknown field type: %s", f.field.Desc.Kind().String())
	}

	if f.field.Desc.IsList() {
		// Scan repeated field
		g.P("result, hasResult, err := ", scanPackage.Ident("Repeated"), "(prefix+\"", f.envKey, "\", ", scanner, ")")
		return nil
	}

	// Normal non-repeated field
	g.P("result, hasResult, err := ", scanner, "(prefix+\"", f.envKey, "\")")
	return nil
}
