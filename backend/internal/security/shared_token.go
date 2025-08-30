package security

import "net/http"

// SharedTokenVerifier validates shared-token header on incoming requests.
type SharedTokenVerifier interface {
    Verify(r *http.Request) error
}

type HeaderSharedTokenVerifier struct {
    Token      string
    HeaderName string
}

func NewHeaderSharedTokenVerifier(token string) *HeaderSharedTokenVerifier {
    return &HeaderSharedTokenVerifier{Token: token, HeaderName: "X-Shared-Token"}
}

func (v *HeaderSharedTokenVerifier) Verify(r *http.Request) error {
    if v.Token == "" { 
        // not configured -> allow
        return nil
    }
    if r.Header.Get(v.HeaderName) != v.Token {
        // generic error
        return http.ErrNoCookie 
    }
    return nil
}

