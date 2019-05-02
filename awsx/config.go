package awsx

// Config for Awesome AWS builder
// This config can be loaded from yaml or environment variables using go-common/env package
type Config struct {
	CredentialsType     string `yaml:"credentials_type" env:"AWS_CREDENTIALS_TYPE"`
	CredentialsFilename string `yaml:"credentials_filename" env:"AWS_CREDENTIALS_FILENAME"`
	Profile             string `env:"AWS_PROFILE"`
	SecretKey           string `yaml:"secret_key" env:"AWS_SECRET_KEY"`
	AccessKey           string `yaml:"access_key" env:"AWS_ACCESS_KEY"`
	Token               string `env:"AWS_TOKEN"`
	Region              string `env:"AWS_REGION"`
	Endpoint            string `env:"AWS_ENDPOINT"`

	// Services overrides configuration for a specific AWS service
	Services map[string]struct {
		Endpoint       *string
		Region         *string
		ForcePathStyle *bool `yaml:"force_path_style"`
	}
}
