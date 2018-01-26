package model

import (
	"time"
	"sync"
)

type Throttler struct {
	last map[string]time.Time
	mu   sync.Mutex
}

func (t *Throttler) Delay(host string) time.Duration {
	now := time.Now()
	t.mu.Lock()
	delay, ok := t.last[host]
	if !ok || delay.Before(now) {
		delay = now
	}
	t.last[host] = delay.Add(time.Second * 1)
	t.mu.Unlock()
	return delay.Sub(now)
}

func NewThrottler() *Throttler {
	return &Throttler{
		last: make(map[string]time.Time),
	}
}
