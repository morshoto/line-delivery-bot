package httpapi

import (
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/go-chi/chi/v5"
)

// test doubles
type stubStore struct{ seen map[string]bool }

func (s *stubStore) Seen(key string) bool {
    if s.seen == nil {
        s.seen = map[string]bool{}
    }
    v := s.seen[key]
    s.seen[key] = true
    return v
}

type stubParser struct{ c, t string }

func (p stubParser) Parse(text string) (string, string) { return p.c, p.t }

type stubPusher struct {
    group string
    text  string
    err   error
}

func (p *stubPusher) Push(g, t string) error { p.group, p.text = g, t; return p.err }

type stubSharedVerifier struct{ err error }

func (v stubSharedVerifier) Verify(r *http.Request) error { return v.err }

type stubSigVerifier struct{ err error }

func (v stubSigVerifier) Verify(string, []byte) error { return v.err }

type nopLogger struct{}

func (nopLogger) Info(map[string]any)  {}
func (nopLogger) Warn(map[string]any)  {}
func (nopLogger) Error(map[string]any) {}

func TestHandleScan_SuccessAndRescan(t *testing.T) {
    store := &stubStore{}
    pusher := &stubPusher{}
    h := &Handler{
        Dedupe:         store,
        Parser:         stubParser{c: "yamato", t: "1234567890"},
        Pusher:         pusher,
        SharedVerifier: stubSharedVerifier{err: nil},
        Logger:         nopLogger{},
    }
    r := chi.NewRouter()
    h.Register(r)

    // first scan
    body1 := map[string]string{
        "group_id":     "G",
        "qr_text":      "text",
        "display_name": "Alice",
    }
    b, _ := json.Marshal(body1)
    req := httptest.NewRequest(http.MethodPost, "/api/scan", bytes.NewReader(b))
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("status: got %d body=%s", w.Code, w.Body.String())
    }
    if pusher.group != "G" || !strings.Contains(pusher.text, "配送スキャン") || !strings.Contains(pusher.text, "伝票番号: 1234567890") || !strings.Contains(pusher.text, "配送会社: yamato") {
        t.Fatalf("pushed text unexpected: %q", pusher.text)
    }
    var resp1 ScanResponse
    if err := json.Unmarshal(w.Body.Bytes(), &resp1); err != nil || resp1.Rescan {
        t.Fatalf("unexpected response: %+v err=%v", resp1, err)
    }

    // second scan triggers rescan
    req2 := httptest.NewRequest(http.MethodPost, "/api/scan", bytes.NewReader(b))
    w2 := httptest.NewRecorder()
    r.ServeHTTP(w2, req2)
    if w2.Code != http.StatusOK {
        t.Fatalf("status: got %d", w2.Code)
    }
    if !strings.Contains(pusher.text, "（再スキャン）") {
        t.Fatalf("expected rescan tag in message: %q", pusher.text)
    }
    var resp2 ScanResponse
    _ = json.Unmarshal(w2.Body.Bytes(), &resp2)
    if !resp2.Rescan {
        t.Fatalf("expected rescan true, got %+v", resp2)
    }
}

func TestHandleScan_Errors(t *testing.T) {
    // shared token failure -> 403
    h := &Handler{SharedVerifier: stubSharedVerifier{err: errors.New("nope")}, Logger: nopLogger{}}
    r := chi.NewRouter()
    h.Register(r)
    req := httptest.NewRequest(http.MethodPost, "/api/scan", strings.NewReader(`{}`))
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusForbidden {
        t.Fatalf("expected 403, got %d", w.Code)
    }

    // bad json -> 400
    h = &Handler{SharedVerifier: stubSharedVerifier{err: nil}, Logger: nopLogger{}}
    r = chi.NewRouter()
    h.Register(r)
    req = httptest.NewRequest(http.MethodPost, "/api/scan", strings.NewReader("not json"))
    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest {
        t.Fatalf("expected 400, got %d", w.Code)
    }

    // missing fields -> 400
    req = httptest.NewRequest(http.MethodPost, "/api/scan", strings.NewReader(`{"group_id":"","qr_text":"","display_name":""}`))
    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadRequest {
        t.Fatalf("expected 400, got %d", w.Code)
    }

    // push fail -> 502
    h = &Handler{
        Dedupe:         &stubStore{},
        Parser:         stubParser{c: "yamato", t: "123"},
        Pusher:         &stubPusher{err: errors.New("push fail")},
        SharedVerifier: stubSharedVerifier{err: nil},
        Logger:         nopLogger{},
    }
    r = chi.NewRouter()
    h.Register(r)
    req = httptest.NewRequest(http.MethodPost, "/api/scan", strings.NewReader(`{"group_id":"g","qr_text":"q","display_name":"d"}`))
    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusBadGateway {
        t.Fatalf("expected 502, got %d", w.Code)
    }
}

func TestHandleCallback(t *testing.T) {
    // success
    h := &Handler{SigVerifier: stubSigVerifier{err: nil}, Logger: nopLogger{}}
    r := chi.NewRouter()
    h.Register(r)
    req := httptest.NewRequest(http.MethodPost, "/callback", strings.NewReader("body"))
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK || strings.TrimSpace(w.Body.String()) != "ok" {
        t.Fatalf("expected 200 ok, got %d %q", w.Code, w.Body.String())
    }

    // not configured -> 503
    h = &Handler{SigVerifier: stubSigVerifier{err: errors.New("not configured")}, Logger: nopLogger{}}
    r = chi.NewRouter()
    h.Register(r)
    req = httptest.NewRequest(http.MethodPost, "/callback", strings.NewReader("b"))
    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusServiceUnavailable {
        t.Fatalf("expected 503, got %d", w.Code)
    }

    // other error -> 403
    h = &Handler{SigVerifier: stubSigVerifier{err: errors.New("bad")}, Logger: nopLogger{}}
    r = chi.NewRouter()
    h.Register(r)
    req = httptest.NewRequest(http.MethodPost, "/callback", strings.NewReader("b"))
    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusForbidden {
        t.Fatalf("expected 403, got %d", w.Code)
    }
}

