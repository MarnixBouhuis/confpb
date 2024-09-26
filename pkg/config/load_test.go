package config_test

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
	"github.com/marnixbouhuis/confpb/pkg/config"
	"github.com/marnixbouhuis/confpb/pkg/config/internal/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

//go:embed internal/testdata
var testDataFS embed.FS

func protoEqual(t *testing.T, expected proto.Message, actual proto.Message) {
	t.Helper()
	if diff := cmp.Diff(expected, actual, protocmp.Transform()); diff != "" {
		t.Errorf("Different protobuf messages (- expected, + actual):\n%s", diff)
	}
}

func TestFromFile(t *testing.T) {
	t.Parallel()

	// Create a temp direrctory containing all test data files
	tempDir, err := os.MkdirTemp("", "confpb_test_load_")
	require.NoError(t, err)
	t.Cleanup(func() {
		t.Logf("Deleting temp directory: %s", tempDir)
		assert.NoError(t, os.RemoveAll(tempDir))
	})
	testutil.CopyAllFromFS(t, testDataFS, tempDir)

	t.Run("Valid config file", func(t *testing.T) {
		t.Parallel()

		files := []string{
			"internal/testdata/config.yaml",
			"internal/testdata/config.yml",
			"internal/testdata/config.json",
			"internal/testdata/config.pb",
			"internal/testdata/config.pb_text",
		}

		for _, file := range files {
			t.Run("Load file: "+file, func(t *testing.T) {
				t.Parallel()

				path := filepath.Join(tempDir, file)
				actual, err := config.FromFile[testdata.TestMessage](path)
				require.NoError(t, err)
				protoEqual(t, &testdata.TestMessage{
					ScalarField:   "foobar",
					RepeatedField: []string{"item1", "item2", "item3"},
					MapField: map[string]string{
						"key1": "value1",
						"key2": "value2",
						"key3": "value3",
					},
					NestedField: &testdata.TestMessage{
						ScalarField: "nested",
					},
					BytesField: []byte("some-bytes"),
				}, actual)
			})
		}
	})

	t.Run("Unknown file extension", func(t *testing.T) {
		t.Parallel()
		path := filepath.Join(tempDir, "unknown.foo")
		actual, err := config.FromFile[testdata.TestMessage](path)
		assert.Nil(t, actual)
		require.Error(t, err)
	})

	t.Run("Missing file (valid extension)", func(t *testing.T) {
		t.Parallel()
		path := filepath.Join(tempDir, "missing.yaml")
		actual, err := config.FromFile[testdata.TestMessage](path)
		assert.Nil(t, actual)
		require.Error(t, err)
	})

	t.Run("Malformed config file", func(t *testing.T) {
		t.Parallel()

		files := []string{
			"internal/testdata/invalid.yaml",
			"internal/testdata/invalid.yml",
			"internal/testdata/invalid.json",
			"internal/testdata/invalid.pb",
			"internal/testdata/invalid.pb_text",
		}

		for _, file := range files {
			t.Run("Load file: "+file, func(t *testing.T) {
				t.Parallel()

				path := filepath.Join(tempDir, file)
				actual, err := config.FromFile[testdata.TestMessage](path)
				require.Error(t, err)
				assert.Nil(t, actual)
			})
		}
	})
}

func TestFromFileFS(t *testing.T) {
	t.Parallel()

	t.Run("Valid config file", func(t *testing.T) {
		t.Parallel()

		files := []string{
			"internal/testdata/config.yaml",
			"internal/testdata/config.yml",
			"internal/testdata/config.json",
			"internal/testdata/config.pb",
			"internal/testdata/config.pb_text",
		}

		for _, file := range files {
			t.Run("Load file: "+file, func(t *testing.T) {
				t.Parallel()

				actual, err := config.FromFileFS[testdata.TestMessage](testDataFS, file)
				require.NoError(t, err)
				protoEqual(t, &testdata.TestMessage{
					ScalarField:   "foobar",
					RepeatedField: []string{"item1", "item2", "item3"},
					MapField: map[string]string{
						"key1": "value1",
						"key2": "value2",
						"key3": "value3",
					},
					NestedField: &testdata.TestMessage{
						ScalarField: "nested",
					},
					BytesField: []byte("some-bytes"),
				}, actual)
			})
		}
	})

	t.Run("Unknown file extension", func(t *testing.T) {
		t.Parallel()
		actual, err := config.FromFileFS[testdata.TestMessage](testDataFS, "unknown.foo")
		assert.Nil(t, actual)
		require.Error(t, err)
	})

	t.Run("Missing file (valid extension)", func(t *testing.T) {
		t.Parallel()
		actual, err := config.FromFileFS[testdata.TestMessage](testDataFS, "missing.yaml")
		assert.Nil(t, actual)
		require.Error(t, err)
	})

	t.Run("Malformed config file", func(t *testing.T) {
		t.Parallel()

		files := []string{
			"internal/testdata/invalid.yaml",
			"internal/testdata/invalid.yml",
			"internal/testdata/invalid.json",
			"internal/testdata/invalid.pb",
			"internal/testdata/invalid.pb_text",
		}

		for _, file := range files {
			t.Run("Load file: "+file, func(t *testing.T) {
				t.Parallel()

				actual, err := config.FromFileFS[testdata.TestMessage](testDataFS, file)
				require.Error(t, err)
				assert.Nil(t, actual)
			})
		}
	})
}

