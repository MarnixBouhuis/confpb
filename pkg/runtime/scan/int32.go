package scan

import (
	"fmt"
	"os"
	"strconv"
)

var _ Scanner[int32] = Int32

func Int32(envKey string) (result int32, hasResult bool, err error) {
	str, hasValue := os.LookupEnv(envKey)
	if !hasValue {
		return 0, false, nil
	}

	i64, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, true, fmt.Errorf("failed to process environment variable \"%s\": invalid value \"%s\", unable to parse int32: %w", envKey, str, err)
	}

	return int32(i64), true, nil
}
