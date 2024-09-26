package scan_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/marnixbouhuis/confpb/pkg/runtime/scan"
	"github.com/stretchr/testify/assert"
)

func TestRepeated(t *testing.T) {
	type args struct {
		envKey      string
		environment map[string]string
	}
	tests := []*struct {
		name              string
		args              args
		expectedResult    []int32
		expectedHasResult bool
		expectedErr       error
	}{
		{
			name: "reading multiple values from the environment",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO_1": "123",
					"FOO_2": "456",
					"FOO_3": "789",
				},
			},
			expectedResult:    []int32{123, 456, 789},
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "first value missing should result in an empty list",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO_2": "456",
					"FOO_3": "789",
				},
			},
			expectedResult:    nil,
			expectedHasResult: false,
			expectedErr:       nil,
		},
		{
			name: "reading multiple values from the environment with a gap in the list should only return up to the gap",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO_1": "123",
					"FOO_2": "456",
					"FOO_4": "012",
				},
			},
			expectedResult:    []int32{123, 456},
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "no environment variables set",
			args: args{
				envKey:      "FOO",
				environment: map[string]string{},
			},
			expectedResult:    nil,
			expectedHasResult: false,
			expectedErr:       nil,
		},
		{
			name: "should error if it encounters an invalid value",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO_1": "123",
					"FOO_2": "456",
					"FOO_3": "this-value-is-invalid",
					"FOO_4": "012",
				},
			},
			expectedResult:    nil,
			expectedHasResult: true,
			expectedErr: fmt.Errorf(
				"failed to process environment variable group \"FOO\": %w",
				fmt.Errorf("failed to process environment variable \"FOO_3\": invalid value \"this-value-is-invalid\", unable to parse int32: %w", &strconv.NumError{
					Func: "ParseInt",
					Num:  "this-value-is-invalid",
					Err:  errors.New("invalid syntax"),
				}),
			),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.args.environment {
				t.Setenv(k, v)
			}

			// For easy testing we use an int32 scanner
			result, hasResult, err := scan.Repeated(test.args.envKey, scan.Int32)
			assert.Equal(t, test.expectedResult, result)
			assert.Equal(t, test.expectedHasResult, hasResult)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
