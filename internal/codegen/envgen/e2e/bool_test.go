package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/envgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestBoolField(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, envgen.GenerateFile, testDataFS, "testdata/bool.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"github.com/stretchr/testify/assert"
			"github.com/stretchr/testify/require"
			"testing"
		)

		func TestNormalField(t *testing.T) {
			t.Setenv("BOOL", "true")

			actual, err := BoolFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Bool{
				Normal: true,
			}, actual)
		}

		func TestPresenceField(t *testing.T) {
			t.Setenv("BOOL_WITH_PRESENCE", "true")

			actual, err := BoolFromEnv()
			require.NoError(t, err)

			expectedValue := true
			protoEqual(t, &Bool{
				WithPresence: &expectedValue,
			}, actual)
		}

		func TestList(t *testing.T) {
			t.Setenv("BOOL_LIST_1", "true")
			t.Setenv("BOOL_LIST_2", "false")
			t.Setenv("BOOL_LIST_3", "no")

			actual, err := BoolFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Bool{
				List: []bool{true, false, false},
			}, actual)
		}

		func TestOneOfOneOptionSet(t *testing.T) {
			t.Setenv("BOOL_ONEOF_A", "false")

			actual, err := BoolFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Bool{
				OneofTest: &Bool_OneofOptionA{
					OneofOptionA: false,
				},
			}, actual)
		}

		func TestOneOfMultipleSet(t *testing.T) {
			t.Setenv("BOOL_ONEOF_A", "false")
			t.Setenv("BOOL_ONEOF_B", "false")

			actual, err := BoolFromEnv()
			assert.Error(t, err)
			assert.Nil(t, actual)
		}
	`)
}
