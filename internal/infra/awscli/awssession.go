package awscli

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func newSession(region, endpoint string) *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
	}))
}
