package scan

import (
	"encoding/base64"
	"fmt"
	"os"
)

var _ Scanner[[]byte] = Bytes

func Bytes(envKey string) (result []byte, hasResult bool, err error) {
	str, hasValue := os.LookupEnv(envKey)
	if !hasValue {
		return nil, false, nil
	}

	bytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, true, fmt.Errorf("failed to process environment variable \"%s\": invalid value \"%s\", unable to base64 decode value: %w", envKey, str, err)
	}

	return bytes, true, nil
}
