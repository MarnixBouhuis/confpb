package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/envgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestUint64Field(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, envgen.GenerateFile, testDataFS, "testdata/uint64.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"github.com/stretchr/testify/assert"
			"github.com/stretchr/testify/require"
			"testing"
		)

		func TestNormalField(t *testing.T) {
			t.Setenv("UINT64", "123")

			actual, err := Uint64FromEnv()
			require.NoError(t, err)

			protoEqual(t, &Uint64{
				Normal: 123,
			}, actual)
		}

		func TestPresenceField(t *testing.T) {
			t.Setenv("UINT64_WITH_PRESENCE", "123")

			actual, err := Uint64FromEnv()
			require.NoError(t, err)

			expectedValue := uint64(123)
			protoEqual(t, &Uint64{
				WithPresence: &expectedValue,
			}, actual)
		}

		func TestList(t *testing.T) {
			t.Setenv("UINT64_LIST_1", "123")
			t.Setenv("UINT64_LIST_2", "456")
			t.Setenv("UINT64_LIST_3", "987654321")

			actual, err := Uint64FromEnv()
			require.NoError(t, err)

			protoEqual(t, &Uint64{
				List: []uint64{123, 456, 987654321},
			}, actual)
		}

		func TestOneOfOneOptionSet(t *testing.T) {
			t.Setenv("UINT64_ONEOF_A", "123")

			actual, err := Uint64FromEnv()
			require.NoError(t, err)

			protoEqual(t, &Uint64{
				OneofTest: &Uint64_OneofOptionA{
					OneofOptionA: 123,
				},
			}, actual)
		}

		func TestOneOfMultipleSet(t *testing.T) {
			t.Setenv("UINT64_ONEOF_A", "123")
			t.Setenv("UINT64_ONEOF_B", "123")

			actual, err := Uint64FromEnv()
			assert.Error(t, err)
			assert.Nil(t, actual)
		}
	`)
}
