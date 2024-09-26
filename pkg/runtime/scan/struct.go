package scan

import (
	"encoding/json"
	"fmt"
	"os"

	"google.golang.org/protobuf/types/known/structpb"
)

var _ Scanner[*structpb.Struct] = Struct

func Struct(envKey string) (result *structpb.Struct, hasResult bool, err error) {
	str, hasValue := os.LookupEnv(envKey)
	if !hasValue {
		return nil, false, nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(str), &data); err != nil {
		return nil, true, fmt.Errorf("failed to process environment variable \"%s\": invalid value \"%s\", unable to JSON decode object: %w", envKey, str, err)
	}

	s, err := structpb.NewStruct(data)
	if err != nil {
		return nil, true, fmt.Errorf("failed to process environment variable \"%s\": invalid value \"%s\", unable to convert decoded data into struct: %w", envKey, str, err)
	}

	return s, true, nil
}
