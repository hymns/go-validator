package validator_test

import (
	"testing"

	validator "github.com/hymns/go-validator"
)

func TestUniqueNoDB(t *testing.T) {
	cases := []struct {
		name string
		rule string
	}{
		{"basic", "required|email|unique:users,email"},
		{"with ignore id", "required|email|unique:users,email,42"},
		{"with ignore col", "required|email|unique:users,email,42,user_id"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.Make(
				validator.Input{"email": "test@example.com"},
				validator.Rules{"email": tt.rule},
			)
			if v.Passes() {
				t.Error("unique without DB should fail with db_required message")
			}
			if !v.Errors().Has("email") {
				t.Error("expected error on email field")
			}
		})
	}
}

func TestExistsNoDB(t *testing.T) {
	cases := []struct {
		name string
		rule string
	}{
		{"basic", "required|exists:roles,id"},
		{"with ignore", "required|exists:roles,id,99"},
		{"with ignore col", "required|exists:roles,id,99,parent_id"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.Make(
				validator.Input{"role_id": "1"},
				validator.Rules{"role_id": tt.rule},
			)
			if v.Passes() {
				t.Error("exists without DB should fail with db_required message")
			}
			if !v.Errors().Has("role_id") {
				t.Error("expected error on role_id field")
			}
		})
	}
}
