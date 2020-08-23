package infra

import (
	"reflect"
	"testing"
)

const line = `SQS_TEST_QUEUE='test-queue'`

func TestGetKey(t *testing.T) {
	if got := getKey(line); got != "SQS_TEST_QUEUE" {
		t.Errorf("want SQS_TEST_QUEUE; got: %s", got)
	}
}

func TestGetValue(t *testing.T) {
	if got := getValue(line); got != "test-queue" {
		t.Errorf("want test-queue; got: %s", got)
	}
}

func TestMerge(t *testing.T) {
	m1 := map[string]string{
		"k1": "v1",
	}
	m2 := map[string]string{
		"k2": "v2",
	}
	want := map[string]string{
		"k1": "v1",
		"k2": "v2",
	}

	if got := merge(m1, m2); reflect.DeepEqual(got, want) {
		t.Errorf("want %v; got: %v", want, got)
	}
}

func TestLoadFileNotFound(t *testing.T) {
	if _, err := loadFile("unreachablePath"); err == nil {
		t.Error("expected an error when file is not found")
	}
}

func TestLoadConfigNotFound(t *testing.T) {
	if _, err := LoadConf("unreachableDirectory"); err == nil {
		t.Error("expected an error when directory not exists")
	}
}
