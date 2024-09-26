package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/defaultgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestInt64Field(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, defaultgen.GenerateFile, testDataFS, "testdata/int64.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"testing"
		)

		func TestDefaults(t *testing.T) {
			t.Parallel()
			actual := Int64FromDefault()

			x := int64(456)
			protoEqual(t, &Int64{
				Normal: int64(123),
				WithPresence: &x,
				List: []int64{123, 456, 789},
				OneofTest: &Int64_OneofOption{
					OneofOption: 100,
				},
				Map: map[int64]int64{
					12: 34,
					56: 78,
				},
			}, actual)
		}
	`)
}
