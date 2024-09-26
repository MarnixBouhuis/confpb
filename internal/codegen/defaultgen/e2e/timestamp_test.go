package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/defaultgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestTimestampField(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, defaultgen.GenerateFile, testDataFS, "testdata/timestamp.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"google.golang.org/protobuf/types/known/timestamppb"
			"testing"
			"time"
		)

		func TestDefaults(t *testing.T) {
			t.Parallel()
			actual := TimestampFromDefault()
			protoEqual(t, &Timestamp{
				Normal: timestamppb.New(time.Date(1985, 04, 12, 23, 20, 50, 520000000, time.UTC)),
				WithPresence: timestamppb.New(time.Date(1985, 04, 12, 23, 20, 50, 520000000, time.UTC)),
				List: []*timestamppb.Timestamp{
					timestamppb.New(time.Date(1985, 04, 12, 23, 20, 50, 520000000, time.UTC)),
					timestamppb.New(time.Date(1937, 01, 01, 11, 40, 27, 870000000, time.UTC)),
					timestamppb.New(time.Date(1990, 12, 31, 23, 59, 59, 0, time.UTC)),
				},
				OneofTest: &Timestamp_OneofOption{
					OneofOption: timestamppb.New(time.Date(1985, 04, 12, 23, 20, 50, 520000000, time.UTC)),
				},
				Map: map[string]*timestamppb.Timestamp{
					"key1": timestamppb.New(time.Date(1985, 04, 12, 23, 20, 50, 520000000, time.UTC)),
					"key2": timestamppb.New(time.Date(1937, 01, 01, 11, 40, 27, 870000000, time.UTC)),
				},
			}, actual)
		}
	`)
}
