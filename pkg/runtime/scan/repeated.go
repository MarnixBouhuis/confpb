package scan

import (
	"fmt"
	"strconv"
)

func Repeated[T any](envKeyPrefix string, scanner Scanner[T]) (result []T, hasResult bool, err error) {
	results := make([]T, 0)
	for i := 1; ; i++ {
		r, present, scanErr := scanner(envKeyPrefix + "_" + strconv.Itoa(i))
		if scanErr != nil {
			fieldsPresent := present || len(result) > 0
			return nil, fieldsPresent, fmt.Errorf("failed to process environment variable group \"%s\": %w", envKeyPrefix, scanErr)
		}

		if !present {
			// No more fields
			break
		}

		results = append(results, r)
	}

	if len(results) == 0 {
		return nil, false, nil
	}

	return results, true, nil
}
