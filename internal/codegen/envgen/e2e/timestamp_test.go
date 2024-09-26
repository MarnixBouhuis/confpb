package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/envgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestTimestampField(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, envgen.GenerateFile, testDataFS, "testdata/timestamp.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"github.com/stretchr/testify/assert"
			"github.com/stretchr/testify/require"
			"google.golang.org/protobuf/types/known/timestamppb"
			"testing"
			"time"
		)

		func TestNormalField(t *testing.T) {
			t.Setenv("TIMESTAMP", "1985-04-12T23:20:50.52Z")

			actual, err := TimestampFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Timestamp{
				Normal: timestamppb.New(time.Date(1985, 04, 12, 23, 20, 50, 520000000, time.UTC)),
			}, actual)
		}

		func TestPresenceField(t *testing.T) {
			t.Setenv("TIMESTAMP_WITH_PRESENCE", "1985-04-12T23:20:50.52Z")

			actual, err := TimestampFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Timestamp{
				WithPresence: timestamppb.New(time.Date(1985, 04, 12, 23, 20, 50, 520000000, time.UTC)),
			}, actual)
		}

		func TestList(t *testing.T) {
			t.Setenv("TIMESTAMP_LIST_1", "1985-04-12T23:20:50.52Z")
			t.Setenv("TIMESTAMP_LIST_2", "1937-01-01T12:00:27.87+00:20")
			t.Setenv("TIMESTAMP_LIST_3", "1990-12-31T23:59:59Z")

			actual, err := TimestampFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Timestamp{
				List: []*timestamppb.Timestamp{
					timestamppb.New(time.Date(1985, 04, 12, 23, 20, 50, 520000000, time.UTC)),
					timestamppb.New(time.Date(1937, 01, 01, 11, 40, 27, 870000000, time.UTC)),
					timestamppb.New(time.Date(1990, 12, 31, 23, 59, 59, 0, time.UTC)),
				},
			}, actual)
		}

		func TestOneOfOneOptionSet(t *testing.T) {
			t.Setenv("TIMESTAMP_ONEOF_A", "1985-04-12T23:20:50.52Z")

			actual, err := TimestampFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Timestamp{
				OneofTest: &Timestamp_OneofOptionA{
					OneofOptionA: timestamppb.New(time.Date(1985, 04, 12, 23, 20, 50, 520000000, time.UTC)),
				},
			}, actual)
		}

		func TestOneOfMultipleSet(t *testing.T) {
			t.Setenv("TIMESTAMP_ONEOF_A", "1985-04-12T23:20:50.52Z")
			t.Setenv("TIMESTAMP_ONEOF_B", "1985-04-12T23:20:50.52Z")

			actual, err := TimestampFromEnv()
			assert.Error(t, err)
			assert.Nil(t, actual)
		}
	`)
}
