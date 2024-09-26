package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/envgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestFixed64Field(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, envgen.GenerateFile, testDataFS, "testdata/fixed64.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"github.com/stretchr/testify/assert"
			"github.com/stretchr/testify/require"
			"testing"
		)

		func TestNormalField(t *testing.T) {
			t.Setenv("FIXED64", "123")

			actual, err := Fixed64FromEnv()
			require.NoError(t, err)

			protoEqual(t, &Fixed64{
				Normal: 123,
			}, actual)
		}

		func TestPresenceField(t *testing.T) {
			t.Setenv("FIXED64_WITH_PRESENCE", "123")

			actual, err := Fixed64FromEnv()
			require.NoError(t, err)

			expectedValue := uint64(123)
			protoEqual(t, &Fixed64{
				WithPresence: &expectedValue,
			}, actual)
		}

		func TestList(t *testing.T) {
			t.Setenv("FIXED64_LIST_1", "123")
			t.Setenv("FIXED64_LIST_2", "456")
			t.Setenv("FIXED64_LIST_3", "987654321")

			actual, err := Fixed64FromEnv()
			require.NoError(t, err)

			protoEqual(t, &Fixed64{
				List: []uint64{123, 456, 987654321},
			}, actual)
		}

		func TestOneOfOneOptionSet(t *testing.T) {
			t.Setenv("FIXED64_ONEOF_A", "123")

			actual, err := Fixed64FromEnv()
			require.NoError(t, err)

			protoEqual(t, &Fixed64{
				OneofTest: &Fixed64_OneofOptionA{
					OneofOptionA: 123,
				},
			}, actual)
		}

		func TestOneOfMultipleSet(t *testing.T) {
			t.Setenv("FIXED64_ONEOF_A", "123")
			t.Setenv("FIXED64_ONEOF_B", "123")

			actual, err := Fixed64FromEnv()
			assert.Error(t, err)
			assert.Nil(t, actual)
		}
	`)
}
