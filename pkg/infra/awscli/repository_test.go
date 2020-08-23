package awscli

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"testing"
)

type Item struct {
	DisputeId int
	Timestamp string
}

func TestPutIntegration(t *testing.T) {
	skipShort(t)
	setupTable()
	defer cleanupTable()
	i := Item{
		DisputeId: 666,
		Timestamp: "2020-04-17T17:19:19.831Z",
	}
	table := newRepository(newLocalSession())

	if err := table.put(i); err != nil {
		t.Errorf("put fails: %d", err)
	}
}

func setupTable() {
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

	r := newRepository(newLocalSession())
	svc := r.svc()
	if _, err := svc.CreateTable(input); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("created the repository", tableName)
}

func cleanupTable() {
	r := newRepository(newLocalSession())
	svc := r.svc()

	input := &dynamodb.DeleteTableInput{
		TableName: tableName,
	}
	if _, err := svc.DeleteTable(input); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("deleted the repository", tableName)
}
