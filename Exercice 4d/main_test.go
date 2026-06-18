package main

import (
	"testing"
	"time"
)

func TestRandomDurationBounds(t *testing.T) {
	for i := 0; i < 100; i++ {
		duration := randomDuration(1, 3)
		if duration < time.Second || duration > 3*time.Second {
			t.Fatalf("randomDuration(1, 3) = %s, want between 1s and 3s", duration)
		}
	}
}

func TestStopAfterClosesQuitChannel(t *testing.T) {
	quitChannel := make(chan struct{})
	go stopAfter(10*time.Millisecond, quitChannel)

	select {
	case <-quitChannel:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("quitChannel was not closed")
	}
}
