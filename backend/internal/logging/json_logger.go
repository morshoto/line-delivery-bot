package logging

import (
    "encoding/json"
    "log"
)

// Logger is a minimal structured logger interface.
type Logger interface {
    Info(fields map[string]any)
    Warn(fields map[string]any)
    Error(fields map[string]any)
}

type JSONLogger struct{}

func (JSONLogger) logWith(level string, fields map[string]any) {
    if fields == nil {
        fields = map[string]any{}
    }
    fields["level"] = level
    b, _ := json.Marshal(fields)
    log.Println(string(b))
}

func (l JSONLogger) Info(fields map[string]any)  { l.logWith("info", fields) }
func (l JSONLogger) Warn(fields map[string]any)  { l.logWith("warn", fields) }
func (l JSONLogger) Error(fields map[string]any) { l.logWith("error", fields) }

