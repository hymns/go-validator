package validator

import (
	"strconv"
	"strings"
)

// RuleBuilder builds a pipe-separated rule string in a compile-safe, chainable way.
// Call R() to start, chain methods, then call Build() to get the rule string for use in Rules{}.
//
//	validator.Rules{
//	    "email": validator.R().Required().Email().Max(100).Build(),
//	    "age":   validator.R().Required().Integer().Min(18).Build(),
//	    "role":  validator.R().Required().In("admin", "user").Build(),
//	}
type RuleBuilder struct {
	parts []string
}

// R creates a new RuleBuilder.
func R() *RuleBuilder {
	return &RuleBuilder{}
}

func (rb *RuleBuilder) add(name string) *RuleBuilder {
	rb.parts = append(rb.parts, name)
	return rb
}

func (rb *RuleBuilder) addParam(name, param string) *RuleBuilder {
	rb.parts = append(rb.parts, name+":"+param)
	return rb
}

// Build returns the pipe-separated rule string.
func (rb *RuleBuilder) Build() string {
	return strings.Join(rb.parts, "|")
}

// ── Presence / Flow ───────────────────────────────────────────────────────────

func (rb *RuleBuilder) Required() *RuleBuilder  { return rb.add("required") }
func (rb *RuleBuilder) Nullable() *RuleBuilder  { return rb.add("nullable") }
func (rb *RuleBuilder) Present() *RuleBuilder   { return rb.add("present") }
func (rb *RuleBuilder) Filled() *RuleBuilder    { return rb.add("filled") }
func (rb *RuleBuilder) Bail() *RuleBuilder      { return rb.add("bail") }
func (rb *RuleBuilder) Confirmed() *RuleBuilder { return rb.add("confirmed") }

// RequiredIf requires the field when otherField equals value.
func (rb *RuleBuilder) RequiredIf(otherField, value string) *RuleBuilder {
	return rb.addParam("required_if", otherField+","+value)
}

// RequiredUnless requires the field unless otherField equals value.
func (rb *RuleBuilder) RequiredUnless(otherField, value string) *RuleBuilder {
	return rb.addParam("required_unless", otherField+","+value)
}

// RequiredWith requires the field when any of the given fields are present.
func (rb *RuleBuilder) RequiredWith(fields ...string) *RuleBuilder {
	return rb.addParam("required_with", strings.Join(fields, ","))
}

// RequiredWithout requires the field when any of the given fields are absent.
func (rb *RuleBuilder) RequiredWithout(fields ...string) *RuleBuilder {
	return rb.addParam("required_without", strings.Join(fields, ","))
}

// RequiredWithAll requires the field when all of the given fields are present.
func (rb *RuleBuilder) RequiredWithAll(fields ...string) *RuleBuilder {
	return rb.addParam("required_with_all", strings.Join(fields, ","))
}

// RequiredWithoutAll requires the field when none of the given fields are present.
func (rb *RuleBuilder) RequiredWithoutAll(fields ...string) *RuleBuilder {
	return rb.addParam("required_without_all", strings.Join(fields, ","))
}

// ── Type checks ───────────────────────────────────────────────────────────────

func (rb *RuleBuilder) Str() *RuleBuilder     { return rb.add("string") }
func (rb *RuleBuilder) Integer() *RuleBuilder { return rb.add("integer") }
func (rb *RuleBuilder) Numeric() *RuleBuilder { return rb.add("numeric") }
func (rb *RuleBuilder) Boolean() *RuleBuilder { return rb.add("boolean") }
func (rb *RuleBuilder) Array() *RuleBuilder   { return rb.add("array") }

// ── Size / Length ─────────────────────────────────────────────────────────────

func (rb *RuleBuilder) Min(n int) *RuleBuilder {
	return rb.addParam("min", strconv.Itoa(n))
}

func (rb *RuleBuilder) Max(n int) *RuleBuilder {
	return rb.addParam("max", strconv.Itoa(n))
}

func (rb *RuleBuilder) Size(n int) *RuleBuilder {
	return rb.addParam("size", strconv.Itoa(n))
}

