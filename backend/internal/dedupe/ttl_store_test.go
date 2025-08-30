package dedupe

import (
    "testing"
    "time"
)

func TestTTLStore_Seen(t *testing.T) {
    s := NewTTLStore(50 * time.Millisecond)

    if seen := s.Seen("a"); seen {
        t.Fatalf("first time should not be seen")
    }
    if seen := s.Seen("a"); !seen {
        t.Fatalf("second time within TTL should be seen")
    }
}

func TestTTLStore_Expiry(t *testing.T) {
    s := NewTTLStore(10 * time.Millisecond)
    // first insert
    _ = s.Seen("x")
    // wait beyond TTL
    time.Sleep(20 * time.Millisecond)
    if seen := s.Seen("x"); seen {
        t.Fatalf("after TTL, should not be seen")
    }
}

