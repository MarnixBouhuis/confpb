package scan

import (
	"os"
)

var _ Scanner[string] = String

func String(envKey string) (result string, hasResult bool, err error) {
	str, hasValue := os.LookupEnv(envKey)
	return str, hasValue, nil
}
