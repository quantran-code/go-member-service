package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestsAreLoggedWithTheRequestID(t *testing.T) {
	t.Run("request carrying a request id", func(t *testing.T) {
		var buf bytes.Buffer
		previous := log.Writer()
		log.SetOutput(&buf)
		defer log.SetOutput(previous)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		req.Header.Set(RequestIDHeader, "abc")
		rec := httptest.NewRecorder()
		NewRouter().ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status %d", rec.Code)
		}
		if got := rec.Header().Get(RequestIDHeader); got != "abc" {
			t.Fatalf("response header %s = %q, want %q", RequestIDHeader, got, "abc")
		}

		logged := buf.String()
		if n := len(strings.Split(strings.TrimRight(logged, "\n"), "\n")); n != 1 {
			t.Fatalf("want exactly one log line, got %d: %q", n, logged)
		}
		if !strings.Contains(logged, "request_id=abc") {
			t.Fatalf("log line missing request id: %q", logged)
		}
		if !strings.Contains(logged, "method=GET") {
			t.Fatalf("log line missing method: %q", logged)
		}
		if !strings.Contains(logged, "path=/health") {
			t.Fatalf("log line missing path: %q", logged)
		}
	})

	t.Run("request without a request id", func(t *testing.T) {
		var buf bytes.Buffer
		previous := log.Writer()
		log.SetOutput(&buf)
		defer log.SetOutput(previous)

		rec := httptest.NewRecorder()
		NewRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))

		if rec.Code != http.StatusOK {
			t.Fatalf("status %d", rec.Code)
		}
		if got := rec.Header().Get(RequestIDHeader); got != "" {
			t.Fatalf("response header %s = %q, want it to be absent", RequestIDHeader, got)
		}

		logged := buf.String()
		if n := len(strings.Split(strings.TrimRight(logged, "\n"), "\n")); n != 1 {
			t.Fatalf("want exactly one log line, got %d: %q", n, logged)
		}
		if !strings.Contains(logged, "request_id=") {
			t.Fatalf("log line missing request id field: %q", logged)
		}
		if !strings.Contains(logged, "method=GET") {
			t.Fatalf("log line missing method: %q", logged)
		}
	})
}
