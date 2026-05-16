package validator_test

import (
	"testing"

	validator "github.com/hymns/go-validator"
)

func TestDate(t *testing.T) {
	cases := []struct {
		val   string
		fails bool
	}{
		{"2024-01-15", false},
		{"2024-01-15 10:30:00", false},
		{"2024-01-15T10:30:00Z", false},
		{"not-a-date", true},
		{"15/01/2024", true},
		{"2024-13-01", true},
	}
	for _, tt := range cases {
		v := validator.Make(validator.Input{"dob": tt.val}, validator.Rules{"dob": "required|date"})
		if v.Fails() != tt.fails {
			t.Errorf("date=%q Fails()=%v want %v, errors: %v", tt.val, v.Fails(), tt.fails, v.Errors())
		}
	}
}

func TestBefore(t *testing.T) {
	v := validator.Make(validator.Input{"start": "2020-01-01"}, validator.Rules{"start": "required|before:2025-01-01"})
	if v.Fails() {
		t.Errorf("before should pass, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"start": "2030-01-01"}, validator.Rules{"start": "required|before:2025-01-01"})
	if v2.Passes() {
		t.Error("before should fail for date after reference")
	}
	v3 := validator.Make(validator.Input{"dob": "1990-06-15"}, validator.Rules{"dob": "required|before:today"})
	if v3.Fails() {
		t.Errorf("before:today should pass for past date, got: %v", v3.Errors())
	}
}

func TestAfter(t *testing.T) {
	v := validator.Make(validator.Input{"expiry": "2099-12-31"}, validator.Rules{"expiry": "required|after:today"})
	if v.Fails() {
		t.Errorf("after:today should pass for future date, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"expiry": "2000-01-01"}, validator.Rules{"expiry": "required|after:today"})
	if v2.Passes() {
		t.Error("after:today should fail for past date")
	}
	v3 := validator.Make(validator.Input{"event": "2030-06-01"}, validator.Rules{"event": "required|after:2025-01-01"})
	if v3.Fails() {
		t.Errorf("after:2025-01-01 should pass for future date, got: %v", v3.Errors())
	}
}

func TestBeforeOrEqual(t *testing.T) {
	v := validator.Make(validator.Input{"d": "2025-01-01"}, validator.Rules{"d": "required|before_or_equal:2025-01-01"})
	if v.Fails() {
		t.Errorf("before_or_equal should pass on equal date, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"d": "2026-01-01"}, validator.Rules{"d": "required|before_or_equal:2025-01-01"})
	if v2.Passes() {
		t.Error("before_or_equal should fail for date after reference")
	}
}

func TestAfterOrEqual(t *testing.T) {
	v := validator.Make(validator.Input{"d": "2025-01-01"}, validator.Rules{"d": "required|after_or_equal:2025-01-01"})
	if v.Fails() {
		t.Errorf("after_or_equal should pass on equal date, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"d": "2024-01-01"}, validator.Rules{"d": "required|after_or_equal:2025-01-01"})
	if v2.Passes() {
		t.Error("after_or_equal should fail for date before reference")
	}
}

func TestDateEquals(t *testing.T) {
	v := validator.Make(validator.Input{"d": "2025-06-01"}, validator.Rules{"d": "required|date_equals:2025-06-01"})
	if v.Fails() {
		t.Errorf("date_equals should pass, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"d": "2025-06-02"}, validator.Rules{"d": "required|date_equals:2025-06-01"})
	if v2.Passes() {
		t.Error("date_equals should fail for different date")
	}
}

func TestDateFormat(t *testing.T) {
	v := validator.Make(validator.Input{"d": "2025-06-01"}, validator.Rules{"d": "required|date_format:2006-01-02"})
	if v.Fails() {
		t.Errorf("date_format should pass, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"d": "01/06/2025"}, validator.Rules{"d": "required|date_format:2006-01-02"})
	if v2.Passes() {
		t.Error("date_format should fail for wrong format")
	}
}

func TestTimezone(t *testing.T) {
	cases := []struct {
		val   string
		fails bool
	}{
		{"Asia/Kuala_Lumpur", false},
		{"UTC", false},
		{"America/New_York", false},
		{"Invalid/Zone", true},
		{"not-a-timezone", true},
	}
	for _, tt := range cases {
		v := validator.Make(validator.Input{"tz": tt.val}, validator.Rules{"tz": "required|timezone"})
		if v.Fails() != tt.fails {
			t.Errorf("timezone=%q Fails()=%v want %v", tt.val, v.Fails(), tt.fails)
		}
	}
}
