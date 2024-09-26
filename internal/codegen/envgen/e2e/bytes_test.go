package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/envgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestBytesField(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, envgen.GenerateFile, testDataFS, "testdata/bytes.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"github.com/stretchr/testify/assert"
			"github.com/stretchr/testify/require"
			"testing"
		)

		func TestNormalField(t *testing.T) {
			t.Setenv("BYTES", "c29tZS1ieXRlcw==")

			actual, err := BytesFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Bytes{
				Normal: []byte("some-bytes"),
			}, actual)
		}

		func TestPresenceField(t *testing.T) {
			t.Setenv("BYTES_WITH_PRESENCE", "c29tZS1ieXRlcw==")

			actual, err := BytesFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Bytes{
				WithPresence: []byte("some-bytes"),
			}, actual)
		}

		func TestList(t *testing.T) {
			t.Setenv("BYTES_LIST_1", "c29tZS1ieXRlcw==")
			t.Setenv("BYTES_LIST_2", "b3RoZXItYnl0ZXM=")
			t.Setenv("BYTES_LIST_3", "")

			actual, err := BytesFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Bytes{
				List: [][]byte{[]byte("some-bytes"), []byte("other-bytes"), {}},
			}, actual)
		}

		func TestOneOfOneOptionSet(t *testing.T) {
			t.Setenv("BYTES_ONEOF_A", "c29tZS1ieXRlcw==")

			actual, err := BytesFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Bytes{
				OneofTest: &Bytes_OneofOptionA{
					OneofOptionA: []byte("some-bytes"),
				},
			}, actual)
		}

		func TestOneOfMultipleSet(t *testing.T) {
			t.Setenv("BYTES_ONEOF_A", "c29tZS1ieXRlcw==")
			t.Setenv("BYTES_ONEOF_B", "c29tZS1ieXRlcw==")

			actual, err := BytesFromEnv()
			assert.Error(t, err)
			assert.Nil(t, actual)
		}
	`)
}
