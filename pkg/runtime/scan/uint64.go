package scan

import (
	"fmt"
	"os"
	"strconv"
)

var _ Scanner[uint64] = Uint64

func Uint64(envKey string) (result uint64, hasResult bool, err error) {
	str, hasValue := os.LookupEnv(envKey)
	if !hasValue {
		return 0, false, nil
	}

	result, err = strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, true, fmt.Errorf("failed to process environment variable \"%s\": invalid value \"%s\", unable to parse uint64: %w", envKey, str, err)
	}

	return result, true, nil
}
