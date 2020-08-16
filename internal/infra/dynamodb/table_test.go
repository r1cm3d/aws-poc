package dynamodb

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"testing"
	"time"
)

func TestPut(t *testing.T) {
	skipShort(t)
	setup()
	defer teardown()

	i := Item{
		DisputeId: 666,
		Timestamp: "2020-04-17T17:19:19.831Z",
	}

	put(i)
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

	time.Sleep(5 * time.Second)

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
