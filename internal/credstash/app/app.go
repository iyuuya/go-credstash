package app

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/kms"

	"github.com/iyuuya/go-credstash/aws"
	"github.com/iyuuya/go-credstash/item"
	"github.com/iyuuya/go-credstash/secret"
)

type App struct {
	DynamoDB *item.DynamoDB
	KMS      *kms.Client
}

func NewApp(endpoint string) (*App, error) {
	d, err := aws.NewDynamoDBClient(endpoint)
	if err != nil {
		return nil, err
	}
	k, err := aws.NewKMSClient(endpoint)
	if err != nil {
		return nil, err
	}

	return &App{item.NewDynamoDB(d), k}, nil
}

func NewAppWithClients(d *dynamodb.Client, k *kms.Client) *App {
	return &App{DynamoDB: item.NewDynamoDB(d), KMS: k}
}

func (a *App) List() ([]*item.Item, error) {
	return a.DynamoDB.List()
}

func (a *App) Get(name string, ctx map[string]string, version string) (string, error) {
	sec, err := secret.Find(a.DynamoDB, a.KMS, name, ctx, version)
	if err != nil {
		return "", err
	}
	if sec.IsFalsified() {
		return "", fmt.Errorf("Invalid secret. %s has falsified", name)
	}
	return sec.DecryptedValue()
}

func (a *App) Put(name, value, kmsKeyId string, ctx map[string]string) error {
	var id *string
	if kmsKeyId != "" {
		id = &kmsKeyId
	}

	s := secret.NewSecret(name, value, nil, nil, nil, ctx)
	if err := s.Encrypt(a.KMS, id, ctx); err != nil {
		return err
	}
	return s.Save(a.DynamoDB)
}

func (a *App) Delete(name string, version string) error {
	items, err := a.DynamoDB.Select(name, -1, version)
	if err != nil {
		return err
	}
	if err != nil || len(items) == 0 {
		return fmt.Errorf("Item not found: %s", name)
	}

	return a.DynamoDB.Delete(items[0])
}
