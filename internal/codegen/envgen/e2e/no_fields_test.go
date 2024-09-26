package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/envgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestNoFieldsWithOption(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, envgen.GenerateFile, testDataFS, "testdata/no_fields.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"github.com/stretchr/testify/require"
			"testing"
		)

		func TestNoFields(t *testing.T) {
			actual, err := NoFieldsWithEnvOptionFromEnv()
			require.NoError(t, err)

			protoEqual(t, &NoFieldsWithEnvOption{}, actual)
		}
	`)
}
