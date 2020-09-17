package awscli

import (
	"aws-poc/pkg/config"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

const (
	region   = "sa-east-1"
	endpoint = "http://localhost:1234"
)

func newLocalSession() (sess *session.Session) {
	env, _ := config.LoadDefaultConf()
	sess = newSession(env["REGION"], env["ENDPOINT"])
	return
}

func newLocalSessionWithS3ForcePathStyle() (sess *session.Session) {
	sess = newLocalSession()
	sess.Config.S3ForcePathStyle = aws.Bool(true)
	return
}

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
