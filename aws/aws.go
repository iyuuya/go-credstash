package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

func NewDynamoDBClient(endpoint string) (*dynamodb.Client, error) {
	cfg, err := newConfig(endpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	return dynamodb.NewFromConfig(cfg), nil
}

func NewKMSClient(endpoint string) (*kms.Client, error) {
	cfg, err := newConfig(endpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	return kms.NewFromConfig(cfg), nil
}

func newConfig(endpoint string) (aws.Config, error) {
	if endpoint == "" {
		return config.LoadDefaultConfig(context.TODO())
	}

	return config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: endpoint}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)
}
