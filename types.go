package validator

// Input is the data map to validate.
type Input map[string]any

// Rules maps field names to pipe-separated rule strings.
type Rules map[string]string

// Messages maps "field.rule" or "rule" keys to custom error messages.
type Messages map[string]string

// ErrorBag holds validation errors per field.
type ErrorBag map[string][]string

// First returns the first error message for the given field, or empty string.
func (e ErrorBag) First(field string) string {
	if msgs, ok := e[field]; ok && len(msgs) > 0 {
		return msgs[0]
	}
	return ""
}

// Has reports whether the given field has any errors.
func (e ErrorBag) Has(field string) bool {
	msgs, ok := e[field]
	return ok && len(msgs) > 0
}
