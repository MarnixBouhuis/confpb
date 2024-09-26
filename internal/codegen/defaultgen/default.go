package defaultgen

import (
	"fmt"
	"slices"

	"github.com/marnixbouhuis/confpb/internal/codegen"
	confpbv1 "github.com/marnixbouhuis/confpb/pkg/gen/confpb/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	timePackage        = protogen.GoImportPath("time")
	timestampPbPackage = protogen.GoImportPath("google.golang.org/protobuf/types/known/timestamppb")
	durationPbPackage  = protogen.GoImportPath("google.golang.org/protobuf/types/known/durationpb")
	structPbPackage    = protogen.GoImportPath("google.golang.org/protobuf/types/known/structpb")
	runtimePackage     = protogen.GoImportPath("github.com/marnixbouhuis/confpb/pkg/runtime")
)

var _ codegen.FileGeneratorFunc = GenerateFile

func GenerateFile(plugin *protogen.Plugin, file *protogen.File) error {
	fileName := file.GeneratedFilenamePrefix + ".defaultpb.go"
	g := codegen.NewFileWithBoilerplate(plugin, file, fileName)
	//nolint:wrapcheck // IterateMessages already supplies all the message info
	return codegen.IterateMessages(file.Messages, func(message *protogen.Message) error {
		return processMessage(g, message)
	})
}

func processMessage(g *protogen.GeneratedFile, message *protogen.Message) error {
	g.P("// ", message.GoIdent, "FromDefault returns a new instance of ", message.GoIdent, " containing only default values")
	g.P("func ", message.GoIdent, "FromDefault() *", message.GoIdent, " {")
	g.P("return (*", message.GoIdent, ")(nil).Default()")
	g.P("}")
	g.P()

	g.P("// Default returns a new instance of ", message.GoIdent, " containing only default values")
	g.P("func (*", message.GoIdent, ") Default() *", message.GoIdent, " {")
	g.P("x := &", message.GoIdent, "{}")

	oneofGroupsProcessed := make(map[string]*protogen.Field)
	for _, field := range message.Fields {
		defaultOption, isDefault := proto.GetExtension(field.Desc.Options(), confpbv1.E_Default).(*confpbv1.Default)
		if !isDefault || defaultOption == nil {
			// Field does not have env option set, no need to generate code
			continue
		}

		if field.Desc.IsWeak() {
			return fmt.Errorf("field \"%s\" is invalid, weak fields are not supported", field.Desc.FullName())
		}

		// Make sure that only one field in an oneof group can have the default tag.
		if field.Oneof != nil {
			// We are in an oneof group, check that we did not process another field from this group already.
			if conflictingField, alreadyProcessedOneof := oneofGroupsProcessed[field.Oneof.GoName]; alreadyProcessedOneof {
				return fmt.Errorf(
					"field \"%s\" is invalid. Field is part of oneof group \"%s\","+
						" only one field in an oneof group can have a default value set (conflicting field: \"%s\")",
					field.Desc.FullName(),
					field.Oneof.GoName,
					conflictingField.Desc.FullName(),
				)
			}
			oneofGroupsProcessed[field.Oneof.GoName] = field
		}

		if field.Desc.IsMap() {
			if err := processMapField(g, field, defaultOption); err != nil {
				return fmt.Errorf("field \"%s\" is invalid: %w", field.Desc.FullName(), err)
			}
			continue
		}

		if field.Desc.IsList() {
			if err := processRepeatedField(g, field, defaultOption); err != nil {
				return fmt.Errorf("field \"%s\" is invalid: %w", field.Desc.FullName(), err)
			}
			continue
		}

		if err := processSingleField(g, field, defaultOption); err != nil {
			return fmt.Errorf("field \"%s\" is invalid: %w", field.Desc.FullName(), err)
		}
	}

	g.P("return x")
	g.P("}")
	g.P()

	return nil
}

type pointerOf[T any] interface {
	~*T
}

// getDefaultForFieldType gets a default field option from the field_type oneof in *confpbv1.Default.
// The expected oneof type can be supplied with TPtr. It returns an error if the actual oneof option is different from
// the expected option.
func getDefaultForFieldType[TPtr pointerOf[T], T any](defaultOption *confpbv1.Default) (TPtr, error) {
	d, ok := defaultOption.FieldType.(TPtr)
	if !ok {
		// Get the expected type of the generic TPtr
		var expected TPtr
		return nil, fmt.Errorf("default option \"%T\" is invalid for this field, expected to have option \"%T\"", defaultOption.FieldType, expected)
	}
	return d, nil
}

