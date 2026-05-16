package validator_test

import (
	"testing"

	validator "github.com/hymns/go-validator"
)

func TestEmail(t *testing.T) {
	cases := []struct {
		email string
		fails bool
	}{
		{"test@example.com", false},
		{"user@mail.example.com", false},
		{"not-an-email", true},
		{"test@", true},
		{"@example.com", true},
	}
	for _, tt := range cases {
		t.Run(tt.email, func(t *testing.T) {
			v := validator.Make(validator.Input{"email": tt.email}, validator.Rules{"email": "required|email"})
			if v.Fails() != tt.fails {
				t.Errorf("email=%q Fails()=%v want %v, errors: %v", tt.email, v.Fails(), tt.fails, v.Errors())
			}
		})
	}
}

func TestURL(t *testing.T) {
	cases := []struct {
		val   string
		fails bool
	}{
		{"https://example.com", false},
		{"http://sub.domain.com/path?q=1", false},
		{"ftp://files.example.com", false},
		{"not-a-url", true},
		{"//missing-scheme.com", true},
	}
	for _, tt := range cases {
		v := validator.Make(validator.Input{"link": tt.val}, validator.Rules{"link": "required|url"})
		if v.Fails() != tt.fails {
			t.Errorf("url=%q Fails()=%v want %v, errors: %v", tt.val, v.Fails(), tt.fails, v.Errors())
		}
	}
}

func TestMinMaxString(t *testing.T) {
	rules := validator.Rules{"name": "required|min:3|max:10"}
	cases := []struct {
		val   string
		fails bool
	}{
		{"Jo", true},
		{"Jon", false},
		{"HelloWorld", false},
		{"HelloWorldXY", true},
	}
	for _, tt := range cases {
		v := validator.Make(validator.Input{"name": tt.val}, rules)
		if v.Fails() != tt.fails {
			t.Errorf("name=%q Fails()=%v want %v, errors: %v", tt.val, v.Fails(), tt.fails, v.Errors())
		}
	}
}

func TestStringType(t *testing.T) {
	v := validator.Make(validator.Input{"val": 42}, validator.Rules{"val": "required|string"})
	if v.Passes() {
		t.Error("expected non-string to fail string rule")
	}
	v2 := validator.Make(validator.Input{"val": "hello"}, validator.Rules{"val": "required|string"})
	if v2.Fails() {
		t.Errorf("string should pass for string value, got: %v", v2.Errors())
	}
}

func TestRegex(t *testing.T) {
	rules := validator.Rules{"code": "required|regex:^[A-Z]{2}[0-9]{4}$"}
	cases := []struct {
		val   string
		fails bool
	}{
		{"AB1234", false},
		{"ab1234", true},
		{"AB12", true},
		{"ABCD1234", true},
	}
	for _, tt := range cases {
		v := validator.Make(validator.Input{"code": tt.val}, rules)
		if v.Fails() != tt.fails {
			t.Errorf("regex(%q) Fails()=%v want %v", tt.val, v.Fails(), tt.fails)
		}
	}
}

