package roxy

import "testing"

func TestIPAllowed(t *testing.T) {
	tests := map[string][]string{
		"192.168.0.1":   {"192.168.0.1"},
		"192.168.0.10":  {"192.168.0.1/24"},
		"::ffff:c0a8:1": {"192.168.0.1/24"}, // IP6 format
	}

	for remoteAddr, allowList := range tests {
		if !IPAllowed(remoteAddr, allowList) {
			t.Errorf("IP Address: %s was expected to pass through IP Filter", remoteAddr)
		}
	}
}

func TestIPDenied(t *testing.T) {
	tests := map[string][]string{
		"192.168.0.1":  {},
		"192.168.0.10": {"192.169.0.1/24"},
	}

	for remoteAddr, allowList := range tests {
		if IPAllowed(remoteAddr, allowList) {
			t.Errorf("IP Address: %s should not pass through IP Filter", remoteAddr)
		}
	}
}
