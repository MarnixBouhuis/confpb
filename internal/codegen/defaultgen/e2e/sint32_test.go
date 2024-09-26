package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/defaultgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestSint32Field(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, defaultgen.GenerateFile, testDataFS, "testdata/sint32.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"testing"
		)

		func TestDefaults(t *testing.T) {
			t.Parallel()
			actual := Sint32FromDefault()

			x := int32(456)
			protoEqual(t, &Sint32{
				Normal: int32(123),
				WithPresence: &x,
				List: []int32{123, 456, 789},
				OneofTest: &Sint32_OneofOption{
					OneofOption: 100,
				},
				Map: map[int32]int32{
					12: 34,
					56: 78,
				},
			}, actual)
		}
	`)
}
