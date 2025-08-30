package security

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "testing"
)

func TestLineSignatureVerifier_Verify(t *testing.T) {
    secret := "shhhh"
    body := []byte("hello world")
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write(body)
    sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))

    v := NewLineSignatureVerifier(secret)

    if err := v.Verify(sig, body); err != nil {
        t.Fatalf("expected valid signature, got %v", err)
    }

    if err := v.Verify("invalid", body); err == nil {
        t.Fatalf("expected invalid signature error")
    }
}

func TestLineSignatureVerifier_NotConfigured(t *testing.T) {
    v := NewLineSignatureVerifier("")
    if err := v.Verify("anything", []byte("body")); err == nil || err.Error() != "not configured" {
        t.Fatalf("expected not configured error, got %v", err)
    }
}

