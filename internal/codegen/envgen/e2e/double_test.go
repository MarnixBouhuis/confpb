package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/envgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestDoubleField(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, envgen.GenerateFile, testDataFS, "testdata/double.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"github.com/stretchr/testify/assert"
			"github.com/stretchr/testify/require"
			"testing"
		)

		func TestNormalField(t *testing.T) {
			t.Setenv("DOUBLE", "123.456")

			actual, err := DoubleFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Double{
				Normal: 123.456,
			}, actual)
		}

		func TestPresenceField(t *testing.T) {
			t.Setenv("DOUBLE_WITH_PRESENCE", "123.456")

			actual, err := DoubleFromEnv()
			require.NoError(t, err)

			expectedValue := float64(123.456)
			protoEqual(t, &Double{
				WithPresence: &expectedValue,
			}, actual)
		}

		func TestList(t *testing.T) {
			t.Setenv("DOUBLE_LIST_1", "123.456")
			t.Setenv("DOUBLE_LIST_2", "-123")
			t.Setenv("DOUBLE_LIST_3", "987654321.19")

			actual, err := DoubleFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Double{
				List: []float64{123.456, -123, 987654321.19},
			}, actual)
		}

		func TestOneOfOneOptionSet(t *testing.T) {
			t.Setenv("DOUBLE_ONEOF_A", "123")

			actual, err := DoubleFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Double{
				OneofTest: &Double_OneofOptionA{
					OneofOptionA: 123,
				},
			}, actual)
		}

		func TestOneOfMultipleSet(t *testing.T) {
			t.Setenv("DOUBLE_ONEOF_A", "123")
			t.Setenv("DOUBLE_ONEOF_B", "123")

			actual, err := DoubleFromEnv()
			assert.Error(t, err)
			assert.Nil(t, actual)
		}
	`)
}
