package awsx

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

// Defined errors
var (
	ErrInvalidCredentialsType = errors.New("credentials type is invalid (should be: shared, static or environment)")
)

// List of crendential types
const (
	// Shard credentials load ~/.aws/credentials file
	CredentialsTypeShared string = "shared"

	// Static credentials use Access Key and Secret Key from configuration
	CredentialsTypeStatic string = "static"

	// Environment credentials use environment variables
	CredentialsTypeEnvironment string = "environment"

	// Gets credentials from IAM Role assigned to EC2 instance
	CredentialsTypeRole string = "role"
)

// Credentials from AWS Config
func Credentials(c Config) (*credentials.Credentials, error) {
	switch c.CredentialsType {
	case CredentialsTypeStatic:
		return credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, c.Token), nil
	case CredentialsTypeEnvironment:
		return credentials.NewEnvCredentials(), nil
	case CredentialsTypeShared:
		return credentials.NewSharedCredentials(c.CredentialsFilename, c.Profile), nil
	case CredentialsTypeRole, "":
		return nil, nil
	}

	return nil, ErrInvalidCredentialsType
}
