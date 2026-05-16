package validator

// RuleFunc is the signature for a custom validation rule.
type RuleFunc func(field string, value any, param string) error

var customRules = map[string]RuleFunc{}

// Extend registers a named custom validation rule globally.
func Extend(name string, fn RuleFunc) {
	customRules[name] = fn
}
