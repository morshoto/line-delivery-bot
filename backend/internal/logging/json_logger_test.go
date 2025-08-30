package logging

import (
    "bytes"
    "encoding/json"
    "log"
    "testing"
)

func captureLog(f func()) string {
    var buf bytes.Buffer
    oldOut := log.Writer()
    oldFlags := log.Flags()
    log.SetOutput(&buf)
    log.SetFlags(0)
    defer func() {
        log.SetOutput(oldOut)
        log.SetFlags(oldFlags)
    }()
    f()
    return buf.String()
}

func TestJSONLogger_Levels(t *testing.T) {
    l := JSONLogger{}
    tests := []struct{ level string }{
        {"info"}, {"warn"}, {"error"},
    }
    for _, tc := range tests {
        out := captureLog(func() {
            switch tc.level {
            case "info":
                l.Info(map[string]any{"msg": "hello"})
            case "warn":
                l.Warn(map[string]any{"msg": "hello"})
            case "error":
                l.Error(map[string]any{"msg": "hello"})
            }
        })
        // Trim newline and parse JSON
        var m map[string]any
        if err := json.Unmarshal(bytes.TrimSpace([]byte(out)), &m); err != nil {
            t.Fatalf("invalid json output: %v, raw=%q", err, out)
        }
        if m["level"] != tc.level {
            t.Fatalf("level missing or incorrect: got %v want %v", m["level"], tc.level)
        }
        if m["msg"] != "hello" {
            t.Fatalf("fields not preserved: %v", m)
        }
    }
}

