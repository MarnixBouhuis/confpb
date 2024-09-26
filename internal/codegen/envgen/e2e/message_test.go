package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/envgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestNestedMessageField(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, envgen.GenerateFile, testDataFS, "testdata/message.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main
	
		import (
			"github.com/stretchr/testify/assert"
			"github.com/stretchr/testify/require"
			"testing"
		)

		func TestNormal(t *testing.T) {
			t.Setenv("NESTED_TEST", "foobar")

			actual, err := NestedFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Nested{
				Message: &Message{
					Test: "foobar",
				},
			}, actual)
		}

		func TestWithPresence(t *testing.T) {
			t.Setenv("NESTED_PRESENCE_TEST", "foobar")

			actual, err := NestedFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Nested{
				WithPresence: &Message{
					Test: "foobar",
				},
			}, actual)
		}

		func TestRepeated(t *testing.T) {
			t.Setenv("NESTED_LIST_1_TEST", "foobar")
			t.Setenv("NESTED_LIST_2_TEST", "second")
			t.Setenv("NESTED_LIST_3_TEST", "third")

			actual, err := NestedFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Nested{
				List: []*Message{
					{Test: "foobar"},
					{Test: "second"},
					{Test: "third"},
				},
			}, actual)
		}

		func TestOneOfOneOptionSet(t *testing.T) {
			t.Setenv("NESTED_ONEOF_A_TEST", "foobar")

			actual, err := NestedFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Nested{
				OneofTest: &Nested_OneofOptionA{
					OneofOptionA: &Message{
						Test: "foobar",
					},
				},
			}, actual)
		}

		func TestOneOfMultipleSet(t *testing.T) {
			t.Setenv("NESTED_ONEOF_A_TEST", "123")
			t.Setenv("NESTED_ONEOF_B_TEST", "123")

			actual, err := NestedFromEnv()
			assert.Error(t, err)
			assert.Nil(t, actual)
		}
	`)
}

func TestRecursiveNestedMessageField(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, envgen.GenerateFile, testDataFS, "testdata/message.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"github.com/stretchr/testify/assert"
			"github.com/stretchr/testify/require"
			"testing"
		)

		func TestNormal(t *testing.T) {
			t.Setenv("TEST", "level_0")
			t.Setenv("NESTED_TEST", "level_1")
			t.Setenv("NESTED_NESTED_TEST", "level_2")

			actual, err := NestedRecursiveFromEnv()
			require.NoError(t, err)

			protoEqual(t, &NestedRecursive{
				Test: "level_0",
				Message: &NestedRecursive{
					Test: "level_1",
					Message: &NestedRecursive{
						Test: "level_2",
					},
				},
			}, actual)
		}

		func TestWithPresence(t *testing.T) {
			t.Setenv("TEST", "level_0")
			t.Setenv("NESTED_PRESENCE_TEST", "level_1")
			t.Setenv("NESTED_PRESENCE_NESTED_PRESENCE_TEST", "level_2")

			actual, err := NestedRecursiveFromEnv()
			require.NoError(t, err)

			protoEqual(t, &NestedRecursive{
				Test: "level_0",
				WithPresence: &NestedRecursive{
					Test: "level_1",
					WithPresence: &NestedRecursive{
						Test: "level_2",
					},
				},
			}, actual)
		}

		func TestRepeated(t *testing.T) {
			t.Setenv("TEST", "level_0")
			t.Setenv("NESTED_LIST_1_TEST", "level_1_item_1")
			t.Setenv("NESTED_LIST_2_TEST", "level_1_item_2")
			t.Setenv("NESTED_LIST_2_NESTED_LIST_1_TEST", "level_1_item_2_item_1")
			t.Setenv("NESTED_LIST_2_NESTED_LIST_2_TEST", "level_1_item_2_item_2")

			actual, err := NestedRecursiveFromEnv()
			require.NoError(t, err)

			protoEqual(t, &NestedRecursive{
				Test: "level_0",
				List: []*NestedRecursive{
					{
						Test: "level_1_item_1",
					},
					{
						Test: "level_1_item_2",
						List: []*NestedRecursive{
							{Test: "level_1_item_2_item_1"},
							{Test: "level_1_item_2_item_2"},
						},
					},
				},
			}, actual)
		}

		func TestOneOfOneOptionSet(t *testing.T) {
			t.Setenv("TEST", "level_0")
			t.Setenv("NESTED_ONEOF_A_TEST", "level_1")
			t.Setenv("NESTED_ONEOF_A_NESTED_ONEOF_A_TEST", "level_2")

			actual, err := NestedRecursiveFromEnv()
			require.NoError(t, err)

			protoEqual(t, &NestedRecursive{
				Test: "level_0",
				OneofTest: &NestedRecursive_OneofOptionA{
					OneofOptionA: &NestedRecursive{
						Test: "level_1",
						OneofTest: &NestedRecursive_OneofOptionA{
							OneofOptionA: &NestedRecursive{
								Test: "level_2",
							},
						},
					},
				},
			}, actual)
		}

		func TestOneOfMultipleSet(t *testing.T) {
			t.Setenv("NESTED_ONEOF_A_TEST", "123")
			t.Setenv("NESTED_ONEOF_B_TEST", "123")

			actual, err := NestedRecursiveFromEnv()
			assert.Error(t, err)
			assert.Nil(t, actual)
		}
	`)
}

