package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/defaultgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestBytesField(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, defaultgen.GenerateFile, testDataFS, "testdata/bytes.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"testing"
		)

		func TestDefaults(t *testing.T) {
			t.Parallel()
			actual := BytesFromDefault()
			protoEqual(t, &Bytes{
				Normal: []byte("some-bytes"),
				WithPresence: []byte("some-bytes"),
				List: [][]byte{[]byte("some-bytes"), []byte("other-bytes"), {}},
				OneofTest: &Bytes_OneofOption{
					OneofOption: []byte("some-bytes"),
				},
				Map: map[string][]byte{
					"key1": []byte("some-bytes"),
					"key2": []byte("other-bytes"),
				},
			}, actual)
		}
	`)
}
