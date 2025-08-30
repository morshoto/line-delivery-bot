package security

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "errors"
)

// SignatureVerifier validates LINE signature for webhook callbacks.
type SignatureVerifier interface {
    Verify(signatureHeader string, body []byte) error
}

type LineSignatureVerifier struct {
    Secret string
}

func NewLineSignatureVerifier(secret string) *LineSignatureVerifier {
    return &LineSignatureVerifier{Secret: secret}
}

func (v *LineSignatureVerifier) Verify(signatureHeader string, body []byte) error {
    if v.Secret == "" {
        return errors.New("not configured")
    }
    mac := hmac.New(sha256.New, []byte(v.Secret))
    mac.Write(body)
    expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
    if !hmac.Equal([]byte(signatureHeader), []byte(expected)) {
        return errors.New("invalid signature")
    }
    return nil
}

