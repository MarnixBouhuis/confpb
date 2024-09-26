package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/defaultgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestFloatField(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, defaultgen.GenerateFile, testDataFS, "testdata/float.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"testing"
		)

		func TestDefaults(t *testing.T) {
			t.Parallel()
			actual := FloatFromDefault()

			x := float32(456)
			protoEqual(t, &Float{
				Normal: float32(123),
				WithPresence: &x,
				List: []float32{123, 456, 789},
				OneofTest: &Float_OneofOption{
					OneofOption: 100,
				},
				Map: map[string]float32{
					"key1": 34,
					"key2": 78,
				},
			}, actual)
		}
	`)
}
