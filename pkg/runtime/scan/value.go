package scan

import (
	"encoding/json"
	"fmt"
	"os"

	"google.golang.org/protobuf/types/known/structpb"
)

var _ Scanner[*structpb.Value] = Value

func Value(envKey string) (result *structpb.Value, hasResult bool, err error) {
	str, hasValue := os.LookupEnv(envKey)
	if !hasValue {
		return nil, false, nil
	}

	var data interface{}
	if err := json.Unmarshal([]byte(str), &data); err != nil {
		return nil, true, fmt.Errorf("failed to process environment variable \"%s\": invalid value \"%s\", unable to JSON decode value: %w", envKey, str, err)
	}

	v, err := structpb.NewValue(data)
	if err != nil {
		return nil, true, fmt.Errorf("failed to process environment variable \"%s\": invalid value \"%s\", unable to convert decoded data into value: %w", envKey, str, err)
	}

	return v, true, nil
}
