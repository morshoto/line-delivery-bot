package line

import (
    "bytes"
    "encoding/json"
    "errors"
    "io"
    "net/http"
    "time"
)

// Pusher abstracts sending messages to a LINE group.
type Pusher interface {
    Push(groupID, text string) error
}

// HTTPPusher sends messages via LINE Messaging API over HTTP.
type HTTPPusher struct {
    Token      string
    HTTPClient *http.Client
    BaseURL    string
}

func NewHTTPPusher(token string) *HTTPPusher {
    return &HTTPPusher{
        Token:      token,
        HTTPClient: &http.Client{Timeout: 10 * time.Second},
        BaseURL:    "https://api.line.me",
    }
}

func (p *HTTPPusher) Push(groupID, text string) error {
    if p == nil {
        return errors.New("nil pusher")
    }
    if p.Token == "" {
        return errors.New("missing token")
    }
    payload := map[string]any{
        "to": groupID,
        "messages": []map[string]string{{
            "type": "text",
            "text": text,
        }},
    }
    var buf bytes.Buffer
    if err := json.NewEncoder(&buf).Encode(payload); err != nil {
        return err
    }
    req, err := http.NewRequest(http.MethodPost, p.BaseURL+"/v2/bot/message/push", &buf)
    if err != nil {
        return err
    }
    req.Header.Set("Authorization", "Bearer "+p.Token)
    req.Header.Set("Content-Type", "application/json")
    resp, err := p.HTTPClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode >= 400 {
        body, _ := io.ReadAll(resp.Body)
        return &HTTPError{StatusCode: resp.StatusCode, Body: string(body)}
    }
    return nil
}

// NoopPusher implements Pusher but does nothing (useful when token is not configured).
type NoopPusher struct{}

func (NoopPusher) Push(string, string) error { return nil }

type HTTPError struct {
    StatusCode int
    Body       string
}

func (e *HTTPError) Error() string { return http.StatusText(e.StatusCode) }

