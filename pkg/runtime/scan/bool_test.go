package scan_test

import (
	"errors"
	"testing"

	"github.com/marnixbouhuis/confpb/pkg/runtime/scan"
	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	type args struct {
		envKey      string
		environment map[string]string
	}
	tests := []*struct {
		name              string
		args              args
		expectedResult    bool
		expectedHasResult bool
		expectedErr       error
	}{
		{
			name: "reading a valid variable from environment (TRUE)",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "TRUE",
				},
			},
			expectedResult:    true,
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "reading a valid variable from environment (FALSE)",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "FALSE",
				},
			},
			expectedResult:    false,
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "reading a valid variable from environment (1)",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "1",
				},
			},
			expectedResult:    true,
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "reading a valid variable from environment (0)",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "0",
				},
			},
			expectedResult:    false,
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "reading a valid variable from environment (yes)",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "yes",
				},
			},
			expectedResult:    true,
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "reading a valid variable from environment (no)",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "no",
				},
			},
			expectedResult:    false,
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "reading a valid variable from environment (y)",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "y",
				},
			},
			expectedResult:    true,
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "reading a valid variable from environment (n)",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "n",
				},
			},
			expectedResult:    false,
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "environment variable missing",
			args: args{
				envKey:      "FOO",
				environment: map[string]string{},
			},
			expectedResult:    false,
			expectedHasResult: false,
			expectedErr:       nil,
		},
		{
			name: "should ignore case of environment variable value",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "tRuE",
				},
			},
			expectedResult:    true,
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "should error on invalid value",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "this-is-not-a-valid-bool",
				},
			},
			expectedResult:    false,
			expectedHasResult: true,
			expectedErr:       errors.New("failed to process environment variable \"FOO\": unknown value \"this-is-not-a-valid-bool\", expected (true, 1, yes, y, false, 0, no, n)"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.args.environment {
				t.Setenv(k, v)
			}

			result, hasResult, err := scan.Bool(test.args.envKey)
			assert.Equal(t, test.expectedResult, result)
			assert.Equal(t, test.expectedHasResult, hasResult)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
