package infra

import (
	"github.com/stretchr/testify/assert"
	"reflect"
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

func TestMerge(t *testing.T) {
	m1 := map[string]string{
		"k1": "v1",
	}
	m2 := map[string]string{
		"k2": "v2",
	}
	exp := map[string]string{
		"k1": "v1",
		"k2": "v2",
	}

	act := merge(m1, m2)

	assert.True(t, reflect.DeepEqual(exp, act))
}

func TestLoadFileNotFound(t *testing.T) {
	if _, err := loadFile("unreachablePath"); err == nil {
		assert.Fail(t, "expected an error when file is not found")
	}
}

func TestLoadConfigNotFound(t *testing.T) {
	if _, err := loadConf("unreachableDirectory"); err == nil {
		assert.Fail(t, "expected an error when directory not exists")
	}
}
