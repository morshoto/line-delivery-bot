package qr

import "strings"

// Parser extracts carrier and tracking information from a QR text payload.
type Parser interface {
    Parse(text string) (carrier, tracking string)
}

// NaiveParser implements a very simple heuristic parser.
type NaiveParser struct{}

func (NaiveParser) Parse(text string) (carrier, tracking string) {
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

