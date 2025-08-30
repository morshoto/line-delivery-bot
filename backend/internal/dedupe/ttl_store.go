package dedupe

import (
    "sync"
    "time"
)

// Store defines an interface for duplicate detection with TTL.
type Store interface {
    // Seen returns true if the key has been seen within TTL and refreshes its timestamp.
    Seen(key string) bool
}

// TTLStore is a threadsafe in-memory TTL-based dedupe store.
type TTLStore struct {
    mu   sync.Mutex
    data map[string]time.Time
    ttl  time.Duration
}

func NewTTLStore(ttl time.Duration) *TTLStore {
    return &TTLStore{data: make(map[string]time.Time), ttl: ttl}
}

// Seen returns whether the key was seen within TTL and refreshes the timestamp.
func (s *TTLStore) Seen(key string) bool {
    s.mu.Lock()
    defer s.mu.Unlock()

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

