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
    // Extract candidates: any 10-14 digit substring. Prefer 12-digit if present,
    // otherwise prefer longer lengths, then first occurrence.
    var runs []string
    cur := make([]rune, 0, len(text))
    flush := func() {
        if len(cur) > 0 {
            runs = append(runs, string(cur))
            cur = cur[:0]
        }
    }
    for _, r := range text {
        if r >= '0' && r <= '9' {
            cur = append(cur, r)
        } else {
            flush()
        }
    }
    flush()

    // collect substrings of length 10-14 from each run (sliding window)
    type candidate struct{
        val string
        length int
        index int // order encountered
    }
    cands := make([]candidate, 0)
    idx := 0
    for _, run := range runs {
        n := len(run)
        if n == 0 {
            continue
        }
        for l := 10; l <= 14; l++ {
            if n < l {
                continue
            }
            for i := 0; i+l <= n; i++ {
                cands = append(cands, candidate{val: run[i : i+l], length: l, index: idx})
                idx++
            }
        }
    }
    // selection preference: 12, then 14, 13, 11, 10. If multiple, earliest.
    pref := []int{12, 14, 13, 11, 10}
    best := ""
    bestIdx := int(^uint(0) >> 1) // max int
    for _, p := range pref {
        for _, c := range cands {
            if c.length == p {
                if best == "" || c.index < bestIdx {
                    best = c.val
                    bestIdx = c.index
                }
            }
        }
        if best != "" {
            break
        }
    }
    tracking = best
    return
}
