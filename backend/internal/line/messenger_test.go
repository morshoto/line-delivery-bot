package line

import (
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestHTTPPusher_Push_Success(t *testing.T) {
    var gotAuth string
    var gotPath string
    var gotPayload struct {
        To       string `json:"to"`
        Messages []struct {
            Type string `json:"type"`
            Text string `json:"text"`
        } `json:"messages"`
    }
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        gotAuth = r.Header.Get("Authorization")
        gotPath = r.URL.Path
        if err := json.NewDecoder(r.Body).Decode(&gotPayload); err != nil {
            t.Fatalf("bad json: %v", err)
        }
        w.WriteHeader(http.StatusOK)
    }))
    defer srv.Close()

    p := NewHTTPPusher("token123")
    p.BaseURL = srv.URL
    p.HTTPClient = srv.Client()

    if err := p.Push("group1", "hello"); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if gotAuth != "Bearer token123" {
        t.Fatalf("auth header: got %q", gotAuth)
    }
    if gotPath != "/v2/bot/message/push" {
        t.Fatalf("path: got %q", gotPath)
    }
    if gotPayload.To != "group1" || len(gotPayload.Messages) != 1 || gotPayload.Messages[0].Text != "hello" {
        t.Fatalf("payload mismatch: %+v", gotPayload)
    }
}

func TestHTTPPusher_Push_Errors(t *testing.T) {
    var pNil *HTTPPusher
    if err := pNil.Push("g", "t"); err == nil || !strings.Contains(err.Error(), "nil pusher") {
        t.Fatalf("expected nil pusher error, got %v", err)
    }

    p := NewHTTPPusher("")
    if err := p.Push("g", "t"); err == nil || !strings.Contains(err.Error(), "missing token") {
        t.Fatalf("expected missing token error, got %v", err)
    }

    // server returns error
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        http.Error(w, "boom", http.StatusBadRequest)
    }))
    defer srv.Close()
    p = NewHTTPPusher("tok")
    p.BaseURL = srv.URL
    p.HTTPClient = srv.Client()
    err := p.Push("g", "t")
    var httpErr *HTTPError
    if err == nil || !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusBadRequest {
        t.Fatalf("expected HTTPError 400, got %v", err)
    }
}

