package scan_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/marnixbouhuis/confpb/pkg/runtime/scan"
	"github.com/stretchr/testify/assert"
)

func TestBytes(t *testing.T) {
	type args struct {
		envKey      string
		environment map[string]string
	}
	tests := []*struct {
		name              string
		args              args
		expectedResult    []byte
		expectedHasResult bool
		expectedErr       error
	}{
		{
			name: "reading a valid variable from environment",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "c2FtcGxlLWRhdGE=",
				},
			},
			expectedResult:    []byte("sample-data"),
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
			expectedResult:    []byte{},
			expectedHasResult: true,
			expectedErr:       nil,
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
			name: "reading invalid base64 data from environment",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "!!!!!!!!!",
				},
			},
			expectedResult:    nil,
			expectedHasResult: true,
			expectedErr:       fmt.Errorf("failed to process environment variable \"FOO\": invalid value \"!!!!!!!!!\", unable to base64 decode value: %w", base64.CorruptInputError(0)),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.args.environment {
				t.Setenv(k, v)
			}

			result, hasResult, err := scan.Bytes(test.args.envKey)
			assert.Equal(t, test.expectedResult, result)
			assert.Equal(t, test.expectedHasResult, hasResult)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