// Between validates that the value is between min and max (inclusive).
func (rb *RuleBuilder) Between(min, max int) *RuleBuilder {
	return rb.addParam("between", strconv.Itoa(min)+","+strconv.Itoa(max))
}

func (rb *RuleBuilder) MinDigits(n int) *RuleBuilder {
	return rb.addParam("min_digits", strconv.Itoa(n))
}

func (rb *RuleBuilder) MaxDigits(n int) *RuleBuilder {
	return rb.addParam("max_digits", strconv.Itoa(n))
}

func (rb *RuleBuilder) Digits(n int) *RuleBuilder {
	return rb.addParam("digits", strconv.Itoa(n))
}

// DigitsBetween validates the digit count is between min and max.
func (rb *RuleBuilder) DigitsBetween(min, max int) *RuleBuilder {
	return rb.addParam("digits_between", strconv.Itoa(min)+","+strconv.Itoa(max))
}

// MultipleOf validates that the value is a multiple of divisor.
func (rb *RuleBuilder) MultipleOf(divisor float64) *RuleBuilder {
	return rb.addParam("multiple_of", strconv.FormatFloat(divisor, 'f', -1, 64))
}

// ── Comparison ────────────────────────────────────────────────────────────────

// Gt validates the value is greater than param (literal or field name).
func (rb *RuleBuilder) Gt(param string) *RuleBuilder  { return rb.addParam("gt", param) }
func (rb *RuleBuilder) Gte(param string) *RuleBuilder { return rb.addParam("gte", param) }
func (rb *RuleBuilder) Lt(param string) *RuleBuilder  { return rb.addParam("lt", param) }
func (rb *RuleBuilder) Lte(param string) *RuleBuilder { return rb.addParam("lte", param) }

// Same validates the field matches another field.
func (rb *RuleBuilder) Same(otherField string) *RuleBuilder {
	return rb.addParam("same", otherField)
}

// Different validates the field differs from another field.
func (rb *RuleBuilder) Different(otherField string) *RuleBuilder {
	return rb.addParam("different", otherField)
}

// ── String rules ──────────────────────────────────────────────────────────────

func (rb *RuleBuilder) Email() *RuleBuilder      { return rb.add("email") }
func (rb *RuleBuilder) URL() *RuleBuilder        { return rb.add("url") }
func (rb *RuleBuilder) UUID() *RuleBuilder       { return rb.add("uuid") }
func (rb *RuleBuilder) ULID() *RuleBuilder       { return rb.add("ulid") }
func (rb *RuleBuilder) JSON() *RuleBuilder       { return rb.add("json") }
func (rb *RuleBuilder) Alpha() *RuleBuilder      { return rb.add("alpha") }
func (rb *RuleBuilder) AlphaNum() *RuleBuilder   { return rb.add("alpha_num") }
func (rb *RuleBuilder) AlphaDash() *RuleBuilder  { return rb.add("alpha_dash") }
func (rb *RuleBuilder) ASCII() *RuleBuilder      { return rb.add("ascii") }
func (rb *RuleBuilder) Lowercase() *RuleBuilder  { return rb.add("lowercase") }
func (rb *RuleBuilder) Uppercase() *RuleBuilder  { return rb.add("uppercase") }
func (rb *RuleBuilder) HexColor() *RuleBuilder   { return rb.add("hex_color") }
func (rb *RuleBuilder) MACAddress() *RuleBuilder { return rb.add("mac_address") }
func (rb *RuleBuilder) IP() *RuleBuilder         { return rb.add("ip") }
func (rb *RuleBuilder) IPv4() *RuleBuilder       { return rb.add("ipv4") }
func (rb *RuleBuilder) IPv6() *RuleBuilder       { return rb.add("ipv6") }
func (rb *RuleBuilder) Timezone() *RuleBuilder   { return rb.add("timezone") }
func (rb *RuleBuilder) Distinct() *RuleBuilder   { return rb.add("distinct") }

func (rb *RuleBuilder) Regex(pattern string) *RuleBuilder {
	return rb.addParam("regex", pattern)
}

func (rb *RuleBuilder) NotRegex(pattern string) *RuleBuilder {
	return rb.addParam("not_regex", pattern)
}

