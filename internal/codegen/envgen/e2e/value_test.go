package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/envgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestValueField(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, envgen.GenerateFile, testDataFS, "testdata/value.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (		
			"github.com/stretchr/testify/assert"
			"github.com/stretchr/testify/require"
			"google.golang.org/protobuf/types/known/structpb"
			"testing"
		)

		func TestNormalField(t *testing.T) {
			t.Setenv("VALUE", "123")

			actual, err := ValueFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Value{
				Normal: &structpb.Value{
					Kind: &structpb.Value_NumberValue{
						NumberValue: 123,
					},
				},
			}, actual)
		}

		func TestPresenceField(t *testing.T) {
			t.Setenv("VALUE_WITH_PRESENCE", "123")

			actual, err := ValueFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Value{
				WithPresence: &structpb.Value{
					Kind: &structpb.Value_NumberValue{
						NumberValue: 123,
					},
				},
			}, actual)
		}

		func TestList(t *testing.T) {
			t.Setenv("VALUE_LIST_1", "123")
			t.Setenv("VALUE_LIST_2", "null")
			t.Setenv("VALUE_LIST_3", "\"some-string\"")

			actual, err := ValueFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Value{
				List: []*structpb.Value{{
					Kind: &structpb.Value_NumberValue{
						NumberValue: 123,
					},
				}, {
					Kind: &structpb.Value_NullValue{
						NullValue: structpb.NullValue_NULL_VALUE,
					},
				}, {
					Kind: &structpb.Value_StringValue{
						StringValue: "some-string",
					},
				}},
			}, actual)
		}

		func TestOneOfOneOptionSet(t *testing.T) {
			t.Setenv("VALUE_ONEOF_A", "123")

			actual, err := ValueFromEnv()
			require.NoError(t, err)

			protoEqual(t, &Value{
				OneofTest: &Value_OneofOptionA{
					OneofOptionA: &structpb.Value{
						Kind: &structpb.Value_NumberValue{
							NumberValue: 123,
						},
					},
				},
			}, actual)
		}

		func TestOneOfMultipleSet(t *testing.T) {
			t.Setenv("VALUE_ONEOF_A", "123")
			t.Setenv("VALUE_ONEOF_B", "123")

			actual, err := ValueFromEnv()
			assert.Error(t, err)
			assert.Nil(t, actual)
		}
	`)
}