func TestEmbeddedNestedMessageField(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, envgen.GenerateFile, testDataFS, "testdata/message.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main
	
		import (
			"github.com/stretchr/testify/assert"
			"github.com/stretchr/testify/require"
			"testing"
		)

		func TestNormal(t *testing.T) {
			t.Setenv("NESTED_TEST", "foobar")

			actual, err := WithEmbeddedFromEnv()
			require.NoError(t, err)

			protoEqual(t, &WithEmbedded{
				Message: &WithEmbedded_EmbeddedMessage{
					Test: "foobar",
				},
			}, actual)
		}

		func TestWithPresence(t *testing.T) {
			t.Setenv("NESTED_PRESENCE_TEST", "foobar")

			actual, err := WithEmbeddedFromEnv()
			require.NoError(t, err)

			protoEqual(t, &WithEmbedded{
				WithPresence: &WithEmbedded_EmbeddedMessage{
					Test: "foobar",
				},
			}, actual)
		}

		func TestRepeated(t *testing.T) {
			t.Setenv("NESTED_LIST_1_TEST", "foobar")
			t.Setenv("NESTED_LIST_2_TEST", "second")
			t.Setenv("NESTED_LIST_3_TEST", "third")

			actual, err := WithEmbeddedFromEnv()
			require.NoError(t, err)

			protoEqual(t, &WithEmbedded{
				List: []*WithEmbedded_EmbeddedMessage{
					{Test: "foobar"},
					{Test: "second"},
					{Test: "third"},
				},
			}, actual)
		}

		func TestOneOfOneOptionSet(t *testing.T) {
			t.Setenv("NESTED_ONEOF_A_TEST", "foobar")

			actual, err := WithEmbeddedFromEnv()
			require.NoError(t, err)

			protoEqual(t, &WithEmbedded{
				OneofTest: &WithEmbedded_OneofOptionA{
					OneofOptionA: &WithEmbedded_EmbeddedMessage{
						Test: "foobar",
					},
				},
			}, actual)
		}

		func TestOneOfMultipleSet(t *testing.T) {
			t.Setenv("NESTED_ONEOF_A_TEST", "123")
			t.Setenv("NESTED_ONEOF_B_TEST", "123")

			actual, err := WithEmbeddedFromEnv()
			assert.Error(t, err)
			assert.Nil(t, actual)
		}
	`)
}
