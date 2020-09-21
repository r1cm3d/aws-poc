package integration

import "testing"

func SkipShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
}
