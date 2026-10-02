package state

import (
	"testing"
	"time"
)

// waitFor polls condition until it returns true. If timeout elapses first, it fails the
// test with msg, so the failure is more informative than a generic timeout notice.
func waitFor(t *testing.T, timeout time.Duration, condition func() bool, msg string) {
	t.Helper()

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(time.Millisecond)
	}

	t.Fatal(msg)
}
