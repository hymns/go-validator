package validator

import "strings"

var defaultMessages = map[string]string{
	"required":   "The :field field is required.",
	"string":     "The :field must be a string.",
	"min":        "The :field must be at least :param characters.",
	"min_num":    "The :field must be at least :param.",
	"max":        "The :field must not exceed :param characters.",
	"max_num":    "The :field must not exceed :param.",
	"email":      "The :field must be a valid email address.",
	"url":        "The :field must be a valid URL.",
	"in":         "The selected :field is invalid.",
	"not_in":     "The selected :field is invalid.",
	"integer":    "The :field must be an integer.",
	"numeric":    "The :field must be a number.",
	"boolean":    "The :field must be true or false.",
	"confirmed":  "The :field confirmation does not match.",
	"min_digits": "The :field must have at least :param digits.",
	"max_digits": "The :field must not have more than :param digits.",
	"regex":      "The :field format is invalid.",
	"date":       "The :field must be a valid date.",
	"before":     "The :field must be a date before :param.",
	"after":      "The :field must be a date after :param.",
	"unique":     "The :field has already been taken.",
	"exists":     "The selected :field is invalid.",
	"db_required": "The :field rule requires a database connection. Call WithDB().",

	// Boolean variants
	"accepted":    "The :field must be accepted.",
	"declined":    "The :field must be declined.",
	"accepted_if": "The :field must be accepted when :other is :value.",
	"declined_if": "The :field must be declined when :other is :value.",

	// String checks
	"alpha":              "The :field must only contain letters.",
	"alpha_ascii":        "The :field must only contain ASCII letters.",
	"alpha_num":          "The :field must only contain letters and numbers.",
	"alpha_dash":         "The :field must only contain letters, numbers, dashes, and underscores.",
	"ascii":              "The :field must only contain ASCII characters.",
	"starts_with":        "The :field must start with one of the following: :param.",
	"ends_with":          "The :field must end with one of the following: :param.",
	"doesnt_start_with":  "The :field must not start with one of the following: :param.",
	"doesnt_end_with":    "The :field must not end with one of the following: :param.",
	"lowercase":          "The :field must be lowercase.",
	"uppercase":          "The :field must be uppercase.",
	"json":               "The :field must be a valid JSON string.",
	"uuid":               "The :field must be a valid UUID.",
	"ulid":               "The :field must be a valid ULID.",
	"hex_color":          "The :field must be a valid hexadecimal color.",
	"mac_address":        "The :field must be a valid MAC address.",
	"ip":                 "The :field must be a valid IP address.",
	"ipv4":               "The :field must be a valid IPv4 address.",
	"ipv6":               "The :field must be a valid IPv6 address.",
	"not_regex":          "The :field format is invalid.",
	"timezone":           "The :field must be a valid timezone.",
	"same":               "The :field and :param must match.",
	"different":          "The :field and :param must be different.",

	// Numeric
	"between":       "The :field must be between :min and :max.",
	"size":          "The :field must be :param characters.",
	"size_num":      "The :field must be :param.",
	"digits":        "The :field must be :param digits.",
	"digits_between": "The :field must be between :min and :max digits.",
	"multiple_of":   "The :field must be a multiple of :param.",
	"gt":            "The :field must be greater than :param.",
	"gte":           "The :field must be greater than or equal to :param.",
	"lt":            "The :field must be less than :param.",
	"lte":           "The :field must be less than or equal to :param.",

	// Date
	"after_or_equal":  "The :field must be a date on or after :param.",
	"before_or_equal": "The :field must be a date on or before :param.",
	"date_equals":     "The :field must be a date equal to :param.",
	"date_format":     "The :field does not match the format :param.",

	// Array
	"array":    "The :field must be an array.",
	"distinct": "The :field must not have duplicate values.",

	// Presence/prohibition
	"prohibited":        "The :field field is prohibited.",
	"prohibited_if":     "The :field field is prohibited when :other is :value.",
	"prohibited_unless": "The :field field is prohibited unless :other is :value.",
	"filled":            "The :field field must not be empty when present.",
	"present":           "The :field field must be present.",

	// Conditional required
	"required_if":           "The :field field is required when :other is :value.",
	"required_unless":       "The :field field is required unless :other is :value.",
	"required_with":         "The :field field is required when :other is present.",
	"required_without":      "The :field field is required when :other is not present.",
	"required_with_all":     "The :field field is required when :other are present.",
	"required_without_all":  "The :field field is required when none of :other are present.",
}

// buildMsgWith builds an error message by replacing multiple named placeholders.
func buildMsgWith(key string, replacements map[string]string) string {
	tmpl, ok := defaultMessages[key]
	if !ok {
		if field, exists := replacements[":field"]; exists {
			return "The " + humanize(field) + " is invalid."
		}
		return "The field is invalid."
	}
	msg := tmpl
	for placeholder, val := range replacements {
		msg = strings.ReplaceAll(msg, placeholder, val)
	}
	return msg
}

// buildMsg builds an error message replacing :field and :param.
// It delegates to buildMsgWith for consistency.
func buildMsg(key, field, param string) string {
	return buildMsgWith(key, map[string]string{
		":field": humanize(field),
		":param": param,
	})
}

func humanize(field string) string {
	return strings.ReplaceAll(field, "_", " ")
}
