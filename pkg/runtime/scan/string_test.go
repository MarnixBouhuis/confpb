package scan_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/pkg/runtime/scan"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	type args struct {
		envKey      string
		environment map[string]string
	}
	tests := []*struct {
		name              string
		args              args
		expectedResult    string
		expectedHasResult bool
		expectedErr       error
	}{
		{
			name: "reading a valid variable from environment",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "bar",
				},
			},
			expectedResult:    "bar",
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
			expectedResult:    "",
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "environment variable missing",
			args: args{
				envKey:      "FOO",
				environment: map[string]string{},
			},
			expectedResult:    "",
			expectedHasResult: false,
			expectedErr:       nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.args.environment {
				t.Setenv(k, v)
			}

			result, hasResult, err := scan.String(test.args.envKey)
			assert.Equal(t, test.expectedResult, result)
			assert.Equal(t, test.expectedHasResult, hasResult)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
