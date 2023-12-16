package credstash

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/kms"

	"github.com/iyuuya/go-credstash/internal/credstash/app"
	"github.com/iyuuya/go-credstash/item"
)

type Credstash struct {
	app *app.App
}

func NewCredstash(d *dynamodb.Client, k *kms.Client) *Credstash {
	return &Credstash{app: app.NewAppWithClients(d, k)}
}

func (c *Credstash) List() ([]*item.Item, error) {
	return c.app.List()
}

func (c *Credstash) Get(name string, ctx map[string]string, version string) (string, error) {
	return c.app.Get(name, ctx, version)
}

func (c *Credstash) Put(name, value, kmsKeyId string, ctx map[string]string) error {
	return c.app.Put(name, value, kmsKeyId, ctx)
}

func (c *Credstash) Delete(name string, version string) error {
	return c.app.Delete(name, version)
}
