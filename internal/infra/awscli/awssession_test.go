package awscli

import (
	"aws-poc/internal/infra"
	"testing"
)

func TestNewSession(t *testing.T) {
	env, _ := infra.LoadDefaultConf()
	region, endpoint := env["REGION"], env["ENDPOINT"]

	exp := newSession(region, endpoint)

	if region != *exp.Config.Region {
		t.Errorf("region exp: %v, got: %v", region, exp.Config.Region)
	}
	if endpoint != *exp.Config.Endpoint {
		t.Errorf("endpoint exp: %v, got: %v", endpoint, exp.Config.Endpoint)
	}
}
