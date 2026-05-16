package validator_test

import (
	"testing"

	validator "github.com/hymns/go-validator"
)

func TestInteger(t *testing.T) {
	cases := []struct {
		val   any
		fails bool
	}{
		{42, false},
		{float64(42), false},
		{float64(42.5), true},
		{"99", false},
		{"99.9", true},
		{true, true},
	}
	for _, tt := range cases {
		v := validator.Make(validator.Input{"n": tt.val}, validator.Rules{"n": "required|integer"})
		if v.Fails() != tt.fails {
			t.Errorf("integer(%v) Fails()=%v want %v, errors: %v", tt.val, v.Fails(), tt.fails, v.Errors())
		}
	}
}

func TestNumeric(t *testing.T) {
	cases := []struct {
		val   any
		fails bool
	}{
		{42, false},
		{3.14, false},
		{"42.5", false},
		{"hello", true},
		{true, true},
	}
	for _, tt := range cases {
		v := validator.Make(validator.Input{"val": tt.val}, validator.Rules{"val": "required|numeric"})
		if v.Fails() != tt.fails {
			t.Errorf("numeric(%v) Fails()=%v want %v, errors: %v", tt.val, v.Fails(), tt.fails, v.Errors())
		}
	}
}

func TestMinMaxNumeric(t *testing.T) {
	rules := validator.Rules{"age": "required|integer|min:18|max:100"}
	cases := []struct {
		val   float64
		fails bool
	}{
		{17, true},
		{18, false},
		{100, false},
		{101, true},
	}
	for _, tt := range cases {
		v := validator.Make(validator.Input{"age": tt.val}, rules)
		if v.Fails() != tt.fails {
			t.Errorf("age=%v Fails()=%v want %v, errors: %v", tt.val, v.Fails(), tt.fails, v.Errors())
		}
	}
}

func TestMinMaxDigits(t *testing.T) {
	v := validator.Make(validator.Input{"code": "12"}, validator.Rules{"code": "required|min_digits:4"})
	if v.Passes() {
		t.Error("min_digits:4 should fail for '12'")
	}
	v2 := validator.Make(validator.Input{"code": "1234"}, validator.Rules{"code": "required|min_digits:4|max_digits:6"})
	if v2.Fails() {
		t.Errorf("min_digits:4|max_digits:6 should pass for '1234', got: %v", v2.Errors())
	}
	v3 := validator.Make(validator.Input{"code": "1234567"}, validator.Rules{"code": "required|max_digits:6"})
	if v3.Passes() {
		t.Error("max_digits:6 should fail for 7-digit value")
	}
}