// StartsWith validates the string starts with one of the given prefixes.
func (rb *RuleBuilder) StartsWith(prefixes ...string) *RuleBuilder {
	return rb.addParam("starts_with", strings.Join(prefixes, ","))
}

// EndsWith validates the string ends with one of the given suffixes.
func (rb *RuleBuilder) EndsWith(suffixes ...string) *RuleBuilder {
	return rb.addParam("ends_with", strings.Join(suffixes, ","))
}

// DoesntStartWith validates the string does not start with any of the given prefixes.
func (rb *RuleBuilder) DoesntStartWith(prefixes ...string) *RuleBuilder {
	return rb.addParam("doesnt_start_with", strings.Join(prefixes, ","))
}

// DoesntEndWith validates the string does not end with any of the given suffixes.
func (rb *RuleBuilder) DoesntEndWith(suffixes ...string) *RuleBuilder {
	return rb.addParam("doesnt_end_with", strings.Join(suffixes, ","))
}

// ── Enum ──────────────────────────────────────────────────────────────────────

// In validates the value is one of the given options.
func (rb *RuleBuilder) In(values ...string) *RuleBuilder {
	return rb.addParam("in", strings.Join(values, ","))
}

// NotIn validates the value is not one of the given options.
func (rb *RuleBuilder) NotIn(values ...string) *RuleBuilder {
	return rb.addParam("not_in", strings.Join(values, ","))
}

// ── Boolean variants ──────────────────────────────────────────────────────────

func (rb *RuleBuilder) Accepted() *RuleBuilder { return rb.add("accepted") }
func (rb *RuleBuilder) Declined() *RuleBuilder { return rb.add("declined") }

func (rb *RuleBuilder) AcceptedIf(otherField, value string) *RuleBuilder {
	return rb.addParam("accepted_if", otherField+","+value)
}

func (rb *RuleBuilder) DeclinedIf(otherField, value string) *RuleBuilder {
	return rb.addParam("declined_if", otherField+","+value)
}

// ── Prohibition ───────────────────────────────────────────────────────────────

func (rb *RuleBuilder) Prohibited() *RuleBuilder { return rb.add("prohibited") }

func (rb *RuleBuilder) ProhibitedIf(otherField, value string) *RuleBuilder {
	return rb.addParam("prohibited_if", otherField+","+value)
}

func (rb *RuleBuilder) ProhibitedUnless(otherField, value string) *RuleBuilder {
	return rb.addParam("prohibited_unless", otherField+","+value)
}

// ── Date ──────────────────────────────────────────────────────────────────────

func (rb *RuleBuilder) Date() *RuleBuilder { return rb.add("date") }

func (rb *RuleBuilder) Before(date string) *RuleBuilder     { return rb.addParam("before", date) }
func (rb *RuleBuilder) After(date string) *RuleBuilder      { return rb.addParam("after", date) }
func (rb *RuleBuilder) DateEquals(date string) *RuleBuilder { return rb.addParam("date_equals", date) }
func (rb *RuleBuilder) DateFormat(layout string) *RuleBuilder {
	return rb.addParam("date_format", layout)
}

func (rb *RuleBuilder) BeforeOrEqual(date string) *RuleBuilder {
	return rb.addParam("before_or_equal", date)
}

func (rb *RuleBuilder) AfterOrEqual(date string) *RuleBuilder {
	return rb.addParam("after_or_equal", date)
}

// ── Database ──────────────────────────────────────────────────────────────────

// Unique validates the value does not exist in table.column.
// Optional: ignoreValue and ignoreColumn (defaults to "id").
func (rb *RuleBuilder) Unique(table string, column ...string) *RuleBuilder {
	param := table
	if len(column) > 0 {
		param += "," + strings.Join(column, ",")
	}
	return rb.addParam("unique", param)
}

// Exists validates the value exists in table.column.
func (rb *RuleBuilder) Exists(table string, column ...string) *RuleBuilder {
	param := table
	if len(column) > 0 {
		param += "," + strings.Join(column, ",")
	}
	return rb.addParam("exists", param)
}
