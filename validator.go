package validator

import (
	"database/sql"
	"fmt"
	"strings"
)

// Validator holds the validation state.
type Validator struct {
	input     Input
	rules     Rules
	customMsg Messages
	errors    ErrorBag
	validated bool
	bailMode  bool
	db        *sql.DB
}

// WithDB sets the database connection used by the unique and exists rules.
func (v *Validator) WithDB(db *sql.DB) *Validator {
	v.db = db
	return v
}

// Make creates a new Validator. Validation runs lazily on first call to Fails or Errors.
func Make(input Input, rules Rules) *Validator {
	return &Validator{
		input:     input,
		rules:     rules,
		customMsg: make(Messages),
		errors:    make(ErrorBag),
	}
}

// Messages sets custom error messages and returns the validator for chaining.
func (v *Validator) Messages(msgs Messages) *Validator {
	v.customMsg = msgs
	return v
}

// Bail stops validation after the first field that has an error.
func (v *Validator) Bail() *Validator {
	v.bailMode = true
	return v
}

// Fails reports whether any validation rule failed.
func (v *Validator) Fails() bool {
	v.run()
	return len(v.errors) > 0
}

// Passes reports whether all validation rules passed.
func (v *Validator) Passes() bool {
	return !v.Fails()
}

// Errors returns all validation errors grouped by field.
func (v *Validator) Errors() ErrorBag {
	v.run()
	return v.errors
}

func (v *Validator) run() {
	if v.validated {
		return
	}
	v.validated = true

	for field, ruleStr := range v.rules {
		// ── Array wildcard (e.g. "tags.*") ────────────────────────────────────
		if strings.HasSuffix(field, ".*") {
			baseField := strings.TrimSuffix(field, ".*")
			arr, ok := getNestedValue(v.input, baseField)
			if !ok {
				continue
			}
			items, ok := arr.([]any)
			if !ok {
				continue
			}
			for i, item := range items {
				elemField := fmt.Sprintf("%s.%d", baseField, i)
				v.validateSingle(elemField, item, ruleStr)
			}
			if v.bailMode && len(v.errors) > 0 {
				return
			}
			continue
		}

		// ── Normal field validation ────────────────────────────────────────────
		val, fieldPresent := getNestedValue(v.input, field)
		v.validateSingle(field, val, ruleStr)
		_ = fieldPresent

		if v.bailMode && len(v.errors) > 0 {
			return
		}
	}
}

// validateSingle runs all rules for a single field/value pair.
func (v *Validator) validateSingle(field string, val any, ruleStr string) {
	parts := parseAndCache(ruleStr)

	_, fieldPresent := getNestedValue(v.input, field)

	hasRequired := false
	isNullable := false
	effectiveRequired := false

	// Parse rule names first pass — detect nullable/required/flow conditions.
	for _, p := range parts {
		switch p.name {
		case "required":
			hasRequired = true
			effectiveRequired = true
		case "nullable":
			isNullable = true
		}
	}

	// ── Presence rules (present, filled) ──────────────────────────────────
	for _, p := range parts {
		if p.name == "present" {
			if !fieldPresent {
				msg := buildMsg("present", field, "")
				v.errors[field] = append(v.errors[field], v.resolve(field, "present", msg))
			}
			continue
		}

		if p.name == "filled" {
			if fieldPresent && isEmpty(val) {
				msg := buildMsg("filled", field, "")
				v.errors[field] = append(v.errors[field], v.resolve(field, "filled", msg))
				return
			}
			continue
		}
	}

	// ── Conditional required rules ────────────────────────────────────────
	for _, p := range parts {
		switch p.name {
		case "required_if":
			ps := splitParam(p.param)
			if len(ps) >= 2 && fmt.Sprintf("%v", v.input[ps[0]]) == ps[1] {
				effectiveRequired = true
			}
		case "required_unless":
			ps := splitParam(p.param)
			if len(ps) >= 2 && fmt.Sprintf("%v", v.input[ps[0]]) != ps[1] {
				effectiveRequired = true
			}
		case "required_with":
			for _, f := range splitParam(p.param) {
				if !isEmpty(v.input[f]) {
					effectiveRequired = true
					break
				}
			}
		case "required_without":
			for _, f := range splitParam(p.param) {
				if isEmpty(v.input[f]) {
					effectiveRequired = true
					break
				}
			}
		case "required_with_all":
			allPresent := true
			for _, f := range splitParam(p.param) {
				if isEmpty(v.input[f]) {
					allPresent = false
					break
				}
			}
			if allPresent {
				effectiveRequired = true
			}
		case "required_without_all":
			allAbsent := true
			for _, f := range splitParam(p.param) {
				if !isEmpty(v.input[f]) {
					allAbsent = false
					break
				}
			}
			if allAbsent {
				effectiveRequired = true
			}
		}
	}

	// ── Handle nullable nil shortcut ──────────────────────────────────────
	if isNullable && val == nil {
		return
	}

	// ── Handle empty value ────────────────────────────────────────────────
	if isEmpty(val) {
		if effectiveRequired {
			ruleName, msg := v.buildRequiredMsg(field, parts)
			v.errors[field] = append(v.errors[field], v.resolve(field, ruleName, msg))
		} else if !hasRequired {
			// Optional field — skip further validation.
		}
		return
	}

	// ── Run check rules ───────────────────────────────────────────────────
	fieldBail := hasBail(parts)
	for _, p := range parts {
		if isFlowRule(p.name) {
			continue
		}
		if msg := v.check(p.name, p.param, field, val); msg != "" {
			v.errors[field] = append(v.errors[field], v.resolve(field, p.name, msg))
			if fieldBail {
				return
			}
			break
		}
	}
}

