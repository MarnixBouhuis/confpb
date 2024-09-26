package scan_test

import (
	"errors"
	"testing"

	"github.com/marnixbouhuis/confpb/pkg/runtime/scan"
	"github.com/stretchr/testify/assert"
)

func TestEnum(t *testing.T) {
	type args struct {
		envKey      string
		enumMap     map[string]int32
		environment map[string]string
	}
	tests := []*struct {
		name              string
		args              args
		expectedResult    int32
		expectedHasResult bool
		expectedErr       error
	}{
		{
			name: "valid enum value",
			args: args{
				envKey: "FOO",
				enumMap: map[string]int32{
					"VALUE_UNSPECIFIED": 0,
					"VALUE_1":           1,
					"VALUE_2":           2,
				},
				environment: map[string]string{
					"FOO": "VALUE_1",
				},
			},
			expectedResult:    1,
			expectedHasResult: true,
			expectedErr:       nil,
		},
		{
			name: "invalid enum value",
			args: args{
				envKey: "FOO",
				enumMap: map[string]int32{
					"VALUE_UNSPECIFIED": 0,
					"VALUE_2":           1,
					"VALUE_3":           2,
				},
				environment: map[string]string{
					"FOO": "NOT_A_VALUE",
				},
			},
			expectedResult:    0,
			expectedHasResult: true,
			expectedErr:       errors.New("failed to process environment variable \"FOO\": invalid value \"NOT_A_VALUE\", value must be one of: [VALUE_2, VALUE_3, VALUE_UNSPECIFIED]"),
		},
		{
			name: "empty environment value",
			args: args{
				envKey: "FOO",
				enumMap: map[string]int32{
					"VALUE_UNSPECIFIED": 0,
					"VALUE_2":           1,
					"VALUE_3":           2,
				},
				environment: map[string]string{
					"FOO": "",
				},
			},
			expectedResult:    0,
			expectedHasResult: true,
			expectedErr:       errors.New("failed to process environment variable \"FOO\": invalid value \"\", value must be one of: [VALUE_2, VALUE_3, VALUE_UNSPECIFIED]"),
		},
		{
			name: "missing environment value",
			args: args{
				envKey: "FOO",
				enumMap: map[string]int32{
					"VALUE_UNSPECIFIED": 0,
					"VALUE_2":           1,
					"VALUE_3":           2,
				},
				environment: map[string]string{},
			},
			expectedResult:    0,
			expectedHasResult: false,
			expectedErr:       nil,
		},
		{
			name: "empty enum map",
			args: args{
				envKey:  "FOO",
				enumMap: map[string]int32{},
				environment: map[string]string{
					"FOO": "VALUE_1",
				},
			},
			expectedResult:    0,
			expectedHasResult: true,
			expectedErr:       errors.New("failed to process environment variable \"FOO\": invalid value \"VALUE_1\", value must be one of: []"),
		},
		{
			name: "nil enum map",
			args: args{
				envKey:  "FOO",
				enumMap: nil,
				environment: map[string]string{
					"FOO": "VALUE_1",
				},
			},
			expectedResult:    0,
			expectedHasResult: false,
			expectedErr:       errors.New("codegen error, enumMap is nil"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.args.environment {
				t.Setenv(k, v)
			}

			enumScanner := scan.NewEnumScanner[int32](test.args.enumMap)
			result, hasResult, err := enumScanner(test.args.envKey)
			assert.Equal(t, test.expectedResult, result)
			assert.Equal(t, test.expectedHasResult, hasResult)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
