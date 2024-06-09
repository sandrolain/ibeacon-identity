package main

import (
	"fmt"
	"time"
)

var cache = make(map[string]time.Time)

func isNewIBeacon(ib *IBeacon, to time.Duration) bool {
	key := fmt.Sprintf("%x", ib.UMMID)
	now := time.Now()
	if t, ok := cache[key]; ok {
		if now.Sub(t) < to {
			return false
		}
	}
	cache[key] = now
	return true
}
