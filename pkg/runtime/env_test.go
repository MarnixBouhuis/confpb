package runtime_test

import (
	"testing"

	"github.com/marnixbouhuis/confpb/pkg/runtime"
	"github.com/stretchr/testify/assert"
)

func TestHasEnvKeyWithPrefix(t *testing.T) {
	type args struct {
		prefix      string
		environment map[string]string
	}
	tests := []struct {
		name           string
		args           args
		expectedResult bool
	}{
		{
			name: "has exact key",
			args: args{
				prefix: "SOME_PREFIX",
				environment: map[string]string{
					"SOME_PREFIX": "foo",
				},
			},
			expectedResult: true,
		},
		{
			name: "has key with prefix",
			args: args{
				prefix: "SOME_PREFIX",
				environment: map[string]string{
					"SOME_PREFIX_TEST": "foo",
				},
			},
			expectedResult: true,
		},
		{
			name: "doesn't have key",
			args: args{
				prefix: "SOME_PREFIX",
				environment: map[string]string{
					"DIFFERENT_SOME_PREFIX": "foo",
				},
			},
			expectedResult: false,
		},
		{
			name: "empty environment",
			args: args{
				prefix:      "SOME_PREFIX",
				environment: map[string]string{},
			},
			expectedResult: false,
		},
		{
			name: "empty prefix",
			args: args{
				prefix: "",
				environment: map[string]string{
					"SOME_PREFIX": "foo",
				},
			},
			expectedResult: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for k, v := range test.args.environment {
				t.Setenv(k, v)
			}

			result := runtime.HasEnvKeyWithPrefix(test.args.prefix)
			assert.Equal(t, test.expectedResult, result)
		})
	}
}