// hasBail reports whether "bail" appears in a rule set.
func hasBail(parts []parsedRule) bool {
	for _, p := range parts {
		if p.name == "bail" {
			return true
		}
	}
	return false
}

// buildRequiredMsg finds the first triggered required-style rule and returns
// the rule name and its formatted message.
func (v *Validator) buildRequiredMsg(field string, parts []parsedRule) (ruleName, msg string) {
	for _, p := range parts {
		ps := splitParam(p.param)
		switch p.name {
		case "required":
			return "required", buildMsg("required", field, p.param)
		case "required_if":
			if len(ps) >= 2 && fmt.Sprintf("%v", v.input[ps[0]]) == ps[1] {
				return "required_if", buildMsgWith("required_if", map[string]string{
					":field": humanize(field),
					":other": humanize(ps[0]),
					":value": ps[1],
				})
			}
		case "required_unless":
			if len(ps) >= 2 && fmt.Sprintf("%v", v.input[ps[0]]) != ps[1] {
				return "required_unless", buildMsgWith("required_unless", map[string]string{
					":field": humanize(field),
					":other": humanize(ps[0]),
					":value": ps[1],
				})
			}
		case "required_with":
			for _, f := range ps {
				if !isEmpty(v.input[f]) {
					return "required_with", buildMsgWith("required_with", map[string]string{
						":field": humanize(field),
						":other": humanize(f),
					})
				}
			}
		case "required_without":
			for _, f := range ps {
				if isEmpty(v.input[f]) {
					return "required_without", buildMsgWith("required_without", map[string]string{
						":field": humanize(field),
						":other": humanize(f),
					})
				}
			}
		case "required_with_all":
			allPresent := true
			for _, f := range ps {
				if isEmpty(v.input[f]) {
					allPresent = false
					break
				}
			}
			if allPresent {
				return "required_with_all", buildMsgWith("required_with_all", map[string]string{
					":field": humanize(field),
					":other": strings.Join(ps, ", "),
				})
			}
		case "required_without_all":
			allAbsent := true
			for _, f := range ps {
				if !isEmpty(v.input[f]) {
					allAbsent = false
					break
				}
			}
			if allAbsent {
				return "required_without_all", buildMsgWith("required_without_all", map[string]string{
					":field": humanize(field),
					":other": strings.Join(ps, ", "),
				})
			}
		}
	}
	return "required", buildMsg("required", field, "")
}

// isFlowRule returns true for rules that are handled in validateSingle() and should not
// be passed to check().
func isFlowRule(name string) bool {
	switch name {
	case "nullable", "filled", "present", "bail",
		"required_if", "required_unless",
		"required_with", "required_without",
		"required_with_all", "required_without_all":
		return true
	}
	return false
}

func (v *Validator) resolve(field, rule, fallback string) string {
	if msg, ok := v.customMsg[field+"."+rule]; ok {
		return msg
	}
	if msg, ok := v.customMsg[rule]; ok {
		return msg
	}
	return fallback
}

func parseRule(s string) (name, param string) {
	parts := strings.SplitN(s, ":", 2)
	name = strings.TrimSpace(parts[0])
	if len(parts) == 2 {
		param = strings.TrimSpace(parts[1])
	}
	return
}
