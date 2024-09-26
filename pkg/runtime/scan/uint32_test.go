package scan_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/marnixbouhuis/confpb/pkg/runtime/scan"
	"github.com/stretchr/testify/assert"
)

func TestUint32(t *testing.T) {
	type args struct {
		envKey      string
		environment map[string]string
	}
	tests := []*struct {
		name              string
		args              args
		expectedResult    uint32
		expectedHasResult bool
		expectedErr       error
	}{
		{
			name: "reading a valid variable from environment",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "123",
				},
			},
			expectedResult:    123,
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "reading an empty value from environment",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "",
				},
			},
			expectedResult:    0,
			expectedHasResult: true,
			expectedErr: fmt.Errorf("failed to process environment variable \"FOO\": invalid value \"\", unable to parse uint32: %w", &strconv.NumError{
				Func: "ParseUint",
				Num:  "",
				Err:  errors.New("invalid syntax"),
			}),
		},
		{
			name: "environment variable missing",
			args: args{
				envKey:      "FOO",
				environment: map[string]string{},
			},
			expectedResult:    0,
			expectedHasResult: false,
			expectedErr:       nil,
		},
		{
			name: "reading an invalid number from environment",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "not-a-number",
				},
			},
			expectedResult:    0,
			expectedHasResult: true,
			expectedErr: fmt.Errorf("failed to process environment variable \"FOO\": invalid value \"not-a-number\", unable to parse uint32: %w", &strconv.NumError{
				Func: "ParseUint",
				Num:  "not-a-number",
				Err:  errors.New("invalid syntax"),
			}),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.args.environment {
				t.Setenv(k, v)
			}

			result, hasResult, err := scan.Uint32(test.args.envKey)
			assert.Equal(t, test.expectedResult, result)
			assert.Equal(t, test.expectedHasResult, hasResult)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
