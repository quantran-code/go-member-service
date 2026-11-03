package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetMemberReturnsTheMember(t *testing.T) {
	store := NewMemberStore()
	want, ok := store.Get("1")
	if !ok {
		t.Fatal(`seed data has no member with id "1"`)
	}

	var buf bytes.Buffer
	previous := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(previous)

	req := httptest.NewRequest(http.MethodGet, "/members/1", nil)
	req.Header.Set(RequestIDHeader, "abc")
	rec := httptest.NewRecorder()
	NewRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d, body %s", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}

	var got Member
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decoding body %s: %v", rec.Body.String(), err)
	}
	if got != want {
		t.Fatalf("member = %+v, want %+v", got, want)
	}

	logged := buf.String()
	if !strings.Contains(logged, "request_id=abc") {
		t.Fatalf("log output missing the request id: %q", logged)
	}
	lines := strings.Split(strings.TrimRight(logged, "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("want the handler line and the middleware line, got %d: %q", len(lines), logged)
	}
	for _, line := range lines {
		if !strings.Contains(line, "request_id=abc") {
			t.Fatalf("log line missing the request id: %q", line)
		}
		if !strings.Contains(line, "path=/members/1") {
			t.Fatalf("log line missing the path: %q", line)
		}
	}
}

func TestGetMemberReturns404ForUnknownID(t *testing.T) {
	rec := httptest.NewRecorder()
	NewRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/members/does-not-exist", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status %d, body %s", rec.Code, rec.Body.String())
	}

	var got errorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("decoding body %s: %v", rec.Body.String(), err)
	}
	if got.Error != "member not found" {
		t.Fatalf("error = %q, want %q", got.Error, "member not found")
	}
}
