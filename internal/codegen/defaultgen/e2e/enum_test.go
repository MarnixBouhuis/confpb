package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/defaultgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestEnumField(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, defaultgen.GenerateFile, testDataFS, "testdata/enum.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"testing"
		)

		func TestDefaults(t *testing.T) {
			t.Parallel()
			actual := EnumMessageFromDefault()

			x := Enum_ENUM_OPTION_B
			protoEqual(t, &EnumMessage{
				Normal: Enum_ENUM_OPTION_A,
				WithPresence: &x,
				List: []Enum{Enum_ENUM_OPTION_B, Enum_ENUM_OPTION_A, Enum_ENUM_UNSPECIFIED},
				OneofTest: &EnumMessage_OneofOption{
					OneofOption: Enum_ENUM_OPTION_A,
				},
				Map: map[string]Enum{
					"key1": Enum_ENUM_OPTION_A,
					"key2": Enum_ENUM_OPTION_B,
				},
			}, actual)
		}

		func TestEmbeddedEnumDefaults(t *testing.T) {
			t.Parallel()
			actual := EmbeddedEnumMessageFromDefault()

			x := EmbeddedEnumMessage_ENUM_OPTION_B
			protoEqual(t, &EmbeddedEnumMessage{
				Normal: EmbeddedEnumMessage_ENUM_OPTION_A,
				WithPresence: &x,
				List: []EmbeddedEnumMessage_EmbeddedEnum{
					EmbeddedEnumMessage_ENUM_OPTION_B,
					EmbeddedEnumMessage_ENUM_OPTION_A,
					EmbeddedEnumMessage_ENUM_UNSPECIFIED,
				},
				OneofTest: &EmbeddedEnumMessage_OneofOption{
					OneofOption: EmbeddedEnumMessage_ENUM_OPTION_A,
				},
				Map: map[string]EmbeddedEnumMessage_EmbeddedEnum{
					"key1": EmbeddedEnumMessage_ENUM_OPTION_A,
					"key2": EmbeddedEnumMessage_ENUM_OPTION_B,
				},
			}, actual)
		}
	`)
}
