package scan

import (
	"fmt"
	"os"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ Scanner[*timestamppb.Timestamp] = Timestamp

func Timestamp(envKey string) (result *timestamppb.Timestamp, hasResult bool, err error) {
	str, hasValue := os.LookupEnv(envKey)
	if !hasValue {
		return nil, false, nil
	}

	timestamp, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return nil, true, fmt.Errorf("failed to process environment variable \"%s\": invalid value \"%s\", unable to parse value as RFC3339 time string: %w", envKey, str, err)
	}

	return timestamppb.New(timestamp), true, nil
}
