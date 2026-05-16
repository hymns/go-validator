package validator_test

import (
	"testing"

	validator "github.com/hymns/go-validator"
)

func TestRequiredIf(t *testing.T) {
	v := validator.Make(
		validator.Input{"role": "admin"},
		validator.Rules{"permissions": "required_if:role,admin"},
	)
	if v.Passes() {
		t.Error("required_if should fail when condition met and field missing")
	}
	v2 := validator.Make(
		validator.Input{"role": "user"},
		validator.Rules{"permissions": "required_if:role,admin"},
	)
	if v2.Fails() {
		t.Errorf("required_if should pass when condition not met, got: %v", v2.Errors())
	}
	v3 := validator.Make(
		validator.Input{"role": "admin", "permissions": "read,write"},
		validator.Rules{"permissions": "required_if:role,admin"},
	)
	if v3.Fails() {
		t.Errorf("required_if should pass when condition met and field present, got: %v", v3.Errors())
	}
}

func TestRequiredUnless(t *testing.T) {
	v := validator.Make(
		validator.Input{"plan": "free"},
		validator.Rules{"card_number": "required_unless:plan,premium"},
	)
	if v.Passes() {
		t.Error("required_unless should fail when unless-condition not met")
	}
	v2 := validator.Make(
		validator.Input{"plan": "premium"},
		validator.Rules{"card_number": "required_unless:plan,premium"},
	)
	if v2.Fails() {
		t.Errorf("required_unless should pass when unless-condition is met, got: %v", v2.Errors())
	}
}

func TestRequiredWith(t *testing.T) {
	// required when any listed sibling is present
	v := validator.Make(
		validator.Input{"city": "KL"},
		validator.Rules{"state": "required_with:city"},
	)
	if v.Passes() {
		t.Error("required_with should fail when sibling is present but field missing")
	}
	v2 := validator.Make(
		validator.Input{},
		validator.Rules{"state": "required_with:city"},
	)
	if v2.Fails() {
		t.Errorf("required_with should pass when sibling absent, got: %v", v2.Errors())
	}
	v3 := validator.Make(
		validator.Input{"city": "KL", "state": "Selangor"},
		validator.Rules{"state": "required_with:city"},
	)
	if v3.Fails() {
		t.Errorf("required_with should pass when both present, got: %v", v3.Errors())
	}
}

func TestRequiredWithout(t *testing.T) {
	// required when any listed sibling is absent
	v := validator.Make(
		validator.Input{},
		validator.Rules{"email": "required_without:phone"},
	)
	if v.Passes() {
		t.Error("required_without should fail when sibling absent and field missing")
	}
	v2 := validator.Make(
		validator.Input{"phone": "0123456789"},
		validator.Rules{"email": "required_without:phone"},
	)
	if v2.Fails() {
		t.Errorf("required_without should pass when sibling present, got: %v", v2.Errors())
	}
}

func TestRequiredWithAll(t *testing.T) {
	// required only when ALL listed siblings present
	v := validator.Make(
		validator.Input{"lat": "3.1", "lng": "101.7"},
		validator.Rules{"location_name": "required_with_all:lat,lng"},
	)
	if v.Passes() {
		t.Error("required_with_all should fail when all siblings present but field missing")
	}
	v2 := validator.Make(
		validator.Input{"lat": "3.1"},
		validator.Rules{"location_name": "required_with_all:lat,lng"},
	)
	if v2.Fails() {
		t.Errorf("required_with_all should pass when not all siblings present, got: %v", v2.Errors())
	}
	v3 := validator.Make(
		validator.Input{"lat": "3.1", "lng": "101.7", "location_name": "KLCC"},
		validator.Rules{"location_name": "required_with_all:lat,lng"},
	)
	if v3.Fails() {
		t.Errorf("required_with_all should pass when all present, got: %v", v3.Errors())
	}
}

func TestRequiredWithoutAll(t *testing.T) {
	// required only when ALL listed siblings absent
	v := validator.Make(
		validator.Input{},
		validator.Rules{"contact": "required_without_all:email,phone"},
	)
	if v.Passes() {
		t.Error("required_without_all should fail when all siblings absent")
	}
	v2 := validator.Make(
		validator.Input{"email": "a@b.com"},
		validator.Rules{"contact": "required_without_all:email,phone"},
	)
	if v2.Fails() {
		t.Errorf("required_without_all should pass when at least one sibling present, got: %v", v2.Errors())
	}
}
