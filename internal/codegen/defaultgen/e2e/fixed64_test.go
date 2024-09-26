package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/defaultgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestFixed64Field(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, defaultgen.GenerateFile, testDataFS, "testdata/fixed64.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"testing"
		)

		func TestDefaults(t *testing.T) {
			t.Parallel()
			actual := Fixed64FromDefault()

			x := uint64(456)
			protoEqual(t, &Fixed64{
				Normal: uint64(123),
				WithPresence: &x,
				List: []uint64{123, 456, 789},
				OneofTest: &Fixed64_OneofOption{
					OneofOption: 100,
				},
				Map: map[uint64]uint64{
					12: 34,
					56: 78,
				},
			}, actual)
		}
	`)
}
