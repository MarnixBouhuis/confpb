package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/envgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestInt64Field(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, envgen.GenerateFile, testDataFS, "testdata/int64.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"github.com/stretchr/testify/assert"
			"github.com/stretchr/testify/require"
			"testing"
		)

		func TestNormalField(t *testing.T) {
			t.Setenv("INT64", "123")

			actual, err := Int64FromEnv()
			require.NoError(t, err)

			protoEqual(t, &Int64{
				Normal: 123,
			}, actual)
		}

		func TestPresenceField(t *testing.T) {
			t.Setenv("INT64_WITH_PRESENCE", "123")

			actual, err := Int64FromEnv()
			require.NoError(t, err)

			expectedValue := int64(123)
			protoEqual(t, &Int64{
				WithPresence: &expectedValue,
			}, actual)
		}

		func TestList(t *testing.T) {
			t.Setenv("INT64_LIST_1", "123")
			t.Setenv("INT64_LIST_2", "-123")
			t.Setenv("INT64_LIST_3", "987654321")

			actual, err := Int64FromEnv()
			require.NoError(t, err)

			protoEqual(t, &Int64{
				List: []int64{123, -123, 987654321},
			}, actual)
		}

		func TestOneOfOneOptionSet(t *testing.T) {
			t.Setenv("INT64_ONEOF_A", "123")

			actual, err := Int64FromEnv()
			require.NoError(t, err)

			protoEqual(t, &Int64{
				OneofTest: &Int64_OneofOptionA{
					OneofOptionA: 123,
				},
			}, actual)
		}

		func TestOneOfMultipleSet(t *testing.T) {
			t.Setenv("INT64_ONEOF_A", "123")
			t.Setenv("INT64_ONEOF_B", "123")

			actual, err := Int64FromEnv()
			assert.Error(t, err)
			assert.Nil(t, actual)
		}
	`)
}
