package runtime

import (
	"os"
	"strings"
)

// HasEnvKeyWithPrefix returns if there exists any environment variables where the key starts with "prefix".
func HasEnvKeyWithPrefix(prefix string) bool {
	vars := os.Environ()
	for _, v := range vars {
		kv := strings.SplitN(v, "=", 2)
		if strings.HasPrefix(kv[0], prefix) {
			return true
		}
	}
	return false
}
