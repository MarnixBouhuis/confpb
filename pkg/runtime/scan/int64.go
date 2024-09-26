package scan

import (
	"fmt"
	"os"
	"strconv"
)

var _ Scanner[int64] = Int64

func Int64(envKey string) (result int64, hasResult bool, err error) {
	str, hasValue := os.LookupEnv(envKey)
	if !hasValue {
		return 0, false, nil
	}

	result, err = strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, true, fmt.Errorf("failed to process environment variable \"%s\": invalid value \"%s\", unable to parse int64: %w", envKey, str, err)
	}

	return result, true, nil
}
