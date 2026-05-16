package validator_test

import (
	"testing"

	validator "github.com/hymns/go-validator"
)

func TestIP(t *testing.T) {
	cases := []struct {
		val   string
		fails bool
	}{
		{"192.168.1.1", false},
		{"::1", false},
		{"2001:db8::1", false},
		{"not-an-ip", true},
		{"999.999.999.999", true},
	}
	for _, tt := range cases {
		v := validator.Make(validator.Input{"addr": tt.val}, validator.Rules{"addr": "required|ip"})
		if v.Fails() != tt.fails {
			t.Errorf("ip=%q Fails()=%v want %v", tt.val, v.Fails(), tt.fails)
		}
	}
}

func TestIPv4(t *testing.T) {
	cases := []struct {
		val   string
		fails bool
	}{
		{"192.168.1.1", false},
		{"127.0.0.1", false},
		{"::1", true},
		{"2001:db8::1", true},
		{"not-an-ip", true},
	}
	for _, tt := range cases {
		v := validator.Make(validator.Input{"addr": tt.val}, validator.Rules{"addr": "required|ipv4"})
		if v.Fails() != tt.fails {
			t.Errorf("ipv4=%q Fails()=%v want %v", tt.val, v.Fails(), tt.fails)
		}
	}
}

func TestIPv6(t *testing.T) {
	cases := []struct {
		val   string
		fails bool
	}{
		{"::1", false},
		{"2001:db8::1", false},
		{"192.168.1.1", true},
		{"not-an-ip", true},
	}
	for _, tt := range cases {
		v := validator.Make(validator.Input{"addr": tt.val}, validator.Rules{"addr": "required|ipv6"})
		if v.Fails() != tt.fails {
			t.Errorf("ipv6=%q Fails()=%v want %v", tt.val, v.Fails(), tt.fails)
		}
	}
}

func TestMACAddress(t *testing.T) {
	cases := []struct {
		val   string
		fails bool
	}{
		{"00:1A:2B:3C:4D:5E", false},
		{"00-1A-2B-3C-4D-5E", false},
		{"not-a-mac", true},
		{"ZZ:ZZ:ZZ:ZZ:ZZ:ZZ", true},
		{"00:1A:2B:3C:4D", true},
	}
	for _, tt := range cases {
		v := validator.Make(validator.Input{"mac": tt.val}, validator.Rules{"mac": "required|mac_address"})
		if v.Fails() != tt.fails {
			t.Errorf("mac=%q Fails()=%v want %v, errors: %v", tt.val, v.Fails(), tt.fails, v.Errors())
		}
	}
}
