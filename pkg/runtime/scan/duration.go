package scan

import (
	"fmt"
	"os"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
)

var _ Scanner[*durationpb.Duration] = Duration

func Duration(envKey string) (result *durationpb.Duration, hasResult bool, err error) {
	str, hasValue := os.LookupEnv(envKey)
	if !hasValue {
		return nil, false, nil
	}

	duration, err := time.ParseDuration(str)
	if err != nil {
		return nil, true, fmt.Errorf("failed to process environment variable \"%s\": invalid value \"%s\", unable to parse value as duration: %w", envKey, str, err)
	}

	return durationpb.New(duration), true, nil
}
