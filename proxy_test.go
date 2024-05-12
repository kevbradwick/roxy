package roxy

import "testing"

func TestGenerateRequestID(t *testing.T) {
	if id := generateRequestID(8); len(id) != 8 {
		t.Fatalf("the request id %s is not 8 characters in length", id)
	}
}
