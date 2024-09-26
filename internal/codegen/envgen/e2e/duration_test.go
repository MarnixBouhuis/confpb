package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/envgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestDurationField(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, envgen.GenerateFile, testDataFS, "testdata/duration.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (		
			"github.com/stretchr/testify/assert"
			"github.com/stretchr/testify/require"
			"google.golang.org/protobuf/types/known/durationpb"
			"testing"
			"time"
		)

		func TestNormalField(t *testing.T) {
			t.Setenv("DURATION", "10s")

			actual, err := DurationFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Duration{
				Normal: durationpb.New(10 * time.Second),
			}, actual)
		}

		func TestPresenceField(t *testing.T) {
			t.Setenv("DURATION_WITH_PRESENCE", "10s")

			actual, err := DurationFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Duration{
				WithPresence: durationpb.New(10 * time.Second),
			}, actual)
		}

		func TestList(t *testing.T) {
			t.Setenv("DURATION_LIST_1", "10s")
			t.Setenv("DURATION_LIST_2", "1m")
			t.Setenv("DURATION_LIST_3", "1h")

			actual, err := DurationFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Duration{
				List: []*durationpb.Duration{durationpb.New(10 * time.Second), durationpb.New(time.Minute), durationpb.New(time.Hour)},
			}, actual)
		}

		func TestOneOfOneOptionSet(t *testing.T) {
			t.Setenv("DURATION_ONEOF_A", "10s")

			actual, err := DurationFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Duration{
				OneofTest: &Duration_OneofOptionA{
					OneofOptionA: durationpb.New(10 * time.Second),
				},
			}, actual)
		}

		func TestOneOfMultipleSet(t *testing.T) {
			t.Setenv("DURATION_ONEOF_A", "10s")
			t.Setenv("DURATION_ONEOF_B", "10s")

			actual, err := DurationFromEnv()
			assert.Error(t, err)
			assert.Nil(t, actual)
		}
	`)
}
