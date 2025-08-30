package security

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHeaderSharedTokenVerifier_Verify_Configured(t *testing.T) {
    v := NewHeaderSharedTokenVerifier("tok")
    r := httptest.NewRequest(http.MethodGet, "/", nil)
    // missing header -> error
    if err := v.Verify(r); err == nil {
        t.Fatalf("expected error for missing/incorrect header")
    }
    r.Header.Set(v.HeaderName, "tok")
    if err := v.Verify(r); err != nil {
        t.Fatalf("expected success, got %v", err)
    }
}

func TestHeaderSharedTokenVerifier_Verify_NotConfigured(t *testing.T) {
    v := NewHeaderSharedTokenVerifier("")
    r := httptest.NewRequest(http.MethodGet, "/", nil)
    if err := v.Verify(r); err != nil {
        t.Fatalf("expected allow when not configured, got %v", err)
    }
}

