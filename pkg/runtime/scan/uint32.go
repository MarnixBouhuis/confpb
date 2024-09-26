package scan

import (
	"fmt"
	"os"
	"strconv"
)

var _ Scanner[uint32] = Uint32

func Uint32(envKey string) (result uint32, hasResult bool, err error) {
	str, hasValue := os.LookupEnv(envKey)
	if !hasValue {
		return 0, false, nil
	}

	u64, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, true, fmt.Errorf("failed to process environment variable \"%s\": invalid value \"%s\", unable to parse uint32: %w", envKey, str, err)
	}

	return uint32(u64), true, nil
}
