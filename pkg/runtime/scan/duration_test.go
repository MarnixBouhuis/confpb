package scan_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/marnixbouhuis/confpb/pkg/runtime/scan"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TestDuration(t *testing.T) {
	type args struct {
		envKey      string
		environment map[string]string
	}
	tests := []*struct {
		name              string
		args              args
		expectedResult    *durationpb.Duration
		expectedHasResult bool
		expectedErr       error
	}{
		{
			name: "reading a valid variable from environment (1h30s200ms)",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "1h30s200ms",
				},
			},
			expectedResult:    durationpb.New(time.Hour + (30 * time.Second) + (200 * time.Millisecond)),
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "reading a valid variable from environment (180s)",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "180s",
				},
			},
			expectedResult:    durationpb.New(180 * time.Second),
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "reading a valid variable from environment (negative value: -180s)",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "-180s",
				},
			},
			expectedResult:    durationpb.New(-180 * time.Second),
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "reading a valid variable from environment (zero value: 0)",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "0",
				},
			},
			expectedResult:    durationpb.New(0),
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "missing environment variable",
			args: args{
				envKey:      "FOO",
				environment: map[string]string{},
			},
			expectedResult:    nil,
			expectedHasResult: false,
			expectedErr:       nil,
		},
		{
			name: "empty environment variable",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "",
				},
			},
			expectedResult:    nil,
			expectedHasResult: true,
			expectedErr:       fmt.Errorf("failed to process environment variable \"FOO\": invalid value \"\", unable to parse value as duration: %w", errors.New("time: invalid duration \"\"")),
		},
		{
			name: "invalid duration",
			args: args{
				envKey: "FOO",
				environment: map[string]string{
					"FOO": "this-is-not-a-duration",
				},
			},
			expectedResult:    nil,
			expectedHasResult: true,
			expectedErr:       fmt.Errorf("failed to process environment variable \"FOO\": invalid value \"this-is-not-a-duration\", unable to parse value as duration: %w", errors.New("time: invalid duration \"this-is-not-a-duration\"")),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.args.environment {
				t.Setenv(k, v)
			}

			result, hasResult, err := scan.Duration(test.args.envKey)
			assert.Equal(t, test.expectedResult, result)
			assert.Equal(t, test.expectedHasResult, hasResult)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
