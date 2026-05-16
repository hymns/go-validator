package validator_test

import (
	"testing"

	validator "github.com/hymns/go-validator"
)

func TestIn(t *testing.T) {
	rules := validator.Rules{"role": "required|in:admin,user,editor"}

	v := validator.Make(validator.Input{"role": "superadmin"}, rules)
	if v.Passes() {
		t.Error("in should fail for 'superadmin'")
	}
	v2 := validator.Make(validator.Input{"role": "admin"}, rules)
	if v2.Fails() {
		t.Errorf("in should pass for 'admin', got: %v", v2.Errors())
	}
	v3 := validator.Make(validator.Input{"role": "editor"}, rules)
	if v3.Fails() {
		t.Errorf("in should pass for 'editor', got: %v", v3.Errors())
	}
}

func TestNotIn(t *testing.T) {
	rules := validator.Rules{"status": "required|not_in:banned,deleted"}

	v := validator.Make(validator.Input{"status": "banned"}, rules)
	if v.Passes() {
		t.Error("not_in should fail for 'banned'")
	}
	v2 := validator.Make(validator.Input{"status": "active"}, rules)
	if v2.Fails() {
		t.Errorf("not_in should pass for 'active', got: %v", v2.Errors())
	}
}

func TestArray(t *testing.T) {
	v := validator.Make(validator.Input{"tags": []any{"go", "php"}}, validator.Rules{"tags": "required|array"})
	if v.Fails() {
		t.Errorf("array should pass for []any, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"tags": "not-array"}, validator.Rules{"tags": "required|array"})
	if v2.Passes() {
		t.Error("array should fail for string")
	}
	v3 := validator.Make(validator.Input{"tags": map[string]any{"a": 1}}, validator.Rules{"tags": "required|array"})
	if v3.Passes() {
		t.Error("array should fail for map")
	}
}

func TestDistinct(t *testing.T) {
	v := validator.Make(validator.Input{"tags": []any{"go", "php", "python"}}, validator.Rules{"tags": "required|array|distinct"})
	if v.Fails() {
		t.Errorf("distinct should pass for unique values, got: %v", v.Errors())
	}
	v2 := validator.Make(validator.Input{"tags": []any{"go", "php", "go"}}, validator.Rules{"tags": "required|array|distinct"})
	if v2.Passes() {
		t.Error("distinct should fail for duplicate values")
	}
	v3 := validator.Make(validator.Input{"tags": []any{1, 2, 1}}, validator.Rules{"tags": "required|array|distinct"})
	if v3.Passes() {
		t.Error("distinct should fail for duplicate numeric values")
	}
}