func TestNotRegex(t *testing.T) {
	v := validator.Make(validator.Input{"name": "Ali123"}, validator.Rules{"name": "required|not_regex:^[0-9]+$"})
	if v.Fails() {
		t.Errorf("not_regex should pass when pattern not matched, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"name": "12345"}, validator.Rules{"name": "required|not_regex:^[0-9]+$"})
	if v2.Passes() {
		t.Error("not_regex should fail when pattern matches")
	}
}

func TestAlpha(t *testing.T) {
	cases := []struct {
		rule  string
		val   string
		fails bool
	}{
		{"alpha", "HelloWorld", false},
		{"alpha", "Hello123", true},
		{"alpha", "Héllo", false},
		{"alpha:ascii", "Hello", false},
		{"alpha:ascii", "Héllo", true},
		{"alpha_num", "Hello123", false},
		{"alpha_num", "Hello!", true},
		{"alpha_dash", "hello-world_123", false},
		{"alpha_dash", "hello world", true},
		{"ascii", "Hello", false},
		{"ascii", "Héllo", true},
	}
	for _, tt := range cases {
		v := validator.Make(validator.Input{"val": tt.val}, validator.Rules{"val": "required|" + tt.rule})
		if v.Fails() != tt.fails {
			t.Errorf("rule=%q val=%q Fails()=%v want %v, errors: %v", tt.rule, tt.val, v.Fails(), tt.fails, v.Errors())
		}
	}
}

func TestStartsWith(t *testing.T) {
	v := validator.Make(validator.Input{"code": "IMG_001"}, validator.Rules{"code": "required|starts_with:IMG,VID"})
	if v.Fails() {
		t.Errorf("starts_with should pass, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"code": "DOC_001"}, validator.Rules{"code": "required|starts_with:IMG,VID"})
	if v2.Passes() {
		t.Error("starts_with should fail for 'DOC_001'")
	}
}

func TestEndsWith(t *testing.T) {
	v := validator.Make(validator.Input{"file": "report.pdf"}, validator.Rules{"file": "required|ends_with:.pdf,.doc"})
	if v.Fails() {
		t.Errorf("ends_with should pass, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"file": "report.txt"}, validator.Rules{"file": "required|ends_with:.pdf,.doc"})
	if v2.Passes() {
		t.Error("ends_with should fail for .txt")
	}
}

func TestDoesntStartWith(t *testing.T) {
	v := validator.Make(validator.Input{"url": "mailto:x@x.com"}, validator.Rules{"url": "required|doesnt_start_with:http,ftp"})
	if v.Fails() {
		t.Errorf("doesnt_start_with should pass, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"url": "https://x.com"}, validator.Rules{"url": "required|doesnt_start_with:http,ftp"})
	if v2.Passes() {
		t.Error("doesnt_start_with should fail for https://")
	}
}

func TestDoesntEndWith(t *testing.T) {
	v := validator.Make(validator.Input{"file": "file.pdf"}, validator.Rules{"file": "required|doesnt_end_with:.exe,.bat"})
	if v.Fails() {
		t.Errorf("doesnt_end_with should pass, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"file": "file.exe"}, validator.Rules{"file": "required|doesnt_end_with:.exe,.bat"})
	if v2.Passes() {
		t.Error("doesnt_end_with should fail for .exe")
	}
}

func TestLowercase(t *testing.T) {
	v := validator.Make(validator.Input{"code": "hello"}, validator.Rules{"code": "required|lowercase"})
	if v.Fails() {
		t.Errorf("lowercase should pass, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"code": "Hello"}, validator.Rules{"code": "required|lowercase"})
	if v2.Passes() {
		t.Error("lowercase should fail for 'Hello'")
	}
}

func TestUppercase(t *testing.T) {
	v := validator.Make(validator.Input{"code": "HELLO"}, validator.Rules{"code": "required|uppercase"})
	if v.Fails() {
		t.Errorf("uppercase should pass, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"code": "Hello"}, validator.Rules{"code": "required|uppercase"})
	if v2.Passes() {
		t.Error("uppercase should fail for 'Hello'")
	}
}

func TestJSON(t *testing.T) {
	cases := []struct {
		val   string
		fails bool
	}{
		{`{"key":"value"}`, false},
		{`[1,2,3]`, false},
		{`"string"`, false},
		{"not json", true},
		{"{bad}", true},
	}
	for _, tt := range cases {
		v := validator.Make(validator.Input{"data": tt.val}, validator.Rules{"data": "required|json"})
		if v.Fails() != tt.fails {
			t.Errorf("json(%q) Fails()=%v want %v", tt.val, v.Fails(), tt.fails)
		}
	}
}

func TestUUID(t *testing.T) {
	v := validator.Make(validator.Input{"id": "550e8400-e29b-41d4-a716-446655440000"}, validator.Rules{"id": "required|uuid"})
	if v.Fails() {
		t.Errorf("uuid should pass, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"id": "not-a-uuid"}, validator.Rules{"id": "required|uuid"})
	if v2.Passes() {
		t.Error("uuid should fail for invalid value")
	}
}

func TestULID(t *testing.T) {
	v := validator.Make(validator.Input{"id": "01ARZ3NDEKTSV4RRFFQ69G5FAV"}, validator.Rules{"id": "required|ulid"})
	if v.Fails() {
		t.Errorf("ulid should pass, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"id": "not-a-ulid"}, validator.Rules{"id": "required|ulid"})
	if v2.Passes() {
		t.Error("ulid should fail for invalid value")
	}
}

func TestHexColor(t *testing.T) {
	cases := []struct {
		val   string
		fails bool
	}{
		{"#fff", false},
		{"#ff5733", false},
		{"#FF5733", false},
		{"#ff573300", false},
		{"red", true},
		{"ff5733", true},
	}
	for _, tt := range cases {
		v := validator.Make(validator.Input{"color": tt.val}, validator.Rules{"color": "required|hex_color"})
		if v.Fails() != tt.fails {
			t.Errorf("hex_color(%q) Fails()=%v want %v", tt.val, v.Fails(), tt.fails)
		}
	}
}

func TestSame(t *testing.T) {
	input := validator.Input{"password": "secret", "password_repeat": "secret"}
	v := validator.Make(input, validator.Rules{"password_repeat": "required|same:password"})
	if v.Fails() {
		t.Errorf("same should pass, got: %v", v.Errors())
	}
	input2 := validator.Input{"password": "secret", "password_repeat": "different"}
	v2 := validator.Make(input2, validator.Rules{"password_repeat": "required|same:password"})
	if v2.Passes() {
		t.Error("same should fail when values differ")
	}
}

func TestDifferent(t *testing.T) {
	input := validator.Input{"new_email": "new@x.com", "old_email": "old@x.com"}
	v := validator.Make(input, validator.Rules{"new_email": "required|different:old_email"})
	if v.Fails() {
		t.Errorf("different should pass, got: %v", v.Errors())
	}
	input2 := validator.Input{"new_email": "same@x.com", "old_email": "same@x.com"}
	v2 := validator.Make(input2, validator.Rules{"new_email": "required|different:old_email"})
	if v2.Passes() {
		t.Error("different should fail when values are equal")
	}
}

func TestConfirmed(t *testing.T) {
	rules := validator.Rules{"password": "required|min:8|confirmed"}

	v := validator.Make(validator.Input{
		"password":              "secret123",
		"password_confirmation": "wrong",
	}, rules)
	if v.Passes() {
		t.Error("confirmed should fail when passwords don't match")
	}

	v2 := validator.Make(validator.Input{
		"password":              "secret123",
		"password_confirmation": "secret123",
	}, rules)
	if v2.Fails() {
		t.Errorf("confirmed should pass when passwords match, got: %v", v2.Errors())
	}
}
