package dynamodb

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

type Item struct {
	DisputeId int
	Timestamp string
}

type errMarshal struct{}
type errPutItem struct{}

func (m errMarshal) marshalMap(_ interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return nil, errors.New("error on marshalMap")
}

func (m errPutItem) marshalMap(_ interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return nil, nil
}

func (m errPutItem) putItem(_ *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return nil, errors.New("error on put item")
}

func TestPutIntegration(t *testing.T) {
	skipShort(t)
	setup()
	defer teardown()
	i := Item{
		DisputeId: 666,
		Timestamp: "2020-04-17T17:19:19.831Z",
	}
	table := newTable()

	err := table.put(i)
	assert.Nil(t, err)
}

func TestPut_Error(t *testing.T) {
	tables := [2]table{{
		marshaller: errMarshal{},
	}, {
		marshaller: errPutItem{},
		persistent: errPutItem{},
	}}

	for _, ta := range tables {
		err := ta.put(nil)

		assert.NotNil(t, err)
	}
}

func setup() {
	svc := svc()

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: disputeId,
				AttributeType: numberType,
			},
			{
				AttributeName: timestamp,
				AttributeType: stringType,
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: disputeId,
				KeyType:       hashKeyType,
			},
			{
				AttributeName: timestamp,
				KeyType:       rangeKeyType,
			},
		},
		BillingMode: payPerRequest,
		TableName:   tableName,
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("created the table", tableName)
}

func teardown() {
	svc := svc()

	input := &dynamodb.DeleteTableInput{
		TableName: tableName,
	}
	_, err := svc.DeleteTable(input)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("deleted the table", tableName)
}

func skipShort(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
}
