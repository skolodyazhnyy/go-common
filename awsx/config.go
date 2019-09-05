package awsx

// Config for Awesome AWS builder
// This config can be loaded from yaml or environment variables using go-common/env package
type Config struct {
	CredentialsType     string `yaml:"credentials_type" env:"AWS_CREDENTIALS_TYPE" config:"credentials_type"`
	CredentialsFilename string `yaml:"credentials_filename" env:"AWS_CREDENTIALS_FILENAME" config:"credentials_filename"`
	Profile             string `env:"AWS_PROFILE" config:"profile"`
	SecretKey           string `yaml:"secret_key" env:"AWS_SECRET_KEY" config:"secret_key"`
	AccessKey           string `yaml:"access_key" env:"AWS_ACCESS_KEY" config:"access_key"`
	Token               string `env:"AWS_TOKEN" config:"token"`
	Region              string `env:"AWS_REGION" config:"region"`
	Endpoint            string `env:"AWS_ENDPOINT" config:"endpoint"`

	// Services overrides configuration for a specific AWS service
	Services map[string]struct {
		Endpoint       *string `config:"endpoint"`
		Region         *string `config:"region"`
		ForcePathStyle *bool   `yaml:"force_path_style" config:"force_path_style"`
	}
}
