package validator_test

import (
	"fmt"
	"testing"

	validator "github.com/hymns/go-validator"
)

func TestRequired(t *testing.T) {
	cases := []struct {
		name  string
		input validator.Input
		fails bool
	}{
		{"empty string fails", validator.Input{"name": ""}, true},
		{"missing field fails", validator.Input{}, true},
		{"whitespace fails", validator.Input{"name": "   "}, true},
		{"nil fails", validator.Input{"name": nil}, true},
		{"value passes", validator.Input{"name": "John"}, false},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.Make(tt.input, validator.Rules{"name": "required"})
			if v.Fails() != tt.fails {
				t.Errorf("Fails()=%v want %v, errors: %v", v.Fails(), tt.fails, v.Errors())
			}
		})
	}
}

func TestOptionalField(t *testing.T) {
	v := validator.Make(validator.Input{}, validator.Rules{"bio": "min:10|max:500"})
	if v.Fails() {
		t.Errorf("optional empty field should pass, got: %v", v.Errors())
	}
}

func TestNullable(t *testing.T) {
	v := validator.Make(validator.Input{"bio": nil}, validator.Rules{"bio": "nullable|min:10"})
	if v.Fails() {
		t.Errorf("nullable nil should pass, got: %v", v.Errors())
	}
}

func TestNullableWithOtherRules(t *testing.T) {
	// nil skips other rules
	v := validator.Make(validator.Input{"bio": nil}, validator.Rules{"bio": "nullable|min:10|max:500"})
	if v.Fails() {
		t.Errorf("nullable nil should skip other rules, got: %v", v.Errors())
	}
	// non-nil value still validates
	v2 := validator.Make(validator.Input{"bio": "Hi"}, validator.Rules{"bio": "nullable|min:10"})
	if v2.Passes() {
		t.Error("nullable with short value should still fail min:10")
	}
}

func TestStopsAtFirstError(t *testing.T) {
	v := validator.Make(validator.Input{"email": ""}, validator.Rules{"email": "required|email"})
	errs := v.Errors()["email"]
	if len(errs) != 1 {
		t.Errorf("expected exactly 1 error, got %d: %v", len(errs), errs)
	}
}

func TestBailBehavior(t *testing.T) {
	v := validator.Make(validator.Input{"email": ""}, validator.Rules{"email": "bail|required|email"})
	errs := v.Errors()["email"]
	if len(errs) != 1 {
		t.Errorf("bail: expected exactly 1 error, got %d: %v", len(errs), errs)
	}
}

func TestCustomMessages(t *testing.T) {
	v := validator.Make(
		validator.Input{"email": ""},
		validator.Rules{"email": "required|email"},
	).Messages(validator.Messages{
		"email.required": "Alamat emel diperlukan.",
		"email.email":    "Format emel tidak sah.",
	})

	if v.Passes() {
		t.Fatal("expected validation to fail")
	}
	got := v.Errors().First("email")
	if got != "Alamat emel diperlukan." {
		t.Errorf("expected custom message, got: %q", got)
	}
}

func TestCustomMessagesForConditionalRules(t *testing.T) {
	v := validator.Make(
		validator.Input{"role": "admin"},
		validator.Rules{"permissions": "required_if:role,admin"},
	).Messages(validator.Messages{
		"permissions.required_if": "Permissions diperlukan untuk admin.",
	})
	if v.Passes() {
		t.Fatal("expected validation to fail")
	}
	got := v.Errors().First("permissions")
	if got != "Permissions diperlukan untuk admin." {
		t.Errorf("expected custom message, got: %q", got)
	}
}

func TestCustomRule(t *testing.T) {
	validator.Extend("starts_with_a", func(field string, value any, param string) error {
		s, _ := value.(string)
		if len(s) == 0 || s[0] != 'A' {
			return fmt.Errorf("The %s must start with 'A'.", field)
		}
		return nil
	})

	v := validator.Make(validator.Input{"name": "Bob"}, validator.Rules{"name": "required|starts_with_a"})
	if v.Passes() {
		t.Error("expected custom rule to fail for 'Bob'")
	}

	v2 := validator.Make(validator.Input{"name": "Alice"}, validator.Rules{"name": "required|starts_with_a"})
	if v2.Fails() {
		t.Errorf("expected custom rule to pass for 'Alice', got: %v", v2.Errors())
	}
}

func TestErrorBagHelpers(t *testing.T) {
	v := validator.Make(validator.Input{}, validator.Rules{
		"name":  "required",
		"email": "required|email",
	})
	errs := v.Errors()

	if !errs.Has("name") {
		t.Error("expected Has('name') to be true")
	}
	if errs.Has("nonexistent") {
		t.Error("expected Has('nonexistent') to be false")
	}
	if errs.First("name") == "" {
		t.Error("expected First('name') to be non-empty")
	}
	if errs.First("nonexistent") != "" {
		t.Error("expected First('nonexistent') to be empty string")
	}
}

func TestFullRegistrationFlow(t *testing.T) {
	v := validator.Make(validator.Input{
		"name":                  "Ali",
		"email":                 "ali@example.com",
		"password":              "secret123",
		"password_confirmation": "secret123",
		"role":                  "admin",
	}, validator.Rules{
		"name":                  "required|min:3|max:50",
		"email":                 "required|email",
		"password":              "required|min:8|confirmed",
		"password_confirmation": "required",
		"role":                  "required|in:admin,user,editor",
	})

	if v.Fails() {
		t.Errorf("expected full registration to pass, got: %v", v.Errors())
	}
}
