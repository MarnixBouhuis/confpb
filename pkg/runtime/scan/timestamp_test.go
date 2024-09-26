package scan_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/marnixbouhuis/confpb/pkg/runtime/scan"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTimestamp(t *testing.T) {
	type args struct {
		envKey      string
		environment map[string]string
	}
	tests := []*struct {
		name              string
		args              args
		expectedResult    *timestamppb.Timestamp
		expectedHasResult bool
		expectedErr       error
	}{
		{
			name: "reading a valid variable from environment formatted according to RFC3339",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "1937-01-01T12:00:27.87+00:20",
				},
			},
			expectedResult:    timestamppb.New(time.Date(1937, 1, 1, 11, 40, 27, 870000000, time.UTC)),
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
			expectedErr: fmt.Errorf("failed to process environment variable \"FOO\": invalid value \"\", unable to parse value as RFC3339 time string: %w", &time.ParseError{
				Layout:     time.RFC3339,
				Value:      "",
				LayoutElem: "2006",
				ValueElem:  "",
				Message:    "",
			}),
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
			name: "reading an invalid time",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "invalid-time",
				},
			},
			expectedResult:    nil,
			expectedHasResult: true,
			expectedErr: fmt.Errorf("failed to process environment variable \"FOO\": invalid value \"invalid-time\", unable to parse value as RFC3339 time string: %w", &time.ParseError{
				Layout:     time.RFC3339,
				Value:      "invalid-time",
				LayoutElem: "2006",
				ValueElem:  "invalid-time",
				Message:    "",
			}),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.args.environment {
				t.Setenv(k, v)
			}

			result, hasResult, err := scan.Timestamp(test.args.envKey)
			assert.Equal(t, test.expectedResult, result)
			assert.Equal(t, test.expectedHasResult, hasResult)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
