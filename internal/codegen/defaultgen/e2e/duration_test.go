package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/defaultgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestDurationField(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, defaultgen.GenerateFile, testDataFS, "testdata/duration.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"google.golang.org/protobuf/types/known/durationpb"
			"testing"
			"time"
		)

		func TestDefaults(t *testing.T) {
			t.Parallel()
			actual := DurationFromDefault()
			protoEqual(t, &Duration{
				Normal: durationpb.New(time.Second*10),
				WithPresence: durationpb.New(time.Second*10),
				List: []*durationpb.Duration{
					durationpb.New(time.Second*10),
					durationpb.New(time.Minute),
					durationpb.New(time.Hour),
				},
				OneofTest: &Duration_OneofOption{
					OneofOption: durationpb.New(time.Second*10),
				},
				Map: map[string]*durationpb.Duration{
					"key1": durationpb.New(time.Second*10),
					"key2": durationpb.New(time.Minute),
				},
			}, actual)
		}
	`)
}
