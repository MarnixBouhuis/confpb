package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/defaultgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestMessageField(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, defaultgen.GenerateFile, testDataFS, "testdata/message.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"testing"
		)

		func TestDefaults(t *testing.T) {
			t.Parallel()
			actual := NestedFromDefault()

			protoEqual(t, &Nested{
				Normal: &Message{Test: "foo"},
				DontFill: &Message{},
				WithPresence: &Message{Test: "foo"},
				List: []*Message{{
					Test: "foo",
				}, {
					// No defaults filled
				}, {
					Test: "foo",
				}},
				OneofTest: &Nested_OneofOption{
					OneofOption: &Message{Test: "foo"},
				},
				Map: map[string]*Message{
					"key1": &Message{Test: "foo"},
					"key2": &Message{},
					"key3": &Message{Test: "foo"},
				},
			}, actual)
		}

		func TestEmbeddedMessageDefaults(t *testing.T) {
			t.Parallel()
			actual := WithEmbeddedFromDefault()

			protoEqual(t, &WithEmbedded{
				Normal: &WithEmbedded_EmbeddedMessage{Test: "foo"},
				DontFill: &WithEmbedded_EmbeddedMessage{},
				WithPresence: &WithEmbedded_EmbeddedMessage{Test: "foo"},
				List: []*WithEmbedded_EmbeddedMessage{{
					Test: "foo",
				}, {
					// No defaults filled
				}, {
					Test: "foo",
				}},
				OneofTest: &WithEmbedded_OneofOption{
					OneofOption: &WithEmbedded_EmbeddedMessage{Test: "foo"},
				},
				Map: map[string]*WithEmbedded_EmbeddedMessage{
					"key1": {Test: "foo"},
					"key2": {},
					"key3": {Test: "foo"},
				},
			}, actual)
		}
	`)
}
