package qr

import "testing"

func TestNaiveParser_Parse_CarrierAndTracking(t *testing.T) {
    p := NaiveParser{}

    tests := []struct {
        name     string
        text     string
        carrier  string
        tracking string
    }{
        {"yamato_keyword", "Kuroneko Yamato: 1234567890", "yamato", "1234567890"},
        {"sagawa_keyword", "Sagawa 荷物 0987654321", "sagawa", "0987654321"},
        {"japanpost_keyword", "https://post.japanpost.jp/tracking 1112223334", "japanpost", "1112223334"},
        {"unknown", "Some QR with no known carrier 1234", "unknown", ""},
        {"longest_10_14_digits", "foo 123456789012 1234567890 bar 1234567", "unknown", "123456789012"},
        {"ignore_too_short_or_long", "a 123456789 a 123456789012345 b", "unknown", "123456789012"},
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            c, tr := p.Parse(tc.text)
            if c != tc.carrier {
                t.Fatalf("carrier: got %q want %q", c, tc.carrier)
            }
            if tr != tc.tracking {
                t.Fatalf("tracking: got %q want %q", tr, tc.tracking)
            }
        })
    }
}

