package awscli

import (
	"testing"
)

const (
	region   = "sa-east-1"
	endpoint = "http://localhost:1234"
)

func TestNewSession(t *testing.T) {
	exp := newSession(region, endpoint)

	if region != *exp.Config.Region {
		t.Errorf("region exp: %v, got: %v", region, exp.Config.Region)
	}
	if endpoint != *exp.Config.Endpoint {
		t.Errorf("endpoint exp: %v, got: %v", endpoint, exp.Config.Endpoint)
	}
}

func TestNewSessionWithS3ForcePathStyle(t *testing.T) {
	exp := newSessionWithS3ForcePathStyle(region, endpoint)

	if region != *exp.Config.Region {
		t.Errorf("region exp: %v, got: %v", region, *exp.Config.Region)
	}
	if endpoint != *exp.Config.Endpoint {
		t.Errorf("endpoint exp: %v, got: %v", endpoint, *exp.Config.Endpoint)
	}
	if !*exp.Config.S3ForcePathStyle {
		t.Errorf("S3ForcePathStyle exp: true, got: %v", *exp.Config.S3ForcePathStyle)
	}
}
