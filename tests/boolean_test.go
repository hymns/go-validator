package validator_test

import (
	"testing"

	validator "github.com/hymns/go-validator"
)

func TestBoolean(t *testing.T) {
	v := validator.Make(validator.Input{"active": true}, validator.Rules{"active": "required|boolean"})
	if v.Fails() {
		t.Errorf("boolean true should pass, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"active": false}, validator.Rules{"active": "required|boolean"})
	if v2.Fails() {
		t.Errorf("boolean false should pass, got: %v", v2.Errors())
	}
	v3 := validator.Make(validator.Input{"active": "yes"}, validator.Rules{"active": "required|boolean"})
	if v3.Passes() {
		t.Error("boolean should fail for string 'yes' (strict Go bool required)")
	}
}

func TestAccepted(t *testing.T) {
	passing := []any{"yes", "on", "1", "true", true, 1}
	for _, val := range passing {
		v := validator.Make(validator.Input{"terms": val}, validator.Rules{"terms": "accepted"})
		if v.Fails() {
			t.Errorf("accepted: %v (%T) should pass, got: %v", val, val, v.Errors())
		}
	}
	failing := []any{"no", "off", "0", "false", false, 0}
	for _, val := range failing {
		v := validator.Make(validator.Input{"terms": val}, validator.Rules{"terms": "accepted"})
		if v.Passes() {
			t.Errorf("accepted: %v (%T) should fail", val, val)
		}
	}
}

func TestDeclined(t *testing.T) {
	passing := []any{"no", "off", "0", "false", false, 0}
	for _, val := range passing {
		v := validator.Make(validator.Input{"active": val}, validator.Rules{"active": "declined"})
		if v.Fails() {
			t.Errorf("declined: %v (%T) should pass, got: %v", val, val, v.Errors())
		}
	}
	failing := []any{"yes", "on", "1", "true", true, 1}
	for _, val := range failing {
		v := validator.Make(validator.Input{"active": val}, validator.Rules{"active": "declined"})
		if v.Passes() {
			t.Errorf("declined: %v (%T) should fail", val, val)
		}
	}
}

func TestAcceptedIf(t *testing.T) {
	// condition met, value not accepted → fail
	v := validator.Make(
		validator.Input{"is_adult": "yes", "terms": "no"},
		validator.Rules{"terms": "accepted_if:is_adult,yes"},
	)
	if v.Passes() {
		t.Error("accepted_if should fail when condition met but not accepted")
	}

	// condition met, value accepted → pass
	v2 := validator.Make(
		validator.Input{"is_adult": "yes", "terms": "yes"},
		validator.Rules{"terms": "accepted_if:is_adult,yes"},
	)
	if v2.Fails() {
		t.Errorf("accepted_if should pass when condition met and accepted, got: %v", v2.Errors())
	}

	// condition not met → pass regardless
	v3 := validator.Make(
		validator.Input{"is_adult": "no", "terms": "no"},
		validator.Rules{"terms": "accepted_if:is_adult,yes"},
	)
	if v3.Fails() {
		t.Errorf("accepted_if should pass when condition not met, got: %v", v3.Errors())
	}
}

func TestDeclinedIf(t *testing.T) {
	// condition not met → pass regardless
	v := validator.Make(
		validator.Input{"plan": "free", "ads": "true"},
		validator.Rules{"ads": "declined_if:plan,premium"},
	)
	if v.Fails() {
		t.Errorf("declined_if should pass when condition not met, got: %v", v.Errors())
	}

	// condition met, value not declined → fail
	v2 := validator.Make(
		validator.Input{"plan": "premium", "ads": "true"},
		validator.Rules{"ads": "declined_if:plan,premium"},
	)
	if v2.Passes() {
		t.Error("declined_if should fail when condition met and not declined")
	}

	// condition met, value declined → pass
	v3 := validator.Make(
		validator.Input{"plan": "premium", "ads": "no"},
		validator.Rules{"ads": "declined_if:plan,premium"},
	)
	if v3.Fails() {
		t.Errorf("declined_if should pass when condition met and declined, got: %v", v3.Errors())
	}
}
