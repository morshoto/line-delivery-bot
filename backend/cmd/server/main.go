package main

import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "io"
    "log"
    "net/http"
    "os"
    "strings"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

type scanRequest struct {
    GroupID     string `json:"group_id"`
    QRText      string `json:"qr_text"`
    DisplayName string `json:"display_name"`
}

type scanResponse struct {
    Status   string `json:"status"`
    Rescan   bool   `json:"rescan"`
    Carrier  string `json:"carrier"`
    Tracking string `json:"tracking_no"`
}

// simple in-memory TTL store for dedupe
type ttlStore struct {
    data map[string]time.Time
    ttl  time.Duration
}

func newTTLStore(ttl time.Duration) *ttlStore {
    return &ttlStore{data: make(map[string]time.Time), ttl: ttl}
}

func (s *ttlStore) seen(key string) (dup bool) {
    now := time.Now()
    // cleanup expired
    for k, t := range s.data {
        if now.Sub(t) > s.ttl {
            delete(s.data, k)
        }
    }
    if t, ok := s.data[key]; ok {
        if now.Sub(t) <= s.ttl {
            s.data[key] = now // refresh
            return true
        }
    }
    s.data[key] = now
    return false
}

// very naive parser to guess carrier and tracking number
func parseQR(text string) (carrier, tracking string) {
    lower := strings.ToLower(text)
    switch {
    case strings.Contains(lower, "kuroneko") || strings.Contains(lower, "yamato"):
        carrier = "yamato"
    case strings.Contains(lower, "sagawa"):
        carrier = "sagawa"
    case strings.Contains(lower, "japanpost") || strings.Contains(lower, "post.japanpost.jp"):
        carrier = "japanpost"
    default:
        carrier = "unknown"
    }
    // pick the longest 10-14 digit sequence as tracking candidate
    digits := make([]rune, 0, len(text))
    best := ""
    for _, r := range text {
        if r >= '0' && r <= '9' {
            digits = append(digits, r)
        } else {
            if l := len(digits); l >= 10 && l <= 14 && l > len(best) {
                best = string(digits)
            }
            digits = digits[:0]
        }
    }
    if l := len(digits); l >= 10 && l <= 14 && l > len(best) {
        best = string(digits)
    }
    tracking = best
    return
}

func jsonLog(fields map[string]any) {
    b, _ := json.Marshal(fields)
    log.Println(string(b))
}

func verifySharedToken(r *http.Request) error {
    token := os.Getenv("SHARED_TOKEN")
    if token == "" { // not configured -> allow
        return nil
    }
    if r.Header.Get("X-Shared-Token") != token {
        return http.ErrNoCookie // use as generic error
    }
    return nil
}

func pushMessage(groupID, text string) error {
    token := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
    if token == "" {
        // Not configured, log and skip
        jsonLog(map[string]any{"level": "warn", "event": "push_skip", "reason": "missing_token"})
        return nil
    }
    payload := map[string]any{
        "to": groupID,
        "messages": []map[string]string{{
            "type": "text",
            "text": text,
        }},
    }
    b, _ := json.Marshal(payload)
    req, _ := http.NewRequest(http.MethodPost, "https://api.line.me/v2/bot/message/push", strings.NewReader(string(b)))
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")
    httpClient := &http.Client{Timeout: 10 * time.Second}
    resp, err := httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode >= 400 {
        body, _ := io.ReadAll(resp.Body)
        return &httpError{StatusCode: resp.StatusCode, Body: string(body)}
    }
    return nil
}

type httpError struct {
    StatusCode int
    Body       string
}

func (e *httpError) Error() string { return "http error: " + http.StatusText(e.StatusCode) }

func main() {
    r := chi.NewRouter()
    r.Use(middleware.RealIP, middleware.Logger, middleware.Recoverer)

    // 健康チェック
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    })

    store := newTTLStore(30 * time.Minute)

    // POST /api/scan
    r.Post("/api/scan", func(w http.ResponseWriter, r *http.Request) {
        if err := verifySharedToken(r); err != nil {
            http.Error(w, "forbidden", http.StatusForbidden)
            return
        }
        var req scanRequest
        dec := json.NewDecoder(r.Body)
        dec.DisallowUnknownFields()
        if err := dec.Decode(&req); err != nil {
            http.Error(w, "bad request", http.StatusBadRequest)
            return
        }
        if req.GroupID == "" || req.QRText == "" || req.DisplayName == "" {
            http.Error(w, "missing fields", http.StatusBadRequest)
            return
        }
        carrier, tracking := parseQR(req.QRText)
        key := carrier + ":" + tracking
        isDup := store.seen(key)

        // assemble message
        tag := ""
        if isDup {
            tag = "（再スキャン）"
        }
        lines := []string{
            "配送スキャン" + tag,
        }
        if carrier != "" && carrier != "unknown" {
            lines = append(lines, "配送会社: "+carrier)
        }
        if tracking != "" {
            lines = append(lines, "伝票番号: "+tracking)
        }
        lines = append(lines, "送信者: "+req.DisplayName)
        // include raw QR for transparency
        lines = append(lines, "QR: "+req.QRText)
        message := strings.Join(lines, "\n")

        if err := pushMessage(req.GroupID, message); err != nil {
            jsonLog(map[string]any{
                "level": "error", "event": "push_fail", "error": err.Error(),
            })
            http.Error(w, "failed to push", http.StatusBadGateway)
            return
        }

        jsonLog(map[string]any{
            "event": "scan", "group_id": req.GroupID, "carrier": carrier, "tracking_no": tracking, "dedupe": isDup,
        })

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(scanResponse{
            Status:   "ok",
            Rescan:   isDup,
            Carrier:  carrier,
            Tracking: tracking,
        })
    })

    // POST /callback: signature verification only
    r.Post("/callback", func(w http.ResponseWriter, r *http.Request) {
        secret := os.Getenv("LINE_CHANNEL_SECRET")
        if secret == "" {
            http.Error(w, "not configured", http.StatusServiceUnavailable)
            return
        }
        body, err := io.ReadAll(r.Body)
        if err != nil {
            http.Error(w, "bad request", http.StatusBadRequest)
            return
        }
        sig := r.Header.Get("X-Line-Signature")
        mac := hmac.New(sha256.New, []byte(secret))
        mac.Write(body)
        expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
        if !hmac.Equal([]byte(sig), []byte(expected)) {
            http.Error(w, "forbidden", http.StatusForbidden)
            return
        }
        // no registration flow — just 200 OK
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("ok"))
    })

    // PORT 環境変数（Renderは自動で渡す）
    port := os.Getenv("PORT")
    if port == "" {
        port = "10000"
    }
    addr := ":" + port
    log.Printf("listening on %s", addr)
    if err := http.ListenAndServe(addr, r); err != nil {
        log.Fatal(err)
    }
}
