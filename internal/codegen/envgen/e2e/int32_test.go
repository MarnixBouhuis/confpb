package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/envgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestInt32Field(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, envgen.GenerateFile, testDataFS, "testdata/int32.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"github.com/stretchr/testify/assert"
			"github.com/stretchr/testify/require"
			"testing"
		)

		func TestNormalField(t *testing.T) {
			t.Setenv("INT32", "123")

			actual, err := Int32FromEnv()
			require.NoError(t, err)

			protoEqual(t, &Int32{
				Normal: 123,
			}, actual)
		}

		func TestPresenceField(t *testing.T) {
			t.Setenv("INT32_WITH_PRESENCE", "123")

			actual, err := Int32FromEnv()
			require.NoError(t, err)

			expectedValue := int32(123)
			protoEqual(t, &Int32{
				WithPresence: &expectedValue,
			}, actual)
		}

		func TestList(t *testing.T) {
			t.Setenv("INT32_LIST_1", "123")
			t.Setenv("INT32_LIST_2", "-123")
			t.Setenv("INT32_LIST_3", "987654321")

			actual, err := Int32FromEnv()
			require.NoError(t, err)

			protoEqual(t, &Int32{
				List: []int32{123, -123, 987654321},
			}, actual)
		}

		func TestOneOfOneOptionSet(t *testing.T) {
			t.Setenv("INT32_ONEOF_A", "123")

			actual, err := Int32FromEnv()
			require.NoError(t, err)

			protoEqual(t, &Int32{
				OneofTest: &Int32_OneofOptionA{
					OneofOptionA: 123,
				},
			}, actual)
		}

		func TestOneOfMultipleSet(t *testing.T) {
			t.Setenv("INT32_ONEOF_A", "123")
			t.Setenv("INT32_ONEOF_B", "123")

			actual, err := Int32FromEnv()
			assert.Error(t, err)
			assert.Nil(t, actual)
		}
	`)
}
