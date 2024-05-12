package roxy

import (
	"net/http"
	"testing"
)

func TestGenerateRequestID(t *testing.T) {
	if id := generateRequestID(8); len(id) != 8 {
		t.Fatalf("the request id %s is not 8 characters in length", id)
	}
}

func TestRequestIDFromHeaders(t *testing.T) {
	headers := http.Header{}
	headers.Add("X-B3-TraceId", "1234")

	if id := requestIdFromHeaders(headers); id != "1234" {
		t.Fatalf("Failed to get TraceID from headers when present. Expected '1234', Got: %s", id)
	}

	headers.Del("X-B3-TraceId")
	if id := requestIdFromHeaders(headers); len(id) != 8 {
		t.Fatalf("Expected a randomly generated request id, Got: %s", id)
	}
}
