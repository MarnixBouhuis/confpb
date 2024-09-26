package e2e_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/internal/codegen/defaultgen"
	"github.com/marnixbouhuis/confpb/internal/codegen/testutil"
)

func TestValueField(t *testing.T) {
	t.Parallel()

	res := testutil.RunGeneratorForFiles(t, defaultgen.GenerateFile, testDataFS, "testdata/value.proto")
	testutil.RunTestInE2ERunner(t, res, `
		package main

		import (
			"google.golang.org/protobuf/types/known/structpb"
			"testing"
		)

		func TestDefaults(t *testing.T) {
			t.Parallel()
			actual := ValueFromDefault()
			protoEqual(t, &Value{
				Normal: &structpb.Value{
					Kind: &structpb.Value_NumberValue{
						NumberValue: 123,
					},
				},
				WithPresence: &structpb.Value{
					Kind: &structpb.Value_NumberValue{
						NumberValue: 123,
					},
				},
				List: []*structpb.Value{
					&structpb.Value{
						Kind: &structpb.Value_NumberValue{
							NumberValue: 123,
						},
					},
					&structpb.Value{
						Kind: &structpb.Value_NullValue{
							NullValue: structpb.NullValue_NULL_VALUE,
						},
					},
					&structpb.Value{
						Kind: &structpb.Value_StringValue{
							StringValue: "some-string",
						},
					},
				},
				OneofTest: &Value_OneofOption{
					OneofOption: &structpb.Value{
						Kind: &structpb.Value_NumberValue{
							NumberValue: 123,
						},
					},
				},
				Map: map[string]*structpb.Value{
					"key1": &structpb.Value{
						Kind: &structpb.Value_NumberValue{
							NumberValue: 123,
						},
					},
					"key2": &structpb.Value{
						Kind: &structpb.Value_BoolValue{
							BoolValue: true,
						},
					},
					"key3": &structpb.Value{
						Kind: &structpb.Value_BoolValue{
							BoolValue: false,
						},
					},
				},
			}, actual)
		}
	`)
}
