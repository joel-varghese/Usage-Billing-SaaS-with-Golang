package config

import (
    "context"
    "os"

    "github.com/aws/aws-sdk-go-v2/aws"
    awscfg "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
)

func LoadAWSConfig(endpoint, region string) aws.Config {
    opts := []func(*awscfg.LoadOptions) error{
		awscfg.WithRegion(region),
	}

    accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	sessionToken := os.Getenv("AWS_SESSION_TOKEN")

    // For LocalStack (when endpoint is provided), use test credentials if not set
    if endpoint != "" {
        if accessKey == "" {
            accessKey = "test"
        }
        if secretKey == "" {
            secretKey = "test"
        }
    }

    if accessKey != "" && secretKey != "" {
        opts = append(opts, awscfg.WithCredentialsProvider(
            credentials.NewStaticCredentialsProvider(accessKey, secretKey, sessionToken),
        ))
	}
    
    cfg, err := awscfg.LoadDefaultConfig(context.TODO(), opts...)
	if err != nil {
		panic(err)
	}

    // Override endpoint for LocalStack if needed
    if endpoint != "" {
        cfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(
            func(service, region string, options ...interface{}) (aws.Endpoint, error) {
                return aws.Endpoint{
                    URL:               endpoint,
                    SigningRegion:     region,
                    HostnameImmutable: true,
                }, nil
            })
    }

	return cfg
}
