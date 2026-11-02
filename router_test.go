package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestRouter(t *testing.T) http.Handler {
	t.Helper()
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	return NewRouter(db)
}

func do(t *testing.T, router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestCreateThenReadAMember(t *testing.T) {
	router := newTestRouter(t)
	rec := do(t, router, http.MethodPost, "/members", `{"name":"Ada","email":"Ada@Example.com"}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: %d %s", rec.Code, rec.Body.String())
	}
	var member Member
	if err := json.Unmarshal(rec.Body.Bytes(), &member); err != nil {
		t.Fatal(err)
	}
	if member.Email != "ada@example.com" {
		t.Fatalf("email was not normalised: %q", member.Email)
	}
	rec = do(t, router, http.MethodGet, "/members/1", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("get: %d", rec.Code)
	}
}

func TestADuplicateEmailIs409(t *testing.T) {
	router := newTestRouter(t)
	if rec := do(t, router, http.MethodPost, "/members", `{"name":"Ada","email":"ada@example.com"}`); rec.Code != 201 {
		t.Fatalf("first create: %d", rec.Code)
	}
	rec := do(t, router, http.MethodPost, "/members", `{"name":"Grace","email":"ada@example.com"}`)
	if rec.Code != http.StatusConflict {
		t.Fatalf("second create: %d %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "email already exists") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestValidationAndMissingMembers(t *testing.T) {
	router := newTestRouter(t)
	for body, want := range map[string]string{
		`{"name":"Ada"}`:                "email is required",
		`{"email":"ada@example.com"}`:   "name is required",
		`{"name":"Ada","email":"nope"}`: "email must be a valid address",
	} {
		rec := do(t, router, http.MethodPost, "/members", body)
		if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), want) {
			t.Fatalf("%s: %d %s", body, rec.Code, rec.Body.String())
		}
	}
	if rec := do(t, router, http.MethodGet, "/members/42", ""); rec.Code != http.StatusNotFound {
		t.Fatalf("missing member: %d", rec.Code)
	}
	if rec := do(t, router, http.MethodGet, "/members/abc", ""); rec.Code != http.StatusBadRequest {
		t.Fatalf("bad id: %d", rec.Code)
	}
}

func TestUpdateAndDelete(t *testing.T) {
	router := newTestRouter(t)
	do(t, router, http.MethodPost, "/members", `{"name":"Ada","email":"ada@example.com"}`)

	rec := do(t, router, http.MethodPut, "/members/1", `{"name":"Ada Lovelace","email":"ada@example.com"}`)
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Ada Lovelace") {
		t.Fatalf("update: %d %s", rec.Code, rec.Body.String())
	}
	if rec := do(t, router, http.MethodDelete, "/members/1", ""); rec.Code != http.StatusNoContent {
		t.Fatalf("delete: %d", rec.Code)
	}
	if rec := do(t, router, http.MethodDelete, "/members/1", ""); rec.Code != http.StatusNotFound {
		t.Fatalf("second delete: %d", rec.Code)
	}
}
