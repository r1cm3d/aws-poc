package awssession

import (
	"aws-poc/pkg/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// NewSession returns a AWS session given a region and an endpoint
func NewSession(region, endpoint string) *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint),
	}))
}

// NewLocalSession returns a AWS session loading the default configurations according config.LoadDefaultConf() function
func NewLocalSession() (sess *session.Session) {
	env, _ := config.LoadDefaultConf()
	sess = NewSession(env["REGION"], env["ENDPOINT"])
	return
}

// NewSessionWithS3ForcePathStyle returns a AWS session given a region and an endpoint with S3ForcePathStyle true.
// This property is used to work with S3 storage locally
func NewSessionWithS3ForcePathStyle(region, endpoint string) (s *session.Session) {
	s = NewSession(region, endpoint)
	s.Config.S3ForcePathStyle = aws.Bool(true)
	return
}

// NewLocalSessionWithS3ForcePathStyle returns a AWS session loading the default configurations according config.LoadDefaultConf() function with S3ForcePathStyle true.
// This property is used to work with S3 storage locally
func NewLocalSessionWithS3ForcePathStyle() (sess *session.Session) {
	sess = NewLocalSession()
	sess.Config.S3ForcePathStyle = aws.Bool(true)
	return
}