func TestBetweenNumeric(t *testing.T) {
	v := validator.Make(validator.Input{"age": float64(25)}, validator.Rules{"age": "required|between:18,60"})
	if v.Fails() {
		t.Errorf("between should pass, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"age": float64(10)}, validator.Rules{"age": "required|between:18,60"})
	if v2.Passes() {
		t.Error("between should fail for 10")
	}
	v3 := validator.Make(validator.Input{"age": float64(60)}, validator.Rules{"age": "required|between:18,60"})
	if v3.Fails() {
		t.Errorf("between should pass for max boundary, got: %v", v3.Errors())
	}
}

func TestBetweenStringLength(t *testing.T) {
	v := validator.Make(validator.Input{"name": "Ali"}, validator.Rules{"name": "required|between:3,10"})
	if v.Fails() {
		t.Errorf("between string len should pass, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"name": "Al"}, validator.Rules{"name": "required|between:3,10"})
	if v2.Passes() {
		t.Error("between string len should fail for too-short string")
	}
}

func TestSizeString(t *testing.T) {
	v := validator.Make(validator.Input{"pin": "1234"}, validator.Rules{"pin": "required|size:4"})
	if v.Fails() {
		t.Errorf("size:4 should pass for '1234', got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"pin": "123"}, validator.Rules{"pin": "required|size:4"})
	if v2.Passes() {
		t.Error("size:4 should fail for 3-char pin")
	}
}

func TestSizeNumeric(t *testing.T) {
	v := validator.Make(validator.Input{"score": float64(100)}, validator.Rules{"score": "required|size:100"})
	if v.Fails() {
		t.Errorf("size numeric should pass, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"score": float64(99)}, validator.Rules{"score": "required|size:100"})
	if v2.Passes() {
		t.Error("size numeric should fail for wrong value")
	}
}

func TestDigits(t *testing.T) {
	v := validator.Make(validator.Input{"otp": "123456"}, validator.Rules{"otp": "required|digits:6"})
	if v.Fails() {
		t.Errorf("digits:6 should pass, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"otp": "12345"}, validator.Rules{"otp": "required|digits:6"})
	if v2.Passes() {
		t.Error("digits:6 should fail for 5-digit value")
	}
	v3 := validator.Make(validator.Input{"otp": "12345a"}, validator.Rules{"otp": "required|digits:6"})
	if v3.Passes() {
		t.Error("digits:6 should fail when non-digit chars present")
	}
}

func TestDigitsBetween(t *testing.T) {
	v := validator.Make(validator.Input{"n": "12"}, validator.Rules{"n": "required|digits_between:2,5"})
	if v.Fails() {
		t.Errorf("digits_between:2,5 should pass for '12', got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"n": "1"}, validator.Rules{"n": "required|digits_between:2,5"})
	if v2.Passes() {
		t.Error("digits_between:2,5 should fail for 1-digit value")
	}
	v3 := validator.Make(validator.Input{"n": "123456"}, validator.Rules{"n": "required|digits_between:2,5"})
	if v3.Passes() {
		t.Error("digits_between:2,5 should fail for 6-digit value")
	}
}

func TestMultipleOf(t *testing.T) {
	v := validator.Make(validator.Input{"qty": float64(15)}, validator.Rules{"qty": "required|multiple_of:5"})
	if v.Fails() {
		t.Errorf("multiple_of:5 should pass for 15, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"qty": float64(7)}, validator.Rules{"qty": "required|multiple_of:5"})
	if v2.Passes() {
		t.Error("multiple_of:5 should fail for 7")
	}
	v3 := validator.Make(validator.Input{"qty": float64(0)}, validator.Rules{"qty": "required|multiple_of:5"})
	if v3.Fails() {
		t.Errorf("multiple_of:5 should pass for 0, got: %v", v3.Errors())
	}
}

func TestGT(t *testing.T) {
	v := validator.Make(validator.Input{"n": float64(11)}, validator.Rules{"n": "required|gt:10"})
	if v.Fails() {
		t.Errorf("gt:10 should pass for 11, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"n": float64(10)}, validator.Rules{"n": "required|gt:10"})
	if v2.Passes() {
		t.Error("gt:10 should fail for 10 (not strictly greater)")
	}
}

func TestGTE(t *testing.T) {
	v := validator.Make(validator.Input{"n": float64(10)}, validator.Rules{"n": "required|gte:10"})
	if v.Fails() {
		t.Errorf("gte:10 should pass for 10, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"n": float64(9)}, validator.Rules{"n": "required|gte:10"})
	if v2.Passes() {
		t.Error("gte:10 should fail for 9")
	}
}

func TestLT(t *testing.T) {
	v := validator.Make(validator.Input{"n": float64(9)}, validator.Rules{"n": "required|lt:10"})
	if v.Fails() {
		t.Errorf("lt:10 should pass for 9, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"n": float64(10)}, validator.Rules{"n": "required|lt:10"})
	if v2.Passes() {
		t.Error("lt:10 should fail for 10")
	}
}

func TestLTE(t *testing.T) {
	v := validator.Make(validator.Input{"n": float64(10)}, validator.Rules{"n": "required|lte:10"})
	if v.Fails() {
		t.Errorf("lte:10 should pass for 10, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"n": float64(11)}, validator.Rules{"n": "required|lte:10"})
	if v2.Passes() {
		t.Error("lte:10 should fail for 11")
	}
}

func TestGTWithField(t *testing.T) {
	input := validator.Input{"min_age": float64(18), "max_age": float64(25)}
	v := validator.Make(input, validator.Rules{"max_age": "required|gt:min_age"})
	if v.Fails() {
		t.Errorf("gt:field should pass when max_age > min_age, got: %v", v.Errors())
	}
	input2 := validator.Input{"min_age": float64(30), "max_age": float64(25)}
	v2 := validator.Make(input2, validator.Rules{"max_age": "required|gt:min_age"})
	if v2.Passes() {
		t.Error("gt:field should fail when max_age < min_age")
	}
}
