package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/envgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestUint32Field(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, envgen.GenerateFile, testDataFS, "testdata/uint32.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"github.com/stretchr/testify/assert"
			"github.com/stretchr/testify/require"
			"testing"
		)

		func TestNormalField(t *testing.T) {
			t.Setenv("UINT32", "123")

			actual, err := Uint32FromEnv()
			require.NoError(t, err)

			protoEqual(t, &Uint32{
				Normal: 123,
			}, actual)
		}

		func TestPresenceField(t *testing.T) {
			t.Setenv("UINT32_WITH_PRESENCE", "123")

			actual, err := Uint32FromEnv()
			require.NoError(t, err)

			expectedValue := uint32(123)
			protoEqual(t, &Uint32{
				WithPresence: &expectedValue,
			}, actual)
		}

		func TestList(t *testing.T) {
			t.Setenv("UINT32_LIST_1", "123")
			t.Setenv("UINT32_LIST_2", "456")
			t.Setenv("UINT32_LIST_3", "987654321")

			actual, err := Uint32FromEnv()
			require.NoError(t, err)

			protoEqual(t, &Uint32{
				List: []uint32{123, 456, 987654321},
			}, actual)
		}

		func TestOneOfOneOptionSet(t *testing.T) {
			t.Setenv("UINT32_ONEOF_A", "123")

			actual, err := Uint32FromEnv()
			require.NoError(t, err)

			protoEqual(t, &Uint32{
				OneofTest: &Uint32_OneofOptionA{
					OneofOptionA: 123,
				},
			}, actual)
		}

		func TestOneOfMultipleSet(t *testing.T) {
			t.Setenv("UINT32_ONEOF_A", "123")
			t.Setenv("UINT32_ONEOF_B", "123")

			actual, err := Uint32FromEnv()
			assert.Error(t, err)
			assert.Nil(t, actual)
		}
	`)
}
