package main

import (
	"testing"
	"time"
)

func TestMainFinishesUnderThreeSeconds(t *testing.T) {
	done := make(chan bool)

	go func() {
		main()
		close(done)
	}()

	select {
	case <-done:
		// main() finished within 3 seconds
	case <-time.After(2 * time.Second):
		t.Error("main() did not finish within 3 seconds")
	}
}
