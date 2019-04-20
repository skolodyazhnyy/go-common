package awsx

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// AWS service builder
type AWS struct {
	session *session.Session
	config  Config
}

// New AWS service builder
func New(c Config) (*AWS, error) {
	sess, err := Session(c)
	if err != nil {
		return nil, err
	}

	return &AWS{config: c, session: sess}, nil
}

// Session for AWS
// Session carries global configuration for all services, things like credentials, AWS region etc
func (a *AWS) Session() *session.Session {
	return a.session
}

// Config for specific service
// This method returns configuration specific to particular AWS service. Name of the service (service argument) should
// repeat golang package name of the service, for example: s3, sqs, kinesis, dynamodb etc.
func (a *AWS) Config(service string) *aws.Config {
	config := aws.NewConfig()

	specific, ok := a.config.Services[service]
	if !ok {
		return config
	}

	if specific.Endpoint != nil {
		config = config.WithEndpoint(aws.StringValue(specific.Endpoint))
	}

	if specific.ForcePathStyle != nil {
		config = config.WithS3ForcePathStyle(aws.BoolValue(specific.ForcePathStyle))
	}

	if specific.Region != nil {
		config = config.WithRegion(aws.StringValue(specific.Region))
	}

	return config
}
