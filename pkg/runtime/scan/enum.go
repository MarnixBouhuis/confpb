package scan

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
)

func NewEnumScanner[TEnum ~int32](enumMap map[string]int32) Scanner[TEnum] {
	return func(envKey string) (result TEnum, hasResult bool, err error) {
		if enumMap == nil {
			// This should not happen, this indicates that the generated code is wrong and does not correctly pass the enum
			// value map to the enum scanner constructor.
			return 0, false, errors.New("codegen error, enumMap is nil")
		}

		str, hasValue := os.LookupEnv(envKey)
		if !hasValue {
			return 0, false, nil
		}

		value, hasKey := enumMap[str]
		if !hasKey {
			supportedValues := make([]string, 0, len(enumMap))
			for k := range enumMap {
				supportedValues = append(supportedValues, k)
			}
			// Sort supported values so the error message is stable between compiles
			slices.Sort(supportedValues)
			return 0, true, fmt.Errorf("failed to process environment variable \"%s\": invalid value \"%s\", value must be one of: [%s]", envKey, str, strings.Join(supportedValues, ", "))
		}

		return TEnum(value), true, nil
	}
}