// getMapKeyDefaultForFieldType gets a default field option from the key_type oneof in *confpbv1.Default_Map_Value.
// The expected oneof type can be supplied with TPtr. It returns an error if the actual oneof option is different from
// the expected option.
func getMapKeyDefaultForFieldType[TPtr pointerOf[T], T any](defaultMapOption *confpbv1.Default_Map_Value) (TPtr, error) {
	d, ok := defaultMapOption.KeyType.(TPtr)
	if !ok {
		// Get the expected type of the generic TPtr
		var expected TPtr
		return nil, fmt.Errorf("default map key type \"%T\" is invalid for this field, expected map key type to be \"%T\"", defaultMapOption.KeyType, expected)
	}
	return d, nil
}

// getMapValueDefaultForFieldType gets a default field option from the value_type oneof in *confpbv1.Default_Map_Value.
// The expected oneof type can be supplied with TPtr. It returns an error if the actual oneof option is different from
// the expected option.
func getMapValueDefaultForFieldType[TPtr pointerOf[T], T any](defaultMapOption *confpbv1.Default_Map_Value) (TPtr, error) {
	d, ok := defaultMapOption.ValueType.(TPtr)
	if !ok {
		// Get the expected type of the generic TPtr
		var expected TPtr
		return nil, fmt.Errorf("default map value type \"%T\" is invalid for this field, expected map value type to be \"%T\"", defaultMapOption.ValueType, expected)
	}
	return d, nil
}

// getValidOptionsForEnum returns all valid string options for a protobuf enum.
func getValidOptionsForEnum(enum *protogen.Enum) []string {
	supportedValues := make([]string, 0, len(enum.Values))
	for _, v := range enum.Values {
		supportedValues = append(supportedValues, string(v.Desc.Name()))
	}
	// Sort supported values so the error message is stable between compiles
	slices.Sort(supportedValues)
	return supportedValues
}

// fieldHasCircularDefaultMessageDependency checks if the message contains a circular dependency in its sub messages.
// It only checks messages that it needs to set the default value of.
// This can happen if the field is of type message, and somewhere in the dependency chain it points back to a previous
// message in the same chain.
func fieldHasCircularDefaultMessageDependency(parent *protogen.Field) error {
	if parent.Message == nil {
		return nil
	}

	visitedMessages := make(map[protoreflect.FullName]struct{})

	var checkCircularDependency func(field *protogen.Field, dependencyChain []protoreflect.FullName) ([]protoreflect.FullName, bool)
	checkCircularDependency = func(field *protogen.Field, dependencyChain []protoreflect.FullName) ([]protoreflect.FullName, bool) {
		// Only consider fields that are of type message
		if field.Message == nil {
			return nil, false
		}

		// Only consider message if it has the tag set to fill defaults for this message
		defaultOption, isDefault := proto.GetExtension(field.Desc.Options(), confpbv1.E_Default).(*confpbv1.Default)
		if !isDefault || defaultOption == nil {
			// Field does not have env option set, no need to check default
			return nil, false
		}

		if d, ok := defaultOption.FieldType.(*confpbv1.Default_Message_); !ok || !d.Message.FillDefaults {
			// No need to check message, we don't need to fill the defaults for this field
			return nil, false
		}

		name := field.Message.Desc.FullName()
		dependencyChain = append(dependencyChain, name)

		if _, visited := visitedMessages[name]; visited {
			// We already visited this message, we have a circular dependency
			return dependencyChain, true
		}
		visitedMessages[name] = struct{}{}

		for _, childField := range field.Message.Fields {
			if path, isCircular := checkCircularDependency(childField, dependencyChain); isCircular {
				return path, true
			}
		}

		return nil, false
	}

	if path, isCircular := checkCircularDependency(parent, []protoreflect.FullName{}); isCircular {
		// Pretty print circular dependency graph in the format:
		// ↳ first.Message
		//  ↳ other.A
		//    ↳ other.B
		//    ↑ ↳ other.C
		//    ↑   ↳ other.B
		//    ↑←←←←←↲
		indentation := 2

		var pathGraph string
		var circularImportIndex int
		var hasPassedCircularImport bool
		last := string(path[len(path)-1])
		for i, item := range path {
			for j := range i * indentation {
				if hasPassedCircularImport && j == circularImportIndex {
					pathGraph += "↑"
				} else {
					pathGraph += " "
				}
			}
			pathGraph += "↳ " + string(item) + "\n"

			if string(item) == last && !hasPassedCircularImport {
				circularImportIndex = i * indentation
				hasPassedCircularImport = true
			}
		}

		for i := range len(path) * indentation {
			if i < circularImportIndex {
				pathGraph += " "
				continue
			}
			if i == circularImportIndex {
				pathGraph += "↑"
				continue
			}
			pathGraph += "←"
		}
		pathGraph += "↲\n"

		return fmt.Errorf("circular dependency: \n%s", pathGraph)
	}
	return nil
}
