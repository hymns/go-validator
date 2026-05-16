package validator_test

import (
	"testing"

	validator "github.com/hymns/go-validator"
)

func TestFilled(t *testing.T) {
	// present and empty → fail
	v := validator.Make(validator.Input{"bio": ""}, validator.Rules{"bio": "filled"})
	if v.Passes() {
		t.Error("filled should fail for empty string when key present")
	}
	// absent → pass (not a required field)
	v2 := validator.Make(validator.Input{}, validator.Rules{"bio": "filled"})
	if v2.Fails() {
		t.Errorf("filled should pass when key absent, got: %v", v2.Errors())
	}
	// present and non-empty → pass
	v3 := validator.Make(validator.Input{"bio": "Hello"}, validator.Rules{"bio": "filled"})
	if v3.Fails() {
		t.Errorf("filled should pass for non-empty value, got: %v", v3.Errors())
	}
}

func TestPresent(t *testing.T) {
	// key exists with empty value → pass
	v := validator.Make(validator.Input{"bio": ""}, validator.Rules{"bio": "present"})
	if v.Fails() {
		t.Errorf("present should pass when key exists (even empty), got: %v", v.Errors())
	}
	// key exists with nil → pass
	v2 := validator.Make(validator.Input{"bio": nil}, validator.Rules{"bio": "present"})
	if v2.Fails() {
		t.Errorf("present should pass when key exists with nil, got: %v", v2.Errors())
	}
	// key absent → fail
	v3 := validator.Make(validator.Input{}, validator.Rules{"bio": "present"})
	if v3.Passes() {
		t.Error("present should fail when key missing")
	}
}

func TestProhibited(t *testing.T) {
	// absent → pass
	v := validator.Make(validator.Input{}, validator.Rules{"admin": "prohibited"})
	if v.Fails() {
		t.Errorf("prohibited should pass when field absent, got: %v", v.Errors())
	}
	// present with value → fail
	v2 := validator.Make(validator.Input{"admin": "true"}, validator.Rules{"admin": "prohibited"})
	if v2.Passes() {
		t.Error("prohibited should fail when field has value")
	}
	// present but empty → pass (empty is considered absent)
	v3 := validator.Make(validator.Input{"admin": ""}, validator.Rules{"admin": "prohibited"})
	if v3.Fails() {
		t.Errorf("prohibited should pass when field is empty string, got: %v", v3.Errors())
	}
}

func TestProhibitedIf(t *testing.T) {
	// condition met, field present → fail
	v := validator.Make(
		validator.Input{"role": "guest", "admin": "yes"},
		validator.Rules{"admin": "prohibited_if:role,guest"},
	)
	if v.Passes() {
		t.Error("prohibited_if should fail when condition met and field present")
	}
	// condition not met, field present → pass
	v2 := validator.Make(
		validator.Input{"role": "admin", "admin": "yes"},
		validator.Rules{"admin": "prohibited_if:role,guest"},
	)
	if v2.Fails() {
		t.Errorf("prohibited_if should pass when condition not met, got: %v", v2.Errors())
	}
	// condition met, field absent → pass
	v3 := validator.Make(
		validator.Input{"role": "guest"},
		validator.Rules{"admin": "prohibited_if:role,guest"},
	)
	if v3.Fails() {
		t.Errorf("prohibited_if should pass when condition met but field absent, got: %v", v3.Errors())
	}
}

func TestProhibitedUnless(t *testing.T) {
	// condition not met (role != admin), field present → fail
	v := validator.Make(
		validator.Input{"role": "user", "admin_code": "secret"},
		validator.Rules{"admin_code": "prohibited_unless:role,admin"},
	)
	if v.Passes() {
		t.Error("prohibited_unless should fail when condition not met and field present")
	}
	// condition met (role == admin), field present → pass
	v2 := validator.Make(
		validator.Input{"role": "admin", "admin_code": "secret"},
		validator.Rules{"admin_code": "prohibited_unless:role,admin"},
	)
	if v2.Fails() {
		t.Errorf("prohibited_unless should pass when condition met, got: %v", v2.Errors())
	}
	// condition not met, field absent → pass
	v3 := validator.Make(
		validator.Input{"role": "user"},
		validator.Rules{"admin_code": "prohibited_unless:role,admin"},
	)
	if v3.Fails() {
		t.Errorf("prohibited_unless should pass when field absent, got: %v", v3.Errors())
	}
}
