package scan

import (
	"fmt"
	"os"
	"strconv"
)

var _ Scanner[float32] = Float

func Float(envKey string) (result float32, hasResult bool, err error) {
	str, hasValue := os.LookupEnv(envKey)
	if !hasValue {
		return 0, false, nil
	}

	f64, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return 0, true, fmt.Errorf("failed to process environment variable \"%s\": invalid value \"%s\", unable to parse float: %w", envKey, str, err)
	}

	return float32(f64), true, nil
}
