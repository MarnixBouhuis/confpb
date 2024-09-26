package scan_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/marnixbouhuis/confpb/pkg/runtime/scan"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestValue(t *testing.T) {
	newValue := func(value interface{}) *structpb.Value {
		t.Helper()
		v, err := structpb.NewValue(value)
		require.NoError(t, err)
		return v
	}

	// newJSONSyntaxErr makes a new JSON syntax error.
	// The msg field on the json.SyntaxError struct is private, so we sadly have to use reflection to set it here.
	// Normally this would be wrong, but since this is testing code it should be acceptable.
	newJSONSyntaxErr := func(msg string, offset int64) *json.SyntaxError {
		t.Helper()

		err := &json.SyntaxError{
			Offset: offset,
		}

		v := reflect.ValueOf(err).Elem()
		field := v.FieldByName("msg")
		reflect.
			NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
			Elem().
			SetString(msg)

		return err
	}

	type args struct {
		envKey      string
		environment map[string]string
	}
	tests := []*struct {
		name              string
		args              args
		expectedResult    *structpb.Value
		expectedHasResult bool
		expectedErr       error
	}{
		{
			name: "reading a valid variable from environment (number)",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "123",
				},
			},
			expectedResult:    newValue(123),
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "reading a valid variable from environment (null)",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "null",
				},
			},
			expectedResult:    newValue(nil),
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "reading a valid variable from environment (boolean array)",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "[true, false, true]",
				},
			},
			expectedResult:    newValue([]interface{}{true, false, true}),
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "reading an empty string from environment",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "",
				},
			},
			expectedResult:    nil,
			expectedHasResult: true,
			expectedErr:       fmt.Errorf("failed to process environment variable \"FOO\": invalid value \"\", unable to JSON decode value: %w", newJSONSyntaxErr("unexpected end of JSON input", 0)),
		},
		{
			name: "environment variable missing",
			args: args{
				envKey:      "FOO",
				environment: map[string]string{},
			},
			expectedResult:    nil,
			expectedHasResult: false,
			expectedErr:       nil,
		},
		{
			name: "reading invalid JSON",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "not-json",
				},
			},
			expectedResult:    nil,
			expectedHasResult: true,
			expectedErr:       fmt.Errorf("failed to process environment variable \"FOO\": invalid value \"not-json\", unable to JSON decode value: %w", newJSONSyntaxErr("invalid character 'o' in literal null (expecting 'u')", 2)),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.args.environment {
				t.Setenv(k, v)
			}

			result, hasResult, err := scan.Value(test.args.envKey)
			assert.Equal(t, test.expectedResult, result)
			assert.Equal(t, test.expectedHasResult, hasResult)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
