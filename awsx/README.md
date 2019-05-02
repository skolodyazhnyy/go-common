# AWSx

AWSx package provides AWS service builder. It helps to avoid repeating writing boilerplate to construct AWS config and AWS session again and again. Instead it provides a handy configuration structure which can be embedded into your YAML configuration and then would allow you to build AWS services based on this configuration.

## Usage

You can populate `awesome.Config` from environment variables using `go-common/env`, or make it part of YAML configuration structure.

```go
type Config struct {
	AWS awsx.Config
}
```

Use configuration to construct awsx service:

```go
aws, err := awsx.New(config.AWS)
if err != nil {
	log.Fatal("An error occurred while initializing AWS session:", err)
}
```

Use awsx service to retrieve AWS session and AWS config specific to the service:

```go
s3Service := s3.New(aws.Session(), aws.Config("s3"))
sqsService := sqs.New(aws.Session(), aws.Config("sqs"))
// ...
```

Function `aws.Session()` will return AWS session object with global configuration, things like credentials, region etc. 
Function `aws.Config(string)` will return additional configuration specific to the service you are building, for example custom endpoint.