func TestFromYAML(t *testing.T) {
	t.Parallel()

	t.Run("Valid config data", func(t *testing.T) {
		t.Parallel()

		b, err := testDataFS.ReadFile("internal/testdata/config.yaml")
		require.NoError(t, err)

		actual, err := config.FromYAML[testdata.TestMessage](b)
		require.NoError(t, err)
		protoEqual(t, &testdata.TestMessage{
			ScalarField:   "foobar",
			RepeatedField: []string{"item1", "item2", "item3"},
			MapField: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			NestedField: &testdata.TestMessage{
				ScalarField: "nested",
			},
			BytesField: []byte("some-bytes"),
		}, actual)
	})

	t.Run("Invalid config data", func(t *testing.T) {
		t.Parallel()

		b, err := testDataFS.ReadFile("internal/testdata/invalid.yaml")
		require.NoError(t, err)

		actual, err := config.FromYAML[testdata.TestMessage](b)
		require.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestFromJSON(t *testing.T) {
	t.Parallel()

	t.Run("Valid config data", func(t *testing.T) {
		t.Parallel()

		b, err := testDataFS.ReadFile("internal/testdata/config.json")
		require.NoError(t, err)

		actual, err := config.FromJSON[testdata.TestMessage](b)
		require.NoError(t, err)
		protoEqual(t, &testdata.TestMessage{
			ScalarField:   "foobar",
			RepeatedField: []string{"item1", "item2", "item3"},
			MapField: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			NestedField: &testdata.TestMessage{
				ScalarField: "nested",
			},
			BytesField: []byte("some-bytes"),
		}, actual)
	})

	t.Run("Invalid config data", func(t *testing.T) {
		t.Parallel()

		b, err := testDataFS.ReadFile("internal/testdata/invalid.json")
		require.NoError(t, err)

		actual, err := config.FromJSON[testdata.TestMessage](b)
		require.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestFromPb(t *testing.T) {
	t.Parallel()

	t.Run("Valid config data", func(t *testing.T) {
		t.Parallel()

		b, err := testDataFS.ReadFile("internal/testdata/config.pb")
		require.NoError(t, err)

		actual, err := config.FromPb[testdata.TestMessage](b)
		require.NoError(t, err)
		protoEqual(t, &testdata.TestMessage{
			ScalarField:   "foobar",
			RepeatedField: []string{"item1", "item2", "item3"},
			MapField: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			NestedField: &testdata.TestMessage{
				ScalarField: "nested",
			},
			BytesField: []byte("some-bytes"),
		}, actual)
	})

	t.Run("Invalid config data", func(t *testing.T) {
		t.Parallel()

		b, err := testDataFS.ReadFile("internal/testdata/invalid.pb")
		require.NoError(t, err)

		actual, err := config.FromPb[testdata.TestMessage](b)
		require.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestFromPbText(t *testing.T) {
	t.Parallel()

	t.Run("Valid config data", func(t *testing.T) {
		t.Parallel()

		b, err := testDataFS.ReadFile("internal/testdata/config.pb_text")
		require.NoError(t, err)

		actual, err := config.FromPbText[testdata.TestMessage](b)
		require.NoError(t, err)
		protoEqual(t, &testdata.TestMessage{
			ScalarField:   "foobar",
			RepeatedField: []string{"item1", "item2", "item3"},
			MapField: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			NestedField: &testdata.TestMessage{
				ScalarField: "nested",
			},
			BytesField: []byte("some-bytes"),
		}, actual)
	})

	t.Run("Invalid config data", func(t *testing.T) {
		t.Parallel()

		b, err := testDataFS.ReadFile("internal/testdata/invalid.pb_text")
		require.NoError(t, err)

		actual, err := config.FromPbText[testdata.TestMessage](b)
		require.Error(t, err)
		assert.Nil(t, actual)
	})
}
