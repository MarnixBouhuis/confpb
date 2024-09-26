package config_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/pkg/config"
	"github.com/marnixbouhuis/confpb/pkg/config/internal/testdata"
	"github.com/stretchr/testify/assert"
)

func TestOverlay(t *testing.T) {
	t.Parallel()

	t.Run("Merging all fields", func(t *testing.T) {
		t.Parallel()

		base := &testdata.TestMessage{
			ScalarField:   "base",
			RepeatedField: []string{"base1", "base2", "base3"},
			MapField: map[string]string{
				"key1": "base1",
				"key2": "base2",
			},
			NestedField: &testdata.TestMessage{
				ScalarField: "base-nested1",
			},
			BytesField: []byte("base-bytes"),
		}

		overlay := &testdata.TestMessage{
			ScalarField:   "overlay",
			RepeatedField: []string{"overlay1", "overlay2", "overlay3"},
			MapField: map[string]string{
				"key1": "overlay1",
				"key2": "overlay2",
			},
			NestedField: &testdata.TestMessage{
				ScalarField: "overlay-nested1",
			},
			BytesField: []byte("overlay-bytes"),
		}

		actual := config.Overlay(base, overlay)
		protoEqual(t, &testdata.TestMessage{
			ScalarField:   "overlay",
			RepeatedField: []string{"overlay1", "overlay2", "overlay3"},
			MapField: map[string]string{
				"key1": "overlay1",
				"key2": "overlay2",
			},
			NestedField: &testdata.TestMessage{
				ScalarField: "overlay-nested1",
			},
			BytesField: []byte("overlay-bytes"),
		}, actual)
	})

	t.Run("Merging some field that are already set in the base", func(t *testing.T) {
		t.Parallel()

		base := &testdata.TestMessage{
			ScalarField:   "base",
			RepeatedField: []string{"base1", "base2", "base3"},
			MapField: map[string]string{
				"key1": "base1",
				"key2": "base2",
			},
			NestedField: &testdata.TestMessage{
				ScalarField: "base-nested1",
			},
			BytesField: []byte("base-bytes"),
		}

		overlay := &testdata.TestMessage{
			NestedField: &testdata.TestMessage{
				ScalarField: "overlay-nested1",
			},
			BytesField: []byte("overlay-bytes"),
		}

		actual := config.Overlay(base, overlay)
		protoEqual(t, &testdata.TestMessage{
			ScalarField:   "base",
			RepeatedField: []string{"base1", "base2", "base3"},
			MapField: map[string]string{
				"key1": "base1",
				"key2": "base2",
			},
			NestedField: &testdata.TestMessage{
				ScalarField: "overlay-nested1",
			},
			BytesField: []byte("overlay-bytes"),
		}, actual)
	})

	t.Run("Merging some field that are not yet set in the base", func(t *testing.T) {
		t.Parallel()

		base := &testdata.TestMessage{
			ScalarField:   "base",
			RepeatedField: []string{"base1", "base2", "base3"},
			MapField: map[string]string{
				"key1": "base1",
				"key2": "base2",
			},
			NestedField: &testdata.TestMessage{
				ScalarField: "base-nested1",
			},
			BytesField: []byte("base-bytes"),
		}

		overlay := &testdata.TestMessage{
			NestedField: &testdata.TestMessage{
				ScalarField: "overlay-nested1",
			},
			BytesField: []byte("overlay-bytes"),
		}

		actual := config.Overlay(base, overlay)
		protoEqual(t, &testdata.TestMessage{
			ScalarField:   "base",
			RepeatedField: []string{"base1", "base2", "base3"},
			MapField: map[string]string{
				"key1": "base1",
				"key2": "base2",
			},
			NestedField: &testdata.TestMessage{
				ScalarField: "overlay-nested1",
			},
			BytesField: []byte("overlay-bytes"),
		}, actual)
	})

	t.Run("Merging pointer fields with a nil overlay should not overwrite the base", func(t *testing.T) {
		t.Parallel()

		base := &testdata.TestMessage{
			NestedField: &testdata.TestMessage{
				ScalarField: "base-nested1",
			},
		}

		overlay := &testdata.TestMessage{
			NestedField: nil,
		}

		actual := config.Overlay(base, overlay)
		protoEqual(t, &testdata.TestMessage{
			NestedField: &testdata.TestMessage{
				ScalarField: "base-nested1",
			},
		}, actual)
	})

	t.Run("Merging pointer fields with a nil base and non nil overlay should overwrite the base", func(t *testing.T) {
		t.Parallel()

		base := &testdata.TestMessage{
			NestedField: nil,
		}

		overlay := &testdata.TestMessage{
			NestedField: &testdata.TestMessage{
				ScalarField: "overlay-nested1",
			},
		}

		actual := config.Overlay(base, overlay)
		protoEqual(t, &testdata.TestMessage{
			NestedField: &testdata.TestMessage{
				ScalarField: "overlay-nested1",
			},
		}, actual)
	})

	t.Run("Merging with a nil base", func(t *testing.T) {
		t.Parallel()

		actual := config.Overlay(nil, &testdata.TestMessage{
			ScalarField: "overlay",
		})
		protoEqual(t, &testdata.TestMessage{
			ScalarField: "overlay",
		}, actual)
	})

	t.Run("Merging with a nil overlay", func(t *testing.T) {
		t.Parallel()

		actual := config.Overlay(&testdata.TestMessage{
			ScalarField: "base",
		}, nil)
		protoEqual(t, &testdata.TestMessage{
			ScalarField: "base",
		}, actual)
	})

	t.Run("Merging with a nil base and overlay", func(t *testing.T) {
		t.Parallel()

		actual := config.Overlay[testdata.TestMessage](nil, nil)
		assert.Nil(t, actual)
	})

	t.Run("Modifying the base after the merge should not change the result struct", func(t *testing.T) {
		t.Parallel()

		base := &testdata.TestMessage{
			NestedField: &testdata.TestMessage{
				ScalarField: "original",
			},
		}
		overlay := &testdata.TestMessage{}
		merged := config.Overlay(base, overlay)

		// Modify the base
		base.NestedField.ScalarField = "modified"

		// Make sure the merged config still has the old base value
		assert.Equal(t, "original", merged.NestedField.ScalarField)
	})

	t.Run("Modifying the overlay after the merge should not change the result struct", func(t *testing.T) {
		t.Parallel()

		base := &testdata.TestMessage{}
		overlay := &testdata.TestMessage{
			NestedField: &testdata.TestMessage{
				ScalarField: "original",
			},
		}
		merged := config.Overlay(base, overlay)

		// Modify the base
		overlay.NestedField.ScalarField = "modified"

		// Make sure the merged config still has the old base value
		assert.Equal(t, "original", merged.NestedField.ScalarField)
	})
}
