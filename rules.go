package validator

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

// Precompiled regex patterns.
var (
	reUUID     = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	reULID     = regexp.MustCompile(`^[0-9A-HJKMNP-TV-Z]{26}$`)
	reHexColor = regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$`)
	reEmail    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

// parsedRule holds a pre-parsed rule name and optional parameter.
type parsedRule struct{ name, param string }

// ruleCache caches the parsed form of rule strings to avoid re-splitting on every validation call.
var ruleCache sync.Map // key: string → []parsedRule

// parseAndCache parses a pipe-separated rule string and caches the result.
func parseAndCache(ruleStr string) []parsedRule {
	if v, ok := ruleCache.Load(ruleStr); ok {
		return v.([]parsedRule)
	}
	parts := strings.Split(ruleStr, "|")
	result := make([]parsedRule, 0, len(parts))
	for _, p := range parts {
		n, param := parseRule(p)
		result = append(result, parsedRule{n, param})
	}
	ruleCache.Store(ruleStr, result)
	return result
}

// regexCache caches compiled user-supplied regex patterns (used by "regex" and "not_regex" rules).
var regexCache sync.Map // key: string → *regexp.Regexp

func compileRegex(pattern string) (*regexp.Regexp, error) {
	if v, ok := regexCache.Load(pattern); ok {
		return v.(*regexp.Regexp), nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	regexCache.Store(pattern, re)
	return re, nil
}

func (v *Validator) check(name, param, field string, value any) string {
	if fn, ok := customRules[name]; ok {
		if err := fn(field, value, param); err != nil {
			return err.Error()
		}
		return ""
	}

	str := fmt.Sprintf("%v", value)

	switch name {
	case "required":
		if isEmpty(value) {
			return buildMsg("required", field, param)
		}

	case "string":
		if _, ok := value.(string); !ok {
			return buildMsg("string", field, param)
		}

	case "min":
		n, _ := strconv.ParseFloat(param, 64)
		switch val := value.(type) {
		case float64:
			if val < n {
				return buildMsg("min_num", field, param)
			}
		case int:
			if float64(val) < n {
				return buildMsg("min_num", field, param)
			}
		case int64:
			if float64(val) < n {
				return buildMsg("min_num", field, param)
			}
		default:
			if utf8.RuneCountInString(str) < int(n) {
				return buildMsg("min", field, param)
			}
		}

	case "max":
		n, _ := strconv.ParseFloat(param, 64)
		switch val := value.(type) {
		case float64:
			if val > n {
				return buildMsg("max_num", field, param)
			}
		case int:
			if float64(val) > n {
				return buildMsg("max_num", field, param)
			}
		case int64:
			if float64(val) > n {
				return buildMsg("max_num", field, param)
			}
		default:
			if utf8.RuneCountInString(str) > int(n) {
				return buildMsg("max", field, param)
			}
		}

	case "email":
		if !reEmail.MatchString(str) {
			return buildMsg("email", field, param)
		}

	case "url":
		u, err := url.ParseRequestURI(str)
		if err != nil || u.Host == "" {
			return buildMsg("url", field, param)
		}

	case "in":
		for _, opt := range strings.Split(param, ",") {
			if strings.TrimSpace(opt) == str {
				return ""
			}
		}
		return buildMsg("in", field, param)

	case "not_in":
		for _, opt := range strings.Split(param, ",") {
			if strings.TrimSpace(opt) == str {
				return buildMsg("not_in", field, param)
			}
		}

	case "integer":
		if !isValidInteger(value) {
			return buildMsg("integer", field, param)
		}

	case "numeric":
		switch value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			// ok
		case string:
			if _, err := strconv.ParseFloat(str, 64); err != nil {
				return buildMsg("numeric", field, param)
			}
		default:
			return buildMsg("numeric", field, param)
		}

	case "boolean":
		if _, ok := value.(bool); !ok {
			return buildMsg("boolean", field, param)
		}

	case "confirmed":
		confirm := fmt.Sprintf("%v", v.input[field+"_confirmation"])
		if str != confirm {
			return buildMsg("confirmed", field, param)
		}

	case "min_digits":
		n, _ := strconv.Atoi(param)
		if countDigits(str) < n {
			return buildMsg("min_digits", field, param)
		}

	case "max_digits":
		n, _ := strconv.Atoi(param)
		if countDigits(str) > n {
			return buildMsg("max_digits", field, param)
		}

	case "regex":
		re, err := compileRegex(param)
		if err != nil || !re.MatchString(str) {
			return buildMsg("regex", field, param)
		}

	case "date":
		if _, err := parseDate(str); err != nil {
			return buildMsg("date", field, param)
		}

	case "before":
		t, err := parseDate(str)
		if err != nil {
			return buildMsg("date", field, param)
		}
		ref, err := parseDate(param)
		if err != nil || !t.Before(ref) {
			return buildMsg("before", field, param)
		}

	case "after":
		t, err := parseDate(str)
		if err != nil {
			return buildMsg("date", field, param)
		}
		ref, err := parseDate(param)
		if err != nil || !t.After(ref) {
			return buildMsg("after", field, param)
		}

	case "unique":
		if v.db == nil {
			return buildMsg("db_required", field, param)
		}
		// Format: unique:table,column,ignore_value,ignore_column
		// ignore_column defaults to "id" — same as Laravel
		parts := splitParam(param)
		table, column := parts[0], field
		if len(parts) >= 2 && parts[1] != "" {
			column = parts[1]
		}
		// table/column are developer-defined rule params, not user input — safe to interpolate
		var count int
		var err error
		if len(parts) >= 3 && parts[2] != "" {
			ignoreCol := "id"
			if len(parts) >= 4 && parts[3] != "" {
				ignoreCol = parts[3]
			}
			q := "SELECT COUNT(*) FROM " + table + " WHERE " + column + " = ? AND " + ignoreCol + " != ?"
			err = v.db.QueryRow(q, str, parts[2]).Scan(&count)
		} else {
			err = v.db.QueryRow("SELECT COUNT(*) FROM "+table+" WHERE "+column+" = ?", str).Scan(&count)
		}
		if err != nil || count > 0 {
			return buildMsg("unique", field, param)
		}

	case "exists":
		if v.db == nil {
			return buildMsg("db_required", field, param)
		}
		// Format: exists:table,column,ignore_value,ignore_column
		parts := splitParam(param)
		table, column := parts[0], field
		if len(parts) >= 2 && parts[1] != "" {
			column = parts[1]
		}
		var count int
		var err error
		if len(parts) >= 3 && parts[2] != "" {
			ignoreCol := "id"
			if len(parts) >= 4 && parts[3] != "" {
				ignoreCol = parts[3]
			}
			q := "SELECT COUNT(*) FROM " + table + " WHERE " + column + " = ? AND " + ignoreCol + " != ?"
			err = v.db.QueryRow(q, str, parts[2]).Scan(&count)
		} else {
			err = v.db.QueryRow("SELECT COUNT(*) FROM "+table+" WHERE "+column+" = ?", str).Scan(&count)
		}
		if err != nil || count == 0 {
			return buildMsg("exists", field, param)
		}

	// ── Boolean variants ──────────────────────────────────────────────────────

	case "accepted":
		if !isAccepted(value) {
			return buildMsg("accepted", field, param)
		}

	case "declined":
		if !isDeclined(value) {
			return buildMsg("declined", field, param)
		}

	case "accepted_if":
		// param: other,value
		parts := splitParam(param)
		if len(parts) < 2 {
			return buildMsg("accepted_if", field, param)
		}
		otherField, otherVal := parts[0], parts[1]
		if fmt.Sprintf("%v", v.input[otherField]) == otherVal {
			if !isAccepted(value) {
				return buildMsgWith("accepted_if", map[string]string{
					":field": humanize(field),
					":other": humanize(otherField),
					":value": otherVal,
				})
			}
		}

	case "declined_if":
		// param: other,value
		parts := splitParam(param)
		if len(parts) < 2 {
			return buildMsg("declined_if", field, param)
		}
		otherField, otherVal := parts[0], parts[1]
		if fmt.Sprintf("%v", v.input[otherField]) == otherVal {
			if !isDeclined(value) {
				return buildMsgWith("declined_if", map[string]string{
					":field": humanize(field),
					":other": humanize(otherField),
					":value": otherVal,
				})
			}
		}

	// ── String checks ─────────────────────────────────────────────────────────

	case "alpha":
		if param == "ascii" {
			for _, ch := range str {
				if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')) {
					return buildMsg("alpha_ascii", field, param)
				}
			}
		} else {
			for _, ch := range str {
				if !unicode.IsLetter(ch) {
					return buildMsg("alpha", field, param)
				}
			}
		}

	case "alpha_num":
		if param == "ascii" {
			for _, ch := range str {
				if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')) {
					return buildMsg("alpha_num", field, param)
				}
			}
		} else {
			for _, ch := range str {
				if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
					return buildMsg("alpha_num", field, param)
				}
			}
		}

	case "alpha_dash":
		for _, ch := range str {
			if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '-' && ch != '_' {
				return buildMsg("alpha_dash", field, param)
			}
		}

	case "ascii":
		for _, ch := range str {
			if ch >= 128 {
				return buildMsg("ascii", field, param)
			}
		}

	case "starts_with":
		matched := false
		for _, opt := range splitParam(param) {
			if strings.HasPrefix(str, opt) {
				matched = true
				break
			}
		}
		if !matched {
			return buildMsg("starts_with", field, param)
		}

	case "ends_with":
		matched := false
		for _, opt := range splitParam(param) {
			if strings.HasSuffix(str, opt) {
				matched = true
				break
			}
		}
		if !matched {
			return buildMsg("ends_with", field, param)
		}

	case "doesnt_start_with":
		for _, opt := range splitParam(param) {
			if strings.HasPrefix(str, opt) {
				return buildMsg("doesnt_start_with", field, param)
			}
		}

	case "doesnt_end_with":
		for _, opt := range splitParam(param) {
			if strings.HasSuffix(str, opt) {
				return buildMsg("doesnt_end_with", field, param)
			}
		}

	case "lowercase":
		if str != strings.ToLower(str) {
			return buildMsg("lowercase", field, param)
		}

	case "uppercase":
		if str != strings.ToUpper(str) {
			return buildMsg("uppercase", field, param)
		}

	case "json":
		if !json.Valid([]byte(str)) {
			return buildMsg("json", field, param)
		}

	case "uuid":
		if !reUUID.MatchString(str) {
			return buildMsg("uuid", field, param)
		}

	case "ulid":
		if !reULID.MatchString(str) {
			return buildMsg("ulid", field, param)
		}

	case "hex_color":
		if !reHexColor.MatchString(str) {
			return buildMsg("hex_color", field, param)
		}

	case "mac_address":
		if _, err := net.ParseMAC(str); err != nil {
			return buildMsg("mac_address", field, param)
		}

	case "ip":
		if net.ParseIP(str) == nil {
			return buildMsg("ip", field, param)
		}

	case "ipv4":
		ip := net.ParseIP(str)
		if ip == nil || ip.To4() == nil {
			return buildMsg("ipv4", field, param)
		}

	case "ipv6":
		ip := net.ParseIP(str)
		if ip == nil || ip.To4() != nil {
			return buildMsg("ipv6", field, param)
		}

	case "not_regex":
		re, err := compileRegex(param)
		if err == nil && re.MatchString(str) {
			return buildMsg("not_regex", field, param)
		}

	case "timezone":
		if _, err := time.LoadLocation(str); err != nil {
			return buildMsg("timezone", field, param)
		}

	case "same":
		otherVal := fmt.Sprintf("%v", v.input[param])
		if str != otherVal {
			return buildMsg("same", field, param)
		}

	case "different":
		otherVal := fmt.Sprintf("%v", v.input[param])
		if str == otherVal {
			return buildMsg("different", field, param)
		}

	// ── Numeric ───────────────────────────────────────────────────────────────

	case "between":
		parts := splitParam(param)
		if len(parts) < 2 {
			break
		}
		minVal, err1 := strconv.ParseFloat(parts[0], 64)
		maxVal, err2 := strconv.ParseFloat(parts[1], 64)
		if err1 != nil || err2 != nil {
			break
		}
		msg := buildMsgWith("between", map[string]string{
			":field": humanize(field),
			":min":   parts[0],
			":max":   parts[1],
		})
		switch val := value.(type) {
		case float64:
			if val < minVal || val > maxVal {
				return msg
			}
		case int:
			if float64(val) < minVal || float64(val) > maxVal {
				return msg
			}
		case int64:
			if float64(val) < minVal || float64(val) > maxVal {
				return msg
			}
		default:
			l := float64(utf8.RuneCountInString(str))
			if l < minVal || l > maxVal {
				return msg
			}
		}

	case "size":
		n, err := strconv.ParseFloat(param, 64)
		if err != nil {
			break
		}
		switch val := value.(type) {
		case float64:
			if val != n {
				return buildMsg("size_num", field, param)
			}
		case int:
			if float64(val) != n {
				return buildMsg("size_num", field, param)
			}
		case int64:
			if float64(val) != n {
				return buildMsg("size_num", field, param)
			}
		default:
			if float64(utf8.RuneCountInString(str)) != n {
				return buildMsg("size", field, param)
			}
		}

	case "digits":
		n, _ := strconv.Atoi(param)
		for _, ch := range str {
			if ch < '0' || ch > '9' {
				return buildMsg("digits", field, param)
			}
		}
		if utf8.RuneCountInString(str) != n {
			return buildMsg("digits", field, param)
		}

	case "digits_between":
		parts := splitParam(param)
		if len(parts) < 2 {
			break
		}
		minN, err1 := strconv.Atoi(parts[0])
		maxN, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			break
		}
		for _, ch := range str {
			if ch < '0' || ch > '9' {
				return buildMsgWith("digits_between", map[string]string{
					":field": humanize(field),
					":min":   parts[0],
					":max":   parts[1],
				})
			}
		}
		l := utf8.RuneCountInString(str)
		if l < minN || l > maxN {
			return buildMsgWith("digits_between", map[string]string{
				":field": humanize(field),
				":min":   parts[0],
				":max":   parts[1],
			})
		}

	case "multiple_of":
		divisor, err := strconv.ParseFloat(param, 64)
		if err != nil || divisor == 0 {
			break
		}
		f, err := toFloat(value)
		if err != nil {
			return buildMsg("multiple_of", field, param)
		}
		if math.Mod(f, divisor) != 0 {
			return buildMsg("multiple_of", field, param)
		}

	case "gt":
		ref, refIsField := fieldOrFloat(v.input, param)
		f, err := toFloat(value)
		if err != nil {
			return buildMsg("gt", field, param)
		}
		_ = refIsField
		if f <= ref {
			return buildMsg("gt", field, param)
		}

	case "gte":
		ref, _ := fieldOrFloat(v.input, param)
		f, err := toFloat(value)
		if err != nil {
			return buildMsg("gte", field, param)
		}
		if f < ref {
			return buildMsg("gte", field, param)
		}

	case "lt":
		ref, _ := fieldOrFloat(v.input, param)
		f, err := toFloat(value)
		if err != nil {
			return buildMsg("lt", field, param)
		}
		if f >= ref {
			return buildMsg("lt", field, param)
		}

	case "lte":
		ref, _ := fieldOrFloat(v.input, param)
		f, err := toFloat(value)
		if err != nil {
			return buildMsg("lte", field, param)
		}
		if f > ref {
			return buildMsg("lte", field, param)
		}

	// ── Date ──────────────────────────────────────────────────────────────────

	case "after_or_equal":
		t, err := parseDate(str)
		if err != nil {
			return buildMsg("date", field, param)
		}
		ref, err := parseDate(param)
		if err != nil || t.Before(ref) {
			return buildMsg("after_or_equal", field, param)
		}

	case "before_or_equal":
		t, err := parseDate(str)
		if err != nil {
			return buildMsg("date", field, param)
		}
		ref, err := parseDate(param)
		if err != nil || t.After(ref) {
			return buildMsg("before_or_equal", field, param)
		}

	case "date_equals":
		t, err := parseDate(str)
		if err != nil {
			return buildMsg("date", field, param)
		}
		ref, err := parseDate(param)
		if err != nil || !t.Equal(ref) {
			return buildMsg("date_equals", field, param)
		}

	case "date_format":
		if _, err := time.Parse(param, str); err != nil {
			return buildMsg("date_format", field, param)
		}

	// ── Array ─────────────────────────────────────────────────────────────────

	case "array":
		if _, ok := value.([]any); !ok {
			return buildMsg("array", field, param)
		}

	case "distinct":
		arr, ok := value.([]any)
		if !ok {
			return buildMsg("array", field, param)
		}
		seen := make(map[string]bool)
		for _, item := range arr {
			key := fmt.Sprintf("%v", item)
			if seen[key] {
				return buildMsg("distinct", field, param)
			}
			seen[key] = true
		}

	// ── Prohibited (in-check variants) ────────────────────────────────────────

	case "prohibited":
		// If we reach here, value is present and non-empty — prohibited.
		return buildMsg("prohibited", field, param)

	case "prohibited_if":
		parts := splitParam(param)
		if len(parts) < 2 {
			break
		}
		otherField, otherVal := parts[0], parts[1]
		if fmt.Sprintf("%v", v.input[otherField]) == otherVal {
			return buildMsgWith("prohibited_if", map[string]string{
				":field": humanize(field),
				":other": humanize(otherField),
				":value": otherVal,
			})
		}

	case "prohibited_unless":
		parts := splitParam(param)
		if len(parts) < 2 {
			break
		}
		otherField, otherVal := parts[0], parts[1]
		if fmt.Sprintf("%v", v.input[otherField]) != otherVal {
			return buildMsgWith("prohibited_unless", map[string]string{
				":field": humanize(field),
				":other": humanize(otherField),
				":value": otherVal,
			})
		}
	}

	return ""
}

var dateLayouts = []string{
	"2006-01-02",
	"2006-01-02 15:04:05",
	time.RFC3339,
}

func parseDate(s string) (time.Time, error) {
	switch s {
	case "today":
		return time.Now().Truncate(24 * time.Hour), nil
	case "tomorrow":
		return time.Now().Add(24 * time.Hour).Truncate(24 * time.Hour), nil
	case "yesterday":
		return time.Now().Add(-24 * time.Hour).Truncate(24 * time.Hour), nil
	}
	for _, layout := range dateLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse date %q", s)
}

// getNestedValue resolves a dot-notation path (e.g. "user.address.postcode") from input.
func getNestedValue(input Input, dotPath string) (any, bool) {
	parts := strings.SplitN(dotPath, ".", 2)
	val, ok := input[parts[0]]
	if !ok || len(parts) == 1 {
		return val, ok
	}
	if nested, ok := val.(map[string]any); ok {
		return getNestedValue(Input(nested), parts[1])
	}
	return nil, false
}

func isEmpty(value any) bool {
	if value == nil {
		return true
	}
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) == ""
	case []any:
		return len(v) == 0
	case map[string]any:
		return len(v) == 0
	}
	return false
}

func isValidInteger(value any) bool {
	switch val := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	case float32:
		return float64(val) == math.Trunc(float64(val))
	case float64:
		return val == math.Trunc(val)
	case string:
		_, err := strconv.ParseInt(val, 10, 64)
		return err == nil
	default:
		return false
	}
}

func splitParam(param string) []string {
	parts := strings.Split(param, ",")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return parts
}

func countDigits(s string) int {
	n := 0
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			n++
		}
	}
	return n
}

// isAccepted returns true if value is one of: "yes", "on", "1", "true", true, 1.
func isAccepted(value any) bool {
	switch v := value.(type) {
	case bool:
		return v
	case int:
		return v == 1
	case int64:
		return v == 1
	case float64:
		return v == 1
	case string:
		s := strings.ToLower(strings.TrimSpace(v))
		return s == "yes" || s == "on" || s == "1" || s == "true"
	}
	return false
}

// isDeclined returns true if value is one of: "no", "off", "0", "false", false, 0.
func isDeclined(value any) bool {
	switch v := value.(type) {
	case bool:
		return !v
	case int:
		return v == 0
	case int64:
		return v == 0
	case float64:
		return v == 0
	case string:
		s := strings.ToLower(strings.TrimSpace(v))
		return s == "no" || s == "off" || s == "0" || s == "false"
	}
	return false
}

// toFloat converts a value to float64.
func toFloat(value any) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	}
	return 0, fmt.Errorf("cannot convert %T to float64", value)
}

// fieldOrFloat resolves param either as a field name in input or as a literal float.
// Returns the float value and whether it was resolved from a field.
func fieldOrFloat(input Input, param string) (float64, bool) {
	if fieldVal, ok := input[param]; ok {
		if f, err := toFloat(fieldVal); err == nil {
			return f, true
		}
	}
	f, err := strconv.ParseFloat(param, 64)
	if err == nil {
		return f, false
	}
	return 0, false
}
