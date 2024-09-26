package scan

// Scanner takes an environment variable key and scans it into a result T.
// It returns (via hasResult) if the environment variable was set. If scanning goes wrong it returns an error.
// Note: if an error is returned the values of "result" and "hasResult" should be ignored.
type Scanner[T any] func(envKey string) (result T, hasResult bool, err error)
