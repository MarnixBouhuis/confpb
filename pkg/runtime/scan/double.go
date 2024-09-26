package scan

import (
	"fmt"
	"os"
	"strconv"
)

var _ Scanner[float64] = Double

func Double(envKey string) (result float64, hasResult bool, err error) {
	str, hasValue := os.LookupEnv(envKey)
	if !hasValue {
		return 0, false, nil
	}

	result, err = strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, true, fmt.Errorf("failed to process environment variable \"%s\": invalid value \"%s\", unable to parse double: %w", envKey, str, err)
	}

	return result, true, nil
}
