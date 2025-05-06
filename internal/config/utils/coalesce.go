package utils

// Coalesce returns the first non-nil value from the provided pointers.
func Coalesce[T any](primary, fallback *T) *T {
	if primary != nil {
		return primary
	}
	return fallback
}
