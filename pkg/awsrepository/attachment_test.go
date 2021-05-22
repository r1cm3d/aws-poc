package awsrepository

import (
	"aws-poc/pkg/test/integration"
	"testing"
)

// TODO: implement it
func TestGet(t *testing.T) {
	integration.SkipShort(t)
	setupTable()
	defer cleanupTable()

}
