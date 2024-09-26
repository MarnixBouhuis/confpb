package scan

import (
	"fmt"
	"os"
	"strings"
)

var _ Scanner[bool] = Bool

func Bool(envKey string) (result bool, hasResult bool, err error) {
	str, hasValue := os.LookupEnv(envKey)
	if !hasValue {
		return false, false, nil
	}

	switch strings.ToLower(str) {
	case "true", "1", "yes", "y":
		return true, true, nil
	case "false", "0", "no", "n":
		return false, true, nil
	default:
		return false, true, fmt.Errorf("failed to process environment variable \"%s\": unknown value \"%s\", expected (true, 1, yes, y, false, 0, no, n)", envKey, str)
	}
}
