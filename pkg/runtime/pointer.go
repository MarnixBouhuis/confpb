package runtime

// Pointer is a utility function that returns v as a pointer to v.
func Pointer[T any](v T) *T {
	return &v
}
