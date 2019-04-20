package awsx

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Session from awsx.Config
func Session(c Config) (*session.Session, error) {
	credentials, err := Credentials(c)
	if err != nil {
		return nil, err
	}

	config := aws.NewConfig().WithCredentials(credentials)

	if c.Endpoint != "" {
		config = config.WithEndpoint(c.Endpoint)
	}

	if c.Region != "" {
		config = config.WithRegion(c.Region)
	}

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return sess, nil
}
