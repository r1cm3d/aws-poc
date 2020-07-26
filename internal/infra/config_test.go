package infra

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const line = `SQS_TEST_QUEUE='test-queue'`

func TestGetKey(t *testing.T) {
	act := getKey(line)

	assert.Equal(t, "SQS_TEST_QUEUE", act)
}

func TestGetValue(t *testing.T) {
	act := getValue(line)

	assert.Equal(t, "test-queue", act)
}
