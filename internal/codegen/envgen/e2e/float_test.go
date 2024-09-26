package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/envgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestFloatField(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, envgen.GenerateFile, testDataFS, "testdata/float.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"github.com/stretchr/testify/assert"
			"github.com/stretchr/testify/require"
			"testing"
		)

		func TestNormalField(t *testing.T) {
			t.Setenv("FLOAT", "123.456")

			actual, err := FloatFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Float{
				Normal: 123.456,
			}, actual)
		}

		func TestPresenceField(t *testing.T) {
			t.Setenv("FLOAT_WITH_PRESENCE", "123.456")

			actual, err := FloatFromEnv()
			require.NoError(t, err)

			expectedValue := float32(123.456)
			protoEqual(t, &Float{
				WithPresence: &expectedValue,
			}, actual)
		}

		func TestList(t *testing.T) {
			t.Setenv("FLOAT_LIST_1", "123.456")
			t.Setenv("FLOAT_LIST_2", "-123")
			t.Setenv("FLOAT_LIST_3", "987654321.19")

			actual, err := FloatFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Float{
				List: []float32{123.456, -123, 987654321.19},
			}, actual)
		}

		func TestOneOfOneOptionSet(t *testing.T) {
			t.Setenv("FLOAT_ONEOF_A", "123")

			actual, err := FloatFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Float{
				OneofTest: &Float_OneofOptionA{
					OneofOptionA: 123,
				},
			}, actual)
		}

		func TestOneOfMultipleSet(t *testing.T) {
			t.Setenv("FLOAT_ONEOF_A", "123")
			t.Setenv("FLOAT_ONEOF_B", "123")

			actual, err := FloatFromEnv()
			assert.Error(t, err)
			assert.Nil(t, actual)
		}
	`)
}
