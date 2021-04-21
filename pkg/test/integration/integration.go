package integration

import "testing"

// SkipShort is used to in the beginning of the integration test.
// Unit tests are executed with -short flag. This is necessary to avoid execute
// tests that depends on external resources.
func SkipShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}
}
