package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/envgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestEnumField(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, envgen.GenerateFile, testDataFS, "testdata/enum.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"github.com/stretchr/testify/assert"
			"github.com/stretchr/testify/require"
			"testing"
		)

		func TestNormalField(t *testing.T) {
			t.Setenv("ENUM", "ENUM_OPTION_A")

			actual, err := EnumMessageFromEnv()
			require.NoError(t, err)

			protoEqual(t, &EnumMessage{
				Normal: Enum_ENUM_OPTION_A,
			}, actual)
		}

		func TestPresenceField(t *testing.T) {
			t.Setenv("ENUM_WITH_PRESENCE", "ENUM_OPTION_A")

			actual, err := EnumMessageFromEnv()
			require.NoError(t, err)

			expectedValue := Enum_ENUM_OPTION_A
			protoEqual(t, &EnumMessage{
				WithPresence: &expectedValue,
			}, actual)
		}

		func TestList(t *testing.T) {
			t.Setenv("ENUM_LIST_1", "ENUM_OPTION_A")
			t.Setenv("ENUM_LIST_2", "ENUM_OPTION_B")
			t.Setenv("ENUM_LIST_3", "ENUM_UNSPECIFIED")

			actual, err := EnumMessageFromEnv()
			require.NoError(t, err)

			protoEqual(t, &EnumMessage{
				List: []Enum{Enum_ENUM_OPTION_A, Enum_ENUM_OPTION_B, Enum_ENUM_UNSPECIFIED},
			}, actual)
		}

		func TestOneOfOneOptionSet(t *testing.T) {
			t.Setenv("ENUM_ONEOF_A", "ENUM_OPTION_A")

			actual, err := EnumMessageFromEnv()
			require.NoError(t, err)

			protoEqual(t, &EnumMessage{
				OneofTest: &EnumMessage_OneofOptionA{
					OneofOptionA: Enum_ENUM_OPTION_A,
				},
			}, actual)
		}

		func TestOneOfMultipleSet(t *testing.T) {
			t.Setenv("ENUM_ONEOF_A", "ENUM_OPTION_B")
			t.Setenv("ENUM_ONEOF_B", "ENUM_OPTION_B")

			actual, err := EnumMessageFromEnv()
			assert.Error(t, err)
			assert.Nil(t, actual)
		}
	`)
}

func TestEmbeddedEnumField(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, envgen.GenerateFile, testDataFS, "testdata/enum.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"github.com/stretchr/testify/assert"
			"github.com/stretchr/testify/require"
			"testing"
		)

		func TestNormalField(t *testing.T) {
			t.Setenv("ENUM", "ENUM_OPTION_A")

			actual, err := EmbeddedEnumMessageFromEnv()
			require.NoError(t, err)

			protoEqual(t, &EmbeddedEnumMessage{
				Normal: EmbeddedEnumMessage_ENUM_OPTION_A,
			}, actual)
		}

		func TestPresenceField(t *testing.T) {
			t.Setenv("ENUM_WITH_PRESENCE", "ENUM_OPTION_A")

			actual, err := EmbeddedEnumMessageFromEnv()
			require.NoError(t, err)

			expectedValue := EmbeddedEnumMessage_ENUM_OPTION_A
			protoEqual(t, &EmbeddedEnumMessage{
				WithPresence: &expectedValue,
			}, actual)
		}

		func TestList(t *testing.T) {
			t.Setenv("ENUM_LIST_1", "ENUM_OPTION_A")
			t.Setenv("ENUM_LIST_2", "ENUM_OPTION_B")
			t.Setenv("ENUM_LIST_3", "ENUM_UNSPECIFIED")

			actual, err := EmbeddedEnumMessageFromEnv()
			require.NoError(t, err)

			protoEqual(t, &EmbeddedEnumMessage{
				List: []EmbeddedEnumMessage_EmbeddedEnum{
					EmbeddedEnumMessage_ENUM_OPTION_A,
					EmbeddedEnumMessage_ENUM_OPTION_B,
					EmbeddedEnumMessage_ENUM_UNSPECIFIED,
				},
			}, actual)
		}

		func TestOneOfOneOptionSet(t *testing.T) {
			t.Setenv("ENUM_ONEOF_A", "ENUM_OPTION_A")

			actual, err := EmbeddedEnumMessageFromEnv()
			require.NoError(t, err)

			protoEqual(t, &EmbeddedEnumMessage{
				OneofTest: &EmbeddedEnumMessage_OneofOptionA{
					OneofOptionA: EmbeddedEnumMessage_ENUM_OPTION_A,
				},
			}, actual)
		}

		func TestOneOfMultipleSet(t *testing.T) {
			t.Setenv("ENUM_ONEOF_A", "ENUM_OPTION_B")
			t.Setenv("ENUM_ONEOF_B", "ENUM_OPTION_B")

			actual, err := EmbeddedEnumMessageFromEnv()
			assert.Error(t, err)
			assert.Nil(t, actual)
		}
	`)
}
