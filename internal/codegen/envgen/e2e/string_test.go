package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/envgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestStringField(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, envgen.GenerateFile, testDataFS, "testdata/string.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"github.com/stretchr/testify/assert"
			"github.com/stretchr/testify/require"
			"testing"
		)

		func TestNormalField(t *testing.T) {
			t.Setenv("STRING", "test-data")

			actual, err := StringFromEnv()
			require.NoError(t, err)

			protoEqual(t, &String{
				Normal: "test-data",
			}, actual)
		}

		func TestPresenceField(t *testing.T) {
			t.Setenv("STRING_WITH_PRESENCE", "test-data")

			actual, err := StringFromEnv()
			require.NoError(t, err)

			expectedValue := "test-data"
			protoEqual(t, &String{
				WithPresence: &expectedValue,
			}, actual)
		}

		func TestList(t *testing.T) {
			t.Setenv("STRING_LIST_1", "test-data")
			t.Setenv("STRING_LIST_2", "something-else")
			t.Setenv("STRING_LIST_3", "")

			actual, err := StringFromEnv()
			require.NoError(t, err)

			protoEqual(t, &String{
				List: []string{"test-data", "something-else", ""},
			}, actual)
		}

		func TestOneOfOneOptionSet(t *testing.T) {
			t.Setenv("STRING_ONEOF_A", "some-data")

			actual, err := StringFromEnv()
			require.NoError(t, err)

			protoEqual(t, &String{
				OneofTest: &String_OneofOptionA{
					OneofOptionA: "some-data",
				},
			}, actual)
		}

		func TestOneOfMultipleSet(t *testing.T) {
			t.Setenv("STRING_ONEOF_A", "some-data")
			t.Setenv("STRING_ONEOF_B", "some-data")

			actual, err := StringFromEnv()
			assert.Error(t, err)
			assert.Nil(t, actual)
		}
	`)
}
