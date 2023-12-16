package item

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDB struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDB(c *dynamodb.Client) *DynamoDB {
	return &DynamoDB{c, "credential-store"}
}

func (db *DynamoDB) Get(name, version string) (*Item, error) {
	items, err := db.Select(name, 1, version)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		if version == "" {
			return nil, fmt.Errorf("%s is not found", name)
		} else {
			return nil, fmt.Errorf("%s --version: %s is not found", name, version)
		}
	}

	return items[0], nil
}

func (db *DynamoDB) Select(name string, limit int, version string) ([]*Item, error) {
	var l *int
	var v *string

	if limit > 0 {
		l = &limit
	}
	if version != "" {
		v = &version
	}
	params := db.buildParams(name, nil, l, v)

	out, err := db.client.Query(context.TODO(), params)
	if err != nil {
		return nil, err
	}
	items := out.Items

	res := make([]*Item, len(items))

	for i, item := range items {
		sv := func(k string, i map[string]types.AttributeValue) *string {
			switch i[k].(type) {
			case *types.AttributeValueMemberS:
				return &i[k].(*types.AttributeValueMemberS).Value
			case *types.AttributeValueMemberNULL:
				return nil
			default:
				return nil
			}
		}

		key := sv("key", item)
		if key == nil {
			k := ""
			key = &k
		}

		res[i] = &Item{
			key:      *key,
			contents: sv("contents", item),
			name:     sv("name", item),
			version:  sv("version", item),
		}
	}
	return res, nil
}

func (db *DynamoDB) Put(i *Item) error {
	sv := func(v *string) types.AttributeValue {
		if v == nil {
			return &types.AttributeValueMemberNULL{}
		} else {
			return &types.AttributeValueMemberS{Value: *v}
		}
	}

	out, err := db.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName:                aws.String(db.tableName),
		ConditionExpression:      aws.String("attribute_not_exists(#name)"),
		ExpressionAttributeNames: map[string]string{"#name": "name"},
		Item: map[string]types.AttributeValue{
			"name":     sv(i.GetName()),
			"version":  sv(i.GetVersion()),
			"key":      &types.AttributeValueMemberS{Value: i.GetKey()},
			"contents": sv(i.GetContents()),
			"hmac":     sv(i.GetHMAC()),
		},
	})
	if err != nil {
		return err
	}
	log.Printf("%#v", *out)
	return nil
}

func (db *DynamoDB) List() ([]*Item, error) {
	items, err := db.fetchAllItems()
	if err != nil {
		log.Fatalf("failed to scan, %v", err)
	}

	res := make([]*Item, len(items))

	for i, item := range items {
		sv := func(k string, i map[string]types.AttributeValue) *string {
			switch i[k].(type) {
			case *types.AttributeValueMemberS:
				return &i[k].(*types.AttributeValueMemberS).Value
			case *types.AttributeValueMemberNULL:
				return nil
			default:
				return nil
			}
		}

		key := sv("key", item)
		if key == nil {
			k := ""
			key = &k
		}

		res[i] = &Item{
			key:      *key,
			contents: sv("contents", item),
			name:     sv("name", item),
			version:  sv("version", item),
		}
	}
	return res, nil
}

func (db *DynamoDB) Delete(item *Item) error {
	_, err := db.client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(db.tableName),
		Key: map[string]types.AttributeValue{
			"name":    &types.AttributeValueMemberS{Value: *item.GetName()},
			"version": &types.AttributeValueMemberS{Value: *item.GetVersion()},
		},
	})
	return err
}

func (db *DynamoDB) fetchAllItems() ([]map[string]types.AttributeValue, error) {
	items := make([]map[string]types.AttributeValue, 0)
	var lastKey map[string]types.AttributeValue

	for {
		out, err := db.client.Scan(
			context.TODO(),
			&dynamodb.ScanInput{
				TableName:            aws.String(db.tableName),
				ProjectionExpression: aws.String("#name, version"),
				ExpressionAttributeNames: map[string]string{
					"#name": "name",
				},
				ExclusiveStartKey: lastKey,
			},
		)
		if err != nil {
			return nil, err
		}
		for _, i := range out.Items {
			items = append(items, i)
		}
		lastKey = out.LastEvaluatedKey
		if lastKey == nil {
			break
		}
	}
	return items, nil
}

func (db *DynamoDB) Setup() {
	db.client.CreateTable(
		context.TODO(),
		&dynamodb.CreateTableInput{
			TableName: aws.String(db.tableName),
			KeySchema: []types.KeySchemaElement{
				{AttributeName: aws.String("name"), KeyType: types.KeyTypeHash},
				{AttributeName: aws.String("version"), KeyType: types.KeyTypeRange},
			},
			AttributeDefinitions: []types.AttributeDefinition{
				{AttributeName: aws.String("name"), AttributeType: types.ScalarAttributeTypeS},
				{AttributeName: aws.String("version"), AttributeType: types.ScalarAttributeTypeS},
			},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(1),
				WriteCapacityUnits: aws.Int64(1),
			},
		},
	)
}

func (db *DynamoDB) buildParams(name string, pluck *string, limit *int, version *string) *dynamodb.QueryInput {
	params := &dynamodb.QueryInput{
		TableName:                aws.String(db.tableName),
		ConsistentRead:           aws.Bool(true),
		KeyConditionExpression:   aws.String("#name = :name"),
		ExpressionAttributeNames: map[string]string{"#name": "name"},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name": &types.AttributeValueMemberS{Value: name},
		},
	}
	if pluck != nil {
		params.ProjectionExpression = pluck
	}
	if limit != nil {
		params.Limit = aws.Int32(int32(*limit))
		params.ScanIndexForward = aws.Bool(false)
	}
	if version != nil {
		params.KeyConditionExpression = aws.String("#name = :name AND #version = :version")
		params.ExpressionAttributeNames["#version"] = "version"
		params.ExpressionAttributeValues[":version"] = &types.AttributeValueMemberS{Value: *version}
	}

	return params
}
