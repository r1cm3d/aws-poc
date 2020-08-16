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

type errMarshaller struct{}

func (m errMarshaller) MarshalMap(_ interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return nil, errors.New("error on marshalMap")
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

func TestPut_MarshalMapError(t *testing.T) {
	table := table{
		errMarshaller{},
	}

	err := table.put(nil)

	assert.NotNil(t, err)
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
